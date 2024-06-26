package outputsink

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	iamqp "github.com/supragya/EtherScope/libs/amqp"
	cfg "github.com/supragya/EtherScope/libs/config"
	logger "github.com/supragya/EtherScope/libs/log"
	"github.com/supragya/EtherScope/libs/service"
	"github.com/supragya/EtherScope/version"
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
			Default:   "escope_processed",
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
	disconnectTime   time.Time
	connecting       bool
	amqpImpl         iamqp.AMQP

	// Connections
	connection iamqp.AMQPConnection
	channel    iamqp.AMQPChannel
}

// OnStart starts the rabbitmq OutputSink. It implements service.Service.
func (n *RabbitMQOutputSinkImpl) OnStart(ctx context.Context) error {
	if err := n.connect(); err != nil {
		n.disconnectTime = time.Now()
		n.log.Info(fmt.Sprintf("Unable to connect to RabbitMQ: %s", err))
		return fmt.Errorf("OutputSinkStartupError: %w", err)
	}
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

func (n *RabbitMQOutputSinkImpl) getConnectionString() string {
	connPrefix := "amqp"
	if viper.GetBool("mq.secureConnection") {
		connPrefix = "amqps"
	}
	return fmt.Sprintf("%s://%s:%s@%s:%d/", connPrefix, n.user, url.QueryEscape(n.pass), n.host, n.port)
}

func (n *RabbitMQOutputSinkImpl) connect() error {
	if n.connecting {
		return nil
	}
	n.connecting = true
	defer func() {
		n.connecting = false
	}()

	mqConnStr := n.getConnectionString()
	connectRabbitMQ, err := n.amqpImpl.Dial(mqConnStr)
	if err != nil {
		if n.disconnectTime.IsZero() {
			n.disconnectTime = time.Now()
		}

		return fmt.Errorf("OutputSinkDialError caused by: %w", err)
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		if n.disconnectTime.IsZero() {
			n.disconnectTime = time.Now()
		}
		return fmt.Errorf("OutputSinkChannelError caused by: %w", err)
	}

	if n.disconnectTime.IsZero() {
		n.log.Info("RabbitMQ connected")
	} else {
		n.log.Info(fmt.Sprintf("RabbitMQ reconnected. Downtime: %dms",
			time.Since(n.disconnectTime).Milliseconds()))
		n.disconnectTime = time.Time{}
	}

	n.connection = connectRabbitMQ
	n.channel = channelRabbitMQ
	return nil
}

func (n *RabbitMQOutputSinkImpl) Send(payload interface{}) error {
	if n.connection == nil || n.connection.IsClosed() {
		if err := n.connect(); err != nil {
			return fmt.Errorf("OutputSinkUnavailable Caused By: %w", err)
		}
	}

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
		n.log.Warn("Error publishing message to RabbitMQ: " + fmt.Sprint(err) + ", caching message")
		if n.disconnectTime.IsZero() {
			n.disconnectTime = time.Now()
		}
		return fmt.Errorf("OutputSinkPublishError Caused by: %w", err)
	}

	n.log.Debug("sent message onto outputsink rmq",
		"msglen", len(item),
		"queue", n.queueName)

	return nil
}

func NewRabbitMQOutputSinkWithViperFields(log logger.Logger, amqpImpl iamqp.AMQP) (OutputSink, error) {
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
	outs.amqpImpl = amqpImpl
	return outs, nil
}
