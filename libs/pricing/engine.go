package priceresolver

import (
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/gograph"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/ethrpc"
	lb "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type we = gograph.WeightedEdge[common.Address, int64, string, interface{}]
type gg = gograph.Graph[common.Address, int64, string, interface{}]
type ggt = itypes.Tuple2[uint64, *gg]

// Enhanced cached, multistep, graph based batch pricing resolver system
type Engine struct {
	log          logger.Logger
	EthRPC       ethrpc.EthRPC
	LocalBackend lb.LocalBackend

	// On-disk dump information
	chainlinkOraclesDumpFile string
	dexDumpFile              string

	// Internal Data Structures
	lastGraph  *gg
	lastHeight uint64
}

// DefaultEngine is default form of enhanced pricing engine
func NewDefaultEngine(log logger.Logger,
	chainlinkOracledDumpFile string,
	dexDumpFile string,
	ethrpcBackend ethrpc.EthRPC,
	localBackend lb.LocalBackend) *Engine {
	return &Engine{
		log:                      log,
		EthRPC:                   ethrpcBackend,
		LocalBackend:             localBackend,
		chainlinkOraclesDumpFile: chainlinkOracledDumpFile,
		dexDumpFile:              dexDumpFile,
		lastGraph:                nil,
		lastHeight:               0,
	}
}

func (n *Engine) Resolve(resHeight uint64, items *[]interface{}) error {
	// fetch the graph in question
PRICING_GRAPH_FETCH:
	lbLatestHeight, ok, err := n.LocalBackend.Get(lb.KeyLatestHeight)
	if err != nil {
		return err
	}
	if !ok {
		// This means localbackend does not know any graph
		// just yet. Go through dumps
		n.log.Warn("no pricingGraph found in localbackend. reading dump files and syncing")
		err := n.syncDump(resHeight)
		if err != nil {
			return err
		}
		goto PRICING_GRAPH_FETCH
	}

	// ensure lbLatestHeight is less than resHeight
	var lbLatestHeightUint64 uint64
	if err := util.GobDecode(lbLatestHeight, &lbLatestHeightUint64); err != nil {
		n.log.Fatal(fmt.Sprintf("wrong latest height encoding: %s", err))
	}
	if resHeight < lbLatestHeightUint64 {
		return fmt.Errorf("requested resolution for %d while lb already has %d",
			resHeight,
			lbLatestHeightUint64)
	}

	// Get the graph
	graph := gograph.NewGraph[common.Address, int64, string, interface{}](true)
	lbLatestGraph, ok, err := n.LocalBackend.Get(lb.KeyLatestPricingGraph)
	if err := util.GobDecode(lbLatestGraph, graph); err != nil {
		return err
	}

	return n.resolvei(graph, items)
}

func (n *Engine) resolvei(graph *gg, items *[]interface{}) error {
	n.log.Info("found resolvei")
	return nil
}

func (n *Engine) syncDump(resHeight uint64) error {
	n.log.Info("undertaking sync dump")
	chainlinkRecords, err := loadChainlinkCSV(n.chainlinkOraclesDumpFile)
	if err != nil {
		return err
	}
	dexRecords, err := loadDexCSV(n.dexDumpFile)
	if err != nil {
		return err
	}

	n.log.Info("read dump file done, loaded chainlink and dexes in memory",
		"chainlinkrec", len(chainlinkRecords),
		"dexrec", len(dexRecords))

	startTime := time.Now()
	graph := n.genGraph(chainlinkRecords, dexRecords, resHeight)
	genTime := time.Since(startTime)
	n.log.Info("graph gen step completed",
		"edges", graph.GetEdgeCount(),
		"vert", graph.GetVertexCount(),
		"_time", genTime)

	n.log.Info("syncing step to backfill in DB")
	err = n.syncGraphToLB(graph, resHeight)
	if err != nil {
		return err
	}

	n.log.Info("first run syncing complete")
	return nil
}

func (n *Engine) syncGraphToLB(graph *gg, height uint64) error {
	if err := n.LocalBackend.Set(lb.KeyLatestHeight, util.GobEncode(height)); err != nil {
		return err
	}
	if err := n.LocalBackend.Set(lb.KeyLatestPricingGraph, util.GobEncode(graph)); err != nil {
		return err
	}

	n.log.Info("final syncing to disk")
	return n.LocalBackend.Sync()
}

// generates graphs for blocks where
// new chainlink oracles came alive or dex pools were
// made. Ideally the steps should be sparse enough
// (not every block has dex creation or oracle setup)
// to justify memory footprint.
func (n *Engine) genGraph(chainlinkRecords ChainlinkRecords,
	dexRecords DexRecords, height uint64) *gg {
	var runningGraph = gograph.NewGraph[common.Address, int64, string, interface{}](true)
	callopts := bind.CallOpts{BlockNumber: big.NewInt(int64(height))}

	maxCnt := 100
	mut := sync.Mutex{}

	n.log.Info("syncing info for chainlink records")
	wg := sync.WaitGroup{}

	for _, _rec := range chainlinkRecords {
		if uint64(_rec.StartBlock) > height {
			break
		}
		wg.Add(1)
		go func(rec ChainlinkRecord) {
			defer wg.Done()
			oracleMetadata, err := n.EthRPC.GetChainlinkRoundData(rec.Oracle, &callopts)
			if err != nil {
				n.log.Fatal("cannot retrieve metadata for cl oracle, skipping",
					"oracle", rec.Oracle,
					"height", height)
			}
			mut.Lock()
			defer mut.Unlock()
			runningGraph.AddWeightedEdge(rec.From,
				rec.To,
				math.MaxInt64,
				"chainlink",
				oracleMetadata)
		}(_rec)
	}
	wg.Wait()
	n.log.Info("syncing info for chainlink records", "edges", runningGraph.GetEdgeCount())

	n.log.Info("syncing info for dex records")
	for _, _rec := range dexRecords {
		maxCnt--
		if maxCnt == 0 {
			break
		}
		if uint64(_rec.StartBlock) > height {
			break
		}
		wg.Add(1)
		go func(rec DexRecord) {
			defer wg.Done()
			t0, t1, err := n.EthRPC.GetTokensUniV2(rec.Pair, &callopts)
			if err != nil {
				n.log.Warn("not a univ2 pair, skipping",
					"pair", rec.Pair)
				return
			}
			if t0 != rec.Token0 || t1 != rec.Token1 {
				n.log.Warn("token mismatch, skipping",
					"pair", rec.Pair,
					"expected0", rec.Token0,
					"expected1", rec.Token1,
					"rpc0", t0,
					"rpc1", t1,
				)
				return
			}
			res0, err := n.EthRPC.GetERC20Balances([]itypes.Tuple2[common.Address, common.Address]{
				{rec.Pair, t0}, {rec.Pair, t1},
			}, &callopts)
			if err != nil {
				n.log.Warn("cannot retrieve balances for pool", "pair", rec.Pair)
				return
			}
			if len(res0) != 2 {
				panic(len(res0))
			}

			mut.Lock()
			defer mut.Unlock()
			runningGraph.AddWeightedEdge(rec.Token0,
				rec.Token1,
				1, // fetch
				"dex",
				itypes.UniV2Metadata{res0[0], res0[1]})
		}(_rec)
	}
	wg.Wait()
	n.log.Info("syncing info for dex records", "edges", runningGraph.GetEdgeCount())

	return runningGraph
}
