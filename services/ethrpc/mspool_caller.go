package ethrpc

import (
	"context"
	"errors"
	"time"

	"github.com/supragya/EtherScope/libs/util"
	"golang.org/x/sync/semaphore"
)

func Do[C any, T any](upstreams *MasterSlavePool[C],
	sem *semaphore.Weighted,
	foo func(context.Context, *C) (T, error),
	defaultVal T) (T, error) {

	util.ENOK(sem.Acquire(context.Background(), 1))
	defer sem.Release(1)

	var gerr error = errors.New("")
	maxRetries := (len(upstreams.Slaves) + 1) * int(DefaultMSPoolConfig.WindowSize)

	for retries := 0; retries < maxRetries; retries++ {
		client, _ := upstreams.GetItem()
		// log.Info(client)
		ctx := util.NewCtx(upstreams.RPCTimeout)
		// d, _ := ctx.Deadline()
		// log.Info("client found: ", meta.Identity, " with ", meta.Reports, " time rem context: ", d.Sub(time.Now()))
		out, _err := foo(ctx, client)
		if _err == nil {
			return out, nil
		}
		// In case of failure
		gerr = _err

		// Node failure
		if util.IsRPCCallTimedOut(_err) {
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
