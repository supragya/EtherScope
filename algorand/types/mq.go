package types

import (
	"github.com/google/uuid"
)

type Payload struct {
	NodeMoniker   string
	NodeID        uuid.UUID
	NodeVersion   string
	Environment   string
	Network       string
	BlockSynopsis *BlockSynopsis
	Items         []interface{}
}

func (p *Payload) Add(item interface{}) {
	p.Items = append(p.Items, item)
}
