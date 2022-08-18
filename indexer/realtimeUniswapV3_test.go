package indexer

import (
	"sync"
	"testing"

	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/stretchr/testify/assert"
)

func TestUniswapV3Swap(t *testing.T) {
	var (
		_log  = loadLog(t, "../test/uniswapV3SwapExample.json")
		ri    = NewRealtimeIndexer(0, []string{"https://rpc.ankr.com/eth"}, &db.DBConn{ChainID: 1}, []string{})
		bm    = itypes.BlockSynopsis{}
		mt    = sync.Mutex{}
		items []interface{}
	)

	ri.processUniV3Swap(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one swap not found")
	assert.Equal(t, itypes.BlockSynopsis{SwapLogs: 1, TotalLogs: 1}, bm, "one swap not found")
}
