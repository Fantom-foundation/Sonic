// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package invalidstart

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

// InvalidstartMetaData contains all meta data concerning the Invalidstart contract.
var InvalidstartMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"create2WithInvalidCode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createEmptyContractAndTransferToIt\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createWithInvalidCode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b5060e58061001b5f395ff3fe6080604052348015600e575f5ffd5b5060043610603a575f3560e01c80630e655a6d14603e578063a685771c146046578063dfb7eeb314604e575b5f5ffd5b60446056565b005b604c6068565b005b60546079565b005b5f355f525f365f5ff56066575f5ffd5b565b5f355f52365f5ff06077575f5ffd5b565b5f60405180602001604052805f81525090505f5f8251602084015ff091505f5f5f5f6001865af190508060aa575f5ffd5b50505056fea264697066735822122043af6541a27f32ceba9b2c4f5016e3f404a4fce52ae54054f337a0af682ae96b64736f6c634300081c0033",
}

// InvalidstartABI is the input ABI used to generate the binding from.
// Deprecated: Use InvalidstartMetaData.ABI instead.
var InvalidstartABI = InvalidstartMetaData.ABI

// InvalidstartBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use InvalidstartMetaData.Bin instead.
var InvalidstartBin = InvalidstartMetaData.Bin

// DeployInvalidstart deploys a new Ethereum contract, binding an instance of Invalidstart to it.
func DeployInvalidstart(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Invalidstart, error) {
	parsed, err := InvalidstartMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(InvalidstartBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Invalidstart{InvalidstartCaller: InvalidstartCaller{contract: contract}, InvalidstartTransactor: InvalidstartTransactor{contract: contract}, InvalidstartFilterer: InvalidstartFilterer{contract: contract}}, nil
}

// Invalidstart is an auto generated Go binding around an Ethereum contract.
type Invalidstart struct {
	InvalidstartCaller     // Read-only binding to the contract
	InvalidstartTransactor // Write-only binding to the contract
	InvalidstartFilterer   // Log filterer for contract events
}

// InvalidstartCaller is an auto generated read-only Go binding around an Ethereum contract.
type InvalidstartCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvalidstartTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InvalidstartTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvalidstartFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InvalidstartFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvalidstartSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InvalidstartSession struct {
	Contract     *Invalidstart     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InvalidstartCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InvalidstartCallerSession struct {
	Contract *InvalidstartCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// InvalidstartTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InvalidstartTransactorSession struct {
	Contract     *InvalidstartTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// InvalidstartRaw is an auto generated low-level Go binding around an Ethereum contract.
type InvalidstartRaw struct {
	Contract *Invalidstart // Generic contract binding to access the raw methods on
}

// InvalidstartCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InvalidstartCallerRaw struct {
	Contract *InvalidstartCaller // Generic read-only contract binding to access the raw methods on
}

// InvalidstartTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InvalidstartTransactorRaw struct {
	Contract *InvalidstartTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInvalidstart creates a new instance of Invalidstart, bound to a specific deployed contract.
func NewInvalidstart(address common.Address, backend bind.ContractBackend) (*Invalidstart, error) {
	contract, err := bindInvalidstart(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Invalidstart{InvalidstartCaller: InvalidstartCaller{contract: contract}, InvalidstartTransactor: InvalidstartTransactor{contract: contract}, InvalidstartFilterer: InvalidstartFilterer{contract: contract}}, nil
}

// NewInvalidstartCaller creates a new read-only instance of Invalidstart, bound to a specific deployed contract.
func NewInvalidstartCaller(address common.Address, caller bind.ContractCaller) (*InvalidstartCaller, error) {
	contract, err := bindInvalidstart(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InvalidstartCaller{contract: contract}, nil
}

// NewInvalidstartTransactor creates a new write-only instance of Invalidstart, bound to a specific deployed contract.
func NewInvalidstartTransactor(address common.Address, transactor bind.ContractTransactor) (*InvalidstartTransactor, error) {
	contract, err := bindInvalidstart(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InvalidstartTransactor{contract: contract}, nil
}

// NewInvalidstartFilterer creates a new log filterer instance of Invalidstart, bound to a specific deployed contract.
func NewInvalidstartFilterer(address common.Address, filterer bind.ContractFilterer) (*InvalidstartFilterer, error) {
	contract, err := bindInvalidstart(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InvalidstartFilterer{contract: contract}, nil
}

// bindInvalidstart binds a generic wrapper to an already deployed contract.
func bindInvalidstart(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InvalidstartMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Invalidstart *InvalidstartRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Invalidstart.Contract.InvalidstartCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Invalidstart *InvalidstartRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invalidstart.Contract.InvalidstartTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Invalidstart *InvalidstartRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Invalidstart.Contract.InvalidstartTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Invalidstart *InvalidstartCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Invalidstart.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Invalidstart *InvalidstartTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invalidstart.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Invalidstart *InvalidstartTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Invalidstart.Contract.contract.Transact(opts, method, params...)
}

// Create2WithInvalidCode is a paid mutator transaction binding the contract method 0x0e655a6d.
//
// Solidity: function create2WithInvalidCode() returns()
func (_Invalidstart *InvalidstartTransactor) Create2WithInvalidCode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invalidstart.contract.Transact(opts, "create2WithInvalidCode")
}

// Create2WithInvalidCode is a paid mutator transaction binding the contract method 0x0e655a6d.
//
// Solidity: function create2WithInvalidCode() returns()
func (_Invalidstart *InvalidstartSession) Create2WithInvalidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.Create2WithInvalidCode(&_Invalidstart.TransactOpts)
}

// Create2WithInvalidCode is a paid mutator transaction binding the contract method 0x0e655a6d.
//
// Solidity: function create2WithInvalidCode() returns()
func (_Invalidstart *InvalidstartTransactorSession) Create2WithInvalidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.Create2WithInvalidCode(&_Invalidstart.TransactOpts)
}

// CreateEmptyContractAndTransferToIt is a paid mutator transaction binding the contract method 0xdfb7eeb3.
//
// Solidity: function createEmptyContractAndTransferToIt() returns()
func (_Invalidstart *InvalidstartTransactor) CreateEmptyContractAndTransferToIt(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invalidstart.contract.Transact(opts, "createEmptyContractAndTransferToIt")
}

// CreateEmptyContractAndTransferToIt is a paid mutator transaction binding the contract method 0xdfb7eeb3.
//
// Solidity: function createEmptyContractAndTransferToIt() returns()
func (_Invalidstart *InvalidstartSession) CreateEmptyContractAndTransferToIt() (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateEmptyContractAndTransferToIt(&_Invalidstart.TransactOpts)
}

// CreateEmptyContractAndTransferToIt is a paid mutator transaction binding the contract method 0xdfb7eeb3.
//
// Solidity: function createEmptyContractAndTransferToIt() returns()
func (_Invalidstart *InvalidstartTransactorSession) CreateEmptyContractAndTransferToIt() (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateEmptyContractAndTransferToIt(&_Invalidstart.TransactOpts)
}

// CreateWithInvalidCode is a paid mutator transaction binding the contract method 0xa685771c.
//
// Solidity: function createWithInvalidCode() returns()
func (_Invalidstart *InvalidstartTransactor) CreateWithInvalidCode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invalidstart.contract.Transact(opts, "createWithInvalidCode")
}

// CreateWithInvalidCode is a paid mutator transaction binding the contract method 0xa685771c.
//
// Solidity: function createWithInvalidCode() returns()
func (_Invalidstart *InvalidstartSession) CreateWithInvalidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateWithInvalidCode(&_Invalidstart.TransactOpts)
}

// CreateWithInvalidCode is a paid mutator transaction binding the contract method 0xa685771c.
//
// Solidity: function createWithInvalidCode() returns()
func (_Invalidstart *InvalidstartTransactorSession) CreateWithInvalidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateWithInvalidCode(&_Invalidstart.TransactOpts)
}
