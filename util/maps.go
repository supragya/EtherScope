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
	avax := make(map[common.Address]string)
	ftm := make(map[common.Address]string)
	optimism := make(map[common.Address]string)
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

	avax[common.HexToAddress("0x5947bb275c521040051d82396192181b413227a3")] = "0x49ccd9ca821efeab2b98c60dc60f518e765ede9a" // LINK / USD
	avax[common.HexToAddress("0xd586e7f844cea2f87f50152665bcbc2c279d8d70")] = "0x51d7180eda2260cc4f6e4eebb82fef5c3c2b8300" // DAI / USD
	avax[common.HexToAddress("0xc7198437980c041c805a1edcba50c1ce5db95118")] = "0xebe676ee90fe1112671f19b6b7459bc678b67e8a" // USDT / USD
	avax[common.HexToAddress("0x50b7545627a5162f82a992c33b87adc75187b218")] = "0x2779d32d5166baaa2b2b658333ba7e6ec0c65743" // BTC / USD
	avax[common.HexToAddress("0x152b9d0FdC40C096757F570A51E494bd4b943E50")] = "0x2779d32d5166baaa2b2b658333ba7e6ec0c65743" // BTC / USD
	avax[common.HexToAddress("0x49d5c2bdffac6ce2bfdb6640f4f80f226bc10bab")] = "0x976b3d034e162d8bd72d6b9c989d545b839003b0" // ETH / USD
	avax[common.HexToAddress("0xb31f66aa3c1e785363f0875a1b74e27b85fd66c7")] = "0x0a77230d17318075983913bc2145db16c7366156" // AVAX / USD
	avax[common.HexToAddress("0x63a72806098bd3d9520cc43356dd78afe5d386d9")] = "0x3ca13391e9fb38a75330fb28f8cc2eb3d9ceceed" // AAVE / USD
	avax[common.HexToAddress("0xb97ef9ef8734c71904d8002f8b6bc66dd9c48a6e")] = "0xf096872672f44d6eba71458d74fe67f9a77a23b9" // USDC / USD
	avax[common.HexToAddress("0xa7d7079b0fead91f3e65f86e8915cb59c1a4c664")] = "0xf096872672f44d6eba71458d74fe67f9a77a23b9" // USDC / USD

	ftm[common.HexToAddress("0xae75A438b2E0cB8Bb01Ec1E1e376De11D44477CC")] = "0xccc059a1a17577676c8673952dc02070d29e5a66" // SUSHI / USD
	ftm[common.HexToAddress("0x56ee926bD8c72B2d5fa1aF4d9E4Cbb515a1E3Adc")] = "0x2eb00cc9db7a7e0a013a49b3f6ac66008d1456f7" // SNX / USD
	ftm[common.HexToAddress("0x657A1861c15A3deD9AF0B6799a195a249ebdCbc6")] = "0xd2ffccfa0934cafda647c5ff8e7918a10103c01c" // CREAM / USD
	ftm[common.HexToAddress("0x6a07A792ab2965C72a5B8088d3a069A7aC3a993B")] = "0xe6ecf7d2361b6459cbb3b4fb065e0ef4b175fe74" // AAVE / USD
	ftm[common.HexToAddress("0xb3654dc3d10ea7645f8319668e8f54d2574fbdc8")] = "0x221c773d8647bc3034e91a0c47062e26d20d97b4" // LINK / USD
	ftm[common.HexToAddress("0x27f26f00e1605903645bbabc0a73e35027dccd45")] = "0x6de70f4791c4151e00ad02e969bd900dc961f92a" // BNB / USD
	ftm[common.HexToAddress("0x8d11ec38a3eb5e956b052f67da8bdc9bef8abf3e")] = "0x91d5defaffe2854c7d02f50c80fa1fdc8a721e52" // DAI / USD
	ftm[common.HexToAddress("0x04068da6c83afcfa0e13ba15a6696662335d5b75")] = "0x2553f4eeb82d5a26427b8d1106c51499cba5d99c" // USDC / USD
	ftm[common.HexToAddress("0xe1146b9ac456fcbb60644c36fd3f868a9072fc6e")] = "0x8e94c22142f4a64b99022ccdd994f4e9ec86e4b4" // BTC / USD
	ftm[common.HexToAddress("0x321162Cd933E2Be498Cd2267a90534A804051b11")] = "0x8e94c22142f4a64b99022ccdd994f4e9ec86e4b4" // BTC / USD
	ftm[common.HexToAddress("0x658b0c7613e890ee50b8c4bc6a3f41ef411208ad")] = "0x11ddd3d147e5b83d01cee7070027092397d63658" // ETH / USD
	ftm[common.HexToAddress("0x21be370d5312f44cb42ce377bc9b8a0cef1a4c83")] = "0xf4766552d15ae4d256ad41b6cf2933482b0680dc" // FTM / USD

	optimism[common.HexToAddress("0x68f180fcce6836688e9084f035309e29bf0a2095")] = "0xd702dd976fb76fffc2d3963d037dfdae5b04e593" // BTC / USD
	optimism[common.HexToAddress("0x4200000000000000000000000000000000000006")] = "0x13e3ee699d1909e989722e753853ae30b17e08c5" // ETH / USD
	optimism[common.HexToAddress("0x7f5c764cbc14f9669b88837ca1490cca17c31607")] = "0x16a9fa2fda030272ce99b29cf780dfa30361e0f3" // USDC / USD
	optimism[common.HexToAddress("0xda10009cbd5d07dd0cecc66161fc93d7c9000da1")] = "0x8dba75e83da73cc766a7e5a0ee71f656bab470d6" // DAI / USD
	optimism[common.HexToAddress("0x94b008aa00579c1307b0ef2c499ad98a8ce58e58")] = "0xecef79e109e997bca29c1c0897ec9d7b03647f5e" // USDT / USD

	switch chain {
	case 1:
		return ethereum, nil
	case 56:
		return bsc, nil
	case 1284:
		return moonbeam, nil
	case 137:
		return polygon, nil
	case 43114:
		return avax, nil
	case 10:
		return optimism, nil
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
		"0xc907e116054ad103354f2d350fd2514433d57f6f",
		"0x49ccd9ca821efeab2b98c60dc60f518e765ede9a",
		"0x51d7180eda2260cc4f6e4eebb82fef5c3c2b8300",
		"0xebe676ee90fe1112671f19b6b7459bc678b67e8a",
		"0x2779d32d5166baaa2b2b658333ba7e6ec0c65743",
		"0x976b3d034e162d8bd72d6b9c989d545b839003b0",
		"0x0a77230d17318075983913bc2145db16c7366156",
		"0x3ca13391e9fb38a75330fb28f8cc2eb3d9ceceed",
		"0xf096872672f44d6eba71458d74fe67f9a77a23b9",
		"0xf4766552d15ae4d256ad41b6cf2933482b0680dc",
		"0x11ddd3d147e5b83d01cee7070027092397d63658",
		"0x8e94c22142f4a64b99022ccdd994f4e9ec86e4b4",
		"0x2553f4eeb82d5a26427b8d1106c51499cba5d99c",
		"0x91d5defaffe2854c7d02f50c80fa1fdc8a721e52",
		"0xf64b636c5dfe1d3555a847341cdc449f612307d0",
		"0x6de70f4791c4151e00ad02e969bd900dc961f92a",
		"0x221c773d8647bc3034e91a0c47062e26d20d97b4",
		"0xe6ecf7d2361b6459cbb3b4fb065e0ef4b175fe74",
		"0xd2ffccfa0934cafda647c5ff8e7918a10103c01c",
		"0x2eb00cc9db7a7e0a013a49b3f6ac66008d1456f7",
		"0xccc059a1a17577676c8673952dc02070d29e5a66",
		"0xd702dd976fb76fffc2d3963d037dfdae5b04e593",
		"0x13e3ee699d1909e989722e753853ae30b17e08c5",
		"0x16a9fa2fda030272ce99b29cf780dfa30361e0f3",
		"0x8dba75e83da73cc766a7e5a0ee71f656bab470d6",
		"0xecef79e109e997bca29c1c0897ec9d7b03647f5e":

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
	case 43114:
		return "0x0a77230d17318075983913bc2145db16c7366156"
	case 250:
		return "0xf4766552d15ae4d256ad41b6cf2933482b0680dc"
	case 10:
		return "0x13e3ee699d1909e989722e753853ae30b17e08c5"
	}
	return "null"
}
