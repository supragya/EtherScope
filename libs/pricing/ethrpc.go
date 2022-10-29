package priceresolver

import (
	"github.com/ethereum/go-ethereum/common"
)

type EthRPC interface {
	Service
	GetTokenName(common.Address) (string, error)
	GetCLLatestAnswer(common.Address) error
	GetUniswapLatestPrice(common.Address, common.Address) error
}

// DefaultEthRPC is default form of eth-rpc access for enhanced pricing engine
type DefaultEthRPC struct{}

func NewDefaultEthRPC() *DefaultEthRPC {
	return &DefaultEthRPC{}
}

// Implements EthRPC

// To be implemented
func (e *DefaultEthRPC) GetTokenName(common.Address) (string, error) {
	panic("unimplemented")
}

// To be implemented
func (e *DefaultEthRPC) GetCLLatestAnswer(common.Address) error {
	panic("unimplemented")
}

// To be implemented
func (e *DefaultEthRPC) GetUniswapLatestPrice(a, b common.Address) error {
	panic("unimplemented")
}
