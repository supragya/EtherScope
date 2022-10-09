package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/Blockpour/Blockpour-Geth-Indexer/abi/chainlink"
	"github.com/Blockpour/Blockpour-Geth-Indexer/gograph"
	"github.com/Blockpour/Blockpour-Geth-Indexer/mspool"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type OracleMap struct {
	Network        string `json:"Network"`
	ChainID        int    `json:"ChainID"`
	StableCoinsUSD []struct {
		Contract string     `json:"contract"`
		Price    *big.Float `json:"price"`
	} `json:"StableCoinsUSD"`
	Tokens []struct {
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
	stableCoins       map[common.Address]*big.Float
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
		stableCoins:       make(map[common.Address]*big.Float),
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
		pricing.stableCoins[common.HexToAddress(sc.Contract)] = sc.Price
	}

	// Setup tokenMap
	for _, token := range pricing.oracleMap.Tokens {
		pricing.tokenMap[common.HexToAddress(token.Contract)] = token.ID
	}

	return &pricing
}

func (d *EthRPC) GetRates2Tokens(
	callopts *bind.CallOpts,
	l types.Log,
	token0Address common.Address,
	token1Address common.Address,
	token0Amount *big.Float,
	token1Amount *big.Float) (token0PriceResult *itypes.PriceResult,
	token1Price *itypes.PriceResult,
	amountUSD *big.Float) {
	prs := d.GetPriceResults(callopts, []util.Tuple2[common.Address, *big.Float]{
		{token0Address, token0Amount},
		{token1Address, token1Amount},
	})

	// prevRates := []bool{rates[0] == nil, rates[1] == nil}

	if prs[0] == nil && prs[1] == nil {
		return nil, nil, nil
	} else if prs[0] == nil && token0Amount.Cmp(ZeroFloat) != 0 {
		cpdi := itypes.CounterpartyPriceDerivationInfo{
			CalculationTx:     l.TxHash,
			LogIdx:            l.Index,
			CounterpartyToken: token1Address,
			CounterpartyQty:   big.NewFloat(0.0).Set(token1Amount),
			CounterpartyPrice: big.NewFloat(0.0).Set(prs[1].Price),
			SelfQty:           big.NewFloat(0.0).Set(token0Amount),
		}

		numerator := big.NewFloat(1.0).Mul(cpdi.CounterpartyPrice, cpdi.CounterpartyQty)
		denominator := cpdi.SelfQty
		derivedPrice := big.NewFloat(1.0).Quo(numerator, denominator)

		entry := itypes.PriceResult{
			Price:                     derivedPrice,
			IsStablecoin:              false,
			IsDerivedFromCounterparty: true,
			CounterpartyInfo:          &cpdi,
			DerivationInfo:            nil,
		}

		// cache derived rate
		lookupKey := util.Tuple2[common.Address, bind.CallOpts]{token0Address, *callopts}
		d.PriceCache.Add(lookupKey, &entry)
		prs[0] = &entry

	} else if prs[1] == nil && token1Amount.Cmp(ZeroFloat) != 0 {
		cpdi := itypes.CounterpartyPriceDerivationInfo{
			CalculationTx:     l.TxHash,
			LogIdx:            l.Index,
			CounterpartyToken: token0Address,
			CounterpartyQty:   big.NewFloat(0.0).Set(token0Amount),
			CounterpartyPrice: big.NewFloat(0.0).Set(prs[0].Price),
			SelfQty:           big.NewFloat(0.0).Set(token1Amount),
		}

		numerator := big.NewFloat(1.0).Mul(cpdi.CounterpartyPrice, cpdi.CounterpartyQty)
		denominator := cpdi.SelfQty
		derivedPrice := big.NewFloat(1.0).Quo(numerator, denominator)

		entry := itypes.PriceResult{
			Price:                     derivedPrice,
			IsStablecoin:              false,
			IsDerivedFromCounterparty: true,
			CounterpartyInfo:          &cpdi,
			DerivationInfo:            nil,
		}

		// cache derived rate
		lookupKey := util.Tuple2[common.Address, bind.CallOpts]{token0Address, *callopts}
		d.PriceCache.Add(lookupKey, &entry)
		prs[1] = &entry
	}

	amountUSD = big.NewFloat(0.0)
	if prs[0] != nil && prs[1] != nil {
		if token0Amount.Cmp(ZeroFloat) != 0 {
			amountUSD.Mul(prs[0].Price, token0Amount)
		} else {
			amountUSD.Mul(prs[1].Price, token1Amount)
		}
	}

	return prs[0], prs[1], amountUSD
}

func (d *EthRPC) GetPriceResults(
	callopts *bind.CallOpts,
	requests []util.Tuple2[common.Address, *big.Float]) []*itypes.PriceResult {
	response := []*itypes.PriceResult{}
	for _, req := range requests {
		response = append(response, d.GetRateForBlock(callopts, req))
	}
	return response
}

func (d *EthRPC) GetRateForBlock(
	callopts *bind.CallOpts,
	request util.Tuple2[common.Address, *big.Float]) *itypes.PriceResult {

	// cache lookup
	lookupKey := util.Tuple2[common.Address, bind.CallOpts]{request.First, *callopts}

	if val, ok := d.PriceCache.Get(lookupKey); ok {
		return val.(*itypes.PriceResult)
	}

	// Is stablecoin
	if stablecoinPrice, ok := d.pricing.stableCoins[request.First]; ok {
		return &itypes.PriceResult{
			Price:                     big.NewFloat(1.0).Set(stablecoinPrice),
			IsStablecoin:              true,
			IsDerivedFromCounterparty: false,
			CounterpartyInfo:          nil,
			DerivationInfo:            nil,
		}
	}

	// if a known token
	if tokenID, ok := d.pricing.tokenMap[request.First]; ok {
		route := d.pricing.graph.GetShortestRoute(tokenID, "USD")
		price := big.NewFloat(1.0)

		if len(route.Edges)+1 != len(route.Vertices) {
			panic("error in APSP routing logic")
		}

		routingMetadata := make(map[string]itypes.DirectPriceDerivationInfo, len(route.Edges))

		for idx, edge := range route.Edges {
			oracleContractAddress := common.HexToAddress(edge.Metadata)
			edgeHumanReadable := fmt.Sprintf("chainlink %s-%s(%s)", route.Vertices[idx], route.Vertices[idx+1], oracleContractAddress)

			latestRoundData, err := mspool.Do(d.upstreams,
				func(ctx context.Context, c *ethclient.Client) (itypes.ChainlinkLatestRoundData, error) {
					oracle, err := chainlink.NewChainlink(oracleContractAddress, c)
					if err != nil {
						return itypes.ChainlinkLatestRoundData{}, err
					}
					callopts.Context = ctx
					return oracle.LatestRoundData(callopts)
				}, itypes.ChainlinkLatestRoundData{})
			util.ENOK(err)

			decimals, err := mspool.Do(d.upstreams,
				func(ctx context.Context, c *ethclient.Client) (uint8, error) {
					oracle, err := chainlink.NewChainlink(oracleContractAddress, c)
					if err != nil {
						return 0, err
					}
					callopts.Context = ctx
					return oracle.Decimals(callopts)
				}, 0)
			util.ENOK(err)

			tokenFormatted := util.DivideBy10pow(latestRoundData.Answer, decimals)

			// Assuming base currency to be worth 1USD
			price = price.Mul(price, tokenFormatted)

			// Add to metadata
			routingMetadata[edgeHumanReadable] = itypes.DirectPriceDerivationInfo{
				LatestRoundData: latestRoundData,
				Decimals:        decimals,
				ConversionPrice: tokenFormatted,
			}
		}

		result := itypes.PriceResult{
			Price:                     price,
			IsStablecoin:              false,
			IsDerivedFromCounterparty: false,
			CounterpartyInfo:          nil,
			DerivationInfo:            routingMetadata,
		}

		// Cache insert
		d.PriceCache.Add(lookupKey, &result)
		return &result
	}

	return nil
}

func init() {
	ZeroFloat = big.NewFloat(0.0)
}
