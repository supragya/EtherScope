package oraclenode

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/ethrpc"
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

type OracleNodeImpl struct {
	service.BaseService

	log        logger.Logger
	EthRPC     ethrpc.EthRPC   // HA upstream connection to rpc nodes, uses mspool
	OutputSink outs.OutputSink // Consumer for offloading processed data

	// Configs
	startBlock          uint64 // User defined startBlock
	skipResumeRemote    bool   // skip checking remote for resume height
	remoteResumeURL     string // URL to use for resume height GET request
	prodcheck           bool   // checks for prod grade settings
	maxCPUParallels     int    // user requested CPU threads to allocate to the process
	maxBlockSpanPerCall uint64 // max block spans to log per initial filtering call
	feedRegistry        common.Address
	feedFile            string

	// Internal Data Structures
	moniker       string    // user defined moniker for this node
	network       string    // user defined evm compatible network name
	nodeID        uuid.UUID // system generated node identifier unique for each run
	indexedHeight uint64
	currentHeight uint64
	feedMap       map[itypes.Tuple2[common.Address, common.Address]]common.Address
	feedMapRev    map[common.Address]itypes.Tuple2[common.Address, common.Address]
	quitCh        chan struct{}

	// Backoff configuration
	backoff *backoff.ConstantBackOff
}

// OnStart starts the Node. It implements service.Service.
func (n *OracleNodeImpl) OnStart(ctx context.Context) error {
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

	if err := n.OutputSink.Start(ctx); err != nil {
		n.log.Info("Error initializing output sink, will reattempt connection until ready")
	}

	// startHeight, err := n.getResumeHeight()
	n.indexedHeight = n.syncStartHeight()

	// Loop for impl
	n.setupInitial()
	go n.loop()

	return nil
}

// OnStop stops the Node. It implements service.Service
func (n *OracleNodeImpl) OnStop() {
	n.quitCh <- struct{}{}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
	}()

	wg.Wait()
}

// SetupInitial sets up initial information that the oracleindexer needs
func (n *OracleNodeImpl) setupInitial() {
	start, end, stride := uint64(12864088), n.indexedHeight, uint64(1000) // Start of feed registry to now, each call indexing 100,000 blocks
	n.log.Infof("need to index %d blocks for feed data, %d to %d", end-start, start, end)

	if _, err := os.Stat(n.feedFile); err == nil {
		// This means that feedfile exists already
		n.log.Info("feed file already exists, reading")
		f, err := os.Open(n.feedFile)
		if err != nil {
			n.log.Fatal(err.Error())
		}
		defer f.Close()
		reader := csv.NewReader(f)
		reader.Read() // First one for heading, don't care
		count := 0

		for {
			items, err := reader.Read()
			if err != nil {
				n.log.Fatal(err.Error())
			}
			height, err := strconv.Atoi(items[0])
			if err != nil {
				n.log.Fatal(err.Error())
			}
			if height > int(n.indexedHeight) {
				break
			}

			asset := common.HexToAddress(items[1])
			denomination := common.HexToAddress(items[2])
			latestAggregator := common.HexToAddress(items[3])
			n.log.Infof("Read Aggregator found (%x:%x) %x", asset, denomination, latestAggregator)

			n.feedMap[itypes.Tuple2[common.Address, common.Address]{
				asset, denomination,
			}] = latestAggregator // aggregator
			n.feedMapRev[latestAggregator] = itypes.Tuple2[common.Address, common.Address]{
				asset, denomination,
			}
			count += 1
		}
		n.log.Infof("Added %d count aggregator", count)
		return
	}

	// otherwise
	f, err := os.Create(n.feedFile) // TODO: make it configurable, check err
	if err != nil {
		n.log.Fatal("error opening csv file")
	}
	defer f.Close()

	w := csv.NewWriter(f)

	if err := w.Write([]string{"height", "asset", "denomination", "latestAggregator"}); err != nil {
		n.log.Fatal("error writing record to file", err)
	}

	for {
		w.Flush()
		if start > end {
			break
		}

		callStart, callEnd := start, start+stride
		if callEnd > end {
			callEnd = end
		}

		logs, err := n.EthRPC.GetFilteredLogs(ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(callStart)),
			ToBlock:   big.NewInt(int64(callEnd)),
			Topics:    [][]common.Hash{[]common.Hash{itypes.ChainLinkFeedConfirmed}},
			Addresses: []common.Address{n.feedRegistry},
		})

		if err != nil {
			n.log.Fatal("error while getting logs for bootup", "error", err)
		}

		for _, log := range logs {
			asset := util.ExtractAddressFromLogTopic(log.Topics[1])
			denomination := util.ExtractAddressFromLogTopic(log.Topics[2])
			latestAggregator := util.ExtractAddressFromLogTopic(log.Topics[3])

			record := []string{strconv.Itoa(int(log.BlockNumber)), asset.Hex(), denomination.Hex(), latestAggregator.Hex()}
			if err := w.Write(record); err != nil {
				n.log.Fatal("error writing record to file", err)
			}

			n.log.Infof("Aggregator found @ %d (%x:%x) %x", log.BlockNumber, asset, denomination, latestAggregator)

			n.feedMap[itypes.Tuple2[common.Address, common.Address]{asset, denomination}] = latestAggregator
			n.feedMapRev[latestAggregator] = itypes.Tuple2[common.Address, common.Address]{asset, denomination}
		}

		start = start + uint64(stride)
	}
}

// Loop implements core indexing logic
func (n *OracleNodeImpl) loop() {
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

				// instrumentation.CurrentBlock.Set(float64(n.currentHeight))

				logs, err := n.EthRPC.GetFilteredLogs(ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(n.indexedHeight + 1)),
					ToBlock:   big.NewInt(int64(endingBlock)),
					Topics: [][]common.Hash{[]common.Hash{
						itypes.ChainLinkFeedConfirmed,
						itypes.ChainLinkAnswerUpdated,
					}},
				})

				if err != nil {
					n.log.Error("encountered error", "error", err)
					continue
				}

				n.processBatchedBlockLogs(logs, n.indexedHeight+1, endingBlock)

				n.indexedHeight = endingBlock

				if isOnHead {
					break
				}
			}
		case <-n.quitCh:
			n.log.Info("quitting realtime indexer")
		}
	}
}

func (n *OracleNodeImpl) processBatchedBlockLogs(logs []types.Log, start uint64, end uint64) {
	// Assuming for any height H, either we will have all the concerned logs
	// or not even one
	kv := GroupByBlockNumber(logs)

	for block := start; block <= end; block++ {
		backoff.Retry(func() error { return n.processBlock(kv, block) }, n.backoff)
	}
}

type ChainLinkUpdate struct {
	Asset        common.Address
	Denomination common.Address
	Aggregator   common.Address
	Answer       itypes.ChainlinkLatestRoundData
}

func (n *OracleNodeImpl) processBlock(kv map[uint64]CLogType, block uint64) error {
	n.log.Info(fmt.Sprintf("processing block %d", block))
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

	var processedItems []ChainLinkUpdate

	for _, _log := range logs {
		roundData, isRoundData := n.decodeLog(_log)
		if !isRoundData {
			continue
		}
		processedItems = append(processedItems, roundData)
	}

	if err != nil {
		n.log.Debug(fmt.Sprintf("Error processing block %d. Retrying. Error caused by: %s", block, err))
		return err
	}

	// // Package processedItems into payload for output
	payload := n.genPayload(&blockSynopis, processedItems)

	for {
		err = n.OutputSink.Send(payload)
		if err == nil {
			break
		}
		n.log.Warn("Error sending message to output sink: " + fmt.Sprint(err))
		time.Sleep(2 * time.Second)
	}

	return nil
}

func (n *OracleNodeImpl) decodeLog(l types.Log) (ChainLinkUpdate, bool) {
	primaryTopic := l.Topics[0]
	switch primaryTopic {
	case itypes.ChainLinkAnswerUpdated:
		// Ensure that this aggregator exists
		info, ok := n.feedMapRev[l.Address]
		if !ok {
			n.log.Infof("Unknown feed answer update for %s, skipping", l.Address.Hex())
		}
		ans, err := n.EthRPC.GetChainlinkRoundData(l.Address, util.GetBlockCallOpts(l.BlockNumber))
		if err != nil {
			n.log.Infof("Error while getting round data for %x: %x, skipping", l.Address.Hex(), err)
			return ChainLinkUpdate{}, false
		}
		n.log.Infof("Rev map: %s, %s %s", l.Address.Hex(), info.First.Hex(), info.Second.Hex())
		return ChainLinkUpdate{
			Asset:        info.First,
			Denomination: info.Second,
			Aggregator:   l.Address,
			Answer:       ans,
		}, true

	case itypes.ChainLinkFeedConfirmed:
		asset := util.ExtractAddressFromLogTopic(l.Topics[1])
		denomination := util.ExtractAddressFromLogTopic(l.Topics[2])
		latestAggregator := util.ExtractAddressFromLogTopic(l.Topics[3])

		record := []string{strconv.Itoa(int(l.BlockNumber)), asset.Hex(), denomination.Hex(), latestAggregator.Hex()}

		n.log.Infof("Aggregator found @ %d (%x:%x) %x", l.BlockNumber, asset, denomination, latestAggregator)

		f, err := os.OpenFile(n.feedFile, os.O_APPEND, 0777)
		if err != nil {
			n.log.Fatal("error opening csv file")
		}
		defer f.Close()

		w := csv.NewWriter(f)
		if err := w.Write(record); err != nil {
			n.log.Fatal("error writing record to file", err)
		}

		n.feedMap[itypes.Tuple2[common.Address, common.Address]{asset, denomination}] = latestAggregator
		n.feedMapRev[latestAggregator] = itypes.Tuple2[common.Address, common.Address]{asset, denomination}
	}
	return ChainLinkUpdate{}, false
}

type Payload struct {
	NodeMoniker string
	NodeID      uuid.UUID
	NodeVersion string
	Environment string
	Network     string
	Items       []ChainLinkUpdate
}

func (n *OracleNodeImpl) genPayload(bs *itypes.BlockSynopsis,
	items []ChainLinkUpdate) *Payload {
	env := "staging"
	if n.prodcheck {
		env = "production tagged " + version.GetGitTag()
	}
	return &Payload{
		NodeMoniker: n.moniker,
		NodeID:      n.nodeID,
		Environment: env,
		NodeVersion: strings.Trim(cfg.SFmt(version.GetVersionStrings()), " "),
		Network:     n.network,
		Items:       items,
	}
}

func (n *OracleNodeImpl) syncStartHeight() uint64 {
	// start by assuming cfg height is correct
	startBlock := n.startBlock

	// Check resume URL
	if !n.skipResumeRemote {
		remoteLatestHeight, err := n.getRemoteLatestheight()
		if err != nil {
			n.log.Fatal(fmt.Sprintf("error while fetching latest height from remote: %v", err))
		}
		if remoteLatestHeight < startBlock {
			n.log.Fatal(fmt.Sprintf("remote reports latest height as %v but either cfg start height or localBackend height disallows this", remoteLatestHeight),
				"cfg start", n.startBlock)
		}
		if remoteLatestHeight > startBlock {
			startBlock = remoteLatestHeight
		}
	}

	n.log.Info("start block height set", "start", startBlock)
	return startBlock
}

func (n *OracleNodeImpl) getRemoteLatestheight() (uint64, error) {
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
// CONTRACT: OracleNodeCFGFields enlists all the fields to be accessed in this function
func NewOracleNodeWithViperFields(log logger.Logger) (service.Service, error) {
	// ensure field integrity for viper
	for _, mf := range OracleNodeCFGFields {
		err := cfg.EnsureFieldIntegrity(OracleNodeCFGSection, mf)
		if err != nil {
			return nil, err
		}
	}

	// if viper.GetBool(OracleNodeCFGSection + ".prodcheck") {
	// 	if version.GetGitTag() == version.UNTAGGED_GITTAG {
	// 		log.Fatal("cannot run a untagged indexer on production. exiting")
	// 	}
	// } else {
	// 	log.Warn("prodcheck is unset, make sure this indexer does not run in production")
	// }

	var (
		// lbType     = viper.GetString(OracleNodeCFGSection + ".localBackendType")
		outsType   = viper.GetString(OracleNodeCFGSection + ".outputSinkType")
		ethrpcType = viper.GetString(OracleNodeCFGSection + ".ethRPCType")
	)

	// // Setup local backend
	// if lbType != "badgerdb" {
	// 	log.Fatal("unsupported localbackend: " + lbType)
	// }
	// localBackend, err := lb.NewBadgerDBWithViperFields(log.With("service", "localbackend"))
	// if err != nil {
	// 	return nil, err
	// }

	// Setup output link
	if outsType != "rabbitmq" {
		log.Fatal("unsupported outputsink: " + outsType)
	}
	outputSink, err := outs.NewRabbitMQOutputSinkWithViperFields(log.With("service", "outputsink"))
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

	node := &OracleNodeImpl{
		log:                 log.With("service", "oraclenode"),
		EthRPC:              _ethrpc,
		OutputSink:          outputSink,
		startBlock:          viper.GetUint64(OracleNodeCFGSection + ".startBlock"),
		skipResumeRemote:    viper.GetBool(OracleNodeCFGSection + ".skipResumeRemote"),
		remoteResumeURL:     viper.GetString(OracleNodeCFGSection + ".remoteResumeURL"),
		maxCPUParallels:     viper.GetInt(OracleNodeCFGSection + ".maxCPUParallels"),
		maxBlockSpanPerCall: viper.GetUint64(OracleNodeCFGSection + ".maxBlockSpanPerCall"),
		quitCh:              make(chan struct{}, 1),
		moniker:             viper.GetString(OracleNodeCFGSection + ".moniker"),
		network:             viper.GetString(OracleNodeCFGSection + ".network"),
		prodcheck:           viper.GetBool(OracleNodeCFGSection + ".prodcheck"),
		feedRegistry:        common.HexToAddress(viper.GetString(OracleNodeCFGSection + ".chainlinkFeedRegistry")),
		feedFile:            viper.GetString(OracleNodeCFGSection + ".feedFile"),

		feedMap:    make(map[itypes.Tuple2[common.Address, common.Address]]common.Address),
		feedMapRev: make(map[common.Address]itypes.Tuple2[common.Address, common.Address]),
	}
	node.BaseService = *service.NewBaseService(log, "oraclenode", node)
	return node, nil
}

func (n *OracleNodeImpl) configureBackoff() {
	n.backoff = backoff.NewConstantBackOff(time.Second * 2)
}
