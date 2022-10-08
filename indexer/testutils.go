package indexer

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/config"
	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	"github.com/Blockpour/Blockpour-Geth-Indexer/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func testingSetup() *RealtimeIndexer {
	util.ENOK(logger.SetLogLevel("error"))
	util.ENOK(config.LoadViperConfig("testdata/configs/testcfg.yaml"))
	return NewRealtimeIndexer(0,
		"https://rpc.ankr.com/eth",
		[]string{},
		time.Second,
		false,
		&db.DBConn{ChainID: 1},
		[]string{})
}

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

func loadJson(t *testing.T, file string, item interface{}) {
	jsonFile, err := os.Open(file)
	if err != nil {
		t.Error(err)
	}
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, item)
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

func assertBigFloatClose(t *testing.T, expected *big.Float, real *big.Float, threshold *big.Float) {
	if threshold == nil {
		threshold = big.NewFloat(0.0001)
	}
	diff := big.NewFloat(0).Sub(expected, real)
	cmp := big.NewFloat(0).Sub(threshold, diff.Abs(diff))
	withinThreshold := cmp.Cmp(big.NewFloat(0)) == 1
	assert.True(t, withinThreshold, "not close enough")
}
