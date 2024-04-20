package priceresolver

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"

	gograph "github.com/supragya/EtherScope/libs/oldgograph"
	"github.com/supragya/EtherScope/libs/util"
	"github.com/supragya/EtherScope/services/ethrpc"
	itypes "github.com/supragya/EtherScope/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
)

type OracleMap struct {
	Network        string `json:"Network"`
	ChainID        int    `json:"ChainID"`
	StableCoinsUSD []struct {
		ID             string `json:"id"`
		Contract       string `json:"contract"`
		OracleContract string `json:"oracleContract"`
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

type ChainlinkLatestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

type DirectPriceDerivationInfo struct {
	LatestRoundData itypes.ChainlinkLatestRoundData
	Decimals        uint8
	ConversionPrice *big.Float
}

type CounterpartyPriceDerivationInfo struct {
	CalculationTx     common.Hash
	LogIdx            uint
	CounterpartyToken common.Address
	CounterpartyQty   *big.Float
	CounterpartyPrice *big.Float
	SelfQty           *big.Float
}

type PriceResult struct {
	Price                     *big.Float
	IsStablecoin              bool
	IsDerivedFromCounterparty bool
	CounterpartyInfo          *CounterpartyPriceDerivationInfo
	DerivationInfo            map[string]DirectPriceDerivationInfo
}

type Pricing struct {
	oracleFile  string
	oracleHash  string
	graph       *gograph.Graph[string, string]
	stableCoins map[common.Address]itypes.Tuple2[string, common.Address]
	tokenMap    map[common.Address]string
	oracleMap   OracleMap
	EthRPC      ethrpc.EthRPC

	contractTokensCache *lru.ARCCache
	ERC20Cache          *lru.ARCCache
	PriceCache          *lru.ARCCache
}

func GetPricingEngine(oracleMapLocation string, rpc ethrpc.EthRPC) *Pricing {
	ctcache, err := lru.NewARC(1024) // Hardcoded 1024
	util.ENOK(err)

	erc20cache, err := lru.NewARC(1024) // Hardcoded 1024
	util.ENOK(err)

	ratecache, err := lru.NewARC(1024) // Hardcoded 1024
	util.ENOK(err)

	pricing := Pricing{
		oracleFile:          oracleMapLocation,
		oracleHash:          "", // Computed below
		graph:               gograph.NewGraphStringUintString(false),
		stableCoins:         make(map[common.Address]itypes.Tuple2[string, common.Address]),
		tokenMap:            make(map[common.Address]string),
		oracleMap:           OracleMap{},
		contractTokensCache: ctcache,
		ERC20Cache:          erc20cache,
		PriceCache:          ratecache,
		EthRPC:              rpc,
	}

	// Open oracle file
	fd, err := os.Open(pricing.oracleFile)
	util.ENOK(err)
	defer fd.Close()
	fileBytes, err := io.ReadAll(fd)
	util.ENOK(err)

	// Load bare oracle Map
	util.ENOK(json.Unmarshal(fileBytes, &pricing.oracleMap))

	for _, o := range pricing.oracleMap.Oracles {
		if err := pricing.graph.AddEdge(o.From, o.To, 1, o.Contract); err != nil {
			panic(fmt.Sprintf("error while adding edge: %s %s->%s as %s", err, o.From, o.To, o.Contract))
		}
	}

	// Setup stablecoins
	for _, sc := range pricing.oracleMap.StableCoinsUSD {
		pricing.stableCoins[common.HexToAddress(sc.Contract)] = itypes.Tuple2[string, common.Address]{
			sc.ID,
			common.HexToAddress(sc.OracleContract),
		}
	}

	// Setup tokenMap
	for _, token := range pricing.oracleMap.Tokens {
		pricing.tokenMap[common.HexToAddress(token.Contract)] = token.ID
	}

	pricing.graph.CalculateAllPairShortestPath()

	return &pricing
}

func (d *Pricing) GetRates2Tokens(
	callopts *bind.CallOpts,
	token0Address common.Address,
	token1Address common.Address,
	token0Amount *big.Float,
	token1Amount *big.Float) (token0PriceResult *itypes.PriceResult,
	token1Price *itypes.PriceResult,
	amountUSD *big.Float) {
	prs := d.GetPriceResults(callopts, []itypes.Tuple2[common.Address, *big.Float]{
		{token0Address, token0Amount},
		{token1Address, token1Amount},
	})

	amountUSD = big.NewFloat(0.0)

	if prs[0] == nil && prs[1] == nil {
		return nil, nil, nil
	} else if prs[0] == nil && token0Amount.Cmp(big.NewFloat(0.0)) != 0 {
		cpdi := CounterpartyPriceDerivationInfo{
			CounterpartyToken: token1Address,
			CounterpartyQty:   big.NewFloat(0.0).Set(token1Amount),
			CounterpartyPrice: big.NewFloat(0.0).Set(prs[1].Price),
			SelfQty:           big.NewFloat(0.0).Set(token0Amount),
		}

		numerator := big.NewFloat(1.0).Mul(cpdi.CounterpartyPrice, cpdi.CounterpartyQty)
		denominator := cpdi.SelfQty
		derivedPrice := big.NewFloat(1.0).Quo(numerator, denominator)
		derivedPrice = derivedPrice.Abs(derivedPrice)

		entry := itypes.PriceResult{
			Price: derivedPrice,
			// We don't have path anymore TODO
		}

		// cache derived rate
		// lookupKey := itypes.Tuple2[common.Address, bind.CallOpts]{token0Address, *callopts}
		// d.PriceCache.Add(lookupKey, &entry)
		prs[0] = &entry
		amountUSD.Mul(prs[1].Price, token1Amount)
	} else if prs[1] == nil && token1Amount.Cmp(big.NewFloat(0.0)) != 0 {
		cpdi := CounterpartyPriceDerivationInfo{
			CounterpartyToken: token0Address,
			CounterpartyQty:   big.NewFloat(0.0).Set(token0Amount),
			CounterpartyPrice: big.NewFloat(0.0).Set(prs[0].Price),
			SelfQty:           big.NewFloat(0.0).Set(token1Amount),
		}

		numerator := big.NewFloat(1.0).Mul(cpdi.CounterpartyPrice, cpdi.CounterpartyQty)
		denominator := cpdi.SelfQty
		derivedPrice := big.NewFloat(1.0).Quo(numerator, denominator)
		derivedPrice = derivedPrice.Abs(derivedPrice)

		entry := itypes.PriceResult{
			Price: derivedPrice,
			// We don't have path anymore TODO
		}

		// cache derived rate
		// lookupKey := itypes.Tuple2[common.Address, bind.CallOpts]{token1Address, *callopts}
		// d.PriceCache.Add(lookupKey, &entry)
		prs[1] = &entry
		amountUSD.Mul(prs[0].Price, token0Amount)
	}

	if amountUSD.Cmp(big.NewFloat(0.0)) == 0 && prs[0] != nil && prs[1] != nil {
		if token0Amount.Cmp(big.NewFloat(0.0)) != 0 {
			amountUSD.Mul(prs[0].Price, token0Amount)
		} else {
			amountUSD.Mul(prs[1].Price, token1Amount)
		}
	}

	return prs[0], prs[1], amountUSD.Abs(amountUSD)
}

func (d *Pricing) GetPriceResults(
	callopts *bind.CallOpts,
	requests []itypes.Tuple2[common.Address, *big.Float]) []*itypes.PriceResult {
	response := []*itypes.PriceResult{}
	for _, req := range requests {
		response = append(response, d.GetRateForBlock(callopts, req))
	}
	return response
}

func (d *Pricing) GetRateForBlock(
	callopts *bind.CallOpts,
	request itypes.Tuple2[common.Address, *big.Float]) *itypes.PriceResult {

	// cache lookup
	lookupKey := itypes.Tuple2[common.Address, bind.CallOpts]{request.First, *callopts}

	if val, ok := d.PriceCache.Get(lookupKey); ok {
		return val.(*itypes.PriceResult)
	}

	// Is stablecoin
	if info, ok := d.stableCoins[request.First]; ok {
		var (
			oracleContractAddress = info.Second
			edgeHumanReadable     = fmt.Sprintf("chainlink %s-USD(%s)", info.First, oracleContractAddress)
		)

		latestRoundData, err := d.EthRPC.GetChainlinkRoundData(oracleContractAddress, callopts)
		if err != nil {
			// Return unknown result
			return &itypes.PriceResult{}
		}

		decimals, err := d.EthRPC.GetChainlinkDecimals(oracleContractAddress, callopts)
		if err != nil {
			// Return unknown result
			return &itypes.PriceResult{}
		}

		tokenFormatted := util.DivideBy10pow(latestRoundData.Answer, decimals)
		routingMetadata := make(map[string]DirectPriceDerivationInfo, 1)
		routingMetadata[edgeHumanReadable] = DirectPriceDerivationInfo{
			LatestRoundData: latestRoundData,
			Decimals:        decimals,
			ConversionPrice: tokenFormatted,
		}

		result := itypes.PriceResult{
			Price: tokenFormatted,
			// We don't have path anymore TODO
		}
		d.PriceCache.Add(lookupKey, &result)

		return &result
	}

	// if a known token
	if tokenID, ok := d.tokenMap[request.First]; ok {
		route := d.graph.GetShortestRoute(tokenID, "USD")
		price := big.NewFloat(1.0)

		if len(route.Edges)+1 != len(route.Vertices) {
			panic(fmt.Sprintf("error in APSP routing logic token: %s edges: %d vs vertices: %d, totals: %d %d",
				tokenID, len(route.Edges), len(route.Vertices), d.graph.GetEdgeCount(), d.graph.GetVertexCount()))
		}

		routingMetadata := make(map[string]DirectPriceDerivationInfo, len(route.Edges))

		for idx, edge := range route.Edges {
			oracleContractAddress := common.HexToAddress(edge.Metadata)
			edgeHumanReadable := fmt.Sprintf("chainlink %s-%s(%s)", route.Vertices[idx], route.Vertices[idx+1], oracleContractAddress)

			latestRoundData, err := d.EthRPC.GetChainlinkRoundData(oracleContractAddress, callopts)
			if err != nil {
				// Return unknown result
				return &itypes.PriceResult{}
			}

			decimals, err := d.EthRPC.GetChainlinkDecimals(oracleContractAddress, callopts)
			if err != nil {
				// Return unknown result
				return &itypes.PriceResult{}
			}

			tokenFormatted := util.DivideBy10pow(latestRoundData.Answer, decimals)

			// Assuming base currency to be worth 1USD
			price = price.Mul(price, tokenFormatted)

			// Add to metadata
			routingMetadata[edgeHumanReadable] = DirectPriceDerivationInfo{
				LatestRoundData: latestRoundData,
				Decimals:        decimals,
				ConversionPrice: tokenFormatted,
			}
		}
		price = price.Abs(price)

		result := itypes.PriceResult{
			Price: price,
			// We don't have path anymore TODO
		}

		// Cache insert
		d.PriceCache.Add(lookupKey, &result)
		return &result
	}

	return nil
}

func (d *Pricing) Resolve(resHeight uint64, items []interface{}) error {
	callopts := &bind.CallOpts{BlockNumber: big.NewInt(int64(resHeight))}
	for _, item := range items {
		switch i := item.(type) {
		case *itypes.Mint:
			if i.ProcessingType != itypes.UserRequested {
				continue
			}
			i.Price0, i.Price1, i.AmountUSD = d.GetRates2Tokens(callopts, i.Token0, i.Token1, i.Amount0, i.Amount1)

		case *itypes.Burn:
			if i.ProcessingType != itypes.UserRequested {
				continue
			}
			i.Price0, i.Price1, i.AmountUSD = d.GetRates2Tokens(callopts, i.Token0, i.Token1, i.Amount0, i.Amount1)

		case *itypes.Swap:
			if i.ProcessingType != itypes.UserRequested {
				continue
			}
			i.Price0, i.Price1, i.AmountUSD = d.GetRates2Tokens(callopts, i.Token0, i.Token1, i.Amount0, i.Amount1)

		}
	}

	return nil
}
