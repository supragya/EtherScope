package uniswapv2

import (
	"math/big"

	"github.com/supragya/EtherScope/libs/util"
	"github.com/supragya/EtherScope/services/ethrpc"
	"github.com/supragya/EtherScope/services/instrumentation"
	itypes "github.com/supragya/EtherScope/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type UniswapV3Processor struct {
	Topics map[common.Hash]itypes.ProcessingType
	EthRPC ethrpc.EthRPC
}

func (n *UniswapV3Processor) ProcessUniV3Mint(
	l types.Log,
	items []interface{},
	idx int,
	blockTime uint64,
) error {
	prcType, ok := n.Topics[itypes.UniV3MintTopic]
	if !ok {
		return nil
	}

	callopts := util.GetBlockCallOpts(l.BlockNumber)

	// Test if the contract is a UniswapV2 type contract
	if !n.isUniswapV3Pair(l.Address, callopts) {
		return nil
	}

	mint := itypes.Mint{
		Type:           "uniswapv3mint",
		ProcessingType: prcType,
		LogIdx:         l.Index,
		Transaction:    l.TxHash,
		Time:           blockTime,
		Height:         l.BlockNumber,
		Sender:         common.Address{}, // To be filled if UserRequested processing
		TxSender:       common.Address{}, // To be filled if UserRequested processing
		PairContract:   l.Address,
		Token0:         common.Address{}, // To be filled if PricingEngineRequest processing
		Token1:         common.Address{}, // To be filled if PricingEngineRequest processing
		Amount0:        nil,              // To be filled if UserRequested processing
		Amount1:        nil,              // To be filled if UserRequested processing
		Reserve0:       nil,              // To be filled if PricingEngineRequest processing
		Reserve1:       nil,              // To be filled if PricingEngineRequest processing
		AmountUSD:      nil,              // To be filled if UserRequested processing
		Price0:         nil,              // To be filled later by Pricing Engine
		Price1:         nil,              // To be filled later by Pricing Engine
	}

	// Fill up the fields needed by pricing engine
	var err error
	mint.Token0, mint.Token1, err = n.EthRPC.GetTokensUniV3(l.Address, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	token0decimals, err := n.EthRPC.GetERC20Decimals(mint.Token0, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	token1decimals, err := n.EthRPC.GetERC20Decimals(mint.Token1, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	reserves, err := n.EthRPC.GetERC20Balances([]itypes.Tuple2[common.Address, common.Address]{
		{l.Address, mint.Token0}, {l.Address, mint.Token1},
	}, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	mint.Reserve0 = util.DivideBy10pow(reserves[0], token0decimals)
	mint.Reserve1 = util.DivideBy10pow(reserves[1], token1decimals)

	// Fill up the fields needed by user
	if prcType == itypes.UserRequested {
		ok, sender, _, am0, am1 := InfoUniV3Mint(l)
		if !ok {
			return nil
		}

		txSender, err := n.EthRPC.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
		if util.IsEthErr(err) {
			return nil
		}
		if err != nil {
			return err
		}

		mint.Sender = sender
		mint.TxSender = txSender
		mint.Amount0 = util.DivideBy10pow(am0, token0decimals)
		mint.Amount1 = util.DivideBy10pow(am1, token1decimals)
	}

	items[idx] = &mint
	instrumentation.MintV3Processed.Inc()
	return nil
}

func (n *UniswapV3Processor) ProcessUniV3Burn(
	l types.Log,
	items []interface{},
	idx int,
	blockTime uint64,
) error {
	prcType, ok := n.Topics[itypes.UniV3BurnTopic]
	if !ok {
		return nil
	}

	callopts := util.GetBlockCallOpts(l.BlockNumber)

	// Test if the contract is a UniswapV2 type contract
	if !n.isUniswapV3Pair(l.Address, callopts) {
		return nil
	}

	burn := itypes.Burn{
		Type:           "uniswapv3burn",
		ProcessingType: prcType,
		LogIdx:         l.Index,
		Transaction:    l.TxHash,
		Time:           blockTime,
		Height:         l.BlockNumber,
		Sender:         common.Address{}, // To be filled if UserRequested processing
		TxSender:       common.Address{}, // To be filled if UserRequested processing
		PairContract:   l.Address,
		Token0:         common.Address{}, // To be filled if PricingEngineRequest processing
		Token1:         common.Address{}, // To be filled if PricingEngineRequest processing
		Amount0:        nil,              // To be filled if UserRequested processing
		Amount1:        nil,              // To be filled if UserRequested processing
		Reserve0:       nil,              // To be filled if PricingEngineRequest processing
		Reserve1:       nil,              // To be filled if PricingEngineRequest processing
		AmountUSD:      nil,              // To be filled if UserRequested processing
		Price0:         nil,              // To be filled later by Pricing Engine
		Price1:         nil,              // To be filled later by Pricing Engine
	}

	// Fill up the fields needed by pricing engine
	var err error
	burn.Token0, burn.Token1, err = n.EthRPC.GetTokensUniV3(l.Address, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	token0decimals, err := n.EthRPC.GetERC20Decimals(burn.Token0, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	token1decimals, err := n.EthRPC.GetERC20Decimals(burn.Token1, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	reserves, err := n.EthRPC.GetERC20Balances([]itypes.Tuple2[common.Address, common.Address]{
		{l.Address, burn.Token0}, {l.Address, burn.Token1},
	}, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	burn.Reserve0 = util.DivideBy10pow(reserves[0], token0decimals)
	burn.Reserve1 = util.DivideBy10pow(reserves[1], token1decimals)

	// Fill up the fields needed by user
	if prcType == itypes.UserRequested {
		ok, sender, _, am0, am1 := InfoUniV3Burn(l)
		if !ok {
			return nil
		}

		txSender, err := n.EthRPC.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
		if util.IsEthErr(err) {
			return nil
		}
		if err != nil {
			return err
		}

		burn.Sender = sender
		burn.TxSender = txSender
		burn.Amount0 = util.DivideBy10pow(am0, token0decimals)
		burn.Amount1 = util.DivideBy10pow(am1, token1decimals)
	}

	items[idx] = &burn
	instrumentation.BurnV3Processed.Inc()
	return nil
}

func (n *UniswapV3Processor) ProcessUniV3Swap(
	l types.Log,
	items []interface{},
	idx int,
	blockTime uint64,
) error {
	prcType, ok := n.Topics[itypes.UniV3SwapTopic]
	if !ok {
		return nil
	}

	callopts := util.GetBlockCallOpts(l.BlockNumber)

	// Test if the contract is a UniswapV2 type contract
	if !n.isUniswapV3Pair(l.Address, callopts) {
		return nil
	}

	swap := itypes.Swap{
		Type:           "uniswapv3swap",
		ProcessingType: prcType,
		LogIdx:         l.Index,
		Transaction:    l.TxHash,
		Time:           blockTime,
		Height:         l.BlockNumber,
		Sender:         common.Address{}, // To be filled if UserRequested processing
		TxSender:       common.Address{}, // To be filled if UserRequested processing
		Receiver:       common.Address{}, // To be filled if UserRequested processing
		PairContract:   l.Address,
		Token0:         common.Address{}, // To be filled if PricingEngineRequest processing
		Token1:         common.Address{}, // To be filled if PricingEngineRequest processing
		Amount0:        nil,              // To be filled if UserRequested processing
		Amount1:        nil,              // To be filled if UserRequested processing
		Reserve0:       nil,              // To be filled if PricingEngineRequest processing
		Reserve1:       nil,              // To be filled if PricingEngineRequest processing
		AmountUSD:      nil,              // To be filled if UserRequested processing
		Price0:         nil,              // To be filled later by Pricing Engine
		Price1:         nil,              // To be filled later by Pricing Engine
	}

	// Fill up the fields needed by pricing engine
	var err error
	swap.Token0, swap.Token1, err = n.EthRPC.GetTokensUniV3(l.Address, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	token0decimals, err := n.EthRPC.GetERC20Decimals(swap.Token0, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	token1decimals, err := n.EthRPC.GetERC20Decimals(swap.Token1, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	reserves, err := n.EthRPC.GetERC20Balances([]itypes.Tuple2[common.Address, common.Address]{
		{l.Address, swap.Token0}, {l.Address, swap.Token1},
	}, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	swap.Reserve0 = util.DivideBy10pow(reserves[0], token0decimals)
	swap.Reserve1 = util.DivideBy10pow(reserves[1], token1decimals)

	// Fill up the fields needed by user
	if prcType == itypes.UserRequested {
		ok, sender, receiver, am0, am1 := InfoUniV3Swap(l)
		if !ok {
			return nil
		}

		txSender, err := n.EthRPC.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
		if util.IsEthErr(err) {
			return nil
		}
		if err != nil {
			return err
		}

		swap.Sender = sender
		swap.Receiver = receiver
		swap.TxSender = txSender
		swap.Amount0 = util.DivideBy10pow(am0, token0decimals)
		swap.Amount1 = util.DivideBy10pow(am1, token1decimals)
	}

	items[idx] = &swap

	instrumentation.SwapV3Processed.Inc()
	return nil
}

func (n *UniswapV3Processor) isUniswapV3Pair(address common.Address,
	callopts *bind.CallOpts) bool {
	_, _, err := n.EthRPC.GetTokensUniV3(address, callopts)
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

func InfoUniV3Mint(l types.Log) (hasSufficientData bool,
	sender common.Address,
	amount *big.Int,
	amount0 *big.Int,
	amount1 *big.Int) {
	if !util.HasSufficientData(l, 4, 128) {
		return false,
			common.Address{},
			big.NewInt(0),
			big.NewInt(0),
			big.NewInt(0)
	}
	return true,
		util.ExtractAddressFromLogTopic(l.Topics[1]),
		util.ExtractIntFromBytes(l.Data[32:64]),
		util.ExtractIntFromBytes(l.Data[64:96]),
		util.ExtractIntFromBytes(l.Data[96:128])
}

func InfoUniV3Burn(l types.Log) (hasSufficientData bool,
	sender common.Address,
	amount *big.Int,
	amount0 *big.Int,
	amount1 *big.Int) {
	if !util.HasSufficientData(l, 4, 96) {
		return false,
			common.Address{},
			big.NewInt(0),
			big.NewInt(0),
			big.NewInt(0)
	}
	return true,
		util.ExtractAddressFromLogTopic(l.Topics[1]),
		util.ExtractIntFromBytes(l.Data[:32]),
		util.ExtractIntFromBytes(l.Data[32:64]),
		util.ExtractIntFromBytes(l.Data[64:96])
}

func InfoUniV3Swap(l types.Log) (hasSufficientData bool,
	sender common.Address,
	recipient common.Address,
	amount0 *big.Int,
	amount1 *big.Int) {
	if !util.HasSufficientData(l, 3, 160) {
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
