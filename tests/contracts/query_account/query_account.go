// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package query_account

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

// QueryAccountMetaData contains all meta data concerning the QueryAccount contract.
var QueryAccountMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getCodeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getInternalValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506102c18061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c8063b51c4f9614610043578063c560907814610073578063f8b2cb4f146100a3575b5f5ffd5b61005d600480360381019061005891906101da565b6100d3565b60405161006a919061021d565b60405180910390f35b61008d600480360381019061008891906101da565b6100e2565b60405161009a919061021d565b60405180910390f35b6100bd60048036038101906100b891906101da565b61015c565b6040516100ca919061021d565b60405180910390f35b5f5f823b905080915050919050565b5f5f8290508073ffffffffffffffffffffffffffffffffffffffff1663209652556040518163ffffffff1660e01b8152600401602060405180830381865afa158015610130573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101549190610260565b915050919050565b5f8173ffffffffffffffffffffffffffffffffffffffff16319050919050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101a982610180565b9050919050565b6101b98161019f565b81146101c3575f5ffd5b50565b5f813590506101d4816101b0565b92915050565b5f602082840312156101ef576101ee61017c565b5b5f6101fc848285016101c6565b91505092915050565b5f819050919050565b61021781610205565b82525050565b5f6020820190506102305f83018461020e565b92915050565b61023f81610205565b8114610249575f5ffd5b50565b5f8151905061025a81610236565b92915050565b5f602082840312156102755761027461017c565b5b5f6102828482850161024c565b9150509291505056fea2646970667358221220ea08a8667ba721a668c7a303fb36d7ec876636130a31856dd36c9daf57da77b764736f6c634300081c0033",
}

// QueryAccountABI is the input ABI used to generate the binding from.
// Deprecated: Use QueryAccountMetaData.ABI instead.
var QueryAccountABI = QueryAccountMetaData.ABI

// QueryAccountBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use QueryAccountMetaData.Bin instead.
var QueryAccountBin = QueryAccountMetaData.Bin

// DeployQueryAccount deploys a new Ethereum contract, binding an instance of QueryAccount to it.
func DeployQueryAccount(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *QueryAccount, error) {
	parsed, err := QueryAccountMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(QueryAccountBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &QueryAccount{QueryAccountCaller: QueryAccountCaller{contract: contract}, QueryAccountTransactor: QueryAccountTransactor{contract: contract}, QueryAccountFilterer: QueryAccountFilterer{contract: contract}}, nil
}

// QueryAccount is an auto generated Go binding around an Ethereum contract.
type QueryAccount struct {
	QueryAccountCaller     // Read-only binding to the contract
	QueryAccountTransactor // Write-only binding to the contract
	QueryAccountFilterer   // Log filterer for contract events
}

// QueryAccountCaller is an auto generated read-only Go binding around an Ethereum contract.
type QueryAccountCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QueryAccountTransactor is an auto generated write-only Go binding around an Ethereum contract.
type QueryAccountTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QueryAccountFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QueryAccountFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QueryAccountSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QueryAccountSession struct {
	Contract     *QueryAccount     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QueryAccountCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QueryAccountCallerSession struct {
	Contract *QueryAccountCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// QueryAccountTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QueryAccountTransactorSession struct {
	Contract     *QueryAccountTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// QueryAccountRaw is an auto generated low-level Go binding around an Ethereum contract.
type QueryAccountRaw struct {
	Contract *QueryAccount // Generic contract binding to access the raw methods on
}

// QueryAccountCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QueryAccountCallerRaw struct {
	Contract *QueryAccountCaller // Generic read-only contract binding to access the raw methods on
}

// QueryAccountTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QueryAccountTransactorRaw struct {
	Contract *QueryAccountTransactor // Generic write-only contract binding to access the raw methods on
}

// NewQueryAccount creates a new instance of QueryAccount, bound to a specific deployed contract.
func NewQueryAccount(address common.Address, backend bind.ContractBackend) (*QueryAccount, error) {
	contract, err := bindQueryAccount(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &QueryAccount{QueryAccountCaller: QueryAccountCaller{contract: contract}, QueryAccountTransactor: QueryAccountTransactor{contract: contract}, QueryAccountFilterer: QueryAccountFilterer{contract: contract}}, nil
}

// NewQueryAccountCaller creates a new read-only instance of QueryAccount, bound to a specific deployed contract.
func NewQueryAccountCaller(address common.Address, caller bind.ContractCaller) (*QueryAccountCaller, error) {
	contract, err := bindQueryAccount(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QueryAccountCaller{contract: contract}, nil
}

// NewQueryAccountTransactor creates a new write-only instance of QueryAccount, bound to a specific deployed contract.
func NewQueryAccountTransactor(address common.Address, transactor bind.ContractTransactor) (*QueryAccountTransactor, error) {
	contract, err := bindQueryAccount(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QueryAccountTransactor{contract: contract}, nil
}

// NewQueryAccountFilterer creates a new log filterer instance of QueryAccount, bound to a specific deployed contract.
func NewQueryAccountFilterer(address common.Address, filterer bind.ContractFilterer) (*QueryAccountFilterer, error) {
	contract, err := bindQueryAccount(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QueryAccountFilterer{contract: contract}, nil
}

// bindQueryAccount binds a generic wrapper to an already deployed contract.
func bindQueryAccount(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := QueryAccountMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QueryAccount *QueryAccountRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QueryAccount.Contract.QueryAccountCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QueryAccount *QueryAccountRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QueryAccount.Contract.QueryAccountTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QueryAccount *QueryAccountRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QueryAccount.Contract.QueryAccountTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QueryAccount *QueryAccountCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QueryAccount.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QueryAccount *QueryAccountTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QueryAccount.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QueryAccount *QueryAccountTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QueryAccount.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_QueryAccount *QueryAccountCaller) GetBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _QueryAccount.contract.Call(opts, &out, "getBalance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_QueryAccount *QueryAccountSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _QueryAccount.Contract.GetBalance(&_QueryAccount.CallOpts, addr)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_QueryAccount *QueryAccountCallerSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _QueryAccount.Contract.GetBalance(&_QueryAccount.CallOpts, addr)
}

// GetCodeSize is a free data retrieval call binding the contract method 0xb51c4f96.
//
// Solidity: function getCodeSize(address addr) view returns(uint256)
func (_QueryAccount *QueryAccountCaller) GetCodeSize(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _QueryAccount.contract.Call(opts, &out, "getCodeSize", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCodeSize is a free data retrieval call binding the contract method 0xb51c4f96.
//
// Solidity: function getCodeSize(address addr) view returns(uint256)
func (_QueryAccount *QueryAccountSession) GetCodeSize(addr common.Address) (*big.Int, error) {
	return _QueryAccount.Contract.GetCodeSize(&_QueryAccount.CallOpts, addr)
}

// GetCodeSize is a free data retrieval call binding the contract method 0xb51c4f96.
//
// Solidity: function getCodeSize(address addr) view returns(uint256)
func (_QueryAccount *QueryAccountCallerSession) GetCodeSize(addr common.Address) (*big.Int, error) {
	return _QueryAccount.Contract.GetCodeSize(&_QueryAccount.CallOpts, addr)
}

// GetInternalValue is a free data retrieval call binding the contract method 0xc5609078.
//
// Solidity: function getInternalValue(address account) view returns(uint256)
func (_QueryAccount *QueryAccountCaller) GetInternalValue(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _QueryAccount.contract.Call(opts, &out, "getInternalValue", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetInternalValue is a free data retrieval call binding the contract method 0xc5609078.
//
// Solidity: function getInternalValue(address account) view returns(uint256)
func (_QueryAccount *QueryAccountSession) GetInternalValue(account common.Address) (*big.Int, error) {
	return _QueryAccount.Contract.GetInternalValue(&_QueryAccount.CallOpts, account)
}

// GetInternalValue is a free data retrieval call binding the contract method 0xc5609078.
//
// Solidity: function getInternalValue(address account) view returns(uint256)
func (_QueryAccount *QueryAccountCallerSession) GetInternalValue(account common.Address) (*big.Int, error) {
	return _QueryAccount.Contract.GetInternalValue(&_QueryAccount.CallOpts, account)
}
