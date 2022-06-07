package indexer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/ERC20"
	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/univ2pair"
	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	EUninitialized = errors.New("uninitialized indexer")
)

type RealtimeIndexer struct {
	currentHeight uint64
	indexedHeight uint64
	upstreams     *LatencySortedPool
	dbconn        *db.DBConn

	quitCh chan struct{}
}

type BlockSynopsis struct {
	totalLogs uint64
	MintLogs  uint64
	BurnLogs  uint64
}

func NewRealtimeIndexer(indexedHeight uint64, upstreams []string, dbconn *db.DBConn) *RealtimeIndexer {
	return &RealtimeIndexer{
		currentHeight: 0,
		indexedHeight: indexedHeight,
		upstreams:     NewLatencySortedPool(upstreams),
		dbconn:        dbconn,
		quitCh:        make(chan struct{}),
	}
}

func (r *RealtimeIndexer) Start() error {
	if r.indexedHeight == 0 || r.upstreams.Len() == 0 {
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
			util.ENOK(r.populateCurrentHeight())
			if r.currentHeight == r.indexedHeight {
				continue
			}
			endingBlock := r.currentHeight
			if (endingBlock - r.indexedHeight) > maxBlockSpanPerCall {
				endingBlock = r.indexedHeight + maxBlockSpanPerCall
			}

			log.Info(fmt.Sprintf("sync up: %d, indexed: %d, to: %d, dist: %d",
				r.currentHeight, r.indexedHeight, endingBlock, r.currentHeight-r.indexedHeight))

			logs, err := r.getLogs(ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(r.indexedHeight + 1)),
				ToBlock:   big.NewInt(int64(endingBlock)),
				Topics:    [][]common.Hash{{MintTopic, BurnTopic}},
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
	kv := GroupByBlockNumber(logs)

	var dbwg sync.WaitGroup
	for block := start; block <= end; block++ {
		logs, ok := kv[block]
		if !ok {
			continue
		}
		var wg sync.WaitGroup
		var mt sync.Mutex
		var items []interface{}
		blockMeta := BlockSynopsis{}
		for _, log := range logs {
			blockMeta.totalLogs++
			go r.DecodeLog(&log, &mt, &items, &blockMeta, &wg)
		}
		wg.Wait()
		log.Info(block, " done blockmeta ", blockMeta)
		go r.persistToDB(&dbwg, items, blockMeta)
	}
	dbwg.Wait()
}

func (r *RealtimeIndexer) persistToDB(dbwg *sync.WaitGroup, items []interface{}, bm BlockSynopsis) {
	dbwg.Add(1)
	defer dbwg.Done()

	// ctx := context.Background()
	// tx, err := r.dbconn.Conn.BeginTx(ctx, nil)
	// util.ENOK(err)
	// for _, item := range items {
	// 	switch item.(type) {
	// 	case Mint:
	// 		tx.ExecContext(ctx, "INSERT INTO ")
	// 	}
	// }
}

func (r *RealtimeIndexer) DecodeLog(l *types.Log,
	mt *sync.Mutex,
	items *[]interface{},
	bm *BlockSynopsis,
	wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	primaryTopic := l.Topics[0]
	switch primaryTopic {
	case MintTopic:
		// Get an upstream
		cl := r.upstreams.GetItem()
		paircontract, err := univ2pair.NewUniv2pair(l.Address, cl)
		util.ENOK(err)
		callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(l.BlockNumber))}
		// check if token0 exists
		// Break if topic is not a DEX liquidity add / removal
		// Many DeFi apps that aren't dexes use Mint & Burn events
		token0, err := paircontract.Token0(callopts)
		if err != nil {
			log.Trace("error while getting token0 ", err)
			return
		}
		token1, err := paircontract.Token1(callopts)
		if err != nil {
			log.Trace("error while getting token1 ", err)
			return
		}

		t0Contract, err := ERC20.NewERC20(token0, cl)
		util.ENOK(err)
		t1Contract, err := ERC20.NewERC20(token1, cl)
		util.ENOK(err)

		t0Decimals, err := t0Contract.Decimals(callopts)
		if err != nil {
			log.Trace("error while getting token0 decimals ", err)
			return
		}

		t1Decimals, err := t1Contract.Decimals(callopts)
		if err != nil {
			log.Trace("error while getting token1 decimals ", err)
			return
		}

		reserves, err := paircontract.GetReserves(callopts)
		if err != nil {
			log.Trace("error while retrieving reserves ", err)
			return
		}

		t0reserves := util.DivideBy10pow(reserves.Reserve0, t0Decimals)
		t1reserves := util.DivideBy10pow(reserves.Reserve1, t1Decimals)

		mint := Mint{
			logIdx:       l.Index,
			transaction:  l.TxHash,
			height:       l.BlockNumber,
			sender:       l.Address, // FIXME
			pairContract: l.Address,
			token0:       token0,
			token1:       token1,
			amount0:      0, // FIXME
			amount1:      0, // FIXME
			reserve0:     t0reserves,
			reserve1:     t1reserves,
		}
		mt.Lock()
		defer mt.Unlock()
		*items = append(*items, mint)
		bm.MintLogs++
	}
}

func (r *RealtimeIndexer) Stop() error {
	return nil
}

func (r *RealtimeIndexer) Init() error {
	if err := r.populateCurrentHeight(); err != nil {
		return err
	}
	log.Info("initializing realtime indexer, indexedHeight: "+fmt.Sprint(r.indexedHeight),
		" currentHeight: "+fmt.Sprint(r.currentHeight))
	return nil
}

func (r *RealtimeIndexer) getLogs(fq ethereum.FilterQuery) ([]types.Log, error) {
	var logs []types.Log
	var retries = 0
	var err error
	for {
		if retries == WD {
			return logs, errors.New("could not fetch logs, retried " + fmt.Sprint(WD) + " times. Last err: " + err.Error())
		}
		cl := r.upstreams.GetItem()

		start := time.Now()
		logs, err = cl.FilterLogs(context.Background(), fq)
		r.upstreams.Report(cl, time.Now().Sub(start).Seconds(), err != nil)
		if err == nil {
			break
		}
		retries++
	}
	return logs, nil
}

func (r *RealtimeIndexer) populateCurrentHeight() error {
	var currentHeight uint64 = 0
	var retries = 0
	var err error
	for {
		if retries == WD {
			return errors.New("could not init realtime indexer, retried " + fmt.Sprint(WD) + " times. Last err: " + err.Error())
		}
		cl := r.upstreams.GetItem()

		start := time.Now()
		currentHeight, err = cl.BlockNumber(context.Background())
		r.upstreams.Report(cl, time.Now().Sub(start).Seconds(), err != nil)
		if err == nil {
			break
		}
		retries++
	}
	r.currentHeight = currentHeight
	return nil
}

func (r *RealtimeIndexer) Status() interface{} {
	return nil
}

func (r *RealtimeIndexer) Quit() {
	r.quitCh <- struct{}{}
}
