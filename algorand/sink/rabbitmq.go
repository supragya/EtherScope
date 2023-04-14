package sink

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/config"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/service"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/version"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var (
	RabbitMQCFGSection   = "mq"
	RabbitMQCFGNecessity = "needed if `node.outputSinkType` == rabbitmq"
	RabbitMQCFGHeader    = "RabbitMQ output sink config"
	RabbitMQCFGFields    = [...]config.Field{
		{
			Name:      "queue",
			Type:      "string",
			Necessity: "always needed",
			Default:   "bgidx_processed",
		},
		{
			Name:      "secureConnection",
			Type:      "bool",
			Necessity: "always needed",
			Info:      "if set to true use secure connection (amqps)",
			Default:   false,
		},
		{
			Name:      "host",
			Type:      "string",
			Necessity: "always needed",
			Default:   "127.0.0.1",
		},
		{
			Name:      "port",
			Type:      "uint64",
			Necessity: "always needed",
			Default:   5672,
		},
		{
			Name:      "user",
			Type:      "string",
			Necessity: "always needed",
			Default:   "devuser",
		},
		{
			Name:      "pass",
			Type:      "string",
			Necessity: "always needed",
			Default:   "devpass",
		},
		{
			Name:      "queueIsDurable",
			Type:      "bool",
			Necessity: "always needed",
			Default:   true,
		},
		{
			Name:      "queueAutoDelete",
			Type:      "bool",
			Necessity: "always needed",
			Default:   false,
		},
		{
			Name:      "queueExclusive",
			Type:      "bool",
			Necessity: "always needed",
			Default:   false,
		},
		{
			Name:      "queueNoWait",
			Type:      "bool",
			Necessity: "always needed",
			Default:   false,
		},
	}
)

type RabbitMQOutputSinkImpl struct {
	service.BaseService

	// Parameters
	log              logger.Logger
	queueName        string
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

	fmt.Print(mqConnStr)

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

type WrappedPayload struct {
	PersistenceVersion uint8
	Data               interface{}
}

func (n *RabbitMQOutputSinkImpl) Send(payload interface{}) error {
	item, err := json.MarshalIndent(WrappedPayload{version.PersistenceVersion, payload}, "", " ")
	if err != nil {
		return err
	}

	err = n.channel.Publish(
		"",          // exchange
		n.queueName, // queue name
		true,        // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "application/json",
			Timestamp:       time.Now(),
			Body:            item,
		}, // message to publish
	)
	if err != nil {
		return err
	}
	n.log.Info("sent message onto outputsink rmq",
		"msglen", len(item),
		"queue", n.queueName)
	return nil
}

func NewRabbitMQOutputSinkWithViperFields(log logger.Logger) (OutputSink, error) {
	outs := &RabbitMQOutputSinkImpl{
		log:              log,
		queueName:        viper.GetString(RabbitMQCFGSection + ".queue"),
		secureConnection: viper.GetBool(RabbitMQCFGSection + ".secureConnection"),
		host:             viper.GetString(RabbitMQCFGSection + ".host"),
		port:             viper.GetUint64(RabbitMQCFGSection + ".port"),
		user:             viper.GetString(RabbitMQCFGSection + ".user"),
		pass:             viper.GetString(RabbitMQCFGSection + ".pass"),
		durable:          viper.GetBool(RabbitMQCFGSection + ".queueIsDurable"),  // durable
		autoDelete:       viper.GetBool(RabbitMQCFGSection + ".queueAutoDelete"), // auto delete
		exclusive:        viper.GetBool(RabbitMQCFGSection + ".queueExclusive"),  // exclusive
		noWait:           viper.GetBool(RabbitMQCFGSection + ".queueNoWait"),     // no wait
	}
	outs.BaseService = *service.NewBaseService(log, "outputsink", outs)
	return outs, nil
}
