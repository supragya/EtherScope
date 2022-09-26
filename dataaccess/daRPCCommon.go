package dataaccess

import (
	"context"
	"math/big"

	"github.com/Blockpour/Blockpour-Geth-Indexer/mspool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Non-cached RPC access to get sender address for any eth transaction
func (d *DataAccess) GetTxSender(txHash common.Hash,
	blockHash common.Hash,
	txIdx uint) (common.Address, error) {
	tx, err := mspool.Do(d.upstreams,
		func(c *ethclient.Client) (*types.Transaction, error) {
			tx, _, err := c.TransactionByHash(context.Background(), txHash)
			return tx, err
		}, nil)
	if err != nil {
		return common.Address{}, err
	}
	sender, err := mspool.Do(d.upstreams,
		func(c *ethclient.Client) (common.Address, error) {
			return c.TransactionSender(context.Background(), tx, blockHash, txIdx)
		}, common.Address{})
	return sender, err
}

// Non-cached RPC access to get current block height
func (d *DataAccess) GetCurrentBlockHeight() (uint64, error) {
	return mspool.Do(d.upstreams,
		func(c *ethclient.Client) (uint64, error) {
			return c.BlockNumber(context.Background())
		}, 0)
}

// Non-cached RPC access to get block timestamp
func (d *DataAccess) GetBlockTimestamp(height uint64) (uint64, error) {
	header, err := mspool.Do(d.upstreams,
		func(c *ethclient.Client) (*types.Header, error) {
			return c.HeaderByNumber(context.Background(), big.NewInt(int64(height)))
		}, nil)
	return header.Time, err
}
