package node

import (
	"fmt"
	"strconv"
	"sync"

	types "github.com/Blockpour/Blockpour-Geth-Indexer/algorand/types"
)

func (n *Node) processTinymanV1Swap(g types.TxGroup, index int, results map[int][]types.Swap, wg *sync.WaitGroup, mt *sync.Mutex) error {
	defer wg.Done()
	// tinyman V1 divides swaps into 4 different transactions
	// groupId := string(g.GroupID)
	height := g.Transactions[0].ConfirmedRound
	time := g.Transactions[0].RoundTime
	txns := g.Transactions
	swaps := []types.Swap{}

	for k := 0; k < len(txns); k++ {
		// TODO: replace this with checking the txn "function signature"
		// TODO: verify that the sender is always the tx sender
		if txns[k].Type == "appl" && string(txns[k].ApplicationTransaction.ApplicationArgs[0]) == "swap" && txns[k].ApplicationTransaction.ApplicationId == 552635992 {
			// feeTx := txns[k-1]
			txId := txns[k].Id
			poolId := txns[k].ApplicationTransaction.ApplicationId
			txSender := txns[k+1].Sender

			var asset0 = types.TokenInfo{}
			var asset1 = types.TokenInfo{}
			var amount0 float64
			var amount1 float64

			assetIn, amountIn, poolAddress, err := n.rpc.ParseTxn(txns[k+1])
			if err != nil {
				return err
			}

			assetOut, amountOut, receiver, err := n.rpc.ParseTxn(txns[k+2])
			if err != nil {
				return err
			}

			asset0Id, asset1Id, err := n.rpc.GetTinymanV1PoolInfo(poolAddress)
			if err != nil {
				return err
			}

			if asset0Id == assetIn.Id && asset1Id == assetOut.Id {
				asset0 = assetIn
				asset1 = assetOut
				amount0 = -amountIn
				amount1 = amountOut
			} else {
				asset0 = assetOut
				asset1 = assetIn
				amount0 = amountOut
				amount1 = -amountIn
			}

			reserves0, err := n.rpc.GetBalance(poolAddress, asset0.Id, asset0.Decimals)
			if err != nil {
				return err
			}
			reserves1, err := n.rpc.GetBalance(poolAddress, asset1.Id, asset1.Decimals)
			if err != nil {
				return err
			}

			price0, price1, err := n.pricing.GetPrices(poolAddress, asset0, asset1, reserves0, reserves1)
			if err != nil {
				return err
			}

			amountUSD, err := n.pricing.GetAmountUSD(asset0, asset1, amount0, amount1, price0, price1)
			if err != nil {
				fmt.Println(err)
				return err
			}

			swap := types.Swap{
				Type:         "tinymanv1swap",
				Network:      99990,
				Transaction:  txId,
				Height:       height,
				Time:         time,
				Sender:       txSender,
				Receiver:     receiver,
				Amount0:      amount0,
				Amount1:      amount1,
				PairContract: poolAddress,
				Token0:       strconv.FormatUint(asset0.Id, 10),
				Token1:       strconv.FormatUint(asset1.Id, 10),
				Name0:        asset0.Name,
				Name1:        asset1.Name,
				Decimals0:    asset0.Decimals,
				Decimals1:    asset1.Decimals,
				// GroupId:      groupId,
				AmountUSD: amountUSD,
				TxSender:  txSender,
				PoolAppId: poolId,
			}

			/*
				need to differentiate between price = 0 and no price so
				we need to return an empty price result struct in case we fail to
				determine the price, hence this weird logic
			*/
			if price0 != 0 {
				swap.Price0 = &types.PriceResult{Price: price0}
			}

			if price1 != 0 {
				swap.Price1 = &types.PriceResult{Price: price1}
			}

			swaps = append(swaps, swap)
		}

	}

	mt.Lock()
	results[index] = swaps
	mt.Unlock()
	return nil
}
