package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	iamqp "github.com/Blockpour/Blockpour-Geth-Indexer/libs/amqp"
	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	oldpriceresolver "github.com/Blockpour/Blockpour-Geth-Indexer/libs/oldpricing"
	priceresolver "github.com/Blockpour/Blockpour-Geth-Indexer/libs/pricing"
	uniswapv2 "github.com/Blockpour/Blockpour-Geth-Indexer/libs/processors/uniswapV2"
	uniswapv3 "github.com/Blockpour/Blockpour-Geth-Indexer/libs/processors/uniswapV3"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/ethrpc"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/instrumentation"
	lb "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	outs "github.com/Blockpour/Blockpour-Geth-Indexer/services/output_sink"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/version"
	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type NodeImpl struct {
	service.BaseService

	log          logger.Logger
	EthRPC       ethrpc.EthRPC   // HA upstream connection to rpc nodes, uses mspool
	LocalBackend lb.LocalBackend // Local database for caching / processing
	OutputSink   outs.OutputSink // Consumer for offloading processed data

	// Configs
	startBlock                      uint64   // User defined startBlock
	skipResumeRemote                bool     // skip checking remote for resume height
	skipResumeLocal                 bool     // skip checking localbackend for resume height
	remoteResumeURL                 string   // URL to use for resume height GET request
	prodcheck                       bool     // checks for prod grade settings
	eventsToIndex                   []string // user requested events to index in string form
	maxCPUParallels                 int      // user requested CPU threads to allocate to the process
	maxBlockSpanPerCall             uint64   // max block spans to log per initial filtering call
	pricingChainlinkOraclesDumpFile string   // user provided chainlink oracles to trust
	pricingDexDumpFile              string   // user provided dexes for faster catchup

	// Internal Data Structures
	moniker            string                                // user defined moniker for this node
	network            string                                // user defined evm compatible network name
	nodeID             uuid.UUID                             // system generated node identifier unique for each run
	mergedTopics       map[common.Hash]itypes.ProcessingType // information on topics to index
	mergedTopicsKeys   []common.Hash                         // cached keys of mergedTopics
	indexedHeight      uint64
	currentHeight      uint64
	allowPricingState  bool
	oldPricerOracleMap string
	quitCh             chan struct{}

	// Backoff configuration
	backoff *backoff.ConstantBackOff

	// Library instances
	procUniV2 uniswapv2.UniswapV2Processor
	procUniV3 uniswapv3.UniswapV3Processor
	pricer    *priceresolver.Engine
	oldpricer *oldpriceresolver.Pricing
}

// OnStart starts the Node. It implements service.Service.
func (n *NodeImpl) OnStart(ctx context.Context) error {
	if int(n.maxCPUParallels) > runtime.NumCPU() {
		n.log.Warn("running on fewer threads than requested parallels",
			"parallels", runtime.NumCPU(),
			"requested", n.maxCPUParallels)
		n.maxCPUParallels = runtime.NumCPU()
	}

	runtime.GOMAXPROCS(n.maxCPUParallels)
	n.log.Info("set runtime max parallelism",
		"parallels", n.maxCPUParallels)

	n.configureBackoff()

	n.nodeID = uuid.New()

	n.log.Info("node identifier generated", "moniker", n.moniker, "nodeID", n.nodeID)

	if err := n.EthRPC.Start(ctx); err != nil {
		return err
	}

	if err := n.LocalBackend.Start(ctx); err != nil {
		return err
	}

	if err := n.OutputSink.Start(ctx); err != nil {
		n.log.Info("Error initializing output sink, will reattempt connection until ready")
	}

	// Setup what to index
	requestedEvents, err := util.ConstructTopics(n.eventsToIndex)
	if err != nil {
		return err
	}
	// required topics by the pricing engine
	var extraRequiredEvents []common.Hash
	if n.allowPricingState {
		extraRequiredEvents = []common.Hash{itypes.UniV2MintTopic,
			itypes.UniV2BurnTopic,
			itypes.UniV2SwapTopic}
	}

	n.mergedTopics = mergeTopics(requestedEvents, extraRequiredEvents)
	keys := make([]common.Hash, len(n.mergedTopics))

	i := 0
	for val, ptype := range n.mergedTopics {
		// Display
		str, _ := itypes.GetStringForTopic(val)
		vh := val.Hex()
		fingerPrint := vh[:7] + ".." + vh[len(vh)-3:]
		reason := "pricing"
		if ptype == itypes.UserRequested {
			reason = "indexing"
		}
		n.log.Info(fmt.Sprintf("enabled %s(%s) for %s", str, fingerPrint, reason), "ptype", ptype)

		// Set val
		keys[i] = val
		i++
	}
	n.mergedTopicsKeys = keys

	// Setup processors
	n.procUniV2 = uniswapv2.UniswapV2Processor{n.mergedTopics, n.EthRPC}
	n.procUniV3 = uniswapv3.UniswapV3Processor{n.mergedTopics, n.EthRPC}
	if n.allowPricingState {
		n.pricer = priceresolver.NewDefaultEngine(n.log.With("module", "pricing"),
			n.pricingChainlinkOraclesDumpFile,
			n.pricingDexDumpFile,
			n.EthRPC,
			n.LocalBackend)
	} else {
		n.oldpricer = oldpriceresolver.GetPricingEngine(n.oldPricerOracleMap, n.EthRPC)
	}

	// TODO: Do height syncup using both LocalBackend and remote http
	// startHeight, err := n.getResumeHeight()
	n.indexedHeight = n.syncStartHeight()

	// Loop for impl
	go n.loop()

	return nil
}

// OnStop stops the Node. It implements service.Service
func (n *NodeImpl) OnStop() {
	n.quitCh <- struct{}{}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		n.LocalBackend.Stop()
		n.LocalBackend.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		n.OutputSink.Stop()
		n.OutputSink.Stop()
	}()

	wg.Wait()
}

// Loop implements core indexing logic
func (n *NodeImpl) loop() {
	for {
		select {
		case <-time.After(time.Second * 2):
			// Loop in case we are lagging, so we dont wait 3 secs between epochs
			for {
				height, err := n.EthRPC.GetCurrentBlockHeight()

				if err != nil {
					n.log.Warn(fmt.Sprintf("Error retrieving block height, retrying. Caused by: %s", err))
					break
				}
				n.currentHeight = height

				if n.currentHeight < n.indexedHeight {
					n.log.Warn(fmt.Sprintf("rpc height (%d) is less than indexed height (%d), possible n/w reorg or p2p failure",
						n.currentHeight, n.indexedHeight))
				}
				if n.currentHeight == n.indexedHeight {
					continue
				}
				endingBlock := n.currentHeight
				isOnHead := true
				if (endingBlock - n.indexedHeight) > n.maxBlockSpanPerCall {
					isOnHead = false
					endingBlock = n.indexedHeight + n.maxBlockSpanPerCall
				}

				n.log.Info(fmt.Sprintf("chainhead: %d (+%d away), indexed: %d",
					n.currentHeight, n.currentHeight-n.indexedHeight, n.indexedHeight))

				instrumentation.CurrentBlock.Set(float64(n.currentHeight))

				logs, err := n.EthRPC.GetFilteredLogs(ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(n.indexedHeight + 1)),
					ToBlock:   big.NewInt(int64(endingBlock)),
					Topics:    [][]common.Hash{n.mergedTopicsKeys},
				})

				if err != nil {
					n.log.Error("encountered error", "error", err)
					continue
				}

				n.processBatchedBlockLogs(logs, n.indexedHeight+1, endingBlock)

				n.indexedHeight = endingBlock
				instrumentation.ProcessedBlock.Set(float64(n.currentHeight))

				if isOnHead {
					break
				}
			}
		case <-n.quitCh:
			n.log.Info("quitting realtime indexer")
		}
	}
}

func (n *NodeImpl) processBatchedBlockLogs(logs []types.Log, start uint64, end uint64) {
	// Assuming for any height H, either we will have all the concerned logs
	// or not even one
	kv := GroupByBlockNumber(logs)

	for block := start; block <= end; block++ {
		backoff.Retry(func() error { return n.processBlock(kv, block) }, n.backoff)
	}
}

func (n *NodeImpl) processBlock(kv map[uint64]CLogType, block uint64) error {
	n.log.Info(fmt.Sprintf("processing block %d", block))
	startTime := time.Now()
	_time, err := n.EthRPC.GetBlockTimestamp(block)
	if err != nil {
		n.log.Warn(fmt.Sprintf("Error retrieving timestamp for block %d. Caused by: %s", block, err))
		return err
	}

	logs := kv[block]
	blockSynopis := itypes.BlockSynopsis{
		Height:        block,
		BlockTime:     _time,
		EventsScanned: uint64(logs.Len()),
	}

	var wg sync.WaitGroup

	var processedItems []interface{} = make([]interface{}, len(logs))
	for idx, _log := range logs {
		wg.Add(1)
		go func(_log types.Log, idx int) {
			n.decodeLog(_log, processedItems, idx, blockSynopis.BlockTime)
			wg.Done()
		}(_log, idx)
	}

	wg.Wait()

	// if err != nil {
	// 	n.log.Debug(fmt.Sprintf("Error processing block %d. Retrying. Error caused by: %s", block, err))
	// 	return err
	// }

	processingTime := time.Now()

	// Run processedItems through pricing engine
	var newDexes []itypes.UniV2Metadata
	if n.allowPricingState {
		n.log.Info("allowed pricing state")
		newDexes, err = backoff.RetryWithData(
			func() ([]itypes.UniV2Metadata, error) {
				newDexes, err := n.pricer.Resolve(block, processedItems)
				if err != nil {
					if errors.Is(err, priceresolver.ErrorRequestedResolutionPresent) {
						n.log.Debugf("%s", err)
						return newDexes, nil
					}
					n.log.Infof("Error resolving dex. Caused by: %s\n", err)
				}
				return newDexes, err
			}, n.backoff)
		if err != nil {
			return err
		}
	} else {
		err = n.oldpricer.Resolve(block, processedItems)

		if err != nil {
			return err
		}
	}

	pricingTime := time.Now()

	// Package processedItems into payload for output
	populateBlockSynopsis(&blockSynopis, processedItems, startTime, processingTime, pricingTime)
	payload := n.genPayload(&blockSynopis, processedItems, newDexes)
	payload.allowPricingState = n.allowPricingState
	n.log.Debug("Sending data to output sink")
	for {
		err = n.OutputSink.Send(payload)
		if err == nil {
			break
		}
		n.log.Warn("Error sending message to output sink: " + fmt.Sprint(err))
		time.Sleep(2 * time.Second)
	}

	n.log.Debug("Syncing local backend")
	// Sync localBackend states
	backoff.Retry(func() error { return n.LocalBackend.Sync() }, n.backoff)
	return nil
}

func (n *NodeImpl) decodeLog(l types.Log,
	items []interface{},
	idx int,
	blockTime uint64,
) error {

	primaryTopic := l.Topics[0]
	switch primaryTopic {
	// ---- Uniswap V2 ----
	case itypes.UniV2MintTopic:
		instrumentation.MintV2Found.Inc()
		return n.procUniV2.ProcessUniV2Mint(l, items, idx, blockTime)
	case itypes.UniV2BurnTopic:
		instrumentation.BurnV2Found.Inc()
		return n.procUniV2.ProcessUniV2Burn(l, items, idx, blockTime)
	case itypes.UniV2SwapTopic:
		instrumentation.SwapV2Found.Inc()
		return n.procUniV2.ProcessUniV2Swap(l, items, idx, blockTime)

	// // ---- Uniswap V3 ----
	case itypes.UniV3MintTopic:
		instrumentation.MintV3Found.Inc()
		return n.procUniV3.ProcessUniV3Mint(l, items, idx, blockTime)
	case itypes.UniV3BurnTopic:
		instrumentation.BurnV3Found.Inc()
		return n.procUniV3.ProcessUniV3Burn(l, items, idx, blockTime)
	case itypes.UniV3SwapTopic:
		instrumentation.SwapV3Found.Inc()
		return n.procUniV3.ProcessUniV3Swap(l, items, idx, blockTime)

		// // ---- ERC 20 ----
		// case itypes.ERC20TransferTopic:
		// 	// instrumentation.TfrFound.Inc()
		// 	n.processERC20Transfer(l, items, bm, mt)
	}
	return nil
}

func populateBlockSynopsis(bs *itypes.BlockSynopsis,
	items []interface{},
	startTime time.Time,
	processingTime time.Time,
	pricingTime time.Time) {
	distribution := make(map[string]uint64, len(items))
	defaultKey := ""
	for _, item := range items {
		if item == nil {
			continue
		}

		itemKey := defaultKey
		isPricedCorrectly := false

		switch i := item.(type) {
		case *itypes.Mint:
			itemKey = fmt.Sprintf("(%v, %v)", i.Type, i.ProcessingType.ToString())
			isPricedCorrectly = i.Price0 != nil && i.Price1 != nil && i.Amount0 != nil && i.Amount1 != nil && i.AmountUSD != nil
		case *itypes.Burn:
			itemKey = fmt.Sprintf("(%v, %v)", i.Type, i.ProcessingType.ToString())
			isPricedCorrectly = i.Price0 != nil && i.Price1 != nil && i.Amount0 != nil && i.Amount1 != nil && i.AmountUSD != nil
		case *itypes.Swap:
			itemKey = fmt.Sprintf("(%v, %v)", i.Type, i.ProcessingType.ToString())
			isPricedCorrectly = i.Price0 != nil && i.Price1 != nil && i.Amount0 != nil && i.Amount1 != nil && i.AmountUSD != nil
		case *itypes.Transfer:
			itemKey = fmt.Sprintf("(%v, %v)", i.Type, i.ProcessingType.ToString())
			isPricedCorrectly = i.AmountUSD != nil
		}

		if itemKey != defaultKey {
			count, ok := distribution[itemKey]
			if !ok {
				distribution[itemKey] = 1
			} else {
				distribution[itemKey] = count + 1
			}
		}

		if isPricedCorrectly {
			bs.EventsPriced += 1
		}
	}
	bs.EventsUserDistribution = distribution
	bs.IndexingTimeNanos = uint64(pricingTime.UnixNano())
	bs.ProcessingDurationNanos = uint64(processingTime.Sub(startTime).Nanoseconds())
	bs.PricingDurationNanos = uint64(pricingTime.Sub(processingTime).Nanoseconds())
}

type Payload struct {
	NodeMoniker       string
	NodeID            uuid.UUID
	NodeVersion       string
	Environment       string
	Network           string
	allowPricingState bool
	BlockSynopsis     *itypes.BlockSynopsis
	NewDexes          []itypes.UniV2Metadata
	Items             []interface{}
}

func (n *NodeImpl) genPayload(bs *itypes.BlockSynopsis,
	items []interface{},
	newDexes []itypes.UniV2Metadata) *Payload {
	nonNilUserItems := []interface{}{}
	for _, item := range items {
		if item == nil {
			continue
		}
		switch i := item.(type) {
		case *itypes.Mint:
			if i.ProcessingType == itypes.UserRequested {
				nonNilUserItems = append(nonNilUserItems, i)
			}
		case *itypes.Burn:
			if i.ProcessingType == itypes.UserRequested {
				nonNilUserItems = append(nonNilUserItems, i)
			}
		case *itypes.Swap:
			if i.ProcessingType == itypes.UserRequested {
				nonNilUserItems = append(nonNilUserItems, i)
			}
		case *itypes.Transfer:
			if i.ProcessingType == itypes.UserRequested {
				nonNilUserItems = append(nonNilUserItems, i)
			}
		}
	}
	env := "staging"
	if n.prodcheck {
		env = "production tagged " + version.GetGitTag()
	}
	return &Payload{
		NodeMoniker:   n.moniker,
		NodeID:        n.nodeID,
		Environment:   env,
		NodeVersion:   strings.Trim(cfg.SFmt(version.GetVersionStrings()), " "),
		Network:       n.network,
		BlockSynopsis: bs,
		Items:         nonNilUserItems,
		NewDexes:      newDexes,
	}
}

func (n *NodeImpl) syncStartHeight() uint64 {
	// start by assuming cfg height is correct
	startBlock := n.startBlock
	var lbLatestHeightUint64 uint64

	// Check localBackend
	if !n.skipResumeLocal {
		lbLatestHeight, ok, err := n.LocalBackend.Get(lb.KeyLatestHeight)
		if err != nil {
			n.log.Fatal(fmt.Sprintf("error while fetching latest height from localbackend: %v", err))
		}
		if !ok {
			n.log.Warn("local backend does not have record for latest height, assuming new LB")
		} else {
			if err := util.GobDecode(lbLatestHeight, &lbLatestHeightUint64); err != nil {
				n.log.Fatal(fmt.Sprintf("wrong latest height encoding: %s", err))
			}
			if lbLatestHeightUint64 > startBlock {
				n.log.Warn(fmt.Sprintf("local backend reports latest height as %v but config requested %v. overriding config",
					lbLatestHeightUint64,
					n.startBlock))
				startBlock = lbLatestHeightUint64
			}
		}
	}

	// Check resume URL
	if !n.skipResumeRemote {
		remoteLatestHeight, err := n.getRemoteLatestheight()
		if err != nil {
			n.log.Fatal(fmt.Sprintf("error while fetching latest height from remote: %v", err))
		}
		if remoteLatestHeight < startBlock {
			n.log.Fatal(fmt.Sprintf("remote reports latest height as %v but either cfg start height or localBackend height disallows this", remoteLatestHeight),
				"cfg start", n.startBlock,
				"lb latest", lbLatestHeightUint64)
		}
		if remoteLatestHeight > startBlock {
			startBlock = remoteLatestHeight
		}
	}

	n.log.Info("start block height set", "start", startBlock)
	return startBlock
}

func (n *NodeImpl) getRemoteLatestheight() (uint64, error) {
	resp, err := http.Get(n.remoteResumeURL)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var responseObject itypes.ResumeAPIResponse
	if err := json.Unmarshal(body, &responseObject); err != nil {
		return 0, err
	}

	n.log.Info("resuming from block height (via API response): ", responseObject.Data.Height)
	return responseObject.Data.Height, nil
}

// Creates a new node service with spf13/viper fields (yaml)
// CONTRACT: NodeCFGFields enlists all the fields to be accessed in this function
func NewNodeWithViperFields(log logger.Logger) (service.Service, error) {
	// ensure field integrity for viper
	for _, mf := range NodeCFGFields {
		err := cfg.EnsureFieldIntegrity(NodeCFGSection, mf)
		if err != nil {
			return nil, err
		}
	}

	if viper.GetBool(NodeCFGSection + ".prodcheck") {
		if version.GetGitTag() == version.UNTAGGED_GITTAG {
			log.Fatal("cannot run a untagged indexer on production. exiting")
		}
	} else {
		log.Warn("prodcheck is unset, make sure this indexer does not run in production")
	}

	var (
		lbType     = viper.GetString(NodeCFGSection + ".localBackendType")
		outsType   = viper.GetString(NodeCFGSection + ".outputSinkType")
		ethrpcType = viper.GetString(NodeCFGSection + ".ethRPCType")
	)

	// Setup local backend
	var (
		localBackend         lb.LocalBackend
		err                  error
		isLocalBackendNoneDB = false
	)
	switch lbType {
	case "badgerdb":
		localBackend, err = lb.NewBadgerDBWithViperFields(log.With("service", "localbackend"))
		if err != nil {
			return nil, err
		}
	case "none":
		localBackend, err = lb.NewNoneDB(log.With("service", "localbackend"))
		if err != nil {
			return nil, err
		}
		isLocalBackendNoneDB = true
	default:
		log.Fatal("unsupported localbackend: " + lbType)
	}

	// Setup output link
	if outsType != "rabbitmq" {
		log.Fatal("unsupported outputsink: " + outsType)
	}
	outputSink, err := outs.NewRabbitMQOutputSinkWithViperFields(log.With("service", "outputsink"), &iamqp.AMQPImpl{})
	if err != nil {
		return nil, err
	}

	// Setup ethrpc
	if ethrpcType != "mspool" {
		log.Fatal("unsupported ethrpc: " + ethrpcType)
	}
	_ethrpc, err := ethrpc.NewMSPoolEthRPCWithViperFields(log.With("service", "ethrpc"))
	if err != nil {
		return nil, err
	}

	node := &NodeImpl{
		log:                             log.With("service", "node"),
		EthRPC:                          _ethrpc,
		LocalBackend:                    localBackend,
		OutputSink:                      outputSink,
		startBlock:                      viper.GetUint64(NodeCFGSection + ".startBlock"),
		skipResumeRemote:                viper.GetBool(NodeCFGSection + ".skipResumeRemote"),
		skipResumeLocal:                 viper.GetBool(NodeCFGSection + ".skipResumeLocal"),
		remoteResumeURL:                 viper.GetString(NodeCFGSection + ".remoteResumeURL"),
		eventsToIndex:                   viper.GetStringSlice(NodeCFGSection + ".eventsToIndex"),
		maxCPUParallels:                 viper.GetInt(NodeCFGSection + ".maxCPUParallels"),
		maxBlockSpanPerCall:             viper.GetUint64(NodeCFGSection + ".maxBlockSpanPerCall"),
		quitCh:                          make(chan struct{}, 1),
		moniker:                         viper.GetString(NodeCFGSection + ".moniker"),
		network:                         viper.GetString(NodeCFGSection + ".network"),
		allowPricingState:               !isLocalBackendNoneDB,
		pricingChainlinkOraclesDumpFile: viper.GetString(NodeCFGSection + ".pricingChainlinkOraclesDumpFile"),
		oldPricerOracleMap:              viper.GetString(NodeCFGSection + ".oldPricerOracleMap"),
		pricingDexDumpFile:              viper.GetString(NodeCFGSection + ".pricingDexDumpFile"),
		prodcheck:                       viper.GetBool(NodeCFGSection + ".prodcheck"),
	}
	node.BaseService = *service.NewBaseService(log, "node", node)
	return node, nil
}

func mergeTopics(requestedEvents, requiredEvents []common.Hash) map[common.Hash]itypes.ProcessingType {
	maxEvents := len(requestedEvents)
	if len(requiredEvents) > maxEvents {
		maxEvents = len(requiredEvents)
	}
	mergedMap := make(map[common.Hash]itypes.ProcessingType, maxEvents)

	for _, item := range requiredEvents {
		mergedMap[item] = itypes.PricingEngineRequest
	}
	for _, item := range requestedEvents {
		mergedMap[item] = itypes.UserRequested
	}
	return mergedMap
}

func (n *NodeImpl) configureBackoff() {
	n.backoff = backoff.NewConstantBackOff(time.Second * 2)
}
