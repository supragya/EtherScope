package outputsink

import (
	"context"
	"fmt"

	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type OutputSink interface {
	service.Service
}

type RabbitMQOutputSinkImpl struct {
	service.BaseService

	// Parameters
	log              logger.Logger
	queueName        string
	network          string
	secureConnection bool
	host             string
	port             uint64
	user             string
	pass             string
	durable          bool
	autoDelete       bool
	exclusive        bool
	noWait           bool

	// Connections
	connection *amqp.Connection
	channel    *amqp.Channel
}

// OnStart starts the rabbitmq OutputSink. It implements service.Service.
func (n *RabbitMQOutputSinkImpl) OnStart(ctx context.Context) error {
	connPrefix := "amqp"
	if viper.GetBool("mq.secureConnection") {
		connPrefix = "amqps"
	}

	mqConnStr := fmt.Sprintf("%s://%s:%s@%s:%d/", connPrefix, n.user, n.pass, n.host, n.port)

	connectRabbitMQ, err := amqp.Dial(mqConnStr)
	if err != nil {
		return err
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return err
	}

	n.connection = connectRabbitMQ
	n.channel = channelRabbitMQ
	return nil
}

// OnStop stops the rabbitmq OutputSink. It implements service.Service
func (n *RabbitMQOutputSinkImpl) OnStop() {
	n.channel.Close()
	n.connection.Close()
}

func NewRabbitMQOutputSinkWithViperFields(log logger.Logger) (OutputSink, error) {
	outs := &RabbitMQOutputSinkImpl{
		log:              log,
		queueName:        viper.GetString("mq.queue"),
		network:          viper.GetString("mq.network"),
		secureConnection: viper.GetBool("mq.secureConnection"),
		host:             viper.GetString("mq.host"),
		port:             viper.GetUint64("mq.port"),
		user:             viper.GetString("mq.user"),
		pass:             viper.GetString("mq.pass"),
		durable:          viper.GetBool("mq.queueIsDurable"),  // durable
		autoDelete:       viper.GetBool("mq.queueAutoDelete"), // auto delete
		exclusive:        viper.GetBool("mq.queueExclusive"),  // exclusive
		noWait:           viper.GetBool("mq.queueNoWait"),     // no wait
	}
	outs.BaseService = *service.NewBaseService(log, "outputsink", outs)
	return outs, nil
}
