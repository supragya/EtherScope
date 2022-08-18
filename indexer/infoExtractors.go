package indexer

import (
	"math/big"

	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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

func InfoUniV2Mint(l types.Log) (hasSufficientData bool,
	sender common.Address,
	amount0 *big.Int,
	amount1 *big.Int) {
	if !HasSufficientData(l, 2, 64) {
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
	if !HasSufficientData(l, 3, 64) {
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
	if !HasSufficientData(l, 1, 128) {
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

	am0, am1 := big.NewInt(0), big.NewInt(0)
	if am0In.Cmp(big.NewInt(0)) == 0 {
		am0 = am0.Neg(am0Out)
		am1 = am1In
	} else {
		am0 = am0In
		am1 = am1.Neg(am1Out)
	}

	return true, am0, am1
}

func InfoUniV3Mint(l types.Log) (hasSufficientData bool,
	tokenID *big.Int,
	amount0 *big.Int,
	amount1 *big.Int) {
	if !HasSufficientData(l, 2, 64) {
		return false,
			big.NewInt(0),
			big.NewInt(0),
			big.NewInt(0)
	}
	return true,
		util.ExtractIntFromBytes(l.Topics[1][:]),
		util.ExtractIntFromBytes(l.Data[32:64]),
		util.ExtractIntFromBytes(l.Data[64:96])
}
