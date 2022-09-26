package mspool

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"

	"github.com/stretchr/testify/assert"
)

type IS struct {
	item int
}

func TestBasic(t *testing.T) {
	logger.SetupLog()
	util.ENOK(logger.SetLogLevel("fatal"))
	pool := newIntegerPool(1)
	item := pool.GetItem()
	assert.Equal(t, *item, IS{0})
}

func TestSwitchoverMasterFailureSimulated(t *testing.T) {
	logger.SetupLog()
	util.ENOK(logger.SetLogLevel("fatal"))
	pool := newIntegerPool(1)
	item := pool.GetItem()
	assert.Equal(t, *item, IS{0})
	for i := 0; i < 9; i++ {
		assert.Equal(t, nil, pool.Report(item, true))
		time.Sleep(time.Millisecond * 11)
	}
	newItem := pool.GetItem()
	assert.Equal(t, *newItem, IS{1})
}

func TestNoSwitchoverMasterIntermittentSimulated0(t *testing.T) {
	logger.SetupLog()
	util.ENOK(logger.SetLogLevel("fatal"))
	pool := newIntegerPool(1)
	item := pool.GetItem()
	assert.Equal(t, *item, IS{0})
	for i := 0; i < 20; i++ {
		assert.Equal(t, nil, pool.Report(item, true))
		time.Sleep(time.Millisecond * 25)
	}
	newItem := pool.GetItem()
	assert.Equal(t, *newItem, IS{0})
}

func TestNoSwitchoverMasterIntermittentSimulated1(t *testing.T) {
	logger.SetupLog()
	util.ENOK(logger.SetLogLevel("fatal"))
	pool := newIntegerPool(1)
	item := pool.GetItem()
	assert.Equal(t, *item, IS{0})
	for i := 0; i < 1000; i++ {
		assert.Equal(t, nil, pool.Report(item, true))
	}
	newItem := pool.GetItem()
	assert.Equal(t, *newItem, IS{0})
}

func TestNoSwitchoverMasterIntermittentSimulated2(t *testing.T) {
	logger.SetupLog()
	util.ENOK(logger.SetLogLevel("fatal"))
	pool := newIntegerPool(1)
	item := pool.GetItem()
	assert.Equal(t, *item, IS{0})
	for i := 0; i < 10000; i++ {
		go func(pool *MasterSlavePool[IS]) {
			assert.Equal(t, nil, pool.Report(item, true))
			threadItem := pool.GetItem()
			assert.Equal(t, *threadItem, IS{0})
		}(&pool)
	}
	newItem := pool.GetItem()
	assert.Equal(t, *newItem, IS{0})
}

func TestSwitchoverMasterFailureReal(t *testing.T) {
	logger.SetupLog()
	util.ENOK(logger.SetLogLevel("fatal"))
	pool := newIntegerPool(1)
	item := pool.GetItem()
	assert.Equal(t, *item, IS{0})
	startTime := time.Now()
	for time.Since(startTime) < time.Millisecond*100 {
		assert.Equal(t, nil, pool.Report(item, true))
	}
	newItem := pool.GetItem()
	assert.Equal(t, *newItem, IS{1})
}

func TestRetryMasterSimulated(t *testing.T) {
	logger.SetupLog()
	util.ENOK(logger.SetLogLevel("fatal"))
	pool := newIntegerPool(1)
	item := pool.GetItem()
	assert.Equal(t, *item, IS{0})
	for i := 0; i < 9; i++ {
		assert.Equal(t, nil, pool.Report(item, true))
		time.Sleep(time.Millisecond * 11)
	}
	newItem := pool.GetItem()
	assert.Equal(t, *newItem, IS{1})
	time.Sleep(time.Millisecond * 10 * 100)
	finalItem := pool.GetItem()
	assert.Equal(t, *finalItem, IS{0})
}

func TestApplicationBlockageShouldProgress(t *testing.T) {
	logger.SetupLog()
	util.ENOK(logger.SetLogLevel("fatal"))
	pool := newIntegerPool(1)
	item := pool.GetItem()
	assert.Equal(t, *item, IS{0})
	for i := 0; i < 3000; i++ {
		go func(pool *MasterSlavePool[IS]) {
			threadItem := pool.GetItem()
			assert.Equal(t, nil, pool.Report(threadItem, true))
		}(&pool)
		time.Sleep(time.Millisecond)
	}
	time.Sleep(time.Millisecond * 10 * 100)
	newItem := pool.GetItem()
	assert.Equal(t, *newItem, IS{0})
}

// Returns a 10ms timestep MS pool
func newIntegerPool(slaveCount int) MasterSlavePool[IS] {
	itemMap := make(map[*IS]*PoolNode[*IS], slaveCount+1)
	masterNode := &IS{0}
	master := NewNode(masterNode, "master")
	itemMap[masterNode] = &master

	slaves := []*PoolNode[*IS]{}
	for idx := 1; idx < slaveCount+1; idx++ {
		slaveNode := &IS{idx}
		slave := NewNode(slaveNode, fmt.Sprintf("slave%v", idx-1))
		itemMap[slaveNode] = &slave
		slaves = append(slaves, &slave)
	}

	var cfg MSPoolConfig = MSPoolConfig{
		WindowSize:     20,
		ToleranceCount: 8,
		TimeStep:       time.Millisecond * 10,
		RetryTimesteps: 100,
	}

	return MasterSlavePool[IS]{
		config:               cfg,
		rwlock:               sync.RWMutex{},
		allFailureLogTime:    time.Time{},
		allFailureCachedItem: nil,
		itemMap:              itemMap,
		Master:               &master,
		Slaves:               slaves,
	}
}
