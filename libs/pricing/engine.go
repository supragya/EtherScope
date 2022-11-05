package priceresolver

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/gograph"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/ethrpc"
	localbackend "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	"github.com/ethereum/go-ethereum/common"
)

// Enhanced cached, multistep, graph based batch pricing resolver system
type Engine struct {
	log          logger.Logger
	EthRPC       ethrpc.EthRPC
	LocalBackend localbackend.LocalBackend

	// On-disk dump information
	chainlinkOraclesDumpFile string
	dexDumpFile              string

	// Internal Data Structures
	lastGraph  *gograph.Graph[common.Address, interface{}]
	lastHeight uint64
}

// DefaultEngine is default form of enhanced pricing engine
func NewDefaultEngine(log logger.Logger,
	chainlinkOracledDumpFile string,
	dexDumpFile string,
	ethrpcBackend ethrpc.EthRPC,
	localBackend localbackend.LocalBackend) *Engine {
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
		_, ok, err := n.LocalBackend.Get("latestPricingGraph")
		if err != nil {
			return err
		}
		if !ok {
			// This means localbackend does not know any graph
			// just yet. Go through dumps
			n.log.Warn("no pricingGraph found in localbackend. reading dump files")
			neededGraph, err := n.syncDump(resHeight)
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

func (n *Engine) syncDump(resHeight uint64) (*gograph.Graph[common.Address, interface{}], error) {
	chainlinkRecords, err := loadChainlinkCSV(n.chainlinkOraclesDumpFile)
	if err != nil {
		return nil, err
	}
	dexRecords, err := loadDexCSV(n.dexDumpFile)
	if err != nil {
		return nil, err
	}

	graphSteps := genGraphSteps(chainlinkRecords, dexRecords)
	neededHeight := graphSteps[len(graphSteps)-1].First
	neededGraph := graphSteps[len(graphSteps)-1].Second

	idx := len(graphSteps) - 1
	for idx > 0 && neededHeight > resHeight {
		idx--
		neededHeight = graphSteps[idx].First
		neededGraph = graphSteps[idx].Second
	}

	err = n.syncGraphStepsToLB(graphSteps)
	if err != nil {
		return nil, err
	}
	return neededGraph, nil
}
