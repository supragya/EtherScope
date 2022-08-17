package indexer

import (
	"fmt"
	"math"
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
	eventsToIndex []common.Hash

	quitCh chan struct{}
}

func NewRealtimeIndexer(indexedHeight uint64,
	upstreams []string,
	dbconn *db.DBConn,
	eventsToIndex []string) *RealtimeIndexer {
	return &RealtimeIndexer{
		currentHeight: 0,
		indexedHeight: indexedHeight,
		dbconn:        dbconn,
		da:            NewDataAccess(upstreams),
		eventsToIndex: util.ConstructTopics(eventsToIndex),

		quitCh: make(chan struct{}),
	}
}

func (r *RealtimeIndexer) Start() error {
	if r.da.Len() == 0 {
		return EUninitialized
	}
	r.ridxLoop()
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

				log.Info(fmt.Sprintf("sync curr: %d (+%d), processing [%d - %d]",
					r.currentHeight, r.currentHeight-r.indexedHeight, r.indexedHeight, endingBlock))

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
		time, err := r.da.GetBlockTimestamp(block)
		util.ENOK(err)

		logs, ok := kv[block]
		blockMeta := itypes.BlockSynopsis{
			Type:    "stats",
			Network: r.dbconn.ChainID,
			Height:  block,
			Time:    time,
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
	case itypes.TransferTopic:
		r.processTransfer(l, items, bm, mt)
	case itypes.MintTopic:
		r.processMint(l, items, bm, mt)
	case itypes.IncreaseLiquidityTopic:
		r.processMintV3(l, items, bm, mt)
	case itypes.DecreaseLiquidityTopic:
		r.processBurnV3(l, items, bm, mt)
	case itypes.BurnTopic:
		r.processBurn(l, items, bm, mt)
	case itypes.UniV2Swap:
		r.processUniV2Swap(l, items, bm, mt)
	case itypes.UniV3Swap:
		r.processUniV3Swap(l, items, bm, mt)
	}
}

func (r *RealtimeIndexer) processTransfer(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	ok, sender, recv, amt := InfoTransfer(l)

	if !ok {
		return
	}

	callopts := GetBlockCallOpts(l.BlockNumber)
	ok, formattedAmount := r.GetFormattedAmount(amt, callopts, l.Address)

	if !ok {
		return
	}

	transfer := itypes.Transfer{
		Type:        "transfer",
		Network:     r.dbconn.ChainID,
		LogIdx:      l.Index,
		Transaction: l.TxHash,
		Time:        bm.Time,
		Height:      l.BlockNumber,
		Token:       l.Address,
		Sender:      sender,
		Receiver:    recv,
		Amount:      formattedAmount,
		AmountUSD:   0, // TODO
	}

	AddToSynopsis(mt, bm, transfer, items, "transfer", true)
}

func (r *RealtimeIndexer) processMint(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {

	callopts := GetBlockCallOpts(l.BlockNumber)

	// Test if the contract is a UniswapV2 type contract
	if !r.isUniswapV2Pair(l.Address, callopts) {
		return
	}

	ok, sender, am0, am1 := r.InfoUniV2Mint(l)
	if !ok {
		return
	}

	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV2(am0, am1, callopts, l.Address)
	if !ok {
		return
	}

	// Assumed infallible since if err != nil, this code should not be reachable
	// due to above condition
	t0, t1, _ := r.da.GetTokensUniV2(l.Address, callopts)

	reserves, err := r.da.GetBalances([]Tuple2[common.Address, common.Address]{
		{l.Address, t0}, {l.Address, t1},
	}, callopts)
	util.ENOK(err)

	token0Price, token1Price, amountusd, tokenMeta := r.da.GetPricesForBlock(r.dbconn.ChainID, callopts, t0, t1, f0, f1)

	mint := itypes.Mint{
		Type:         "mint",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       sender,
		Receiver:     sender,
		PairContract: l.Address,
		Token0:       t0,
		Token1:       t1,
		Amount0:      f0,
		Amount1:      f1,
		Reserve0:     util.DivideBy10pow(reserves[0].Second, t0d),
		Reserve1:     util.DivideBy10pow(reserves[1].Second, t1d),
		AmountUSD:    amountusd,
		Price0:       token0Price,
		Price1:       token1Price,
		Meta:         tokenMeta,
	}

	AddToSynopsis(mt, bm, mint, items, "mint", true)
}

func (r *RealtimeIndexer) processMintV3(
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

	tokenID := util.ExtractIntFromBytes(l.Topics[1][:])

	// Test if the contract is a UniswapV2 type contract
	token0, token1, err := r.da.GetTokensUniV3(l.Address, tokenID, callopts)

	// Check if we have enough data to retrieve amount of token being minted
	if len(l.Data) < 96 {
		return
	}

	if util.IsExecutionReverted(err) {
		return
	}

	am0 := util.ExtractIntFromBytes(l.Data[32:64])
	am1 := util.ExtractIntFromBytes(l.Data[32:64])

	if len(am0.Bits()) == 0 || len(am1.Bits()) == 0 {
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

	formattedAmount0 := util.DivideBy10pow(am0, token0Decimals)
	formattedAmount1 := util.DivideBy10pow(am1, token1Decimals)
	token0Price, token1Price, amountusd, tokenMeta := r.da.GetPricesForBlock(r.dbconn.ChainID, callopts, token0, token1, formattedAmount0, formattedAmount1)

	mint := itypes.Mint{
		Type:         "mint",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       sender,
		Receiver:     sender,
		PairContract: l.Address,
		Token0:       token0,
		Token1:       token1,
		Amount0:      formattedAmount0,
		Amount1:      formattedAmount1,
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve1, token1Decimals),
		AmountUSD:    amountusd,
		Price0:       token0Price,
		Price1:       token1Price,
		Meta:         tokenMeta,
	}
	mt.Lock()
	defer mt.Unlock()
	is0Nan := math.IsInf(token0Price, 0)
	is1Nan := math.IsInf(token1Price, 0)
	if amountusd > -1 && !is0Nan && !is1Nan {
		*items = append(*items, mint)
		bm.MintLogs++
		bm.TotalLogs++
	}
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

	am0 := util.ExtractIntFromBytes(l.Data[:32])
	am1 := util.ExtractIntFromBytes(l.Data[32:64])

	if len(am0.Bits()) == 0 || len(am1.Bits()) == 0 {
		return
	}

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

	formattedAmount0 := util.DivideBy10pow(am0, token0Decimals)
	formattedAmount1 := util.DivideBy10pow(am1, token1Decimals)
	token0Price, token1Price, amountusd, tokenMeta := r.da.GetPricesForBlock(r.dbconn.ChainID, callopts, token0, token1, formattedAmount0, formattedAmount1)

	burn := itypes.Burn{
		Type:         "burn",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       sender,
		Receiver:     recipient,
		PairContract: l.Address,
		Token0:       token0,
		Token1:       token1,
		Amount0:      formattedAmount0,
		Amount1:      formattedAmount1,
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve1, token1Decimals),
		AmountUSD:    amountusd,
		Price0:       token0Price,
		Price1:       token1Price,
		Meta:         tokenMeta,
	}

	mt.Lock()
	defer mt.Unlock()
	is0Nan := math.IsInf(token0Price, 0)
	is1Nan := math.IsInf(token1Price, 0)
	if amountusd > -1 && !is0Nan && !is1Nan {
		*items = append(*items, burn)
		bm.BurnLogs++
		bm.TotalLogs++
	}
}

func (r *RealtimeIndexer) processBurnV3(
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

	tokenID := util.ExtractIntFromBytes(l.Topics[1][:])

	// Test if the contract is a UniswapV2 type contract
	token0, token1, err := r.da.GetTokensUniV3(l.Address, tokenID, callopts)

	// Check if we have enough data to retrieve amount of token being minted
	if len(l.Data) < 96 {
		return
	}

	if util.IsExecutionReverted(err) {
		return
	}

	am0 := util.ExtractIntFromBytes(l.Data[32:64])
	am1 := util.ExtractIntFromBytes(l.Data[32:64])

	if len(am0.Bits()) == 0 || len(am1.Bits()) == 0 {
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

	formattedAmount0 := util.DivideBy10pow(am0, token0Decimals)
	formattedAmount1 := util.DivideBy10pow(am1, token1Decimals)
	token0Price, token1Price, amountusd, tokenMeta := r.da.GetPricesForBlock(r.dbconn.ChainID, callopts, token0, token1, formattedAmount0, formattedAmount1)

	mint := itypes.Burn{
		Type:         "burn",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       sender,
		Receiver:     sender,
		PairContract: l.Address,
		Token0:       token0,
		Token1:       token1,
		Amount0:      formattedAmount0,
		Amount1:      formattedAmount1,
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve1, token1Decimals),
		AmountUSD:    amountusd,
		Price0:       token0Price,
		Price1:       token1Price,
		Meta:         tokenMeta,
	}
	mt.Lock()
	defer mt.Unlock()
	is0Nan := math.IsInf(token0Price, 0)
	is1Nan := math.IsInf(token1Price, 0)
	if amountusd > -1 && !is0Nan && !is1Nan {
		*items = append(*items, mint)
		bm.BurnLogs++
		bm.TotalLogs++
	}
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

	formattedAmount0 := util.DivideBy10pow(am0, token0Decimals)
	formattedAmount1 := util.DivideBy10pow(am1, token1Decimals)
	token0Price, token1Price, amountusd, tokenMeta := r.da.GetPricesForBlock(r.dbconn.ChainID, callopts, token0, token1, formattedAmount0, formattedAmount1)

	swap := itypes.Swap{
		Type:         "swap",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       util.ExtractAddressFromLogTopic(l.Topics[1]),
		Receiver:     util.ExtractAddressFromLogTopic(l.Topics[2]),
		PairContract: l.Address,
		Token0:       token0,
		Token1:       token1,
		Amount0:      formattedAmount0,
		Amount1:      formattedAmount1,
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve1, token1Decimals),
		AmountUSD:    amountusd,
		Price0:       token0Price,
		Price1:       token1Price,
		Meta:         tokenMeta,
	}

	mt.Lock()
	defer mt.Unlock()
	is0Nan := math.IsInf(token0Price, 0)
	is1Nan := math.IsInf(token1Price, 0)
	if amountusd > -1 && !is0Nan && !is1Nan {
		*items = append(*items, swap)
		bm.SwapLogs++
		bm.TotalLogs++
	}
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

	formattedAmount0 := util.DivideBy10pow(am0, token0Decimals)
	formattedAmount1 := util.DivideBy10pow(am1, token1Decimals)
	token0Price, token1Price, amountusd, tokenMeta := r.da.GetPricesForBlock(r.dbconn.ChainID, callopts, token0, token1, formattedAmount0, formattedAmount1)

	swap := itypes.Swap{
		Type:         "swap",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       util.ExtractAddressFromLogTopic(l.Topics[1]),
		Receiver:     util.ExtractAddressFromLogTopic(l.Topics[2]),
		PairContract: l.Address,
		Token0:       token0,
		Token1:       token1,
		Amount0:      formattedAmount0,
		Amount1:      formattedAmount1,
		Reserve0:     util.DivideBy10pow(reserves.Reserve0, token0Decimals),
		Reserve1:     util.DivideBy10pow(reserves.Reserve1, token1Decimals),
		AmountUSD:    amountusd,
		Price0:       token0Price,
		Price1:       token1Price,
		Meta:         tokenMeta,
	}

	mt.Lock()
	defer mt.Unlock()
	is0Nan := math.IsInf(token0Price, 0)
	is1Nan := math.IsInf(token1Price, 0)
	if amountusd > -1 && !is0Nan && !is1Nan {
		*items = append(*items, swap)
		bm.SwapLogs++
		bm.TotalLogs++
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

func HasSufficientData(l types.Log,
	requiredTopicLen int,
	requiredDataLen int) bool {
	return len(l.Topics) == requiredTopicLen && len(l.Data) == requiredDataLen
}

func InfoTransfer(l types.Log) (hasSufficientData bool,
	sender common.Address,
	receiver common.Address,
	amount *big.Int) {
	if !HasSufficientData(l, 3, 32) {
		return false,
			common.Address{},
			common.Address{},
			big.NewInt(0)
	}
	return true,
		util.ExtractAddressFromLogTopic(l.Topics[1]),
		util.ExtractAddressFromLogTopic(l.Topics[2]),
		util.ExtractIntFromBytes(l.Data[:32])
}

func (r *RealtimeIndexer) InfoUniV2Mint(l types.Log) (hasSufficientData bool,
	sender common.Address,
	amount0 *big.Int,
	amount1 *big.Int) {
	sender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if !HasSufficientData(l, 1, 64) || err != nil {
		if !util.IsEthErr(err) {
			util.ENOKS(2, err)
		}
		return false,
			common.Address{},
			big.NewInt(0),
			big.NewInt(0)
	}
	return true,
		sender,
		util.ExtractIntFromBytes(l.Data[:32]),
		util.ExtractIntFromBytes(l.Data[32:64])
}

func (r *RealtimeIndexer) GetFormattedAmount(amount *big.Int,
	callopts *bind.CallOpts,
	erc20Address common.Address) (ok bool,
	formattedAmount *big.Float) {
	erc, client := r.da.GetERC20(erc20Address)

	tokenDecimals, err := r.da.GetERC20Decimals(erc, client, callopts)
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
		default:
			util.ENOKS(2, fmt.Errorf("unknown add to synopsis: %s", _type))
		}
		bm.TotalLogs++
	}
}

func (r *RealtimeIndexer) isUniswapV2Pair(address common.Address,
	callopts *bind.CallOpts) bool {
	_, _, err := r.da.GetTokensUniV2(address, callopts)
	if err != nil {
		return true
	}
	// Execution Revert: Could be a non uniswap contract (like AAVE V2). Example log:
	// https://etherscan.io/tx/0x65ed6ba09f2a22805b772ff607f81fa4bb5d93ce287ecf05ab5ad97cab34c97c#eventlog logIdx 180
	// not handled currently
	if !util.IsExecutionReverted(err) {
		util.ENOKS(2, err)
	}
	return false
}

// TODO: refactor this
func (r *RealtimeIndexer) GetFormattedAmountsUniV2(amount0 *big.Int,
	amount1 *big.Int,
	callopts *bind.CallOpts,
	address common.Address) (ok bool,
	formattedAmount0 *big.Float,
	formattedAmount1 *big.Float,
	token0Decimals uint8,
	token1Decimals uint8) {
	t0, t1, err := r.da.GetTokensUniV2(address, callopts)
	if err != nil {
		return false,
			big.NewFloat(0.0),
			big.NewFloat(0.0),
			0,
			0
	}

	erc0, client0 := r.da.GetERC20(t0)

	token0Decimals, err = r.da.GetERC20Decimals(erc0, client0, callopts)
	if util.IsExecutionReverted(err) {
		// Non ERC-20 contract
		token0Decimals = 0
	} else {
		if util.IsEthErr(err) {
			return false,
				big.NewFloat(0.0),
				big.NewFloat(0.0),
				0,
				0
		}
		util.ENOKS(2, err)
	}

	erc1, client1 := r.da.GetERC20(t1)

	token1Decimals, err = r.da.GetERC20Decimals(erc1, client1, callopts)
	if util.IsExecutionReverted(err) {
		// Non ERC-20 contract
		token1Decimals = 0
	} else {
		if util.IsEthErr(err) {
			return false,
				big.NewFloat(0.0),
				big.NewFloat(0.0),
				0,
				0
		}
		util.ENOKS(2, err)
	}

	return true,
		util.DivideBy10pow(amount0, token0Decimals),
		util.DivideBy10pow(amount1, token1Decimals),
		token0Decimals,
		token1Decimals
}
