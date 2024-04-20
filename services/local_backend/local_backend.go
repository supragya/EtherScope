package localbackend

import (
	"github.com/supragya/EtherScope/libs/service"
)

type LocalBackend interface {
	service.Service

	Get(key string) ([]byte, bool, error)
	Set(key string, val []byte) error
	Sync() error
}

const (
	// Provides a uint64 of latest height
	KeyLatestHeight = "lh"

	// Provides a tuple containing latest pricing graph
	// height and graph definition
	KeyLatestPricingGraph = "lapg"

	// Provides lowest height for which pricing graph is
	// available
	KeyLowestPricingHeight = "loph"

	// Graph prefix for height
	KeyGraphPrefix = "pg"
)
