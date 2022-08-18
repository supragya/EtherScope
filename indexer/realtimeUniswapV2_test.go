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

	ri.processUniV2Swap(_log, &items, &bm, &mt)
	assert.Equal(t, len(items), 1, "one swap not found")
	assert.Equal(t, bm, itypes.BlockSynopsis{SwapLogs: 1, TotalLogs: 1}, "one swap not found")
}
