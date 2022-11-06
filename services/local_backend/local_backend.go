package localbackend

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
)

type LocalBackend interface {
	service.Service

	Get(key string) ([]byte, bool, error)
	Set(key string, val interface{}) error
	Sync() error
}

const (
	// Provides a tuple containing latest pricing graph
	// height and graph definition
	KeyLatestPricingGraph = "lapg"

	// Provides lowest height for which pricing graph is
	// available
	KeyLowestPricingHeight = "loph"

	// Graph prefix for height
	KeyGraphPrefix = "pg"
)
