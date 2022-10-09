package ethrpc

import (
	"context"
	"math/big"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/univ3pair"
	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/univ3positionsnft"
	"github.com/Blockpour/Blockpour-Geth-Indexer/mspool"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (d *EthRPC) GetTokensUniV3(pairContract common.Address,
	callopts *bind.CallOpts) (common.Address, common.Address, error) {
	lookupKey := util.Tuple2[common.Address, bind.CallOpts]{pairContract, *callopts}
	if ret, ok := d.contractTokensCache.Get(lookupKey); ok {
		retI := ret.(util.Tuple2[common.Address, common.Address])
		return retI.First, retI.Second, nil
	}

	token0, err := mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			pc, err := univ3pair.NewUniv3pair(pairContract, c)
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
			pc, err := univ3pair.NewUniv3pair(pairContract, c)
			if err != nil {
				return common.Address{}, err
			}
			callopts.Context = ctx
			return pc.Token1(callopts)
		}, common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	d.contractTokensCache.Add(lookupKey, util.Tuple2[common.Address, common.Address]{token0, token1})
	return token0, token1, nil
}

func (d *EthRPC) GetTokensUniV3NFT(nftContract common.Address, tokenID *big.Int, callopts *bind.CallOpts) (common.Address, common.Address, error) {
	// Cache checkup
	lookupKey := util.Tuple2[common.Address, bind.CallOpts]{nftContract, *callopts}
	if ret, ok := d.contractTokensCache.Get(lookupKey); ok {
		retI := ret.(util.Tuple2[common.Address, common.Address])
		return retI.First, retI.Second, nil
	}

	tokens, err := mspool.Do(d.upstreams,
		func(ctx context.Context, c *ethclient.Client) (util.Tuple2[common.Address, common.Address], error) {
			pc, err := univ3positionsnft.NewUniv3positionsnft(nftContract, c)
			if err != nil {
				return util.Tuple2[common.Address, common.Address]{}, err
			}
			callopts.Context = ctx
			positions, err := pc.Positions(callopts, tokenID)
			return util.Tuple2[common.Address, common.Address]{positions.Token0, positions.Token1}, err
		}, util.Tuple2[common.Address, common.Address]{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	d.contractTokensCache.Add(lookupKey, tokens)
	return tokens.First, tokens.Second, nil

}
