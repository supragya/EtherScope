package util

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

func GetOracleContracts(chain uint) (map[common.Address]string, error) {
	// Maps token address -> smart contract address for the oracle price of that token
	ethereum := make(map[common.Address]string)
	bsc := make(map[common.Address]string)
	moonbeam := make(map[common.Address]string)
	polygon := make(map[common.Address]string)
	ethereum[common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")] = "0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419" // WETH / USD
	ethereum[common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599")] = "0xf4030086522a5beea4988f8ca5b36dbc97bee88c" // WBTC / USD
	ethereum[common.HexToAddress("0x514910771af9ca656af840dff83e8264ecf986ca")] = "0x2c1d072e956affc0d435cb7ac38ef18d24d9127c" // LINK / USD
	ethereum[common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")] = "0xee9f2375b4bdf6387aa8265dd4fb8f16512a1d46" // USDT / ETH
	ethereum[common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")] = "0x986b5e1e1755e3c2440e960477f25201b0a8bbd4" // USDC / ETH
	ethereum[common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")] = "0x773616e4d11a78f511299002da57a0a94577f1f4" // DAI / ETH
	ethereum[common.HexToAddress("0xc011a73ee8576fb46f5e1c5751ca3b9fe0af2a6f")] = "0x79291a9d692df95334b1a0b3b4ae6bc606782f8c" // SNX / ETH
	ethereum[common.HexToAddress("0xbbbbca6a901c926f240b89eacb641d8aec7aeafd")] = "0x160ac928a16c93ed4895c2de6f81ecce9a7eb7b4" // LRC / ETH
	ethereum[common.HexToAddress("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2")] = "0x24551a8fb2a7211a25a17b1481f043a8a8adc7f2" // MKR / ETH
	ethereum[common.HexToAddress("0x0f5d2fb29fb7d3cfee444a200298f468908cc942")] = "0x82a44d92d6c329826dc557c5e1be6ebec5d5feb9" // MANA/ ETH
	ethereum[common.HexToAddress("0xdd974d5c2e2928dea5f71b9825b8b646686bd200")] = "0x656c0544ef4c98a6a98491833a89204abb045d6b" // KNC / ETH
	ethereum[common.HexToAddress("0x0d8775f648430679a709e98d2b0cb6250d2887ef")] = "0x0d16d4528239e9ee52fa531af613acdb23d88c94" // BAT / ETH
	ethereum[common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")] = "0x2da4983a622a8498bb1a21fae9d8f6c664939962" // ZRX / ETH
	ethereum[common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376")] = "0x3886ba987236181d98f2401c507fb8bea7871df2" // MANA/ ETH
	ethereum[common.HexToAddress("0x57ab1ec28d129707052df4df418d58a2d46d5f51")] = "0x8e0b7e6062272b5ef4524250bfff8e5bd3497757" // KNC / ETH

	bsc[common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")] = "0x0567f2323251f0aab15c8dfb1967e4e8a7d42aee" // BNB / USD
	bsc[common.HexToAddress("0x2170ed0880ac9a755fd29b2688956bd959f933f8")] = "0x9ef1b8c0e4f7dc8bf5719ea496883dc6401d5b2e" // ETH / USD
	bsc[common.HexToAddress("0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c")] = "0x264990fbd0a4796a3e3d8e37c4d5f87a3aca5ebf" // BTC / USD
	bsc[common.HexToAddress("0x55d398326f99059ff775485246999027b3197955")] = "0xd5c40f5144848bd4ef08a9605d860e727b991513" // USDT / BNB
	bsc[common.HexToAddress("0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d")] = "0x45f86ca2a8bc9ebd757225b19a1a0d7051be46db" // USDC / BNB

	moonbeam[common.HexToAddress("0x818ec0a7fe18ff94269904fced6ae3dae6d6dc0b")] = "0xA122591F60115D63421f66F752EF9f6e0bc73abC" // USDC / USD
	moonbeam[common.HexToAddress("0xacc15dc74880c9944775448304b263d191c6077f")] = "0x4497B606be93e773bbA5eaCFCb2ac5E2214220Eb" // GLMR / USD
	moonbeam[common.HexToAddress("0xfa9343c3897324496a05fc75abed6bac29f8a40f")] = "0x9ce2388a1696e22F870341C3FC1E89710C7569B5" // ETH / USD
	moonbeam[common.HexToAddress("0x922d641a426dcffaef11680e5358f34d97d112e1")] = "0x8c4425e141979c66423A83bE2ee59135864487Eb" // BTC / USD
	// moonbeam[common.HexToAddress("0x922d641a426dcffaef11680e5358f34d97d112e1")] = "0xd61D7398B7734aBe7C4B143fE57dC666D2fe83aD" // LINK / USD

	polygon[common.HexToAddress("0xc2132d05d31c914a87c6611c10748aeb04b58e8f")] = "0x0a6513e40db6eb1b165753ad52e80663aea50545" // USDT / USD
	polygon[common.HexToAddress("0x8f3cf7ad23cd3cadbd9735aff958023239c6a063")] = "0x4746dec9e833a82ec7c2c1356372ccf2cfcd2f3d" // DAI / USD
	polygon[common.HexToAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")] = "0xfe4a8cc5b5b2366c1b58bea3858e81843581b2f7" // USDC / USD
	polygon[common.HexToAddress("0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270")] = "0xab594600376ec9fd91f8e885dadf0ce036862de0" // MATIC / USD
	polygon[common.HexToAddress("0x7ceb23fd6bc0add59e62ac25578270cff1b9f619")] = "0xf9680d99d6c9589e2a93a78a04a279e509205945" // ETH / USD
	polygon[common.HexToAddress("0x1bfd67037b42cf73acf2047067bd4f2c47d9bfd6")] = "0xc907e116054ad103354f2d350fd2514433d57f6f" // WBTC / USD

	switch chain {
	case 1:
		return ethereum, nil
	case 56:
		return bsc, nil
	case 1284:
		return moonbeam, nil
	case 137:
		return polygon, nil
	}
	err := errors.New("Cannot find oracle map for provided Chain")
	return nil, err
}

func IsUSDOracle(contract string) bool {
	// Includes USD based oracle smart contracts across all networks
	switch contract {
	case
		"0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419",
		"0xf4030086522a5beea4988f8ca5b36dbc97bee88c",
		"0x2c1d072e956affc0d435cb7ac38ef18d24d9127c",
		"0x0567f2323251f0aab15c8dfb1967e4e8a7d42aee",
		"0x9ef1b8c0e4f7dc8bf5719ea496883dc6401d5b2e",
		"0x264990fbd0a4796a3e3d8e37c4d5f87a3aca5ebf",
		"0xA122591F60115D63421f66F752EF9f6e0bc73abC",
		"0x4497B606be93e773bbA5eaCFCb2ac5E2214220Eb",
		"0x9ce2388a1696e22F870341C3FC1E89710C7569B5",
		"0x8c4425e141979c66423A83bE2ee59135864487Eb",
		"0x0a6513e40db6eb1b165753ad52e80663aea50545",
		"0x4746dec9e833a82ec7c2c1356372ccf2cfcd2f3d",
		"0xfe4a8cc5b5b2366c1b58bea3858e81843581b2f7",
		"0xab594600376ec9fd91f8e885dadf0ce036862de0",
		"0xf9680d99d6c9589e2a93a78a04a279e509205945",
		"0xc907e116054ad103354f2d350fd2514433d57f6f":
		return true
	}
	return false
}

func BaseNativeToken(chain uint) string {
	switch chain {
	case
		1:
		return "0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419" // WETH / USD Oracle
	case 56:
		return "0x0567f2323251f0aab15c8dfb1967e4e8a7d42aee"
	case 1284:
		return "0x4497B606be93e773bbA5eaCFCb2ac5E2214220Eb"
	case 137:
		return "0xab594600376ec9fd91f8e885dadf0ce036862de0"
	}
	return "null"
}
