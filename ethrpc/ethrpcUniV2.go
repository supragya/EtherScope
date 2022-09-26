package ethrpc

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/univ2pair"
	"github.com/Blockpour/Blockpour-Geth-Indexer/mspool"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (d *EthRPC) GetTokensUniV2(pairContract common.Address, callopts *bind.CallOpts) (common.Address, common.Address, error) {
	// Cache checkup
	lookupKey := util.Tuple2[common.Address, bind.CallOpts]{pairContract, *callopts}
	if ret, ok := d.contractTokensCache.Get(lookupKey); ok {
		retI := ret.(util.Tuple2[common.Address, common.Address])
		return retI.First, retI.Second, nil
	}

	token0, err := mspool.Do(d.upstreams,
		func(c *ethclient.Client) (common.Address, error) {
			pc, err := univ2pair.NewUniv2pair(pairContract, c)
			if err != nil {
				return common.Address{}, err
			}
			return pc.Token0(callopts)
		}, common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	token1, err := mspool.Do(d.upstreams,
		func(c *ethclient.Client) (common.Address, error) {
			pc, err := univ2pair.NewUniv2pair(pairContract, c)
			if err != nil {
				return common.Address{}, err
			}
			return pc.Token1(callopts)
		}, common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	d.contractTokensCache.Add(lookupKey, util.Tuple2[common.Address, common.Address]{token0, token1})
	return token0, token1, nil
}
