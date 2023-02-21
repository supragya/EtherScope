package iamqp

import (
	"github.com/streadway/amqp"
)

type (
	AMQP interface {
		Dial(address string) (AMQPConnection, error)
	}

	AMQPConnection interface {
		Channel() (AMQPChannel, error)
		Close() error
		IsClosed() bool
	}

	AMQPChannel interface {
		Publish(string, string, bool, bool, amqp.Publishing) error
		Close() error
	}
)

type AMQPImpl struct{}

func (amqpImpl *AMQPImpl) Dial(address string) (AMQPConnection, error) {
	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, err
	}
	return NewAMQPWrappedConnection(conn), nil
}

type AMQPWrappedConnection struct {
	connection *amqp.Connection
}

func (conn AMQPWrappedConnection) Channel() (AMQPChannel, error) {
	return conn.connection.Channel()
}

func (conn AMQPWrappedConnection) Close() error {
	return conn.connection.Close()
}

func (conn AMQPWrappedConnection) IsClosed() bool {
	return conn.connection.IsClosed()
}

func NewAMQPWrappedConnection(conn *amqp.Connection) AMQPWrappedConnection {
	return AMQPWrappedConnection{
		connection: conn,
	}
}
