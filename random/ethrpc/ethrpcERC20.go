package ethrpc

import (
	"context"
	"math/big"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/ERC20"
	"github.com/Blockpour/Blockpour-Geth-Indexer/mspool"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Cached RPC access to get decimals for ERC20 addresses
func (d *EthRPC) GetERC20Decimals(erc20Address common.Address, callopts *bind.CallOpts) (uint8, error) {
	lookupKey := erc20Address
	if ret, ok := d.ERC20Cache.Get(lookupKey); ok {
		retI := ret.(uint8)
		return retI, nil
	}

	decimals, err := mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) (uint8, error) {
			erc20, err := ERC20.NewERC20(erc20Address, c)
			if err != nil {
				return 0, nil
			}
			callopts.Context = ctx
			return erc20.Decimals(callopts)
		}, 0)

	if err != nil {
		return 0, err
	}

	d.ERC20Cache.Add(lookupKey, decimals)
	return decimals, nil
}

// Non-cached RPC access to get balances for tuple (holderAddress, tokenAddress)
func (d *EthRPC) GetERC20Balances(requests []itypes.Tuple2[common.Address, common.Address],
	callopts *bind.CallOpts) ([]itypes.Tuple2[common.Address, *big.Int], error) {
	results := []itypes.Tuple2[common.Address, *big.Int]{}

	for _, req := range requests {
		balance, err := mspool.Do(d.upstreams,
			func(ctx context.Context, c *ethclient.Client) (*big.Int, error) {
				token, err := ERC20.NewERC20(req.Second, c)
				if err != nil {
					return big.NewInt(0), nil
				}
				callopts.Context = ctx
				return token.BalanceOf(callopts, req.First)
			}, nil)
		if err != nil {
			return results, err
		}
		results = append(results, itypes.Tuple2[common.Address, *big.Int]{req.First, balance})
	}
	return results, nil
}
