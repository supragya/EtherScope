package node

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/ethrpc"
	lb "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	outs "github.com/Blockpour/Blockpour-Geth-Indexer/services/output_sink"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

type processingType uint

const (
	OnlyPricingEngine processingType = iota
	UserRequested
)

var (
	NodeCFGSection   = "node"
	NodeCFGNecessity = "always needed"
	NodeCFGHeader    = cfg.SArr("node is core indexing service for bgidx",
		"node is tasked with initiating other services such as",
		"localbackend (badger-db) and outputsink (rabbitmq)")
	NodeCFGFields = [...]cfg.Field{
		{
			Name:      "maxCPUParallels",
			Type:      "uint",
			Necessity: "always needed",
			Info:      cfg.SArr("maximum number of CPU threads to give to bgidx"),
			Default:   4,
		},
		{
			Name:      "startBlock",
			Type:      "uint64",
			Necessity: "always needed",
			Info: cfg.SArr("user defined blockheight to start sync from.",
				"this may be overriden at runtime using resume from localbackend",
				"and remoteHTTP endpoint"),
			Default: 15865859,
		},
		{
			Name:      "skipResumeRemote",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("disables fetch for blockheight from remoteHTTP"),
			Default:   false,
		},
		{
			Name:      "skipResumeLocal",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("disables fetch for blockheight from localbackend"),
			Default:   false,
		},
		{
			Name:      "remoteResumeURL",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("remoteHTTP URL for fetching blockheight to resume from"),
			Default:   "https://myremote.blockpour.com",
		},
		{
			Name:      "localBackendType",
			Type:      "string",
			Necessity: "always needed",
			Info: cfg.SArr("type of local backend indexer should use.",
				"only possible type right now is `badgerdb`"),
			Default: "badgerdb",
		},
		{
			Name:      "outputSinkType",
			Type:      "string",
			Necessity: "always needed",
			Info: cfg.SArr("type of output sink backend indexer should",
				"offload indexed information to. only possible type",
				"right now is `rabbitmq`"),
			Default: "rabbitmq",
		},
		{
			Name:      "ethRPCType",
			Type:      "string",
			Necessity: "always needed",
			Info: cfg.SArr("type of ethrpc handler to route requests",
				"through. only possible type right now is `mspool`"),
			Default: "mspool",
		},
		{
			Name:      "eventsToIndex",
			Type:      "[]string",
			Necessity: "always needed",
			Info: cfg.SArr("ethereum events to index. events listed here",
				"are not guaranteed to be the only calls made",
				"to underlying rpc for processing, but are guaranteed",
				"to be the only events presented to the output sink",
				"could be one or many of the following:",
				"- UniswapV2Swap",
				"- UniswapV2Mint",
				"- UniswapV2Burn",
				"- UniswapV3Swap",
				"- UniswapV3Mint",
				"- UniswapV3Burn",
				"- ERC20Transfer"),
			Default: "\n    - ERC20Transfer\n    - UniswapV2Swap",
		},
		{
			Name:      "maxBlockSpanPerCall",
			Type:      "uint64",
			Necessity: "always needed",
			Info: cfg.SArr("number of blocks to fetch logs for at the",
				"beginning of processing loop"),
			Default: 5,
		},
	}
)

type NodeImpl struct {
	service.BaseService

	log          logger.Logger
	EthRPC       ethrpc.EthRPC   // HA upstream connection to rpc nodes, uses mspool
	LocalBackend lb.LocalBackend // Local database for caching / processing
	OutputSink   outs.OutputSink // Consumer for offloading processed data

	// Configs
	startBlock          uint64   // User defined startBlock
	skipResumeRemote    bool     // skip checking remote for resume height
	skipResumeLocal     bool     // skip checking localbackend for resume height
	remoteResumeURL     string   // URL to use for resume height GET request
	eventsToIndex       []string // user requested events to index in string form
	maxCPUParallels     int      // user requested CPU threads to allocate to the process
	maxBlockSpanPerCall uint64   // max block spans to log per initial filtering call

	// Internal Data Structures
	mergedTopics     map[common.Hash]processingType // information on topics to index
	mergedTopicsKeys []common.Hash                  // cached keys of mergedTopics
	indexedHeight    uint64
	currentHeight    uint64
	quitCh           chan struct{}
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

	if err := n.EthRPC.Start(ctx); err != nil {
		return err
	}

	if err := n.LocalBackend.Start(ctx); err != nil {
		return err
	}

	if err := n.OutputSink.Start(ctx); err != nil {
		return err
	}

	// Setup what to index
	requestedEvents, err := util.ConstructTopics(n.eventsToIndex)
	if err != nil {
		return err
	}
	// required topics by the pricing engine
	requiredEvents := []common.Hash{itypes.UniV2MintTopic,
		itypes.UniV2BurnTopic,
		itypes.UniV2SwapTopic}

	n.mergedTopics = mergeTopics(requestedEvents, requiredEvents)
	keys := make([]common.Hash, len(n.mergedTopics))

	i := 0
	for val, _ := range n.mergedTopics {
		keys[i] = val
		i++
	}
	n.mergedTopicsKeys = keys

	// TODO: Do height syncup using both LocalBackend and remote http
	// startHeight, err := n.getResumeHeight()
	n.indexedHeight = n.startBlock

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

				util.ENOK(err)
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

				n.log.Info(fmt.Sprintf("chainhead: %d (+%d away), indexing [%d - %d]",
					n.currentHeight, n.currentHeight-n.indexedHeight, n.indexedHeight, endingBlock))

				// instrumentation.CurrentBlock.Set(float64(n.currentHeight))

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
				// instrumentation.ProcessedBlock.Set(float64(r.indexedHeight))

				if isOnHead {
					break
				}
			}
		case <-n.quitCh:
			// TODO: Graceful exit
			n.log.Info("quitting realtime indexer")
		}
	}
}

func (n *NodeImpl) processBatchedBlockLogs(logs []types.Log, start uint64, end uint64) {
	// Assuming for any height H, either we will have all the concerned logs
	// or not even one
	kv := GroupByBlockNumber(logs)

	for block := start; block <= end; block++ {
		_time, err := n.EthRPC.GetBlockTimestamp(block)
		util.ENOK(err)

		logs, _ := kv[block]
		blockMeta := itypes.BlockSynopsis{
			Type:   "stats",
			Height: block,
			Time:   _time,
		}

		var wg sync.WaitGroup
		var mt sync.Mutex
		var items []interface{} = make([]interface{}, len(logs))

		for idx, _log := range logs {
			wg.Add(1)
			go n.decodeLog(_log, &mt, &items, idx, &blockMeta, &wg)
		}
		wg.Wait()
	}
}

func (n *NodeImpl) decodeLog(l types.Log,
	mt *sync.Mutex,
	items *[]interface{},
	idx int,
	bm *itypes.BlockSynopsis,
	wg *sync.WaitGroup) {
	defer wg.Done()

	// primaryTopic := l.Topics[0]
	// switch primaryTopic {
	// // ---- Uniswap V2 ----
	// case itypes.UniV2MintTopic:
	// 	instrumentation.MintV2Found.Inc()
	// 	r.processUniV2Mint(l, items, bm, mt)
	// case itypes.UniV2BurnTopic:
	// 	instrumentation.BurnV2Found.Inc()
	// 	r.processUniV2Burn(l, items, bm, mt)
	// case itypes.UniV2SwapTopic:
	// 	instrumentation.SwapV2Found.Inc()
	// 	r.processUniV2Swap(l, items, bm, mt)

	// // ---- Uniswap V3 ----
	// case itypes.UniV3MintTopic:
	// 	instrumentation.MintV3Found.Inc()
	// 	r.processUniV3Mint(l, items, bm, mt)
	// case itypes.UniV3BurnTopic:
	// 	instrumentation.BurnV3Found.Inc()
	// 	r.processUniV3Burn(l, items, bm, mt)
	// case itypes.UniV3SwapTopic:
	// 	instrumentation.SwapV3Found.Inc()
	// 	r.processUniV3Swap(l, items, bm, mt)

	// // ---- ERC 20 ----
	// case itypes.ERC20TransferTopic:
	// 	instrumentation.TfrFound.Inc()
	// 	r.processERC20Transfer(l, items, bm, mt)

	// }
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

	var (
		lbType     = viper.GetString(NodeCFGSection + ".localBackendType")
		outsType   = viper.GetString(NodeCFGSection + ".outputSinkType")
		ethrpcType = viper.GetString(NodeCFGSection + ".ethRPCType")
	)

	// Setup local backend
	if lbType != "badgerdb" {
		log.Fatal("unsupported localbackend: " + lbType)
	}
	localBackend, err := lb.NewBadgerDBWithViperFields(log.With("service", "localbackend"))
	if err != nil {
		return nil, err
	}

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
	_ethrpc, err := ethrpc.NewMSPoolEthRPCWithViperFields(log.With("service", "ethrpc"), localBackend)
	if err != nil {
		return nil, err
	}

	node := &NodeImpl{
		log:                 log.With("service", "node"),
		EthRPC:              _ethrpc,
		LocalBackend:        localBackend,
		OutputSink:          outputSink,
		startBlock:          viper.GetUint64(NodeCFGSection + ".startBlock"),
		skipResumeRemote:    viper.GetBool(NodeCFGSection + ".skipResumeRemote"),
		skipResumeLocal:     viper.GetBool(NodeCFGSection + ".skipResumeLocal"),
		remoteResumeURL:     viper.GetString(NodeCFGSection + ".remoteResumeURL"),
		eventsToIndex:       viper.GetStringSlice(NodeCFGSection + ".eventsToIndex"),
		maxCPUParallels:     viper.GetInt(NodeCFGSection + ".maxCPUParallels"),
		maxBlockSpanPerCall: viper.GetUint64(NodeCFGSection + ".maxBlockSpanPerCall"),
		quitCh:              make(chan struct{}, 1),
	}
	node.BaseService = *service.NewBaseService(log, "node", node)
	return node, nil
}

func mergeTopics(requestedEvents, requiredEvents []common.Hash) map[common.Hash]processingType {
	maxEvents := len(requestedEvents)
	if len(requiredEvents) > maxEvents {
		maxEvents = len(requiredEvents)
	}
	mergedMap := make(map[common.Hash]processingType)

	for _, item := range requiredEvents {
		mergedMap[item] = OnlyPricingEngine
	}
	for _, item := range requestedEvents {
		mergedMap[item] = UserRequested
	}
	return mergedMap
}
