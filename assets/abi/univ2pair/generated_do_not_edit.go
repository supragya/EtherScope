// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package univ2pair

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

// Univ2pairMetaData contains all meta data concerning the Univ2pair contract.
var Univ2pairMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint112\",\"name\":\"reserve0\",\"type\":\"uint112\"},{\"indexed\":false,\"internalType\":\"uint112\",\"name\":\"reserve1\",\"type\":\"uint112\"}],\"name\":\"Sync\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MINIMUM_LIQUIDITY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"PERMIT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint112\",\"name\":\"_reserve0\",\"type\":\"uint112\"},{\"internalType\":\"uint112\",\"name\":\"_reserve1\",\"type\":\"uint112\"},{\"internalType\":\"uint32\",\"name\":\"_blockTimestampLast\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token1\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"kLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"price0CumulativeLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"price1CumulativeLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"skim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"sync\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token0\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token1\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// Univ2pairABI is the input ABI used to generate the binding from.
// Deprecated: Use Univ2pairMetaData.ABI instead.
var Univ2pairABI = Univ2pairMetaData.ABI

// Univ2pair is an auto generated Go binding around an Ethereum contract.
type Univ2pair struct {
	Univ2pairCaller     // Read-only binding to the contract
	Univ2pairTransactor // Write-only binding to the contract
	Univ2pairFilterer   // Log filterer for contract events
}

// Univ2pairCaller is an auto generated read-only Go binding around an Ethereum contract.
type Univ2pairCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Univ2pairTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Univ2pairTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Univ2pairFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Univ2pairFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Univ2pairSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Univ2pairSession struct {
	Contract     *Univ2pair        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Univ2pairCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Univ2pairCallerSession struct {
	Contract *Univ2pairCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// Univ2pairTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Univ2pairTransactorSession struct {
	Contract     *Univ2pairTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// Univ2pairRaw is an auto generated low-level Go binding around an Ethereum contract.
type Univ2pairRaw struct {
	Contract *Univ2pair // Generic contract binding to access the raw methods on
}

// Univ2pairCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Univ2pairCallerRaw struct {
	Contract *Univ2pairCaller // Generic read-only contract binding to access the raw methods on
}

// Univ2pairTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Univ2pairTransactorRaw struct {
	Contract *Univ2pairTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniv2pair creates a new instance of Univ2pair, bound to a specific deployed contract.
func NewUniv2pair(address common.Address, backend bind.ContractBackend) (*Univ2pair, error) {
	contract, err := bindUniv2pair(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Univ2pair{Univ2pairCaller: Univ2pairCaller{contract: contract}, Univ2pairTransactor: Univ2pairTransactor{contract: contract}, Univ2pairFilterer: Univ2pairFilterer{contract: contract}}, nil
}

// NewUniv2pairCaller creates a new read-only instance of Univ2pair, bound to a specific deployed contract.
func NewUniv2pairCaller(address common.Address, caller bind.ContractCaller) (*Univ2pairCaller, error) {
	contract, err := bindUniv2pair(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Univ2pairCaller{contract: contract}, nil
}

// NewUniv2pairTransactor creates a new write-only instance of Univ2pair, bound to a specific deployed contract.
func NewUniv2pairTransactor(address common.Address, transactor bind.ContractTransactor) (*Univ2pairTransactor, error) {
	contract, err := bindUniv2pair(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Univ2pairTransactor{contract: contract}, nil
}

// NewUniv2pairFilterer creates a new log filterer instance of Univ2pair, bound to a specific deployed contract.
func NewUniv2pairFilterer(address common.Address, filterer bind.ContractFilterer) (*Univ2pairFilterer, error) {
	contract, err := bindUniv2pair(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Univ2pairFilterer{contract: contract}, nil
}

// bindUniv2pair binds a generic wrapper to an already deployed contract.
func bindUniv2pair(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Univ2pairMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Univ2pair *Univ2pairRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Univ2pair.Contract.Univ2pairCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Univ2pair *Univ2pairRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Univ2pair.Contract.Univ2pairTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Univ2pair *Univ2pairRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Univ2pair.Contract.Univ2pairTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Univ2pair *Univ2pairCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Univ2pair.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Univ2pair *Univ2pairTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Univ2pair.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Univ2pair *Univ2pairTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Univ2pair.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Univ2pair *Univ2pairCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Univ2pair *Univ2pairSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Univ2pair.Contract.DOMAINSEPARATOR(&_Univ2pair.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Univ2pair *Univ2pairCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Univ2pair.Contract.DOMAINSEPARATOR(&_Univ2pair.CallOpts)
}

// MINIMUMLIQUIDITY is a free data retrieval call binding the contract method 0xba9a7a56.
//
// Solidity: function MINIMUM_LIQUIDITY() view returns(uint256)
func (_Univ2pair *Univ2pairCaller) MINIMUMLIQUIDITY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "MINIMUM_LIQUIDITY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINIMUMLIQUIDITY is a free data retrieval call binding the contract method 0xba9a7a56.
//
// Solidity: function MINIMUM_LIQUIDITY() view returns(uint256)
func (_Univ2pair *Univ2pairSession) MINIMUMLIQUIDITY() (*big.Int, error) {
	return _Univ2pair.Contract.MINIMUMLIQUIDITY(&_Univ2pair.CallOpts)
}

// MINIMUMLIQUIDITY is a free data retrieval call binding the contract method 0xba9a7a56.
//
// Solidity: function MINIMUM_LIQUIDITY() view returns(uint256)
func (_Univ2pair *Univ2pairCallerSession) MINIMUMLIQUIDITY() (*big.Int, error) {
	return _Univ2pair.Contract.MINIMUMLIQUIDITY(&_Univ2pair.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_Univ2pair *Univ2pairCaller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "PERMIT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_Univ2pair *Univ2pairSession) PERMITTYPEHASH() ([32]byte, error) {
	return _Univ2pair.Contract.PERMITTYPEHASH(&_Univ2pair.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_Univ2pair *Univ2pairCallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _Univ2pair.Contract.PERMITTYPEHASH(&_Univ2pair.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_Univ2pair *Univ2pairCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_Univ2pair *Univ2pairSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Univ2pair.Contract.Allowance(&_Univ2pair.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_Univ2pair *Univ2pairCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Univ2pair.Contract.Allowance(&_Univ2pair.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_Univ2pair *Univ2pairCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_Univ2pair *Univ2pairSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _Univ2pair.Contract.BalanceOf(&_Univ2pair.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_Univ2pair *Univ2pairCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _Univ2pair.Contract.BalanceOf(&_Univ2pair.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Univ2pair *Univ2pairCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Univ2pair *Univ2pairSession) Decimals() (uint8, error) {
	return _Univ2pair.Contract.Decimals(&_Univ2pair.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Univ2pair *Univ2pairCallerSession) Decimals() (uint8, error) {
	return _Univ2pair.Contract.Decimals(&_Univ2pair.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_Univ2pair *Univ2pairCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_Univ2pair *Univ2pairSession) Factory() (common.Address, error) {
	return _Univ2pair.Contract.Factory(&_Univ2pair.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_Univ2pair *Univ2pairCallerSession) Factory() (common.Address, error) {
	return _Univ2pair.Contract.Factory(&_Univ2pair.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 _reserve0, uint112 _reserve1, uint32 _blockTimestampLast)
func (_Univ2pair *Univ2pairCaller) GetReserves(opts *bind.CallOpts) (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "getReserves")

	outstruct := new(struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reserve0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserve1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockTimestampLast = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 _reserve0, uint112 _reserve1, uint32 _blockTimestampLast)
func (_Univ2pair *Univ2pairSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {
	return _Univ2pair.Contract.GetReserves(&_Univ2pair.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 _reserve0, uint112 _reserve1, uint32 _blockTimestampLast)
func (_Univ2pair *Univ2pairCallerSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {
	return _Univ2pair.Contract.GetReserves(&_Univ2pair.CallOpts)
}

// KLast is a free data retrieval call binding the contract method 0x7464fc3d.
//
// Solidity: function kLast() view returns(uint256)
func (_Univ2pair *Univ2pairCaller) KLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "kLast")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// KLast is a free data retrieval call binding the contract method 0x7464fc3d.
//
// Solidity: function kLast() view returns(uint256)
func (_Univ2pair *Univ2pairSession) KLast() (*big.Int, error) {
	return _Univ2pair.Contract.KLast(&_Univ2pair.CallOpts)
}

// KLast is a free data retrieval call binding the contract method 0x7464fc3d.
//
// Solidity: function kLast() view returns(uint256)
func (_Univ2pair *Univ2pairCallerSession) KLast() (*big.Int, error) {
	return _Univ2pair.Contract.KLast(&_Univ2pair.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Univ2pair *Univ2pairCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Univ2pair *Univ2pairSession) Name() (string, error) {
	return _Univ2pair.Contract.Name(&_Univ2pair.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Univ2pair *Univ2pairCallerSession) Name() (string, error) {
	return _Univ2pair.Contract.Name(&_Univ2pair.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Univ2pair *Univ2pairCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Univ2pair *Univ2pairSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Univ2pair.Contract.Nonces(&_Univ2pair.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Univ2pair *Univ2pairCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Univ2pair.Contract.Nonces(&_Univ2pair.CallOpts, arg0)
}

// Price0CumulativeLast is a free data retrieval call binding the contract method 0x5909c0d5.
//
// Solidity: function price0CumulativeLast() view returns(uint256)
func (_Univ2pair *Univ2pairCaller) Price0CumulativeLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "price0CumulativeLast")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Price0CumulativeLast is a free data retrieval call binding the contract method 0x5909c0d5.
//
// Solidity: function price0CumulativeLast() view returns(uint256)
func (_Univ2pair *Univ2pairSession) Price0CumulativeLast() (*big.Int, error) {
	return _Univ2pair.Contract.Price0CumulativeLast(&_Univ2pair.CallOpts)
}

// Price0CumulativeLast is a free data retrieval call binding the contract method 0x5909c0d5.
//
// Solidity: function price0CumulativeLast() view returns(uint256)
func (_Univ2pair *Univ2pairCallerSession) Price0CumulativeLast() (*big.Int, error) {
	return _Univ2pair.Contract.Price0CumulativeLast(&_Univ2pair.CallOpts)
}

// Price1CumulativeLast is a free data retrieval call binding the contract method 0x5a3d5493.
//
// Solidity: function price1CumulativeLast() view returns(uint256)
func (_Univ2pair *Univ2pairCaller) Price1CumulativeLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "price1CumulativeLast")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Price1CumulativeLast is a free data retrieval call binding the contract method 0x5a3d5493.
//
// Solidity: function price1CumulativeLast() view returns(uint256)
func (_Univ2pair *Univ2pairSession) Price1CumulativeLast() (*big.Int, error) {
	return _Univ2pair.Contract.Price1CumulativeLast(&_Univ2pair.CallOpts)
}

// Price1CumulativeLast is a free data retrieval call binding the contract method 0x5a3d5493.
//
// Solidity: function price1CumulativeLast() view returns(uint256)
func (_Univ2pair *Univ2pairCallerSession) Price1CumulativeLast() (*big.Int, error) {
	return _Univ2pair.Contract.Price1CumulativeLast(&_Univ2pair.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Univ2pair *Univ2pairCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Univ2pair *Univ2pairSession) Symbol() (string, error) {
	return _Univ2pair.Contract.Symbol(&_Univ2pair.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Univ2pair *Univ2pairCallerSession) Symbol() (string, error) {
	return _Univ2pair.Contract.Symbol(&_Univ2pair.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_Univ2pair *Univ2pairCaller) Token0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "token0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_Univ2pair *Univ2pairSession) Token0() (common.Address, error) {
	return _Univ2pair.Contract.Token0(&_Univ2pair.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_Univ2pair *Univ2pairCallerSession) Token0() (common.Address, error) {
	return _Univ2pair.Contract.Token0(&_Univ2pair.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_Univ2pair *Univ2pairCaller) Token1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "token1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_Univ2pair *Univ2pairSession) Token1() (common.Address, error) {
	return _Univ2pair.Contract.Token1(&_Univ2pair.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_Univ2pair *Univ2pairCallerSession) Token1() (common.Address, error) {
	return _Univ2pair.Contract.Token1(&_Univ2pair.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Univ2pair *Univ2pairCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Univ2pair.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Univ2pair *Univ2pairSession) TotalSupply() (*big.Int, error) {
	return _Univ2pair.Contract.TotalSupply(&_Univ2pair.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Univ2pair *Univ2pairCallerSession) TotalSupply() (*big.Int, error) {
	return _Univ2pair.Contract.TotalSupply(&_Univ2pair.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Univ2pair *Univ2pairTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Univ2pair *Univ2pairSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Univ2pair.Contract.Approve(&_Univ2pair.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Univ2pair *Univ2pairTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Univ2pair.Contract.Approve(&_Univ2pair.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_Univ2pair *Univ2pairTransactor) Burn(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "burn", to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_Univ2pair *Univ2pairSession) Burn(to common.Address) (*types.Transaction, error) {
	return _Univ2pair.Contract.Burn(&_Univ2pair.TransactOpts, to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_Univ2pair *Univ2pairTransactorSession) Burn(to common.Address) (*types.Transaction, error) {
	return _Univ2pair.Contract.Burn(&_Univ2pair.TransactOpts, to)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _token0, address _token1) returns()
func (_Univ2pair *Univ2pairTransactor) Initialize(opts *bind.TransactOpts, _token0 common.Address, _token1 common.Address) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "initialize", _token0, _token1)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _token0, address _token1) returns()
func (_Univ2pair *Univ2pairSession) Initialize(_token0 common.Address, _token1 common.Address) (*types.Transaction, error) {
	return _Univ2pair.Contract.Initialize(&_Univ2pair.TransactOpts, _token0, _token1)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _token0, address _token1) returns()
func (_Univ2pair *Univ2pairTransactorSession) Initialize(_token0 common.Address, _token1 common.Address) (*types.Transaction, error) {
	return _Univ2pair.Contract.Initialize(&_Univ2pair.TransactOpts, _token0, _token1)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_Univ2pair *Univ2pairTransactor) Mint(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "mint", to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_Univ2pair *Univ2pairSession) Mint(to common.Address) (*types.Transaction, error) {
	return _Univ2pair.Contract.Mint(&_Univ2pair.TransactOpts, to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_Univ2pair *Univ2pairTransactorSession) Mint(to common.Address) (*types.Transaction, error) {
	return _Univ2pair.Contract.Mint(&_Univ2pair.TransactOpts, to)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Univ2pair *Univ2pairTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Univ2pair *Univ2pairSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Univ2pair.Contract.Permit(&_Univ2pair.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Univ2pair *Univ2pairTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Univ2pair.Contract.Permit(&_Univ2pair.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_Univ2pair *Univ2pairTransactor) Skim(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "skim", to)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_Univ2pair *Univ2pairSession) Skim(to common.Address) (*types.Transaction, error) {
	return _Univ2pair.Contract.Skim(&_Univ2pair.TransactOpts, to)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_Univ2pair *Univ2pairTransactorSession) Skim(to common.Address) (*types.Transaction, error) {
	return _Univ2pair.Contract.Skim(&_Univ2pair.TransactOpts, to)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_Univ2pair *Univ2pairTransactor) Swap(opts *bind.TransactOpts, amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "swap", amount0Out, amount1Out, to, data)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_Univ2pair *Univ2pairSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _Univ2pair.Contract.Swap(&_Univ2pair.TransactOpts, amount0Out, amount1Out, to, data)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_Univ2pair *Univ2pairTransactorSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _Univ2pair.Contract.Swap(&_Univ2pair.TransactOpts, amount0Out, amount1Out, to, data)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_Univ2pair *Univ2pairTransactor) Sync(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "sync")
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_Univ2pair *Univ2pairSession) Sync() (*types.Transaction, error) {
	return _Univ2pair.Contract.Sync(&_Univ2pair.TransactOpts)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_Univ2pair *Univ2pairTransactorSession) Sync() (*types.Transaction, error) {
	return _Univ2pair.Contract.Sync(&_Univ2pair.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Univ2pair *Univ2pairTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Univ2pair *Univ2pairSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Univ2pair.Contract.Transfer(&_Univ2pair.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Univ2pair *Univ2pairTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Univ2pair.Contract.Transfer(&_Univ2pair.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Univ2pair *Univ2pairTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Univ2pair.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Univ2pair *Univ2pairSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Univ2pair.Contract.TransferFrom(&_Univ2pair.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Univ2pair *Univ2pairTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Univ2pair.Contract.TransferFrom(&_Univ2pair.TransactOpts, from, to, value)
}

// Univ2pairApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Univ2pair contract.
type Univ2pairApprovalIterator struct {
	Event *Univ2pairApproval // Event containing the contract specifics and raw log

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
func (it *Univ2pairApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Univ2pairApproval)
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
		it.Event = new(Univ2pairApproval)
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
func (it *Univ2pairApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Univ2pairApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Univ2pairApproval represents a Approval event raised by the Univ2pair contract.
type Univ2pairApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Univ2pair *Univ2pairFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*Univ2pairApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Univ2pair.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &Univ2pairApprovalIterator{contract: _Univ2pair.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Univ2pair *Univ2pairFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Univ2pairApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Univ2pair.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Univ2pairApproval)
				if err := _Univ2pair.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Univ2pair *Univ2pairFilterer) ParseApproval(log types.Log) (*Univ2pairApproval, error) {
	event := new(Univ2pairApproval)
	if err := _Univ2pair.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Univ2pairBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the Univ2pair contract.
type Univ2pairBurnIterator struct {
	Event *Univ2pairBurn // Event containing the contract specifics and raw log

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
func (it *Univ2pairBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Univ2pairBurn)
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
		it.Event = new(Univ2pairBurn)
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
func (it *Univ2pairBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Univ2pairBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Univ2pairBurn represents a Burn event raised by the Univ2pair contract.
type Univ2pairBurn struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	To      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_Univ2pair *Univ2pairFilterer) FilterBurn(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*Univ2pairBurnIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Univ2pair.contract.FilterLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Univ2pairBurnIterator{contract: _Univ2pair.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_Univ2pair *Univ2pairFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *Univ2pairBurn, sender []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Univ2pair.contract.WatchLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Univ2pairBurn)
				if err := _Univ2pair.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_Univ2pair *Univ2pairFilterer) ParseBurn(log types.Log) (*Univ2pairBurn, error) {
	event := new(Univ2pairBurn)
	if err := _Univ2pair.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Univ2pairMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the Univ2pair contract.
type Univ2pairMintIterator struct {
	Event *Univ2pairMint // Event containing the contract specifics and raw log

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
func (it *Univ2pairMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Univ2pairMint)
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
		it.Event = new(Univ2pairMint)
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
func (it *Univ2pairMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Univ2pairMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Univ2pairMint represents a Mint event raised by the Univ2pair contract.
type Univ2pairMint struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_Univ2pair *Univ2pairFilterer) FilterMint(opts *bind.FilterOpts, sender []common.Address) (*Univ2pairMintIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Univ2pair.contract.FilterLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return &Univ2pairMintIterator{contract: _Univ2pair.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_Univ2pair *Univ2pairFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *Univ2pairMint, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Univ2pair.contract.WatchLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Univ2pairMint)
				if err := _Univ2pair.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_Univ2pair *Univ2pairFilterer) ParseMint(log types.Log) (*Univ2pairMint, error) {
	event := new(Univ2pairMint)
	if err := _Univ2pair.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Univ2pairSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the Univ2pair contract.
type Univ2pairSwapIterator struct {
	Event *Univ2pairSwap // Event containing the contract specifics and raw log

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
func (it *Univ2pairSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Univ2pairSwap)
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
		it.Event = new(Univ2pairSwap)
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
func (it *Univ2pairSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Univ2pairSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Univ2pairSwap represents a Swap event raised by the Univ2pair contract.
type Univ2pairSwap struct {
	Sender     common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To         common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_Univ2pair *Univ2pairFilterer) FilterSwap(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*Univ2pairSwapIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Univ2pair.contract.FilterLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Univ2pairSwapIterator{contract: _Univ2pair.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_Univ2pair *Univ2pairFilterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *Univ2pairSwap, sender []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Univ2pair.contract.WatchLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Univ2pairSwap)
				if err := _Univ2pair.contract.UnpackLog(event, "Swap", log); err != nil {
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

// ParseSwap is a log parse operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_Univ2pair *Univ2pairFilterer) ParseSwap(log types.Log) (*Univ2pairSwap, error) {
	event := new(Univ2pairSwap)
	if err := _Univ2pair.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Univ2pairSyncIterator is returned from FilterSync and is used to iterate over the raw logs and unpacked data for Sync events raised by the Univ2pair contract.
type Univ2pairSyncIterator struct {
	Event *Univ2pairSync // Event containing the contract specifics and raw log

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
func (it *Univ2pairSyncIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Univ2pairSync)
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
		it.Event = new(Univ2pairSync)
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
func (it *Univ2pairSyncIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Univ2pairSyncIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Univ2pairSync represents a Sync event raised by the Univ2pair contract.
type Univ2pairSync struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSync is a free log retrieval operation binding the contract event 0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1.
//
// Solidity: event Sync(uint112 reserve0, uint112 reserve1)
func (_Univ2pair *Univ2pairFilterer) FilterSync(opts *bind.FilterOpts) (*Univ2pairSyncIterator, error) {

	logs, sub, err := _Univ2pair.contract.FilterLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return &Univ2pairSyncIterator{contract: _Univ2pair.contract, event: "Sync", logs: logs, sub: sub}, nil
}

// WatchSync is a free log subscription operation binding the contract event 0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1.
//
// Solidity: event Sync(uint112 reserve0, uint112 reserve1)
func (_Univ2pair *Univ2pairFilterer) WatchSync(opts *bind.WatchOpts, sink chan<- *Univ2pairSync) (event.Subscription, error) {

	logs, sub, err := _Univ2pair.contract.WatchLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Univ2pairSync)
				if err := _Univ2pair.contract.UnpackLog(event, "Sync", log); err != nil {
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

// ParseSync is a log parse operation binding the contract event 0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1.
//
// Solidity: event Sync(uint112 reserve0, uint112 reserve1)
func (_Univ2pair *Univ2pairFilterer) ParseSync(log types.Log) (*Univ2pairSync, error) {
	event := new(Univ2pairSync)
	if err := _Univ2pair.contract.UnpackLog(event, "Sync", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Univ2pairTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Univ2pair contract.
type Univ2pairTransferIterator struct {
	Event *Univ2pairTransfer // Event containing the contract specifics and raw log

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
func (it *Univ2pairTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Univ2pairTransfer)
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
		it.Event = new(Univ2pairTransfer)
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
func (it *Univ2pairTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Univ2pairTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Univ2pairTransfer represents a Transfer event raised by the Univ2pair contract.
type Univ2pairTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Univ2pair *Univ2pairFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*Univ2pairTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Univ2pair.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Univ2pairTransferIterator{contract: _Univ2pair.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Univ2pair *Univ2pairFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Univ2pairTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Univ2pair.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Univ2pairTransfer)
				if err := _Univ2pair.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Univ2pair *Univ2pairFilterer) ParseTransfer(log types.Log) (*Univ2pairTransfer, error) {
	event := new(Univ2pairTransfer)
	if err := _Univ2pair.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
