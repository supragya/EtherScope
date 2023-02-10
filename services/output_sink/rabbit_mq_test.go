package outputsink_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"

	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	outs "github.com/Blockpour/Blockpour-Geth-Indexer/services/output_sink"
)

type TestAMQPImpl struct {
	DialImpl func(address string) (outs.AMQPConnection, error)
}

func (amqpImpl TestAMQPImpl) Dial(address string) (outs.AMQPConnection, error) {
	if amqpImpl.DialImpl != nil {
		return amqpImpl.DialImpl(address)
	} else {
		return nil, nil
	}
}

type TestAMQPConnectionImpl struct {
	ChannelImpl  func() (outs.AMQPChannel, error)
	CloseImpl    func() error
	IsClosedImpl func() bool
}

func (testConnection TestAMQPConnectionImpl) Channel() (outs.AMQPChannel, error) {
	if testConnection.ChannelImpl != nil {
		return testConnection.ChannelImpl()
	}
	return nil, nil
}

func (testConnection TestAMQPConnectionImpl) Close() error {
	if testConnection.CloseImpl != nil {
		return testConnection.CloseImpl()
	}
	return nil
}

func (testConnection TestAMQPConnectionImpl) IsClosed() bool {
	if testConnection.IsClosedImpl != nil {
		return testConnection.IsClosedImpl()
	}
	return false
}

type TestAMQPChannelImpl struct {
	PublishImpl func(string, string, bool, bool, amqp.Publishing) error
	CloseImpl   func() error
}

func (channel TestAMQPChannelImpl) Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	if channel.PublishImpl != nil {
		return channel.PublishImpl(exchange, key, mandatory, immediate, msg)
	}
	return nil
}

func (channel TestAMQPChannelImpl) Close() error {
	if channel.PublishImpl != nil {
		return channel.Close()
	}
	return nil
}

var _ = Describe("RabbitMq", func() {
	var testLogger logger.Logger
	var testAMQP TestAMQPImpl
	var testConnection TestAMQPConnectionImpl
	var testOutputSinkRMQ outs.OutputSink
	var testRabbitMQChannel TestAMQPChannelImpl

	BeforeEach(func() {
		viper.GetViper().Set(outs.RabbitMQCFGSection+".queue", "testqueue")
		viper.GetViper().Set(outs.RabbitMQCFGSection+".secureConnection", false)
		viper.GetViper().Set(outs.RabbitMQCFGSection+".host", "localhost")
		viper.GetViper().Set(outs.RabbitMQCFGSection+".port", 1111)
		viper.GetViper().Set(outs.RabbitMQCFGSection+".user", "testuser")
		viper.GetViper().Set(outs.RabbitMQCFGSection+".pass", "testpass")
		viper.GetViper().Set(outs.RabbitMQCFGSection+".queueIsDurable", false)
		viper.GetViper().Set(outs.RabbitMQCFGSection+".queueAutoDelete", true)
		viper.GetViper().Set(outs.RabbitMQCFGSection+".queueExclusive", false)
		viper.GetViper().Set(outs.RabbitMQCFGSection+".queueNoWait", false)

		testAMQP = TestAMQPImpl{
			DialImpl: nil,
		}

		testLogger = logger.NewNopLogger()
		testOutputSinkRMQ, _ = outs.NewRabbitMQOutputSinkWithViperFields(testLogger, &testAMQP)
		testConnection = TestAMQPConnectionImpl{}
		testAMQP.DialImpl = func(address string) (outs.AMQPConnection, error) {
			return testConnection, nil
		}

		testRabbitMQChannel = TestAMQPChannelImpl{}
		testConnection.ChannelImpl = func() (outs.AMQPChannel, error) {
			return testRabbitMQChannel, nil
		}
	})

	Context("Start", func() {
		It("should pass on err when connect fails", func() {
			testAMQP.DialImpl = func(address string) (outs.AMQPConnection, error) { return nil, errors.New("Test dial error") }
			Expect(testOutputSinkRMQ.Start(context.Background())).To(MatchError("Test dial error"))
		})

		It("should return err when connect succeeds and Channel errors", func() {
			testConnection.ChannelImpl = func() (outs.AMQPChannel, error) { return nil, errors.New("Test Channel error") }
			Expect(testOutputSinkRMQ.Start(context.Background())).To(MatchError("Test Channel error"))
		})

		It("should return nil error following successful connection", func() {
			Expect(testOutputSinkRMQ.Start(context.Background())).To(BeNil())
		})
	})

	Context("Send", func() {
		It("should return OutputSinkUnavailable err when connection is nil and Dial returns err", func() {
			testAMQP.DialImpl = func(address string) (outs.AMQPConnection, error) { return nil, errors.New("Test dial error") }
			Expect(testOutputSinkRMQ.Send(map[string]interface{}{"a": "B"})).To(MatchError(
				ContainSubstring("OutputSinkUnavailable"),
			))
		})

		It("should return OutputSinkUnavailable err when connection.IsClosed returns true and Dial returns err", func() {
			// Initialize connection
			testAMQP.Dial("somewhere")

			// Configure to appear as a closed connection and subsequent Dial call to fail
			testConnection.IsClosedImpl = func() bool { return true }
			testAMQP.DialImpl = func(address string) (outs.AMQPConnection, error) { return nil, errors.New("Test dial error") }

			Expect(testOutputSinkRMQ.Send(map[string]interface{}{"a": "B"})).To(MatchError(
				ContainSubstring("OutputSinkUnavailable"),
			))
		})

		It("should return OutputSinkPublishError err when Publish returns err", func() {
			testRabbitMQChannel.PublishImpl = func(s1, s2 string, b1, b2 bool, p amqp.Publishing) error { return errors.New("test publish error") }
			Expect(testOutputSinkRMQ.Send(map[string]interface{}{"a": "B"})).To(MatchError(
				ContainSubstring("OutputSinkPublishError"),
			))
		})

		It("should return nil when connected and Publish returns normally", func() {
			Expect(testOutputSinkRMQ.Send(map[string]interface{}{"a": "B"})).To(BeNil())
		})
	})

})
