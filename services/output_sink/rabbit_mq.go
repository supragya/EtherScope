package outputsink

import (
	"context"
	"fmt"

	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var (
	RabbitMQCFGSection   = "outputSinkRabbitMQ"
	RabbitMQCFGNecessity = "needed if `node.outputSinkType` == rabbitmq"
	RabbitMQCFGHeader    = cfg.SArr("rabbitmq is an impl for OutputSink used",
		"by indexer to send indexed information to the",
		"backend")
	RabbitMQCFGFields = [...]cfg.Field{
		{
			Name:      "queue",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("queue to output processed info onto"),
			Default:   "bgidx_processed",
		},
		{
			Name:      "secureConnection",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("if set to true use secure connection (amqps)"),
			Default:   false,
		},
		{
			Name:      "host",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("rabbitmq host"),
			Default:   "127.0.0.1",
		},
		{
			Name:      "port",
			Type:      "uint64",
			Necessity: "always needed",
			Info:      cfg.SArr("rabbitmq port"),
			Default:   5672,
		},
		{
			Name:      "user",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("rabbitmq user"),
			Default:   "devuser",
		},
		{
			Name:      "pass",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("rabbitmq pass"),
			Default:   "devpass",
		},
		{
			Name:      "queueIsDurable",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("queue durability"),
			Default:   true,
		},
		{
			Name:      "queueAutoDelete",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("queue autodelete"),
			Default:   false,
		},
		{
			Name:      "queueExclusive",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("queue exclusivity"),
			Default:   false,
		},
		{
			Name:      "queueNoWait",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("queue no wait"),
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
