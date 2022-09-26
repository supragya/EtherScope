package indexer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/config"
	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func loadLog(t *testing.T, file string) types.Log {
	_log := types.Log{}
	jsonFile, err := os.Open(file)
	if err != nil {
		t.Error(err)
	}
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &_log)
	return _log
}

func loadRawJSON(t *testing.T, file string) string {
	jsonFile, err := os.Open(file)
	if err != nil {
		t.Error(err)
	}
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return string(byteValue)
}

func TestUniswapV2Mint(t *testing.T) {
	util.ENOK(logger.SetLogLevel("error"))
	util.ENOK(config.LoadViperConfig("testdata/configs/testcfg.yaml"))
	var (
		_log = loadLog(t, "testdata/uniswapV2MintExample.json")
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

	ri.processMint(_log, &items, &bm, &mt)

	assert.Equal(t, 1, len(items), "one mint not found")
	u, _ := json.MarshalIndent(items[0], " ", "  ")
	assert.JSONEq(t, loadRawJSON(t, "testdata/uniswapV2MintExampleExpected.json"), string(u))
	assert.Equal(t, itypes.BlockSynopsis{MintLogs: 1, TotalLogs: 1}, bm, "one mint not found")
}

func TestUniswapV2Burn(t *testing.T) {
	util.ENOK(logger.SetLogLevel("error"))
	util.ENOK(config.LoadViperConfig("testdata/configs/testcfg.yaml"))
	var (
		_log = loadLog(t, "testdata/uniswapV2BurnExample.json")
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

	ri.processBurn(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one burn not found")
	assert.Equal(t, itypes.BlockSynopsis{BurnLogs: 1, TotalLogs: 1}, bm, "one burn not found")
}

func TestUniswapV2Swap(t *testing.T) {
	util.ENOK(logger.SetLogLevel("error"))
	util.ENOK(config.LoadViperConfig("testdata/configs/testcfg.yaml"))
	var (
		_log = loadLog(t, "testdata/uniswapV2SwapExample.json")
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

	ri.processUniV2Swap(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one swap not found")
	// u, _ := json.MarshalIndent(items[0], " ", "  ")
	// assert.JSONEq(t, loadRawJSON(t, "testdata/uniswapV2SwapExampleExpected.json"), string(u))
	assert.Equal(t, itypes.BlockSynopsis{SwapLogs: 1, TotalLogs: 1}, bm, "one swap not found")
}
