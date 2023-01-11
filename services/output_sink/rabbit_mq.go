package outputsink

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	"github.com/Blockpour/Blockpour-Geth-Indexer/version"
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
	disconnectTime   time.Time
	cachedMessages   []interface{}
	connecting       bool

	// Connections
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (n *RabbitMQOutputSinkImpl) IsReady() bool {
	ready := len(n.cachedMessages) == 0 && n.connection != nil && !n.connection.IsClosed() && n.channel != nil
	if !ready && n.disconnectTime.IsZero() {
		n.log.Warn("Setting disconnect time")
		n.disconnectTime = time.Now()
	}

	return ready
}

// OnStart starts the rabbitmq OutputSink. It implements service.Service.
func (n *RabbitMQOutputSinkImpl) OnStart(ctx context.Context) error {
	n.connecting = true
	defer func() {
		n.connecting = false
	}()
	if err := n.connect(); err != nil {
		n.disconnectTime = time.Now()
		n.log.Info(fmt.Sprintf("Unable to connect to RabbitMQ: %s", fmt.Sprint(err)))
		return err
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

	return fmt.Sprintf("%s://%s:%s@%s:%d/", connPrefix, n.user, n.pass, n.host, n.port)
}

func (n *RabbitMQOutputSinkImpl) connect() error {
	mqConnStr := n.getConnectionString()

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

/*
Callable function which will trigger the service to reconnect to the MQ.
*/
func (n *RabbitMQOutputSinkImpl) Reconnect() error {
	if n.connecting {
		return nil
	}
	n.connecting = true

	err := n.connect()
	if err != nil {
		n.connecting = false
		return err
	}

	n.log.Info(fmt.Sprintf("RabbitMQ reconnected. Downtime: %dms",
		time.Since(n.disconnectTime).Milliseconds()))

	n.disconnectTime = time.Time{}
	err = n.onReconnect()
	n.connecting = false
	return err
}

/*
Behavior triggered following MQ reconnection
*/
func (n *RabbitMQOutputSinkImpl) onReconnect() error {
	// Handle cached messages
	if len(n.cachedMessages) > 0 {
		if err := n.pushCachedMessages(); err != nil {
			return err
		}
	}
	return nil
}

/*
Pushes cached messages to RabbitMQ
*/
func (n *RabbitMQOutputSinkImpl) pushCachedMessages() error {
	n.log.Info(fmt.Sprintf("Attempting to push %d cached messages to RabbitMQ", len(n.cachedMessages)))
	messageCount := 0
	for !n.connection.IsClosed() && len(n.cachedMessages) > 0 {
		messageCount += 1
		n.log.Debug(fmt.Sprintf("Pushing cached message #%d", messageCount))
		message := n.cachedMessages[0]
		err := n.Send(message)
		if err != nil {
			return err
		}
		n.cachedMessages = n.cachedMessages[1:]
	}
	return nil
}

/*
Places a message in the cache. Messages in the cache are pushed to RabbitMQ
upon reconnection
*/
func (n *RabbitMQOutputSinkImpl) cacheMessage(payload interface{}) {
	n.log.Info("Caching message")
	n.cachedMessages = append(n.cachedMessages, payload)
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
		n.log.Warn("Error publishing message to RabbitMQ: " + fmt.Sprint(err) + ", caching message")
		if n.disconnectTime.IsZero() {
			n.disconnectTime = time.Now()
		}
		n.cacheMessage(payload)
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
