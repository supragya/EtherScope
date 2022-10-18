package ethrpc

import (
	"context"
	"math/big"

	"github.com/Blockpour/Blockpour-Geth-Indexer/mspool"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Non-cached RPC access to get sender address for any eth transaction
func (d *EthRPC) GetTxSender(txHash common.Hash,
	blockHash common.Hash,
	txIdx uint) (common.Address, error) {
	tx, err := mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) (*types.Transaction, error) {
			tx, _, err := c.TransactionByHash(ctx, txHash)
			return tx, err
		}, nil)
	if err != nil {
		return common.Address{}, err
	}

	sender, err := mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			return c.TransactionSender(ctx, tx, blockHash, txIdx)
		}, common.Address{})
	return sender, err
}

// Non-cached RPC access to get current block height
func (d *EthRPC) GetCurrentBlockHeight() (uint64, error) {
	return mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) (uint64, error) {
			return c.BlockNumber(ctx)
		}, 0)
}

// Non-cached RPC access to get block timestamp
func (d *EthRPC) GetBlockTimestamp(height uint64) (uint64, error) {
	header, err := mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) (*types.Header, error) {
			return c.HeaderByNumber(ctx, big.NewInt(int64(height)))
		}, nil)
	if err != nil {
		return 0, err
	}
	return header.Time, err
}

// Non-cached RPC access to get filtered logs
func (d *EthRPC) GetFilteredLogs(fq ethereum.FilterQuery) ([]types.Log, error) {
	return mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) ([]types.Log, error) {
			return c.FilterLogs(ctx, fq)
		}, []types.Log{})
}
