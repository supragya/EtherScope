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
	uniswapv2 "github.com/Blockpour/Blockpour-Geth-Indexer/libs/processors/uniswapV2"
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
	mergedTopics     map[common.Hash]itypes.ProcessingType // information on topics to index
	mergedTopicsKeys []common.Hash                         // cached keys of mergedTopics
	indexedHeight    uint64
	currentHeight    uint64
	quitCh           chan struct{}

	procUniV2 uniswapv2.UniswapV2Processor
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
			Type:   "statistics",
			Height: block,
			Time:   _time,
		}

		var wg sync.WaitGroup
		var items []interface{} = make([]interface{}, len(logs))

		for idx, _log := range logs {
			go n.decodeLog(_log, &items, idx, &blockMeta, &wg)
		}
		wg.Wait()
		for idx, _ := range logs {
			if items[idx] == nil {
				continue
			}
			n.log.Info("logg", idx, items[idx])
		}
	}
}

func (n *NodeImpl) decodeLog(l types.Log,
	items *[]interface{},
	idx int,
	bm *itypes.BlockSynopsis,
	wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	primaryTopic := l.Topics[0]
	switch primaryTopic {
	// ---- Uniswap V2 ----
	case itypes.UniV2MintTopic:
		// instrumentation.MintV2Found.Inc()
		n.procUniV2.ProcessUniV2Mint(l, items, idx, bm)
	case itypes.UniV2BurnTopic:
		// instrumentation.BurnV2Found.Inc()
		n.procUniV2.ProcessUniV2Burn(l, items, idx, bm)
	case itypes.UniV2SwapTopic:
		// instrumentation.SwapV2Found.Inc()
		n.procUniV2.ProcessUniV2Swap(l, items, idx, bm)

		// // ---- Uniswap V3 ----
		// case itypes.UniV3MintTopic:
		// 	// instrumentation.MintV3Found.Inc()
		// 	n.processUniV3Mint(l, items, bm, mt)
		// case itypes.UniV3BurnTopic:
		// 	// instrumentation.BurnV3Found.Inc()
		// 	n.processUniV3Burn(l, items, bm, mt)
		// case itypes.UniV3SwapTopic:
		// 	// instrumentation.SwapV3Found.Inc()
		// 	n.processUniV3Swap(l, items, bm, mt)

		// // ---- ERC 20 ----
		// case itypes.ERC20TransferTopic:
		// 	// instrumentation.TfrFound.Inc()
		// 	n.processERC20Transfer(l, items, bm, mt)
	}
}

// // For Uniswap V3

// func (r *NodeImpl) processUniV3Mint(
// 	l types.Log,
// 	items *[]interface{},
// 	bm *itypes.BlockSynopsis,
// 	mt *sync.Mutex,
// ) {
// 	callopts := GetBlockCallOpts(l.BlockNumber)
// 	// Test if the contract is a UniswapV3Pair type contract
// 	if !r.isUniswapV3(l.Address, callopts) {
// 		return
// 	}

// 	sender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	ok, _, am0, am1 := InfoUniV3Mint(l)
// 	if !ok {
// 		return
// 	}

// 	// Test if the contract is a UniswapV3NFT type contract
// 	t0, t1, err := r.da.GetTokensUniV3(l.Address, callopts)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV3(am0, am1, callopts, l.Address)
// 	if !ok {
// 		return
// 	}

// 	reserves, err := r.da.GetERC20Balances([]util.Tuple2[common.Address, common.Address]{
// 		{l.Address, t0}, {l.Address, t1},
// 	}, callopts)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, l, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

// 	mint := itypes.Mint{
// 		Type:         "uniswapv3mint",
// 		Network:      r.dbconn.ChainID,
// 		LogIdx:       l.Index,
// 		Transaction:  l.TxHash,
// 		Time:         bm.Time,
// 		Height:       l.BlockNumber,
// 		Sender:       sender,
// 		TxSender:     sender,
// 		PairContract: l.Address,
// 		Token0:       t0,
// 		Token1:       t1,
// 		Amount0:      f0,
// 		Amount1:      f1,
// 		Reserve0:     util.DivideBy10pow(reserves[0].Second, t0d),
// 		Reserve1:     util.DivideBy10pow(reserves[1].Second, t1d),
// 		AmountUSD:    amountusd,
// 		Price0:       token0Price,
// 		Price1:       token1Price,
// 	}

// 	AddToSynopsis(mt, bm, mint, items, "mint", true)
// 	instrumentation.MintV3Processed.Inc()
// }

// func (r *NodeImpl) processUniV3Burn(
// 	l types.Log,
// 	items *[]interface{},
// 	bm *itypes.BlockSynopsis,
// 	mt *sync.Mutex,
// ) {
// 	callopts := GetBlockCallOpts(l.BlockNumber)

// 	// Test if the contract is a UniswapV3Pair type contract
// 	if !r.isUniswapV3(l.Address, callopts) {
// 		return
// 	}

// 	sender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	ok, _, am0, am1 := InfoUniV3Burn(l)
// 	if !ok {
// 		return
// 	}

// 	t0, t1, err := r.da.GetTokensUniV3(l.Address, callopts)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV3(am0, am1, callopts, l.Address)
// 	if !ok {
// 		return
// 	}

// 	reserves, err := r.da.GetERC20Balances([]util.Tuple2[common.Address, common.Address]{
// 		{l.Address, t0}, {l.Address, t1},
// 	}, callopts)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, l, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

// 	burn := itypes.Burn{
// 		Type:         "uniswapv3burn",
// 		Network:      r.dbconn.ChainID,
// 		LogIdx:       l.Index,
// 		Transaction:  l.TxHash,
// 		Time:         bm.Time,
// 		Height:       l.BlockNumber,
// 		Sender:       sender,
// 		TxSender:     sender,
// 		PairContract: l.Address,
// 		Token0:       t0,
// 		Token1:       t1,
// 		Amount0:      f0,
// 		Amount1:      f1,
// 		Reserve0:     util.DivideBy10pow(reserves[0].Second, t0d),
// 		Reserve1:     util.DivideBy10pow(reserves[1].Second, t1d),
// 		AmountUSD:    amountusd,
// 		Price0:       token0Price,
// 		Price1:       token1Price,
// 	}

// 	AddToSynopsis(mt, bm, burn, items, "burn", true)
// 	instrumentation.BurnV3Processed.Inc()
// }

// func (r *NodeImpl) processUniV3Swap(
// 	l types.Log,
// 	items *[]interface{},
// 	bm *itypes.BlockSynopsis,
// 	mt *sync.Mutex,
// ) {
// 	callopts := GetBlockCallOpts(l.BlockNumber)

// 	// Test if the contract is a UniswapV3 NFT type contract
// 	if !r.isUniswapV3(l.Address, callopts) {
// 		return
// 	}

// 	ok, sender, receiver, am0, am1 := InfoUniV3Swap(l)
// 	if !ok {
// 		return
// 	}

// 	// Test if the contract is a UniswapV3NFT type contract
// 	t0, t1, err := r.da.GetTokensUniV3(l.Address, callopts)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV3(am0, am1, callopts, l.Address)
// 	if !ok {
// 		return
// 	}

// 	reserves, err := r.da.GetERC20Balances([]util.Tuple2[common.Address, common.Address]{
// 		{l.Address, t0}, {l.Address, t1},
// 	}, callopts)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, l, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

// 	txSender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	swap := itypes.Swap{
// 		Type:         "uniswapv3swap",
// 		Network:      r.dbconn.ChainID,
// 		LogIdx:       l.Index,
// 		Transaction:  l.TxHash,
// 		Time:         bm.Time,
// 		Height:       l.BlockNumber,
// 		Sender:       sender,
// 		TxSender:     txSender,
// 		Receiver:     receiver,
// 		PairContract: l.Address,
// 		Token0:       t0,
// 		Token1:       t1,
// 		Amount0:      f0,
// 		Amount1:      f1,
// 		Reserve0:     util.DivideBy10pow(reserves[0].Second, t0d),
// 		Reserve1:     util.DivideBy10pow(reserves[0].Second, t1d),
// 		AmountUSD:    amountusd,
// 		Price0:       token0Price,
// 		Price1:       token1Price,
// 	}

// 	AddToSynopsis(mt, bm, swap, items, "swap", true)
// 	instrumentation.SwapV3Processed.Inc()
// }

// func (r *NodeImpl) isUniswapV3(address common.Address,
// 	callopts *bind.CallOpts) bool {
// 	_, _, err := r.da.GetTokensUniV3(address, callopts)
// 	if err == nil {
// 		return true
// 	}

// 	if !util.IsExecutionReverted(err) {
// 		util.ENOKS(2, err)
// 	}
// 	return false
// }

// func (r *NodeImpl) isUniswapV3NFT(address common.Address,
// 	callopts *bind.CallOpts) bool {
// 	_, _, err := r.da.GetTokensUniV3NFT(address, big.NewInt(1), callopts)
// 	if err == nil {
// 		return true
// 	}

// 	if !util.IsExecutionReverted(err) {
// 		util.ENOKS(2, err)
// 	}
// 	return false
// }

// func (r *NodeImpl) GetFormattedAmountsUniV3NFT(amount0 *big.Int,
// 	amount1 *big.Int,
// 	tokenID *big.Int,
// 	callopts *bind.CallOpts,
// 	address common.Address) (ok bool,
// 	formattedAmount0 *big.Float,
// 	formattedAmount1 *big.Float,
// 	token0Decimals uint8,
// 	token1Decimals uint8) {
// 	t0, t1, err := r.da.GetTokensUniV3NFT(address, tokenID, callopts)
// 	if err != nil {
// 		return false,
// 			big.NewFloat(0.0),
// 			big.NewFloat(0.0),
// 			0,
// 			0
// 	}

// 	token0Decimals, err = r.da.GetERC20Decimals(t0, callopts)
// 	if util.IsExecutionReverted(err) {
// 		// Non ERC-20 contract
// 		token0Decimals = 0
// 	} else {
// 		if util.IsEthErr(err) {
// 			return false,
// 				big.NewFloat(0.0),
// 				big.NewFloat(0.0),
// 				0,
// 				0
// 		}
// 		util.ENOKS(2, err)
// 	}

// 	token1Decimals, err = r.da.GetERC20Decimals(t1, callopts)
// 	if util.IsExecutionReverted(err) {
// 		// Non ERC-20 contract
// 		token1Decimals = 0
// 	} else {
// 		if util.IsEthErr(err) {
// 			return false,
// 				big.NewFloat(0.0),
// 				big.NewFloat(0.0),
// 				0,
// 				0
// 		}
// 		util.ENOKS(2, err)
// 	}

// 	return true,
// 		util.DivideBy10pow(amount0, token0Decimals),
// 		util.DivideBy10pow(amount1, token1Decimals),
// 		token0Decimals,
// 		token1Decimals
// }

// func (r *NodeImpl) GetFormattedAmountsUniV3(amount0 *big.Int,
// 	amount1 *big.Int,
// 	callopts *bind.CallOpts,
// 	address common.Address) (ok bool,
// 	formattedAmount0 *big.Float,
// 	formattedAmount1 *big.Float,
// 	token0Decimals uint8,
// 	token1Decimals uint8) {
// 	t0, t1, err := r.da.GetTokensUniV3(address, callopts)
// 	if err != nil {
// 		return false,
// 			big.NewFloat(0.0),
// 			big.NewFloat(0.0),
// 			0,
// 			0
// 	}

// 	token0Decimals, err = r.da.GetERC20Decimals(t0, callopts)
// 	if util.IsExecutionReverted(err) {
// 		// Non ERC-20 contract
// 		token0Decimals = 0
// 	} else {
// 		if util.IsEthErr(err) {
// 			return false,
// 				big.NewFloat(0.0),
// 				big.NewFloat(0.0),
// 				0,
// 				0
// 		}
// 		util.ENOKS(2, err)
// 	}

// 	token1Decimals, err = r.da.GetERC20Decimals(t1, callopts)
// 	if util.IsExecutionReverted(err) {
// 		// Non ERC-20 contract
// 		token1Decimals = 0
// 	} else {
// 		if util.IsEthErr(err) {
// 			return false,
// 				big.NewFloat(0.0),
// 				big.NewFloat(0.0),
// 				0,
// 				0
// 		}
// 		util.ENOKS(2, err)
// 	}

// 	return true,
// 		util.DivideBy10pow(amount0, token0Decimals),
// 		util.DivideBy10pow(amount1, token1Decimals),
// 		token0Decimals,
// 		token1Decimals
// }

// // For ERC20
// func setupERC20TransferRestrictions(events []common.Hash) *ERC20TransferRestrictions {
// 	isERC20TransferToBeIndexed := false
// 	for _, e := range events {
// 		if e == itypes.ERC20TransferTopic {
// 			isERC20TransferToBeIndexed = true
// 		}
// 	}

// 	if !isERC20TransferToBeIndexed {
// 		return nil
// 	}

// 	var (
// 		restrictionType = viper.GetString("erc20transfer.restrictionType")
// 		whitelistFile   = viper.GetString("erc20transfer.whitelistFile")
// 		whitelistMap    = make(map[common.Address]bool)
// 	)

// 	var _type ERC20RestrictionType
// 	switch restrictionType {
// 	case "none":
// 		_type = None
// 		// short circuit
// 		return &ERC20TransferRestrictions{_type, &whitelistMap}
// 	case "to":
// 		_type = To
// 	case "from":
// 		_type = From
// 	case "both":
// 		_type = Both
// 	case "either":
// 		_type = Either
// 	default:
// 		panic("unknown ERC20RestrictionType")
// 	}

// 	file, err := os.Open(whitelistFile)
// 	util.ENOK(err)

// 	_bytes, err := ioutil.ReadAll(file)
// 	util.ENOK(err)

// 	whitelist := []common.Address{}
// 	util.ENOK(json.Unmarshal(_bytes, &whitelist))

// 	for _, ra := range whitelist {
// 		whitelistMap[ra] = true
// 	}

// 	return &ERC20TransferRestrictions{_type, &whitelistMap}
// }

// func allowIndexing(r *ERC20TransferRestrictions, from common.Address, to common.Address) bool {
// 	if r._type == None {
// 		return true
// 	}
// 	var (
// 		whFrom = false
// 		whTo   = false
// 	)
// 	if _, ok := (*r.whitelist)[from]; ok {
// 		whFrom = true
// 	}
// 	if _, ok := (*r.whitelist)[to]; ok {
// 		whTo = true
// 	}
// 	switch r._type {
// 	case None:
// 		return true
// 	case To:
// 		return whTo
// 	case From:
// 		return whFrom
// 	case Both:
// 		return whTo && whFrom
// 	case Either:
// 		return whTo || whFrom
// 	}
// 	return false
// }

// func (r *NodeImpl) processERC20Transfer(
// 	l types.Log,
// 	items *[]interface{},
// 	bm *itypes.BlockSynopsis,
// 	mt *sync.Mutex,
// ) {
// 	ok, sender, recv, amt := InfoTransfer(l)
// 	if !ok {
// 		return
// 	}

// 	if !allowIndexing(r.erc20TransferRestrictions, sender, recv) {
// 		return
// 	}

// 	callopts := GetBlockCallOpts(l.BlockNumber)

// 	ok, formattedAmount := r.GetFormattedAmount(amt, callopts, l.Address)
// 	if !ok {
// 		return
// 	}

// 	tokenPrice := r.da.GetRateForBlock(callopts, util.Tuple2[common.Address, *big.Float]{l.Address, formattedAmount})

// 	txSender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	transfer := itypes.Transfer{
// 		Type:                "erc20transfer",
// 		Network:             r.dbconn.ChainID,
// 		LogIdx:              l.Index,
// 		Transaction:         l.TxHash,
// 		Time:                bm.Time,
// 		Height:              l.BlockNumber,
// 		Token:               l.Address,
// 		Sender:              sender,
// 		TxSender:            txSender,
// 		Receiver:            recv,
// 		Amount:              formattedAmount,
// 		AmountUSD:           tokenPrice.Price,
// 		PriceDerivationMeta: tokenPrice,
// 	}

// 	AddToSynopsis(mt, bm, transfer, items, "transfer", true)
// 	instrumentation.TfrProcessed.Inc()
// }

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
