// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package coinbase

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

// CoinbaseMetaData contains all meta data concerning the Coinbase contract.
var CoinbaseMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fee\",\"type\":\"address\"}],\"name\":\"LogCoinbase\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getCoinbase\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"logCoinbase\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506101a28061001c5f395ff3fe608060405260043610610028575f3560e01c80639c945ef11461002c578063d1a82a9d14610036575b5f5ffd5b610034610060565b005b348015610041575f5ffd5b5061004a610099565b60405161005791906100df565b60405180910390f35b7f3c7a1a2619f5180be3c6c302a3b4ec2f9051a356c12c319b58b01ce38bcd9c784160405161008f9190610153565b60405180910390a1565b5f41905090565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100c9826100a0565b9050919050565b6100d9816100bf565b82525050565b5f6020820190506100f25f8301846100d0565b92915050565b5f819050919050565b5f61011b610116610111846100a0565b6100f8565b6100a0565b9050919050565b5f61012c82610101565b9050919050565b5f61013d82610122565b9050919050565b61014d81610133565b82525050565b5f6020820190506101665f830184610144565b9291505056fea2646970667358221220dcb064a104ed862d139b1868757a2704ba859692abbe0234e64a31eba57d039864736f6c634300081c0033",
}

// CoinbaseABI is the input ABI used to generate the binding from.
// Deprecated: Use CoinbaseMetaData.ABI instead.
var CoinbaseABI = CoinbaseMetaData.ABI

// CoinbaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CoinbaseMetaData.Bin instead.
var CoinbaseBin = CoinbaseMetaData.Bin

// DeployCoinbase deploys a new Ethereum contract, binding an instance of Coinbase to it.
func DeployCoinbase(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Coinbase, error) {
	parsed, err := CoinbaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CoinbaseBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Coinbase{CoinbaseCaller: CoinbaseCaller{contract: contract}, CoinbaseTransactor: CoinbaseTransactor{contract: contract}, CoinbaseFilterer: CoinbaseFilterer{contract: contract}}, nil
}

// Coinbase is an auto generated Go binding around an Ethereum contract.
type Coinbase struct {
	CoinbaseCaller     // Read-only binding to the contract
	CoinbaseTransactor // Write-only binding to the contract
	CoinbaseFilterer   // Log filterer for contract events
}

// CoinbaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type CoinbaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoinbaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CoinbaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoinbaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CoinbaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoinbaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CoinbaseSession struct {
	Contract     *Coinbase         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CoinbaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CoinbaseCallerSession struct {
	Contract *CoinbaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// CoinbaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CoinbaseTransactorSession struct {
	Contract     *CoinbaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// CoinbaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type CoinbaseRaw struct {
	Contract *Coinbase // Generic contract binding to access the raw methods on
}

// CoinbaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CoinbaseCallerRaw struct {
	Contract *CoinbaseCaller // Generic read-only contract binding to access the raw methods on
}

// CoinbaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CoinbaseTransactorRaw struct {
	Contract *CoinbaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCoinbase creates a new instance of Coinbase, bound to a specific deployed contract.
func NewCoinbase(address common.Address, backend bind.ContractBackend) (*Coinbase, error) {
	contract, err := bindCoinbase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Coinbase{CoinbaseCaller: CoinbaseCaller{contract: contract}, CoinbaseTransactor: CoinbaseTransactor{contract: contract}, CoinbaseFilterer: CoinbaseFilterer{contract: contract}}, nil
}

// NewCoinbaseCaller creates a new read-only instance of Coinbase, bound to a specific deployed contract.
func NewCoinbaseCaller(address common.Address, caller bind.ContractCaller) (*CoinbaseCaller, error) {
	contract, err := bindCoinbase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CoinbaseCaller{contract: contract}, nil
}

// NewCoinbaseTransactor creates a new write-only instance of Coinbase, bound to a specific deployed contract.
func NewCoinbaseTransactor(address common.Address, transactor bind.ContractTransactor) (*CoinbaseTransactor, error) {
	contract, err := bindCoinbase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CoinbaseTransactor{contract: contract}, nil
}

// NewCoinbaseFilterer creates a new log filterer instance of Coinbase, bound to a specific deployed contract.
func NewCoinbaseFilterer(address common.Address, filterer bind.ContractFilterer) (*CoinbaseFilterer, error) {
	contract, err := bindCoinbase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CoinbaseFilterer{contract: contract}, nil
}

// bindCoinbase binds a generic wrapper to an already deployed contract.
func bindCoinbase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CoinbaseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Coinbase *CoinbaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Coinbase.Contract.CoinbaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Coinbase *CoinbaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coinbase.Contract.CoinbaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Coinbase *CoinbaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Coinbase.Contract.CoinbaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Coinbase *CoinbaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Coinbase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Coinbase *CoinbaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coinbase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Coinbase *CoinbaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Coinbase.Contract.contract.Transact(opts, method, params...)
}

// GetCoinbase is a free data retrieval call binding the contract method 0xd1a82a9d.
//
// Solidity: function getCoinbase() view returns(address)
func (_Coinbase *CoinbaseCaller) GetCoinbase(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Coinbase.contract.Call(opts, &out, "getCoinbase")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetCoinbase is a free data retrieval call binding the contract method 0xd1a82a9d.
//
// Solidity: function getCoinbase() view returns(address)
func (_Coinbase *CoinbaseSession) GetCoinbase() (common.Address, error) {
	return _Coinbase.Contract.GetCoinbase(&_Coinbase.CallOpts)
}

// GetCoinbase is a free data retrieval call binding the contract method 0xd1a82a9d.
//
// Solidity: function getCoinbase() view returns(address)
func (_Coinbase *CoinbaseCallerSession) GetCoinbase() (common.Address, error) {
	return _Coinbase.Contract.GetCoinbase(&_Coinbase.CallOpts)
}

// LogCoinbase is a paid mutator transaction binding the contract method 0x9c945ef1.
//
// Solidity: function logCoinbase() payable returns()
func (_Coinbase *CoinbaseTransactor) LogCoinbase(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coinbase.contract.Transact(opts, "logCoinbase")
}

// LogCoinbase is a paid mutator transaction binding the contract method 0x9c945ef1.
//
// Solidity: function logCoinbase() payable returns()
func (_Coinbase *CoinbaseSession) LogCoinbase() (*types.Transaction, error) {
	return _Coinbase.Contract.LogCoinbase(&_Coinbase.TransactOpts)
}

// LogCoinbase is a paid mutator transaction binding the contract method 0x9c945ef1.
//
// Solidity: function logCoinbase() payable returns()
func (_Coinbase *CoinbaseTransactorSession) LogCoinbase() (*types.Transaction, error) {
	return _Coinbase.Contract.LogCoinbase(&_Coinbase.TransactOpts)
}

// CoinbaseLogCoinbaseIterator is returned from FilterLogCoinbase and is used to iterate over the raw logs and unpacked data for LogCoinbase events raised by the Coinbase contract.
type CoinbaseLogCoinbaseIterator struct {
	Event *CoinbaseLogCoinbase // Event containing the contract specifics and raw log

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
func (it *CoinbaseLogCoinbaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoinbaseLogCoinbase)
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
		it.Event = new(CoinbaseLogCoinbase)
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
func (it *CoinbaseLogCoinbaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoinbaseLogCoinbaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoinbaseLogCoinbase represents a LogCoinbase event raised by the Coinbase contract.
type CoinbaseLogCoinbase struct {
	Fee common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogCoinbase is a free log retrieval operation binding the contract event 0x3c7a1a2619f5180be3c6c302a3b4ec2f9051a356c12c319b58b01ce38bcd9c78.
//
// Solidity: event LogCoinbase(address fee)
func (_Coinbase *CoinbaseFilterer) FilterLogCoinbase(opts *bind.FilterOpts) (*CoinbaseLogCoinbaseIterator, error) {

	logs, sub, err := _Coinbase.contract.FilterLogs(opts, "LogCoinbase")
	if err != nil {
		return nil, err
	}
	return &CoinbaseLogCoinbaseIterator{contract: _Coinbase.contract, event: "LogCoinbase", logs: logs, sub: sub}, nil
}

// WatchLogCoinbase is a free log subscription operation binding the contract event 0x3c7a1a2619f5180be3c6c302a3b4ec2f9051a356c12c319b58b01ce38bcd9c78.
//
// Solidity: event LogCoinbase(address fee)
func (_Coinbase *CoinbaseFilterer) WatchLogCoinbase(opts *bind.WatchOpts, sink chan<- *CoinbaseLogCoinbase) (event.Subscription, error) {

	logs, sub, err := _Coinbase.contract.WatchLogs(opts, "LogCoinbase")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoinbaseLogCoinbase)
				if err := _Coinbase.contract.UnpackLog(event, "LogCoinbase", log); err != nil {
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

// ParseLogCoinbase is a log parse operation binding the contract event 0x3c7a1a2619f5180be3c6c302a3b4ec2f9051a356c12c319b58b01ce38bcd9c78.
//
// Solidity: event LogCoinbase(address fee)
func (_Coinbase *CoinbaseFilterer) ParseLogCoinbase(log types.Log) (*CoinbaseLogCoinbase, error) {
	event := new(CoinbaseLogCoinbase)
	if err := _Coinbase.contract.UnpackLog(event, "LogCoinbase", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
