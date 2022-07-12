package indexer

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RealtimeIndexer struct {
	currentHeight uint64
	indexedHeight uint64
	dbconn        *db.DBConn
	da            *DataAccess

	quitCh chan struct{}
}

func NewRealtimeIndexer(indexedHeight uint64, upstreams []string, dbconn *db.DBConn) *RealtimeIndexer {
	return &RealtimeIndexer{
		currentHeight: 0,
		indexedHeight: indexedHeight,
		dbconn:        dbconn,
		da:            NewDataAccess(upstreams),

		quitCh: make(chan struct{}),
	}
}

func (r *RealtimeIndexer) Start() error {
	if r.indexedHeight == 0 || r.da.Len() == 0 {
		return EUninitialized
	}
	r.ridxLoop()
	time.Sleep(time.Second * 2)
	return nil
}

func (r *RealtimeIndexer) ridxLoop() {
	maxBlockSpanPerCall := viper.GetUint64("general.maxBlockSpanPerCall")
	for {
		select {
		case <-time.After(time.Second):
			height, err := r.da.GetCurrentBlockHeight()
			util.ENOK(err)
			r.currentHeight = height

			if r.currentHeight == r.indexedHeight {
				continue
			}
			endingBlock := r.currentHeight
			if (endingBlock - r.indexedHeight) > maxBlockSpanPerCall {
				endingBlock = r.indexedHeight + maxBlockSpanPerCall
			}

			log.Info(fmt.Sprintf("sync up: %d, indexed: %d, to: %d, dist: %d",
				r.currentHeight, r.indexedHeight, endingBlock, r.currentHeight-r.indexedHeight))

			logs, err := r.da.GetFilteredLogs(ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(r.indexedHeight + 1)),
				ToBlock:   big.NewInt(int64(endingBlock)),
				Topics:    [][]common.Hash{{itypes.MintTopic, itypes.BurnTopic, itypes.UniV2Swap}},
			})

			if err != nil {
				log.Error(err)
				continue
			}

			r.processBatchedBlockLogs(logs, r.indexedHeight+1, endingBlock)

			r.indexedHeight = endingBlock
		case <-r.quitCh:
			log.Info("quitting realtime indexer")
		}
	}
}

func (r *RealtimeIndexer) processBatchedBlockLogs(logs []types.Log, start uint64, end uint64) {
	// Assuming for any height H, either we will have all the concerned logs
	// or not even one
	log.Info("found logs for batch: ", len(logs))
	kv := GroupByBlockNumber(logs)
	dbCtx, dbTx := r.dbconn.BeginTx()

	for block := start; block <= end; block++ {
		logs, ok := kv[block]
		blockMeta := itypes.BlockSynopsis{}
		if !ok || len(logs) == 0 {
			r.dbconn.AddToTx(&dbCtx, dbTx, nil, blockMeta, block)
			continue
		}
		var wg sync.WaitGroup
		var mt sync.Mutex
		var items []interface{}
		for _, log := range logs {
			go r.DecodeLog(&log, &mt, &items, &blockMeta, &wg)
		}
		wg.Wait()
		r.dbconn.AddToTx(&dbCtx, dbTx, items, blockMeta, block)
	}
	util.ENOK(dbTx.Commit())
}

func (r *RealtimeIndexer) DecodeLog(l *types.Log,
	mt *sync.Mutex,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	primaryTopic := l.Topics[0]
	switch primaryTopic {
	case itypes.MintTopic:
		r.processMint(l, items, bm, mt)
	case itypes.BurnTopic:
		r.processBurn(l, items, bm, mt)
	case itypes.UniV2Swap:
		r.processUniV2Swap(l, items, bm, mt)
	}
}

func (r *RealtimeIndexer) processMint(
	l *types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(l.BlockNumber))}
	token0, token1, err := r.da.GetTokensUniV2(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0Decimals, err := r.da.GetERC20Decimals(token0, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token1Decimals, err := r.da.GetERC20Decimals(token1, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	reserves, err := r.da.GetReservesUniV2(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	sender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	mint := itypes.Mint{
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         time.Now().Unix(),
		Height:       l.BlockNumber,
		Sender:       sender,
		PairContract: l.Address,
		Token0:       token0,
		Token1:       token1,
		Amount0:      big.NewFloat(0.0), // FIXME
		Amount1:      big.NewFloat(0.0), // FIXME
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve0, token1Decimals),
	}
	mt.Lock()
	defer mt.Unlock()
	*items = append(*items, mint)
	bm.MintLogs++
	bm.TotalLogs++
}

func (r *RealtimeIndexer) processBurn(
	l *types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(l.BlockNumber))}
	token0, token1, err := r.da.GetTokensUniV2(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0Decimals, err := r.da.GetERC20Decimals(token0, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token1Decimals, err := r.da.GetERC20Decimals(token1, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	reserves, err := r.da.GetReservesUniV2(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	sender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	mint := itypes.Burn{
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         time.Now().Unix(),
		Height:       l.BlockNumber,
		Sender:       sender,
		PairContract: l.Address,
		Token0:       token0,
		Token1:       token1,
		Amount0:      big.NewFloat(0.0), // FIXME
		Amount1:      big.NewFloat(0.0), // FIXME
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve0, token1Decimals),
	}
	mt.Lock()
	defer mt.Unlock()
	*items = append(*items, mint)
	bm.MintLogs++
	bm.TotalLogs++
}

func (r *RealtimeIndexer) processUniV2Swap(
	l *types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	if len(l.Data) != 128 {
		log.Warn("unknown swap event data len: ", len(l.Data), " expected 128")
		return
	}

	am0In := util.ExtractUintFromBytes(l.Data[0:32])
	am1In := util.ExtractUintFromBytes(l.Data[32:64])
	am0Out := util.ExtractUintFromBytes(l.Data[64:96])
	am1Out := util.ExtractUintFromBytes(l.Data[96:128])

	am0, am1 := big.NewInt(0), big.NewInt(0)
	if am0In.Cmp(big.NewInt(0)) == 0 {
		am0 = am0.Neg(am0Out)
		am1 = am1In
	} else {
		am0 = am0In
		am1 = am1.Neg(am1Out)
	}

	callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(l.BlockNumber))}
	token0, token1, err := r.da.GetTokensUniV2(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0Decimals, err := r.da.GetERC20Decimals(token0, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token1Decimals, err := r.da.GetERC20Decimals(token1, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	reserves, err := r.da.GetReservesUniV2(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	swap := itypes.Swap{
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         time.Now().Unix(),
		Height:       l.BlockNumber,
		Sender:       util.ExtractAddressFromLogTopic(l.Topics[1]),
		Receiver:     util.ExtractAddressFromLogTopic(l.Topics[2]),
		PairContract: l.Address,
		Token0:       token0,
		Token1:       token1,
		Amount0:      util.DivideBy10pow(am0, token0Decimals),
		Amount1:      util.DivideBy10pow(am1, token1Decimals),
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve0, token1Decimals),
	}
	mt.Lock()
	defer mt.Unlock()
	*items = append(*items, swap)
	bm.SwapLogs++
	bm.TotalLogs++
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
