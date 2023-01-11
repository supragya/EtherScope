package outputsink

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
)

type OutputSink interface {
	service.Service

	Send(payload interface{}) error
	/*
		Predicate expression detailing the readiness of the output sink to accept input
		This function will return false if there is no active connection to the output sink
		or there are still cached messages which have not been pushed to the MQ following reconnect
	*/
	IsReady() bool

	/*
		Function which can be called to trigger reconnection to the output sink.
		Initial connection will be automatically created at startup via the service start up hook. This
		method can be used if the initial connection fails or to recover from intermittent connection
		issues.
	*/
	Reconnect() error
}
