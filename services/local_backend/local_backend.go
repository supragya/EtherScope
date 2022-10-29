package localbackend

import "github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"

type LocalBackend interface {
	service.Service

	GetKey()
}

type BadgerDBLocalBackend struct {
}

func NewBadgerDBWithViperFields() (*BadgerDBLocalBackend, error) {
	return &BadgerDBLocalBackend{}, nil
}
