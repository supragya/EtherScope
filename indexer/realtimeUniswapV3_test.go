package indexer

import (
	"sync"
	"testing"

	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/stretchr/testify/assert"
)

func TestUniswapV3Swap(t *testing.T) {
	ri := testingSetup()

	var (
		file     = "testdata/uniswapV3SwapExample"
		_log     = loadLog(t, file+".json")
		expected = itypes.Swap{}
		bm       = itypes.BlockSynopsis{}
		mt       = sync.Mutex{}
		items    []interface{}
	)

	loadJson(t, file+"Expected.json", &expected)
	ri.processUniV3Swap(_log, &items, &bm, &mt)

	assert.Equal(t, 1, len(items), "one transfer not found")
	assert.Equal(t, itypes.BlockSynopsis{SwapLogs: 1, TotalLogs: 1}, bm, "one swap not found")

	item := items[0].(itypes.Swap)
	// u, _ := json.MarshalIndent(item, "", " ")
	// fmt.Printf("%s", u)
	assert.Equal(t, expected.Type, item.Type, "does not match")
	assert.Equal(t, expected.Network, item.Network, "does not match")
	assert.Equal(t, expected.LogIdx, item.LogIdx, "does not match")
	assert.Equal(t, expected.Transaction, item.Transaction, "does not match")
	assert.Equal(t, expected.Time, item.Time, "does not match")
	assert.Equal(t, expected.Height, item.Height, "does not match")
	assert.Equal(t, expected.Sender, item.Sender, "does not match")
	assert.Equal(t, expected.PairContract, item.PairContract, "does not match")
	assert.Equal(t, expected.Token0, item.Token0, "does not match")
	assert.Equal(t, expected.Token1, item.Token1, "does not match")
	assertBigFloatClose(t, expected.Amount0, item.Amount0, nil)
	assertBigFloatClose(t, expected.Amount1, item.Amount1, nil)
	assertBigFloatClose(t, expected.Reserve0, item.Reserve0, nil)
	assertBigFloatClose(t, expected.Reserve1, item.Reserve1, nil)
	assertBigFloatClose(t, expected.AmountUSD, item.AmountUSD, nil)
	assertBigFloatClose(t, expected.Price0, item.Price0, nil)
	assertBigFloatClose(t, expected.Price1, item.Price1, nil)
}

func TestUniswapV3Mint(t *testing.T) {
	ri := testingSetup()

	var (
		file     = "testdata/uniswapV3MintExample"
		_log     = loadLog(t, file+".json")
		expected = itypes.Mint{}
		bm       = itypes.BlockSynopsis{}
		mt       = sync.Mutex{}
		items    []interface{}
		err      error
	)

	bm.Time, err = ri.da.GetBlockTimestamp(_log.BlockNumber)
	assert.Nil(t, err, "error is not nil")

	loadJson(t, file+"Expected.json", &expected)
	ri.processUniV3Mint(_log, &items, &bm, &mt)

	assert.Equal(t, 1, len(items), "one transfer not found")
	assert.Equal(t, itypes.BlockSynopsis{MintLogs: 1, TotalLogs: 1, Time: bm.Time}, bm, "one mint not found")

	item := items[0].(itypes.Mint)
	// u, _ := json.MarshalIndent(item, "", " ")
	// fmt.Printf("%s", u)
	assert.Equal(t, expected.Type, item.Type, "does not match")
	assert.Equal(t, expected.Network, item.Network, "does not match")
	assert.Equal(t, expected.LogIdx, item.LogIdx, "does not match")
	assert.Equal(t, expected.Transaction, item.Transaction, "does not match")
	assert.Equal(t, expected.Time, item.Time, "does not match")
	assert.Equal(t, expected.Height, item.Height, "does not match")
	assert.Equal(t, expected.Sender, item.Sender, "does not match")
	assert.Equal(t, expected.PairContract, item.PairContract, "does not match")
	assert.Equal(t, expected.Token0, item.Token0, "does not match")
	assert.Equal(t, expected.Token1, item.Token1, "does not match")
	assertBigFloatClose(t, expected.Amount0, item.Amount0, nil)
	assertBigFloatClose(t, expected.Amount1, item.Amount1, nil)
	assertBigFloatClose(t, expected.Reserve0, item.Reserve0, nil)
	assertBigFloatClose(t, expected.Reserve1, item.Reserve1, nil)
	assertBigFloatClose(t, expected.AmountUSD, item.AmountUSD, nil)
	assertBigFloatClose(t, expected.Price0, item.Price0, nil)
	assertBigFloatClose(t, expected.Price1, item.Price1, nil)
}

func TestUniswapV3Burn(t *testing.T) {
	ri := testingSetup()

	var (
		file     = "testdata/uniswapV3BurnExample"
		_log     = loadLog(t, file+".json")
		expected = itypes.Mint{}
		bm       = itypes.BlockSynopsis{}
		mt       = sync.Mutex{}
		items    []interface{}
		err      error
	)

	bm.Time, err = ri.da.GetBlockTimestamp(_log.BlockNumber)
	assert.Nil(t, err, "error is not nil")

	loadJson(t, file+"Expected.json", &expected)
	ri.processUniV3Burn(_log, &items, &bm, &mt)

	assert.Equal(t, 1, len(items), "one transfer not found")
	assert.Equal(t, itypes.BlockSynopsis{BurnLogs: 1, TotalLogs: 1, Time: bm.Time}, bm, "one mint not found")

	item := items[0].(itypes.Burn)
	// u, _ := json.MarshalIndent(item, "", " ")
	// fmt.Printf("%s", u)
	assert.Equal(t, expected.Type, item.Type, "does not match")
	assert.Equal(t, expected.Network, item.Network, "does not match")
	assert.Equal(t, expected.LogIdx, item.LogIdx, "does not match")
	assert.Equal(t, expected.Transaction, item.Transaction, "does not match")
	assert.Equal(t, expected.Time, item.Time, "does not match")
	assert.Equal(t, expected.Height, item.Height, "does not match")
	assert.Equal(t, expected.Sender, item.Sender, "does not match")
	assert.Equal(t, expected.PairContract, item.PairContract, "does not match")
	assert.Equal(t, expected.Token0, item.Token0, "does not match")
	assert.Equal(t, expected.Token1, item.Token1, "does not match")
	assertBigFloatClose(t, expected.Amount0, item.Amount0, nil)
	assertBigFloatClose(t, expected.Amount1, item.Amount1, nil)
	assertBigFloatClose(t, expected.Reserve0, item.Reserve0, nil)
	assertBigFloatClose(t, expected.Reserve1, item.Reserve1, nil)
	assertBigFloatClose(t, expected.AmountUSD, item.AmountUSD, nil)
	assertBigFloatClose(t, expected.Price0, item.Price0, nil)
	assertBigFloatClose(t, expected.Price1, item.Price1, nil)
}
