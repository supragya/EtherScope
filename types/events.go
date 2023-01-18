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

	// ---- Topic Map and Reverse Map ---
	topicMap  map[string]common.Hash
	rTopicMap map[common.Hash]string
)

type Transfer struct {
	Type                string
	ProcessingType      ProcessingType `json:"-"`
	LogIdx              uint
	Transaction         common.Hash
	Time                uint64
	Height              uint64
	Token               common.Address
	Sender              common.Address
	TxSender            common.Address
	Receiver            common.Address
	Amount              *big.Float
	AmountUSD           *big.Float
	PriceDerivationMeta *PriceResult
}

type Mint struct {
	Type           string
	ProcessingType ProcessingType `json:"-"`
	LogIdx         uint
	Transaction    common.Hash
	Time           uint64
	Height         uint64
	Sender         common.Address
	TxSender       common.Address
	PairContract   common.Address
	Token0         common.Address
	Token1         common.Address
	Amount0        *big.Float
	Amount1        *big.Float
	Reserve0       *big.Float
	Reserve1       *big.Float
	AmountUSD      *big.Float
	Price0         *PriceResult
	Price1         *PriceResult
}

type Burn struct {
	Type           string
	ProcessingType ProcessingType `json:"-"`
	LogIdx         uint
	Transaction    common.Hash
	Time           uint64
	Height         uint64
	Sender         common.Address
	TxSender       common.Address
	PairContract   common.Address
	Token0         common.Address
	Token1         common.Address
	Amount0        *big.Float
	Amount1        *big.Float
	Reserve0       *big.Float
	Reserve1       *big.Float
	AmountUSD      *big.Float
	Price0         *PriceResult
	Price1         *PriceResult
}

type Swap struct {
	Type           string
	ProcessingType ProcessingType `json:"-"`
	LogIdx         uint
	Transaction    common.Hash
	Time           uint64
	Height         uint64
	Sender         common.Address
	TxSender       common.Address
	Receiver       common.Address
	PairContract   common.Address
	Token0         common.Address
	Token1         common.Address
	Amount0        *big.Float
	Amount1        *big.Float
	Reserve0       *big.Float
	Reserve1       *big.Float
	AmountUSD      *big.Float
	Price0         *PriceResult
	Price1         *PriceResult
}

type BlockSynopsis struct {
	Height                  uint64
	BlockTime               uint64
	IndexingTimeNanos       uint64
	ProcessingDurationNanos uint64
	PricingDurationNanos    uint64
	EventsScanned           uint64
	EventsPriced            uint64
	EventsUserDistribution  map[string]uint64
}

func toHash(str string) common.Hash {
	return *(*common.Hash)(crypto.Keccak256([]byte(str)))
}

func init() {
	topicMap = make(map[string]common.Hash, 10)
	rTopicMap = make(map[common.Hash]string, 10)

	// ---- Uniswap V2 ----
	UniV2MintTopic = setTopic("Mint(address,uint256,uint256)", "UniswapV2Mint")
	UniV2BurnTopic = setTopic("Burn(address,uint256,uint256,address)", "UniswapV2Burn")
	UniV2SwapTopic = setTopic("Swap(address,uint256,uint256,uint256,uint256,address)", "UniswapV2Swap")

	// ---- Uniswap V3 ----
	UniV3MintTopic = setTopic("Mint(address,address,int24,int24,uint128,uint256,uint256)", "UniswapV3Mint")
	UniV3BurnTopic = setTopic("Burn(address,int24,int24,uint128,uint256,uint256)", "UniswapV3Burn")
	UniV3SwapTopic = setTopic("Swap(address,address,int256,int256,uint160,uint128,int24)", "UniswapV3Swap")

	// ---- ERC 20 ---
	ERC20TransferTopic = setTopic("Transfer(address,address,uint256)", "ERC20Transfer")
}

func setTopic(topicString, infoString string) common.Hash {
	topicHash := toHash(topicString)
	topicMap[infoString] = topicHash
	rTopicMap[topicHash] = infoString
	return topicHash
}

func GetTopicForString(topicString string) (common.Hash, bool) {
	val, ok := topicMap[topicString]
	return val, ok
}

func GetStringForTopic(topicHash common.Hash) (string, bool) {
	val, ok := rTopicMap[topicHash]
	return val, ok
}
