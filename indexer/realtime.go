package indexer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
			util.ENOK(r.populateCurrentHeight())
			if r.currentHeight == r.indexedHeight {
				continue
			}
			log.Info("indexing height: ", r.indexedHeight+1, " to ", r.currentHeight)
			logs, err := r.getLogs(ethereum.FilterQuery{
				FromBlock: big.NewInt(r.indexedHeight + 1),
				ToBlock:   big.NewInt(r.currentHeight),
				Topics:    [][]common.Hash{{MintTopic, BurnTopic}},
			})
			if err != nil {
				log.Error(err)
				continue
			}

			for _, l := range logs {
				log.Info("got log: ", l)
			}

			r.indexedHeight = r.currentHeight
		case <-r.quitCh:
			log.Info("quitting realtime indexer")
		}
	}
}

func (r *RealtimeIndexer) Stop() error {
	return nil
}

func (r *RealtimeIndexer) Init() error {
	if err := r.populateCurrentHeight(); err != nil {
		return err
	}
	log.Info("initializing realtime indexer, indexedHeight: "+fmt.Sprint(r.indexedHeight),
		" currentHeight: "+fmt.Sprint(r.currentHeight))
	return nil
}

func (r *RealtimeIndexer) getLogs(fq ethereum.FilterQuery) ([]types.Log, error) {
	var logs []types.Log
	var retries = 0
	for {
		if retries == WD {
			return logs, errors.New("could not fetch logs, retried " + fmt.Sprint(WD) + " times")
		}
		cl := r.upstreams.GetItem()
		var err error

		start := time.Now()
		logs, err = cl.FilterLogs(context.Background(), fq)
		r.upstreams.Report(cl, time.Now().Sub(start).Seconds(), err != nil)
		if err == nil {
			break
		}
		retries++
	}
	return logs, nil
}

func (r *RealtimeIndexer) populateCurrentHeight() error {
	var currentHeight uint64 = 0
	var retries = 0
	for {
		if retries == WD {
			return errors.New("could not init realtime indexer, retried " + fmt.Sprint(WD) + " times")
		}
		cl := r.upstreams.GetItem()
		var err error

		start := time.Now()
		currentHeight, err = cl.BlockNumber(context.Background())
		r.upstreams.Report(cl, time.Now().Sub(start).Seconds(), err != nil)
		if err == nil {
			break
		}
		retries++
	}
	r.currentHeight = int64(currentHeight)
	return nil
}

func (r *RealtimeIndexer) Status() interface{} {
	return nil
}

func (r *RealtimeIndexer) Quit() {
	r.quitCh <- struct{}{}
}
