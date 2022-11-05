package localbackend

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
)

type LocalBackend interface {
	service.Service

	Sync() error
}
