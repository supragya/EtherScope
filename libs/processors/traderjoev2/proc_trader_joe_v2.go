package traderjoev2

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/supragya/EtherScope/libs/util"
	"github.com/supragya/EtherScope/services/ethrpc"
	"github.com/supragya/EtherScope/services/instrumentation"
	itypes "github.com/supragya/EtherScope/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TraderJoeProcessor struct {
	Topics map[common.Hash]itypes.ProcessingType
	EthRPC ethrpc.EthRPC
}

type TraderJoeExtraData struct {
	Id                    *big.Int
	VolatilityAccumulated *big.Int
	Fees                  *big.Int
}

func (n *TraderJoeProcessor) ProcessTokenExchange(
	l types.Log,
	items []interface{},
	idx int,
	blockTime uint64,
) error {
	prcType, ok := n.Topics[itypes.TraderJoeV2SwapTopic]
	if !ok {
		return nil
	}

	callopts := util.GetBlockCallOpts(l.BlockNumber)

	// Test if the contract is a TraderJoePool type contract
	if isTJPool, err := n.isTraderJoePool(l.Address, callopts); err != nil {
		return err
	} else {
		if !isTJPool {
			return nil
		}
	}

	swap := itypes.Swap{
		Type:           "traderjoev2swap",
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
	hasSufficientData, senderAddress, recipientAddress, id, swapForY, amountIn, amountOut, volatilityAccumulated, fees := InfoTraderJoeV2Swap(l)
	if !hasSufficientData {
		return errors.New(fmt.Sprintf("data error: log doesn't have sufficient data. Contract Address: %+v", l.Address))
	}
	extraData := TraderJoeExtraData{id, volatilityAccumulated, fees}

	TokenX, err := n.EthRPC.GetTraderJoeTokenX(l.Address, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}
	TokenY, err := n.EthRPC.GetTraderJoeTokenY(l.Address, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	tokenXdecimals, err := n.EthRPC.GetERC20Decimals(TokenX, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	tokenYdecimals, err := n.EthRPC.GetERC20Decimals(TokenY, callopts)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if !swapForY {
		swap.Token0 = TokenY
		swap.Token1 = TokenX
		swap.Amount0 = util.DivideBy10pow(amountIn, tokenYdecimals)
		swap.Amount1 = util.DivideBy10pow(amountOut, tokenXdecimals)
	} else {
		swap.Token0 = TokenX
		swap.Token1 = TokenY
		swap.Amount0 = util.DivideBy10pow(amountIn, tokenXdecimals)
		swap.Amount1 = util.DivideBy10pow(amountOut, tokenYdecimals)
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

	if !swapForY {
		swap.Reserve0 = util.DivideBy10pow(reserves[0], tokenYdecimals)
		swap.Reserve1 = util.DivideBy10pow(reserves[1], tokenXdecimals)
	} else {
		swap.Reserve0 = util.DivideBy10pow(reserves[0], tokenXdecimals)
		swap.Reserve1 = util.DivideBy10pow(reserves[1], tokenYdecimals)
	}

	txSender, err := n.EthRPC.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return nil
	}
	if err != nil {
		return err
	}

	swap.Sender = senderAddress
	swap.Receiver = recipientAddress
	swap.TxSender = txSender
	swap.ExtraData = extraData

	items[idx] = &swap

	instrumentation.TraderJoeV2SwapProcessed.Inc()
	return nil
}

func (n *TraderJoeProcessor) isTraderJoePool(address common.Address,
	callopts *bind.CallOpts) (bool, error) {
	_, err := n.EthRPC.GetTraderJoeTokenX(address, callopts)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func InfoTraderJoeV2Swap(l types.Log) (hasSufficientData bool,
	senderAddress common.Address,
	recipientAddress common.Address,
	id *big.Int,
	swapForY bool,
	amountIn *big.Int,
	amountOut *big.Int,
	volatilityAccumulated *big.Int,
	fees *big.Int) {
	if !util.HasSufficientData(l, 4, 160) {
		return false,
			common.Address{},
			common.Address{},
			big.NewInt(0),
			false,
			big.NewInt(0),
			big.NewInt(0),
			big.NewInt(0),
			big.NewInt(0)
	}
	senderAddress = util.ExtractAddressFromLogTopic(l.Topics[1])
	recipientAddress = util.ExtractAddressFromLogTopic(l.Topics[2])
	id = util.ExtractIntFromBytes(l.Topics[3].Bytes())

	swapForY = util.ExtractIntFromBytes(l.Data[0:32]).Int64() == 1
	amountIn = util.ExtractIntFromBytes(l.Data[32:64])
	amountOut = util.ExtractIntFromBytes(l.Data[64:96])
	volatilityAccumulated = util.ExtractIntFromBytes(l.Data[96:128])
	fees = util.ExtractIntFromBytes(l.Data[128:160])

	return true,
		senderAddress,
		recipientAddress,
		id,
		swapForY,
		amountIn,
		amountOut,
		volatilityAccumulated,
		fees
}
