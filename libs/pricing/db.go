package priceresolver

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type DBStatus uint

const (
	Unintialized DBStatus = iota
	Connected
	Closed
)

type ChainlinkDirect interface {
	GetChainlinkDirectBatchPrices([][]common.Address) (map[common.Address]*big.Float, error)
	RemChainlinkDirectBatchPrices([][]common.Address)
	SetChainlinkDirectBatchPrices([]common.Address, []*big.Float) error
}

type ChainlinkDerived interface {
	GetChainlinkDerivedBatchPrices([][]common.Address) (map[common.Address]*big.Float, error)
	RemChainlinkDerivedBatchPrices([][]common.Address)
	SetChainlinkDerivedBatchPrices([]common.Address, []*big.Float) error
}

type UniswapDirect interface {
	GetUniswapDirectBatchPrices([][]common.Address) (map[common.Address]*big.Float, error)
	RemUniswapDirectBatchPrices([][]common.Address)
	SetUniswapDirectBatchPrices([]common.Address, []*big.Float) error
}

type UniswapDerived interface {
	GetUniswapDerivedBatchPrices([][]common.Address) (map[common.Address]*big.Float, error)
	RemUniswapDerivedBatchPrices([][]common.Address)
	SetUniswapDerivedBatchPrices([]common.Address, []*big.Float) error
}

type CEXDirect interface {
	GetCEXDirectBatchPrices([][]common.Address) (map[common.Address]*big.Float, error)
	RemCEXDirectBatchPrices([][]common.Address)
	SetCEXDirectBatchPrices([]common.Address, []*big.Float) error
}

type CEXDerived interface {
	GetCEXDerivedBatchPrices([][]common.Address) (map[common.Address]*big.Float, error)
	RemCEXDerivedBatchPrices([][]common.Address)
	SetCEXDerivedBatchPrices([]common.Address, []*big.Float) error
}

type DBTokenNameResolver interface {
	GetTokenNames([]common.Address) ([]string, error)
	SetTokenNames([]common.Address, []string) error
}

type DB interface {
	Service
	DBTokenNameResolver
	ChainlinkDirect
	ChainlinkDerived
	UniswapDirect
	UniswapDerived
	CEXDirect
	CEXDerived
}
