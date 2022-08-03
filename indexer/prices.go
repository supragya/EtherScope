package indexer

import (
	"math"
	"math/big"
	"strconv"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/ChainLink"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func fetchBaseCurrency(callopts *bind.CallOpts, cl *ethclient.Client) float64 {
	networkID := viper.GetUint("general.chainid")

	nativeTokenUSD := util.BaseNativeToken(networkID)

	tokenAddress := common.HexToAddress(nativeTokenUSD)
	instance, err := ChainLink.NewChainLink(tokenAddress, cl)
	decimals, _ := instance.Decimals(callopts)

	if err != nil {
		log.Fatal(err)
	}

	nativeTokenUSDPrice, err := instance.LatestAnswer(callopts)

	if err != nil {
		log.Fatal(err)
	}

	nativeUSDPriceFormatted := util.DivideBy10pow(nativeTokenUSDPrice, decimals)
	nativeUSDPriceParsed, err := strconv.ParseFloat(nativeUSDPriceFormatted.String(), 64)
	return nativeUSDPriceParsed
}

func (d *DataAccess) GetPricesForBlock(
	callopts *bind.CallOpts, token0 common.Address, token1 common.Address, amount0 *big.Float, amount1 *big.Float) (float64, float64, float64) {

	networkID := viper.GetUint("general.chainid")
	oracleMap, err := util.GetOracleContracts(networkID)

	token0Amount, _ := strconv.ParseFloat(amount0.String(), 64)
	token1Amount, _ := strconv.ParseFloat(amount1.String(), 64)
	util.ENOK(err)

	cl := d.upstreams.GetItem()

	var token0Price float64
	var token1Price float64
	var amountusd float64
	var isUSD bool

	// Run function if token0 maps to an oracle
	if token0OracleAddress, token0Oracle := oracleMap[token0]; token0Oracle {

		isUSD = util.IsUSDOracle(token0OracleAddress)

		tokenAddress := common.HexToAddress(token0OracleAddress)
		instance, err := ChainLink.NewChainLink(tokenAddress, cl)
		if err != nil {
			log.Fatal(err)
		}
		token0LastPrice, _ := instance.LatestAnswer(callopts)
		decimals, _ := instance.Decimals(callopts)

		token0Formatted := util.DivideBy10pow(token0LastPrice, decimals)
		token0Price, err = strconv.ParseFloat(token0Formatted.String(), 64)

		// Divides amount0 by amount1 to get ratio of tokens
		ratio := new(big.Float).Quo(amount0, amount1)
		ratioToInt, _ := strconv.ParseFloat(ratio.String(), 64)

		if !isUSD {
			baseCurrency := fetchBaseCurrency(callopts, cl)
			token0Price = baseCurrency * token0Price
			token1Price = ratioToInt
			amountusd = token0Price * token0Amount
		} else {
			token1Price = ratioToInt * token0Price
			amountusd = token0Price * token0Amount
		}
	} else if token1OracleAddress, token1Oracle := oracleMap[token1]; token1Oracle {

		isUSD = util.IsUSDOracle(token1OracleAddress)
		tokenAddress := common.HexToAddress(token1OracleAddress)
		instance, err := ChainLink.NewChainLink(tokenAddress, cl)
		if err != nil {
			log.Fatal(err)
		}
		token1LastPrice, _ := instance.LatestAnswer(callopts)
		decimals, _ := instance.Decimals(callopts)

		token1Formatted := util.DivideBy10pow(token1LastPrice, decimals)
		token1Price, err = strconv.ParseFloat(token1Formatted.String(), 64)

		// Divides amount0 by amount1 to get ratio of tokens
		ratio := new(big.Float).Quo(amount1, amount0)
		ratioToInt, _ := strconv.ParseFloat(ratio.String(), 64)

		if !isUSD {
			baseCurrency := fetchBaseCurrency(callopts, cl)
			token1Price = baseCurrency * token1Price
			token0Price = ratioToInt
			amountusd = token1Price * token1Amount
		} else {
			token0Price = ratioToInt * token1Price
			amountusd = token1Price * token1Amount
		}
	} else {
		token0Price = 1
		token1Price = 1
		amountusd = 1
	}

	return math.Abs(token0Price), math.Abs(token1Price), math.Abs(amountusd)
}
