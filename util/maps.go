package util

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

func GetOracleContracts(chain uint) (map[common.Address]string, error) {
	// Maps token address -> smart contract address for the oracle price of that token
	ethereum := make(map[common.Address]string)
	// WETH
	ethereum[common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")] = "0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419"
	// USDC
	ethereum[common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")] = "0x8fffffd4afb6115b954bd326cbe7b4ba576818f6"
	// USDT
	ethereum[common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")] = "0x3e7d1eab13ad0104d2750b8863b489d65364e32d"
	// DAI
	ethereum[common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")] = "0xaed0c38402a5d19df6e4c03f4e2dced6e29c1ee9"
	// WBTC
	ethereum[common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599")] = "0xf4030086522a5beea4988f8ca5b36dbc97bee88c"

	switch chain {
	case 1:
		return ethereum, nil
	}
	err1 := errors.New("Cannot find oracle map for provided Chain")
	return nil, err1
}

func IsUSDOracle(contract string) bool {
	switch contract {
	case
		"0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419",
		"0xf4030086522a5beea4988f8ca5b36dbc97bee88c":
		return true
	}
	return false
}
