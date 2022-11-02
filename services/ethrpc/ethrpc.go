package ethrpc

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthRPC interface {
	service.Service

	GetTxSender(txHash, blockHash common.Hash, txIdx uint) (common.Address, error)
	GetCurrentBlockHeight() (uint64, error)
	GetBlockTimestamp(height uint64) (uint64, error)
	GetFilteredLogs(ethereum.FilterQuery) ([]types.Log, error)
}
