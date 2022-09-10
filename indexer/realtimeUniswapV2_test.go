package indexer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
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

func TestUniswapV2Mint(t *testing.T) {
	var (
		_log  = loadLog(t, "../test/uniswapV2MintExample.json")
		ri    = NewRealtimeIndexer(0, []string{"https://rpc.ankr.com/eth"}, &db.DBConn{ChainID: 1}, []string{})
		bm    = itypes.BlockSynopsis{}
		mt    = sync.Mutex{}
		items []interface{}
	)

	ri.processMint(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one mint not found")
	assert.Equal(t, itypes.BlockSynopsis{MintLogs: 1, TotalLogs: 1}, bm, "one mint not found")
}

func TestUniswapV2Burn(t *testing.T) {
	var (
		_log  = loadLog(t, "../test/uniswapV2BurnExample.json")
		ri    = NewRealtimeIndexer(0, []string{"https://rpc.ankr.com/eth"}, &db.DBConn{ChainID: 1}, []string{})
		bm    = itypes.BlockSynopsis{}
		mt    = sync.Mutex{}
		items []interface{}
	)

	ri.processBurn(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one burn not found")
	assert.Equal(t, itypes.BlockSynopsis{BurnLogs: 1, TotalLogs: 1}, bm, "one burn not found")
}

func TestUniswapV2Swap(t *testing.T) {
	var (
		_log  = loadLog(t, "../test/uniswapV2SwapExample.json")
		ri    = NewRealtimeIndexer(0, []string{"https://rpc.ankr.com/eth"}, &db.DBConn{ChainID: 1}, []string{})
		bm    = itypes.BlockSynopsis{}
		mt    = sync.Mutex{}
		items []interface{}
	)

	ri.processUniV2Swap(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one swap not found")
	assert.Equal(t, itypes.BlockSynopsis{SwapLogs: 1, TotalLogs: 1}, bm, "one swap not found")
}
