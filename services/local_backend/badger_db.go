package localbackend

import (
	"context"
	"fmt"
	"sync"
	"time"

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
			Name:      "periodicSync",
			Type:      "time.Duration",
			Necessity: "always needed",
			Info: cfg.SArr("time duration after which periodic sync to disk",
				"and call to badgerDB garbage collector is done"),
			Default: "30s",
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

	periodicSync time.Duration
	log          logger.Logger
	lock         *sync.RWMutex
	dbLocation   string
	namespace    string
	inMem        map[string][]byte
	db           *badger.DB
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

	// periodic runtime GC goroutine
	// https://dgraph.io/docs/badger/get-started/#garbage-collection
	go n.loop()

	return nil
}

// OnStop stops the badgerdb LocalBackend. It implements service.Service
func (n *BadgerDBLocalBackendImpl) OnStop() {
	n.db.Close()
}

func (n *BadgerDBLocalBackendImpl) Get(key string) ([]byte, bool, error) {
	n.lock.RLock()
	defer n.lock.RUnlock()

	queryKey := fmt.Sprintf("%s::%s", n.namespace, key)
	if val, ok := n.inMem[queryKey]; ok {
		return val, true, nil
	}
	value := []byte{}
	err := n.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(queryKey))
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(value)
		if err != nil {
			return err
		}
		return nil
	})
	if err == badger.ErrKeyNotFound {
		return []byte{}, false, nil
	}
	return value, true, err
}

func (n *BadgerDBLocalBackendImpl) Set(key string, val []byte) error {
	n.lock.Lock()
	defer n.lock.Unlock()

	queryKey := fmt.Sprintf("%s::%s", n.namespace, key)
	// if len(val) < 20 {
	// 	n.log.Info("set", queryKey, val)
	// }
	n.inMem[queryKey] = val
	return nil
}

func (n *BadgerDBLocalBackendImpl) loop() {
	for {
		<-time.After(n.periodicSync)
		n.log.Info("running periodic on-disk sync")
		if err := n.db.Sync(); err != nil {
			n.log.Fatal(err.Error())
		}

		n.log.Info("running periodic badgerdb garbage collector")
		n.db.RunValueLogGC(0.5)
	}
}

func (n *BadgerDBLocalBackendImpl) Sync() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	// multiple sync call protection
	if len(n.inMem) == 0 {
		// n.log.Info("multiple sync call protection")
		return nil
	}

	// Start a writable transaction.
	txn := n.db.NewTransaction(true)

	// Use the transaction...
	var err error
	count := 0
	for key, val := range n.inMem {
		err = txn.Set([]byte(key), val)
		if err == badger.ErrTxnTooBig {
			n.log.Info("too big of a transaction, splitting")
			err = txn.Commit()
			if err != nil {
				n.log.Warn("too big a transaction, splitting into multiple")
				return err
			}
			txn = n.db.NewTransaction(true)
			err = txn.Set([]byte(key), val)
			if err != nil {
				return err
			}
		}
		count++
	}
	err = txn.Commit()
	if err != nil {
		n.log.Error("transaction commit failed", "error", err)
		return err
	}

	// Very important step apparently
	txn.Discard()

	err = n.db.Sync()
	if err != nil {
		n.log.Error("localbackend on-disk sync failed", "error", err)
		return err
	}

	n.log.Info("localbackend on-disk sync completed", "records", count)

	// Flush inMem db
	n.inMem = make(map[string][]byte, 10000)

	return nil
}

func NewBadgerDBWithViperFields(log logger.Logger) (LocalBackend, error) {
	// ensure field integrity for viper
	for _, mf := range BadgerCFGFields {
		err := cfg.EnsureFieldIntegrity(BadgerCFGSection, mf)
		if err != nil {
			return nil, err
		}
	}

	lb := &BadgerDBLocalBackendImpl{
		log:          log,
		lock:         &sync.RWMutex{},
		periodicSync: viper.GetDuration(BadgerCFGSection + ".periodicSync"),
		dbLocation:   viper.GetString(BadgerCFGSection + ".dbLocation"),
		namespace:    viper.GetString(BadgerCFGSection + ".namespace"),
		inMem:        make(map[string][]byte, 10000),
		db:           nil,
	}
	lb.BaseService = *service.NewBaseService(log, "localbackend", lb)
	return lb, nil
}
