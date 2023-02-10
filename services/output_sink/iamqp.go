package outputsink

import (
	"github.com/streadway/amqp"
)

type (
	AMQP interface {
		Dial(address string) (RabbitMQConnection, error)
	}

	RabbitMQConnection interface {
		Channel() (*amqp.Channel, error)
		Close() error
		IsClosed() bool
	}

	RabbitMQChannel interface {
		Publish(string, string, bool, bool, amqp.Publishing) error
		Close() error
	}
)

type AMQPImpl struct{}

func (amqpImpl *AMQPImpl) Dial(address string) (RabbitMQConnection, error) {
	return amqp.Dial(address)
}
