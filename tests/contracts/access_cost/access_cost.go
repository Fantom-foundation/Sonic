// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package accessCost

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

// AccessCostMetaData contains all meta data concerning the AccessCost contract.
var AccessCostMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cost\",\"type\":\"uint256\"}],\"name\":\"LogCost\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getCoinBaseAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOrigin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"touchAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"touchCoinBase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"touchOrigin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506102db8061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610055575f3560e01c80638102f66114610059578063d847a22b14610063578063df1f29ee1461007f578063ed9435631461009d578063f3c4df1a146100bb575b5f5ffd5b6100616100c5565b005b61007d600480360381019061007891906101c1565b6100d7565b005b610087610143565b60405161009491906101fb565b60405180910390f35b6100a561014a565b6040516100b291906101fb565b60405180910390f35b6100c3610151565b005b6100d56100d0610143565b6100d7565b565b5f5a90505f8273ffffffffffffffffffffffffffffffffffffffff163190507fc3263769c8cc487b9d31817141d5ff9ce159867184d644cc2b238d6c54619fdf5a83610123919061024a565b604051610130919061028c565b60405180910390a1805f81905550505050565b5f32905090565b5f41905090565b61016161015c61014a565b6100d7565b565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61019082610167565b9050919050565b6101a081610186565b81146101aa575f5ffd5b50565b5f813590506101bb81610197565b92915050565b5f602082840312156101d6576101d5610163565b5b5f6101e3848285016101ad565b91505092915050565b6101f581610186565b82525050565b5f60208201905061020e5f8301846101ec565b92915050565b5f819050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61025482610214565b915061025f83610214565b92508282039050818111156102775761027661021d565b5b92915050565b61028681610214565b82525050565b5f60208201905061029f5f83018461027d565b9291505056fea264697066735822122073c49b26d3bd90cd612bebdbb14dd83319893649247f63113daee7d763443afd64736f6c634300081c0033",
}

// AccessCostABI is the input ABI used to generate the binding from.
// Deprecated: Use AccessCostMetaData.ABI instead.
var AccessCostABI = AccessCostMetaData.ABI

// AccessCostBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AccessCostMetaData.Bin instead.
var AccessCostBin = AccessCostMetaData.Bin

// DeployAccessCost deploys a new Ethereum contract, binding an instance of AccessCost to it.
func DeployAccessCost(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AccessCost, error) {
	parsed, err := AccessCostMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AccessCostBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccessCost{AccessCostCaller: AccessCostCaller{contract: contract}, AccessCostTransactor: AccessCostTransactor{contract: contract}, AccessCostFilterer: AccessCostFilterer{contract: contract}}, nil
}

// AccessCost is an auto generated Go binding around an Ethereum contract.
type AccessCost struct {
	AccessCostCaller     // Read-only binding to the contract
	AccessCostTransactor // Write-only binding to the contract
	AccessCostFilterer   // Log filterer for contract events
}

// AccessCostCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccessCostCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessCostTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccessCostTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessCostFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccessCostFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessCostSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccessCostSession struct {
	Contract     *AccessCost       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccessCostCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccessCostCallerSession struct {
	Contract *AccessCostCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AccessCostTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccessCostTransactorSession struct {
	Contract     *AccessCostTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AccessCostRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccessCostRaw struct {
	Contract *AccessCost // Generic contract binding to access the raw methods on
}

// AccessCostCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccessCostCallerRaw struct {
	Contract *AccessCostCaller // Generic read-only contract binding to access the raw methods on
}

// AccessCostTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccessCostTransactorRaw struct {
	Contract *AccessCostTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccessCost creates a new instance of AccessCost, bound to a specific deployed contract.
func NewAccessCost(address common.Address, backend bind.ContractBackend) (*AccessCost, error) {
	contract, err := bindAccessCost(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessCost{AccessCostCaller: AccessCostCaller{contract: contract}, AccessCostTransactor: AccessCostTransactor{contract: contract}, AccessCostFilterer: AccessCostFilterer{contract: contract}}, nil
}

// NewAccessCostCaller creates a new read-only instance of AccessCost, bound to a specific deployed contract.
func NewAccessCostCaller(address common.Address, caller bind.ContractCaller) (*AccessCostCaller, error) {
	contract, err := bindAccessCost(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessCostCaller{contract: contract}, nil
}

// NewAccessCostTransactor creates a new write-only instance of AccessCost, bound to a specific deployed contract.
func NewAccessCostTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessCostTransactor, error) {
	contract, err := bindAccessCost(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessCostTransactor{contract: contract}, nil
}

// NewAccessCostFilterer creates a new log filterer instance of AccessCost, bound to a specific deployed contract.
func NewAccessCostFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessCostFilterer, error) {
	contract, err := bindAccessCost(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessCostFilterer{contract: contract}, nil
}

// bindAccessCost binds a generic wrapper to an already deployed contract.
func bindAccessCost(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AccessCostMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessCost *AccessCostRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessCost.Contract.AccessCostCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessCost *AccessCostRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessCost.Contract.AccessCostTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessCost *AccessCostRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessCost.Contract.AccessCostTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessCost *AccessCostCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessCost.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessCost *AccessCostTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessCost.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessCost *AccessCostTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessCost.Contract.contract.Transact(opts, method, params...)
}

// GetCoinBaseAddress is a free data retrieval call binding the contract method 0xed943563.
//
// Solidity: function getCoinBaseAddress() view returns(address)
func (_AccessCost *AccessCostCaller) GetCoinBaseAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessCost.contract.Call(opts, &out, "getCoinBaseAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetCoinBaseAddress is a free data retrieval call binding the contract method 0xed943563.
//
// Solidity: function getCoinBaseAddress() view returns(address)
func (_AccessCost *AccessCostSession) GetCoinBaseAddress() (common.Address, error) {
	return _AccessCost.Contract.GetCoinBaseAddress(&_AccessCost.CallOpts)
}

// GetCoinBaseAddress is a free data retrieval call binding the contract method 0xed943563.
//
// Solidity: function getCoinBaseAddress() view returns(address)
func (_AccessCost *AccessCostCallerSession) GetCoinBaseAddress() (common.Address, error) {
	return _AccessCost.Contract.GetCoinBaseAddress(&_AccessCost.CallOpts)
}

// GetOrigin is a free data retrieval call binding the contract method 0xdf1f29ee.
//
// Solidity: function getOrigin() view returns(address)
func (_AccessCost *AccessCostCaller) GetOrigin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessCost.contract.Call(opts, &out, "getOrigin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOrigin is a free data retrieval call binding the contract method 0xdf1f29ee.
//
// Solidity: function getOrigin() view returns(address)
func (_AccessCost *AccessCostSession) GetOrigin() (common.Address, error) {
	return _AccessCost.Contract.GetOrigin(&_AccessCost.CallOpts)
}

// GetOrigin is a free data retrieval call binding the contract method 0xdf1f29ee.
//
// Solidity: function getOrigin() view returns(address)
func (_AccessCost *AccessCostCallerSession) GetOrigin() (common.Address, error) {
	return _AccessCost.Contract.GetOrigin(&_AccessCost.CallOpts)
}

// TouchAddress is a paid mutator transaction binding the contract method 0xd847a22b.
//
// Solidity: function touchAddress(address addr) returns()
func (_AccessCost *AccessCostTransactor) TouchAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _AccessCost.contract.Transact(opts, "touchAddress", addr)
}

// TouchAddress is a paid mutator transaction binding the contract method 0xd847a22b.
//
// Solidity: function touchAddress(address addr) returns()
func (_AccessCost *AccessCostSession) TouchAddress(addr common.Address) (*types.Transaction, error) {
	return _AccessCost.Contract.TouchAddress(&_AccessCost.TransactOpts, addr)
}

// TouchAddress is a paid mutator transaction binding the contract method 0xd847a22b.
//
// Solidity: function touchAddress(address addr) returns()
func (_AccessCost *AccessCostTransactorSession) TouchAddress(addr common.Address) (*types.Transaction, error) {
	return _AccessCost.Contract.TouchAddress(&_AccessCost.TransactOpts, addr)
}

// TouchCoinBase is a paid mutator transaction binding the contract method 0xf3c4df1a.
//
// Solidity: function touchCoinBase() returns()
func (_AccessCost *AccessCostTransactor) TouchCoinBase(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessCost.contract.Transact(opts, "touchCoinBase")
}

// TouchCoinBase is a paid mutator transaction binding the contract method 0xf3c4df1a.
//
// Solidity: function touchCoinBase() returns()
func (_AccessCost *AccessCostSession) TouchCoinBase() (*types.Transaction, error) {
	return _AccessCost.Contract.TouchCoinBase(&_AccessCost.TransactOpts)
}

// TouchCoinBase is a paid mutator transaction binding the contract method 0xf3c4df1a.
//
// Solidity: function touchCoinBase() returns()
func (_AccessCost *AccessCostTransactorSession) TouchCoinBase() (*types.Transaction, error) {
	return _AccessCost.Contract.TouchCoinBase(&_AccessCost.TransactOpts)
}

// TouchOrigin is a paid mutator transaction binding the contract method 0x8102f661.
//
// Solidity: function touchOrigin() returns()
func (_AccessCost *AccessCostTransactor) TouchOrigin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessCost.contract.Transact(opts, "touchOrigin")
}

// TouchOrigin is a paid mutator transaction binding the contract method 0x8102f661.
//
// Solidity: function touchOrigin() returns()
func (_AccessCost *AccessCostSession) TouchOrigin() (*types.Transaction, error) {
	return _AccessCost.Contract.TouchOrigin(&_AccessCost.TransactOpts)
}

// TouchOrigin is a paid mutator transaction binding the contract method 0x8102f661.
//
// Solidity: function touchOrigin() returns()
func (_AccessCost *AccessCostTransactorSession) TouchOrigin() (*types.Transaction, error) {
	return _AccessCost.Contract.TouchOrigin(&_AccessCost.TransactOpts)
}

// AccessCostLogCostIterator is returned from FilterLogCost and is used to iterate over the raw logs and unpacked data for LogCost events raised by the AccessCost contract.
type AccessCostLogCostIterator struct {
	Event *AccessCostLogCost // Event containing the contract specifics and raw log

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
func (it *AccessCostLogCostIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessCostLogCost)
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
		it.Event = new(AccessCostLogCost)
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
func (it *AccessCostLogCostIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessCostLogCostIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessCostLogCost represents a LogCost event raised by the AccessCost contract.
type AccessCostLogCost struct {
	Cost *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogCost is a free log retrieval operation binding the contract event 0xc3263769c8cc487b9d31817141d5ff9ce159867184d644cc2b238d6c54619fdf.
//
// Solidity: event LogCost(uint256 cost)
func (_AccessCost *AccessCostFilterer) FilterLogCost(opts *bind.FilterOpts) (*AccessCostLogCostIterator, error) {

	logs, sub, err := _AccessCost.contract.FilterLogs(opts, "LogCost")
	if err != nil {
		return nil, err
	}
	return &AccessCostLogCostIterator{contract: _AccessCost.contract, event: "LogCost", logs: logs, sub: sub}, nil
}

// WatchLogCost is a free log subscription operation binding the contract event 0xc3263769c8cc487b9d31817141d5ff9ce159867184d644cc2b238d6c54619fdf.
//
// Solidity: event LogCost(uint256 cost)
func (_AccessCost *AccessCostFilterer) WatchLogCost(opts *bind.WatchOpts, sink chan<- *AccessCostLogCost) (event.Subscription, error) {

	logs, sub, err := _AccessCost.contract.WatchLogs(opts, "LogCost")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessCostLogCost)
				if err := _AccessCost.contract.UnpackLog(event, "LogCost", log); err != nil {
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
func (_AccessCost *AccessCostFilterer) ParseLogCost(log types.Log) (*AccessCostLogCost, error) {
	event := new(AccessCostLogCost)
	if err := _AccessCost.contract.UnpackLog(event, "LogCost", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
