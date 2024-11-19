// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_gas

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

// BurnGasMetaData contains all meta data concerning the BurnGas contract.
var BurnGasMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"burnGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"data\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506101ae8061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610034575f3560e01c8063b9554c5914610038578063f0ba844014610042575b5f5ffd5b610040610072565b005b61005c600480360381019061005791906100f8565b6100a9565b6040516100699190610132565b60405180910390f35b5f5f90505b60c88110156100a657805f8260c881106100945761009361014b565b5b01819055508080600101915050610077565b50565b5f8160c881106100b7575f80fd5b015f915090505481565b5f5ffd5b5f819050919050565b6100d7816100c5565b81146100e1575f5ffd5b50565b5f813590506100f2816100ce565b92915050565b5f6020828403121561010d5761010c6100c1565b5b5f61011a848285016100e4565b91505092915050565b61012c816100c5565b82525050565b5f6020820190506101455f830184610123565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffdfea264697066735822122065607b666318458dcd03758a68a7cf162ae02e9d8076d54ee5bf4007c4b6419a64736f6c634300081c0033",
}

// BurnGasABI is the input ABI used to generate the binding from.
// Deprecated: Use BurnGasMetaData.ABI instead.
var BurnGasABI = BurnGasMetaData.ABI

// BurnGasBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BurnGasMetaData.Bin instead.
var BurnGasBin = BurnGasMetaData.Bin

// DeployBurnGas deploys a new Ethereum contract, binding an instance of BurnGas to it.
func DeployBurnGas(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BurnGas, error) {
	parsed, err := BurnGasMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnGasBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnGas{BurnGasCaller: BurnGasCaller{contract: contract}, BurnGasTransactor: BurnGasTransactor{contract: contract}, BurnGasFilterer: BurnGasFilterer{contract: contract}}, nil
}

// BurnGas is an auto generated Go binding around an Ethereum contract.
type BurnGas struct {
	BurnGasCaller     // Read-only binding to the contract
	BurnGasTransactor // Write-only binding to the contract
	BurnGasFilterer   // Log filterer for contract events
}

// BurnGasCaller is an auto generated read-only Go binding around an Ethereum contract.
type BurnGasCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnGasTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BurnGasTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnGasFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BurnGasFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnGasSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BurnGasSession struct {
	Contract     *BurnGas          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BurnGasCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BurnGasCallerSession struct {
	Contract *BurnGasCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BurnGasTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BurnGasTransactorSession struct {
	Contract     *BurnGasTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BurnGasRaw is an auto generated low-level Go binding around an Ethereum contract.
type BurnGasRaw struct {
	Contract *BurnGas // Generic contract binding to access the raw methods on
}

// BurnGasCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BurnGasCallerRaw struct {
	Contract *BurnGasCaller // Generic read-only contract binding to access the raw methods on
}

// BurnGasTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BurnGasTransactorRaw struct {
	Contract *BurnGasTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBurnGas creates a new instance of BurnGas, bound to a specific deployed contract.
func NewBurnGas(address common.Address, backend bind.ContractBackend) (*BurnGas, error) {
	contract, err := bindBurnGas(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnGas{BurnGasCaller: BurnGasCaller{contract: contract}, BurnGasTransactor: BurnGasTransactor{contract: contract}, BurnGasFilterer: BurnGasFilterer{contract: contract}}, nil
}

// NewBurnGasCaller creates a new read-only instance of BurnGas, bound to a specific deployed contract.
func NewBurnGasCaller(address common.Address, caller bind.ContractCaller) (*BurnGasCaller, error) {
	contract, err := bindBurnGas(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnGasCaller{contract: contract}, nil
}

// NewBurnGasTransactor creates a new write-only instance of BurnGas, bound to a specific deployed contract.
func NewBurnGasTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnGasTransactor, error) {
	contract, err := bindBurnGas(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnGasTransactor{contract: contract}, nil
}

// NewBurnGasFilterer creates a new log filterer instance of BurnGas, bound to a specific deployed contract.
func NewBurnGasFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnGasFilterer, error) {
	contract, err := bindBurnGas(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnGasFilterer{contract: contract}, nil
}

// bindBurnGas binds a generic wrapper to an already deployed contract.
func bindBurnGas(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnGasMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BurnGas *BurnGasRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnGas.Contract.BurnGasCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BurnGas *BurnGasRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnGas.Contract.BurnGasTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BurnGas *BurnGasRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnGas.Contract.BurnGasTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BurnGas *BurnGasCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnGas.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BurnGas *BurnGasTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnGas.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BurnGas *BurnGasTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnGas.Contract.contract.Transact(opts, method, params...)
}

// Data is a free data retrieval call binding the contract method 0xf0ba8440.
//
// Solidity: function data(uint256 ) view returns(uint256)
func (_BurnGas *BurnGasCaller) Data(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BurnGas.contract.Call(opts, &out, "data", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Data is a free data retrieval call binding the contract method 0xf0ba8440.
//
// Solidity: function data(uint256 ) view returns(uint256)
func (_BurnGas *BurnGasSession) Data(arg0 *big.Int) (*big.Int, error) {
	return _BurnGas.Contract.Data(&_BurnGas.CallOpts, arg0)
}

// Data is a free data retrieval call binding the contract method 0xf0ba8440.
//
// Solidity: function data(uint256 ) view returns(uint256)
func (_BurnGas *BurnGasCallerSession) Data(arg0 *big.Int) (*big.Int, error) {
	return _BurnGas.Contract.Data(&_BurnGas.CallOpts, arg0)
}

// BurnGas is a paid mutator transaction binding the contract method 0xb9554c59.
//
// Solidity: function burnGas() returns()
func (_BurnGas *BurnGasTransactor) BurnGas(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnGas.contract.Transact(opts, "burnGas")
}

// BurnGas is a paid mutator transaction binding the contract method 0xb9554c59.
//
// Solidity: function burnGas() returns()
func (_BurnGas *BurnGasSession) BurnGas() (*types.Transaction, error) {
	return _BurnGas.Contract.BurnGas(&_BurnGas.TransactOpts)
}

// BurnGas is a paid mutator transaction binding the contract method 0xb9554c59.
//
// Solidity: function burnGas() returns()
func (_BurnGas *BurnGasTransactorSession) BurnGas() (*types.Transaction, error) {
	return _BurnGas.Contract.BurnGas(&_BurnGas.TransactOpts)
}
