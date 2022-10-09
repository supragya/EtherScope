package indexer

import (
	"sync"
	"testing"

	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/stretchr/testify/assert"
)

func TestERC20Transfer(t *testing.T) {
	ri := testingSetup()

	var (
		file     = "testdata/erc20TransferExample"
		_log     = loadLog(t, file+".json")
		expected = itypes.Transfer{}
		bm       = itypes.BlockSynopsis{}
		mt       = sync.Mutex{}
		items    []interface{}
	)

	loadJson(t, file+"Expected.json", &expected)
	ri.processERC20Transfer(_log, &items, &bm, &mt)

	assert.Equal(t, 1, len(items), "one transfer not found")
	assert.Equal(t, itypes.BlockSynopsis{TransferLogs: 1, TotalLogs: 1}, bm, "one transfer not found")

	item := items[0].(itypes.Transfer)
	assert.Equal(t, expected.Type, item.Type, "does not match")
	assert.Equal(t, expected.Network, item.Network, "does not match")
	assert.Equal(t, expected.LogIdx, item.LogIdx, "does not match")
	assert.Equal(t, expected.Transaction, item.Transaction, "does not match")
	assert.Equal(t, expected.Time, item.Time, "does not match")
	assert.Equal(t, expected.Height, item.Height, "does not match")
	assert.Equal(t, expected.Token, item.Token, "does not match")
	assert.Equal(t, expected.Sender, item.Sender, "does not match")
	assert.Equal(t, expected.Receiver, item.Receiver, "does not match")
	assertBigFloatClose(t, expected.Amount, item.Amount, nil)
	assertBigFloatClose(t, expected.AmountUSD, item.AmountUSD.Price, nil)
}
