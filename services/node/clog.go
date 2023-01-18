package node

import (
	"sort"

	"github.com/ethereum/go-ethereum/core/types"
)

type CLogType []types.Log

func (l CLogType) Len() int {
	return len(l)
}

func (l CLogType) Less(i, j int) bool {
	if l[i].BlockNumber < l[j].BlockNumber {
		return true
	}
	return l[i].Index < l[j].Index
}

func (l CLogType) Swap(i, j int) {
	temp := l[i]
	l[i] = l[j]
	l[i] = temp
}

type CUint64 []uint64

func (l CUint64) Len() int {
	return len(l)
}

func (l CUint64) Less(i, j int) bool {
	return l[i] < l[j]
}

func (l CUint64) Swap(i, j int) {
	temp := l[i]
	l[j] = l[i]
	l[i] = temp
}

func GroupByBlockNumber(logs []types.Log) map[uint64]CLogType {
	kv := make(map[uint64]CLogType)
	for _, log := range logs {
		block := log.BlockNumber
		if val, ok := kv[block]; !ok {
			kv[block] = []types.Log{log}
		} else {
			kv[block] = append(val, log)
		}
	}
	for k, v := range kv {
		sort.Sort(v)
		kv[k] = v
	}
	return kv
}
