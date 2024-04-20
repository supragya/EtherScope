package common

import "github.com/supragya/EtherScope/algorand/types"

const tinymanV1AppId = 552635992
const tinymanV2AppId = 1002541853

const ALGOUSDC = "FPOU46NBKTWUZCNMNQNXRWNW3SMPOOK4ZJIN5WSILCWP662ANJLTXVRUKA"
const STBL2USDC = 855716333

var ALGO = types.TokenInfo{
	Id:       1,
	Name:     "ALGO",
	Decimals: 6,
}

var USDC = types.TokenInfo{
	Id:       31566704,
	Name:     "USDC",
	Decimals: 6,
}

var STBL = types.TokenInfo{
	Id:       465865291,
	Name:     "STBL",
	Decimals: 6,
}

var STBL2 = types.TokenInfo{
	Id:       841126810,
	Name:     "STBL2",
	Decimals: 6,
}

var PDAT = types.TokenInfo{
	Id:       919889450,
	Name:     "PDAT",
	Decimals: 4,
}

var USDt = types.TokenInfo{
	Id:       312769,
	Name:     "USDt",
	Decimals: 6,
}

var Stablecoins = map[uint64]types.TokenInfo{
	312769:   USDt,
	31566704: USDC,
}

/*
List of quotes that are safe to price against
TODO: Add STBL2 which has enough liquidity (not STBL1)
*/
var Quotes = map[uint64]types.TokenInfo{
	ALGO.Id: ALGO,
	USDC.Id: USDC,
	USDt.Id: USDt,
}

var AlgoFiSwapSignatures = []types.FunctionSignature{
	{
		AppId: 0,
		Key:   "sef",
	},
	{
		AppId: 0,
		Key:   "sfe",
	},
}

var TinymanV1SwapSignature = types.FunctionSignature{
	AppId: tinymanV1AppId,
	Key:   "swap",
}

var TinymanV2SwapSignature = types.FunctionSignature{
	AppId: tinymanV2AppId,
	Key:   "swap",
}

var SupportedFunctionSignatures = []types.FunctionSignature{
	TinymanV1SwapSignature,
	TinymanV2SwapSignature,
	AlgoFiSwapSignatures[0],
	AlgoFiSwapSignatures[1],
}
