package indexer

import (
	"sync"
	"testing"

	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/stretchr/testify/assert"
)

func TestERC20Transfer(t *testing.T) {
	var (
		_log  = loadLog(t, "../test/transferExample.json")
		ri    = NewRealtimeIndexer(0, []string{"https://rpc.ankr.com/eth"}, &db.DBConn{ChainID: 1}, []string{})
		bm    = itypes.BlockSynopsis{}
		mt    = sync.Mutex{}
		items []interface{}
	)

	ri.processTransfer(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one transfer not found")
	assert.Equal(t, itypes.BlockSynopsis{TransferLogs: 1, TotalLogs: 1}, bm, "one transfer not found")
}
