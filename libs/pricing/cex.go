package priceresolver

import (
	"math/big"
	"time"
)

type CEX interface {
	Service
	GetPrice(string, time.Time) (*big.Float, error)
}

// DefaultCEX is default form of eth-rpc access for enhanced pricing engine
type DefaultCEX struct{}

func NewDefaultCEX() *DefaultCEX {
	return &DefaultCEX{}
}

// Implements CEX

// To be implemented
func (e *DefaultCEX) GetPrice(string, time.Time) (*big.Float, error) {
	panic("unimplemented")
}
