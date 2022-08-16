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
		eventsToIndex: ConstructTopics(eventsToIndex),

		quitCh: make(chan struct{}),
	}
}

func ConstructTopics(eventsToIndex []string) []common.Hash {
	topicsList := []common.Hash{}
	for _, t := range eventsToIndex {
		switch t {
		case "UniswapV2Swap":
			topicsList = append(topicsList, itypes.UniV2Swap)
		case "UniswapV2Mint":
			topicsList = append(topicsList, itypes.MintTopic)
		case "UniswapV2Burn":
			topicsList = append(topicsList, itypes.BurnTopic)
		case "UniswapV3Swap":
			topicsList = append(topicsList, itypes.UniV3Swap)
		case "UniswapV3IncreaseLiquidity":
			topicsList = append(topicsList, itypes.IncreaseLiquidityTopic)
		case "UniswapV3DecreaseLiquidity":
			topicsList = append(topicsList, itypes.DecreaseLiquidityTopic)
		case "Transfer":
			topicsList = append(topicsList, itypes.TransferTopic)
		}
	}
	return topicsList
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
	callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(l.BlockNumber))}

	if len(l.Topics) < 3 {
		return
	}

	sender := util.ExtractAddressFromLogTopic(l.Topics[1])
	receiver := util.ExtractAddressFromLogTopic(l.Topics[2])

	// Check if we have enough data to retrieve amount of token being minted
	if len(l.Data) < 32 {
		return
	}

	amount := util.ExtractIntFromBytes(l.Data[:32])

	if len(amount.Bits()) == 0 {
		return
	}

	erc, client := r.da.GetERC20(l.Address)

	tokenDecimals, err := r.da.GetERC20Decimals(erc, client, callopts)
	if util.IsExecutionReverted(err) {
		// Non ERC-20 contract
		tokenDecimals = 0
	} else {
		if util.IsEthErr(err) {
			return
		}
		util.ENOK(err)
	}

	formattedAmount := util.DivideBy10pow(amount, tokenDecimals)
	// token0Price, token1Price, amountusd, tokenMeta := r.da.GetPricesForBlock(r.dbconn.ChainID, callopts, token0, token1, formattedAmount0, formattedAmount1)//TODO
	amountUSD := 0
	tokenPrice := 0

	transfer := itypes.Transfer{
		Type:        "transfer",
		Network:     r.dbconn.ChainID,
		LogIdx:      l.Index,
		Transaction: l.TxHash,
		Time:        bm.Time,
		Height:      l.BlockNumber,
		Token:       l.Address,
		Sender:      sender,
		Receiver:    receiver,
		Amount:      formattedAmount,
		AmountUSD:   0, // TODO
	}
	mt.Lock()
	defer mt.Unlock()
	is0Nan := math.IsInf(float64(tokenPrice), 0)
	if amountUSD > -1 && !is0Nan {
		*items = append(*items, transfer)
		bm.TransferLogs++
		bm.TotalLogs++
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
	if len(l.Data) < 64 {
		return
	}

	am0 := util.ExtractIntFromBytes(l.Data[:32])
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
