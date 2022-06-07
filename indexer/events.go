package indexer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	MintTopic     common.Hash
	BurnTopic     common.Hash
	TransferTopic common.Hash
)

type Mint struct {
	logIdx       uint
	transaction  common.Hash
	height       uint64
	sender       common.Address
	pairContract common.Address
	token0       common.Address
	token1       common.Address
	amount0      float64
	amount1      float64
	reserve0     *big.Float
	reserve1     *big.Float
}

func init() {
	MintTopic = *(*common.Hash)(crypto.Keccak256([]byte("Mint(address,uint256,uint256)")))
	BurnTopic = *(*common.Hash)(crypto.Keccak256([]byte("Burn(address,uint256,uint256,address)")))
	TransferTopic = *(*common.Hash)(crypto.Keccak256([]byte("Transfer(address,address,uint256)")))
}
