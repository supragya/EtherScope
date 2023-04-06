package sink

import "github.com/Blockpour/Blockpour-Geth-Indexer/algorand/service"

type OutputSink interface {
	service.Service

	Send(payload interface{}) error
}
