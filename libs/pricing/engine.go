package priceresolver

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/supragya/EtherScope/libs/gograph"
	logger "github.com/supragya/EtherScope/libs/log"
	"github.com/supragya/EtherScope/libs/util"
	"github.com/supragya/EtherScope/services/ethrpc"
	lb "github.com/supragya/EtherScope/services/local_backend"
	itypes "github.com/supragya/EtherScope/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type we = gograph.WeightedEdge[common.Address, int64, string, interface{}]
type gg = gograph.Graph[common.Address, int64, string, interface{}]
type ggt = itypes.Tuple2[uint64, *gg]
type addrTuple = itypes.Tuple2[common.Address, common.Address]
type resTuple = itypes.Tuple2[*big.Float, *big.Float]

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

var (
	ErrorRequestedResolutionPresent = errors.New("RequestedResolutionPresentError")
)

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

func (n *Engine) FetchLatestBlockHeightOrConstructTill(resHeight uint64) uint64 {
	// fetch the graph in question
PRICING_GRAPH_FETCH:
	lbLatestHeight, ok, err := n.LocalBackend.Get(lb.KeyLatestHeight)
	if err != nil {
		n.log.Fatal(err.Error())
	}
	if !ok {
		// This means localbackend does not know any graph
		// just yet. Go through dumps
		n.log.Warn("no pricingGraph found in localbackend. reading dump files and syncing")
		err := n.syncDump(resHeight)
		if err != nil {
			n.log.Fatal(err.Error())
		}
		goto PRICING_GRAPH_FETCH
	}

	// ensure lbLatestHeight is less than resHeight
	var lbLatestHeightUint64 uint64
	if err := util.GobDecode(lbLatestHeight, &lbLatestHeightUint64); err != nil {
		n.log.Fatal(fmt.Sprintf("wrong latest height encoding: %s", err))
	}
	return lbLatestHeightUint64
}

func (n *Engine) GetLatestGraph() (*gg, error) {
	graph := gograph.NewGraph[common.Address, int64, string, interface{}](true)
	lbLatestGraph, ok, err := n.LocalBackend.Get(lb.KeyLatestPricingGraph)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("cannot get latest pricing graph")
	}
	if err := util.GobDecode(lbLatestGraph, graph); err != nil {
		return nil, err
	}
	return graph, nil
}

func (n *Engine) Resolve(resHeight uint64, items []interface{}) ([]itypes.UniV2Metadata, error) {
	// Ensure entries exist and we are resolving
	// for heights >= height in localbackend
	fetchedHeight := n.FetchLatestBlockHeightOrConstructTill(resHeight)
	if resHeight < fetchedHeight {
		return []itypes.UniV2Metadata{}, fmt.Errorf("requested resolution for %d while lb already has %d: %w",
			resHeight,
			fetchedHeight,
			ErrorRequestedResolutionPresent)
	}

	// Get the graph from localbackend
	graph, err := n.GetLatestGraph()
	if err != nil {
		return []itypes.UniV2Metadata{}, err
	}

	// Get updates to be applied to graph
	resUpdates := n.getReserveUpdates(items)

	// update graph
	newDexes := n.updateGraph(graph, resUpdates, resHeight)

	// Resolve items: best attempt
	n.resolveItems(graph, items, resHeight)

	// Update localbackend
	n.syncGraphToLB(graph, resHeight)

	return newDexes, nil
}

func (n *Engine) resolveItems(graph *gg, items []interface{}, resHeight uint64) {
	// For each item that is UserRequested
	// Do bfs to find the best 3 candidates
	requested, priced, cpriced := 0, 0, 0
	tc := make(map[common.Address]*itypes.PriceResult, len(items))

	for idx, item := range items {
		switch i := item.(type) {
		case *itypes.Mint:
			if i.ProcessingType != itypes.UserRequested {
				continue
			}
			requested++
			if n.tryPricingUSD(i.Token0, &i.Price0, graph, tc, resHeight) {
				priced++
			}
			requested++
			if n.tryPricingUSD(i.Token1, &i.Price1, graph, tc, resHeight) {
				priced++
			}
			if i.Price0 != nil || i.Price1 != nil {
				if i.Price0 == nil {
					p := big.NewFloat(1.0).Set(i.Price1.Price)
					p = p.Mul(p, i.Amount1)
					p = p.Abs(p)
					counterparty := big.NewFloat(1.0).Set(p)
					counterparty = counterparty.Quo(counterparty, i.Amount0)
					counterparty = counterparty.Abs(counterparty)
					i.Price0 = &itypes.PriceResult{Price: counterparty,
						Path: []interface{}{
							itypes.CounterPartyResolutionMetadata{
								Description: "Counterparty resolution",
								Price:       counterparty,
							},
						},
					}
					cpriced++
					priced++
					i.AmountUSD = p
				} else if i.Price1 == nil {
					p := big.NewFloat(1.0).Set(i.Price0.Price)
					p = p.Mul(p, i.Amount0)
					p = p.Abs(p)
					counterparty := big.NewFloat(1.0).Set(p)
					counterparty = counterparty.Quo(counterparty, i.Amount1)
					counterparty = counterparty.Abs(counterparty)
					i.Price1 = &itypes.PriceResult{Price: counterparty,
						Path: []interface{}{
							itypes.CounterPartyResolutionMetadata{
								Description: "Counterparty resolution",
								Price:       counterparty,
							},
						},
					}
					cpriced++
					priced++
					i.AmountUSD = p
				} else {
					p1 := big.NewFloat(1.0).Set(i.Price0.Price)
					p1 = p1.Mul(p1, i.Amount0)
					p1 = p1.Abs(p1)
					i.AmountUSD = p1
				}
			}
			items[idx] = i
		case *itypes.Burn:
			if i.ProcessingType != itypes.UserRequested {
				continue
			}
			requested++
			if n.tryPricingUSD(i.Token0, &i.Price0, graph, tc, resHeight) {
				priced++
			}
			requested++
			if n.tryPricingUSD(i.Token1, &i.Price1, graph, tc, resHeight) {
				priced++
			}
			if i.Price0 != nil || i.Price1 != nil {
				if i.Price0 == nil {
					p := big.NewFloat(1.0).Set(i.Price1.Price)
					p = p.Mul(p, i.Amount1)
					p = p.Abs(p)
					counterparty := big.NewFloat(1.0).Set(p)
					counterparty = counterparty.Quo(counterparty, i.Amount0)
					counterparty = counterparty.Abs(counterparty)
					i.Price0 = &itypes.PriceResult{Price: counterparty,
						Path: []interface{}{
							itypes.CounterPartyResolutionMetadata{
								Description: "Counterparty resolution",
								Price:       counterparty,
							},
						},
					}
					cpriced++
					priced++
					i.AmountUSD = p
				} else if i.Price1 == nil {
					p := big.NewFloat(1.0).Set(i.Price0.Price)
					p = p.Mul(p, i.Amount0)
					p = p.Abs(p)
					counterparty := big.NewFloat(1.0).Set(p)
					counterparty = counterparty.Quo(counterparty, i.Amount1)
					counterparty = counterparty.Abs(counterparty)
					i.Price1 = &itypes.PriceResult{Price: counterparty,
						Path: []interface{}{
							itypes.CounterPartyResolutionMetadata{
								Description: "Counterparty resolution",
								Price:       counterparty,
							},
						},
					}
					cpriced++
					priced++
					i.AmountUSD = p
				} else {
					p1 := big.NewFloat(1.0).Set(i.Price0.Price)
					p1 = p1.Mul(p1, i.Amount0)
					p1 = p1.Abs(p1)
					i.AmountUSD = p1
				}
			}
			items[idx] = i
		case *itypes.Swap:
			if i.ProcessingType != itypes.UserRequested {
				continue
			}
			requested++
			if n.tryPricingUSD(i.Token0, &i.Price0, graph, tc, resHeight) {
				priced++
			}
			requested++
			if n.tryPricingUSD(i.Token1, &i.Price1, graph, tc, resHeight) {
				priced++
			}
			if i.Price0 != nil || i.Price1 != nil {
				if i.Price0 == nil {
					p := big.NewFloat(1.0).Set(i.Price1.Price)
					p = p.Mul(p, i.Amount1)
					p = p.Abs(p)
					counterparty := big.NewFloat(1.0).Set(p)
					counterparty = counterparty.Quo(counterparty, i.Amount0)
					counterparty = counterparty.Abs(counterparty)
					i.Price0 = &itypes.PriceResult{Price: counterparty,
						Path: []interface{}{
							itypes.CounterPartyResolutionMetadata{
								Description: "Counterparty resolution",
								Price:       counterparty,
							},
						},
					}
					cpriced++
					priced++
					i.AmountUSD = p
				} else if i.Price1 == nil {
					p := big.NewFloat(1.0).Set(i.Price0.Price)
					p = p.Mul(p, i.Amount0)
					p = p.Abs(p)
					counterparty := big.NewFloat(1.0).Set(p)
					counterparty = counterparty.Quo(counterparty, i.Amount1)
					counterparty = counterparty.Abs(counterparty)
					i.Price1 = &itypes.PriceResult{Price: counterparty,
						Path: []interface{}{
							itypes.CounterPartyResolutionMetadata{
								Description: "Counterparty resolution",
								Price:       counterparty,
							},
						},
					}
					cpriced++
					priced++
					i.AmountUSD = p
				} else {
					p1 := big.NewFloat(1.0).Set(i.Price0.Price)
					p1 = p1.Mul(p1, i.Amount0)
					p1 = p1.Abs(p1)
					i.AmountUSD = p1
				}
			}
			items[idx] = i
		}
	}
	n.log.Info("pricing engine resolution statistics",
		"height", resHeight,
		"requested", requested,
		"priced", priced,
		"counterpartypriced", cpriced)
}

func (n *Engine) getReserveUpdates(items []interface{}) map[addrTuple]itypes.UniV2Metadata {
	reserveUpdates := make(map[addrTuple]itypes.UniV2Metadata, len(items))

	for _, item := range items {
		switch i := item.(type) {
		case *itypes.Mint:
			reserveUpdates[addrTuple{i.Token0, i.Token1}] = itypes.UniV2Metadata{"", i.PairContract, i.Token0, i.Token1, i.Reserve0, i.Reserve1}
		case *itypes.Burn:
			reserveUpdates[addrTuple{i.Token0, i.Token1}] = itypes.UniV2Metadata{"", i.PairContract, i.Token0, i.Token1, i.Reserve0, i.Reserve1}
		case *itypes.Swap:
			reserveUpdates[addrTuple{i.Token0, i.Token1}] = itypes.UniV2Metadata{"", i.PairContract, i.Token0, i.Token1, i.Reserve0, i.Reserve1}
		}
	}

	return reserveUpdates
}

func (n *Engine) updateGraph(graph *gg, updates map[addrTuple]itypes.UniV2Metadata, resHeight uint64) []itypes.UniV2Metadata {
	wg := sync.WaitGroup{}
	callopts := bind.CallOpts{BlockNumber: big.NewInt(int64(resHeight))}
	newDexes := []itypes.UniV2Metadata{}
	// nullAddr := common.Address{}

	// Ensure swaps exist
	for addrs, val := range updates {
		// expecting bidirectional graph
		if _, ok := graph.Graph[addrs.First]; !ok {
			_, _, err := n.EthRPC.GetTokensUniV2(val.Pair, &callopts)
			if err != nil {
				n.log.Warn("not a univ2 pair, skipping",
					"pair", val.Pair)
				break
			}

			name0, err := n.EthRPC.GetERC20Name(val.Token0, &callopts)
			if err != nil {
				n.log.Warn("unable to get name for token", "error", err, "token", val.Token0)
			}
			name1, err := n.EthRPC.GetERC20Name(val.Token1, &callopts)
			if err != nil {
				n.log.Warn("unable to get name for token", "error", err, "token", val.Token1)
			}
			dex := fmt.Sprintf("(%v, %v) %v", name0, name1, strings.ToLower(val.Pair.Hex()))
			n.log.Info("adding previously unseen dex to pricing graph", "dex", dex)
			newDexes = append(newDexes, val)

			graph.AddWeightedEdge(val.Token0,
				val.Token1,
				1, // fetch
				"dex",
				val)
		}
	}

	mut := sync.Mutex{}
	for from, connections := range graph.Graph {
		for to, edge := range connections {
			switch i := edge.Metadata.(type) {

			case itypes.WrappedCLMetadata:
				wg.Add(1)
				go func(from common.Address, to common.Address, i itypes.WrappedCLMetadata, edge we) {
					defer wg.Done()
					oracleMetadata, err := n.EthRPC.GetChainlinkRoundData(i.Oracle, &callopts)
					if err != nil {
						n.log.Warn("cannot retrieve metadata for cl oracle, skipping",
							"oracle", i.Oracle,
							"height", resHeight)
						return
					}
					i.Data = oracleMetadata
					edge.Metadata = i
					mut.Lock()
					graph.Graph[from][to] = edge
					mut.Unlock()
				}(from, to, i, edge)

			case itypes.UniV2Metadata:
				if update, ok := updates[addrTuple{from, to}]; ok {
					i.Pair = update.Pair
					i.Res0 = update.Res0
					i.Res1 = update.Res1
					edge.Metadata = i
					mut.Lock()
					graph.Graph[from][to] = edge
					mut.Unlock()
				}

			default:
				panic(fmt.Sprintf("unknown type detected: %v", i))
			}
		}
	}
	wg.Wait()

	return newDexes
}

// best effort
func (n *Engine) tryPricingUSD(from common.Address,
	result **itypes.PriceResult,
	graph *gg,
	tc map[common.Address]*itypes.PriceResult,
	resHeight uint64) bool {
	// check cache
	if v, ok := tc[from]; ok {
		*result = v
		return true
	}

	// find cadidates
	maxRoutes := 5
	routes := graph.GetBFSCandidates(maxRoutes,
		from, common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff"))

	callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(resHeight))}

	maxScore := big.NewFloat(0.0)
	var calcResult *itypes.PriceResult = nil

	for _, route := range routes {
		// n.log.Info("route for", "from", from, "route", route)

		multiplier := big.NewFloat(1.0)
		pr := itypes.PriceResult{}
		minScore := big.NewFloat(math.MaxInt64)

		for _, edge := range route {
			switch i := edge.Metadata.(type) {
			case itypes.WrappedCLMetadata:
				name0, err := n.EthRPC.GetERC20Name(i.From, callopts)
				if err != nil {
					n.log.Warn("unable to get name for token", "error", err, "token", i.From)
				}
				name1, err := n.EthRPC.GetERC20Name(i.To, callopts)
				if err != nil {
					n.log.Warn("unable to get name for token", "error", err, "token", i.To)
				}
				i.Description = fmt.Sprintf("Chainlink (%v, %v), rev:%v", name0, name1, edge.IsReverseEdge)
				pr.Path = append(pr.Path, i)
				// TODO: error checks here
				decimals, _ := n.EthRPC.GetERC20Decimals(i.Oracle, callopts)
				if !edge.IsReverseEdge {
					multiplier = multiplier.Mul(multiplier, util.DivideBy10pow(i.Data.Answer, decimals))
				} else {
					multiplier = multiplier.Quo(multiplier, util.DivideBy10pow(i.Data.Answer, decimals))
				}
			case itypes.UniV2Metadata:
				name0, err := n.EthRPC.GetERC20Name(i.Token0, callopts)
				if err != nil {
					n.log.Warn("unable to get name for token", "error", err, "token", i.Token0)
				}
				name1, err := n.EthRPC.GetERC20Name(i.Token1, callopts)
				if err != nil {
					n.log.Warn("unable to get name for token", "error", err, "token", i.Token1)
				}
				i.Description = fmt.Sprintf("UniswapV2Dex (%v, %v), rev:%v", name0, name1, edge.IsReverseEdge)

				pr.Path = append(pr.Path, i)
				ratio := big.NewFloat(1.0).Quo(i.Res1, i.Res0)
				if !edge.IsReverseEdge {
					edgeImpact := big.NewFloat(1.0).Mul(i.Res0, multiplier)
					// edgeImpact = edgeImpact * 2.0 // Implicit
					multiplier = multiplier.Mul(multiplier, ratio)
					if edgeImpact.Cmp(minScore) == -1 {
						minScore.Set(edgeImpact)
					}
				} else {
					edgeImpact := big.NewFloat(1.0).Mul(i.Res1, multiplier)
					multiplier = multiplier.Quo(multiplier, ratio)
					if edgeImpact.Cmp(minScore) == -1 {
						minScore.Set(edgeImpact)
					}
				}
			}
		}
		pr.Price = multiplier

		if maxScore.Cmp(minScore) == -1 {
			maxScore.Set(minScore)
			calcResult = &pr
		}

		// if from == common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f") {
		// 	n.log.Info("price found", "price", pr.Price, "path", pr.Path, "maxScore", maxScore, "minScore", minScore)
		// }
	}

	// Cache result
	if calcResult != nil {
		// n.log.Info("caching result", "token", from, "price", calcResult)
		tc[from] = calcResult
		*result = calcResult
		return true
	}
	return false
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
				itypes.WrappedCLMetadata{
					Data:   oracleMetadata,
					Oracle: rec.Oracle,
					From:   rec.From,
					To:     rec.To,
				})
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
			res, err := n.EthRPC.GetERC20Balances([]itypes.Tuple2[common.Address, common.Address]{
				{rec.Pair, t0}, {rec.Pair, t1},
			}, &callopts)
			if err != nil {
				n.log.Warn("cannot retrieve balances for pool", "pair", rec.Pair)
				return
			}
			if len(res) != 2 {
				panic(len(res))
			}

			t0d, err := n.EthRPC.GetERC20Decimals(t0, &callopts)
			if err != nil {
				return
			}
			t1d, err := n.EthRPC.GetERC20Decimals(t1, &callopts)
			if err != nil {
				return
			}

			mut.Lock()
			defer mut.Unlock()
			runningGraph.AddWeightedEdge(rec.Token0,
				rec.Token1,
				1, // fetch
				"dex",
				itypes.UniV2Metadata{
					Pair:   rec.Pair,
					Token0: rec.Token0,
					Token1: rec.Token1,
					Res0:   util.DivideBy10pow(res[0], t0d),
					Res1:   util.DivideBy10pow(res[1], t1d),
				})
		}(_rec)
	}
	wg.Wait()
	n.log.Info("syncing info for dex records", "edges", runningGraph.GetEdgeCount())

	return runningGraph
}
