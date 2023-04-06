package rpc

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/common"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/util"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"golang.org/x/time/rate"

	"github.com/patrickmn/go-cache"
)

type AlgoRPC struct {
	client  *algod.Client
	indexer *indexer.Client
	cache   *cache.Cache
	limiter *rate.Limiter
}

func NewAlgoRPC(algodUrl string, indexerUrl string, token string) (*AlgoRPC, error) {
	i, err := indexer.MakeClient(indexerUrl, token)
	if err != nil {
		fmt.Println("Could not instantiate indexer client")
		return nil, err
	}

	c, err := algod.MakeClient(algodUrl, token)
	if err != nil {
		fmt.Println("Could not instantiate client")
		return nil, err
	}

	return &AlgoRPC{
		indexer: i,
		client:  c,
		cache:   cache.New(5*time.Minute, 10*time.Minute),
		limiter: rate.NewLimiter(100, 1),
	}, err
}

func (r *AlgoRPC) GetCurrentBlockHeight() (uint64, error) {
	r.limiter.Wait(context.Background())
	status, err := r.client.Status().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return status.LastRound, nil
}

func (r *AlgoRPC) GetBlock(round uint64) (models.Block, error) {
	r.limiter.Wait(context.Background())

	block, err := r.indexer.LookupBlock(round).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return models.Block{}, err
	}

	return block, nil
}

func (r *AlgoRPC) GetBlockTransactions(round uint64) ([]models.Transaction, error) {
	r.limiter.Wait(context.Background())

	block, err := r.indexer.LookupBlock(round).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	txns := []models.Transaction{}
	for _, txn := range block.Transactions {
		txns = append(txns, models.Transaction(txn))
	}

	return txns, nil
}

func (r *AlgoRPC) GetBlockTimestamp(round uint64) (uint64, error) {
	r.limiter.Wait(context.Background())

	block, err := r.indexer.LookupBlock(round).Do(context.Background())
	if err != nil {
		return 0, err
	}

	return block.Timestamp, nil
}

func (r *AlgoRPC) GetAssetInfo(assetId uint64) (types.TokenInfo, error) {
	if assetId == 0 || assetId == 1 {
		return common.ALGO, nil
	}

	if x, found := r.cache.Get(fmt.Sprintf("asset_%d", assetId)); found {
		return x.(types.TokenInfo), nil
	}

	r.limiter.Wait(context.Background())

	asset, err := r.indexer.SearchForAssets().AssetID(assetId).Do(context.Background())
	if err != nil {
		return types.TokenInfo{}, err
	}

	name := asset.Assets[0].Params.UnitName
	decimals := asset.Assets[0].Params.Decimals

	info := types.TokenInfo{
		Id:       assetId,
		Name:     name,
		Decimals: decimals,
	}

	r.cache.Set(fmt.Sprintf("asset_%d", assetId), info, 24*time.Hour)

	return info, nil
}

func (r *AlgoRPC) FindPreviousContractTx(round uint64, appId uint64, assetId uint64, filter []types.FunctionSignature) ([]models.Transaction, error) {
	asset1, asset2, _ := r.GetAlgoFiPoolInfo(appId)

	var key = ""
	if assetId == asset1 {
		key = "a1r"
	} else if assetId == asset2 {
		key = "a2r"
	} else {
		return nil, fmt.Errorf("assetId not found in pool")
	}

	util.EncodeBase64(key)

	for i := round; i > 0; i-- {
		block, err := r.indexer.LookupBlock(i).Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, txn := range block.Transactions {
			if len(txn.ApplicationTransaction.ApplicationArgs) == 0 {
				continue
			}

			txKey := string(txn.ApplicationTransaction.ApplicationArgs[0])
			txAppId := txn.ApplicationTransaction.ApplicationId

			for _, funcsig := range filter {
				if funcsig.Equals(types.FunctionSignature{AppId: txAppId, Key: txKey}) {
					for _, delta := range txn.GlobalStateDelta {
						fmt.Println(delta.Key)
						if delta.Key == util.EncodeBase64(key) {
							fmt.Println("Found key", delta)
						}
					}
				}
			}
		}
	}

	return nil, nil
}

func (r *AlgoRPC) GetAccountAssets(address string) ([]models.AssetHolding, error) {
	r.limiter.Wait(context.Background())

	res, err := r.indexer.LookupAccountAssets(address).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res.Assets, nil
}

func (r *AlgoRPC) GetALGOBalanceWei(address string) (uint64, error) {
	r.limiter.Wait(context.Background())

	_, account, err := r.indexer.LookupAccountByID(address).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return account.Amount, nil
}

func (r *AlgoRPC) GetALGOBalance(address string) (float64, error) {
	r.limiter.Wait(context.Background())

	_, account, err := r.indexer.LookupAccountByID(address).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	amount := float64(account.Amount) / float64(1000000)

	return amount, nil
}

func (r *AlgoRPC) GetBalanceWei(address string, assetId uint64) (uint64, error) {
	if assetId == 1 {
		return r.GetALGOBalanceWei(address)
	}

	assets, err := r.GetAccountAssets(address)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	asset := util.Filter(assets, func(a models.AssetHolding) bool {
		return a.AssetId == assetId
	})

	if len(asset) == 0 {
		return 0, fmt.Errorf("asset not found")
	}

	if len(asset) > 1 {
		return 0, fmt.Errorf("more than one asset found (impossible)")
	}

	return asset[0].Amount, nil
}

func (r *AlgoRPC) GetBalance(address string, assetId uint64, decimals uint64) (float64, error) {
	if assetId == 1 {
		return r.GetALGOBalance(address)
	}

	assets, err := r.GetAccountAssets(address)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	asset := util.Filter(assets, func(a models.AssetHolding) bool {
		return a.AssetId == assetId
	})

	if len(asset) == 0 {
		return 0, fmt.Errorf("asset not found")
	}

	if len(asset) > 1 {
		return 0, fmt.Errorf("more than one asset found (impossible)")
	}

	return float64(asset[0].Amount) / math.Pow(10, float64(decimals)), nil
}

func (r *AlgoRPC) GetTinymanV1PoolInfo(poolId string) (uint64, uint64, error) {
	if x, found := r.cache.Get(fmt.Sprintf("poolinfo_%s", poolId)); found {
		return x.([]uint64)[0], x.([]uint64)[1], nil
	}

	r.limiter.Wait(context.Background())

	result, err := r.indexer.LookupAccountAppLocalStates(poolId).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}

	asset0State := util.Filter(result.AppsLocalStates[0].KeyValue, func(a models.TealKeyValue) bool {
		decoded, _ := util.DecodeBase64(a.Key)
		return decoded == "a1"
	})
	asset1State := util.Filter(result.AppsLocalStates[0].KeyValue, func(a models.TealKeyValue) bool {
		decoded, _ := util.DecodeBase64(a.Key)
		return decoded == "a2"
	})

	// Tinyman uses 0 for the native currency ALGO id. We use the number 1 for ALGO (same as algofi).
	// So we convert the 0 to 1.
	var asset0Id uint64
	var asset1Id uint64
	if asset0State[0].Value.Uint == 0 {
		asset0Id = 1
	} else {
		asset0Id = asset0State[0].Value.Uint
	}

	if asset1State[0].Value.Uint == 0 {
		asset1Id = 1
	} else {
		asset1Id = asset1State[0].Value.Uint
	}

	r.cache.Set(fmt.Sprintf("poolinfo_%s", poolId), []uint64{asset0Id, asset1Id}, 24*time.Hour)
	return asset0Id, asset1Id, nil
}

func (r *AlgoRPC) GetAlgoFiPoolInfo(appId uint64) (uint64, uint64, error) {
	if x, found := r.cache.Get(fmt.Sprintf("poolinfo_%d", appId)); found {
		return x.([]uint64)[0], x.([]uint64)[1], nil
	}

	r.limiter.Wait(context.Background())

	result, err := r.indexer.LookupApplicationByID(appId).Do(context.Background())
	if err != nil {
		return 0, 0, err
	}

	asset0State := util.Filter(result.Application.Params.GlobalState, func(a models.TealKeyValue) bool {
		decoded, _ := util.DecodeBase64(a.Key)
		return decoded == "a1"
	})

	asset1State := util.Filter(result.Application.Params.GlobalState, func(a models.TealKeyValue) bool {
		decoded, _ := util.DecodeBase64(a.Key)
		return decoded == "a2"
	})

	// Tinyman uses 0 for the native currency ALGO id. We use the number 1 for ALGO (same as algofi).
	// So we convert the 0 to 1.
	var asset0Id uint64
	var asset1Id uint64
	if asset0State[0].Value.Uint == 0 {
		asset0Id = 1
	} else {
		asset0Id = asset0State[0].Value.Uint
	}

	if asset1State[0].Value.Uint == 0 {
		asset1Id = 1
	} else {
		asset1Id = asset1State[0].Value.Uint
	}

	r.cache.Set(fmt.Sprintf("poolinfo_%d", appId), []uint64{asset0Id, asset1Id}, 24*time.Hour)

	return asset0Id, asset1Id, nil
}

func (r *AlgoRPC) GetAlgoFiPoolReserves(appId uint64) (uint64, uint64, error) {
	r.limiter.Wait(context.Background())

	result, err := r.indexer.LookupApplicationByID(appId).Do(context.Background())
	if err != nil {
		return 0, 0, err
	}

	asset0State := util.Filter(result.Application.Params.GlobalState, func(a models.TealKeyValue) bool {
		decoded, _ := util.DecodeBase64(a.Key)
		return decoded == "a1r"
	})

	asset1State := util.Filter(result.Application.Params.GlobalState, func(a models.TealKeyValue) bool {
		decoded, _ := util.DecodeBase64(a.Key)
		return decoded == "a2r"
	})

	asset0Id := asset0State[0].Value.Uint
	asset1Id := asset1State[0].Value.Uint

	return asset0Id, asset1Id, nil
}

func (r *AlgoRPC) GetTinymanV2PoolInfo(poolId string) (uint64, uint64, error) {
	if x, found := r.cache.Get(fmt.Sprintf("poolinfo_%s", poolId)); found {
		return x.([]uint64)[0], x.([]uint64)[1], nil
	}

	r.limiter.Wait(context.Background())

	result, err := r.indexer.LookupAccountAppLocalStates(poolId).Do(context.Background())
	if err != nil {
		return 0, 0, err
	}

	asset0State := util.Filter(result.AppsLocalStates[0].KeyValue, func(a models.TealKeyValue) bool {
		decoded, _ := util.DecodeBase64(a.Key)
		return decoded == "asset_1_id"
	})
	asset1State := util.Filter(result.AppsLocalStates[0].KeyValue, func(a models.TealKeyValue) bool {
		decoded, _ := util.DecodeBase64(a.Key)
		return decoded == "asset_2_id"
	})

	// Tinyman uses 0 for the native currency ALGO id. We use the number 1 for ALGO (same as algofi).
	// So we convert the 0 to 1.
	var asset0Id uint64
	var asset1Id uint64
	if asset0State[0].Value.Uint == 0 {
		asset0Id = 1
	} else {
		asset0Id = asset0State[0].Value.Uint
	}

	if asset1State[0].Value.Uint == 0 {
		asset1Id = 1
	} else {
		asset1Id = asset1State[0].Value.Uint
	}

	r.cache.Set(fmt.Sprintf("poolinfo_%s", poolId), []uint64{asset0Id, asset1Id}, 24*time.Hour)
	return asset0Id, asset1Id, nil
}

func (r *AlgoRPC) GetTinymanV1PoolReserves(txn models.Transaction) (uint64, uint64, error) {
	reserve0State := util.Filter(txn.LocalStateDelta[0].Delta, func(kv models.EvalDeltaKeyValue) bool {
		return kv.Key == "czI="
	})

	reserve1State := util.Filter(txn.LocalStateDelta[0].Delta, func(kv models.EvalDeltaKeyValue) bool {
		return kv.Key == "czE="
	})

	reserves0 := reserve0State[0].Value.Uint
	reserves1 := reserve1State[0].Value.Uint
	return reserves0, reserves1, nil
}
