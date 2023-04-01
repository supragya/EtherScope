package node

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/common"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/pricing"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/rpc"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/service"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/sink"
	types "github.com/Blockpour/Blockpour-Geth-Indexer/algorand/types"
	util "github.com/Blockpour/Blockpour-Geth-Indexer/algorand/util"
	"github.com/google/uuid"
)

type Node struct {
	service.BaseService

	log     logger.Logger
	rpc     *rpc.AlgoRPC
	sink    sink.OutputSink
	pricing *pricing.PricingEngine

	//Configs
	startBlock          uint64
	skipResume          bool
	maxBlockSpanPerCall uint64
	resumeURL           string

	// Internal
	currentHeight uint64
	indexedHeight uint64
	quit          chan struct{}

	nodeID  uuid.UUID
	moniker string
}

func NewNode(
	startBlock uint64,
	algodUrl string,
	indexerUrl string,
	rpcToken string,
	maxBlockSpanPerCall uint64,
	skipResume bool,
	resumeURL string,
	log logger.Logger,
) (*Node, error) {
	rpc, err := rpc.NewAlgoRPC(algodUrl, indexerUrl, rpcToken)
	if err != nil {
		return nil, err
	}

	pe := pricing.NewPricingEngine(rpc)

	sink, err := sink.NewRabbitMQOutputSinkWithViperFields(log.With("service", "outputsink"))
	if err != nil {
		return nil, err
	}

	n := &Node{
		currentHeight:       0,
		startBlock:          startBlock,
		indexedHeight:       startBlock,
		maxBlockSpanPerCall: maxBlockSpanPerCall,
		skipResume:          skipResume,
		resumeURL:           resumeURL,
		log:                 log,
		pricing:             pe,
		rpc:                 rpc,
		sink:                sink,
		nodeID:              uuid.New(),
		moniker:             "algorand-indexer",
		quit:                make(chan struct{}, 1),
	}

	n.BaseService = *service.NewBaseService(log, "Indexer", n)
	return n, nil
}

func (n *Node) OnStart(ctx context.Context) error {

	err := n.sink.Start(ctx)
	if err != nil {
		panic(err)
	}

	n.indexedHeight = n.syncStartHeight()
	go n.loop()

	return nil
}

func (n *Node) OnStop() {
	n.quit <- struct{}{}
}

func (n *Node) loop() {
	for {
		select {
		case <-time.After(time.Second * 2):
			// Loop in case we are lagging, so we dont wait 3 secs between epochs
			for {
				height, err := n.rpc.GetCurrentBlockHeight()
				if err != nil {
					n.log.Error("error getting current block height", err)
					break
				}

				n.currentHeight = height

				if n.currentHeight == n.indexedHeight {
					continue
				}

				endHeight := n.currentHeight
				isOnHead := true
				if (endHeight - n.indexedHeight) > n.maxBlockSpanPerCall {
					isOnHead = false
					endHeight = n.indexedHeight + n.maxBlockSpanPerCall
				}

				/*
					The algod sdk and indexer sdk clients are sometimes not in sync. In particular, the last
					block that the algod sdk returns is not always processed by the indexer sdk.
					The block will be processed in the next loop.
				*/
				_, err = n.rpc.GetBlock(endHeight)
				if err != nil {
					n.log.Error("error looking up end height block", err)
					break
				}

				n.log.Info(fmt.Sprintf("chainhead: %d (+%d away), indexed: %d",
					n.currentHeight, n.currentHeight-n.indexedHeight, n.indexedHeight))

				sigs := common.SupportedFunctionSignatures

				f := types.TxFilter{
					StartRound: n.indexedHeight + 1,
					EndRound:   endHeight,
					Signatures: sigs,
				}

				n.processBlocks(f)
				n.indexedHeight = endHeight

				if isOnHead {
					break
				}
			}
		case <-n.quit:
			n.log.Info("quitting realtime indexer")
		}
	}
}

func (n *Node) processBlocks(f types.TxFilter) {
	startTime := time.Now()
	swaps, err := n.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	processingTime := time.Now()

	for r := f.StartRound; r <= f.EndRound; r++ {
		blockSwaps := util.Filter(swaps, func(swap types.Swap) bool {
			return swap.Height == r
		})

		time, err := n.rpc.GetBlockTimestamp(r)
		if err != nil {
			n.log.Error("error getting block timestamp", err)

			// If we cant get the timestamp from rpc, we get the timestamp from the first swap in the block
			if len(blockSwaps) > 0 {
				time = blockSwaps[0].Time
			} else {
				// Otherwise if there are no swaps in the block, we set the time to 0
				time = 0
			}
		}

		meta := types.BlockSynopsis{
			Type:                    "stats",
			Network:                 99990,
			Height:                  r,
			Time:                    time,
			TotalLogs:               uint64(len(blockSwaps)),
			SwapLogs:                uint64(len(blockSwaps)),
			BurnLogs:                0,
			MintLogs:                0,
			IndexingTimeNanos:       uint64(processingTime.UnixNano()),
			ProcessingDurationNanos: uint64(processingTime.Sub(startTime).Nanoseconds()),
		}

		p := types.Payload{
			NodeMoniker:   n.moniker,
			NodeID:        n.nodeID,
			Environment:   "production",
			NodeVersion:   "1.0.0",
			Network:       "algorand",
			BlockSynopsis: &meta,
		}

		p.Add(meta)

		for _, swap := range blockSwaps {
			p.Add(swap)
		}

		err = n.sink.Send(p)
		if err != nil {
			n.log.Error("error sending payload", err)
		}

	}
}

func (n *Node) Status() interface{} {
	return nil
}

func (n *Node) syncStartHeight() uint64 {
	startBlock := n.indexedHeight

	if !n.skipResume {
		remoteLatestHeight, err := n.getRemoteLatestHeight()
		if err != nil {
			n.log.Fatal(fmt.Sprintf("error while fetching latest height from remote: %v", err))
		}
		if remoteLatestHeight < startBlock {
			n.log.Fatal("remote latest height is less than start block height")
		}
		if remoteLatestHeight > startBlock {
			startBlock = remoteLatestHeight
		}
	}

	n.log.Info("start block height set", "start", startBlock)
	return startBlock
}

func (n *Node) getRemoteLatestHeight() (uint64, error) {
	resp, err := http.Get(n.resumeURL)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var r types.ResumeAPIResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, err
	}

	n.log.Info("resuming from block height: ", r.Data.Height)
	return r.Data.Height, nil
}

func (n *Node) ProcessTxns(f types.TxFilter) ([]types.Swap, error) {
	swaps := []types.Swap{}

	for h := f.StartRound; h <= f.EndRound; h++ {
		n.log.Info(fmt.Sprintf("processing block: %d", h))

		txGroups, err := n.rpc.GetTxGroups(h, f.Signatures)
		if err != nil {
			fmt.Println(err)
			return []types.Swap{}, err
		}

		wg := &sync.WaitGroup{}
		mt := &sync.Mutex{}

		// there can be several swaps within the same tx group so we
		// need to use an ordered array to store the results when doing concurrent processing
		// if we want to conserve tx group ordering.
		var results types.OrderedSwaps = make(types.OrderedSwaps, len(txGroups))

		for idx, g := range txGroups {
			wg.Add(1)
			switch g.FunctionSignature {
			case common.TinymanV1SwapSignature:
				go n.processTinymanV1Swap(g, idx, results, wg, mt)
			case common.TinymanV2SwapSignature:
				go n.processTinymanV2Swap(g, idx, results, wg, mt)
			case common.AlgoFiSwapSignatures[0]:
				go n.processAlgofiSwap(g, idx, results, wg, mt)
			case common.AlgoFiSwapSignatures[1]:
				go n.processAlgofiSwap(g, idx, results, wg, mt)
			}
		}

		wg.Wait()

		// logic is different from evm indexer as a tx group can contain several txs
		for _, ss := range results {
			swaps = append(swaps, ss...)
		}

	}

	return swaps, nil
}
