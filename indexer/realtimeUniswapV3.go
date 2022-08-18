package indexer

import (
	"math"
	"math/big"
	"sync"

	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

func (r *RealtimeIndexer) processMintV3(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	callopts := GetBlockCallOpts(l.BlockNumber)

	// Test if the contract is a UniswapV3 NFT type contract
	if !r.isUniswapV3NFT(l.Address, callopts) {
		return
	}

	sender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	ok, tokenID, am0, am1 := InfoUniV3Mint(l)
	if !ok {
		return
	}

	// Test if the contract is a UniswapV3NFT type contract
	t0, t1, err := r.da.GetTokensUniV3(l.Address, tokenID, callopts)

	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV3(am0, am1, tokenID, callopts, l.Address)
	if !ok {
		return
	}

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

func (r *RealtimeIndexer) isUniswapV3NFT(address common.Address,
	callopts *bind.CallOpts) bool {
	_, _, err := r.da.GetTokensUniV3(address, big.NewInt(0), callopts)
	if err != nil {
		return true
	}

	if !util.IsExecutionReverted(err) {
		util.ENOKS(2, err)
	}
	return false
}

// TODO: refactor this with GetFormattedAmountsUniV2
func (r *RealtimeIndexer) GetFormattedAmountsUniV3(amount0 *big.Int,
	amount1 *big.Int,
	tokenID *big.Int,
	callopts *bind.CallOpts,
	address common.Address) (ok bool,
	formattedAmount0 *big.Float,
	formattedAmount1 *big.Float,
	token0Decimals uint8,
	token1Decimals uint8) {
	t0, t1, err := r.da.GetTokensUniV3(address, tokenID, callopts)
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
