package mspool

import (
	"errors"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
)

func Do[C any, T any](upstreams *MasterSlavePool[C],
	foo func(*C) (T, error),
	defaultVal T) (T, error) {
	var gerr error = errors.New("")
	maxRetries := (len(upstreams.Slaves) + 1) * int(DefaultMSPoolConfig.WindowSize)
	for retries := 0; retries < maxRetries; retries++ {
		client := upstreams.GetItem()
		out, _err := foo(client)
		if _err == nil {
			return out, nil
		}
		// In case of failure
		gerr = _err
		if !util.IsEthErr(_err) {
			upstreams.Report(client, true)
		}
		// Backoff
		time.Sleep(DefaultMSPoolConfig.TimeStep + 100)
	}
	return defaultVal, errors.New("Fetch error: " + gerr.Error())
}
