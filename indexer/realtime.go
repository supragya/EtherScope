package indexer

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	EUninitialized = errors.New("uninitialized indexer")
)

type RealtimeIndexer struct {
	currentHeight int64
	indexedHeight int64
	upstreams     *LatencySortedPool

	quitCh chan struct{}
}

func NewRealtimeIndexer(indexedHeight int64, upstreams []string) *RealtimeIndexer {
	return &RealtimeIndexer{
		currentHeight: 0,
		indexedHeight: indexedHeight,
		upstreams:     NewLatencySortedPool(upstreams),
		quitCh:        make(chan struct{}),
	}
}

func (r *RealtimeIndexer) Start() error {
	if r.indexedHeight == 0 || r.upstreams.Len() == 0 {
		return EUninitialized
	}
	r.ridxLoop()
	time.Sleep(time.Second * 2)
	return nil
}

func (r *RealtimeIndexer) ridxLoop() {
	for {
		select {
		case <-time.After(time.Second):
			log.Info("timed chan went off")
		case <-r.quitCh:
			log.Info("quitting realtime indexer")
		}
	}
}

func (r *RealtimeIndexer) Init() error {
	return nil
}

func (r *RealtimeIndexer) Quit() {
	r.quitCh <- struct{}{}
}
