package dataaccess

import (
	"math/big"

	"github.com/Blockpour/Blockpour-Geth-Indexer/mspool"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/ethclient"
	lru "github.com/hashicorp/golang-lru"
)

const WD = 20

type DataAccess struct {
	upstreams           *mspool.MasterSlavePool[ethclient.Client]
	isErigon            bool
	contractTokensCache *lru.ARCCache
	ERC20Cache          *lru.ARCCache
	RateCache           *lru.ARCCache
	pricing             *Pricing
}

type UniV2Reserves struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}

func NewDataAccess(isErigon bool, masterUpstream string, slaveUpstreams []string) *DataAccess {
	ctcache, err := lru.NewARC(1024) // Hardcoded 1024
	util.ENOK(err)

	erc20cache, err := lru.NewARC(1024) // Hardcoded 1024
	util.ENOK(err)

	ratecache, err := lru.NewARC(1024) // Hardcoded 1024
	util.ENOK(err)

	pool, err := mspool.NewEthClientMasterSlavePool(masterUpstream, slaveUpstreams, mspool.DefaultMSPoolConfig)
	util.ENOK(err)

	return &DataAccess{
		upstreams:           pool,
		isErigon:            isErigon,
		contractTokensCache: ctcache,
		ERC20Cache:          erc20cache,
		RateCache:           ratecache,
		pricing:             GetPricingEngine(),
	}
}

func (d *DataAccess) Len() int {
	return len(d.upstreams.Slaves) + 1
}
