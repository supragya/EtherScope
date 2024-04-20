package node

import (
	"fmt"
	"strconv"
	"sync"

	types "github.com/supragya/EtherScope/algorand/types"
)

/*
I) There are 2 types of transaction signatures for Algofi.

1) Transaction without change (sef)
The basic case is one transfer from the sender to the pool and one from the pool
to the sender.

2) Transaction with change (sfe)
In this type of transaction, the sender wishes to get an exact amount of the BUY asset.
So the sender will send a transfer a certain amount of SELL asset to the pool
and the pool will return the amount of BUY asset that the sender wants. The pool will also
return the excess SELL asset to the sender.

# II) Considerations

Every single transaction is different. The basic case is a simple transfer of an
asset to the pool and the pool returns the other asset.

However transaction groups often include more than one swap including things like
- change transactions
- nanoswap transactions (algofi uses a different contract for stablecoin/stablecoin
swaps.)
- cross dex arb
- flashloans

III) Example of transaction groups
1) case where there are 2 transactions.
https://algoexplorer.io/tx/group/mEPxUp%2BeZqGn7sd88ssaUrKGvM%2FbNer1fUsMAwyRpWI%3D

2) case where there are 3 transactions (last tx is the change tx)
the last transaction is the fee transaction.
https://algoexplorer.io/tx/group/ZlVJHnC%2FNux8Ot4bNZB3oCd1Pdghf1uCSj4TIocZbR4%3D

3) case where there more transactions:
algofi + tinyman: kDHSaenH9Qam6yWgFQsBME/+zQH6UUH1dYey35RI3fw=
https://algoexplorer.io/tx/group/kDHSaenH9Qam6yWgFQsBME%2F%2BzQH6UUH1dYey35RI3fw%3D

4) arb between algofi and tinyman
Vu/H8FfYGAhp0PQHM+PXr8oJRG0gAwxCfOxzaTfGbOU=
https://algoexplorer.io/tx/group/Vu%2FH8FfYGAhp0PQHM%2BPXr8oJRG0gAwxCfOxzaTfGbOU%3D

5) ALGO -> STBL -> USDC -> ALGO. Arb.
https://algoexplorer.io/tx/group/2lMQI1PpvETKAasQHTtWLpipNh9ZOXf7bOpU1Yr%2B5BU%3D

6) Transaction groups including nanoswaps (stablecoin swaps)
https://algoexplorer.io/tx/group/ghBYk9Qg%2BIrFXpMLVt3T3PRSmAIRONQlgAuiLwlUVEw%3D
*/
func (n *Node) processAlgofiSwap(g types.TxGroup, index int, results map[int][]types.Swap, wg *sync.WaitGroup, mt *sync.Mutex) error {
	defer wg.Done()

	// groupId := string(g.GroupID)
	height := g.Transactions[0].ConfirmedRound
	time := g.Transactions[0].RoundTime
	txns := g.Transactions
	swaps := []types.Swap{}

	for k := 0; k < len(txns)-1; k++ {
		// TODO: replace this with checking the txn "function signature"
		// case of transactions without change
		if (txns[k].Type == "axfer" || txns[k].Type == "pay") &&
			(txns[k+1].Type == "appl" && (string(txns[k+1].ApplicationTransaction.ApplicationArgs[0]) == "sef")) {
			txId := txns[k+1].Id
			poolId := txns[k+1].ApplicationTransaction.ApplicationId
			txSender := txns[k].Sender

			asset0Id, asset1Id, err := n.rpc.GetAlgoFiPoolInfo(poolId)
			if err != nil {
				fmt.Println(err)
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

			sender, receiver, poolAddress, amount0, amount1, err := n.rpc.SumTransfers(txns[k:k+2], asset0Id, asset1Id)
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

			/*
				need to differentiate between price = 0 and no price so
				we need to return an empty price result struct in case we fail to
				determine the price, hence this weird logic
			*/
			priceResult0 := types.PriceResult{}
			priceResult1 := types.PriceResult{}

			if price0 != 0 {
				priceResult0.Price = price0
			}

			if price1 != 0 {
				priceResult1.Price = price1
			}

			swap := types.Swap{
				Type:         "algofiswap",
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
				AmountUSD:    amountUSD,
				// GroupId:      groupId,
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

			// case of transactions with change
			// TODO: replace this with checking the txn "function signature"
		} else if (txns[k].Type == "axfer" || txns[k].Type == "pay") &&
			(txns[k+1].Type == "appl" && (string(txns[k+1].ApplicationTransaction.ApplicationArgs[0]) == "sfe")) {
			txId := txns[k+1].Id
			poolId := txns[k+1].ApplicationTransaction.ApplicationId
			txSender := txns[k].Sender

			asset0Id, asset1Id, err := n.rpc.GetAlgoFiPoolInfo(poolId)
			if err != nil {
				fmt.Println(err)
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

			sender, receiver, poolAddress, amount0, amount1, err := n.rpc.SumTransfers(txns[k:k+3], asset0Id, asset1Id)
			if err != nil {
				fmt.Println(err)
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
				fmt.Println(err)
				return err
			}

			amountUSD, err := n.pricing.GetAmountUSD(asset0, asset1, amount0, amount1, price0, price1)
			if err != nil {
				fmt.Println(err)
				return err
			}

			swap := types.Swap{
				Type:         "algofiswap",
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
				AmountUSD:    amountUSD,
				// GroupId:      groupId,
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
