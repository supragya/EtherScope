package indexer

import (
	"sync"
	"testing"

	"github.com/Blockpour/Blockpour-Geth-Indexer/config"
	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/stretchr/testify/assert"
)

func TestERC20Transfer(t *testing.T) {
	util.ENOK(logger.SetLogLevel("error"))
	util.ENOK(config.LoadViperConfig("testdata/configs/testcfg.yaml"))
	var (
		_log = loadLog(t, "testdata/transferExample.json")
		ri   = NewRealtimeIndexer(0,
			"https://rpc.ankr.com/eth",
			[]string{},
			false,
			&db.DBConn{ChainID: 1},
			[]string{})
		bm    = itypes.BlockSynopsis{}
		mt    = sync.Mutex{}
		items []interface{}
	)

	ri.processTransfer(_log, &items, &bm, &mt)
	assert.Equal(t, 1, len(items), "one transfer not found")
	assert.Equal(t, itypes.BlockSynopsis{TransferLogs: 1, TotalLogs: 1}, bm, "one transfer not found")
}
