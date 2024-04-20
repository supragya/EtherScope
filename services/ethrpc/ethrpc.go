package ethrpc

import (
	"math/big"

	"github.com/supragya/EtherScope/libs/service"
	itypes "github.com/supragya/EtherScope/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthRPC interface {
	service.Service

	GetTxSender(txHash, blockHash common.Hash, txIdx uint) (common.Address, error)
	GetCurrentBlockHeight() (uint64, error)
	GetBlockTimestamp(height uint64) (uint64, error)
	GetFilteredLogs(ethereum.FilterQuery) ([]types.Log, error)
	GetTokensUniV2(common.Address, *bind.CallOpts) (common.Address, common.Address, error)
	GetERC20Decimals(common.Address, *bind.CallOpts) (uint8, error)
	GetERC20Balances(requests []itypes.Tuple2[common.Address, common.Address],
		callopts *bind.CallOpts) ([]*big.Int, error)
	GetERC20Name(common.Address, *bind.CallOpts) (string, error)
	GetTokensUniV3(pairContract common.Address,
		callopts *bind.CallOpts) (common.Address, common.Address, error)
	GetTokensUniV3NFT(nftContract common.Address, tokenID *big.Int, callopts *bind.CallOpts) (common.Address, common.Address, error)
	GetChainlinkRoundData(contractAddress common.Address, callopts *bind.CallOpts) (itypes.ChainlinkLatestRoundData, error)
	IsContract(address common.Address, callopts *bind.CallOpts) (bool, error)
	GetChainlinkDecimals(contractAddress common.Address, callopts *bind.CallOpts) (uint8, error)
	GetTraderJoeTokenX(contractAddress common.Address, callopts *bind.CallOpts) (common.Address, error)
	GetTraderJoeTokenY(contractAddress common.Address, callopts *bind.CallOpts) (common.Address, error)
}
