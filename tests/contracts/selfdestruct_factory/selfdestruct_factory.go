// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package selfdestruct_factory

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

// SelfdestructFactoryMetaData contains all meta data concerning the SelfdestructFactory contract.
var SelfdestructFactoryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"LogContractStorage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogDeployed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"create\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"beneficiary\",\"type\":\"address\"}],\"name\":\"createAndDestroy\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createAndDestroyWithoutBeneficiary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"beneficiary\",\"type\":\"address\"}],\"name\":\"destroy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"destroyWithoutBeneficiary\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506109a98061001c5f395ff3fe608060405260043610610048575f3560e01c8062f55d9d1461004c5780635244e834146100745780637d3cf16a14610090578063aaff233b146100a6578063efc81a8c146100b0575b5f5ffd5b348015610057575f5ffd5b50610072600480360381019061006d91906104d7565b6100ba565b005b61008e600480360381019061008991906104d7565b610145565b005b34801561009b575f5ffd5b506100a461021c565b005b6100ae6102a7565b005b6100b861039d565b005b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ca29aee95f836040518363ffffffff1660e01b815260040161011592919061052b565b5f604051808303815f87803b15801561012c575f5ffd5b505af115801561013e573d5f5f3e3d5ffd5b5050505050565b61014d61039d565b610156816100ba565b7f3be3676a60e58335e9af2a05f4c6208663920cd9e2cd3900745377049fa721f05f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c73b71756040518163ffffffff1660e01b8152600401602060405180830381865afa1580156101e0573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906102049190610585565b60405161021191906105bf565b60405180910390a150565b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ca29aee960015f6040518363ffffffff1660e01b815260040161027892919061052b565b5f604051808303815f87803b15801561028f575f5ffd5b505af11580156102a1573d5f5f3e3d5ffd5b50505050565b6102af61039d565b6102d85f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff166100ba565b7f3be3676a60e58335e9af2a05f4c6208663920cd9e2cd3900745377049fa721f05f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c73b71756040518163ffffffff1660e01b8152600401602060405180830381865afa158015610362573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103869190610585565b60405161039391906105bf565b60405180910390a1565b345f5f5f6040516103ad9061046c565b6103b9939291906105d8565b6040518091039082f09050801580156103d4573d5f5f3e3d5ffd5b505f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f4a0f3be106bcb1b76a2cf127e86da8ceb3dba6b2d1bb19ec8ffbfa90c4647b515f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16604051610462919061062d565b60405180910390a1565b61032d8061064783390190565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6104a68261047d565b9050919050565b6104b68161049c565b81146104c0575f5ffd5b50565b5f813590506104d1816104ad565b92915050565b5f602082840312156104ec576104eb610479565b5b5f6104f9848285016104c3565b91505092915050565b5f8115159050919050565b61051681610502565b82525050565b6105258161049c565b82525050565b5f60408201905061053e5f83018561050d565b61054b602083018461051c565b9392505050565b5f819050919050565b61056481610552565b811461056e575f5ffd5b50565b5f8151905061057f8161055b565b92915050565b5f6020828403121561059a57610599610479565b5b5f6105a784828501610571565b91505092915050565b6105b981610552565b82525050565b5f6020820190506105d25f8301846105b0565b92915050565b5f6060820190506105eb5f83018661050d565b6105f8602083018561050d565b610605604083018461051c565b949350505050565b5f6106178261047d565b9050919050565b6106278161060d565b82525050565b5f6020820190506106405f83018461061e565b9291505056fe6080604052607b5f5560405161032d38038061032d833981810160405281019061002991906100fe565b82156100405761003f828261004860201b60201c565b5b50505061014e565b8115610052573090505b8073ffffffffffffffffffffffffffffffffffffffff16ff5b5f5ffd5b5f8115159050919050565b6100838161006f565b811461008d575f5ffd5b50565b5f8151905061009e8161007a565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100cd826100a4565b9050919050565b6100dd816100c3565b81146100e7575f5ffd5b50565b5f815190506100f8816100d4565b92915050565b5f5f5f606084860312156101155761011461006b565b5b5f61012286828701610090565b935050602061013386828701610090565b9250506040610144868287016100ea565b9150509250925092565b6101d28061015b5f395ff3fe608060405234801561000f575f5ffd5b5060043610610034575f3560e01c8063c73b717514610038578063ca29aee914610056575b5f5ffd5b610040610072565b60405161004d91906100b2565b60405180910390f35b610070600480360381019061006b919061015e565b610077565b005b5f5481565b8115610081573090505b8073ffffffffffffffffffffffffffffffffffffffff16ff5b5f819050919050565b6100ac8161009a565b82525050565b5f6020820190506100c55f8301846100a3565b92915050565b5f5ffd5b5f8115159050919050565b6100e3816100cf565b81146100ed575f5ffd5b50565b5f813590506100fe816100da565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61012d82610104565b9050919050565b61013d81610123565b8114610147575f5ffd5b50565b5f8135905061015881610134565b92915050565b5f5f60408385031215610174576101736100cb565b5b5f610181858286016100f0565b92505060206101928582860161014a565b915050925092905056fea264697066735822122032c75c0702475f02b57c96638d6fa9b294ddba493eda1ea9cc2486cb5535a3f364736f6c634300081c0033a2646970667358221220efdf6183262bfcafda1fb2db4d7cbb0b68ead12a0bf728eb30a512344723aeb364736f6c634300081c0033",
}

// SelfdestructFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use SelfdestructFactoryMetaData.ABI instead.
var SelfdestructFactoryABI = SelfdestructFactoryMetaData.ABI

// SelfdestructFactoryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SelfdestructFactoryMetaData.Bin instead.
var SelfdestructFactoryBin = SelfdestructFactoryMetaData.Bin

// DeploySelfdestructFactory deploys a new Ethereum contract, binding an instance of SelfdestructFactory to it.
func DeploySelfdestructFactory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SelfdestructFactory, error) {
	parsed, err := SelfdestructFactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SelfdestructFactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SelfdestructFactory{SelfdestructFactoryCaller: SelfdestructFactoryCaller{contract: contract}, SelfdestructFactoryTransactor: SelfdestructFactoryTransactor{contract: contract}, SelfdestructFactoryFilterer: SelfdestructFactoryFilterer{contract: contract}}, nil
}

// SelfdestructFactory is an auto generated Go binding around an Ethereum contract.
type SelfdestructFactory struct {
	SelfdestructFactoryCaller     // Read-only binding to the contract
	SelfdestructFactoryTransactor // Write-only binding to the contract
	SelfdestructFactoryFilterer   // Log filterer for contract events
}

// SelfdestructFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SelfdestructFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfdestructFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SelfdestructFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfdestructFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SelfdestructFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfdestructFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SelfdestructFactorySession struct {
	Contract     *SelfdestructFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// SelfdestructFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SelfdestructFactoryCallerSession struct {
	Contract *SelfdestructFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// SelfdestructFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SelfdestructFactoryTransactorSession struct {
	Contract     *SelfdestructFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// SelfdestructFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SelfdestructFactoryRaw struct {
	Contract *SelfdestructFactory // Generic contract binding to access the raw methods on
}

// SelfdestructFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SelfdestructFactoryCallerRaw struct {
	Contract *SelfdestructFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// SelfdestructFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SelfdestructFactoryTransactorRaw struct {
	Contract *SelfdestructFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSelfdestructFactory creates a new instance of SelfdestructFactory, bound to a specific deployed contract.
func NewSelfdestructFactory(address common.Address, backend bind.ContractBackend) (*SelfdestructFactory, error) {
	contract, err := bindSelfdestructFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SelfdestructFactory{SelfdestructFactoryCaller: SelfdestructFactoryCaller{contract: contract}, SelfdestructFactoryTransactor: SelfdestructFactoryTransactor{contract: contract}, SelfdestructFactoryFilterer: SelfdestructFactoryFilterer{contract: contract}}, nil
}

// NewSelfdestructFactoryCaller creates a new read-only instance of SelfdestructFactory, bound to a specific deployed contract.
func NewSelfdestructFactoryCaller(address common.Address, caller bind.ContractCaller) (*SelfdestructFactoryCaller, error) {
	contract, err := bindSelfdestructFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SelfdestructFactoryCaller{contract: contract}, nil
}

// NewSelfdestructFactoryTransactor creates a new write-only instance of SelfdestructFactory, bound to a specific deployed contract.
func NewSelfdestructFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*SelfdestructFactoryTransactor, error) {
	contract, err := bindSelfdestructFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SelfdestructFactoryTransactor{contract: contract}, nil
}

// NewSelfdestructFactoryFilterer creates a new log filterer instance of SelfdestructFactory, bound to a specific deployed contract.
func NewSelfdestructFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*SelfdestructFactoryFilterer, error) {
	contract, err := bindSelfdestructFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SelfdestructFactoryFilterer{contract: contract}, nil
}

// bindSelfdestructFactory binds a generic wrapper to an already deployed contract.
func bindSelfdestructFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SelfdestructFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SelfdestructFactory *SelfdestructFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SelfdestructFactory.Contract.SelfdestructFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SelfdestructFactory *SelfdestructFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.SelfdestructFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SelfdestructFactory *SelfdestructFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.SelfdestructFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SelfdestructFactory *SelfdestructFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SelfdestructFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SelfdestructFactory *SelfdestructFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SelfdestructFactory *SelfdestructFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.contract.Transact(opts, method, params...)
}

// Create is a paid mutator transaction binding the contract method 0xefc81a8c.
//
// Solidity: function create() payable returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactor) Create(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfdestructFactory.contract.Transact(opts, "create")
}

// Create is a paid mutator transaction binding the contract method 0xefc81a8c.
//
// Solidity: function create() payable returns()
func (_SelfdestructFactory *SelfdestructFactorySession) Create() (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.Create(&_SelfdestructFactory.TransactOpts)
}

// Create is a paid mutator transaction binding the contract method 0xefc81a8c.
//
// Solidity: function create() payable returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactorSession) Create() (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.Create(&_SelfdestructFactory.TransactOpts)
}

// CreateAndDestroy is a paid mutator transaction binding the contract method 0x5244e834.
//
// Solidity: function createAndDestroy(address beneficiary) payable returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactor) CreateAndDestroy(opts *bind.TransactOpts, beneficiary common.Address) (*types.Transaction, error) {
	return _SelfdestructFactory.contract.Transact(opts, "createAndDestroy", beneficiary)
}

// CreateAndDestroy is a paid mutator transaction binding the contract method 0x5244e834.
//
// Solidity: function createAndDestroy(address beneficiary) payable returns()
func (_SelfdestructFactory *SelfdestructFactorySession) CreateAndDestroy(beneficiary common.Address) (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.CreateAndDestroy(&_SelfdestructFactory.TransactOpts, beneficiary)
}

// CreateAndDestroy is a paid mutator transaction binding the contract method 0x5244e834.
//
// Solidity: function createAndDestroy(address beneficiary) payable returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactorSession) CreateAndDestroy(beneficiary common.Address) (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.CreateAndDestroy(&_SelfdestructFactory.TransactOpts, beneficiary)
}

// CreateAndDestroyWithoutBeneficiary is a paid mutator transaction binding the contract method 0xaaff233b.
//
// Solidity: function createAndDestroyWithoutBeneficiary() payable returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactor) CreateAndDestroyWithoutBeneficiary(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfdestructFactory.contract.Transact(opts, "createAndDestroyWithoutBeneficiary")
}

// CreateAndDestroyWithoutBeneficiary is a paid mutator transaction binding the contract method 0xaaff233b.
//
// Solidity: function createAndDestroyWithoutBeneficiary() payable returns()
func (_SelfdestructFactory *SelfdestructFactorySession) CreateAndDestroyWithoutBeneficiary() (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.CreateAndDestroyWithoutBeneficiary(&_SelfdestructFactory.TransactOpts)
}

// CreateAndDestroyWithoutBeneficiary is a paid mutator transaction binding the contract method 0xaaff233b.
//
// Solidity: function createAndDestroyWithoutBeneficiary() payable returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactorSession) CreateAndDestroyWithoutBeneficiary() (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.CreateAndDestroyWithoutBeneficiary(&_SelfdestructFactory.TransactOpts)
}

// Destroy is a paid mutator transaction binding the contract method 0x00f55d9d.
//
// Solidity: function destroy(address beneficiary) returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactor) Destroy(opts *bind.TransactOpts, beneficiary common.Address) (*types.Transaction, error) {
	return _SelfdestructFactory.contract.Transact(opts, "destroy", beneficiary)
}

// Destroy is a paid mutator transaction binding the contract method 0x00f55d9d.
//
// Solidity: function destroy(address beneficiary) returns()
func (_SelfdestructFactory *SelfdestructFactorySession) Destroy(beneficiary common.Address) (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.Destroy(&_SelfdestructFactory.TransactOpts, beneficiary)
}

// Destroy is a paid mutator transaction binding the contract method 0x00f55d9d.
//
// Solidity: function destroy(address beneficiary) returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactorSession) Destroy(beneficiary common.Address) (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.Destroy(&_SelfdestructFactory.TransactOpts, beneficiary)
}

// DestroyWithoutBeneficiary is a paid mutator transaction binding the contract method 0x7d3cf16a.
//
// Solidity: function destroyWithoutBeneficiary() returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactor) DestroyWithoutBeneficiary(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfdestructFactory.contract.Transact(opts, "destroyWithoutBeneficiary")
}

// DestroyWithoutBeneficiary is a paid mutator transaction binding the contract method 0x7d3cf16a.
//
// Solidity: function destroyWithoutBeneficiary() returns()
func (_SelfdestructFactory *SelfdestructFactorySession) DestroyWithoutBeneficiary() (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.DestroyWithoutBeneficiary(&_SelfdestructFactory.TransactOpts)
}

// DestroyWithoutBeneficiary is a paid mutator transaction binding the contract method 0x7d3cf16a.
//
// Solidity: function destroyWithoutBeneficiary() returns()
func (_SelfdestructFactory *SelfdestructFactoryTransactorSession) DestroyWithoutBeneficiary() (*types.Transaction, error) {
	return _SelfdestructFactory.Contract.DestroyWithoutBeneficiary(&_SelfdestructFactory.TransactOpts)
}

// SelfdestructFactoryLogContractStorageIterator is returned from FilterLogContractStorage and is used to iterate over the raw logs and unpacked data for LogContractStorage events raised by the SelfdestructFactory contract.
type SelfdestructFactoryLogContractStorageIterator struct {
	Event *SelfdestructFactoryLogContractStorage // Event containing the contract specifics and raw log

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
func (it *SelfdestructFactoryLogContractStorageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelfdestructFactoryLogContractStorage)
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
		it.Event = new(SelfdestructFactoryLogContractStorage)
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
func (it *SelfdestructFactoryLogContractStorageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SelfdestructFactoryLogContractStorageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SelfdestructFactoryLogContractStorage represents a LogContractStorage event raised by the SelfdestructFactory contract.
type SelfdestructFactoryLogContractStorage struct {
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLogContractStorage is a free log retrieval operation binding the contract event 0x3be3676a60e58335e9af2a05f4c6208663920cd9e2cd3900745377049fa721f0.
//
// Solidity: event LogContractStorage(uint256 value)
func (_SelfdestructFactory *SelfdestructFactoryFilterer) FilterLogContractStorage(opts *bind.FilterOpts) (*SelfdestructFactoryLogContractStorageIterator, error) {

	logs, sub, err := _SelfdestructFactory.contract.FilterLogs(opts, "LogContractStorage")
	if err != nil {
		return nil, err
	}
	return &SelfdestructFactoryLogContractStorageIterator{contract: _SelfdestructFactory.contract, event: "LogContractStorage", logs: logs, sub: sub}, nil
}

// WatchLogContractStorage is a free log subscription operation binding the contract event 0x3be3676a60e58335e9af2a05f4c6208663920cd9e2cd3900745377049fa721f0.
//
// Solidity: event LogContractStorage(uint256 value)
func (_SelfdestructFactory *SelfdestructFactoryFilterer) WatchLogContractStorage(opts *bind.WatchOpts, sink chan<- *SelfdestructFactoryLogContractStorage) (event.Subscription, error) {

	logs, sub, err := _SelfdestructFactory.contract.WatchLogs(opts, "LogContractStorage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SelfdestructFactoryLogContractStorage)
				if err := _SelfdestructFactory.contract.UnpackLog(event, "LogContractStorage", log); err != nil {
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

// ParseLogContractStorage is a log parse operation binding the contract event 0x3be3676a60e58335e9af2a05f4c6208663920cd9e2cd3900745377049fa721f0.
//
// Solidity: event LogContractStorage(uint256 value)
func (_SelfdestructFactory *SelfdestructFactoryFilterer) ParseLogContractStorage(log types.Log) (*SelfdestructFactoryLogContractStorage, error) {
	event := new(SelfdestructFactoryLogContractStorage)
	if err := _SelfdestructFactory.contract.UnpackLog(event, "LogContractStorage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SelfdestructFactoryLogDeployedIterator is returned from FilterLogDeployed and is used to iterate over the raw logs and unpacked data for LogDeployed events raised by the SelfdestructFactory contract.
type SelfdestructFactoryLogDeployedIterator struct {
	Event *SelfdestructFactoryLogDeployed // Event containing the contract specifics and raw log

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
func (it *SelfdestructFactoryLogDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelfdestructFactoryLogDeployed)
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
		it.Event = new(SelfdestructFactoryLogDeployed)
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
func (it *SelfdestructFactoryLogDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SelfdestructFactoryLogDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SelfdestructFactoryLogDeployed represents a LogDeployed event raised by the SelfdestructFactory contract.
type SelfdestructFactoryLogDeployed struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogDeployed is a free log retrieval operation binding the contract event 0x4a0f3be106bcb1b76a2cf127e86da8ceb3dba6b2d1bb19ec8ffbfa90c4647b51.
//
// Solidity: event LogDeployed(address addr)
func (_SelfdestructFactory *SelfdestructFactoryFilterer) FilterLogDeployed(opts *bind.FilterOpts) (*SelfdestructFactoryLogDeployedIterator, error) {

	logs, sub, err := _SelfdestructFactory.contract.FilterLogs(opts, "LogDeployed")
	if err != nil {
		return nil, err
	}
	return &SelfdestructFactoryLogDeployedIterator{contract: _SelfdestructFactory.contract, event: "LogDeployed", logs: logs, sub: sub}, nil
}

// WatchLogDeployed is a free log subscription operation binding the contract event 0x4a0f3be106bcb1b76a2cf127e86da8ceb3dba6b2d1bb19ec8ffbfa90c4647b51.
//
// Solidity: event LogDeployed(address addr)
func (_SelfdestructFactory *SelfdestructFactoryFilterer) WatchLogDeployed(opts *bind.WatchOpts, sink chan<- *SelfdestructFactoryLogDeployed) (event.Subscription, error) {

	logs, sub, err := _SelfdestructFactory.contract.WatchLogs(opts, "LogDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SelfdestructFactoryLogDeployed)
				if err := _SelfdestructFactory.contract.UnpackLog(event, "LogDeployed", log); err != nil {
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

// ParseLogDeployed is a log parse operation binding the contract event 0x4a0f3be106bcb1b76a2cf127e86da8ceb3dba6b2d1bb19ec8ffbfa90c4647b51.
//
// Solidity: event LogDeployed(address addr)
func (_SelfdestructFactory *SelfdestructFactoryFilterer) ParseLogDeployed(log types.Log) (*SelfdestructFactoryLogDeployed, error) {
	event := new(SelfdestructFactoryLogDeployed)
	if err := _SelfdestructFactory.contract.UnpackLog(event, "LogDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
