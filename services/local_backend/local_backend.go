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
