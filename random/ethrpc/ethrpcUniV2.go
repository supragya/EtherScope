package ethrpc

import (
	"context"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/univ2pair"
	"github.com/Blockpour/Blockpour-Geth-Indexer/mspool"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (d *EthRPC) GetTokensUniV2(pairContract common.Address, callopts *bind.CallOpts) (common.Address, common.Address, error) {
	// Cache checkup
	lookupKey := itypes.Tuple2[common.Address, bind.CallOpts]{pairContract, *callopts}
	if ret, ok := d.contractTokensCache.Get(lookupKey); ok {
		retI := ret.(itypes.Tuple2[common.Address, common.Address])
		return retI.First, retI.Second, nil
	}

	token0, err := mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			pc, err := univ2pair.NewUniv2pair(pairContract, c)
			if err != nil {
				return common.Address{}, err
			}
			callopts.Context = ctx
			return pc.Token0(callopts)
		}, common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	token1, err := mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			pc, err := univ2pair.NewUniv2pair(pairContract, c)
			if err != nil {
				return common.Address{}, err
			}
			callopts.Context = ctx
			return pc.Token1(callopts)
		}, common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	d.contractTokensCache.Add(lookupKey, itypes.Tuple2[common.Address, common.Address]{token0, token1})
	return token0, token1, nil
}
