// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package traderjoev2

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// FeeHelperFeeParameters is an auto generated low-level Go binding around an user-defined struct.
type FeeHelperFeeParameters struct {
	BinStep                  uint16
	BaseFactor               uint16
	FilterPeriod             uint16
	DecayPeriod              uint16
	ReductionFactor          uint16
	VariableFeeControl       *big.Int
	ProtocolShare            uint16
	MaxVolatilityAccumulated *big.Int
	VolatilityAccumulated    *big.Int
	VolatilityReference      *big.Int
	IndexRef                 *big.Int
	Time                     *big.Int
}

// Traderjoev2MetaData contains all meta data concerning the Traderjoev2 contract.
var Traderjoev2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractILBFactory\",\"name\":\"_factory\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"bp\",\"type\":\"uint256\"}],\"name\":\"BinHelper__BinStepOverflows\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BinHelper__IdOverflows\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__AddressZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__AddressZeroOrThis\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LBPair__CompositionFactorFlawed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__DistributionsOverflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__FlashLoanCallbackFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__FlashLoanInvalidBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__FlashLoanInvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__InsufficientAmounts\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LBPair__InsufficientLiquidityBurned\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LBPair__InsufficientLiquidityMinted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__OnlyFactory\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"LBPair__OnlyFeeRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__OnlyStrictlyIncreasingId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oracleSize\",\"type\":\"uint256\"}],\"name\":\"LBPair__OracleNewSizeTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBPair__WrongLengths\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LBToken__BurnExceedsBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBToken__BurnFromAddress0\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"accountsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idsLength\",\"type\":\"uint256\"}],\"name\":\"LBToken__LengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBToken__MintToAddress0\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"LBToken__SelfApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"LBToken__SpenderNotApproved\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LBToken__TransferExceedsBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBToken__TransferFromOrToAddress0\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LBToken__TransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"Math128x128__PowerUnderflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"prod1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"Math512Bits__MulDivOverflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"prod1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"}],\"name\":\"Math512Bits__MulShiftOverflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"}],\"name\":\"Math512Bits__OffsetOverflows\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_lookUpTimestamp\",\"type\":\"uint256\"}],\"name\":\"Oracle__LookUpTimestampTooOld\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Oracle__NotInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardUpgradeable__AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardUpgradeable__ReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"SafeCast__Exceeds112Bits\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"SafeCast__Exceeds128Bits\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"SafeCast__Exceeds24Bits\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"SafeCast__Exceeds40Bits\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenHelper__CallFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenHelper__NonContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenHelper__TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TreeMath__ErrorDepthSearch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"feesX\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"feesY\",\"type\":\"uint256\"}],\"name\":\"CompositionFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountX\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountY\",\"type\":\"uint256\"}],\"name\":\"DepositedToBin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountX\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountY\",\"type\":\"uint256\"}],\"name\":\"FeesCollected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractILBFlashLoanCallback\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"FlashLoan\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousSize\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSize\",\"type\":\"uint256\"}],\"name\":\"OracleSizeIncreased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountX\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountY\",\"type\":\"uint256\"}],\"name\":\"ProtocolFeesCollected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"swapForY\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"volatilityAccumulated\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fees\",\"type\":\"uint256\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountX\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountY\",\"type\":\"uint256\"}],\"name\":\"WithdrawnFromBin\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"batchBalances\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountY\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"}],\"name\":\"collectFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountY\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collectProtocolFees\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"amountX\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amountY\",\"type\":\"uint128\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"contractILBFactory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeParameters\",\"outputs\":[{\"components\":[{\"internalType\":\"uint16\",\"name\":\"binStep\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"baseFactor\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"filterPeriod\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"decayPeriod\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"reductionFactor\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"variableFeeControl\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"protocolShare\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"maxVolatilityAccumulated\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"volatilityAccumulated\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"volatilityReference\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"indexRef\",\"type\":\"uint24\"},{\"internalType\":\"uint40\",\"name\":\"time\",\"type\":\"uint40\"}],\"internalType\":\"structFeeHelper.FeeParameters\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"_id\",\"type\":\"uint24\"},{\"internalType\":\"bool\",\"name\":\"_swapForY\",\"type\":\"bool\"}],\"name\":\"findFirstNonEmptyBinId\",\"outputs\":[{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractILBFlashLoanCallback\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"flashLoan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceDecay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"_id\",\"type\":\"uint24\"}],\"name\":\"getBin\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"reserveX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveY\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGlobalFees\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"feesXTotal\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"feesYTotal\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"feesXProtocol\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"feesYProtocol\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracleParameters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"oracleSampleLifetime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oracleSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oracleActiveSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oracleLastTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oracleId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"min\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_timeDelta\",\"type\":\"uint256\"}],\"name\":\"getOracleSampleFrom\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"cumulativeId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cumulativeVolatilityAccumulated\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cumulativeBinCrossed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReservesAndId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"reserveX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"activeId\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_newLength\",\"type\":\"uint16\"}],\"name\":\"increaseOracleLength\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"_tokenX\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_tokenY\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"_activeId\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"_sampleLifetime\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_packedFeeParameters\",\"type\":\"bytes32\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_distributionX\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_distributionY\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"liquidityMinted\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"}],\"name\":\"pendingFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountY\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_packedFeeParameters\",\"type\":\"bytes32\"}],\"name\":\"setFeesParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_swapForY\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"swap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountXOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountYOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenX\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenY\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Traderjoev2ABI is the input ABI used to generate the binding from.
// Deprecated: Use Traderjoev2MetaData.ABI instead.
var Traderjoev2ABI = Traderjoev2MetaData.ABI

// Traderjoev2 is an auto generated Go binding around an Ethereum contract.
type Traderjoev2 struct {
	Traderjoev2Caller     // Read-only binding to the contract
	Traderjoev2Transactor // Write-only binding to the contract
	Traderjoev2Filterer   // Log filterer for contract events
}

// Traderjoev2Caller is an auto generated read-only Go binding around an Ethereum contract.
type Traderjoev2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Traderjoev2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Traderjoev2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Traderjoev2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Traderjoev2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Traderjoev2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Traderjoev2Session struct {
	Contract     *Traderjoev2      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Traderjoev2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Traderjoev2CallerSession struct {
	Contract *Traderjoev2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// Traderjoev2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Traderjoev2TransactorSession struct {
	Contract     *Traderjoev2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// Traderjoev2Raw is an auto generated low-level Go binding around an Ethereum contract.
type Traderjoev2Raw struct {
	Contract *Traderjoev2 // Generic contract binding to access the raw methods on
}

// Traderjoev2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Traderjoev2CallerRaw struct {
	Contract *Traderjoev2Caller // Generic read-only contract binding to access the raw methods on
}

// Traderjoev2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Traderjoev2TransactorRaw struct {
	Contract *Traderjoev2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTraderjoev2 creates a new instance of Traderjoev2, bound to a specific deployed contract.
func NewTraderjoev2(address common.Address, backend bind.ContractBackend) (*Traderjoev2, error) {
	contract, err := bindTraderjoev2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2{Traderjoev2Caller: Traderjoev2Caller{contract: contract}, Traderjoev2Transactor: Traderjoev2Transactor{contract: contract}, Traderjoev2Filterer: Traderjoev2Filterer{contract: contract}}, nil
}

// NewTraderjoev2Caller creates a new read-only instance of Traderjoev2, bound to a specific deployed contract.
func NewTraderjoev2Caller(address common.Address, caller bind.ContractCaller) (*Traderjoev2Caller, error) {
	contract, err := bindTraderjoev2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2Caller{contract: contract}, nil
}

// NewTraderjoev2Transactor creates a new write-only instance of Traderjoev2, bound to a specific deployed contract.
func NewTraderjoev2Transactor(address common.Address, transactor bind.ContractTransactor) (*Traderjoev2Transactor, error) {
	contract, err := bindTraderjoev2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2Transactor{contract: contract}, nil
}

// NewTraderjoev2Filterer creates a new log filterer instance of Traderjoev2, bound to a specific deployed contract.
func NewTraderjoev2Filterer(address common.Address, filterer bind.ContractFilterer) (*Traderjoev2Filterer, error) {
	contract, err := bindTraderjoev2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2Filterer{contract: contract}, nil
}

// bindTraderjoev2 binds a generic wrapper to an already deployed contract.
func bindTraderjoev2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Traderjoev2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Traderjoev2 *Traderjoev2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Traderjoev2.Contract.Traderjoev2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Traderjoev2 *Traderjoev2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Traderjoev2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Traderjoev2 *Traderjoev2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Traderjoev2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Traderjoev2 *Traderjoev2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Traderjoev2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Traderjoev2 *Traderjoev2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Traderjoev2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Traderjoev2 *Traderjoev2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Traderjoev2.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address _account, uint256 _id) view returns(uint256)
func (_Traderjoev2 *Traderjoev2Caller) BalanceOf(opts *bind.CallOpts, _account common.Address, _id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "balanceOf", _account, _id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address _account, uint256 _id) view returns(uint256)
func (_Traderjoev2 *Traderjoev2Session) BalanceOf(_account common.Address, _id *big.Int) (*big.Int, error) {
	return _Traderjoev2.Contract.BalanceOf(&_Traderjoev2.CallOpts, _account, _id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address _account, uint256 _id) view returns(uint256)
func (_Traderjoev2 *Traderjoev2CallerSession) BalanceOf(_account common.Address, _id *big.Int) (*big.Int, error) {
	return _Traderjoev2.Contract.BalanceOf(&_Traderjoev2.CallOpts, _account, _id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] _accounts, uint256[] _ids) view returns(uint256[] batchBalances)
func (_Traderjoev2 *Traderjoev2Caller) BalanceOfBatch(opts *bind.CallOpts, _accounts []common.Address, _ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "balanceOfBatch", _accounts, _ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] _accounts, uint256[] _ids) view returns(uint256[] batchBalances)
func (_Traderjoev2 *Traderjoev2Session) BalanceOfBatch(_accounts []common.Address, _ids []*big.Int) ([]*big.Int, error) {
	return _Traderjoev2.Contract.BalanceOfBatch(&_Traderjoev2.CallOpts, _accounts, _ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] _accounts, uint256[] _ids) view returns(uint256[] batchBalances)
func (_Traderjoev2 *Traderjoev2CallerSession) BalanceOfBatch(_accounts []common.Address, _ids []*big.Int) ([]*big.Int, error) {
	return _Traderjoev2.Contract.BalanceOfBatch(&_Traderjoev2.CallOpts, _accounts, _ids)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_Traderjoev2 *Traderjoev2Caller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_Traderjoev2 *Traderjoev2Session) Factory() (common.Address, error) {
	return _Traderjoev2.Contract.Factory(&_Traderjoev2.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_Traderjoev2 *Traderjoev2CallerSession) Factory() (common.Address, error) {
	return _Traderjoev2.Contract.Factory(&_Traderjoev2.CallOpts)
}

// FeeParameters is a free data retrieval call binding the contract method 0x98c7adf3.
//
// Solidity: function feeParameters() view returns((uint16,uint16,uint16,uint16,uint16,uint24,uint16,uint24,uint24,uint24,uint24,uint40))
func (_Traderjoev2 *Traderjoev2Caller) FeeParameters(opts *bind.CallOpts) (FeeHelperFeeParameters, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "feeParameters")

	if err != nil {
		return *new(FeeHelperFeeParameters), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeHelperFeeParameters)).(*FeeHelperFeeParameters)

	return out0, err

}

// FeeParameters is a free data retrieval call binding the contract method 0x98c7adf3.
//
// Solidity: function feeParameters() view returns((uint16,uint16,uint16,uint16,uint16,uint24,uint16,uint24,uint24,uint24,uint24,uint40))
func (_Traderjoev2 *Traderjoev2Session) FeeParameters() (FeeHelperFeeParameters, error) {
	return _Traderjoev2.Contract.FeeParameters(&_Traderjoev2.CallOpts)
}

// FeeParameters is a free data retrieval call binding the contract method 0x98c7adf3.
//
// Solidity: function feeParameters() view returns((uint16,uint16,uint16,uint16,uint16,uint24,uint16,uint24,uint24,uint24,uint24,uint40))
func (_Traderjoev2 *Traderjoev2CallerSession) FeeParameters() (FeeHelperFeeParameters, error) {
	return _Traderjoev2.Contract.FeeParameters(&_Traderjoev2.CallOpts)
}

// FindFirstNonEmptyBinId is a free data retrieval call binding the contract method 0x8f919a83.
//
// Solidity: function findFirstNonEmptyBinId(uint24 _id, bool _swapForY) view returns(uint24)
func (_Traderjoev2 *Traderjoev2Caller) FindFirstNonEmptyBinId(opts *bind.CallOpts, _id *big.Int, _swapForY bool) (*big.Int, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "findFirstNonEmptyBinId", _id, _swapForY)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FindFirstNonEmptyBinId is a free data retrieval call binding the contract method 0x8f919a83.
//
// Solidity: function findFirstNonEmptyBinId(uint24 _id, bool _swapForY) view returns(uint24)
func (_Traderjoev2 *Traderjoev2Session) FindFirstNonEmptyBinId(_id *big.Int, _swapForY bool) (*big.Int, error) {
	return _Traderjoev2.Contract.FindFirstNonEmptyBinId(&_Traderjoev2.CallOpts, _id, _swapForY)
}

// FindFirstNonEmptyBinId is a free data retrieval call binding the contract method 0x8f919a83.
//
// Solidity: function findFirstNonEmptyBinId(uint24 _id, bool _swapForY) view returns(uint24)
func (_Traderjoev2 *Traderjoev2CallerSession) FindFirstNonEmptyBinId(_id *big.Int, _swapForY bool) (*big.Int, error) {
	return _Traderjoev2.Contract.FindFirstNonEmptyBinId(&_Traderjoev2.CallOpts, _id, _swapForY)
}

// GetBin is a free data retrieval call binding the contract method 0x0abe9688.
//
// Solidity: function getBin(uint24 _id) view returns(uint256 reserveX, uint256 reserveY)
func (_Traderjoev2 *Traderjoev2Caller) GetBin(opts *bind.CallOpts, _id *big.Int) (struct {
	ReserveX *big.Int
	ReserveY *big.Int
}, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "getBin", _id)

	outstruct := new(struct {
		ReserveX *big.Int
		ReserveY *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ReserveX = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ReserveY = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetBin is a free data retrieval call binding the contract method 0x0abe9688.
//
// Solidity: function getBin(uint24 _id) view returns(uint256 reserveX, uint256 reserveY)
func (_Traderjoev2 *Traderjoev2Session) GetBin(_id *big.Int) (struct {
	ReserveX *big.Int
	ReserveY *big.Int
}, error) {
	return _Traderjoev2.Contract.GetBin(&_Traderjoev2.CallOpts, _id)
}

// GetBin is a free data retrieval call binding the contract method 0x0abe9688.
//
// Solidity: function getBin(uint24 _id) view returns(uint256 reserveX, uint256 reserveY)
func (_Traderjoev2 *Traderjoev2CallerSession) GetBin(_id *big.Int) (struct {
	ReserveX *big.Int
	ReserveY *big.Int
}, error) {
	return _Traderjoev2.Contract.GetBin(&_Traderjoev2.CallOpts, _id)
}

// GetGlobalFees is a free data retrieval call binding the contract method 0xa582cdaa.
//
// Solidity: function getGlobalFees() view returns(uint128 feesXTotal, uint128 feesYTotal, uint128 feesXProtocol, uint128 feesYProtocol)
func (_Traderjoev2 *Traderjoev2Caller) GetGlobalFees(opts *bind.CallOpts) (struct {
	FeesXTotal    *big.Int
	FeesYTotal    *big.Int
	FeesXProtocol *big.Int
	FeesYProtocol *big.Int
}, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "getGlobalFees")

	outstruct := new(struct {
		FeesXTotal    *big.Int
		FeesYTotal    *big.Int
		FeesXProtocol *big.Int
		FeesYProtocol *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FeesXTotal = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeesYTotal = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FeesXProtocol = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FeesYProtocol = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetGlobalFees is a free data retrieval call binding the contract method 0xa582cdaa.
//
// Solidity: function getGlobalFees() view returns(uint128 feesXTotal, uint128 feesYTotal, uint128 feesXProtocol, uint128 feesYProtocol)
func (_Traderjoev2 *Traderjoev2Session) GetGlobalFees() (struct {
	FeesXTotal    *big.Int
	FeesYTotal    *big.Int
	FeesXProtocol *big.Int
	FeesYProtocol *big.Int
}, error) {
	return _Traderjoev2.Contract.GetGlobalFees(&_Traderjoev2.CallOpts)
}

// GetGlobalFees is a free data retrieval call binding the contract method 0xa582cdaa.
//
// Solidity: function getGlobalFees() view returns(uint128 feesXTotal, uint128 feesYTotal, uint128 feesXProtocol, uint128 feesYProtocol)
func (_Traderjoev2 *Traderjoev2CallerSession) GetGlobalFees() (struct {
	FeesXTotal    *big.Int
	FeesYTotal    *big.Int
	FeesXProtocol *big.Int
	FeesYProtocol *big.Int
}, error) {
	return _Traderjoev2.Contract.GetGlobalFees(&_Traderjoev2.CallOpts)
}

// GetOracleParameters is a free data retrieval call binding the contract method 0x55182894.
//
// Solidity: function getOracleParameters() view returns(uint256 oracleSampleLifetime, uint256 oracleSize, uint256 oracleActiveSize, uint256 oracleLastTimestamp, uint256 oracleId, uint256 min, uint256 max)
func (_Traderjoev2 *Traderjoev2Caller) GetOracleParameters(opts *bind.CallOpts) (struct {
	OracleSampleLifetime *big.Int
	OracleSize           *big.Int
	OracleActiveSize     *big.Int
	OracleLastTimestamp  *big.Int
	OracleId             *big.Int
	Min                  *big.Int
	Max                  *big.Int
}, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "getOracleParameters")

	outstruct := new(struct {
		OracleSampleLifetime *big.Int
		OracleSize           *big.Int
		OracleActiveSize     *big.Int
		OracleLastTimestamp  *big.Int
		OracleId             *big.Int
		Min                  *big.Int
		Max                  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.OracleSampleLifetime = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.OracleSize = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.OracleActiveSize = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.OracleLastTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.OracleId = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Min = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Max = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetOracleParameters is a free data retrieval call binding the contract method 0x55182894.
//
// Solidity: function getOracleParameters() view returns(uint256 oracleSampleLifetime, uint256 oracleSize, uint256 oracleActiveSize, uint256 oracleLastTimestamp, uint256 oracleId, uint256 min, uint256 max)
func (_Traderjoev2 *Traderjoev2Session) GetOracleParameters() (struct {
	OracleSampleLifetime *big.Int
	OracleSize           *big.Int
	OracleActiveSize     *big.Int
	OracleLastTimestamp  *big.Int
	OracleId             *big.Int
	Min                  *big.Int
	Max                  *big.Int
}, error) {
	return _Traderjoev2.Contract.GetOracleParameters(&_Traderjoev2.CallOpts)
}

// GetOracleParameters is a free data retrieval call binding the contract method 0x55182894.
//
// Solidity: function getOracleParameters() view returns(uint256 oracleSampleLifetime, uint256 oracleSize, uint256 oracleActiveSize, uint256 oracleLastTimestamp, uint256 oracleId, uint256 min, uint256 max)
func (_Traderjoev2 *Traderjoev2CallerSession) GetOracleParameters() (struct {
	OracleSampleLifetime *big.Int
	OracleSize           *big.Int
	OracleActiveSize     *big.Int
	OracleLastTimestamp  *big.Int
	OracleId             *big.Int
	Min                  *big.Int
	Max                  *big.Int
}, error) {
	return _Traderjoev2.Contract.GetOracleParameters(&_Traderjoev2.CallOpts)
}

// GetOracleSampleFrom is a free data retrieval call binding the contract method 0xa21635a7.
//
// Solidity: function getOracleSampleFrom(uint256 _timeDelta) view returns(uint256 cumulativeId, uint256 cumulativeVolatilityAccumulated, uint256 cumulativeBinCrossed)
func (_Traderjoev2 *Traderjoev2Caller) GetOracleSampleFrom(opts *bind.CallOpts, _timeDelta *big.Int) (struct {
	CumulativeId                    *big.Int
	CumulativeVolatilityAccumulated *big.Int
	CumulativeBinCrossed            *big.Int
}, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "getOracleSampleFrom", _timeDelta)

	outstruct := new(struct {
		CumulativeId                    *big.Int
		CumulativeVolatilityAccumulated *big.Int
		CumulativeBinCrossed            *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CumulativeId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.CumulativeVolatilityAccumulated = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.CumulativeBinCrossed = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetOracleSampleFrom is a free data retrieval call binding the contract method 0xa21635a7.
//
// Solidity: function getOracleSampleFrom(uint256 _timeDelta) view returns(uint256 cumulativeId, uint256 cumulativeVolatilityAccumulated, uint256 cumulativeBinCrossed)
func (_Traderjoev2 *Traderjoev2Session) GetOracleSampleFrom(_timeDelta *big.Int) (struct {
	CumulativeId                    *big.Int
	CumulativeVolatilityAccumulated *big.Int
	CumulativeBinCrossed            *big.Int
}, error) {
	return _Traderjoev2.Contract.GetOracleSampleFrom(&_Traderjoev2.CallOpts, _timeDelta)
}

// GetOracleSampleFrom is a free data retrieval call binding the contract method 0xa21635a7.
//
// Solidity: function getOracleSampleFrom(uint256 _timeDelta) view returns(uint256 cumulativeId, uint256 cumulativeVolatilityAccumulated, uint256 cumulativeBinCrossed)
func (_Traderjoev2 *Traderjoev2CallerSession) GetOracleSampleFrom(_timeDelta *big.Int) (struct {
	CumulativeId                    *big.Int
	CumulativeVolatilityAccumulated *big.Int
	CumulativeBinCrossed            *big.Int
}, error) {
	return _Traderjoev2.Contract.GetOracleSampleFrom(&_Traderjoev2.CallOpts, _timeDelta)
}

// GetReservesAndId is a free data retrieval call binding the contract method 0x1b05b83e.
//
// Solidity: function getReservesAndId() view returns(uint256 reserveX, uint256 reserveY, uint256 activeId)
func (_Traderjoev2 *Traderjoev2Caller) GetReservesAndId(opts *bind.CallOpts) (struct {
	ReserveX *big.Int
	ReserveY *big.Int
	ActiveId *big.Int
}, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "getReservesAndId")

	outstruct := new(struct {
		ReserveX *big.Int
		ReserveY *big.Int
		ActiveId *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ReserveX = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ReserveY = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ActiveId = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetReservesAndId is a free data retrieval call binding the contract method 0x1b05b83e.
//
// Solidity: function getReservesAndId() view returns(uint256 reserveX, uint256 reserveY, uint256 activeId)
func (_Traderjoev2 *Traderjoev2Session) GetReservesAndId() (struct {
	ReserveX *big.Int
	ReserveY *big.Int
	ActiveId *big.Int
}, error) {
	return _Traderjoev2.Contract.GetReservesAndId(&_Traderjoev2.CallOpts)
}

// GetReservesAndId is a free data retrieval call binding the contract method 0x1b05b83e.
//
// Solidity: function getReservesAndId() view returns(uint256 reserveX, uint256 reserveY, uint256 activeId)
func (_Traderjoev2 *Traderjoev2CallerSession) GetReservesAndId() (struct {
	ReserveX *big.Int
	ReserveY *big.Int
	ActiveId *big.Int
}, error) {
	return _Traderjoev2.Contract.GetReservesAndId(&_Traderjoev2.CallOpts)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address _owner, address _spender) view returns(bool)
func (_Traderjoev2 *Traderjoev2Caller) IsApprovedForAll(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (bool, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "isApprovedForAll", _owner, _spender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address _owner, address _spender) view returns(bool)
func (_Traderjoev2 *Traderjoev2Session) IsApprovedForAll(_owner common.Address, _spender common.Address) (bool, error) {
	return _Traderjoev2.Contract.IsApprovedForAll(&_Traderjoev2.CallOpts, _owner, _spender)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address _owner, address _spender) view returns(bool)
func (_Traderjoev2 *Traderjoev2CallerSession) IsApprovedForAll(_owner common.Address, _spender common.Address) (bool, error) {
	return _Traderjoev2.Contract.IsApprovedForAll(&_Traderjoev2.CallOpts, _owner, _spender)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Traderjoev2 *Traderjoev2Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Traderjoev2 *Traderjoev2Session) Name() (string, error) {
	return _Traderjoev2.Contract.Name(&_Traderjoev2.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Traderjoev2 *Traderjoev2CallerSession) Name() (string, error) {
	return _Traderjoev2.Contract.Name(&_Traderjoev2.CallOpts)
}

// PendingFees is a free data retrieval call binding the contract method 0xf7cff1f8.
//
// Solidity: function pendingFees(address _account, uint256[] _ids) view returns(uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Caller) PendingFees(opts *bind.CallOpts, _account common.Address, _ids []*big.Int) (struct {
	AmountX *big.Int
	AmountY *big.Int
}, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "pendingFees", _account, _ids)

	outstruct := new(struct {
		AmountX *big.Int
		AmountY *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AmountX = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.AmountY = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PendingFees is a free data retrieval call binding the contract method 0xf7cff1f8.
//
// Solidity: function pendingFees(address _account, uint256[] _ids) view returns(uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Session) PendingFees(_account common.Address, _ids []*big.Int) (struct {
	AmountX *big.Int
	AmountY *big.Int
}, error) {
	return _Traderjoev2.Contract.PendingFees(&_Traderjoev2.CallOpts, _account, _ids)
}

// PendingFees is a free data retrieval call binding the contract method 0xf7cff1f8.
//
// Solidity: function pendingFees(address _account, uint256[] _ids) view returns(uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2CallerSession) PendingFees(_account common.Address, _ids []*big.Int) (struct {
	AmountX *big.Int
	AmountY *big.Int
}, error) {
	return _Traderjoev2.Contract.PendingFees(&_Traderjoev2.CallOpts, _account, _ids)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceId) view returns(bool)
func (_Traderjoev2 *Traderjoev2Caller) SupportsInterface(opts *bind.CallOpts, _interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "supportsInterface", _interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceId) view returns(bool)
func (_Traderjoev2 *Traderjoev2Session) SupportsInterface(_interfaceId [4]byte) (bool, error) {
	return _Traderjoev2.Contract.SupportsInterface(&_Traderjoev2.CallOpts, _interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceId) view returns(bool)
func (_Traderjoev2 *Traderjoev2CallerSession) SupportsInterface(_interfaceId [4]byte) (bool, error) {
	return _Traderjoev2.Contract.SupportsInterface(&_Traderjoev2.CallOpts, _interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Traderjoev2 *Traderjoev2Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Traderjoev2 *Traderjoev2Session) Symbol() (string, error) {
	return _Traderjoev2.Contract.Symbol(&_Traderjoev2.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Traderjoev2 *Traderjoev2CallerSession) Symbol() (string, error) {
	return _Traderjoev2.Contract.Symbol(&_Traderjoev2.CallOpts)
}

// TokenX is a free data retrieval call binding the contract method 0x16dc165b.
//
// Solidity: function tokenX() view returns(address)
func (_Traderjoev2 *Traderjoev2Caller) TokenX(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "tokenX")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenX is a free data retrieval call binding the contract method 0x16dc165b.
//
// Solidity: function tokenX() view returns(address)
func (_Traderjoev2 *Traderjoev2Session) TokenX() (common.Address, error) {
	return _Traderjoev2.Contract.TokenX(&_Traderjoev2.CallOpts)
}

// TokenX is a free data retrieval call binding the contract method 0x16dc165b.
//
// Solidity: function tokenX() view returns(address)
func (_Traderjoev2 *Traderjoev2CallerSession) TokenX() (common.Address, error) {
	return _Traderjoev2.Contract.TokenX(&_Traderjoev2.CallOpts)
}

// TokenY is a free data retrieval call binding the contract method 0xb7d19fc4.
//
// Solidity: function tokenY() view returns(address)
func (_Traderjoev2 *Traderjoev2Caller) TokenY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "tokenY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenY is a free data retrieval call binding the contract method 0xb7d19fc4.
//
// Solidity: function tokenY() view returns(address)
func (_Traderjoev2 *Traderjoev2Session) TokenY() (common.Address, error) {
	return _Traderjoev2.Contract.TokenY(&_Traderjoev2.CallOpts)
}

// TokenY is a free data retrieval call binding the contract method 0xb7d19fc4.
//
// Solidity: function tokenY() view returns(address)
func (_Traderjoev2 *Traderjoev2CallerSession) TokenY() (common.Address, error) {
	return _Traderjoev2.Contract.TokenY(&_Traderjoev2.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 _id) view returns(uint256)
func (_Traderjoev2 *Traderjoev2Caller) TotalSupply(opts *bind.CallOpts, _id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Traderjoev2.contract.Call(opts, &out, "totalSupply", _id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 _id) view returns(uint256)
func (_Traderjoev2 *Traderjoev2Session) TotalSupply(_id *big.Int) (*big.Int, error) {
	return _Traderjoev2.Contract.TotalSupply(&_Traderjoev2.CallOpts, _id)
}

// TotalSupply is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 _id) view returns(uint256)
func (_Traderjoev2 *Traderjoev2CallerSession) TotalSupply(_id *big.Int) (*big.Int, error) {
	return _Traderjoev2.Contract.TotalSupply(&_Traderjoev2.CallOpts, _id)
}

// Burn is a paid mutator transaction binding the contract method 0x0acd451d.
//
// Solidity: function burn(uint256[] _ids, uint256[] _amounts, address _to) returns(uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Transactor) Burn(opts *bind.TransactOpts, _ids []*big.Int, _amounts []*big.Int, _to common.Address) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "burn", _ids, _amounts, _to)
}

// Burn is a paid mutator transaction binding the contract method 0x0acd451d.
//
// Solidity: function burn(uint256[] _ids, uint256[] _amounts, address _to) returns(uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Session) Burn(_ids []*big.Int, _amounts []*big.Int, _to common.Address) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Burn(&_Traderjoev2.TransactOpts, _ids, _amounts, _to)
}

// Burn is a paid mutator transaction binding the contract method 0x0acd451d.
//
// Solidity: function burn(uint256[] _ids, uint256[] _amounts, address _to) returns(uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2TransactorSession) Burn(_ids []*big.Int, _amounts []*big.Int, _to common.Address) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Burn(&_Traderjoev2.TransactOpts, _ids, _amounts, _to)
}

// CollectFees is a paid mutator transaction binding the contract method 0x225b20b9.
//
// Solidity: function collectFees(address _account, uint256[] _ids) returns(uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Transactor) CollectFees(opts *bind.TransactOpts, _account common.Address, _ids []*big.Int) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "collectFees", _account, _ids)
}

// CollectFees is a paid mutator transaction binding the contract method 0x225b20b9.
//
// Solidity: function collectFees(address _account, uint256[] _ids) returns(uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Session) CollectFees(_account common.Address, _ids []*big.Int) (*types.Transaction, error) {
	return _Traderjoev2.Contract.CollectFees(&_Traderjoev2.TransactOpts, _account, _ids)
}

// CollectFees is a paid mutator transaction binding the contract method 0x225b20b9.
//
// Solidity: function collectFees(address _account, uint256[] _ids) returns(uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2TransactorSession) CollectFees(_account common.Address, _ids []*big.Int) (*types.Transaction, error) {
	return _Traderjoev2.Contract.CollectFees(&_Traderjoev2.TransactOpts, _account, _ids)
}

// CollectProtocolFees is a paid mutator transaction binding the contract method 0xa1af5b9a.
//
// Solidity: function collectProtocolFees() returns(uint128 amountX, uint128 amountY)
func (_Traderjoev2 *Traderjoev2Transactor) CollectProtocolFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "collectProtocolFees")
}

// CollectProtocolFees is a paid mutator transaction binding the contract method 0xa1af5b9a.
//
// Solidity: function collectProtocolFees() returns(uint128 amountX, uint128 amountY)
func (_Traderjoev2 *Traderjoev2Session) CollectProtocolFees() (*types.Transaction, error) {
	return _Traderjoev2.Contract.CollectProtocolFees(&_Traderjoev2.TransactOpts)
}

// CollectProtocolFees is a paid mutator transaction binding the contract method 0xa1af5b9a.
//
// Solidity: function collectProtocolFees() returns(uint128 amountX, uint128 amountY)
func (_Traderjoev2 *Traderjoev2TransactorSession) CollectProtocolFees() (*types.Transaction, error) {
	return _Traderjoev2.Contract.CollectProtocolFees(&_Traderjoev2.TransactOpts)
}

// FlashLoan is a paid mutator transaction binding the contract method 0x5cffe9de.
//
// Solidity: function flashLoan(address _receiver, address _token, uint256 _amount, bytes _data) returns()
func (_Traderjoev2 *Traderjoev2Transactor) FlashLoan(opts *bind.TransactOpts, _receiver common.Address, _token common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "flashLoan", _receiver, _token, _amount, _data)
}

// FlashLoan is a paid mutator transaction binding the contract method 0x5cffe9de.
//
// Solidity: function flashLoan(address _receiver, address _token, uint256 _amount, bytes _data) returns()
func (_Traderjoev2 *Traderjoev2Session) FlashLoan(_receiver common.Address, _token common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Traderjoev2.Contract.FlashLoan(&_Traderjoev2.TransactOpts, _receiver, _token, _amount, _data)
}

// FlashLoan is a paid mutator transaction binding the contract method 0x5cffe9de.
//
// Solidity: function flashLoan(address _receiver, address _token, uint256 _amount, bytes _data) returns()
func (_Traderjoev2 *Traderjoev2TransactorSession) FlashLoan(_receiver common.Address, _token common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Traderjoev2.Contract.FlashLoan(&_Traderjoev2.TransactOpts, _receiver, _token, _amount, _data)
}

// ForceDecay is a paid mutator transaction binding the contract method 0xd3b9fbe4.
//
// Solidity: function forceDecay() returns()
func (_Traderjoev2 *Traderjoev2Transactor) ForceDecay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "forceDecay")
}

// ForceDecay is a paid mutator transaction binding the contract method 0xd3b9fbe4.
//
// Solidity: function forceDecay() returns()
func (_Traderjoev2 *Traderjoev2Session) ForceDecay() (*types.Transaction, error) {
	return _Traderjoev2.Contract.ForceDecay(&_Traderjoev2.TransactOpts)
}

// ForceDecay is a paid mutator transaction binding the contract method 0xd3b9fbe4.
//
// Solidity: function forceDecay() returns()
func (_Traderjoev2 *Traderjoev2TransactorSession) ForceDecay() (*types.Transaction, error) {
	return _Traderjoev2.Contract.ForceDecay(&_Traderjoev2.TransactOpts)
}

// IncreaseOracleLength is a paid mutator transaction binding the contract method 0xc7bd6586.
//
// Solidity: function increaseOracleLength(uint16 _newLength) returns()
func (_Traderjoev2 *Traderjoev2Transactor) IncreaseOracleLength(opts *bind.TransactOpts, _newLength uint16) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "increaseOracleLength", _newLength)
}

// IncreaseOracleLength is a paid mutator transaction binding the contract method 0xc7bd6586.
//
// Solidity: function increaseOracleLength(uint16 _newLength) returns()
func (_Traderjoev2 *Traderjoev2Session) IncreaseOracleLength(_newLength uint16) (*types.Transaction, error) {
	return _Traderjoev2.Contract.IncreaseOracleLength(&_Traderjoev2.TransactOpts, _newLength)
}

// IncreaseOracleLength is a paid mutator transaction binding the contract method 0xc7bd6586.
//
// Solidity: function increaseOracleLength(uint16 _newLength) returns()
func (_Traderjoev2 *Traderjoev2TransactorSession) IncreaseOracleLength(_newLength uint16) (*types.Transaction, error) {
	return _Traderjoev2.Contract.IncreaseOracleLength(&_Traderjoev2.TransactOpts, _newLength)
}

// Initialize is a paid mutator transaction binding the contract method 0xd32db437.
//
// Solidity: function initialize(address _tokenX, address _tokenY, uint24 _activeId, uint16 _sampleLifetime, bytes32 _packedFeeParameters) returns()
func (_Traderjoev2 *Traderjoev2Transactor) Initialize(opts *bind.TransactOpts, _tokenX common.Address, _tokenY common.Address, _activeId *big.Int, _sampleLifetime uint16, _packedFeeParameters [32]byte) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "initialize", _tokenX, _tokenY, _activeId, _sampleLifetime, _packedFeeParameters)
}

// Initialize is a paid mutator transaction binding the contract method 0xd32db437.
//
// Solidity: function initialize(address _tokenX, address _tokenY, uint24 _activeId, uint16 _sampleLifetime, bytes32 _packedFeeParameters) returns()
func (_Traderjoev2 *Traderjoev2Session) Initialize(_tokenX common.Address, _tokenY common.Address, _activeId *big.Int, _sampleLifetime uint16, _packedFeeParameters [32]byte) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Initialize(&_Traderjoev2.TransactOpts, _tokenX, _tokenY, _activeId, _sampleLifetime, _packedFeeParameters)
}

// Initialize is a paid mutator transaction binding the contract method 0xd32db437.
//
// Solidity: function initialize(address _tokenX, address _tokenY, uint24 _activeId, uint16 _sampleLifetime, bytes32 _packedFeeParameters) returns()
func (_Traderjoev2 *Traderjoev2TransactorSession) Initialize(_tokenX common.Address, _tokenY common.Address, _activeId *big.Int, _sampleLifetime uint16, _packedFeeParameters [32]byte) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Initialize(&_Traderjoev2.TransactOpts, _tokenX, _tokenY, _activeId, _sampleLifetime, _packedFeeParameters)
}

// Mint is a paid mutator transaction binding the contract method 0x714c8592.
//
// Solidity: function mint(uint256[] _ids, uint256[] _distributionX, uint256[] _distributionY, address _to) returns(uint256, uint256, uint256[] liquidityMinted)
func (_Traderjoev2 *Traderjoev2Transactor) Mint(opts *bind.TransactOpts, _ids []*big.Int, _distributionX []*big.Int, _distributionY []*big.Int, _to common.Address) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "mint", _ids, _distributionX, _distributionY, _to)
}

// Mint is a paid mutator transaction binding the contract method 0x714c8592.
//
// Solidity: function mint(uint256[] _ids, uint256[] _distributionX, uint256[] _distributionY, address _to) returns(uint256, uint256, uint256[] liquidityMinted)
func (_Traderjoev2 *Traderjoev2Session) Mint(_ids []*big.Int, _distributionX []*big.Int, _distributionY []*big.Int, _to common.Address) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Mint(&_Traderjoev2.TransactOpts, _ids, _distributionX, _distributionY, _to)
}

// Mint is a paid mutator transaction binding the contract method 0x714c8592.
//
// Solidity: function mint(uint256[] _ids, uint256[] _distributionX, uint256[] _distributionY, address _to) returns(uint256, uint256, uint256[] liquidityMinted)
func (_Traderjoev2 *Traderjoev2TransactorSession) Mint(_ids []*big.Int, _distributionX []*big.Int, _distributionY []*big.Int, _to common.Address) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Mint(&_Traderjoev2.TransactOpts, _ids, _distributionX, _distributionY, _to)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0xfba0ee64.
//
// Solidity: function safeBatchTransferFrom(address _from, address _to, uint256[] _ids, uint256[] _amounts) returns()
func (_Traderjoev2 *Traderjoev2Transactor) SafeBatchTransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "safeBatchTransferFrom", _from, _to, _ids, _amounts)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0xfba0ee64.
//
// Solidity: function safeBatchTransferFrom(address _from, address _to, uint256[] _ids, uint256[] _amounts) returns()
func (_Traderjoev2 *Traderjoev2Session) SafeBatchTransferFrom(_from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int) (*types.Transaction, error) {
	return _Traderjoev2.Contract.SafeBatchTransferFrom(&_Traderjoev2.TransactOpts, _from, _to, _ids, _amounts)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0xfba0ee64.
//
// Solidity: function safeBatchTransferFrom(address _from, address _to, uint256[] _ids, uint256[] _amounts) returns()
func (_Traderjoev2 *Traderjoev2TransactorSession) SafeBatchTransferFrom(_from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int) (*types.Transaction, error) {
	return _Traderjoev2.Contract.SafeBatchTransferFrom(&_Traderjoev2.TransactOpts, _from, _to, _ids, _amounts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x0febdd49.
//
// Solidity: function safeTransferFrom(address _from, address _to, uint256 _id, uint256 _amount) returns()
func (_Traderjoev2 *Traderjoev2Transactor) SafeTransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "safeTransferFrom", _from, _to, _id, _amount)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x0febdd49.
//
// Solidity: function safeTransferFrom(address _from, address _to, uint256 _id, uint256 _amount) returns()
func (_Traderjoev2 *Traderjoev2Session) SafeTransferFrom(_from common.Address, _to common.Address, _id *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _Traderjoev2.Contract.SafeTransferFrom(&_Traderjoev2.TransactOpts, _from, _to, _id, _amount)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x0febdd49.
//
// Solidity: function safeTransferFrom(address _from, address _to, uint256 _id, uint256 _amount) returns()
func (_Traderjoev2 *Traderjoev2TransactorSession) SafeTransferFrom(_from common.Address, _to common.Address, _id *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _Traderjoev2.Contract.SafeTransferFrom(&_Traderjoev2.TransactOpts, _from, _to, _id, _amount)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address _spender, bool _approved) returns()
func (_Traderjoev2 *Traderjoev2Transactor) SetApprovalForAll(opts *bind.TransactOpts, _spender common.Address, _approved bool) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "setApprovalForAll", _spender, _approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address _spender, bool _approved) returns()
func (_Traderjoev2 *Traderjoev2Session) SetApprovalForAll(_spender common.Address, _approved bool) (*types.Transaction, error) {
	return _Traderjoev2.Contract.SetApprovalForAll(&_Traderjoev2.TransactOpts, _spender, _approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address _spender, bool _approved) returns()
func (_Traderjoev2 *Traderjoev2TransactorSession) SetApprovalForAll(_spender common.Address, _approved bool) (*types.Transaction, error) {
	return _Traderjoev2.Contract.SetApprovalForAll(&_Traderjoev2.TransactOpts, _spender, _approved)
}

// SetFeesParameters is a paid mutator transaction binding the contract method 0x54b5fc87.
//
// Solidity: function setFeesParameters(bytes32 _packedFeeParameters) returns()
func (_Traderjoev2 *Traderjoev2Transactor) SetFeesParameters(opts *bind.TransactOpts, _packedFeeParameters [32]byte) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "setFeesParameters", _packedFeeParameters)
}

// SetFeesParameters is a paid mutator transaction binding the contract method 0x54b5fc87.
//
// Solidity: function setFeesParameters(bytes32 _packedFeeParameters) returns()
func (_Traderjoev2 *Traderjoev2Session) SetFeesParameters(_packedFeeParameters [32]byte) (*types.Transaction, error) {
	return _Traderjoev2.Contract.SetFeesParameters(&_Traderjoev2.TransactOpts, _packedFeeParameters)
}

// SetFeesParameters is a paid mutator transaction binding the contract method 0x54b5fc87.
//
// Solidity: function setFeesParameters(bytes32 _packedFeeParameters) returns()
func (_Traderjoev2 *Traderjoev2TransactorSession) SetFeesParameters(_packedFeeParameters [32]byte) (*types.Transaction, error) {
	return _Traderjoev2.Contract.SetFeesParameters(&_Traderjoev2.TransactOpts, _packedFeeParameters)
}

// Swap is a paid mutator transaction binding the contract method 0x53c059a0.
//
// Solidity: function swap(bool _swapForY, address _to) returns(uint256 amountXOut, uint256 amountYOut)
func (_Traderjoev2 *Traderjoev2Transactor) Swap(opts *bind.TransactOpts, _swapForY bool, _to common.Address) (*types.Transaction, error) {
	return _Traderjoev2.contract.Transact(opts, "swap", _swapForY, _to)
}

// Swap is a paid mutator transaction binding the contract method 0x53c059a0.
//
// Solidity: function swap(bool _swapForY, address _to) returns(uint256 amountXOut, uint256 amountYOut)
func (_Traderjoev2 *Traderjoev2Session) Swap(_swapForY bool, _to common.Address) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Swap(&_Traderjoev2.TransactOpts, _swapForY, _to)
}

// Swap is a paid mutator transaction binding the contract method 0x53c059a0.
//
// Solidity: function swap(bool _swapForY, address _to) returns(uint256 amountXOut, uint256 amountYOut)
func (_Traderjoev2 *Traderjoev2TransactorSession) Swap(_swapForY bool, _to common.Address) (*types.Transaction, error) {
	return _Traderjoev2.Contract.Swap(&_Traderjoev2.TransactOpts, _swapForY, _to)
}

// Traderjoev2ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Traderjoev2 contract.
type Traderjoev2ApprovalForAllIterator struct {
	Event *Traderjoev2ApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2ApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2ApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2ApprovalForAll represents a ApprovalForAll event raised by the Traderjoev2 contract.
type Traderjoev2ApprovalForAll struct {
	Account  common.Address
	Sender   common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed sender, bool approved)
func (_Traderjoev2 *Traderjoev2Filterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, sender []common.Address) (*Traderjoev2ApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "ApprovalForAll", accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2ApprovalForAllIterator{contract: _Traderjoev2.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed sender, bool approved)
func (_Traderjoev2 *Traderjoev2Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *Traderjoev2ApprovalForAll, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "ApprovalForAll", accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2ApprovalForAll)
				if err := _Traderjoev2.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed sender, bool approved)
func (_Traderjoev2 *Traderjoev2Filterer) ParseApprovalForAll(log types.Log) (*Traderjoev2ApprovalForAll, error) {
	event := new(Traderjoev2ApprovalForAll)
	if err := _Traderjoev2.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2CompositionFeeIterator is returned from FilterCompositionFee and is used to iterate over the raw logs and unpacked data for CompositionFee events raised by the Traderjoev2 contract.
type Traderjoev2CompositionFeeIterator struct {
	Event *Traderjoev2CompositionFee // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2CompositionFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2CompositionFee)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2CompositionFee)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2CompositionFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2CompositionFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2CompositionFee represents a CompositionFee event raised by the Traderjoev2 contract.
type Traderjoev2CompositionFee struct {
	Sender    common.Address
	Recipient common.Address
	Id        *big.Int
	FeesX     *big.Int
	FeesY     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCompositionFee is a free log retrieval operation binding the contract event 0x56f8e764728c77dd99ffbc1b64e6d02e227e6ec8214f165d4ef31351de136a0d.
//
// Solidity: event CompositionFee(address indexed sender, address indexed recipient, uint256 indexed id, uint256 feesX, uint256 feesY)
func (_Traderjoev2 *Traderjoev2Filterer) FilterCompositionFee(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, id []*big.Int) (*Traderjoev2CompositionFeeIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "CompositionFee", senderRule, recipientRule, idRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2CompositionFeeIterator{contract: _Traderjoev2.contract, event: "CompositionFee", logs: logs, sub: sub}, nil
}

// WatchCompositionFee is a free log subscription operation binding the contract event 0x56f8e764728c77dd99ffbc1b64e6d02e227e6ec8214f165d4ef31351de136a0d.
//
// Solidity: event CompositionFee(address indexed sender, address indexed recipient, uint256 indexed id, uint256 feesX, uint256 feesY)
func (_Traderjoev2 *Traderjoev2Filterer) WatchCompositionFee(opts *bind.WatchOpts, sink chan<- *Traderjoev2CompositionFee, sender []common.Address, recipient []common.Address, id []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "CompositionFee", senderRule, recipientRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2CompositionFee)
				if err := _Traderjoev2.contract.UnpackLog(event, "CompositionFee", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCompositionFee is a log parse operation binding the contract event 0x56f8e764728c77dd99ffbc1b64e6d02e227e6ec8214f165d4ef31351de136a0d.
//
// Solidity: event CompositionFee(address indexed sender, address indexed recipient, uint256 indexed id, uint256 feesX, uint256 feesY)
func (_Traderjoev2 *Traderjoev2Filterer) ParseCompositionFee(log types.Log) (*Traderjoev2CompositionFee, error) {
	event := new(Traderjoev2CompositionFee)
	if err := _Traderjoev2.contract.UnpackLog(event, "CompositionFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2DepositedToBinIterator is returned from FilterDepositedToBin and is used to iterate over the raw logs and unpacked data for DepositedToBin events raised by the Traderjoev2 contract.
type Traderjoev2DepositedToBinIterator struct {
	Event *Traderjoev2DepositedToBin // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2DepositedToBinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2DepositedToBin)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2DepositedToBin)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2DepositedToBinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2DepositedToBinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2DepositedToBin represents a DepositedToBin event raised by the Traderjoev2 contract.
type Traderjoev2DepositedToBin struct {
	Sender    common.Address
	Recipient common.Address
	Id        *big.Int
	AmountX   *big.Int
	AmountY   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDepositedToBin is a free log retrieval operation binding the contract event 0x4216cc3bd0c40a90259d92f800c06ede5c47765f41a488072b7e7104a1f95841.
//
// Solidity: event DepositedToBin(address indexed sender, address indexed recipient, uint256 indexed id, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) FilterDepositedToBin(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, id []*big.Int) (*Traderjoev2DepositedToBinIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "DepositedToBin", senderRule, recipientRule, idRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2DepositedToBinIterator{contract: _Traderjoev2.contract, event: "DepositedToBin", logs: logs, sub: sub}, nil
}

// WatchDepositedToBin is a free log subscription operation binding the contract event 0x4216cc3bd0c40a90259d92f800c06ede5c47765f41a488072b7e7104a1f95841.
//
// Solidity: event DepositedToBin(address indexed sender, address indexed recipient, uint256 indexed id, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) WatchDepositedToBin(opts *bind.WatchOpts, sink chan<- *Traderjoev2DepositedToBin, sender []common.Address, recipient []common.Address, id []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "DepositedToBin", senderRule, recipientRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2DepositedToBin)
				if err := _Traderjoev2.contract.UnpackLog(event, "DepositedToBin", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDepositedToBin is a log parse operation binding the contract event 0x4216cc3bd0c40a90259d92f800c06ede5c47765f41a488072b7e7104a1f95841.
//
// Solidity: event DepositedToBin(address indexed sender, address indexed recipient, uint256 indexed id, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) ParseDepositedToBin(log types.Log) (*Traderjoev2DepositedToBin, error) {
	event := new(Traderjoev2DepositedToBin)
	if err := _Traderjoev2.contract.UnpackLog(event, "DepositedToBin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2FeesCollectedIterator is returned from FilterFeesCollected and is used to iterate over the raw logs and unpacked data for FeesCollected events raised by the Traderjoev2 contract.
type Traderjoev2FeesCollectedIterator struct {
	Event *Traderjoev2FeesCollected // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2FeesCollectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2FeesCollected)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2FeesCollected)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2FeesCollectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2FeesCollectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2FeesCollected represents a FeesCollected event raised by the Traderjoev2 contract.
type Traderjoev2FeesCollected struct {
	Sender    common.Address
	Recipient common.Address
	AmountX   *big.Int
	AmountY   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFeesCollected is a free log retrieval operation binding the contract event 0x28a87b6059180e46de5fb9ab35eb043e8fe00ab45afcc7789e3934ecbbcde3ea.
//
// Solidity: event FeesCollected(address indexed sender, address indexed recipient, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) FilterFeesCollected(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*Traderjoev2FeesCollectedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "FeesCollected", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2FeesCollectedIterator{contract: _Traderjoev2.contract, event: "FeesCollected", logs: logs, sub: sub}, nil
}

// WatchFeesCollected is a free log subscription operation binding the contract event 0x28a87b6059180e46de5fb9ab35eb043e8fe00ab45afcc7789e3934ecbbcde3ea.
//
// Solidity: event FeesCollected(address indexed sender, address indexed recipient, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) WatchFeesCollected(opts *bind.WatchOpts, sink chan<- *Traderjoev2FeesCollected, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "FeesCollected", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2FeesCollected)
				if err := _Traderjoev2.contract.UnpackLog(event, "FeesCollected", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFeesCollected is a log parse operation binding the contract event 0x28a87b6059180e46de5fb9ab35eb043e8fe00ab45afcc7789e3934ecbbcde3ea.
//
// Solidity: event FeesCollected(address indexed sender, address indexed recipient, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) ParseFeesCollected(log types.Log) (*Traderjoev2FeesCollected, error) {
	event := new(Traderjoev2FeesCollected)
	if err := _Traderjoev2.contract.UnpackLog(event, "FeesCollected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2FlashLoanIterator is returned from FilterFlashLoan and is used to iterate over the raw logs and unpacked data for FlashLoan events raised by the Traderjoev2 contract.
type Traderjoev2FlashLoanIterator struct {
	Event *Traderjoev2FlashLoan // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2FlashLoanIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2FlashLoan)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2FlashLoan)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2FlashLoanIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2FlashLoanIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2FlashLoan represents a FlashLoan event raised by the Traderjoev2 contract.
type Traderjoev2FlashLoan struct {
	Sender   common.Address
	Receiver common.Address
	Token    common.Address
	Amount   *big.Int
	Fee      *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFlashLoan is a free log retrieval operation binding the contract event 0x3659d15bd4bb92ab352a8d35bc3119ec6e7e0ab48e4d46201c8a28e02b6a8a86.
//
// Solidity: event FlashLoan(address indexed sender, address indexed receiver, address token, uint256 amount, uint256 fee)
func (_Traderjoev2 *Traderjoev2Filterer) FilterFlashLoan(opts *bind.FilterOpts, sender []common.Address, receiver []common.Address) (*Traderjoev2FlashLoanIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "FlashLoan", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2FlashLoanIterator{contract: _Traderjoev2.contract, event: "FlashLoan", logs: logs, sub: sub}, nil
}

// WatchFlashLoan is a free log subscription operation binding the contract event 0x3659d15bd4bb92ab352a8d35bc3119ec6e7e0ab48e4d46201c8a28e02b6a8a86.
//
// Solidity: event FlashLoan(address indexed sender, address indexed receiver, address token, uint256 amount, uint256 fee)
func (_Traderjoev2 *Traderjoev2Filterer) WatchFlashLoan(opts *bind.WatchOpts, sink chan<- *Traderjoev2FlashLoan, sender []common.Address, receiver []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "FlashLoan", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2FlashLoan)
				if err := _Traderjoev2.contract.UnpackLog(event, "FlashLoan", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFlashLoan is a log parse operation binding the contract event 0x3659d15bd4bb92ab352a8d35bc3119ec6e7e0ab48e4d46201c8a28e02b6a8a86.
//
// Solidity: event FlashLoan(address indexed sender, address indexed receiver, address token, uint256 amount, uint256 fee)
func (_Traderjoev2 *Traderjoev2Filterer) ParseFlashLoan(log types.Log) (*Traderjoev2FlashLoan, error) {
	event := new(Traderjoev2FlashLoan)
	if err := _Traderjoev2.contract.UnpackLog(event, "FlashLoan", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2OracleSizeIncreasedIterator is returned from FilterOracleSizeIncreased and is used to iterate over the raw logs and unpacked data for OracleSizeIncreased events raised by the Traderjoev2 contract.
type Traderjoev2OracleSizeIncreasedIterator struct {
	Event *Traderjoev2OracleSizeIncreased // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2OracleSizeIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2OracleSizeIncreased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2OracleSizeIncreased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2OracleSizeIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2OracleSizeIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2OracleSizeIncreased represents a OracleSizeIncreased event raised by the Traderjoev2 contract.
type Traderjoev2OracleSizeIncreased struct {
	PreviousSize *big.Int
	NewSize      *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOracleSizeIncreased is a free log retrieval operation binding the contract event 0x525a4241308ea122822834c841f67b00d5efc977ad9118724750f974f7f6531c.
//
// Solidity: event OracleSizeIncreased(uint256 previousSize, uint256 newSize)
func (_Traderjoev2 *Traderjoev2Filterer) FilterOracleSizeIncreased(opts *bind.FilterOpts) (*Traderjoev2OracleSizeIncreasedIterator, error) {

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "OracleSizeIncreased")
	if err != nil {
		return nil, err
	}
	return &Traderjoev2OracleSizeIncreasedIterator{contract: _Traderjoev2.contract, event: "OracleSizeIncreased", logs: logs, sub: sub}, nil
}

// WatchOracleSizeIncreased is a free log subscription operation binding the contract event 0x525a4241308ea122822834c841f67b00d5efc977ad9118724750f974f7f6531c.
//
// Solidity: event OracleSizeIncreased(uint256 previousSize, uint256 newSize)
func (_Traderjoev2 *Traderjoev2Filterer) WatchOracleSizeIncreased(opts *bind.WatchOpts, sink chan<- *Traderjoev2OracleSizeIncreased) (event.Subscription, error) {

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "OracleSizeIncreased")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2OracleSizeIncreased)
				if err := _Traderjoev2.contract.UnpackLog(event, "OracleSizeIncreased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOracleSizeIncreased is a log parse operation binding the contract event 0x525a4241308ea122822834c841f67b00d5efc977ad9118724750f974f7f6531c.
//
// Solidity: event OracleSizeIncreased(uint256 previousSize, uint256 newSize)
func (_Traderjoev2 *Traderjoev2Filterer) ParseOracleSizeIncreased(log types.Log) (*Traderjoev2OracleSizeIncreased, error) {
	event := new(Traderjoev2OracleSizeIncreased)
	if err := _Traderjoev2.contract.UnpackLog(event, "OracleSizeIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2ProtocolFeesCollectedIterator is returned from FilterProtocolFeesCollected and is used to iterate over the raw logs and unpacked data for ProtocolFeesCollected events raised by the Traderjoev2 contract.
type Traderjoev2ProtocolFeesCollectedIterator struct {
	Event *Traderjoev2ProtocolFeesCollected // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2ProtocolFeesCollectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2ProtocolFeesCollected)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2ProtocolFeesCollected)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2ProtocolFeesCollectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2ProtocolFeesCollectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2ProtocolFeesCollected represents a ProtocolFeesCollected event raised by the Traderjoev2 contract.
type Traderjoev2ProtocolFeesCollected struct {
	Sender    common.Address
	Recipient common.Address
	AmountX   *big.Int
	AmountY   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeesCollected is a free log retrieval operation binding the contract event 0x26b782206d6b531bf95d487110cfefdc443291f176f1977e94abcb7e67bd1b79.
//
// Solidity: event ProtocolFeesCollected(address indexed sender, address indexed recipient, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) FilterProtocolFeesCollected(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*Traderjoev2ProtocolFeesCollectedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "ProtocolFeesCollected", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2ProtocolFeesCollectedIterator{contract: _Traderjoev2.contract, event: "ProtocolFeesCollected", logs: logs, sub: sub}, nil
}

// WatchProtocolFeesCollected is a free log subscription operation binding the contract event 0x26b782206d6b531bf95d487110cfefdc443291f176f1977e94abcb7e67bd1b79.
//
// Solidity: event ProtocolFeesCollected(address indexed sender, address indexed recipient, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) WatchProtocolFeesCollected(opts *bind.WatchOpts, sink chan<- *Traderjoev2ProtocolFeesCollected, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "ProtocolFeesCollected", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2ProtocolFeesCollected)
				if err := _Traderjoev2.contract.UnpackLog(event, "ProtocolFeesCollected", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProtocolFeesCollected is a log parse operation binding the contract event 0x26b782206d6b531bf95d487110cfefdc443291f176f1977e94abcb7e67bd1b79.
//
// Solidity: event ProtocolFeesCollected(address indexed sender, address indexed recipient, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) ParseProtocolFeesCollected(log types.Log) (*Traderjoev2ProtocolFeesCollected, error) {
	event := new(Traderjoev2ProtocolFeesCollected)
	if err := _Traderjoev2.contract.UnpackLog(event, "ProtocolFeesCollected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2SwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the Traderjoev2 contract.
type Traderjoev2SwapIterator struct {
	Event *Traderjoev2Swap // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2SwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2Swap)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2Swap)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2SwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2SwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2Swap represents a Swap event raised by the Traderjoev2 contract.
type Traderjoev2Swap struct {
	Sender                common.Address
	Recipient             common.Address
	Id                    *big.Int
	SwapForY              bool
	AmountIn              *big.Int
	AmountOut             *big.Int
	VolatilityAccumulated *big.Int
	Fees                  *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0xc528cda9e500228b16ce84fadae290d9a49aecb17483110004c5af0a07f6fd73.
//
// Solidity: event Swap(address indexed sender, address indexed recipient, uint256 indexed id, bool swapForY, uint256 amountIn, uint256 amountOut, uint256 volatilityAccumulated, uint256 fees)
func (_Traderjoev2 *Traderjoev2Filterer) FilterSwap(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, id []*big.Int) (*Traderjoev2SwapIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "Swap", senderRule, recipientRule, idRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2SwapIterator{contract: _Traderjoev2.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0xc528cda9e500228b16ce84fadae290d9a49aecb17483110004c5af0a07f6fd73.
//
// Solidity: event Swap(address indexed sender, address indexed recipient, uint256 indexed id, bool swapForY, uint256 amountIn, uint256 amountOut, uint256 volatilityAccumulated, uint256 fees)
func (_Traderjoev2 *Traderjoev2Filterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *Traderjoev2Swap, sender []common.Address, recipient []common.Address, id []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "Swap", senderRule, recipientRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2Swap)
				if err := _Traderjoev2.contract.UnpackLog(event, "Swap", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSwap is a log parse operation binding the contract event 0xc528cda9e500228b16ce84fadae290d9a49aecb17483110004c5af0a07f6fd73.
//
// Solidity: event Swap(address indexed sender, address indexed recipient, uint256 indexed id, bool swapForY, uint256 amountIn, uint256 amountOut, uint256 volatilityAccumulated, uint256 fees)
func (_Traderjoev2 *Traderjoev2Filterer) ParseSwap(log types.Log) (*Traderjoev2Swap, error) {
	event := new(Traderjoev2Swap)
	if err := _Traderjoev2.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2TransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the Traderjoev2 contract.
type Traderjoev2TransferBatchIterator struct {
	Event *Traderjoev2TransferBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2TransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2TransferBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2TransferBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2TransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2TransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2TransferBatch represents a TransferBatch event raised by the Traderjoev2 contract.
type Traderjoev2TransferBatch struct {
	Sender  common.Address
	From    common.Address
	To      common.Address
	Ids     []*big.Int
	Amounts []*big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed sender, address indexed from, address indexed to, uint256[] ids, uint256[] amounts)
func (_Traderjoev2 *Traderjoev2Filterer) FilterTransferBatch(opts *bind.FilterOpts, sender []common.Address, from []common.Address, to []common.Address) (*Traderjoev2TransferBatchIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "TransferBatch", senderRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2TransferBatchIterator{contract: _Traderjoev2.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed sender, address indexed from, address indexed to, uint256[] ids, uint256[] amounts)
func (_Traderjoev2 *Traderjoev2Filterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *Traderjoev2TransferBatch, sender []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "TransferBatch", senderRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2TransferBatch)
				if err := _Traderjoev2.contract.UnpackLog(event, "TransferBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed sender, address indexed from, address indexed to, uint256[] ids, uint256[] amounts)
func (_Traderjoev2 *Traderjoev2Filterer) ParseTransferBatch(log types.Log) (*Traderjoev2TransferBatch, error) {
	event := new(Traderjoev2TransferBatch)
	if err := _Traderjoev2.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2TransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the Traderjoev2 contract.
type Traderjoev2TransferSingleIterator struct {
	Event *Traderjoev2TransferSingle // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2TransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2TransferSingle)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2TransferSingle)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2TransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2TransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2TransferSingle represents a TransferSingle event raised by the Traderjoev2 contract.
type Traderjoev2TransferSingle struct {
	Sender common.Address
	From   common.Address
	To     common.Address
	Id     *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed sender, address indexed from, address indexed to, uint256 id, uint256 amount)
func (_Traderjoev2 *Traderjoev2Filterer) FilterTransferSingle(opts *bind.FilterOpts, sender []common.Address, from []common.Address, to []common.Address) (*Traderjoev2TransferSingleIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "TransferSingle", senderRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2TransferSingleIterator{contract: _Traderjoev2.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed sender, address indexed from, address indexed to, uint256 id, uint256 amount)
func (_Traderjoev2 *Traderjoev2Filterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *Traderjoev2TransferSingle, sender []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "TransferSingle", senderRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2TransferSingle)
				if err := _Traderjoev2.contract.UnpackLog(event, "TransferSingle", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed sender, address indexed from, address indexed to, uint256 id, uint256 amount)
func (_Traderjoev2 *Traderjoev2Filterer) ParseTransferSingle(log types.Log) (*Traderjoev2TransferSingle, error) {
	event := new(Traderjoev2TransferSingle)
	if err := _Traderjoev2.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Traderjoev2WithdrawnFromBinIterator is returned from FilterWithdrawnFromBin and is used to iterate over the raw logs and unpacked data for WithdrawnFromBin events raised by the Traderjoev2 contract.
type Traderjoev2WithdrawnFromBinIterator struct {
	Event *Traderjoev2WithdrawnFromBin // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Traderjoev2WithdrawnFromBinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Traderjoev2WithdrawnFromBin)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Traderjoev2WithdrawnFromBin)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Traderjoev2WithdrawnFromBinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Traderjoev2WithdrawnFromBinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Traderjoev2WithdrawnFromBin represents a WithdrawnFromBin event raised by the Traderjoev2 contract.
type Traderjoev2WithdrawnFromBin struct {
	Sender    common.Address
	Recipient common.Address
	Id        *big.Int
	AmountX   *big.Int
	AmountY   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawnFromBin is a free log retrieval operation binding the contract event 0xda5e7177dface55f5e0eff7dfc67420a1db4243ddfcf0ecc84ed93e034dd8cc2.
//
// Solidity: event WithdrawnFromBin(address indexed sender, address indexed recipient, uint256 indexed id, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) FilterWithdrawnFromBin(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, id []*big.Int) (*Traderjoev2WithdrawnFromBinIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Traderjoev2.contract.FilterLogs(opts, "WithdrawnFromBin", senderRule, recipientRule, idRule)
	if err != nil {
		return nil, err
	}
	return &Traderjoev2WithdrawnFromBinIterator{contract: _Traderjoev2.contract, event: "WithdrawnFromBin", logs: logs, sub: sub}, nil
}

// WatchWithdrawnFromBin is a free log subscription operation binding the contract event 0xda5e7177dface55f5e0eff7dfc67420a1db4243ddfcf0ecc84ed93e034dd8cc2.
//
// Solidity: event WithdrawnFromBin(address indexed sender, address indexed recipient, uint256 indexed id, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) WatchWithdrawnFromBin(opts *bind.WatchOpts, sink chan<- *Traderjoev2WithdrawnFromBin, sender []common.Address, recipient []common.Address, id []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Traderjoev2.contract.WatchLogs(opts, "WithdrawnFromBin", senderRule, recipientRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Traderjoev2WithdrawnFromBin)
				if err := _Traderjoev2.contract.UnpackLog(event, "WithdrawnFromBin", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawnFromBin is a log parse operation binding the contract event 0xda5e7177dface55f5e0eff7dfc67420a1db4243ddfcf0ecc84ed93e034dd8cc2.
//
// Solidity: event WithdrawnFromBin(address indexed sender, address indexed recipient, uint256 indexed id, uint256 amountX, uint256 amountY)
func (_Traderjoev2 *Traderjoev2Filterer) ParseWithdrawnFromBin(log types.Log) (*Traderjoev2WithdrawnFromBin, error) {
	event := new(Traderjoev2WithdrawnFromBin)
	if err := _Traderjoev2.contract.UnpackLog(event, "WithdrawnFromBin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
