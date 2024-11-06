// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package initcode

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

// InitcodeMetaData contains all meta data concerning the Initcode contract.
var InitcodeMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cost\",\"type\":\"uint256\"}],\"name\":\"LogCost\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"codeSize\",\"type\":\"uint256\"}],\"name\":\"Create2ContractWith\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"codeSize\",\"type\":\"uint256\"}],\"name\":\"CreateContractWith\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"codeSize\",\"type\":\"uint256\"}],\"name\":\"MeasureAssignGasCost\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506102b48061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c80636e63afac146100435780639c07b0fc1461005f578063f220431d1461007b575b5f5ffd5b61005d600480360381019061005891906101cb565b610097565b005b610079600480360381019061007491906101cb565b6100ec565b005b610095600480360381019061009091906101cb565b61013e565b005b5f5a90505f5f9050825f5ff090507fc3263769c8cc487b9d31817141d5ff9ce159867184d644cc2b238d6c54619fdf5a836100d29190610223565b6040516100df9190610265565b60405180910390a1505050565b5f5a90505f5f90508290507fc3263769c8cc487b9d31817141d5ff9ce159867184d644cc2b238d6c54619fdf5a836101249190610223565b6040516101319190610265565b60405180910390a1505050565b5f5a90505f5f90505f835f5ff590507fc3263769c8cc487b9d31817141d5ff9ce159867184d644cc2b238d6c54619fdf5a8361017a9190610223565b6040516101879190610265565b60405180910390a1505050565b5f5ffd5b5f819050919050565b6101aa81610198565b81146101b4575f5ffd5b50565b5f813590506101c5816101a1565b92915050565b5f602082840312156101e0576101df610194565b5b5f6101ed848285016101b7565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61022d82610198565b915061023883610198565b92508282039050818111156102505761024f6101f6565b5b92915050565b61025f81610198565b82525050565b5f6020820190506102785f830184610256565b9291505056fea26469706673582212207ce7bb87574a43b41161b21d4e055c650817dbe5c679072faab2d2b290c0997164736f6c634300081c0033",
}

// InitcodeABI is the input ABI used to generate the binding from.
// Deprecated: Use InitcodeMetaData.ABI instead.
var InitcodeABI = InitcodeMetaData.ABI

// InitcodeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use InitcodeMetaData.Bin instead.
var InitcodeBin = InitcodeMetaData.Bin

// DeployInitcode deploys a new Ethereum contract, binding an instance of Initcode to it.
func DeployInitcode(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Initcode, error) {
	parsed, err := InitcodeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(InitcodeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Initcode{InitcodeCaller: InitcodeCaller{contract: contract}, InitcodeTransactor: InitcodeTransactor{contract: contract}, InitcodeFilterer: InitcodeFilterer{contract: contract}}, nil
}

// Initcode is an auto generated Go binding around an Ethereum contract.
type Initcode struct {
	InitcodeCaller     // Read-only binding to the contract
	InitcodeTransactor // Write-only binding to the contract
	InitcodeFilterer   // Log filterer for contract events
}

// InitcodeCaller is an auto generated read-only Go binding around an Ethereum contract.
type InitcodeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InitcodeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InitcodeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InitcodeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InitcodeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InitcodeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InitcodeSession struct {
	Contract     *Initcode         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InitcodeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InitcodeCallerSession struct {
	Contract *InitcodeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// InitcodeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InitcodeTransactorSession struct {
	Contract     *InitcodeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// InitcodeRaw is an auto generated low-level Go binding around an Ethereum contract.
type InitcodeRaw struct {
	Contract *Initcode // Generic contract binding to access the raw methods on
}

// InitcodeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InitcodeCallerRaw struct {
	Contract *InitcodeCaller // Generic read-only contract binding to access the raw methods on
}

// InitcodeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InitcodeTransactorRaw struct {
	Contract *InitcodeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInitcode creates a new instance of Initcode, bound to a specific deployed contract.
func NewInitcode(address common.Address, backend bind.ContractBackend) (*Initcode, error) {
	contract, err := bindInitcode(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Initcode{InitcodeCaller: InitcodeCaller{contract: contract}, InitcodeTransactor: InitcodeTransactor{contract: contract}, InitcodeFilterer: InitcodeFilterer{contract: contract}}, nil
}

// NewInitcodeCaller creates a new read-only instance of Initcode, bound to a specific deployed contract.
func NewInitcodeCaller(address common.Address, caller bind.ContractCaller) (*InitcodeCaller, error) {
	contract, err := bindInitcode(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InitcodeCaller{contract: contract}, nil
}

// NewInitcodeTransactor creates a new write-only instance of Initcode, bound to a specific deployed contract.
func NewInitcodeTransactor(address common.Address, transactor bind.ContractTransactor) (*InitcodeTransactor, error) {
	contract, err := bindInitcode(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InitcodeTransactor{contract: contract}, nil
}

// NewInitcodeFilterer creates a new log filterer instance of Initcode, bound to a specific deployed contract.
func NewInitcodeFilterer(address common.Address, filterer bind.ContractFilterer) (*InitcodeFilterer, error) {
	contract, err := bindInitcode(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InitcodeFilterer{contract: contract}, nil
}

// bindInitcode binds a generic wrapper to an already deployed contract.
func bindInitcode(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InitcodeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Initcode *InitcodeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Initcode.Contract.InitcodeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Initcode *InitcodeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Initcode.Contract.InitcodeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Initcode *InitcodeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Initcode.Contract.InitcodeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Initcode *InitcodeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Initcode.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Initcode *InitcodeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Initcode.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Initcode *InitcodeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Initcode.Contract.contract.Transact(opts, method, params...)
}

// Create2ContractWith is a paid mutator transaction binding the contract method 0xf220431d.
//
// Solidity: function Create2ContractWith(uint256 codeSize) returns()
func (_Initcode *InitcodeTransactor) Create2ContractWith(opts *bind.TransactOpts, codeSize *big.Int) (*types.Transaction, error) {
	return _Initcode.contract.Transact(opts, "Create2ContractWith", codeSize)
}

// Create2ContractWith is a paid mutator transaction binding the contract method 0xf220431d.
//
// Solidity: function Create2ContractWith(uint256 codeSize) returns()
func (_Initcode *InitcodeSession) Create2ContractWith(codeSize *big.Int) (*types.Transaction, error) {
	return _Initcode.Contract.Create2ContractWith(&_Initcode.TransactOpts, codeSize)
}

// Create2ContractWith is a paid mutator transaction binding the contract method 0xf220431d.
//
// Solidity: function Create2ContractWith(uint256 codeSize) returns()
func (_Initcode *InitcodeTransactorSession) Create2ContractWith(codeSize *big.Int) (*types.Transaction, error) {
	return _Initcode.Contract.Create2ContractWith(&_Initcode.TransactOpts, codeSize)
}

// CreateContractWith is a paid mutator transaction binding the contract method 0x6e63afac.
//
// Solidity: function CreateContractWith(uint256 codeSize) returns()
func (_Initcode *InitcodeTransactor) CreateContractWith(opts *bind.TransactOpts, codeSize *big.Int) (*types.Transaction, error) {
	return _Initcode.contract.Transact(opts, "CreateContractWith", codeSize)
}

// CreateContractWith is a paid mutator transaction binding the contract method 0x6e63afac.
//
// Solidity: function CreateContractWith(uint256 codeSize) returns()
func (_Initcode *InitcodeSession) CreateContractWith(codeSize *big.Int) (*types.Transaction, error) {
	return _Initcode.Contract.CreateContractWith(&_Initcode.TransactOpts, codeSize)
}

// CreateContractWith is a paid mutator transaction binding the contract method 0x6e63afac.
//
// Solidity: function CreateContractWith(uint256 codeSize) returns()
func (_Initcode *InitcodeTransactorSession) CreateContractWith(codeSize *big.Int) (*types.Transaction, error) {
	return _Initcode.Contract.CreateContractWith(&_Initcode.TransactOpts, codeSize)
}

// MeasureAssignGasCost is a paid mutator transaction binding the contract method 0x9c07b0fc.
//
// Solidity: function MeasureAssignGasCost(uint256 codeSize) returns()
func (_Initcode *InitcodeTransactor) MeasureAssignGasCost(opts *bind.TransactOpts, codeSize *big.Int) (*types.Transaction, error) {
	return _Initcode.contract.Transact(opts, "MeasureAssignGasCost", codeSize)
}

// MeasureAssignGasCost is a paid mutator transaction binding the contract method 0x9c07b0fc.
//
// Solidity: function MeasureAssignGasCost(uint256 codeSize) returns()
func (_Initcode *InitcodeSession) MeasureAssignGasCost(codeSize *big.Int) (*types.Transaction, error) {
	return _Initcode.Contract.MeasureAssignGasCost(&_Initcode.TransactOpts, codeSize)
}

// MeasureAssignGasCost is a paid mutator transaction binding the contract method 0x9c07b0fc.
//
// Solidity: function MeasureAssignGasCost(uint256 codeSize) returns()
func (_Initcode *InitcodeTransactorSession) MeasureAssignGasCost(codeSize *big.Int) (*types.Transaction, error) {
	return _Initcode.Contract.MeasureAssignGasCost(&_Initcode.TransactOpts, codeSize)
}

// InitcodeLogCostIterator is returned from FilterLogCost and is used to iterate over the raw logs and unpacked data for LogCost events raised by the Initcode contract.
type InitcodeLogCostIterator struct {
	Event *InitcodeLogCost // Event containing the contract specifics and raw log

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
func (it *InitcodeLogCostIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InitcodeLogCost)
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
		it.Event = new(InitcodeLogCost)
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
func (it *InitcodeLogCostIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InitcodeLogCostIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InitcodeLogCost represents a LogCost event raised by the Initcode contract.
type InitcodeLogCost struct {
	Cost *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogCost is a free log retrieval operation binding the contract event 0xc3263769c8cc487b9d31817141d5ff9ce159867184d644cc2b238d6c54619fdf.
//
// Solidity: event LogCost(uint256 cost)
func (_Initcode *InitcodeFilterer) FilterLogCost(opts *bind.FilterOpts) (*InitcodeLogCostIterator, error) {

	logs, sub, err := _Initcode.contract.FilterLogs(opts, "LogCost")
	if err != nil {
		return nil, err
	}
	return &InitcodeLogCostIterator{contract: _Initcode.contract, event: "LogCost", logs: logs, sub: sub}, nil
}

// WatchLogCost is a free log subscription operation binding the contract event 0xc3263769c8cc487b9d31817141d5ff9ce159867184d644cc2b238d6c54619fdf.
//
// Solidity: event LogCost(uint256 cost)
func (_Initcode *InitcodeFilterer) WatchLogCost(opts *bind.WatchOpts, sink chan<- *InitcodeLogCost) (event.Subscription, error) {

	logs, sub, err := _Initcode.contract.WatchLogs(opts, "LogCost")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InitcodeLogCost)
				if err := _Initcode.contract.UnpackLog(event, "LogCost", log); err != nil {
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

// ParseLogCost is a log parse operation binding the contract event 0xc3263769c8cc487b9d31817141d5ff9ce159867184d644cc2b238d6c54619fdf.
//
// Solidity: event LogCost(uint256 cost)
func (_Initcode *InitcodeFilterer) ParseLogCost(log types.Log) (*InitcodeLogCost, error) {
	event := new(InitcodeLogCost)
	if err := _Initcode.contract.UnpackLog(event, "LogCost", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
