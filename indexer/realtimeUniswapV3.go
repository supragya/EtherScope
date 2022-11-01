package indexer

import (
	"math/big"
	"sync"

	"github.com/Blockpour/Blockpour-Geth-Indexer/instrumentation"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (r *RealtimeIndexer) processUniV3Mint(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	callopts := GetBlockCallOpts(l.BlockNumber)
	// Test if the contract is a UniswapV3Pair type contract
	if !r.isUniswapV3(l.Address, callopts) {
		return
	}

	sender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	ok, _, am0, am1 := InfoUniV3Mint(l)
	if !ok {
		return
	}

	// Test if the contract is a UniswapV3NFT type contract
	t0, t1, err := r.da.GetTokensUniV3(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV3(am0, am1, callopts, l.Address)
	if !ok {
		return
	}

	reserves, err := r.da.GetERC20Balances([]util.Tuple2[common.Address, common.Address]{
		{l.Address, t0}, {l.Address, t1},
	}, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, l, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

	mint := itypes.Mint{
		Type:         "uniswapv3mint",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       sender,
		TxSender:     sender,
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
	}

	AddToSynopsis(mt, bm, mint, items, "mint", true)
	instrumentation.MintV3Processed.Inc()
}

func (r *RealtimeIndexer) processUniV3Burn(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	callopts := GetBlockCallOpts(l.BlockNumber)

	// Test if the contract is a UniswapV3Pair type contract
	if !r.isUniswapV3(l.Address, callopts) {
		return
	}

	sender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	ok, _, am0, am1 := InfoUniV3Burn(l)
	if !ok {
		return
	}

	t0, t1, err := r.da.GetTokensUniV3(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV3(am0, am1, callopts, l.Address)
	if !ok {
		return
	}

	reserves, err := r.da.GetERC20Balances([]util.Tuple2[common.Address, common.Address]{
		{l.Address, t0}, {l.Address, t1},
	}, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, l, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

	burn := itypes.Burn{
		Type:         "uniswapv3burn",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       sender,
		TxSender:     sender,
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
	}

	AddToSynopsis(mt, bm, burn, items, "burn", true)
	instrumentation.BurnV3Processed.Inc()
}

func (r *RealtimeIndexer) processUniV3Swap(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	callopts := GetBlockCallOpts(l.BlockNumber)

	// Test if the contract is a UniswapV3 NFT type contract
	if !r.isUniswapV3(l.Address, callopts) {
		return
	}

	ok, sender, receiver, am0, am1 := InfoUniV3Swap(l)
	if !ok {
		return
	}

	// Test if the contract is a UniswapV3NFT type contract
	t0, t1, err := r.da.GetTokensUniV3(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV3(am0, am1, callopts, l.Address)
	if !ok {
		return
	}

	reserves, err := r.da.GetERC20Balances([]util.Tuple2[common.Address, common.Address]{
		{l.Address, t0}, {l.Address, t1},
	}, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, l, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

	txSender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	swap := itypes.Swap{
		Type:         "uniswapv3swap",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       sender,
		TxSender:     txSender,
		Receiver:     receiver,
		PairContract: l.Address,
		Token0:       t0,
		Token1:       t1,
		Amount0:      f0,
		Amount1:      f1,
		Reserve0:     util.DivideBy10pow(reserves[0].Second, t0d),
		Reserve1:     util.DivideBy10pow(reserves[0].Second, t1d),
		AmountUSD:    amountusd,
		Price0:       token0Price,
		Price1:       token1Price,
	}

	AddToSynopsis(mt, bm, swap, items, "swap", true)
	instrumentation.SwapV3Processed.Inc()
}

func (r *RealtimeIndexer) isUniswapV3(address common.Address,
	callopts *bind.CallOpts) bool {
	_, _, err := r.da.GetTokensUniV3(address, callopts)
	if err == nil {
		return true
	}

	if !util.IsEthErr(err) {
		util.ENOKS(2, err)
	}
	return false
}

func (r *RealtimeIndexer) isUniswapV3NFT(address common.Address,
	callopts *bind.CallOpts) bool {
	_, _, err := r.da.GetTokensUniV3NFT(address, big.NewInt(1), callopts)
	if err == nil {
		return true
	}

	if !util.IsEthErr(err) {
		util.ENOKS(2, err)
	}
	return false
}

func (r *RealtimeIndexer) GetFormattedAmountsUniV3NFT(amount0 *big.Int,
	amount1 *big.Int,
	tokenID *big.Int,
	callopts *bind.CallOpts,
	address common.Address) (ok bool,
	formattedAmount0 *big.Float,
	formattedAmount1 *big.Float,
	token0Decimals uint8,
	token1Decimals uint8) {
	t0, t1, err := r.da.GetTokensUniV3NFT(address, tokenID, callopts)
	if err != nil {
		return false,
			big.NewFloat(0.0),
			big.NewFloat(0.0),
			0,
			0
	}

	token0Decimals, err = r.da.GetERC20Decimals(t0, callopts)
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

	token1Decimals, err = r.da.GetERC20Decimals(t1, callopts)
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

func (r *RealtimeIndexer) GetFormattedAmountsUniV3(amount0 *big.Int,
	amount1 *big.Int,
	callopts *bind.CallOpts,
	address common.Address) (ok bool,
	formattedAmount0 *big.Float,
	formattedAmount1 *big.Float,
	token0Decimals uint8,
	token1Decimals uint8) {
	t0, t1, err := r.da.GetTokensUniV3(address, callopts)
	if err != nil {
		return false,
			big.NewFloat(0.0),
			big.NewFloat(0.0),
			0,
			0
	}

	token0Decimals, err = r.da.GetERC20Decimals(t0, callopts)
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

	token1Decimals, err = r.da.GetERC20Decimals(t1, callopts)
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
