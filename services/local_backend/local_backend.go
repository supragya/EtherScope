package localbackend

import (
	"context"

	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/spf13/viper"
)

type LocalBackend interface {
	service.Service
}

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

	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	n.db = db

	return nil
}

// OnStop stops the badgerdb LocalBackend. It implements service.Service
func (n *BadgerDBLocalBackendImpl) OnStop() {
}

func NewBadgerDBWithViperFields(log logger.Logger) (LocalBackend, error) {
	lb := &BadgerDBLocalBackendImpl{
		log:        log,
		dbLocation: viper.GetString("localbackend.db"),
		db:         nil,
	}
	lb.BaseService = *service.NewBaseService(log, "localbackend", lb)
	return lb, nil
}
