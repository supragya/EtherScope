package rpc

import (
	"context"
	"fmt"
	"math"

	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/common"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/util"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type Transfer struct {
	From   string
	To     string
	Amount float64
	Asset  types.TokenInfo
}

// Return all transactions in a block if they match the filter.
func (r *AlgoRPC) GetTxGroups(round uint64, funcsigs []types.FunctionSignature) ([]types.TxGroup, error) {
	groups := []types.TxGroup{}

	groupIds := []string{}

	block, err := r.indexer.LookupBlock(round).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return []types.TxGroup{}, nil
	}

	for _, txn := range block.Transactions {
		if len(txn.ApplicationTransaction.ApplicationArgs) == 0 {
			continue
		}

		if util.Contains(groupIds, string(txn.Group)) {
			continue
		}

		txKey := string(txn.ApplicationTransaction.ApplicationArgs[0])
		txAppId := txn.ApplicationTransaction.ApplicationId

		for _, funcsig := range funcsigs {
			if funcsig.Equals(types.FunctionSignature{AppId: txAppId, Key: txKey}) {
				txns := util.Filter(block.Transactions, func(t models.Transaction) bool {
					return string(t.Group) == string(txn.Group)
				})

				batch := []models.Transaction{}

				for _, txn := range txns {
					batch = append(batch, models.Transaction(txn))
				}

				txGroup := types.TxGroup{
					GroupID:           txn.Group,
					Transactions:      batch,
					FunctionSignature: funcsig,
				}

				groups = append(groups, txGroup)
				groupIds = append(groupIds, string(txn.Group))
			}
		}

	}

	return groups, nil
}

func (r *AlgoRPC) ParseTxn(txn models.Transaction) (types.TokenInfo, float64, string, error) {
	var asset types.TokenInfo
	var err error

	if txn.Type == "axfer" {
		assetId := txn.AssetTransferTransaction.AssetId
		asset, err = r.GetAssetInfo(assetId)
		if err != nil {
			return types.TokenInfo{}, 0, "", err
		}
	} else if txn.Type == "pay" {
		asset = common.ALGO
	} else if txn.Type == "appl" {
		for _, innerTx := range txn.InnerTxns {
			if innerTx.Type == "axfer" {
				assetId := innerTx.AssetTransferTransaction.AssetId
				asset, err = r.GetAssetInfo(assetId)
				if err != nil {
					return types.TokenInfo{}, 0, "", err
				}
			} else if innerTx.Type == "pay" {
				asset = common.ALGO
			}
		}
	}

	var receiver string
	if txn.Type == "axfer" {
		receiver = txn.AssetTransferTransaction.Receiver
	} else if txn.Type == "pay" {
		receiver = txn.PaymentTransaction.Receiver
	} else if txn.Type == "appl" {
		if txn.InnerTxns[0].Type == "axfer" {
			receiver = txn.InnerTxns[0].AssetTransferTransaction.Receiver
		} else if txn.InnerTxns[0].Type == "pay" {
			receiver = txn.InnerTxns[0].PaymentTransaction.Receiver
		}
	}

	var amountWei uint64
	if txn.Type == "axfer" {
		amountWei = txn.AssetTransferTransaction.Amount
	} else if txn.Type == "pay" {
		amountWei = txn.PaymentTransaction.Amount
	} else if txn.Type == "appl" {
		innerTx := txn.InnerTxns[0]
		if innerTx.Type == "axfer" {
			amountWei = innerTx.AssetTransferTransaction.Amount
		} else if innerTx.Type == "pay" {
			amountWei = innerTx.PaymentTransaction.Amount
		}
	}

	amount := float64(amountWei) / math.Pow(10, float64(asset.Decimals))
	return asset, amount, receiver, nil
}

func (r *AlgoRPC) ParseTransfers(txn models.Transaction) ([]Transfer, error) {
	var asset types.TokenInfo
	var err error

	if txn.Type == "axfer" {
		assetId := txn.AssetTransferTransaction.AssetId
		amountWei := txn.AssetTransferTransaction.Amount
		receiver := txn.AssetTransferTransaction.Receiver
		asset, err = r.GetAssetInfo(assetId)
		if err != nil {
			return []Transfer{}, err
		}

		transfer := Transfer{
			From:   txn.Sender,
			To:     receiver,
			Asset:  asset,
			Amount: float64(amountWei) / math.Pow(10, float64(asset.Decimals)),
		}

		return []Transfer{transfer}, nil

	} else if txn.Type == "pay" {
		asset = common.ALGO
		amountWei := txn.PaymentTransaction.Amount
		receiver := txn.PaymentTransaction.Receiver
		transfer := Transfer{
			From:   txn.Sender,
			To:     receiver,
			Asset:  asset,
			Amount: float64(amountWei) / math.Pow(10, float64(asset.Decimals)),
		}

		return []Transfer{transfer}, nil

	} else if txn.Type == "appl" {
		transfers := []Transfer{}

		for _, innerTx := range txn.InnerTxns {
			if innerTx.Type == "axfer" {
				assetId := innerTx.AssetTransferTransaction.AssetId
				amountWei := innerTx.AssetTransferTransaction.Amount
				receiver := innerTx.AssetTransferTransaction.Receiver
				asset, err = r.GetAssetInfo(assetId)

				if err != nil {
					return []Transfer{}, err
				}

				transfer := Transfer{
					From:   innerTx.Sender,
					To:     receiver,
					Asset:  asset,
					Amount: float64(amountWei) / math.Pow(10, float64(asset.Decimals)),
				}

				transfers = append(transfers, transfer)
			} else if innerTx.Type == "pay" {
				asset = common.ALGO
				amountWei := innerTx.PaymentTransaction.Amount
				receiver := innerTx.PaymentTransaction.Receiver

				transfer := Transfer{
					From:   innerTx.Sender,
					To:     receiver,
					Asset:  asset,
					Amount: float64(amountWei) / math.Pow(10, float64(asset.Decimals)),
				}

				transfers = append(transfers, transfer)
			}
		}

		return transfers, nil
	}

	return []Transfer{}, nil
}

func (r *AlgoRPC) SumTransfers(txns []models.Transaction, asset0Id uint64, asset1Id uint64) (string, string, string, float64, float64, error) {
	var assetIn = types.TokenInfo{}
	var amount0 float64
	var amount1 float64
	var receiver string
	sender := txns[0].Sender

	// The first transfer is always the one that sends the assets to the pool
	firstTransfers, err := r.ParseTransfers(txns[0])
	if err != nil {
		return "", "", "", 0, 0, err
	}

	assetIn = firstTransfers[0].Asset
	amountIn := firstTransfers[0].Amount
	poolAddress := firstTransfers[0].To

	transfers := []Transfer{}
	for _, txn := range txns[1:] {
		tt, err := r.ParseTransfers(txn)
		if err != nil {
			return "", "", "", 0, 0, err
		}

		transfers = append(transfers, tt...)
	}

	/*
		The sender is always defined as the tx sender of the first transaction.
		But a transaction group can potentially have several receivers.
		We sum the total amount of asset 0 and asset 1 received by different addresses.
		The receiver is defined as the address that received the most amount of assetOut.
		(assetOut can be either asset0 or asset1 which is why we need an if statement)
	*/
	amountDelta0 := map[string]float64{}
	amountDelta1 := map[string]float64{}

	if asset0Id == assetIn.Id {
		amountDelta0[sender] = -amountIn
		amountDelta1[sender] = 0
	} else {
		amountDelta0[sender] = 0
		amountDelta1[sender] = -amountIn
	}

	for _, t := range transfers {
		if t.Asset.Id == asset0Id {
			if t.To == poolAddress {
				amountDelta0[t.From] -= t.Amount
			} else if t.From == poolAddress {
				amountDelta0[t.To] += t.Amount
			}
		} else if t.Asset.Id == asset1Id {
			if t.To == poolAddress {
				amountDelta1[t.From] -= t.Amount
			} else if t.From == poolAddress {
				amountDelta1[t.To] += t.Amount
			}
		}
	}

	// The receiver is defined as the address that received the most amount of assetOut.
	if asset0Id == assetIn.Id {
		amount0 = amountDelta0[sender]
		receiver, amount1 = util.MapMax(amountDelta1)
	} else {
		amount1 = amountDelta1[sender]
		receiver, amount0 = util.MapMax(amountDelta0)
	}

	return sender, receiver, poolAddress, amount0, amount1, nil
}
