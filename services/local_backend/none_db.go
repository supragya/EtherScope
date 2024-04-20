package localbackend

import (
	"context"
	"errors"

	logger "github.com/supragya/EtherScope/libs/log"
	"github.com/supragya/EtherScope/libs/service"
)

type NoneDBImpl struct {
	service.BaseService
}

// OnStart starts the nonedb LocalBackend. It implements service.Service.
func (n *NoneDBImpl) OnStart(ctx context.Context) error {
	return nil
}

// OnStop stops the nonedb LocalBackend. It implements service.Service
func (n *NoneDBImpl) OnStop() {
}

func (n *NoneDBImpl) Get(key string) ([]byte, bool, error) {
	return []byte{}, false, nil
}

func (n *NoneDBImpl) Set(key string, val []byte) error {
	return errors.New("cannot set on none_db")
}

func (n *NoneDBImpl) Sync() error {
	return nil
}

func NewNoneDB(log logger.Logger) (LocalBackend, error) {
	lb := &NoneDBImpl{}
	lb.BaseService = *service.NewBaseService(log, "localbackend", lb)
	return lb, nil
}
