package pricing

import (
	"errors"
	"math"

	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/common"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/rpc"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/types"
)

type PricingEngine struct {
	rpc *rpc.AlgoRPC
}

func NewPricingEngine(rpc *rpc.AlgoRPC) *PricingEngine {
	return &PricingEngine{
		rpc: rpc,
	}
}

// works well for current prices but not for prices in the past as we cannot get reserves at a given block.
func (p *PricingEngine) GetPrices(poolAddress string, asset0 types.TokenInfo, asset1 types.TokenInfo, reserve0 float64, reserve1 float64) (float64, float64, error) {
	_, asset0IsQuote := common.Quotes[asset0.Id]
	_, asset1IsQuote := common.Quotes[asset1.Id]

	if asset0IsQuote && asset1IsQuote {
		price0, err := p.GetQuotePrice(asset0)
		if err != nil {
			return 0, 0, err
		}
		price1, err := p.GetQuotePrice(asset1)
		if err != nil {
			return 0, 0, err
		}

		return price0, price1, nil
	} else if asset0IsQuote && !asset1IsQuote {
		price0, err := p.GetQuotePrice(asset0)
		if err != nil {
			return 0, 0, err
		}

		if reserve1 == 0 || reserve0 == 0 {
			return 0, 0, errors.New("pool reserves are zero")
		}

		poolPrice1 := reserve0 / reserve1
		price1 := poolPrice1 * price0
		return price0, price1, nil
	} else if !asset0IsQuote && asset1IsQuote {
		price1, err := p.GetQuotePrice(asset1)
		if err != nil {
			return 0, 0, err
		}

		if reserve1 == 0 || reserve0 == 0 {
			return 0, 0, errors.New("pool reserves are zero")
		}

		poolPrice0 := reserve1 / reserve0
		price0 := poolPrice0 * price1
		return price0, price1, nil
	} else {
		return 0, 0, nil
	}
}

func (p *PricingEngine) GetAmountUSD(
	asset0 types.TokenInfo,
	asset1 types.TokenInfo,
	price0 float64,
	price1 float64,
	amount0 float64,
	amount1 float64,
) (float64, error) {
	_, asset0IsStablecoin := common.Stablecoins[asset0.Id]
	_, asset1IsStablecoin := common.Stablecoins[asset1.Id]
	_, asset0IsQuote := common.Quotes[asset0.Id]
	_, asset1IsQuote := common.Quotes[asset1.Id]

	// stablecoins provide the most precise pricing
	if asset0IsStablecoin && asset1IsStablecoin {
		return (math.Abs(amount0*price0) + math.Abs(amount1*price1)) / 2, nil
	} else if asset0IsStablecoin {
		return math.Abs(amount0 * price0), nil
	} else if asset1IsStablecoin {
		return math.Abs(amount1 * price1), nil
		// if none of the tokens is a stablecoin, we can use the USD price
		// provided by a pool with a large liquidity. For example ALGO/USDC
	} else if asset0IsQuote && asset1IsQuote {
		return (math.Abs(amount0*price0) + math.Abs(amount1*price1)) / 2, nil
	} else if asset0IsQuote && !asset1IsQuote {
		return math.Abs(amount0 * price0), nil
	} else if !asset0IsQuote && asset1IsQuote {
		return math.Abs(amount1 * price1), nil
		// if no tokens is either a stablecoin or a quote (e.g. token with a large /USDC pool)
		// we set the amountUSD to 0 because there is no easy/precise way to price that token.
	} else {
		return 0, nil
	}

}

func (p *PricingEngine) GetQuotePrice(asset types.TokenInfo) (float64, error) {
	switch asset.Id {
	case common.USDC.Id:
		return 1, nil
	case common.USDt.Id:
		return 1, nil
	case common.ALGO.Id:
		reserveALGO, _ := p.rpc.GetBalance(common.ALGOUSDC, common.ALGO.Id, common.ALGO.Decimals)
		reserveUSDC, _ := p.rpc.GetBalance(common.ALGOUSDC, common.USDC.Id, common.USDC.Decimals)
		if reserveALGO == 0 || reserveUSDC == 0 {
			return 0, errors.New("pool reserves are zero")
		}

		return reserveUSDC / reserveALGO, nil

	// the logic is a bit different since there doesn't seem to a pool address.
	// I think algofi keeps asset in different addresses as collateral for lending pool
	// (and used at the same time for swapping)
	case common.STBL2.Id:
		reserveSTBL2Wei, reserveUSDCWei, err := p.rpc.GetAlgoFiPoolReserves(common.STBL2USDC)
		if err != nil {
			return 0, err
		}

		if reserveSTBL2Wei == 0 || reserveUSDCWei == 0 {
			return 0, errors.New("pool reserves are zero")
		}

		reserveSTBL2 := float64(reserveSTBL2Wei) / math.Pow(10, float64(common.STBL2.Decimals))
		reserveUSDC := float64(reserveUSDCWei) / math.Pow(10, float64(common.USDC.Decimals))
		return float64(reserveUSDC) / float64(reserveSTBL2), nil
	}

	return 0, errors.New("quote not found")
}
