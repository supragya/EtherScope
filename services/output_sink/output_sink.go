package outputsink

import (
	"context"

	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
)

type OutputSink interface {
	service.Service
}

type RabbitMQOutputSinkImpl struct {
	service.BaseService

	log logger.Logger
}

// OnStart starts the rabbitmq OutputSink. It implements service.Service.
func (n *RabbitMQOutputSinkImpl) OnStart(ctx context.Context) error {
	return nil
}

// OnStop stops the rabbitmq OutputSink. It implements service.Service
func (n *RabbitMQOutputSinkImpl) OnStop() {
}

func NewRabbitMQOutputSinkWithViperFields(log logger.Logger) (OutputSink, error) {
	outs := &RabbitMQOutputSinkImpl{
		log: log,
	}
	outs.BaseService = *service.NewBaseService(log, "outputsink", outs)
	return outs, nil
}
