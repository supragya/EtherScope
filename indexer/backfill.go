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

type BackfillIndexer struct {
	currentHeight uint64
	indexedHeight uint64
	dbconn        *db.DBConn
	da            *DataAccess

	quitCh chan struct{}
}

func NewBackfillIndexer(indexedHeight uint64, upstreams []string, dbconn *db.DBConn) *BackfillIndexer {
	return &BackfillIndexer{
		currentHeight: 0,
		indexedHeight: indexedHeight,
		dbconn:        dbconn,
		da:            NewDataAccess(upstreams),

		quitCh: make(chan struct{}),
	}
}

func (r *BackfillIndexer) Start() error {
	if r.indexedHeight == 0 || r.da.Len() == 0 {
		return EUninitialized
	}
	r.ridxLoop()
	time.Sleep(time.Second * 2)
	return nil
}

func (r *BackfillIndexer) ridxLoop() {
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
				Topics:    [][]common.Hash{{itypes.MintTopic, itypes.BurnTopic}},
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

func (r *BackfillIndexer) processBatchedBlockLogs(logs []types.Log, start uint64, end uint64) {
	// Assuming for any height H, either we will have all the concerned logs
	// or not even one
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

func (r *BackfillIndexer) DecodeLog(l *types.Log,
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
	}
}

func (r *BackfillIndexer) processMint(
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
		Amount0:      0, // FIXME
		Amount1:      0, // FIXME
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve0, token1Decimals),
	}
	mt.Lock()
	defer mt.Unlock()
	*items = append(*items, mint)
	bm.MintLogs++
	bm.TotalLogs++
}

func (r *BackfillIndexer) processBurn(
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
		Amount0:      0, // FIXME
		Amount1:      0, // FIXME
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve0, token1Decimals),
	}
	mt.Lock()
	defer mt.Unlock()
	*items = append(*items, mint)
	bm.MintLogs++
	bm.TotalLogs++
}

func (r *BackfillIndexer) Stop() error {
	return nil
}

func (r *BackfillIndexer) Init() error {
	height, err := r.da.GetCurrentBlockHeight()
	util.ENOK(err)
	r.currentHeight = height
	log.Info("initializing realtime indexer, indexedHeight: "+fmt.Sprint(r.indexedHeight),
		" currentHeight: "+fmt.Sprint(r.currentHeight))
	return nil
}

func (r *BackfillIndexer) Status() interface{} {
	return nil
}

func (r *BackfillIndexer) Quit() {
	r.quitCh <- struct{}{}
}
