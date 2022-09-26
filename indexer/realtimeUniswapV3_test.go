package indexer

import (
	"sync"
	"testing"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/config"
	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/stretchr/testify/assert"
)

func TestUniswapV3Swap(t *testing.T) {
	util.ENOK(logger.SetLogLevel("error"))
	util.ENOK(config.LoadViperConfig("testdata/configs/testcfg.yaml"))
	var (
		_log = loadLog(t, "testdata/uniswapV3SwapExample.json")
		ri   = NewRealtimeIndexer(0,
			"https://rpc.ankr.com/eth",
			[]string{},
			time.Second,
			false,
			&db.DBConn{ChainID: 1},
			[]string{})
		bm    = itypes.BlockSynopsis{}
		mt    = sync.Mutex{}
		items []interface{}
	)

	ri.processUniV3Swap(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one swap not found")
	assert.Equal(t, itypes.BlockSynopsis{SwapLogs: 1, TotalLogs: 1}, bm, "one swap not found")
}

func TestUniswapV3IncreaseLiquidity(t *testing.T) {
	util.ENOK(logger.SetLogLevel("error"))
	util.ENOK(config.LoadViperConfig("testdata/configs/testcfg.yaml"))
	var (
		_log = loadLog(t, "testdata/uniswapV3IncreaseLiquidityExample.json")
		ri   = NewRealtimeIndexer(0,
			"https://rpc.ankr.com/eth",
			[]string{},
			time.Second,
			false,
			&db.DBConn{ChainID: 1},
			[]string{})
		bm    = itypes.BlockSynopsis{}
		mt    = sync.Mutex{}
		items []interface{}
	)

	ri.processMintV3(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one mint not found")
	assert.Equal(t, itypes.BlockSynopsis{MintLogs: 1, TotalLogs: 1}, bm, "one mint not found")
}

func TestUniswapV3DecreaseLiquidity(t *testing.T) {
	util.ENOK(logger.SetLogLevel("error"))
	util.ENOK(config.LoadViperConfig("testdata/configs/testcfg.yaml"))
	var (
		_log = loadLog(t, "testdata/uniswapV3DecreaseLiquidityExample.json")
		ri   = NewRealtimeIndexer(0,
			"https://rpc.ankr.com/eth",
			[]string{},
			time.Second,
			false,
			&db.DBConn{ChainID: 1},
			[]string{})
		bm    = itypes.BlockSynopsis{}
		mt    = sync.Mutex{}
		items []interface{}
	)

	ri.processBurnV3(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one burn not found")
	assert.Equal(t, itypes.BlockSynopsis{BurnLogs: 1, TotalLogs: 1}, bm, "one burn not found")
}
