package indexer

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/ERC20"
	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/univ2pair"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"
)

type DataAccess struct {
	upstreams *LatencySortedPool
	txLRU     *lru.Cache
}

type UniV2Reserves struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}

func NewDataAccess(upstreams []string) *DataAccess {
	txLRU, _ := lru.New(1024) // Hardcoded 1024
	lsp := NewLatencySortedPool(upstreams)
	go lsp.ShowStatus()
	return &DataAccess{
		upstreams: lsp,
		txLRU:     txLRU,
	}
}

func (d *DataAccess) Len() int {
	return d.upstreams.Len()
}

func (d *DataAccess) GetFilteredLogs(fq ethereum.FilterQuery) ([]types.Log, error) {
	var logs []types.Log
	var err error

	for retries := 0; retries < WD; retries++ {
		cl := d.upstreams.GetItem()

		start := time.Now()
		logs, err = cl.FilterLogs(context.Background(), fq)
		d.upstreams.Report(cl, time.Now().Sub(start).Seconds(), err != nil)

		if err == nil {
			return logs, nil
		}
	}

	return logs, errors.New("Fetch error: " + err.Error())
}

func (d *DataAccess) GetTokensUniV2(pairContract common.Address, callopts *bind.CallOpts) (common.Address, common.Address, error) {
	var token0, token1 common.Address
	var err error
	var pc *univ2pair.Univ2pair

	for retries := 0; retries < WD; retries++ {
		cl := d.upstreams.GetItem()
		pc, err = univ2pair.NewUniv2pair(pairContract, cl)

		start := time.Now()
		token0, err = pc.Token0(callopts)
		elapsed := time.Now().Sub(start).Seconds()
		if err != nil {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(cl, elapsed, false)
				return token0, token1, err
			}
			continue
		}
		d.upstreams.Report(cl, elapsed, err != nil)

		start = time.Now()
		token1, err = pc.Token1(callopts)
		d.upstreams.Report(cl, time.Now().Sub(start).Seconds(), err != nil)
		if err == nil {
			return token0, token1, nil
		}
	}

	return token0, token1, errors.New("Fetch error: " + err.Error())
}

func (d *DataAccess) GetDEXReserves(pairContract common.Address, token0 common.Address, token1 common.Address, callopts *bind.CallOpts) (UniV2Reserves, error) {
	var reserves UniV2Reserves
	var err error
	var token0Contract *ERC20.ERC20
	var token1Contract *ERC20.ERC20
	for retries := 0; retries < WD; retries++ {
		cl := d.upstreams.GetItem()

		start := time.Now()
		token0Contract, err = ERC20.NewERC20(token0, cl)
		token1Contract, err = ERC20.NewERC20(token1, cl)
		balToken0, err := token0Contract.BalanceOf(callopts, pairContract)
		balToken1, err := token1Contract.BalanceOf(callopts, pairContract)

		elapsed := time.Now().Sub(start).Seconds()
		if err == nil {
			d.upstreams.Report(cl, elapsed, false)
			return UniV2Reserves{
				Reserve0:           balToken0,
				Reserve1:           balToken1,
				BlockTimestampLast: reserves.BlockTimestampLast,
			}, nil
		}

		if err != nil {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(cl, elapsed, false)
				break
			}
			d.upstreams.Report(cl, elapsed, true)
		}
	}

	return reserves, errors.New("Fetch error: " + err.Error())
}

func (d *DataAccess) GetERC20Decimals(erc20Address common.Address, callopts *bind.CallOpts) (uint8, error) {
	var decimals uint8
	var err error
	var token *ERC20.ERC20

	for retries := 0; retries < WD; retries++ {
		cl := d.upstreams.GetItem()
		token, err = ERC20.NewERC20(erc20Address, cl)

		start := time.Now()
		decimals, err = token.Decimals(callopts)
		elapsed := time.Now().Sub(start).Seconds()
		if err == nil {
			d.upstreams.Report(cl, elapsed, false)
			return decimals, nil
		}
		if err != nil {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(cl, elapsed, false)
				break
			}
			d.upstreams.Report(cl, elapsed, true)
		}
	}

	return decimals, errors.New("Fetch error: " + err.Error())
}

func (d *DataAccess) GetTxSender(txHash common.Hash, blockHash common.Hash, txIdx uint) (common.Address, error) {
	var sender common.Address
	var err error
	var tx *types.Transaction

	for retries := 0; retries < WD; retries++ {
		cl := d.upstreams.GetItem()

		start := time.Now()
		tx, _, err = cl.TransactionByHash(context.Background(), txHash)
		elapsed := time.Now().Sub(start).Seconds()
		if err != nil {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(cl, elapsed, false)
				break
			}
			d.upstreams.Report(cl, elapsed, true)
		}

		start = time.Now()
		sender, err = cl.TransactionSender(context.Background(), tx, blockHash, txIdx)
		elapsed = time.Now().Sub(start).Seconds()
		if err == nil {
			d.upstreams.Report(cl, elapsed, false)
			return sender, nil
		}
		if err != nil {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(cl, elapsed, false)
				break
			}
			d.upstreams.Report(cl, elapsed, true)
		}
	}

	return sender, errors.New("Fetch error: " + err.Error())
}

func (d *DataAccess) GetCurrentBlockHeight() (uint64, error) {
	var height uint64
	var err error

	for retries := 0; retries < WD; retries++ {
		cl := d.upstreams.GetItem()

		start := time.Now()
		height, err = cl.BlockNumber(context.Background())
		elapsed := time.Now().Sub(start).Seconds()
		if err == nil {
			d.upstreams.Report(cl, elapsed, false)
			return height, nil
		}
		if err != nil {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(cl, elapsed, false)
				break
			}
			d.upstreams.Report(cl, elapsed, true)
		}
	}

	return height, errors.New("Fetch error: " + err.Error())
}
