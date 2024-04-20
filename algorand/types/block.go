package types

import "github.com/supragya/EtherScope/algorand/util"

type BlockSynopsis struct {
	Type                    string
	Network                 uint
	Height                  uint64
	Time                    uint64
	TotalLogs               uint64
	SwapLogs                uint64
	BurnLogs                uint64
	MintLogs                uint64
	IndexingTimeNanos       uint64
	ProcessingDurationNanos uint64
}

func (b BlockSynopsis) Print() {
	util.PrintJSON(map[string]interface{}{
		"type":                    b.Type,
		"network":                 b.Network,
		"height":                  b.Height,
		"time":                    b.Time,
		"totalLogs":               b.TotalLogs,
		"swapLogs":                b.SwapLogs,
		"indexingTimeNanos":       b.IndexingTimeNanos,
		"processingdurationNanos": b.ProcessingDurationNanos,
	})
}
