package uniswapv2

import (
	"math/big"
	"sync"

	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/ethrpc"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type UniswapV2Processor struct {
	Topics map[common.Hash]itypes.ProcessingType
	EthRPC ethrpc.EthRPC
}

func (n *UniswapV2Processor) ProcessUniV2Mint(
	l types.Log,
	items []interface{},
	idx int,
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	prcType, ok := n.Topics[itypes.UniV2MintTopic]
	if !ok {
		return
	}

	callopts := util.GetBlockCallOpts(l.BlockNumber)

	// Test if the contract is a UniswapV2 type contract
	if !n.isUniswapV2Pair(l.Address, callopts) {
		return
	}

	mint := itypes.Mint{
		Type:         "uniswapv2mint",
		LogIdx:       l.Index,
		Transaction:  l.TxHash,
		Time:         bm.Time,
		Height:       l.BlockNumber,
		Sender:       common.Address{}, // To be filled if UserRequested processing
		TxSender:     common.Address{}, // To be filled if UserRequested processing
		PairContract: l.Address,
		Token0:       common.Address{}, // To be filled if PricingEngineRequest processing
		Token1:       common.Address{}, // To be filled if PricingEngineRequest processing
		Amount0:      nil,              // To be filled if UserRequested processing
		Amount1:      nil,              // To be filled if UserRequested processing
		Reserve0:     nil,              // To be filled if PricingEngineRequest processing
		Reserve1:     nil,              // To be filled if PricingEngineRequest processing
		AmountUSD:    nil,              // To be filled if UserRequested processing
		Price0:       nil,              // To be filled later by Pricing Engine
		Price1:       nil,              // To be filled later by Pricing Engine
	}

	// Fill up the fields needed by pricing engine
	var err error
	mint.Token0, mint.Token1, err = n.EthRPC.GetTokensUniV2(l.Address, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token0decimals, err := n.EthRPC.GetERC20Decimals(mint.Token0, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	token1decimals, err := n.EthRPC.GetERC20Decimals(mint.Token1, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	reserves, err := n.EthRPC.GetERC20Balances([]util.Tuple2[common.Address, common.Address]{
		{l.Address, mint.Token0}, {l.Address, mint.Token1},
	}, callopts)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	mint.Reserve0 = util.DivideBy10pow(reserves[0], token0decimals)
	mint.Reserve1 = util.DivideBy10pow(reserves[1], token1decimals)

	// Fill up the fields needed by user
	if prcType == itypes.UserRequested {
		ok, sender, am0, am1 := InfoUniV2Mint(l)
		if !ok {
			return
		}

		txSender, err := n.EthRPC.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
		if util.IsEthErr(err) {
			return
		}
		util.ENOK(err)

		mint.Sender = sender
		mint.TxSender = txSender
		mint.Amount0 = util.DivideBy10pow(am0, token0decimals)
		mint.Amount1 = util.DivideBy10pow(am1, token1decimals)
	}

	items[idx] = mint
	// instrumentation.MintV2Processed.Inc()
}

// func (r *UniswapV2Processor) processUniV2Burn(
// 	l types.Log,
// 	items *[]interface{},
// 	bm *itypes.BlockSynopsis,
// 	mt *sync.Mutex,
// ) {
// 	callopts := GetBlockCallOpts(l.BlockNumber)

// 	// Test if the contract is a UniswapV2 type contract
// 	if !r.isUniswapV2Pair(l.Address, callopts) {
// 		return
// 	}

// 	ok, sender, _, am0, am1 := InfoUniV2Burn(l)
// 	if !ok {
// 		return
// 	}

// 	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV2(am0, am1, callopts, l.Address)
// 	if !ok {
// 		return
// 	}

// 	// Assumed infallible since if err != nil, this code should not be reachable
// 	// due to above condition
// 	t0, t1, _ := r.da.GetTokensUniV2(l.Address, callopts)

// 	reserves, err := r.da.GetERC20Balances([]util.Tuple2[common.Address, common.Address]{
// 		{l.Address, t0}, {l.Address, t1},
// 	}, callopts)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, l, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

// 	txSender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	burn := itypes.Burn{
// 		Type:         "uniswapv2burn",
// 		Network:      r.dbconn.ChainID,
// 		LogIdx:       l.Index,
// 		Transaction:  l.TxHash,
// 		Time:         bm.Time,
// 		Height:       l.BlockNumber,
// 		Sender:       sender,
// 		TxSender:     txSender,
// 		PairContract: l.Address,
// 		Token0:       t0,
// 		Token1:       t1,
// 		Amount0:      f0,
// 		Amount1:      f1,
// 		Reserve0:     util.DivideBy10pow(reserves[0].Second, t0d),
// 		Reserve1:     util.DivideBy10pow(reserves[1].Second, t1d),
// 		AmountUSD:    amountusd,
// 		Price0:       token0Price,
// 		Price1:       token1Price,
// 	}

// 	AddToSynopsis(mt, bm, burn, items, "burn", true)
// 	instrumentation.BurnV2Processed.Inc()
// }

// func (r *UniswapV2Processor) processUniV2Swap(
// 	l types.Log,
// 	items *[]interface{},
// 	bm *itypes.BlockSynopsis,
// 	mt *sync.Mutex,
// ) {
// 	callopts := GetBlockCallOpts(l.BlockNumber)

// 	// Test if the contract is a UniswapV2 type contract
// 	if !r.isUniswapV2Pair(l.Address, callopts) {
// 		return
// 	}

// 	ok, am0, am1 := InfoUniV2Swap(l)
// 	if !ok {
// 		return
// 	}

// 	ok, f0, f1, t0d, t1d := r.GetFormattedAmountsUniV2(am0, am1, callopts, l.Address)
// 	if !ok {
// 		return
// 	}

// 	// Assumed infallible since if err != nil, this code should not be reachable
// 	// due to above condition
// 	t0, t1, _ := r.da.GetTokensUniV2(l.Address, callopts)

// 	// if t0 != common.HexToAddress("0x79291a9d692df95334b1a0b3b4ae6bc606782f8c") || t1 != common.HexToAddress("0x79291a9d692df95334b1a0b3b4ae6bc606782f8c") {
// 	// 	return
// 	// }

// 	reserves, err := r.da.GetERC20Balances([]util.Tuple2[common.Address, common.Address]{
// 		{l.Address, t0}, {l.Address, t1},
// 	}, callopts)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	token0Price, token1Price, amountusd := r.da.GetRates2Tokens(callopts, l, t0, t1, big.NewFloat(1.0).Abs(f0), big.NewFloat(1.0).Abs(f1))

// 	txSender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
// 	if util.IsEthErr(err) {
// 		return
// 	}
// 	util.ENOK(err)

// 	swap := itypes.Swap{
// 		Type:         "uniswapv2swap",
// 		Network:      r.dbconn.ChainID,
// 		LogIdx:       l.Index,
// 		Transaction:  l.TxHash,
// 		Time:         bm.Time,
// 		Height:       l.BlockNumber,
// 		Sender:       util.ExtractAddressFromLogTopic(l.Topics[1]),
// 		TxSender:     txSender,
// 		Receiver:     util.ExtractAddressFromLogTopic(l.Topics[2]),
// 		PairContract: l.Address,
// 		Token0:       t0,
// 		Token1:       t1,
// 		Amount0:      f0,
// 		Amount1:      f1,
// 		Reserve0:     util.DivideBy10pow(reserves[0].Second, t0d),
// 		Reserve1:     util.DivideBy10pow(reserves[1].Second, t1d),
// 		AmountUSD:    amountusd,
// 		Price0:       token0Price,
// 		Price1:       token1Price,
// 	}

// 	AddToSynopsis(mt, bm, swap, items, "swap", true)
// 	instrumentation.SwapV2Processed.Inc()
// }

func (n *UniswapV2Processor) isUniswapV2Pair(address common.Address,
	callopts *bind.CallOpts) bool {
	_, _, err := n.EthRPC.GetTokensUniV2(address, callopts)
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

func (n *UniswapV2Processor) GetFormattedAmountsUniV2(amount0 *big.Int,
	amount1 *big.Int,
	callopts *bind.CallOpts,
	address common.Address) (ok bool,
	formattedAmount0 *big.Float,
	formattedAmount1 *big.Float,
	token0Decimals uint8,
	token1Decimals uint8) {
	t0, t1, err := n.EthRPC.GetTokensUniV2(address, callopts)
	if err != nil {
		return false,
			big.NewFloat(0.0),
			big.NewFloat(0.0),
			0,
			0
	}

	token0Decimals, err = n.EthRPC.GetERC20Decimals(t0, callopts)
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

	token1Decimals, err = n.EthRPC.GetERC20Decimals(t1, callopts)
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

func InfoUniV2Mint(l types.Log) (hasSufficientData bool,
	sender common.Address,
	amount0 *big.Int,
	amount1 *big.Int) {
	if !util.HasSufficientData(l, 2, 64) {
		return false,
			common.Address{},
			big.NewInt(0),
			big.NewInt(0)
	}
	return true,
		util.ExtractAddressFromLogTopic(l.Topics[1]),
		util.ExtractIntFromBytes(l.Data[:32]),
		util.ExtractIntFromBytes(l.Data[32:64])
}

func InfoUniV2Burn(l types.Log) (hasSufficientData bool,
	sender common.Address,
	recipient common.Address,
	amount0 *big.Int,
	amount1 *big.Int) {
	if !util.HasSufficientData(l, 3, 64) {
		return false,
			common.Address{},
			common.Address{},
			big.NewInt(0),
			big.NewInt(0)
	}
	return true,
		util.ExtractAddressFromLogTopic(l.Topics[1]),
		util.ExtractAddressFromLogTopic(l.Topics[2]),
		util.ExtractIntFromBytes(l.Data[:32]),
		util.ExtractIntFromBytes(l.Data[32:64])
}

func InfoUniV2Swap(l types.Log) (hasSufficientData bool,
	amount0 *big.Int,
	amount1 *big.Int) {
	if !util.HasSufficientData(l, 3, 128) {
		return false,
			big.NewInt(0),
			big.NewInt(0)
	}

	var (
		am0In  = util.ExtractIntFromBytes(l.Data[0:32])
		am1In  = util.ExtractIntFromBytes(l.Data[32:64])
		am0Out = util.ExtractIntFromBytes(l.Data[64:96])
		am1Out = util.ExtractIntFromBytes(l.Data[96:128])
	)

	return true, big.NewInt(0).Sub(am0Out, am0In), big.NewInt(0).Sub(am1Out, am1In)
}
