package indexer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/chainlink"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/supragya/gograph"
)

type OracleMap struct {
	Network        string   `json:"Network"`
	ChainID        int      `json:"ChainID"`
	StableCoinsUSD []string `json:"StableCoinsUSD"`
	Tokens         []struct {
		ID       string `json:"id"`
		Contract string `json:"contract"`
	} `json:"Tokens"`
	Oracles []struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Contract string `json:"contract"`
	} `json:"Oracles"`
}

var (
	ZeroFloat *big.Float
)

type Pricing struct {
	oracleMapsRootDir string
	diskCacheRootDir  string
	networkName       string
	oracleFile        string
	cacheFile         string
	oracleHash        string
	graph             *gograph.Graph[string, string]
	stableCoins       map[common.Address]bool
	tokenMap          map[common.Address]string
	oracleMap         OracleMap
}

func GetPricingEngine() *Pricing {
	pricing := Pricing{
		oracleMapsRootDir: viper.GetString("general.oracleMapsRootDir"),
		diskCacheRootDir:  viper.GetString("general.diskCacheRootDir"),
		networkName:       viper.GetString("general.networkName"),
		oracleFile:        viper.GetString("general.oracleMapsRootDir") + "/oraclemaps_" + viper.GetString("general.networkName") + ".json",
		oracleHash:        "", // Computed below
		graph:             gograph.NewGraphStringUintString(false),
		stableCoins:       make(map[common.Address]bool),
		tokenMap:          make(map[common.Address]string),
		oracleMap:         OracleMap{},
	}

	// Open oracle file
	fd, err := os.Open(pricing.oracleFile)
	util.ENOK(err)
	defer fd.Close()
	fileBytes, err := io.ReadAll(fd)
	util.ENOK(err)

	// Load bare oracle Map
	util.ENOK(json.Unmarshal(fileBytes, &pricing.oracleMap))

	pricing.oracleHash = hex.EncodeToString(util.SHA256Hash(fileBytes))

	pricing.cacheFile = fmt.Sprintf("%s/prcg_%s_%s.bincache",
		pricing.diskCacheRootDir,
		pricing.networkName,
		pricing.oracleHash)

	if _, err := os.Stat(pricing.cacheFile); err != nil {
		// Cache doesn't exists
		log.Info("prewarming pricing graph cache for future runs")
		tempGraph := gograph.NewGraphStringUintString(false)
		for _, o := range pricing.oracleMap.Oracles {
			tempGraph.AddEdge(o.From, o.To, 1, o.Contract)
		}
		tempGraph.SaveToDisk(pricing.cacheFile)
	} else {
		log.Info("loading up prewarmed pricing graph cache")
	}

	// Load up pricing
	util.ENOK(pricing.graph.ReadFromDisk(pricing.cacheFile))
	log.Info("loaded up pricing graph")

	// Setup stablecoins
	for _, sc := range pricing.oracleMap.StableCoinsUSD {
		pricing.stableCoins[common.HexToAddress(sc)] = true
	}

	// Setup tokenMap
	for _, token := range pricing.oracleMap.Tokens {
		pricing.tokenMap[common.HexToAddress(token.Contract)] = token.ID
	}

	return &pricing
}

func (d *DataAccess) GetRates2Tokens(
	callopts *bind.CallOpts,
	token0Address common.Address,
	token1Address common.Address,
	token0Amount *big.Float,
	token1Amount *big.Float) (token0Price *big.Float,
	token1Price *big.Float,
	amountUSD *big.Float) {
	rates := d.GetRatesForBlock(callopts, []Tuple2[common.Address, *big.Float]{
		{token0Address, token0Amount},
		{token1Address, token1Amount},
	})

	if rates[0] == nil && rates[1] == nil {
		return nil, nil, nil
	} else if rates[0] == nil && token0Amount.Cmp(ZeroFloat) != 0 {
		numerator := big.NewFloat(1.0).Mul(rates[1], token1Amount)
		denominator := token0Amount
		rates[0] = big.NewFloat(1.0).Quo(numerator, denominator)
		// cache derived rate
		lookupKey := Tuple2[common.Address, bind.CallOpts]{token0Address, *callopts}
		d.RateCache.Add(lookupKey, rates[0])

	} else if rates[1] == nil && token1Amount.Cmp(ZeroFloat) != 0 {
		numerator := big.NewFloat(1.0).Mul(rates[0], token0Amount)
		denominator := token1Amount
		rates[1] = big.NewFloat(1.0).Quo(numerator, denominator)
		// cache derived rate
		lookupKey := Tuple2[common.Address, bind.CallOpts]{token1Address, *callopts}
		d.RateCache.Add(lookupKey, rates[1])
	}

	amountUSD = big.NewFloat(0.0)
	if rates[0] != nil && rates[1] != nil {
		if token0Amount.Cmp(ZeroFloat) != 0 {
			amountUSD.Mul(rates[0], token0Amount)
		} else {
			amountUSD.Mul(rates[1], token1Amount)
		}
	}

	return rates[0], rates[1], amountUSD
}

func (d *DataAccess) GetRatesForBlock(
	callopts *bind.CallOpts,
	requests []Tuple2[common.Address, *big.Float]) []*big.Float {
	response := []*big.Float{}
	for _, req := range requests {
		response = append(response, d.GetRateForBlock(callopts, req))
	}
	return response
}

func (d *DataAccess) GetRateForBlock(
	callopts *bind.CallOpts,
	request Tuple2[common.Address, *big.Float]) *big.Float {

	// cache lookup
	lookupKey := Tuple2[common.Address, bind.CallOpts]{request.First, *callopts}

	if val, ok := d.RateCache.Get(lookupKey); ok {
		return val.(*big.Float)
	}

	// Is stablecoin
	if _, ok := d.pricing.stableCoins[request.First]; ok {
		return request.Second
	}

	// if a known token
	if tokenID, ok := d.pricing.tokenMap[request.First]; ok {
		route := d.pricing.graph.GetShortestRoute(tokenID, "USD")
		rate := big.NewFloat(1.0)

		for _, edge := range route.Edges {
			oracleContractAddress := common.HexToAddress(edge.Metadata)

			for retries := 0; retries < WD; retries++ {
				cl := d.upstreams.GetItem()
				oracle, err := chainlink.NewChainlink(oracleContractAddress, cl)
				util.ENOK(err)

				latestRoundData, err := oracle.LatestRoundData(callopts)
				if err != nil {
					if util.IsEthErr(err) {
						d.upstreams.Report(cl, false)
						return nil
					}
					d.upstreams.Report(cl, true)
					continue
				}

				decimals, err := oracle.Decimals(callopts)
				if err != nil {
					if util.IsEthErr(err) {
						d.upstreams.Report(cl, false)
						return nil
					}
					d.upstreams.Report(cl, true)
					continue
				}

				tokenFormatted := util.DivideBy10pow(latestRoundData.Answer, decimals)

				// Assuming base currency to be worth 1USD
				rate = rate.Mul(rate, tokenFormatted)
				break
			}
		}

		// Cache insert
		d.RateCache.Add(lookupKey, rate)
		return rate
	}

	return nil
}

func init() {
	ZeroFloat = big.NewFloat(0.0)
}
