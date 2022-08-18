package indexer

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/ERC20"
	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/univ2pair"
	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/univ3positionsnft"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (d *DataAccess) GetTokensUniV3(pairContract common.Address, tokenID *big.Int, callopts *bind.CallOpts) (common.Address, common.Address, error) {
	type Positions struct {
		Nonce                    *big.Int
		Operator                 common.Address
		Token0                   common.Address
		Token1                   common.Address
		Fee                      *big.Int
		TickLower                *big.Int
		TickUpper                *big.Int
		Liquidity                *big.Int
		FeeGrowthInside0LastX128 *big.Int
		FeeGrowthInside1LastX128 *big.Int
		TokensOwed0              *big.Int
		TokensOwed1              *big.Int
	}
	var positions Positions
	var err error
	var pc *univ3positionsnft.Univ3positionsnftCaller

	for retries := 0; retries < WD; retries++ {
		cl := d.upstreams.GetItem()
		pc, err = univ3positionsnft.NewUniv3positionsnftCaller(pairContract, cl)

		start := time.Now()
		positions, err = pc.Positions(callopts, tokenID)
		elapsed := time.Now().Sub(start).Seconds()
		if err != nil {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(cl, elapsed, false)
				return common.Address{}, common.Address{}, err
			}
			continue
		}
		d.upstreams.Report(cl, elapsed, err != nil)
		break
	}

	return positions.Token0, positions.Token1, errors.New("Fetch error: " + err.Error())
}

func (d *DataAccess) GetERC20(erc20Address common.Address) (*ERC20.ERC20, *ethclient.Client) {
	cl := d.upstreams.GetItem()
	obj, err := ERC20.NewERC20(erc20Address, cl)
	util.ENOK(err)
	return obj, cl
}

func (d *DataAccess) GetERC20Decimals(erc20 *ERC20.ERC20, client *ethclient.Client, callopts *bind.CallOpts) (uint8, error) {
	var decimals uint8
	var err error

	for retries := 0; retries < WD; retries++ {
		start := time.Now()
		decimals, err = erc20.Decimals(callopts)
		elapsed := time.Now().Sub(start).Seconds()
		if err == nil {
			d.upstreams.Report(client, elapsed, false)
			return decimals, nil
		}
		if err != nil {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(client, elapsed, false)
				break
			}
			d.upstreams.Report(client, elapsed, true)
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

func (d *DataAccess) GetBlockTimestamp(height uint64) (uint64, error) {
	var err error

	for retries := 0; retries < WD; retries++ {
		cl := d.upstreams.GetItem()

		start := time.Now()
		bl, err := cl.BlockByNumber(context.Background(), big.NewInt(int64(height)))
		elapsed := time.Now().Sub(start).Seconds()
		if err == nil {
			d.upstreams.Report(cl, elapsed, false)
			return bl.Header().Time, nil
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

	return 0, errors.New("Fetch error: " + err.Error())
}

// Take it to utils
type Tuple2[A any, B any] struct {
	First  A
	Second B
}

func (d *DataAccess) GetBalances(requests []Tuple2[common.Address, common.Address],
	callopts *bind.CallOpts) ([]Tuple2[common.Address, *big.Int], error) {
	results := []Tuple2[common.Address, *big.Int]{}

	for _, req := range requests {
		balance, err := d.GetBalance(req.First, req.Second, callopts)
		if err != nil {
			return results, err
		}
		results = append(results, Tuple2[common.Address, *big.Int]{req.First, balance})
	}
	return results, nil
}

func (d *DataAccess) GetBalance(address common.Address,
	tokenAddress common.Address,
	callopts *bind.CallOpts) (*big.Int, error) {
	var balToken *big.Int
	var err error
	var token *ERC20.ERC20

	for retries := 0; retries < WD; retries++ {
		// Get Balance
		client := d.upstreams.GetItem()
		token, err = ERC20.NewERC20(tokenAddress, client)
		if err != nil {
			return big.NewInt(0), err
		}

		start := time.Now()
		balToken, err = token.BalanceOf(callopts, address)
		elapsed := time.Since(start).Seconds()
		if err == nil {
			d.upstreams.Report(client, elapsed, false)
			return balToken, nil
		} else {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(client, elapsed, false)
				break
			}
			d.upstreams.Report(client, elapsed, true)
		}
	}
	return balToken, errors.New("Fetch error: " + err.Error())
}

func (d *DataAccess) GetDEXReserves(
	pairContract common.Address,
	token0 *ERC20.ERC20,
	client0 *ethclient.Client,
	token1 *ERC20.ERC20,
	client1 *ethclient.Client,
	callopts *bind.CallOpts) (UniV2Reserves, error) {
	reserves := UniV2Reserves{}
	var err error
	for retries := 0; retries < WD; retries++ {
		// Get Balance 0
		start := time.Now()
		balToken0, err := token0.BalanceOf(callopts, pairContract)
		elapsed := time.Now().Sub(start).Seconds()
		if err == nil {
			d.upstreams.Report(client0, elapsed, false)
		} else {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(client0, elapsed, false)
				break
			}
			d.upstreams.Report(client0, elapsed, true)
		}

		// Get Balance 1
		start = time.Now()
		balToken1, err := token1.BalanceOf(callopts, pairContract)
		elapsed = time.Now().Sub(start).Seconds()
		if err != nil {
			// Early exit
			if util.IsEthErr(err) {
				d.upstreams.Report(client1, elapsed, false)
				break
			}
			d.upstreams.Report(client1, elapsed, true)
		}
		d.upstreams.Report(client1, elapsed, false)

		return UniV2Reserves{
			Reserve0:           balToken0,
			Reserve1:           balToken1,
			BlockTimestampLast: reserves.BlockTimestampLast,
		}, nil
	}

	return reserves, errors.New("Fetch error: " + err.Error())
}
