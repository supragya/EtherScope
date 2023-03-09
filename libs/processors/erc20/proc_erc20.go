package erc20

import (
	"math/big"

	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/ethrpc"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/instrumentation"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC20Processor struct {
	Topics map[common.Hash]itypes.ProcessingType
	EthRPC ethrpc.EthRPC
}

func (n *ERC20Processor) ProcessERC20Transfer(
	l types.Log,
	items []interface{},
	idx int,
	blockTime uint64,
) error {
	prcType, ok := n.Topics[itypes.ERC20TransferTopic]
	if !ok {
		return nil
	}

	ok, sender, recv, amt := InfoTransfer(l)
	if !ok {
		return nil
	}

	callopts := util.GetBlockCallOpts(l.BlockNumber)
	ok, formattedAmount := n.GetFormattedAmount(amt, callopts, l.Address)
	if !ok {
		return nil
	}

	txSender, err := n.EthRPC.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return nil
	}
	util.ENOK(err)

	transfer := itypes.Transfer{
		Type:           "erc20transfer",
		ProcessingType: prcType,
		// Network:     r.dbconn.ChainID,
		LogIdx:      l.Index,
		Transaction: l.TxHash,
		// Time:        bm.Time,
		Height:   l.BlockNumber,
		Token:    l.Address,
		Sender:   sender,
		TxSender: txSender,
		Receiver: recv,
		Amount:   formattedAmount,
		// AmountUSD:           tokenPrice.Price,
		// PriceDerivationMeta: tokenPrice,
	}

	items[idx] = &transfer
	instrumentation.TfrProcessed.Inc()
	return nil
}

func InfoTransfer(l types.Log) (hasSufficientData bool,
	sender common.Address,
	receiver common.Address,
	amount *big.Int) {
	if !util.HasSufficientData(l, 3, 32) {
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

func (n *ERC20Processor) GetFormattedAmount(amount *big.Int,
	callopts *bind.CallOpts,
	erc20Address common.Address) (ok bool,
	formattedAmount *big.Float) {
	tokenDecimals, err := n.EthRPC.GetERC20Decimals(erc20Address, callopts)
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
