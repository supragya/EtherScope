package outputsink

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
)

type OutputSink interface {
	service.Service

	Send(payload interface{}) error
}
