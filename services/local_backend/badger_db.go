package localbackend

import (
	"context"

	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
	"github.com/spf13/viper"
)

var (
	BadgerCFGSection   = "localBackendBadger"
	BadgerCFGNecessity = "needed if `node.localBackendType` == badgerdb"
	BadgerCFGHeader    = cfg.SArr("badgerdb is an impl for LocalBackend used",
		"by indexer to persist local caches and pricing",
		"graphs")
	BadgerCFGFields = [...]cfg.Field{
		{
			Name:      "dbLocation",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("disk location where database is to be persisted"),
			Default:   "lb.badger.db",
		},
		{
			Name:      "namespace",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("rootlevel namespace on database"),
			Default:   "bp",
		},
	}
)

type BadgerDBLocalBackendImpl struct {
	service.BaseService

	log        logger.Logger
	dbLocation string
	db         *badger.DB
}

// OnStart starts the badgerdb LocalBackend. It implements service.Service.
func (n *BadgerDBLocalBackendImpl) OnStart(ctx context.Context) error {
	// Open the Badger database located in the n.db directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions(n.dbLocation)
	opts.Logger = n.log.With("module", "badgerdb")
	opts.Compression = options.None

	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	n.db = db

	return nil
}

// OnStop stops the badgerdb LocalBackend. It implements service.Service
func (n *BadgerDBLocalBackendImpl) OnStop() {
	n.db.Close()
}

func NewBadgerDBWithViperFields(log logger.Logger) (LocalBackend, error) {
	lb := &BadgerDBLocalBackendImpl{
		log:        log,
		dbLocation: viper.GetString(BadgerCFGSection + ".dbLocation"),
		db:         nil,
	}
	lb.BaseService = *service.NewBaseService(log, "localbackend", lb)
	return lb, nil
}
