// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package selfdestruct

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

// SelfdestructMetaData contains all meta data concerning the Selfdestruct contract.
var SelfdestructMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"destroyNow\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"selfBeneficiary\",\"type\":\"bool\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"LogAfterDestruct\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"selfBeneficiary\",\"type\":\"bool\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"destroyContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"someData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052607b5f5560405161032d38038061032d833981810160405281019061002991906100fe565b82156100405761003f828261004860201b60201c565b5b50505061014e565b8115610052573090505b8073ffffffffffffffffffffffffffffffffffffffff16ff5b5f5ffd5b5f8115159050919050565b6100838161006f565b811461008d575f5ffd5b50565b5f8151905061009e8161007a565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100cd826100a4565b9050919050565b6100dd816100c3565b81146100e7575f5ffd5b50565b5f815190506100f8816100d4565b92915050565b5f5f5f606084860312156101155761011461006b565b5b5f61012286828701610090565b935050602061013386828701610090565b9250506040610144868287016100ea565b9150509250925092565b6101d28061015b5f395ff3fe608060405234801561000f575f5ffd5b5060043610610034575f3560e01c8063c73b717514610038578063ca29aee914610056575b5f5ffd5b610040610072565b60405161004d91906100b2565b60405180910390f35b610070600480360381019061006b919061015e565b610077565b005b5f5481565b8115610081573090505b8073ffffffffffffffffffffffffffffffffffffffff16ff5b5f819050919050565b6100ac8161009a565b82525050565b5f6020820190506100c55f8301846100a3565b92915050565b5f5ffd5b5f8115159050919050565b6100e3816100cf565b81146100ed575f5ffd5b50565b5f813590506100fe816100da565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61012d82610104565b9050919050565b61013d81610123565b8114610147575f5ffd5b50565b5f8135905061015881610134565b92915050565b5f5f60408385031215610174576101736100cb565b5b5f610181858286016100f0565b92505060206101928582860161014a565b915050925092905056fea2646970667358221220dd7cf5a452513fee71dcab61ad515163ca8c83c6551b9500939b5b843661b32964736f6c634300081c0033",
}

// SelfdestructABI is the input ABI used to generate the binding from.
// Deprecated: Use SelfdestructMetaData.ABI instead.
var SelfdestructABI = SelfdestructMetaData.ABI

// SelfdestructBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SelfdestructMetaData.Bin instead.
var SelfdestructBin = SelfdestructMetaData.Bin

// DeploySelfdestruct deploys a new Ethereum contract, binding an instance of Selfdestruct to it.
func DeploySelfdestruct(auth *bind.TransactOpts, backend bind.ContractBackend, destroyNow bool, selfBeneficiary bool, recipient common.Address) (common.Address, *types.Transaction, *Selfdestruct, error) {
	parsed, err := SelfdestructMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SelfdestructBin), backend, destroyNow, selfBeneficiary, recipient)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Selfdestruct{SelfdestructCaller: SelfdestructCaller{contract: contract}, SelfdestructTransactor: SelfdestructTransactor{contract: contract}, SelfdestructFilterer: SelfdestructFilterer{contract: contract}}, nil
}

// Selfdestruct is an auto generated Go binding around an Ethereum contract.
type Selfdestruct struct {
	SelfdestructCaller     // Read-only binding to the contract
	SelfdestructTransactor // Write-only binding to the contract
	SelfdestructFilterer   // Log filterer for contract events
}

// SelfdestructCaller is an auto generated read-only Go binding around an Ethereum contract.
type SelfdestructCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfdestructTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SelfdestructTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfdestructFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SelfdestructFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfdestructSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SelfdestructSession struct {
	Contract     *Selfdestruct     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SelfdestructCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SelfdestructCallerSession struct {
	Contract *SelfdestructCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// SelfdestructTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SelfdestructTransactorSession struct {
	Contract     *SelfdestructTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// SelfdestructRaw is an auto generated low-level Go binding around an Ethereum contract.
type SelfdestructRaw struct {
	Contract *Selfdestruct // Generic contract binding to access the raw methods on
}

// SelfdestructCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SelfdestructCallerRaw struct {
	Contract *SelfdestructCaller // Generic read-only contract binding to access the raw methods on
}

// SelfdestructTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SelfdestructTransactorRaw struct {
	Contract *SelfdestructTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSelfdestruct creates a new instance of Selfdestruct, bound to a specific deployed contract.
func NewSelfdestruct(address common.Address, backend bind.ContractBackend) (*Selfdestruct, error) {
	contract, err := bindSelfdestruct(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Selfdestruct{SelfdestructCaller: SelfdestructCaller{contract: contract}, SelfdestructTransactor: SelfdestructTransactor{contract: contract}, SelfdestructFilterer: SelfdestructFilterer{contract: contract}}, nil
}

// NewSelfdestructCaller creates a new read-only instance of Selfdestruct, bound to a specific deployed contract.
func NewSelfdestructCaller(address common.Address, caller bind.ContractCaller) (*SelfdestructCaller, error) {
	contract, err := bindSelfdestruct(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SelfdestructCaller{contract: contract}, nil
}

// NewSelfdestructTransactor creates a new write-only instance of Selfdestruct, bound to a specific deployed contract.
func NewSelfdestructTransactor(address common.Address, transactor bind.ContractTransactor) (*SelfdestructTransactor, error) {
	contract, err := bindSelfdestruct(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SelfdestructTransactor{contract: contract}, nil
}

// NewSelfdestructFilterer creates a new log filterer instance of Selfdestruct, bound to a specific deployed contract.
func NewSelfdestructFilterer(address common.Address, filterer bind.ContractFilterer) (*SelfdestructFilterer, error) {
	contract, err := bindSelfdestruct(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SelfdestructFilterer{contract: contract}, nil
}

// bindSelfdestruct binds a generic wrapper to an already deployed contract.
func bindSelfdestruct(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SelfdestructMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Selfdestruct *SelfdestructRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Selfdestruct.Contract.SelfdestructCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Selfdestruct *SelfdestructRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Selfdestruct.Contract.SelfdestructTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Selfdestruct *SelfdestructRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Selfdestruct.Contract.SelfdestructTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Selfdestruct *SelfdestructCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Selfdestruct.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Selfdestruct *SelfdestructTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Selfdestruct.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Selfdestruct *SelfdestructTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Selfdestruct.Contract.contract.Transact(opts, method, params...)
}

// SomeData is a free data retrieval call binding the contract method 0xc73b7175.
//
// Solidity: function someData() view returns(uint256)
func (_Selfdestruct *SelfdestructCaller) SomeData(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Selfdestruct.contract.Call(opts, &out, "someData")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SomeData is a free data retrieval call binding the contract method 0xc73b7175.
//
// Solidity: function someData() view returns(uint256)
func (_Selfdestruct *SelfdestructSession) SomeData() (*big.Int, error) {
	return _Selfdestruct.Contract.SomeData(&_Selfdestruct.CallOpts)
}

// SomeData is a free data retrieval call binding the contract method 0xc73b7175.
//
// Solidity: function someData() view returns(uint256)
func (_Selfdestruct *SelfdestructCallerSession) SomeData() (*big.Int, error) {
	return _Selfdestruct.Contract.SomeData(&_Selfdestruct.CallOpts)
}

// DestroyContract is a paid mutator transaction binding the contract method 0xca29aee9.
//
// Solidity: function destroyContract(bool selfBeneficiary, address recipient) returns()
func (_Selfdestruct *SelfdestructTransactor) DestroyContract(opts *bind.TransactOpts, selfBeneficiary bool, recipient common.Address) (*types.Transaction, error) {
	return _Selfdestruct.contract.Transact(opts, "destroyContract", selfBeneficiary, recipient)
}

// DestroyContract is a paid mutator transaction binding the contract method 0xca29aee9.
//
// Solidity: function destroyContract(bool selfBeneficiary, address recipient) returns()
func (_Selfdestruct *SelfdestructSession) DestroyContract(selfBeneficiary bool, recipient common.Address) (*types.Transaction, error) {
	return _Selfdestruct.Contract.DestroyContract(&_Selfdestruct.TransactOpts, selfBeneficiary, recipient)
}

// DestroyContract is a paid mutator transaction binding the contract method 0xca29aee9.
//
// Solidity: function destroyContract(bool selfBeneficiary, address recipient) returns()
func (_Selfdestruct *SelfdestructTransactorSession) DestroyContract(selfBeneficiary bool, recipient common.Address) (*types.Transaction, error) {
	return _Selfdestruct.Contract.DestroyContract(&_Selfdestruct.TransactOpts, selfBeneficiary, recipient)
}

// SelfdestructLogAfterDestructIterator is returned from FilterLogAfterDestruct and is used to iterate over the raw logs and unpacked data for LogAfterDestruct events raised by the Selfdestruct contract.
type SelfdestructLogAfterDestructIterator struct {
	Event *SelfdestructLogAfterDestruct // Event containing the contract specifics and raw log

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
func (it *SelfdestructLogAfterDestructIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelfdestructLogAfterDestruct)
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
		it.Event = new(SelfdestructLogAfterDestruct)
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
func (it *SelfdestructLogAfterDestructIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SelfdestructLogAfterDestructIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SelfdestructLogAfterDestruct represents a LogAfterDestruct event raised by the Selfdestruct contract.
type SelfdestructLogAfterDestruct struct {
	Balance *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogAfterDestruct is a free log retrieval operation binding the contract event 0xf91cb366864fbc31224338e7a7493dfb408ce511551c30da5895d780d726a9b5.
//
// Solidity: event LogAfterDestruct(uint256 balance)
func (_Selfdestruct *SelfdestructFilterer) FilterLogAfterDestruct(opts *bind.FilterOpts) (*SelfdestructLogAfterDestructIterator, error) {

	logs, sub, err := _Selfdestruct.contract.FilterLogs(opts, "LogAfterDestruct")
	if err != nil {
		return nil, err
	}
	return &SelfdestructLogAfterDestructIterator{contract: _Selfdestruct.contract, event: "LogAfterDestruct", logs: logs, sub: sub}, nil
}

// WatchLogAfterDestruct is a free log subscription operation binding the contract event 0xf91cb366864fbc31224338e7a7493dfb408ce511551c30da5895d780d726a9b5.
//
// Solidity: event LogAfterDestruct(uint256 balance)
func (_Selfdestruct *SelfdestructFilterer) WatchLogAfterDestruct(opts *bind.WatchOpts, sink chan<- *SelfdestructLogAfterDestruct) (event.Subscription, error) {

	logs, sub, err := _Selfdestruct.contract.WatchLogs(opts, "LogAfterDestruct")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SelfdestructLogAfterDestruct)
				if err := _Selfdestruct.contract.UnpackLog(event, "LogAfterDestruct", log); err != nil {
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

// ParseLogAfterDestruct is a log parse operation binding the contract event 0xf91cb366864fbc31224338e7a7493dfb408ce511551c30da5895d780d726a9b5.
//
// Solidity: event LogAfterDestruct(uint256 balance)
func (_Selfdestruct *SelfdestructFilterer) ParseLogAfterDestruct(log types.Log) (*SelfdestructLogAfterDestruct, error) {
	event := new(SelfdestructLogAfterDestruct)
	if err := _Selfdestruct.contract.UnpackLog(event, "LogAfterDestruct", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
