package indexer

import (
	"math/big"
	"strconv"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/ChainLink"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (d *DataAccess) GetPricesForBlock(
	callopts *bind.CallOpts, token0 common.Address, token1 common.Address, amount0 *big.Float, amount1 *big.Float) (float64, float64, float64) {

	oracleMap, err := util.GetOracleContracts(viper.GetUint("general.chainid"))

	util.ENOK(err)

	cl := d.upstreams.GetItem()
	var token0Price float64
	var token1Price float64
	var amountusd float64

	if val, ok := oracleMap[token0]; ok {

		isUSD := util.IsUSDOracle(val)

		tokenAddress := common.HexToAddress(val)
		instance, err := ChainLink.NewChainLink(tokenAddress, cl)
		if err != nil {
			log.Fatal(err)
		}
		bal, err := instance.LatestAnswer(callopts)
		if err != nil {
			log.Fatal(err)
		}

		if !isUSD {
			// If the oracle is not derived to USD, we need to derive it via ratio of tokens

		}

		token0RawPrice := util.DivideBy10pow(bal, 8)
		token0Price, err = strconv.ParseFloat(token0RawPrice.String(), 64)
		ratio := new(big.Float).Quo(amount0, amount1)

		ratioToInt, err := strconv.ParseFloat(ratio.String(), 64)

		b, err := strconv.ParseFloat(amount0.String(), 64)
		token1Price = token0Price * b
		amountusd = ratioToInt * token0Price

		// Temp Debugging
	} else if _, ok := oracleMap[token1]; ok {
		token0Price = 1
		token1Price = 1
		amountusd = 1
	} else {
		token0Price = 2
		token1Price = 2
		amountusd = 2
	}

	return token0Price, token1Price, amountusd
}
