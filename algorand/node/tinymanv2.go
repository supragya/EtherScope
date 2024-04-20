package node

import (
	"strconv"
	"sync"

	types "github.com/supragya/EtherScope/algorand/types"
)

func (n *Node) processTinymanV2Swap(g types.TxGroup, index int, results map[int][]types.Swap, wg *sync.WaitGroup, mt *sync.Mutex) error {
	defer wg.Done()
	txns := g.Transactions
	// groupId := string(g.GroupID)
	height := g.Transactions[0].ConfirmedRound
	time := g.Transactions[0].RoundTime
	swaps := []types.Swap{}

	for k := 0; k < len(txns); k++ {
		if txns[k].Type == "appl" && string(txns[k].ApplicationTransaction.ApplicationArgs[0]) == "swap" && txns[k].ApplicationTransaction.ApplicationId == 1002541853 {
			txId := txns[k].Id
			poolId := txns[k].ApplicationTransaction.ApplicationId
			txSender := txns[k].Sender

			// we assume the first transaction contains only one transfer
			transfers, err := n.rpc.ParseTransfers(txns[k-1])
			if err != nil {
				return err
			}

			poolAddress := transfers[0].To

			asset0Id, asset1Id, err := n.rpc.GetTinymanV2PoolInfo(poolAddress)
			if err != nil {
				return err
			}

			asset0, err := n.rpc.GetAssetInfo(asset0Id)
			if err != nil {
				return err
			}

			asset1, err := n.rpc.GetAssetInfo(asset1Id)
			if err != nil {
				return err
			}

			sender, receiver, poolAddress, amount0, amount1, err := n.rpc.SumTransfers(txns[k-1:k+1], asset0Id, asset1Id)
			if err != nil {
				return err
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
				return err
			}

			swap := types.Swap{
				Type:         "tinymanv2swap",
				Network:      99990,
				Transaction:  txId,
				Height:       height,
				Time:         time,
				Sender:       sender,
				Receiver:     receiver,
				Amount0:      amount0,
				Amount1:      amount1,
				Reserve0:     reserves0,
				Reserve1:     reserves1,
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
