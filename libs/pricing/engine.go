package priceresolver

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
)

// Enhanced cached, multistep, graph based batch pricing resolver system
type Engine struct {
	Service
	db     DB
	ethrpc EthRPC
	cex    CEX
}

// DefaultEngine is default form of enhanced pricing engine
func NewDefaultEngine() *Engine {
	var (
		db     = NewDefaultDB()
		ethrpc = NewDefaultEthRPC()
		cex    = NewDefaultCEX()
	)
	if err := db.Start(); err != nil {
		return nil
	}
	if err := ethrpc.Start(); err != nil {
		return nil
	}
	if err := cex.Start(); err != nil {
		return nil
	}
	return &Engine{db, ethrpc, cex}
}

func ConnectDB(path string) (*badger.DB, error) {
	opts := badger.DefaultOptions(path).
		WithCompression(options.None).
		WithSyncWrites(true)
	db, err := badger.Open(opts)
	return db, err
}
