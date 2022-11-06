package priceresolver

import (
	"fmt"
	"math"

	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/gograph"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/ethrpc"
	lb "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
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
	// Check if we even have latest graph to update
	// and price using that
	if n.lastHeight+1 != resHeight {
		// Check what height localBackend has stored
		_, ok, err := n.LocalBackend.Get(lb.KeyLatestPricingGraph)
		if err != nil {
			return err
		}
		if !ok {
			// This means localbackend does not know any graph
			// just yet. Go through dumps
			n.log.Warn("no pricingGraph found in localbackend. reading dump files")
			_, err := n.syncDump(resHeight)
			if err != nil {
				return err
			}
		}
	}

	// We know for certain at this point,
	// n.lastGraph is not nil and
	// n.lastHeight + 1 == resHeight
	return nil
}

func (n *Engine) syncDump(resHeight uint64) (*gg, error) {
	chainlinkRecords, err := loadChainlinkCSV(n.chainlinkOraclesDumpFile)
	if err != nil {
		return nil, err
	}
	dexRecords, err := loadDexCSV(n.dexDumpFile)
	if err != nil {
		return nil, err
	}

	graphSteps := genGraphSteps(chainlinkRecords, dexRecords)
	gHeight := graphSteps[len(graphSteps)-1].First
	gGraph := graphSteps[len(graphSteps)-1].Second

	idx := len(graphSteps) - 1
	for idx > 0 && gHeight > resHeight {
		idx--
		gHeight = graphSteps[idx].First
		gGraph = graphSteps[idx].Second
	}

	err = n.syncGraphStepsToLB(graphSteps)
	if err != nil {
		return nil, err
	}
	return gGraph, nil
}

func (n *Engine) syncGraphStepsToLB(steps []ggt) error {
	ggtLen := len(steps)
	if ggtLen == 0 {
		return nil
	}

	lowestHeight := steps[0].First
	highestHeight := steps[ggtLen-1].First

	if err := n.LocalBackend.Set(lb.KeyLowestPricingHeight, lowestHeight); err != nil {
		return err
	}

	gIdx := 0
	for idx := lowestHeight; idx <= highestHeight; idx++ {
		prev := steps[gIdx].First
		next := uint64(math.MaxUint64)
		if gIdx != ggtLen-1 {
			next = steps[gIdx+1].First
		}
		if idx == prev {
			if err := n.LocalBackend.Set(lb.KeyGraphPrefix+fmt.Sprint(idx), steps[gIdx]); err != nil {
				return err
			}
		} else if prev < idx && idx < next {
			if err := n.LocalBackend.Set(lb.KeyGraphPrefix+fmt.Sprint(idx), prev); err != nil {
				return err
			}
		} else if idx == next {
			if err := n.LocalBackend.Set(lb.KeyGraphPrefix+fmt.Sprint(idx), steps[gIdx+1]); err != nil {
				return err
			}
			gIdx++
		}
	}

	return n.LocalBackend.Sync()
}

// generates independent graphs for blocks where
// new chainlink oracles came alive or dex pools were
// made. Ideally the steps should be sparse enough
// (not every block has dex creation or oracle setup)
// to justify memory footprint.
func genGraphSteps(chainlinkRecords ChainlinkRecords,
	dexRecords DexRecords) []ggt {
	var (
		crec          *ChainlinkRecord = nil
		drec          *DexRecord       = nil
		cidx, clen    int              = 0, chainlinkRecords.Len()
		didx, dlen    int              = 0, dexRecords.Len()
		steps                          = []ggt{}
		runningHeight uint64           = 0
		runningGraph                   = gograph.NewGraph[common.Address, int64, string, interface{}](true)
	)

	for {
		crec, drec = nil, nil

		if cidx < clen {
			crec = &chainlinkRecords[cidx]
		}
		if didx < dlen {
			drec = &dexRecords[didx]
		}

		// At this point, cidx and crec are index and pointer to
		// chainlink record candidate and didx and dlen and drec
		// are dex record candidate
		var (
			candidateStartBlock int64 = -1
			edgeToAdd                 = we{}
		)

		if crec != nil {
			candidateStartBlock = crec.StartBlock
			edgeToAdd = we{
				VertexFrom: crec.From,
				VertexTo:   crec.To,
				Weight:     9223372036854775807, // Max int64
				Hint:       "chainlink",
				Metadata:   crec.Oracle,
			}
			cidx++
		}
		if drec != nil && candidateStartBlock > drec.StartBlock {
			if candidateStartBlock == -1 {
				cidx--
			}
			candidateStartBlock = drec.StartBlock
			edgeToAdd = we{
				VertexFrom: drec.Token0,
				VertexTo:   drec.Token1,
				Weight:     1, // to fetch
				Hint:       "dex",
				Metadata:   drec.Pair,
			}
			didx++
		}

		// At this point, either candidateStartBlock == -1 which means
		// there is no candidate to add (both chainlink edges and dex edges
		// have been exhaused) or candidateBlock != -1 (positive) so, we need to add
		// an edge to appropriate graph
		if candidateStartBlock == -1 {
			// All done
			break
		}

		if candidateStartBlock < 0 {
			panic(fmt.Sprintf("candiateStartBlock(%v) lower than 0", candidateStartBlock))
		}

		csb := uint64(candidateStartBlock)

		if csb < runningHeight {
			panic(fmt.Sprintf("candidateStartBlock(%v) lower than runningHeight(%v). Are the stacks not sorted?",
				candidateStartBlock, runningHeight))
		}

		if csb > runningHeight {
			steps = append(steps, ggt{runningHeight, runningGraph})
			runningGraph = gograph.CopyGraph(runningGraph)
			runningHeight = csb
		}

		runningGraph.AddWeightedEdge(edgeToAdd.VertexFrom,
			edgeToAdd.VertexTo,
			edgeToAdd.Weight,
			edgeToAdd.Hint,
			edgeToAdd.Metadata)
	}

	return steps
}
