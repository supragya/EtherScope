package sink

import "github.com/supragya/EtherScope/algorand/service"

type OutputSink interface {
	service.Service

	Send(payload interface{}) error
}
