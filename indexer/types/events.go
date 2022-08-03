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
	UniV3Swap     common.Hash
)

type tokenMeta struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

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
	AmountUSD    float64
	Price0       float64
	Price1       float64
	Meta         tokenMeta
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
	AmountUSD    float64
	Price0       float64
	Price1       float64
	Meta         tokenMeta
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
	AmountUSD    float64
	Price0       float64
	Price1       float64
	Meta         tokenMeta
}

type BlockSynopsis struct {
	Height    uint64
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
	UniV3Swap = *(*common.Hash)(crypto.Keccak256([]byte("Swap(address,address,int256,int256,uint160,uint128,int24)")))
}
