package indexer

import (
	"math/big"
	"sync"

	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/instrumentation"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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

	ok, sender, am0, am1 := InfoUniV2Mint(l)
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
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

	mint := itypes.Mint{
		Type:         "mint",
		Network:      r.dbconn.ChainID,
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       sender,
		Receiver:     common.Address{},
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
	instrumentation.MintV2Processed.Inc()
}

func (r *RealtimeIndexer) processBurn(
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

	ok, sender, recipient, am0, am1 := InfoUniV2Burn(l)
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
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

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
	instrumentation.BurnV2Processed.Inc()
}

func (r *RealtimeIndexer) processUniV2Swap(
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

	ok, am0, am1 := InfoUniV2Swap(l)
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
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

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

	AddToSynopsis(mt, bm, swap, items, "swap", true)
	instrumentation.SwapV2Processed.Inc()
}

func (r *RealtimeIndexer) isUniswapV2Pair(address common.Address,
	callopts *bind.CallOpts) bool {
	_, _, err := r.da.GetTokensUniV2(address, callopts)
	if err == nil {
		return true
	}

	// Execution Revert: Could be a non uniswap contract (like AAVE V2). Example log:
	// https://etherscan.io/tx/0x65ed6ba09f2a22805b772ff607f81fa4bb5d93ce287ecf05ab5ad97cab34c97c#eventlog logIdx 180
	// not handled currently
	if !util.IsEthErr(err) {
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

	token0Decimals, err = r.da.GetERC20Decimals(erc0, client0, t0, callopts)
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

	token1Decimals, err = r.da.GetERC20Decimals(erc1, client1, t1, callopts)
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
