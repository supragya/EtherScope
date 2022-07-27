package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	MintTopic     common.Hash
	BurnTopic     common.Hash
	TransferTopic common.Hash
	UniV2Swap     common.Hash
)

type Mint struct {
	LogIdx       uint
	Transaction  common.Hash
	Time         int64
	Height       uint64
	Sender       common.Address
	PairContract common.Address
	Token0       common.Address
	Token1       common.Address
	Amount0      *big.Float
	Amount1      *big.Float
	Reserve0     *big.Float
	Reserve1     *big.Float
}

type Burn struct {
	LogIdx       uint
	Transaction  common.Hash
	Time         int64
	Height       uint64
	Sender       common.Address
	Receiver     common.Address
	PairContract common.Address
	Token0       common.Address
	Token1       common.Address
	Amount0      *big.Float
	Amount1      *big.Float
	Reserve0     *big.Float
	Reserve1     *big.Float
}

type Swap struct {
	LogIdx       uint
	Transaction  common.Hash
	Time         int64
	Height       uint64
	Sender       common.Address
	Receiver     common.Address
	PairContract common.Address
	Token0       common.Address
	Token1       common.Address
	Amount0      *big.Float
	Amount1      *big.Float
	Reserve0     *big.Float
	Reserve1     *big.Float
}

type BlockSynopsis struct {
	TotalLogs uint64
	MintLogs  uint64
	BurnLogs  uint64
	SwapLogs  uint64
}

func init() {
	MintTopic = *(*common.Hash)(crypto.Keccak256([]byte("Mint(address,uint256,uint256)")))
	BurnTopic = *(*common.Hash)(crypto.Keccak256([]byte("Burn(address,uint256,uint256,address)")))
	TransferTopic = *(*common.Hash)(crypto.Keccak256([]byte("Transfer(address,address,uint256)")))
	UniV2Swap = *(*common.Hash)(crypto.Keccak256([]byte("Swap(address,uint256,uint256,uint256,uint256,address)")))
}
