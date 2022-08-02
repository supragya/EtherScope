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
	if r.da.Len() == 0 {
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
				syncing := "head"
				if !isOnHead {
					syncing = "sync"
				}

				log.Info(fmt.Sprintf("%s curr: %d (+%d), processing [%d - %d]",
					syncing, r.currentHeight, r.currentHeight-r.indexedHeight, r.indexedHeight, endingBlock))

				logs, err := r.da.GetFilteredLogs(ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(r.indexedHeight + 1)),
					ToBlock:   big.NewInt(int64(endingBlock)),
					Topics: [][]common.Hash{{
						itypes.MintTopic,
						itypes.BurnTopic,
						itypes.UniV2Swap,
						itypes.UniV3Swap,
					}},
				})

				if err != nil {
					log.Error(err)
					continue
				}

				r.processBatchedBlockLogs(logs, r.indexedHeight+1, endingBlock)

				r.indexedHeight = endingBlock

				if isOnHead {
					break
				}
			}
		case <-r.quitCh:
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
		logs, ok := kv[block]
		blockMeta := itypes.BlockSynopsis{
			Height: block,
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
	case itypes.MintTopic:
		r.processMint(l, items, bm, mt)
	case itypes.BurnTopic:
		r.processBurn(l, items, bm, mt)
	case itypes.UniV2Swap:
		r.processUniV2Swap(l, items, bm, mt)
	case itypes.UniV3Swap:
		r.processUniV3Swap(l, items, bm, mt)
	}
}

func (r *RealtimeIndexer) processMint(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(l.BlockNumber))}

	sender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	// Test if the contract is a UniswapV2 type contract
	token0, token1, err := r.da.GetTokensUniV2(l.Address, callopts)

	if util.IsExecutionReverted(err) {
		// Could be a non uniswap contract (like AAVE V2). Example log:
		// https://etherscan.io/tx/0x65ed6ba09f2a22805b772ff607f81fa4bb5d93ce287ecf05ab5ad97cab34c97c#eventlog logIdx 180
		// not handled currently
		return
	}

	// Check if we have enough data to retrieve amount of token being minted
	if len(l.Data) < 32 {
		return
	}

	amount0 := big.NewFloat(0.0).SetInt(big.NewInt(0).SetBytes(l.Data[:32]))

	erc0, client0 := r.da.GetERC20(token0)
	erc1, client1 := r.da.GetERC20(token1)

	token0Decimals, err := r.da.GetERC20Decimals(erc0, client0, callopts)
	if util.IsExecutionReverted(err) {
		// Non ERC-20 contract
		token0Decimals = 0
	} else {
		if util.IsEthErr(err) {
			return
		}
		util.ENOK(err)
	}

	token1Decimals, err := r.da.GetERC20Decimals(erc1, client1, callopts)
	if util.IsExecutionReverted(err) {
		// Non ERC-20 contract
		token1Decimals = 0
	} else {
		if util.IsEthErr(err) {
			return
		}
		util.ENOK(err)
	}

	reserves, err := r.da.GetDEXReserves(l.Address, erc0, client0, erc1, client1, callopts)
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
		Amount0:      amount0,
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve1, token1Decimals),
	}
	mt.Lock()
	defer mt.Unlock()
	*items = append(*items, mint)
	bm.MintLogs++
	bm.TotalLogs++
}

func (r *RealtimeIndexer) processBurn(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(l.BlockNumber))}

	if len(l.Topics) < 3 {
		return
	}

	sender := util.ExtractAddressFromLogTopic(l.Topics[1])
	recipient := util.ExtractAddressFromLogTopic(l.Topics[2])

	// Check if we have enough data to retrieve amount of token being minted
	if len(l.Data) < 64 {
		return
	}

	amount0 := big.NewFloat(0.0).SetInt(big.NewInt(0).SetBytes(l.Data[:32]))
	amount1 := big.NewFloat(0.0).SetInt(big.NewInt(0).SetBytes(l.Data[32:64]))

	// Test if the contract is a UniswapV2 type contract
	token0, token1, err := r.da.GetTokensUniV2(l.Address, callopts)

	if util.IsExecutionReverted(err) {
		// Could be a non uniswap contract (Seen seldom in practice).
		// TODO document any such logs in code

		log.Info("Burn execution revert. ADD THIS TRANSACTION TO CODE COMMENT, CONTACT AUTHOR, details: [", l.TxHash, " idx ", l.Index, "]")
		return
	}

	erc0, client0 := r.da.GetERC20(token0)
	erc1, client1 := r.da.GetERC20(token1)

	token0Decimals, err := r.da.GetERC20Decimals(erc0, client0, callopts)
	if util.IsExecutionReverted(err) {
		// Non ERC-20 contract
		token0Decimals = 0
	} else {
		if util.IsEthErr(err) {
			return
		}
		util.ENOK(err)
	}

	token1Decimals, err := r.da.GetERC20Decimals(erc1, client1, callopts)
	if util.IsExecutionReverted(err) {
		// Non ERC-20 contract
		token1Decimals = 0
	} else {
		if util.IsEthErr(err) {
			return
		}
		util.ENOK(err)
	}

	reserves, err := r.da.GetDEXReserves(l.Address, erc0, client0, erc1, client1, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	burn := itypes.Burn{
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         time.Now().Unix(),
		Height:       l.BlockNumber,
		Sender:       sender,
		Receiver:     recipient,
		PairContract: l.Address,
		Token0:       token0,
		Token1:       token1,
		Amount0:      amount0,
		Amount1:      amount1,
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve1, token1Decimals),
	}

	mt.Lock()
	defer mt.Unlock()
	*items = append(*items, burn)
	bm.BurnLogs++
	bm.TotalLogs++
}

func (r *RealtimeIndexer) processUniV2Swap(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	if len(l.Data) != 128 {
		log.Warn("unknown swap event data len: ", len(l.Data), " expected 128")
		return
	}

	am0In := util.ExtractIntFromBytes(l.Data[0:32])
	am1In := util.ExtractIntFromBytes(l.Data[32:64])
	am0Out := util.ExtractIntFromBytes(l.Data[64:96])
	am1Out := util.ExtractIntFromBytes(l.Data[96:128])

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

	erc0, client0 := r.da.GetERC20(token0)
	erc1, client1 := r.da.GetERC20(token1)

	token0Decimals, err := r.da.GetERC20Decimals(erc0, client0, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token1Decimals, err := r.da.GetERC20Decimals(erc1, client1, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	reserves, err := r.da.GetDEXReserves(l.Address, erc0, client0, erc1, client1, callopts)
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
		Reserve1:     util.DivideBy10pow(reserves.Reserve1, token1Decimals),
	}
	mt.Lock()
	defer mt.Unlock()
	*items = append(*items, swap)
	bm.SwapLogs++
	bm.TotalLogs++
}

func (r *RealtimeIndexer) processUniV3Swap(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	if len(l.Data) != 160 {
		log.Warn("unknown swap event data len: ", len(l.Data), " expected 160")
		return
	}

	am0 := util.ExtractIntFromBytes(l.Data[0:32])
	am1 := util.ExtractIntFromBytes(l.Data[32:64])

	callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(l.BlockNumber))}
	token0, token1, err := r.da.GetTokensUniV2(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	erc0, client0 := r.da.GetERC20(token0)
	erc1, client1 := r.da.GetERC20(token1)

	token0Decimals, err := r.da.GetERC20Decimals(erc0, client0, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token1Decimals, err := r.da.GetERC20Decimals(erc1, client1, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	reserves, err := r.da.GetDEXReserves(l.Address, erc0, client0, erc1, client1, callopts)
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
		Reserve1:     util.DivideBy10pow(reserves.Reserve1, token1Decimals),
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
