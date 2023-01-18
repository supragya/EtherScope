package ethrpc

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PoolNodeMeta struct {
	Identity    string
	IsAlive     bool
	Reports     uint32
	reportMutex sync.Mutex
	FirstReport time.Time
	LastReport  time.Time
	BringAlive  time.Time
}

type PoolNode[I any] struct {
	Item I
	Meta PoolNodeMeta
}

type MSPoolConfig struct {
	WindowSize     uint32
	ToleranceCount uint32
	TimeStep       time.Duration
	RetryTimesteps uint32
}

var DefaultMSPoolConfig MSPoolConfig = MSPoolConfig{
	WindowSize:     80,                     // 10 seconds window size
	ToleranceCount: 5,                      // can tolerate maximum of 5 failures
	TimeStep:       time.Millisecond * 100, // can report a max of 1 failure every 1000 millisec
	RetryTimesteps: 300,                    // retry after 1 minute (60 seconds)
}

type MasterSlavePool[I any] struct {
	config               MSPoolConfig
	rwlock               sync.RWMutex
	log                  logger.Logger
	allFailureLogTime    time.Time
	allFailureCachedItem *I
	itemMap              map[*I]*PoolNode[*I]
	Master               *PoolNode[*I]
	Slaves               []*PoolNode[*I]
	RPCTimeout           time.Duration
}

type DurationTuple[I any] struct {
	Duration time.Duration
	Item     I
}

type DurationTupleList[I any] []DurationTuple[I]

func (a DurationTupleList[I]) Len() int           { return len(a) }
func (a DurationTupleList[I]) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DurationTupleList[I]) Less(i, j int) bool { return a[i].Duration < a[j].Duration }

func NewNode[I any](item *I, identity string) PoolNode[*I] {
	return PoolNode[*I]{
		Item: item,
		Meta: PoolNodeMeta{
			Identity:    identity,
			IsAlive:     true,
			Reports:     0,
			reportMutex: sync.Mutex{},
			FirstReport: time.Time{},
			LastReport:  time.Time{},
			BringAlive:  time.Time{},
		},
	}
}

func NewEthClientMasterSlavePool(masterURL string,
	slaveURLs []string,
	config MSPoolConfig,
	timeout time.Duration,
	log logger.Logger) (*MasterSlavePool[ethclient.Client], error) {
	itemMap := make(map[*ethclient.Client]*PoolNode[*ethclient.Client], len(slaveURLs)+1)

	// Setup master
	ms, err := ethclient.Dial(masterURL)
	if err != nil {
		return nil, err
	}
	master := NewNode(ms, "master")
	itemMap[ms] = &master

	// Setup slaves
	slaves := []*PoolNode[*ethclient.Client]{}
	for idx, url := range slaveURLs {
		cl, err := ethclient.Dial(url)
		if err != nil {
			return nil, err
		}
		slave := NewNode(cl, fmt.Sprintf("slave%v", idx))
		itemMap[cl] = &slave
		slaves = append(slaves, &slave)
	}

	return &MasterSlavePool[ethclient.Client]{
		config:               config,
		log:                  log,
		rwlock:               sync.RWMutex{},
		allFailureLogTime:    time.Time{},
		allFailureCachedItem: nil,
		itemMap:              itemMap,
		Master:               &master,
		Slaves:               slaves,
		RPCTimeout:           timeout,
	}, nil
}

func (m *MasterSlavePool[I]) Report(item *I, timedOut bool) error {
	if !timedOut {
		return nil
	}
	pn, ok := m.itemMap[item]
	if !ok {
		return errors.New("item not found")
	}

	now := time.Now()

	// Short circuit in case of pool being not alive or last report too recently
	if !pn.Meta.IsAlive || now.Sub(pn.Meta.LastReport) < m.config.TimeStep {
		return nil
	}

	pn.Meta.reportMutex.Lock()
	defer pn.Meta.reportMutex.Unlock()

	now = time.Now()
	// Maybe somebody else reported while we were waiting for lock
	// Short circuit in case of pool being not alive or last report too recently
	if !pn.Meta.IsAlive || now.Sub(pn.Meta.LastReport) < m.config.TimeStep {
		return nil
	}

	pn.Meta.LastReport = now

	// We forget all the failures accrued till now if counter addition
	// start time has been since long.
	if now.Sub(pn.Meta.FirstReport) > m.config.TimeStep*time.Duration(m.config.WindowSize) {
		pn.Meta.FirstReport = now
		pn.Meta.Reports = 1
	} else {
		pn.Meta.Reports += 1
	}

	// If more than enough of timeSteps have resulted in failure, go to cooldown
	if pn.Meta.Reports > m.config.ToleranceCount && now.Sub(pn.Meta.FirstReport) < m.config.TimeStep*time.Duration(m.config.WindowSize) {
		m.log.Warn("mspool reports upstream failure", "upstream", pn.Meta.Identity)
		pn.Meta.IsAlive = false
		pn.Meta.BringAlive = now.Add(m.config.TimeStep * time.Duration(m.config.RetryTimesteps))
	}

	return nil
}

func (m *MasterSlavePool[I]) GetItem() (*I, *PoolNodeMeta) {
	// Lock global RW lock for reads
	m.rwlock.RLock()

	// Check if master is alive, if so return master
	if m.Master.Meta.IsAlive {
		m.rwlock.RUnlock()
		return m.Master.Item, &m.Master.Meta
	}

	// If master is not alive, check if time has come to
	// recheck on master
	now := time.Now()
	if m.Master.Meta.BringAlive.Sub(now) <= time.Duration(0) {
		m.rwlock.RUnlock()
		m.rwlock.Lock()
		MakeAlive(&m.Master.Meta)
		m.rwlock.Unlock()
		return m.Master.Item, &m.Master.Meta
	}

	// If master is not alive, nor is the time to bring it
	// back online, check if any of the slaves is ready.
	for _, slave := range m.Slaves {
		sm := &slave.Meta
		if sm.IsAlive {
			m.rwlock.RUnlock()
			return slave.Item, &slave.Meta
		}
		if sm.BringAlive.Sub(now) <= time.Duration(0) {
			m.rwlock.RUnlock()
			m.rwlock.Lock()
			MakeAlive(&slave.Meta)
			m.rwlock.Unlock()
			return slave.Item, &slave.Meta
		}
	}

	// If none of the others were successfully, we may have to
	// wait till first rpc comes back online and send it
	m.rwlock.RUnlock()

	// Very expensive proposition, lot of mutex lock unlocks
	return m.allFailureRecovery()
}

func (m *MasterSlavePool[I]) allFailureRecovery() (*I, *PoolNodeMeta) {
	currentTime := time.Now()

	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	// Critical section below
	// The first thread that enters below does the hefty work of
	// sort and wait, sets for a timeStep a cached item. Rest of
	// the threads pick this item and return

	if m.allFailureLogTime.Sub(currentTime) <= time.Duration(0) {
		// First thread doing hefty work
		m.log.Warn("critical rpc failure. all upstreams in cooldown state. blocking application")
	} else {
		// If allFailureLogTime is in future, it can only be done by
		// another thread which set this up recently. We can use the cached response hence.
		return m.allFailureCachedItem, nil
	}

	list := DurationTupleList[*PoolNode[*I]]{}

	list = append(list, DurationTuple[*PoolNode[*I]]{
		Duration: m.Master.Meta.BringAlive.Sub(currentTime),
		Item:     m.Master,
	})

	for _, slave := range m.Slaves {
		list = append(list, DurationTuple[*PoolNode[*I]]{
			Duration: slave.Meta.BringAlive.Sub(currentTime),
			Item:     slave,
		})
	}

	sort.Sort(list)

	minDuration := list[0].Duration
	time.Sleep(minDuration)

	// Cleanup
	for _, tuple := range list {
		if tuple.Duration == minDuration {
			MakeAlive(&tuple.Item.Meta)
		}
	}

	m.allFailureLogTime = currentTime.Add(m.config.TimeStep)
	m.allFailureCachedItem = list[0].Item.Item

	return m.allFailureCachedItem, &list[0].Item.Meta
}

func MakeAlive(m *PoolNodeMeta) {
	m.reportMutex.Lock()
	defer m.reportMutex.Unlock()

	m.IsAlive = true
	m.Reports = 0
	m.FirstReport = time.Time{}
	m.LastReport = time.Time{}
	m.BringAlive = time.Time{}
}

// No Read lock based approach, may be wrong but avoids
// locking / unlocking saving mutex access for other goroutines
func (m *MasterSlavePool[I]) PeriodicRecording(dur time.Duration) {
	for {
		time.Sleep(dur)
		updates := make([]interface{}, len(m.Slaves)*2)

		updates = append(updates, m.Master.Meta.Identity, alStr(m.Master.Meta.IsAlive))
		for _, slv := range m.Slaves {
			updates = append(updates, slv.Meta.Identity, alStr(slv.Meta.IsAlive))
		}
		m.log.Info("master slave pool upstream report", updates...)
	}
}

func alStr(isAlive bool) string {
	if isAlive {
		return "alive"
	}
	return "down"
}
