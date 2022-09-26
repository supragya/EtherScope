package mspool

import (
	"context"
	"errors"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
)

func Do[C any, T any](upstreams *MasterSlavePool[C],
	foo func(context.Context, *C) (T, error),
	defaultVal T) (T, error) {
	var gerr error = errors.New("")
	maxRetries := (len(upstreams.Slaves) + 1) * int(DefaultMSPoolConfig.WindowSize)
	for retries := 0; retries < maxRetries; retries++ {
		client := upstreams.GetItem()
		// log.Info(client)
		ctx := util.NewCtx(upstreams.RPCTimeout)
		out, _err := foo(ctx, client)
		if _err == nil {
			return out, nil
		}
		// In case of failure
		gerr = _err

		// Node failure
		if util.IsRPCCallTimedOut(_err) {
			// log.Info("hello", _err)
			// log.Info(fmt.Sprintf("%+v", upstreams.Master.Meta))
			upstreams.Report(client, true)
			time.Sleep(upstreams.config.TimeStep)
			continue
		}

		// Any other failure
		if util.IsEthErr(_err) {
			return defaultVal, _err
		}

		// Backoff and retry for any other case
		time.Sleep(upstreams.config.TimeStep)
	}
	return defaultVal, errors.New("Fetch error: " + gerr.Error())
}
