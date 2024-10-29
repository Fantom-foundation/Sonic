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
	ABI: "[{\"inputs\":[],\"name\":\"create2ContractWithInvalidCode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"create2ContractWithValidCode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcd\",\"type\":\"bytes\"}],\"name\":\"create2OrRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createContractWithInvalidCode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createContractWithValidCode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"bytecode\",\"type\":\"bytes\"}],\"name\":\"createOrRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"invalidBytecode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validBytecode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040526040518060400160405280600a81526020017f60ef60005360016000f3000000000000000000000000000000000000000000008152505f908161004791906102db565b506040518060400160405280600a81526020017f60fe60005360016000f3000000000000000000000000000000000000000000008152506001908161008c91906102db565b50348015610098575f5ffd5b506103aa565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061011957607f821691505b60208210810361012c5761012b6100d5565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261018e7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610153565b6101988683610153565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6101dc6101d76101d2846101b0565b6101b9565b6101b0565b9050919050565b5f819050919050565b6101f5836101c2565b610209610201826101e3565b84845461015f565b825550505050565b5f5f905090565b610220610211565b61022b8184846101ec565b505050565b5b8181101561024e576102435f82610218565b600181019050610231565b5050565b601f8211156102935761026481610132565b61026d84610144565b8101602085101561027c578190505b61029061028885610144565b830182610230565b50505b505050565b5f82821c905092915050565b5f6102b35f1984600802610298565b1980831691505092915050565b5f6102cb83836102a4565b9150826002028217905092915050565b6102e48261009e565b67ffffffffffffffff8111156102fd576102fc6100a8565b5b6103078254610102565b610312828285610252565b5f60209050601f831160018114610343575f8415610331578287015190505b61033b85826102c0565b8655506103a2565b601f19841661035186610132565b5f5b8281101561037857848901518255600182019150602085019450602081019050610353565b868310156103955784890151610391601f8916826102a4565b8355505b6001600288020188555050505b505050505050565b610759806103b75f395ff3fe608060405234801561000f575f5ffd5b5060043610610086575f3560e01c8063a2990d4211610059578063a2990d42146100da578063d5a28fdc146100e4578063f3245b9a14610100578063fd3bb91e1461010a57610086565b806302fc6a881461008a5780630da11c4b146100945780634652ff07146100b257806389c0ea30146100d0575b5f5ffd5b610092610126565b005b61009c6101b9565b6040516100a99190610522565b60405180910390f35b6100ba610244565b6040516100c79190610522565b60405180910390f35b6100d86102d0565b005b6100e2610363565b005b6100fe60048036038101906100f9919061067f565b6103f5565b005b61010861040a565b005b610124600480360381019061011f919061067f565b61049c565b005b6101b760018054610136906106f3565b80601f0160208091040260200160405190810160405280929190818152602001828054610162906106f3565b80156101ad5780601f10610184576101008083540402835291602001916101ad565b820191905f5260205f20905b81548152906001019060200180831161019057829003601f168201915b50505050506103f5565b565b5f80546101c5906106f3565b80601f01602080910402602001604051908101604052809291908181526020018280546101f1906106f3565b801561023c5780601f106102135761010080835404028352916020019161023c565b820191905f5260205f20905b81548152906001019060200180831161021f57829003601f168201915b505050505081565b60018054610251906106f3565b80601f016020809104026020016040519081016040528092919081815260200182805461027d906106f3565b80156102c85780601f1061029f576101008083540402835291602001916102c8565b820191905f5260205f20905b8154815290600101906020018083116102ab57829003601f168201915b505050505081565b610361600180546102e0906106f3565b80601f016020809104026020016040519081016040528092919081815260200182805461030c906106f3565b80156103575780601f1061032e57610100808354040283529160200191610357565b820191905f5260205f20905b81548152906001019060200180831161033a57829003601f168201915b505050505061049c565b565b6103f35f8054610372906106f3565b80601f016020809104026020016040519081016040528092919081815260200182805461039e906106f3565b80156103e95780601f106103c0576101008083540402835291602001916103e9565b820191905f5260205f20905b8154815290600101906020018083116103cc57829003601f168201915b50505050506103f5565b565b8051602082015ff080610406575f5ffd5b5050565b61049a5f8054610419906106f3565b80601f0160208091040260200160405190810160405280929190818152602001828054610445906106f3565b80156104905780601f1061046757610100808354040283529160200191610490565b820191905f5260205f20905b81548152906001019060200180831161047357829003601f168201915b505050505061049c565b565b5f8151602083015ff5806104ae575f5ffd5b5050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6104f4826104b2565b6104fe81856104bc565b935061050e8185602086016104cc565b610517816104da565b840191505092915050565b5f6020820190508181035f83015261053a81846104ea565b905092915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610591826104da565b810181811067ffffffffffffffff821117156105b0576105af61055b565b5b80604052505050565b5f6105c2610542565b90506105ce8282610588565b919050565b5f67ffffffffffffffff8211156105ed576105ec61055b565b5b6105f6826104da565b9050602081019050919050565b828183375f83830152505050565b5f61062361061e846105d3565b6105b9565b90508281526020810184848401111561063f5761063e610557565b5b61064a848285610603565b509392505050565b5f82601f83011261066657610665610553565b5b8135610676848260208601610611565b91505092915050565b5f602082840312156106945761069361054b565b5b5f82013567ffffffffffffffff8111156106b1576106b061054f565b5b6106bd84828501610652565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061070a57607f821691505b60208210810361071d5761071c6106c6565b5b5091905056fea2646970667358221220e3274787395684101b7a3076c71a017061e31de71ccf72e12b06eb14634bd27b64736f6c634300081c0033",
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

// InvalidBytecode is a free data retrieval call binding the contract method 0x0da11c4b.
//
// Solidity: function invalidBytecode() view returns(bytes)
func (_Invalidstart *InvalidstartCaller) InvalidBytecode(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Invalidstart.contract.Call(opts, &out, "invalidBytecode")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// InvalidBytecode is a free data retrieval call binding the contract method 0x0da11c4b.
//
// Solidity: function invalidBytecode() view returns(bytes)
func (_Invalidstart *InvalidstartSession) InvalidBytecode() ([]byte, error) {
	return _Invalidstart.Contract.InvalidBytecode(&_Invalidstart.CallOpts)
}

// InvalidBytecode is a free data retrieval call binding the contract method 0x0da11c4b.
//
// Solidity: function invalidBytecode() view returns(bytes)
func (_Invalidstart *InvalidstartCallerSession) InvalidBytecode() ([]byte, error) {
	return _Invalidstart.Contract.InvalidBytecode(&_Invalidstart.CallOpts)
}

// ValidBytecode is a free data retrieval call binding the contract method 0x4652ff07.
//
// Solidity: function validBytecode() view returns(bytes)
func (_Invalidstart *InvalidstartCaller) ValidBytecode(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Invalidstart.contract.Call(opts, &out, "validBytecode")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ValidBytecode is a free data retrieval call binding the contract method 0x4652ff07.
//
// Solidity: function validBytecode() view returns(bytes)
func (_Invalidstart *InvalidstartSession) ValidBytecode() ([]byte, error) {
	return _Invalidstart.Contract.ValidBytecode(&_Invalidstart.CallOpts)
}

// ValidBytecode is a free data retrieval call binding the contract method 0x4652ff07.
//
// Solidity: function validBytecode() view returns(bytes)
func (_Invalidstart *InvalidstartCallerSession) ValidBytecode() ([]byte, error) {
	return _Invalidstart.Contract.ValidBytecode(&_Invalidstart.CallOpts)
}

// Create2ContractWithInvalidCode is a paid mutator transaction binding the contract method 0xf3245b9a.
//
// Solidity: function create2ContractWithInvalidCode() returns()
func (_Invalidstart *InvalidstartTransactor) Create2ContractWithInvalidCode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invalidstart.contract.Transact(opts, "create2ContractWithInvalidCode")
}

// Create2ContractWithInvalidCode is a paid mutator transaction binding the contract method 0xf3245b9a.
//
// Solidity: function create2ContractWithInvalidCode() returns()
func (_Invalidstart *InvalidstartSession) Create2ContractWithInvalidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.Create2ContractWithInvalidCode(&_Invalidstart.TransactOpts)
}

// Create2ContractWithInvalidCode is a paid mutator transaction binding the contract method 0xf3245b9a.
//
// Solidity: function create2ContractWithInvalidCode() returns()
func (_Invalidstart *InvalidstartTransactorSession) Create2ContractWithInvalidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.Create2ContractWithInvalidCode(&_Invalidstart.TransactOpts)
}

// Create2ContractWithValidCode is a paid mutator transaction binding the contract method 0x89c0ea30.
//
// Solidity: function create2ContractWithValidCode() returns()
func (_Invalidstart *InvalidstartTransactor) Create2ContractWithValidCode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invalidstart.contract.Transact(opts, "create2ContractWithValidCode")
}

// Create2ContractWithValidCode is a paid mutator transaction binding the contract method 0x89c0ea30.
//
// Solidity: function create2ContractWithValidCode() returns()
func (_Invalidstart *InvalidstartSession) Create2ContractWithValidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.Create2ContractWithValidCode(&_Invalidstart.TransactOpts)
}

// Create2ContractWithValidCode is a paid mutator transaction binding the contract method 0x89c0ea30.
//
// Solidity: function create2ContractWithValidCode() returns()
func (_Invalidstart *InvalidstartTransactorSession) Create2ContractWithValidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.Create2ContractWithValidCode(&_Invalidstart.TransactOpts)
}

// Create2OrRevert is a paid mutator transaction binding the contract method 0xfd3bb91e.
//
// Solidity: function create2OrRevert(bytes btcd) returns()
func (_Invalidstart *InvalidstartTransactor) Create2OrRevert(opts *bind.TransactOpts, btcd []byte) (*types.Transaction, error) {
	return _Invalidstart.contract.Transact(opts, "create2OrRevert", btcd)
}

// Create2OrRevert is a paid mutator transaction binding the contract method 0xfd3bb91e.
//
// Solidity: function create2OrRevert(bytes btcd) returns()
func (_Invalidstart *InvalidstartSession) Create2OrRevert(btcd []byte) (*types.Transaction, error) {
	return _Invalidstart.Contract.Create2OrRevert(&_Invalidstart.TransactOpts, btcd)
}

// Create2OrRevert is a paid mutator transaction binding the contract method 0xfd3bb91e.
//
// Solidity: function create2OrRevert(bytes btcd) returns()
func (_Invalidstart *InvalidstartTransactorSession) Create2OrRevert(btcd []byte) (*types.Transaction, error) {
	return _Invalidstart.Contract.Create2OrRevert(&_Invalidstart.TransactOpts, btcd)
}

// CreateContractWithInvalidCode is a paid mutator transaction binding the contract method 0xa2990d42.
//
// Solidity: function createContractWithInvalidCode() returns()
func (_Invalidstart *InvalidstartTransactor) CreateContractWithInvalidCode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invalidstart.contract.Transact(opts, "createContractWithInvalidCode")
}

// CreateContractWithInvalidCode is a paid mutator transaction binding the contract method 0xa2990d42.
//
// Solidity: function createContractWithInvalidCode() returns()
func (_Invalidstart *InvalidstartSession) CreateContractWithInvalidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateContractWithInvalidCode(&_Invalidstart.TransactOpts)
}

// CreateContractWithInvalidCode is a paid mutator transaction binding the contract method 0xa2990d42.
//
// Solidity: function createContractWithInvalidCode() returns()
func (_Invalidstart *InvalidstartTransactorSession) CreateContractWithInvalidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateContractWithInvalidCode(&_Invalidstart.TransactOpts)
}

// CreateContractWithValidCode is a paid mutator transaction binding the contract method 0x02fc6a88.
//
// Solidity: function createContractWithValidCode() returns()
func (_Invalidstart *InvalidstartTransactor) CreateContractWithValidCode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invalidstart.contract.Transact(opts, "createContractWithValidCode")
}

// CreateContractWithValidCode is a paid mutator transaction binding the contract method 0x02fc6a88.
//
// Solidity: function createContractWithValidCode() returns()
func (_Invalidstart *InvalidstartSession) CreateContractWithValidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateContractWithValidCode(&_Invalidstart.TransactOpts)
}

// CreateContractWithValidCode is a paid mutator transaction binding the contract method 0x02fc6a88.
//
// Solidity: function createContractWithValidCode() returns()
func (_Invalidstart *InvalidstartTransactorSession) CreateContractWithValidCode() (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateContractWithValidCode(&_Invalidstart.TransactOpts)
}

// CreateOrRevert is a paid mutator transaction binding the contract method 0xd5a28fdc.
//
// Solidity: function createOrRevert(bytes bytecode) returns()
func (_Invalidstart *InvalidstartTransactor) CreateOrRevert(opts *bind.TransactOpts, bytecode []byte) (*types.Transaction, error) {
	return _Invalidstart.contract.Transact(opts, "createOrRevert", bytecode)
}

// CreateOrRevert is a paid mutator transaction binding the contract method 0xd5a28fdc.
//
// Solidity: function createOrRevert(bytes bytecode) returns()
func (_Invalidstart *InvalidstartSession) CreateOrRevert(bytecode []byte) (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateOrRevert(&_Invalidstart.TransactOpts, bytecode)
}

// CreateOrRevert is a paid mutator transaction binding the contract method 0xd5a28fdc.
//
// Solidity: function createOrRevert(bytes bytecode) returns()
func (_Invalidstart *InvalidstartTransactorSession) CreateOrRevert(bytecode []byte) (*types.Transaction, error) {
	return _Invalidstart.Contract.CreateOrRevert(&_Invalidstart.TransactOpts, bytecode)
}
