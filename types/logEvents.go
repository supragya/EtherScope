package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// ---- Uniswap V2 ----
	UniV2MintTopic common.Hash
	UniV2BurnTopic common.Hash
	UniV2SwapTopic common.Hash

	// ---- Uniswap V3 ----
	UniV3IncreaseLiquidityTopic common.Hash // Not in use
	UniV3DecreaseLiquidityTopic common.Hash // Not in use
	UniV3MintTopic              common.Hash
	UniV3BurnTopic              common.Hash
	UniV3SwapTopic              common.Hash

	// ---- ERC 20 ----
	ERC20TransferTopic common.Hash
)

type Transfer struct {
	Type        string
	Network     uint
	LogIdx      uint
	Transaction common.Hash
	Time        uint64
	Height      uint64
	Token       common.Address
	Sender      common.Address
	Receiver    common.Address
	Amount      *big.Float
	AmountUSD   *PriceResult
}

type Mint struct {
	Type         string
	Network      uint
	LogIdx       uint
	Transaction  common.Hash
	Time         uint64
	Height       uint64
	Sender       common.Address
	PairContract common.Address
	Token0       common.Address
	Token1       common.Address
	Amount0      *big.Float
	Amount1      *big.Float
	Reserve0     *big.Float
	Reserve1     *big.Float
	AmountUSD    *big.Float
	Price0       *PriceResult
	Price1       *PriceResult
}

type Burn struct {
	Type         string
	Network      uint
	LogIdx       uint
	Transaction  common.Hash
	Time         uint64
	Height       uint64
	Sender       common.Address
	PairContract common.Address
	Token0       common.Address
	Token1       common.Address
	Amount0      *big.Float
	Amount1      *big.Float
	Reserve0     *big.Float
	Reserve1     *big.Float
	AmountUSD    *big.Float
	Price0       *PriceResult
	Price1       *PriceResult
}

type Swap struct {
	Type         string
	Network      uint
	LogIdx       uint
	Transaction  common.Hash
	Time         uint64
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
	AmountUSD    *big.Float
	Price0       *PriceResult
	Price1       *PriceResult
}

type BlockSynopsis struct {
	Type         string
	Network      uint
	Height       uint64
	Time         uint64
	TotalLogs    uint64
	MintLogs     uint64
	BurnLogs     uint64
	SwapLogs     uint64
	TransferLogs uint64
}

func toHash(str string) common.Hash {
	return *(*common.Hash)(crypto.Keccak256([]byte(str)))
}

func init() {
	// ---- Uniswap V2 ----
	UniV2MintTopic = toHash("Mint(address,uint256,uint256)")
	UniV2BurnTopic = toHash("Burn(address,uint256,uint256,address)")
	UniV2SwapTopic = toHash("Swap(address,uint256,uint256,uint256,uint256,address)")

	// ---- Uniswap V3 ----
	UniV3MintTopic = toHash("Mint(address,address,int24,int24,uint128,uint256,uint256)")
	UniV3BurnTopic = toHash("Burn(address,int24,int24,uint128,uint256,uint256)")
	UniV3SwapTopic = toHash("Swap(address,address,int256,int256,uint160,uint128,int24)")

	// ---- ERC 20 ---
	ERC20TransferTopic = toHash("Transfer(address,address,uint256)")
}
