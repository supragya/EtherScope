package outputsink

import "github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"

type OutputSink interface {
	service.Service

	GetKey()
}

type RabbitMQOutputSink struct {
}

func NewRabbitMQOutputSinkWithViperFields() (*RabbitMQOutputSink, error) {
	return &RabbitMQOutputSink{}, nil
}
