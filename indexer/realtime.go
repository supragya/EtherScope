package indexer

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	"github.com/Blockpour/Blockpour-Geth-Indexer/ethrpc"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/instrumentation"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RealtimeIndexer struct {
	currentHeight    uint64
	indexedHeight    uint64
	dbconn           *db.DBConn
	da               ethrpc.EthRPC
	eventsToIndex    []common.Hash
	eventsToIndexStr []string

	quitCh chan struct{}
}

func NewRealtimeIndexer(indexedHeight uint64,
	masterUpstream string,
	slaveUpstreams []string,
	timeout time.Duration,
	isErigon bool,
	dbconn *db.DBConn,
	eventsToIndex []string) *RealtimeIndexer {
	events, err := util.ConstructTopics(eventsToIndex)
	util.ENOK(err)
	return &RealtimeIndexer{
		currentHeight:    0,
		indexedHeight:    indexedHeight,
		dbconn:           dbconn,
		da:               *ethrpc.NewEthRPC(isErigon, masterUpstream, slaveUpstreams, timeout),
		eventsToIndex:    events,
		eventsToIndexStr: eventsToIndex,

		quitCh: make(chan struct{}),
	}
}

func (r *RealtimeIndexer) Start() error {
	if r.da.Len() == 0 {
		return EUninitialized
	}
	for i := 0; i < len(r.eventsToIndex); i++ {
		log.Info("enabled: ", r.eventsToIndexStr[i], "(", r.eventsToIndex[i], ")")
	}
	r.ridxLoop()
	return nil
}

func (r *RealtimeIndexer) ridxLoop() {
	instrumentation.StartingBlock.Set(float64(r.indexedHeight))
	instrumentation.ProcessedBlock.Set(float64(r.indexedHeight))

	maxBlockSpanPerCall := viper.GetUint64("general.maxBlockSpanPerCall")
	for {
		select {
		case <-time.After(time.Second * 2):
			// Loop in case we are lagging, so we dont wait 3 secs between epochs
			for {
				height, err := r.da.GetCurrentBlockHeight()

				util.ENOK(err)
				r.currentHeight = height

				if r.currentHeight == r.indexedHeight {
					continue
				}
				endingBlock := r.currentHeight
				isOnHead := true
				if (endingBlock - r.indexedHeight) > maxBlockSpanPerCall {
					isOnHead = false
					endingBlock = r.indexedHeight + maxBlockSpanPerCall
				}

				log.Info(fmt.Sprintf("sync curr: %d (+%d), processing [%d - %d]",
					r.currentHeight, r.currentHeight-r.indexedHeight, r.indexedHeight, endingBlock))

				instrumentation.CurrentBlock.Set(float64(r.currentHeight))

				logs, err := r.da.GetFilteredLogs(ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(r.indexedHeight + 1)),
					ToBlock:   big.NewInt(int64(endingBlock)),
					Topics:    [][]common.Hash{r.eventsToIndex},
				})

				if err != nil {
					log.Error(err)
					continue
				}

				r.processBatchedBlockLogs(logs, r.indexedHeight+1, endingBlock)

				r.indexedHeight = endingBlock
				instrumentation.ProcessedBlock.Set(float64(r.indexedHeight))

				if isOnHead {
					break
				}
			}
		case <-r.quitCh:
			// TODO: Graceful exit
			log.Info("quitting realtime indexer")
		}
	}
}

func (r *RealtimeIndexer) processBatchedBlockLogs(logs []types.Log, start uint64, end uint64) {
	// Assuming for any height H, either we will have all the concerned logs
	// or not even one
	kv := GroupByBlockNumber(logs)
	dbCtx, dbTx := r.dbconn.BeginTx()

	for block := start; block <= end; block++ {
		_time, err := r.da.GetBlockTimestamp(block)
		util.ENOK(err)

		logs, ok := kv[block]
		blockMeta := itypes.BlockSynopsis{
			Type:    "stats",
			Network: r.dbconn.ChainID,
			Height:  block,
			Time:    _time,
		}
		if !ok || len(logs) == 0 {
			r.dbconn.AddToTx(&dbCtx, dbTx, nil, blockMeta, block)
			continue
		}
		var wg sync.WaitGroup
		var mt sync.Mutex
		var items []interface{}
		for _, _log := range logs {
			wg.Add(1)
			go r.DecodeLog(_log, &mt, &items, &blockMeta, &wg)
		}
		wg.Wait()
		r.dbconn.AddToTx(&dbCtx, dbTx, items, blockMeta, block)
	}
	util.ENOK(r.dbconn.CommitTx(dbTx))
}

func (r *RealtimeIndexer) DecodeLog(l types.Log,
	mt *sync.Mutex,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	wg *sync.WaitGroup) {
	defer wg.Done()

	primaryTopic := l.Topics[0]
	switch primaryTopic {
	// ---- Uniswap V2 ----
	case itypes.UniV2MintTopic:
		instrumentation.MintV2Found.Inc()
		r.processMint(l, items, bm, mt)
	case itypes.UniV2BurnTopic:
		instrumentation.BurnV2Found.Inc()
		r.processBurn(l, items, bm, mt)
	case itypes.UniV2SwapTopic:
		instrumentation.SwapV2Found.Inc()
		r.processUniV2Swap(l, items, bm, mt)

	// ---- Uniswap V3 ----
	case itypes.UniV3MintTopic:
		instrumentation.MintV3Found.Inc()
		r.processMintV3(l, items, bm, mt)
	case itypes.UniV3BurnTopic:
		instrumentation.BurnV3Found.Inc()
		r.processBurnV3(l, items, bm, mt)
	case itypes.UniV3SwapTopic:
		instrumentation.SwapV3Found.Inc()
		r.processUniV3Swap(l, items, bm, mt)

	// ---- ERC 20 ----
	case itypes.ERC20TransferTopic:
		instrumentation.TfrFound.Inc()
		r.processTransfer(l, items, bm, mt)

	}
}

func (r *RealtimeIndexer) Stop() error {
	return nil
}

func (r *RealtimeIndexer) Init() error {
	height, err := r.da.GetCurrentBlockHeight()
	util.ENOK(err)
	r.currentHeight = height
	log.Info("initializing realtime indexer, indexedHeight: "+fmt.Sprint(r.indexedHeight),
		" currentHeight: "+fmt.Sprint(r.currentHeight))
	return nil
}

func (r *RealtimeIndexer) Status() interface{} {
	return nil
}

func (r *RealtimeIndexer) Quit() {
	r.quitCh <- struct{}{}
}

// ---- NEW

func GetBlockCallOpts(blockNumber uint64) *bind.CallOpts {
	return &bind.CallOpts{BlockNumber: big.NewInt(int64(blockNumber))}
}

func (r *RealtimeIndexer) GetFormattedAmount(amount *big.Int,
	callopts *bind.CallOpts,
	erc20Address common.Address) (ok bool,
	formattedAmount *big.Float) {
	tokenDecimals, err := r.da.GetERC20Decimals(erc20Address, callopts)
	if util.IsExecutionReverted(err) {
		// Non ERC-20 contract
		tokenDecimals = 0
	} else {
		if util.IsEthErr(err) {
			return false, big.NewFloat(0.0)
		}
		util.ENOKS(2, err)
	}

	return true, util.DivideBy10pow(amount, tokenDecimals)
}

func AddToSynopsis(mt *sync.Mutex,
	bm *itypes.BlockSynopsis,
	item interface{},
	items *[]interface{},
	_type string,
	condition bool) {
	mt.Lock()
	defer mt.Unlock()
	if condition {
		*items = append(*items, item)
		switch _type {
		case "transfer":
			bm.TransferLogs++
		case "mint":
			bm.MintLogs++
		case "burn":
			bm.BurnLogs++
		case "swap":
			bm.SwapLogs++
		default:
			util.ENOKS(2, fmt.Errorf("unknown add to synopsis: %s", _type))
		}
		bm.TotalLogs++
	}
}
