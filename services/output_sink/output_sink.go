package outputsink

import (
	"github.com/supragya/EtherScope/libs/service"
)

type OutputSink interface {
	service.Service

	Send(payload interface{}) error
}
