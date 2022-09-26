package ethrpc

import (
	"math/big"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/mspool"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/ethclient"
	lru "github.com/hashicorp/golang-lru"
)

const WD = 20

type EthRPC struct {
	upstreams           *mspool.MasterSlavePool[ethclient.Client]
	isErigon            bool
	contractTokensCache *lru.ARCCache
	ERC20Cache          *lru.ARCCache
	PriceCache          *lru.ARCCache
	pricing             *Pricing
}

type UniV2Reserves struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}

func NewEthRPC(isErigon bool, masterUpstream string, slaveUpstreams []string, timeout time.Duration) *EthRPC {
	ctcache, err := lru.NewARC(1024) // Hardcoded 1024
	util.ENOK(err)

	erc20cache, err := lru.NewARC(1024) // Hardcoded 1024
	util.ENOK(err)

	ratecache, err := lru.NewARC(1024) // Hardcoded 1024
	util.ENOK(err)

	pool, err := mspool.NewEthClientMasterSlavePool(masterUpstream, slaveUpstreams, mspool.DefaultMSPoolConfig, timeout)
	util.ENOK(err)

	return &EthRPC{
		upstreams:           pool,
		isErigon:            isErigon,
		contractTokensCache: ctcache,
		ERC20Cache:          erc20cache,
		PriceCache:          ratecache,
		pricing:             GetPricingEngine(),
	}
}

func (d *EthRPC) Len() int {
	return len(d.upstreams.Slaves) + 1
}
