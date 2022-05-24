package indexer

import (
	"errors"
	"sort"
	"sync"

	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	WD       = 20 // should be even
	TimedOut = 10000
)

type LNode struct {
	Item      *ethclient.Client
	Latency   float64
	Window    [WD]float64
	WindowPtr uint
	TimedOut  uint
	ResetCtr  uint64
}

type LatencySortedPool struct {
	lock     *sync.Mutex
	ctr      uint64
	items    []*LNode
	itemsMap map[*ethclient.Client]*LNode
}

func NewLatencySortedPool(items []string) *LatencySortedPool {
	lsp := &LatencySortedPool{
		lock:     &sync.Mutex{},
		ctr:      0,
		items:    []*LNode{},
		itemsMap: make(map[*ethclient.Client]*LNode),
	}
	for _, item := range items {
		cl, err := ethclient.Dial(item)
		util.ENOK(err)
		node := &LNode{
			Item:      cl,
			Latency:   0,
			Window:    [WD]float64{},
			WindowPtr: 0,
			ResetCtr:  WD / 2,
		}
		lsp.items = append(lsp.items, node)
		lsp.itemsMap[cl] = node
	}
	return lsp
}

func (l *LatencySortedPool) Len() int {
	return len(l.items)
}

func (l *LatencySortedPool) Less(i, j int) bool {
	return l.items[i].Latency < l.items[j].Latency
}

func (l *LatencySortedPool) Swap(i, j int) {
	temp := l.items[i]
	l.items[i] = l.items[j]
	l.items[j] = temp
}

func (l *LatencySortedPool) Report(item *ethclient.Client, latency float64, timedOut bool) error {
	ptr, ok := l.itemsMap[item]
	if !ok {
		return errors.New("unknown item ")
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	l.ctr++
	// Keep a counter of number of timeouts on this upstream
	if timedOut {
		ptr.TimedOut++
	}
	// If timed out enough number of times, do not allow
	// upstream to be reset quickly.
	if ptr.TimedOut > WD/2 {
		ptr.ResetCtr = l.ctr + WD*1000
		ptr.TimedOut = 0
	} else {
		ptr.ResetCtr = l.ctr + WD
	}
	ptr.Latency = ptr.Latency + (latency-ptr.Window[ptr.WindowPtr])/WD
	ptr.Window[ptr.WindowPtr] = latency
	ptr.WindowPtr = (ptr.WindowPtr + 1) % WD

	// Check if any of the nodes needs to be reset
	// Timeout count is not reset however
	for _, lnode := range l.items {
		if lnode.ResetCtr <= l.ctr {
			lnode.Latency = 0
			for i := 0; i < WD; i++ {
				lnode.Window[i] = 0
			}
			lnode.WindowPtr = 0
			lnode.ResetCtr = l.ctr + WD
		}
	}

	sort.Sort(l)
	return nil
}

func (l *LatencySortedPool) GetItem() *ethclient.Client {
	return l.items[0].Item
}
