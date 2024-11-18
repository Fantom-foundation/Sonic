// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package driverauth100

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
)

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_sfc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_driver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newDriverAuth\",\"type\":\"address\"}],\"name\":\"migrateTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"executable\",\"type\":\"address\"}],\"name\":\"execute\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"executable\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"selfCodeHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"driverCodeHash\",\"type\":\"bytes32\"}],\"name\":\"mutExecute\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"diff\",\"type\":\"uint256\"}],\"name\":\"incBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"upgradeCode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"copyCode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"diff\",\"type\":\"uint256\"}],\"name\":\"incNonce\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"diff\",\"type\":\"bytes\"}],\"name\":\"updateNetworkRules\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minGasPrice\",\"type\":\"uint256\"}],\"name\":\"updateMinGasPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"updateNetworkVersion\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"advanceEpochs\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"updateValidatorWeight\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"updateValidatorPubkey\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_auth\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedTime\",\"type\":\"uint256\"}],\"name\":\"setGenesisValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupFromEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupEndTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earlyUnlockPenalty\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewards\",\"type\":\"uint256\"}],\"name\":\"setGenesisDelegation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"name\":\"deactivateValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nextValidatorIDs\",\"type\":\"uint256[]\"}],\"name\":\"sealEpochValidators\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"offlineTimes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"offlineBlocks\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"uptimes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"originatedTxsFee\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"usedGas\",\"type\":\"uint256\"}],\"name\":\"sealEpoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50612269806100206000396000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c806366e7ea0f116100e3578063b9cc6b1c1161008c578063e08d7e6611610066578063e08d7e6614610702578063f2fde38b14610772578063fd1b6ec1146107985761018d565b8063b9cc6b1c1461062c578063c0c53b8b1461069c578063d6a0c7af146106d45761018d565b80638da5cb5b116100bd5780638da5cb5b146105c95780638f32d59b146105ed578063a4066fbe146106095761018d565b806366e7ea0f14610569578063715018a61461059557806379bead381461059d5761018d565b8063242a6e3f116101455780634ddaf8f21161011f5780634ddaf8f21461033f5780634feb92f314610365578063592fe0c0146104035761018d565b8063242a6e3f14610285578063267ab446146102fc5780634b64e492146103195761018d565b806318f628d41161017657806318f628d4146101ce5780631cef4fab146102265780631e702f83146102625761018d565b806307aaf344146101925780630aeeca00146101b1575b600080fd5b6101af600480360360208110156101a857600080fd5b50356107c6565b005b6101af600480360360208110156101c757600080fd5b5035610969565b6101af60048036036101208110156101e557600080fd5b506001600160a01b038135169060208101359060408101359060608101359060808101359060a08101359060c08101359060e0810135906101000135610a28565b6101af6004803603608081101561023c57600080fd5b506001600160a01b03813581169160208101359091169060408101359060600135610b35565b6101af6004803603604081101561027857600080fd5b5080359060200135610ba0565b6101af6004803603604081101561029b57600080fd5b813591908101906040810160208201356401000000008111156102bd57600080fd5b8201836020820111156102cf57600080fd5b803590602001918460018302840111640100000000831117156102f157600080fd5b509092509050610c72565b6101af6004803603602081101561031257600080fd5b5035610d87565b6101af6004803603602081101561032f57600080fd5b50356001600160a01b0316610e46565b6101af6004803603602081101561035557600080fd5b50356001600160a01b0316610ed1565b6101af600480360361010081101561037c57600080fd5b6001600160a01b03823516916020810135918101906060810160408201356401000000008111156103ac57600080fd5b8201836020820111156103be57600080fd5b803590602001918460018302840111640100000000831117156103e057600080fd5b919350915080359060208101359060408101359060608101359060800135610f91565b6101af600480360360a081101561041957600080fd5b81019060208101813564010000000081111561043457600080fd5b82018360208201111561044657600080fd5b8035906020019184602083028401116401000000008311171561046857600080fd5b91939092909160208101903564010000000081111561048657600080fd5b82018360208201111561049857600080fd5b803590602001918460208302840111640100000000831117156104ba57600080fd5b9193909290916020810190356401000000008111156104d857600080fd5b8201836020820111156104ea57600080fd5b8035906020019184602083028401116401000000008311171561050c57600080fd5b91939092909160208101903564010000000081111561052a57600080fd5b82018360208201111561053c57600080fd5b8035906020019184602083028401116401000000008311171561055e57600080fd5b9193509150356110ab565b6101af6004803603604081101561057f57600080fd5b506001600160a01b038135169060200135611270565b6101af611394565b6101af600480360360408110156105b357600080fd5b506001600160a01b03813516906020013561144f565b6105d1611516565b604080516001600160a01b039092168252519081900360200190f35b6105f5611525565b604080519115158252519081900360200190f35b6101af6004803603604081101561061f57600080fd5b5080359060200135611536565b6101af6004803603602081101561064257600080fd5b81019060208101813564010000000081111561065d57600080fd5b82018360208201111561066f57600080fd5b8035906020019184600183028401116401000000008311171561069157600080fd5b509092509050611602565b6101af600480360360608110156106b257600080fd5b506001600160a01b0381358116916020810135821691604090910135166116eb565b6101af600480360360408110156106ea57600080fd5b506001600160a01b0381358116916020013516611838565b6101af6004803603602081101561071857600080fd5b81019060208101813564010000000081111561073357600080fd5b82018360208201111561074557600080fd5b8035906020019184602083028401116401000000008311171561076757600080fd5b509092509050611900565b6101af6004803603602081101561078857600080fd5b50356001600160a01b03166119df565b6101af600480360360408110156107ae57600080fd5b506001600160a01b0381358116916020013516611a41565b6066546001600160a01b03163314610825576040805162461bcd60e51b815260206004820152601e60248201527f63616c6c6572206973206e6f74207468652053464320636f6e74726163740000604482015290519081900360640190fd5b60675460408051808201909152601a81527f7b2245636f6e6f6d79223a7b224d696e4761735072696365223a00000000000060208201526001600160a01b039091169063b9cc6b1c906108b69061087b85611b04565b6040518060400160405280600281526020017f7d7d000000000000000000000000000000000000000000000000000000000000815250611c28565b6040518263ffffffff1660e01b81526004018080602001828103825283818151815260200191508051906020019080838360005b838110156109025781810151838201526020016108ea565b50505050905090810190601f16801561092f5780820380516001836020036101000a031916815260200191505b5092505050600060405180830381600087803b15801561094e57600080fd5b505af1158015610962573d6000803e3d6000fd5b5050505050565b610971611525565b6109c2576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b606754604080517f0aeeca000000000000000000000000000000000000000000000000000000000081526004810184905290516001600160a01b0390921691630aeeca009160248082019260009290919082900301818387803b15801561094e57600080fd5b6067546001600160a01b03163314610a715760405162461bcd60e51b81526004018080602001828103825260258152602001806122106025913960400191505060405180910390fd5b606654604080517f18f628d40000000000000000000000000000000000000000000000000000000081526001600160a01b038c81166004830152602482018c9052604482018b9052606482018a90526084820189905260a4820188905260c4820187905260e482018690526101048201859052915191909216916318f628d49161012480830192600092919082900301818387803b158015610b1257600080fd5b505af1158015610b26573d6000803e3d6000fd5b50505050505050505050505050565b610b3d611525565b610b8e576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b610b9a84848484611dc6565b50505050565b6067546001600160a01b03163314610be95760405162461bcd60e51b81526004018080602001828103825260258152602001806122106025913960400191505060405180910390fd5b606654604080517f1e702f83000000000000000000000000000000000000000000000000000000008152600481018590526024810184905290516001600160a01b0390921691631e702f839160448082019260009290919082900301818387803b158015610c5657600080fd5b505af1158015610c6a573d6000803e3d6000fd5b505050505050565b6066546001600160a01b03163314610cd1576040805162461bcd60e51b815260206004820152601e60248201527f63616c6c6572206973206e6f74207468652053464320636f6e74726163740000604482015290519081900360640190fd5b606754604080517f242a6e3f0000000000000000000000000000000000000000000000000000000081526004810186815260248201928352604482018590526001600160a01b039093169263242a6e3f928792879287929091606401848480828437600081840152601f19601f820116905080830192505050945050505050600060405180830381600087803b158015610d6a57600080fd5b505af1158015610d7e573d6000803e3d6000fd5b50505050505050565b610d8f611525565b610de0576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b606754604080517f267ab4460000000000000000000000000000000000000000000000000000000081526004810184905290516001600160a01b039092169163267ab4469160248082019260009290919082900301818387803b15801561094e57600080fd5b610e4e611525565b610e9f576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b610ece81610eab611516565b610eb430611ef0565b606754610ec9906001600160a01b0316611ef0565b611dc6565b50565b610ed9611525565b610f2a576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b606754604080517fda7fc24f0000000000000000000000000000000000000000000000000000000081526001600160a01b0384811660048301529151919092169163da7fc24f91602480830192600092919082900301818387803b15801561094e57600080fd5b6067546001600160a01b03163314610fda5760405162461bcd60e51b81526004018080602001828103825260258152602001806122106025913960400191505060405180910390fd5b606660009054906101000a90046001600160a01b03166001600160a01b0316634feb92f38a8a8a8a8a8a8a8a8a6040518a63ffffffff1660e01b8152600401808a6001600160a01b03166001600160a01b03168152602001898152602001806020018781526020018681526020018581526020018481526020018381526020018281038252898982818152602001925080828437600081840152601f19601f8201169050808301925050509a5050505050505050505050600060405180830381600087803b158015610b1257600080fd5b6067546001600160a01b031633146110f45760405162461bcd60e51b81526004018080602001828103825260258152602001806122106025913960400191505060405180910390fd5b606660009054906101000a90046001600160a01b03166001600160a01b031663592fe0c08a8a8a8a8a8a8a8a8a6040518a63ffffffff1660e01b8152600401808060200180602001806020018060200186815260200185810385528e8e82818152602001925060200280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910186810385528c8152602090810191508d908d0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910186810384528a8152602090810191508b908b0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169091018681038352888152602090810191508990890280828437600081840152601f19601f8201169050808301925050509d5050505050505050505050505050600060405180830381600087803b158015610b1257600080fd5b6066546001600160a01b031633146112cf576040805162461bcd60e51b815260206004820152601e60248201527f63616c6c6572206973206e6f74207468652053464320636f6e74726163740000604482015290519081900360640190fd5b6066546001600160a01b0383811691161461131b5760405162461bcd60e51b81526004018080602001828103825260218152602001806121ef6021913960400191505060405180910390fd5b6067546001600160a01b039081169063e30443bc908490611345908216318563ffffffff611ef416565b6040518363ffffffff1660e01b815260040180836001600160a01b03166001600160a01b0316815260200182815260200192505050600060405180830381600087803b158015610c5657600080fd5b61139c611525565b6113ed576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6033546040516000916001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3603380547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055565b611457611525565b6114a8576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b606754604080517f79bead380000000000000000000000000000000000000000000000000000000081526001600160a01b03858116600483015260248201859052915191909216916379bead3891604480830192600092919082900301818387803b158015610c5657600080fd5b6033546001600160a01b031690565b6033546001600160a01b0316331490565b6066546001600160a01b03163314611595576040805162461bcd60e51b815260206004820152601e60248201527f63616c6c6572206973206e6f74207468652053464320636f6e74726163740000604482015290519081900360640190fd5b606754604080517fa4066fbe000000000000000000000000000000000000000000000000000000008152600481018590526024810184905290516001600160a01b039092169163a4066fbe9160448082019260009290919082900301818387803b158015610c5657600080fd5b61160a611525565b61165b576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6067546040517fb9cc6b1c000000000000000000000000000000000000000000000000000000008152602060048201908152602482018490526001600160a01b039092169163b9cc6b1c91859185918190604401848480828437600081840152601f19601f8201169050808301925050509350505050600060405180830381600087803b158015610c5657600080fd5b600054610100900460ff16806117045750611704611f55565b80611712575060005460ff16155b61174d5760405162461bcd60e51b815260040180806020018281038252602e8152602001806121c1602e913960400191505060405180910390fd5b600054610100900460ff161580156117b357600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff909116610100171660011790555b6117bc82611f5b565b606780546001600160a01b038086167fffffffffffffffffffffffff00000000000000000000000000000000000000009283161790925560668054928716929091169190911790558015610b9a57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16905550505050565b611840611525565b611891576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b606754604080517fd6a0c7af0000000000000000000000000000000000000000000000000000000081526001600160a01b03858116600483015284811660248301529151919092169163d6a0c7af91604480830192600092919082900301818387803b158015610c5657600080fd5b6067546001600160a01b031633146119495760405162461bcd60e51b81526004018080602001828103825260258152602001806122106025913960400191505060405180910390fd5b6066546040517fe08d7e66000000000000000000000000000000000000000000000000000000008152602060048201818152602483018590526001600160a01b039093169263e08d7e6692869286929182916044909101908590850280828437600081840152601f19601f8201169050808301925050509350505050600060405180830381600087803b158015610c5657600080fd5b6119e7611525565b611a38576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b610ece816120bd565b611a49611525565b611a9a576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b611aa382612176565b8015611ab35750611ab381612176565b611891576040805162461bcd60e51b815260206004820152600e60248201527f6e6f74206120636f6e7472616374000000000000000000000000000000000000604482015290519081900360640190fd5b606081611b45575060408051808201909152600181527f30000000000000000000000000000000000000000000000000000000000000006020820152611c23565b6000611b508361217c565b90506060816040519080825280601f01601f191660200182016040528015611b7f576020820181803883390190505b5090507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82015b8415611c1e57600a850660300160f81b828281518110611bc257fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a850494507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01611ba6565b509150505b919050565b60608084905060608490506060849050606081518351855101016040519080825280601f01601f191660200182016040528015611c6c576020820181803883390190505b509050806000805b8651811015611cdd57868181518110611c8957fe5b602001015160f81c60f81b838380600101945081518110611ca657fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600101611c74565b5060005b8551811015611d4a57858181518110611cf657fe5b602001015160f81c60f81b838380600101945081518110611d1357fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600101611ce1565b5060005b8451811015611db757848181518110611d6357fe5b602001015160f81c60f81b838380600101945081518110611d8057fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600101611d4e565b50909998505050505050505050565b611dcf846120bd565b836001600160a01b031663614619546040518163ffffffff1660e01b8152600401600060405180830381600087803b158015611e0a57600080fd5b505af1158015611e1e573d6000803e3d6000fd5b50505050611e2b836120bd565b81611e3530611ef0565b14611e87576040805162461bcd60e51b815260206004820152601c60248201527f73656c6620636f6465206861736820646f65736e2774206d6174636800000000604482015290519081900360640190fd5b6067548190611e9e906001600160a01b0316611ef0565b14610b9a576040805162461bcd60e51b815260206004820152601e60248201527f64726976657220636f6465206861736820646f65736e2774206d617463680000604482015290519081900360640190fd5b3f90565b600082820183811015611f4e576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b303b1590565b600054610100900460ff1680611f745750611f74611f55565b80611f82575060005460ff16155b611fbd5760405162461bcd60e51b815260040180806020018281038252602e8152602001806121c1602e913960400191505060405180910390fd5b600054610100900460ff1615801561202357600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff909116610100171660011790555b603380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384811691909117918290556040519116906000907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a380156120b957600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690555b5050565b6001600160a01b0381166121025760405162461bcd60e51b815260040180806020018281038252602681526020018061219b6026913960400191505060405180910390fd5b6033546040516001600160a01b038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3603380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b3b151590565b6000805b821561219457600101600a83049250612180565b9291505056fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373436f6e747261637420696e7374616e63652068617320616c7265616479206265656e20696e697469616c697a6564726563697069656e74206973206e6f74207468652053464320636f6e747261637463616c6c6572206973206e6f7420746865204e6f646544726976657220636f6e7472616374a265627a7a723158207c2ee67320c429dee1c46e39a39b15221704649befc2e0e47d40ef199c3bbe1864736f6c63430005110032",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Contract *ContractCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Contract *ContractSession) IsOwner() (bool, error) {
	return _Contract.Contract.IsOwner(&_Contract.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Contract *ContractCallerSession) IsOwner() (bool, error) {
	return _Contract.Contract.IsOwner(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCallerSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// AdvanceEpochs is a paid mutator transaction binding the contract method 0x0aeeca00.
//
// Solidity: function advanceEpochs(uint256 num) returns()
func (_Contract *ContractTransactor) AdvanceEpochs(opts *bind.TransactOpts, num *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "advanceEpochs", num)
}

// AdvanceEpochs is a paid mutator transaction binding the contract method 0x0aeeca00.
//
// Solidity: function advanceEpochs(uint256 num) returns()
func (_Contract *ContractSession) AdvanceEpochs(num *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.AdvanceEpochs(&_Contract.TransactOpts, num)
}

// AdvanceEpochs is a paid mutator transaction binding the contract method 0x0aeeca00.
//
// Solidity: function advanceEpochs(uint256 num) returns()
func (_Contract *ContractTransactorSession) AdvanceEpochs(num *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.AdvanceEpochs(&_Contract.TransactOpts, num)
}

// CopyCode is a paid mutator transaction binding the contract method 0xd6a0c7af.
//
// Solidity: function copyCode(address acc, address from) returns()
func (_Contract *ContractTransactor) CopyCode(opts *bind.TransactOpts, acc common.Address, from common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "copyCode", acc, from)
}

// CopyCode is a paid mutator transaction binding the contract method 0xd6a0c7af.
//
// Solidity: function copyCode(address acc, address from) returns()
func (_Contract *ContractSession) CopyCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _Contract.Contract.CopyCode(&_Contract.TransactOpts, acc, from)
}

// CopyCode is a paid mutator transaction binding the contract method 0xd6a0c7af.
//
// Solidity: function copyCode(address acc, address from) returns()
func (_Contract *ContractTransactorSession) CopyCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _Contract.Contract.CopyCode(&_Contract.TransactOpts, acc, from)
}

// DeactivateValidator is a paid mutator transaction binding the contract method 0x1e702f83.
//
// Solidity: function deactivateValidator(uint256 validatorID, uint256 status) returns()
func (_Contract *ContractTransactor) DeactivateValidator(opts *bind.TransactOpts, validatorID *big.Int, status *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "deactivateValidator", validatorID, status)
}

// DeactivateValidator is a paid mutator transaction binding the contract method 0x1e702f83.
//
// Solidity: function deactivateValidator(uint256 validatorID, uint256 status) returns()
func (_Contract *ContractSession) DeactivateValidator(validatorID *big.Int, status *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DeactivateValidator(&_Contract.TransactOpts, validatorID, status)
}

// DeactivateValidator is a paid mutator transaction binding the contract method 0x1e702f83.
//
// Solidity: function deactivateValidator(uint256 validatorID, uint256 status) returns()
func (_Contract *ContractTransactorSession) DeactivateValidator(validatorID *big.Int, status *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DeactivateValidator(&_Contract.TransactOpts, validatorID, status)
}

// Execute is a paid mutator transaction binding the contract method 0x4b64e492.
//
// Solidity: function execute(address executable) returns()
func (_Contract *ContractTransactor) Execute(opts *bind.TransactOpts, executable common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "execute", executable)
}

// Execute is a paid mutator transaction binding the contract method 0x4b64e492.
//
// Solidity: function execute(address executable) returns()
func (_Contract *ContractSession) Execute(executable common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Execute(&_Contract.TransactOpts, executable)
}

// Execute is a paid mutator transaction binding the contract method 0x4b64e492.
//
// Solidity: function execute(address executable) returns()
func (_Contract *ContractTransactorSession) Execute(executable common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Execute(&_Contract.TransactOpts, executable)
}

// IncBalance is a paid mutator transaction binding the contract method 0x66e7ea0f.
//
// Solidity: function incBalance(address acc, uint256 diff) returns()
func (_Contract *ContractTransactor) IncBalance(opts *bind.TransactOpts, acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "incBalance", acc, diff)
}

// IncBalance is a paid mutator transaction binding the contract method 0x66e7ea0f.
//
// Solidity: function incBalance(address acc, uint256 diff) returns()
func (_Contract *ContractSession) IncBalance(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.IncBalance(&_Contract.TransactOpts, acc, diff)
}

// IncBalance is a paid mutator transaction binding the contract method 0x66e7ea0f.
//
// Solidity: function incBalance(address acc, uint256 diff) returns()
func (_Contract *ContractTransactorSession) IncBalance(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.IncBalance(&_Contract.TransactOpts, acc, diff)
}

// IncNonce is a paid mutator transaction binding the contract method 0x79bead38.
//
// Solidity: function incNonce(address acc, uint256 diff) returns()
func (_Contract *ContractTransactor) IncNonce(opts *bind.TransactOpts, acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "incNonce", acc, diff)
}

// IncNonce is a paid mutator transaction binding the contract method 0x79bead38.
//
// Solidity: function incNonce(address acc, uint256 diff) returns()
func (_Contract *ContractSession) IncNonce(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.IncNonce(&_Contract.TransactOpts, acc, diff)
}

// IncNonce is a paid mutator transaction binding the contract method 0x79bead38.
//
// Solidity: function incNonce(address acc, uint256 diff) returns()
func (_Contract *ContractTransactorSession) IncNonce(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.IncNonce(&_Contract.TransactOpts, acc, diff)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _sfc, address _driver, address _owner) returns()
func (_Contract *ContractTransactor) Initialize(opts *bind.TransactOpts, _sfc common.Address, _driver common.Address, _owner common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "initialize", _sfc, _driver, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _sfc, address _driver, address _owner) returns()
func (_Contract *ContractSession) Initialize(_sfc common.Address, _driver common.Address, _owner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Initialize(&_Contract.TransactOpts, _sfc, _driver, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _sfc, address _driver, address _owner) returns()
func (_Contract *ContractTransactorSession) Initialize(_sfc common.Address, _driver common.Address, _owner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Initialize(&_Contract.TransactOpts, _sfc, _driver, _owner)
}

// MigrateTo is a paid mutator transaction binding the contract method 0x4ddaf8f2.
//
// Solidity: function migrateTo(address newDriverAuth) returns()
func (_Contract *ContractTransactor) MigrateTo(opts *bind.TransactOpts, newDriverAuth common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "migrateTo", newDriverAuth)
}

// MigrateTo is a paid mutator transaction binding the contract method 0x4ddaf8f2.
//
// Solidity: function migrateTo(address newDriverAuth) returns()
func (_Contract *ContractSession) MigrateTo(newDriverAuth common.Address) (*types.Transaction, error) {
	return _Contract.Contract.MigrateTo(&_Contract.TransactOpts, newDriverAuth)
}

// MigrateTo is a paid mutator transaction binding the contract method 0x4ddaf8f2.
//
// Solidity: function migrateTo(address newDriverAuth) returns()
func (_Contract *ContractTransactorSession) MigrateTo(newDriverAuth common.Address) (*types.Transaction, error) {
	return _Contract.Contract.MigrateTo(&_Contract.TransactOpts, newDriverAuth)
}

// MutExecute is a paid mutator transaction binding the contract method 0x1cef4fab.
//
// Solidity: function mutExecute(address executable, address newOwner, bytes32 selfCodeHash, bytes32 driverCodeHash) returns()
func (_Contract *ContractTransactor) MutExecute(opts *bind.TransactOpts, executable common.Address, newOwner common.Address, selfCodeHash [32]byte, driverCodeHash [32]byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "mutExecute", executable, newOwner, selfCodeHash, driverCodeHash)
}

// MutExecute is a paid mutator transaction binding the contract method 0x1cef4fab.
//
// Solidity: function mutExecute(address executable, address newOwner, bytes32 selfCodeHash, bytes32 driverCodeHash) returns()
func (_Contract *ContractSession) MutExecute(executable common.Address, newOwner common.Address, selfCodeHash [32]byte, driverCodeHash [32]byte) (*types.Transaction, error) {
	return _Contract.Contract.MutExecute(&_Contract.TransactOpts, executable, newOwner, selfCodeHash, driverCodeHash)
}

// MutExecute is a paid mutator transaction binding the contract method 0x1cef4fab.
//
// Solidity: function mutExecute(address executable, address newOwner, bytes32 selfCodeHash, bytes32 driverCodeHash) returns()
func (_Contract *ContractTransactorSession) MutExecute(executable common.Address, newOwner common.Address, selfCodeHash [32]byte, driverCodeHash [32]byte) (*types.Transaction, error) {
	return _Contract.Contract.MutExecute(&_Contract.TransactOpts, executable, newOwner, selfCodeHash, driverCodeHash)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// SealEpoch is a paid mutator transaction binding the contract method 0x592fe0c0.
//
// Solidity: function sealEpoch(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee, uint256 usedGas) returns()
func (_Contract *ContractTransactor) SealEpoch(opts *bind.TransactOpts, offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int, usedGas *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "sealEpoch", offlineTimes, offlineBlocks, uptimes, originatedTxsFee, usedGas)
}

// SealEpoch is a paid mutator transaction binding the contract method 0x592fe0c0.
//
// Solidity: function sealEpoch(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee, uint256 usedGas) returns()
func (_Contract *ContractSession) SealEpoch(offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int, usedGas *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SealEpoch(&_Contract.TransactOpts, offlineTimes, offlineBlocks, uptimes, originatedTxsFee, usedGas)
}

// SealEpoch is a paid mutator transaction binding the contract method 0x592fe0c0.
//
// Solidity: function sealEpoch(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee, uint256 usedGas) returns()
func (_Contract *ContractTransactorSession) SealEpoch(offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int, usedGas *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SealEpoch(&_Contract.TransactOpts, offlineTimes, offlineBlocks, uptimes, originatedTxsFee, usedGas)
}

// SealEpochValidators is a paid mutator transaction binding the contract method 0xe08d7e66.
//
// Solidity: function sealEpochValidators(uint256[] nextValidatorIDs) returns()
func (_Contract *ContractTransactor) SealEpochValidators(opts *bind.TransactOpts, nextValidatorIDs []*big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "sealEpochValidators", nextValidatorIDs)
}

// SealEpochValidators is a paid mutator transaction binding the contract method 0xe08d7e66.
//
// Solidity: function sealEpochValidators(uint256[] nextValidatorIDs) returns()
func (_Contract *ContractSession) SealEpochValidators(nextValidatorIDs []*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SealEpochValidators(&_Contract.TransactOpts, nextValidatorIDs)
}

// SealEpochValidators is a paid mutator transaction binding the contract method 0xe08d7e66.
//
// Solidity: function sealEpochValidators(uint256[] nextValidatorIDs) returns()
func (_Contract *ContractTransactorSession) SealEpochValidators(nextValidatorIDs []*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SealEpochValidators(&_Contract.TransactOpts, nextValidatorIDs)
}

// SetGenesisDelegation is a paid mutator transaction binding the contract method 0x18f628d4.
//
// Solidity: function setGenesisDelegation(address delegator, uint256 toValidatorID, uint256 stake, uint256 lockedStake, uint256 lockupFromEpoch, uint256 lockupEndTime, uint256 lockupDuration, uint256 earlyUnlockPenalty, uint256 rewards) returns()
func (_Contract *ContractTransactor) SetGenesisDelegation(opts *bind.TransactOpts, delegator common.Address, toValidatorID *big.Int, stake *big.Int, lockedStake *big.Int, lockupFromEpoch *big.Int, lockupEndTime *big.Int, lockupDuration *big.Int, earlyUnlockPenalty *big.Int, rewards *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setGenesisDelegation", delegator, toValidatorID, stake, lockedStake, lockupFromEpoch, lockupEndTime, lockupDuration, earlyUnlockPenalty, rewards)
}

// SetGenesisDelegation is a paid mutator transaction binding the contract method 0x18f628d4.
//
// Solidity: function setGenesisDelegation(address delegator, uint256 toValidatorID, uint256 stake, uint256 lockedStake, uint256 lockupFromEpoch, uint256 lockupEndTime, uint256 lockupDuration, uint256 earlyUnlockPenalty, uint256 rewards) returns()
func (_Contract *ContractSession) SetGenesisDelegation(delegator common.Address, toValidatorID *big.Int, stake *big.Int, lockedStake *big.Int, lockupFromEpoch *big.Int, lockupEndTime *big.Int, lockupDuration *big.Int, earlyUnlockPenalty *big.Int, rewards *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetGenesisDelegation(&_Contract.TransactOpts, delegator, toValidatorID, stake, lockedStake, lockupFromEpoch, lockupEndTime, lockupDuration, earlyUnlockPenalty, rewards)
}

// SetGenesisDelegation is a paid mutator transaction binding the contract method 0x18f628d4.
//
// Solidity: function setGenesisDelegation(address delegator, uint256 toValidatorID, uint256 stake, uint256 lockedStake, uint256 lockupFromEpoch, uint256 lockupEndTime, uint256 lockupDuration, uint256 earlyUnlockPenalty, uint256 rewards) returns()
func (_Contract *ContractTransactorSession) SetGenesisDelegation(delegator common.Address, toValidatorID *big.Int, stake *big.Int, lockedStake *big.Int, lockupFromEpoch *big.Int, lockupEndTime *big.Int, lockupDuration *big.Int, earlyUnlockPenalty *big.Int, rewards *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetGenesisDelegation(&_Contract.TransactOpts, delegator, toValidatorID, stake, lockedStake, lockupFromEpoch, lockupEndTime, lockupDuration, earlyUnlockPenalty, rewards)
}

// SetGenesisValidator is a paid mutator transaction binding the contract method 0x4feb92f3.
//
// Solidity: function setGenesisValidator(address _auth, uint256 validatorID, bytes pubkey, uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime) returns()
func (_Contract *ContractTransactor) SetGenesisValidator(opts *bind.TransactOpts, _auth common.Address, validatorID *big.Int, pubkey []byte, status *big.Int, createdEpoch *big.Int, createdTime *big.Int, deactivatedEpoch *big.Int, deactivatedTime *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setGenesisValidator", _auth, validatorID, pubkey, status, createdEpoch, createdTime, deactivatedEpoch, deactivatedTime)
}

// SetGenesisValidator is a paid mutator transaction binding the contract method 0x4feb92f3.
//
// Solidity: function setGenesisValidator(address _auth, uint256 validatorID, bytes pubkey, uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime) returns()
func (_Contract *ContractSession) SetGenesisValidator(_auth common.Address, validatorID *big.Int, pubkey []byte, status *big.Int, createdEpoch *big.Int, createdTime *big.Int, deactivatedEpoch *big.Int, deactivatedTime *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetGenesisValidator(&_Contract.TransactOpts, _auth, validatorID, pubkey, status, createdEpoch, createdTime, deactivatedEpoch, deactivatedTime)
}

// SetGenesisValidator is a paid mutator transaction binding the contract method 0x4feb92f3.
//
// Solidity: function setGenesisValidator(address _auth, uint256 validatorID, bytes pubkey, uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime) returns()
func (_Contract *ContractTransactorSession) SetGenesisValidator(_auth common.Address, validatorID *big.Int, pubkey []byte, status *big.Int, createdEpoch *big.Int, createdTime *big.Int, deactivatedEpoch *big.Int, deactivatedTime *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetGenesisValidator(&_Contract.TransactOpts, _auth, validatorID, pubkey, status, createdEpoch, createdTime, deactivatedEpoch, deactivatedTime)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// UpdateMinGasPrice is a paid mutator transaction binding the contract method 0x07aaf344.
//
// Solidity: function updateMinGasPrice(uint256 minGasPrice) returns()
func (_Contract *ContractTransactor) UpdateMinGasPrice(opts *bind.TransactOpts, minGasPrice *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateMinGasPrice", minGasPrice)
}

// UpdateMinGasPrice is a paid mutator transaction binding the contract method 0x07aaf344.
//
// Solidity: function updateMinGasPrice(uint256 minGasPrice) returns()
func (_Contract *ContractSession) UpdateMinGasPrice(minGasPrice *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateMinGasPrice(&_Contract.TransactOpts, minGasPrice)
}

// UpdateMinGasPrice is a paid mutator transaction binding the contract method 0x07aaf344.
//
// Solidity: function updateMinGasPrice(uint256 minGasPrice) returns()
func (_Contract *ContractTransactorSession) UpdateMinGasPrice(minGasPrice *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateMinGasPrice(&_Contract.TransactOpts, minGasPrice)
}

// UpdateNetworkRules is a paid mutator transaction binding the contract method 0xb9cc6b1c.
//
// Solidity: function updateNetworkRules(bytes diff) returns()
func (_Contract *ContractTransactor) UpdateNetworkRules(opts *bind.TransactOpts, diff []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateNetworkRules", diff)
}

// UpdateNetworkRules is a paid mutator transaction binding the contract method 0xb9cc6b1c.
//
// Solidity: function updateNetworkRules(bytes diff) returns()
func (_Contract *ContractSession) UpdateNetworkRules(diff []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpdateNetworkRules(&_Contract.TransactOpts, diff)
}

// UpdateNetworkRules is a paid mutator transaction binding the contract method 0xb9cc6b1c.
//
// Solidity: function updateNetworkRules(bytes diff) returns()
func (_Contract *ContractTransactorSession) UpdateNetworkRules(diff []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpdateNetworkRules(&_Contract.TransactOpts, diff)
}

// UpdateNetworkVersion is a paid mutator transaction binding the contract method 0x267ab446.
//
// Solidity: function updateNetworkVersion(uint256 version) returns()
func (_Contract *ContractTransactor) UpdateNetworkVersion(opts *bind.TransactOpts, version *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateNetworkVersion", version)
}

// UpdateNetworkVersion is a paid mutator transaction binding the contract method 0x267ab446.
//
// Solidity: function updateNetworkVersion(uint256 version) returns()
func (_Contract *ContractSession) UpdateNetworkVersion(version *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateNetworkVersion(&_Contract.TransactOpts, version)
}

// UpdateNetworkVersion is a paid mutator transaction binding the contract method 0x267ab446.
//
// Solidity: function updateNetworkVersion(uint256 version) returns()
func (_Contract *ContractTransactorSession) UpdateNetworkVersion(version *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateNetworkVersion(&_Contract.TransactOpts, version)
}

// UpdateValidatorPubkey is a paid mutator transaction binding the contract method 0x242a6e3f.
//
// Solidity: function updateValidatorPubkey(uint256 validatorID, bytes pubkey) returns()
func (_Contract *ContractTransactor) UpdateValidatorPubkey(opts *bind.TransactOpts, validatorID *big.Int, pubkey []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateValidatorPubkey", validatorID, pubkey)
}

// UpdateValidatorPubkey is a paid mutator transaction binding the contract method 0x242a6e3f.
//
// Solidity: function updateValidatorPubkey(uint256 validatorID, bytes pubkey) returns()
func (_Contract *ContractSession) UpdateValidatorPubkey(validatorID *big.Int, pubkey []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpdateValidatorPubkey(&_Contract.TransactOpts, validatorID, pubkey)
}

// UpdateValidatorPubkey is a paid mutator transaction binding the contract method 0x242a6e3f.
//
// Solidity: function updateValidatorPubkey(uint256 validatorID, bytes pubkey) returns()
func (_Contract *ContractTransactorSession) UpdateValidatorPubkey(validatorID *big.Int, pubkey []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpdateValidatorPubkey(&_Contract.TransactOpts, validatorID, pubkey)
}

// UpdateValidatorWeight is a paid mutator transaction binding the contract method 0xa4066fbe.
//
// Solidity: function updateValidatorWeight(uint256 validatorID, uint256 value) returns()
func (_Contract *ContractTransactor) UpdateValidatorWeight(opts *bind.TransactOpts, validatorID *big.Int, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateValidatorWeight", validatorID, value)
}

// UpdateValidatorWeight is a paid mutator transaction binding the contract method 0xa4066fbe.
//
// Solidity: function updateValidatorWeight(uint256 validatorID, uint256 value) returns()
func (_Contract *ContractSession) UpdateValidatorWeight(validatorID *big.Int, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateValidatorWeight(&_Contract.TransactOpts, validatorID, value)
}

// UpdateValidatorWeight is a paid mutator transaction binding the contract method 0xa4066fbe.
//
// Solidity: function updateValidatorWeight(uint256 validatorID, uint256 value) returns()
func (_Contract *ContractTransactorSession) UpdateValidatorWeight(validatorID *big.Int, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateValidatorWeight(&_Contract.TransactOpts, validatorID, value)
}

// UpgradeCode is a paid mutator transaction binding the contract method 0xfd1b6ec1.
//
// Solidity: function upgradeCode(address acc, address from) returns()
func (_Contract *ContractTransactor) UpgradeCode(opts *bind.TransactOpts, acc common.Address, from common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "upgradeCode", acc, from)
}

// UpgradeCode is a paid mutator transaction binding the contract method 0xfd1b6ec1.
//
// Solidity: function upgradeCode(address acc, address from) returns()
func (_Contract *ContractSession) UpgradeCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _Contract.Contract.UpgradeCode(&_Contract.TransactOpts, acc, from)
}

// UpgradeCode is a paid mutator transaction binding the contract method 0xfd1b6ec1.
//
// Solidity: function upgradeCode(address acc, address from) returns()
func (_Contract *ContractTransactorSession) UpgradeCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _Contract.Contract.UpgradeCode(&_Contract.TransactOpts, acc, from)
}

// ContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Contract contract.
type ContractOwnershipTransferredIterator struct {
	Event *ContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractOwnershipTransferred)
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
		it.Event = new(ContractOwnershipTransferred)
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
func (it *ContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractOwnershipTransferred represents a OwnershipTransferred event raised by the Contract contract.
type ContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractOwnershipTransferredIterator{contract: _Contract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractOwnershipTransferred)
				if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) ParseOwnershipTransferred(log types.Log) (*ContractOwnershipTransferred, error) {
	event := new(ContractOwnershipTransferred)
	if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

var ContractBinRuntime = "0x60806040526004361061013c575f3560e01c806379bead38116100b3578063c0c53b8b1161006d578063c0c53b8b1461036d578063d6a0c7af1461038c578063e08d7e66146103ab578063ebdf104c146103ca578063f2fde38b146103e9578063fd1b6ec114610408575f5ffd5b806379bead38146102885780638da5cb5b146102a7578063a4066fbe146102d3578063a8ab09ba146102f2578063ad3cb1cc14610311578063b9cc6b1c1461034e575f5ffd5b80634b64e492116101045780634b64e492146101dd5780634f1ef286146101fc57806352d1902d1461020f57806366e7ea0f14610236578063715018a61461025557806376fed43a14610269575f5ffd5b80630aeeca00146101405780631cef4fab146101615780631e702f8314610180578063242a6e3f1461019f578063267ab446146101be575b5f5ffd5b34801561014b575f5ffd5b5061015f61015a366004611163565b610427565b005b34801561016c575f5ffd5b5061015f61017b36600461118e565b61048b565b34801561018b575f5ffd5b5061015f61019a3660046111d1565b6104a5565b3480156101aa575f5ffd5b5061015f6101b9366004611235565b610534565b3480156101c9575f5ffd5b5061015f6101d8366004611163565b6105c4565b3480156101e8575f5ffd5b5061015f6101f736600461127c565b6105fd565b61015f61020a3660046112ab565b610628565b34801561021a575f5ffd5b50610223610647565b6040519081526020015b60405180910390f35b348015610241575f5ffd5b5061015f61025036600461136e565b610662565b348015610260575f5ffd5b5061015f61070f565b348015610274575f5ffd5b5061015f610283366004611398565b610722565b348015610293575f5ffd5b5061015f6102a236600461136e565b6107b8565b3480156102b2575f5ffd5b506102bb6107f9565b6040516001600160a01b03909116815260200161022d565b3480156102de575f5ffd5b5061015f6102ed3660046111d1565b610827565b3480156102fd575f5ffd5b5061015f61030c3660046113f6565b610889565b34801561031c575f5ffd5b50610341604051806040016040528060058152602001640352e302e360dc1b81525081565b60405161022d9190611428565b348015610359575f5ffd5b5061015f61036836600461145d565b6108f3565b348015610378575f5ffd5b5061015f61038736600461149b565b61092d565b348015610397575f5ffd5b5061015f6103a63660046114e3565b610a72565b3480156103b6575f5ffd5b5061015f6103c536600461155a565b610ab4565b3480156103d5575f5ffd5b5061015f6103e436600461158c565b610b10565b3480156103f4575f5ffd5b5061015f61040336600461127c565b610baf565b348015610413575f5ffd5b5061015f6104223660046114e3565b610bee565b61042f610c21565b6001546040516205776560e91b8152600481018390526001600160a01b0390911690630aeeca00906024015b5f604051808303815f87803b158015610472575f5ffd5b505af1158015610484573d5f5f3e3d5ffd5b5050505050565b610493610c21565b61049f84848484610c53565b50505050565b6001546001600160a01b031633146104d057604051630607323760e11b815260040160405180910390fd5b5f54604051631e702f8360e01b815260048101849052602481018390526001600160a01b0390911690631e702f83906044015b5f604051808303815f87803b15801561051a575f5ffd5b505af115801561052c573d5f5f3e3d5ffd5b505050505050565b5f546001600160a01b0316331461055e5760405163d42fccad60e01b815260040160405180910390fd5b60015460405163242a6e3f60e01b81526001600160a01b039091169063242a6e3f906105929086908690869060040161167e565b5f604051808303815f87803b1580156105a9575f5ffd5b505af11580156105bb573d5f5f3e3d5ffd5b50505050505050565b6105cc610c21565b60015460405163133d5a2360e11b8152600481018390526001600160a01b039091169063267ab4469060240161045b565b610605610c21565b610625816106116107f9565b303f6001546001600160a01b03163f610c53565b50565b610630610d00565b61063982610da4565b6106438282610dac565b5050565b5f610650610e6d565b505f5160206117e35f395f51905f5290565b5f546001600160a01b0316331461068c5760405163d42fccad60e01b815260040160405180910390fd5b5f546001600160a01b038381169116146106b957604051630ea42ef960e31b815260040160405180910390fd5b6001546001600160a01b039081169063e30443bc9084906106de9085908316316116a0565b6040516001600160e01b031960e085901b1681526001600160a01b0390921660048301526024820152604401610503565b610717610c21565b6107205f610eb6565b565b6001546001600160a01b0316331461074d57604051630607323760e11b815260040160405180910390fd5b5f54604051633b7f6a1d60e11b81526001600160a01b03909116906376fed43a9061078490889088908890889088906004016116bf565b5f604051808303815f87803b15801561079b575f5ffd5b505af11580156107ad573d5f5f3e3d5ffd5b505050505050505050565b6107c0610c21565b600154604051630f37d5a760e31b81526001600160a01b03848116600483015260248201849052909116906379bead3890604401610503565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b5f546001600160a01b031633146108515760405163d42fccad60e01b815260040160405180910390fd5b60015460405163520337df60e11b815260048101849052602481018390526001600160a01b039091169063a4066fbe90604401610503565b6001546001600160a01b031633146108b457604051630607323760e11b815260040160405180910390fd5b5f5460405163545584dd60e11b81526001600160a01b03858116600483015260248201859052604482018490529091169063a8ab09ba90606401610592565b6108fb610c21565b600154604051632e731ac760e21b81526001600160a01b039091169063b9cc6b1c9061050390859085906004016116f8565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff1615906001600160401b03165f811580156109715750825b90505f826001600160401b0316600114801561098c5750303b155b90508115801561099a575080155b156109b85760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156109e257845460ff60401b1916600160401b1785555b6109eb86610f26565b6109f3610f37565b600180546001600160a01b03808a166001600160a01b0319928316179092555f8054928b16929091169190911790558315610a6857845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050505050565b610a7a610c21565b60015460405163d6a0c7af60e01b81526001600160a01b03848116600483015283811660248301529091169063d6a0c7af90604401610503565b6001546001600160a01b03163314610adf57604051630607323760e11b815260040160405180910390fd5b5f54604051637046bf3360e11b81526001600160a01b039091169063e08d7e66906105039085908590600401611743565b6001546001600160a01b03163314610b3b57604051630607323760e11b815260040160405180910390fd5b5f54604051633af7c41360e21b81526001600160a01b039091169063ebdf104c90610b78908b908b908b908b908b908b908b908b90600401611756565b5f604051808303815f87803b158015610b8f575f5ffd5b505af1158015610ba1573d5f5f3e3d5ffd5b505050505050505050505050565b610bb7610c21565b6001600160a01b038116610be557604051631e4fbdf760e01b81525f60048201526024015b60405180910390fd5b61062581610eb6565b610bf6610c21565b813b1580610c035750803b155b15610a7a57604051636f7c43f160e01b815260040160405180910390fd5b33610c2a6107f9565b6001600160a01b0316146107205760405163118cdaa760e01b8152336004820152602401610bdc565b610c5c84610eb6565b836001600160a01b031663614619546040518163ffffffff1660e01b81526004015f604051808303815f87803b158015610c94575f5ffd5b505af1158015610ca6573d5f5f3e3d5ffd5b50505050610cb383610eb6565b81303f14610cd4576040516311387fef60e21b815260040160405180910390fd5b6001546001600160a01b03163f811461049f5760405163f0c300ef60e01b815260040160405180910390fd5b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480610d8657507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610d7a5f5160206117e35f395f51905f52546001600160a01b031690565b6001600160a01b031614155b156107205760405163703e46dd60e11b815260040160405180910390fd5b610625610c21565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610e06575060408051601f3d908101601f19168201909252610e03918101906117b5565b60015b610e2e57604051634c9c8ce360e01b81526001600160a01b0383166004820152602401610bdc565b5f5160206117e35f395f51905f528114610e5e57604051632a87526960e21b815260048101829052602401610bdc565b610e688383610f3f565b505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107205760405163703e46dd60e11b815260040160405180910390fd5b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b610f2e610f94565b61062581610fdd565b610720610f94565b610f4882610fe5565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115610f8c57610e688282611048565b6106436110bc565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054600160401b900460ff1661072057604051631afcd79f60e31b815260040160405180910390fd5b610bb7610f94565b806001600160a01b03163b5f0361101a57604051634c9c8ce360e01b81526001600160a01b0382166004820152602401610bdc565b5f5160206117e35f395f51905f5280546001600160a01b0319166001600160a01b0392909216919091179055565b60605f5f846001600160a01b03168460405161106491906117cc565b5f60405180830381855af49150503d805f811461109c576040519150601f19603f3d011682016040523d82523d5f602084013e6110a1565b606091505b50915091506110b18583836110db565b925050505b92915050565b34156107205760405163b398979f60e01b815260040160405180910390fd5b6060826110f0576110eb8261113a565b611133565b815115801561110757506001600160a01b0384163b155b1561113057604051639996b31560e01b81526001600160a01b0385166004820152602401610bdc565b50805b9392505050565b80511561114a5780518082602001fd5b60405163d6bda27560e01b815260040160405180910390fd5b5f60208284031215611173575f5ffd5b5035919050565b6001600160a01b0381168114610625575f5ffd5b5f5f5f5f608085870312156111a1575f5ffd5b84356111ac8161117a565b935060208501356111bc8161117a565b93969395505050506040820135916060013590565b5f5f604083850312156111e2575f5ffd5b50508035926020909101359150565b5f5f83601f840112611201575f5ffd5b5081356001600160401b03811115611217575f5ffd5b60208301915083602082850101111561122e575f5ffd5b9250929050565b5f5f5f60408486031215611247575f5ffd5b8335925060208401356001600160401b03811115611263575f5ffd5b61126f868287016111f1565b9497909650939450505050565b5f6020828403121561128c575f5ffd5b81356111338161117a565b634e487b7160e01b5f52604160045260245ffd5b5f5f604083850312156112bc575f5ffd5b82356112c78161117a565b915060208301356001600160401b038111156112e1575f5ffd5b8301601f810185136112f1575f5ffd5b80356001600160401b0381111561130a5761130a611297565b604051601f8201601f19908116603f011681016001600160401b038111828210171561133857611338611297565b60405281815282820160200187101561134f575f5ffd5b816020840160208301375f602083830101528093505050509250929050565b5f5f6040838503121561137f575f5ffd5b823561138a8161117a565b946020939093013593505050565b5f5f5f5f5f608086880312156113ac575f5ffd5b85356113b78161117a565b94506020860135935060408601356001600160401b038111156113d8575f5ffd5b6113e4888289016111f1565b96999598509660600135949350505050565b5f5f5f60608486031215611408575f5ffd5b83356114138161117a565b95602085013595506040909401359392505050565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b5f5f6020838503121561146e575f5ffd5b82356001600160401b03811115611483575f5ffd5b61148f858286016111f1565b90969095509350505050565b5f5f5f606084860312156114ad575f5ffd5b83356114b88161117a565b925060208401356114c88161117a565b915060408401356114d88161117a565b809150509250925092565b5f5f604083850312156114f4575f5ffd5b82356114ff8161117a565b9150602083013561150f8161117a565b809150509250929050565b5f5f83601f84011261152a575f5ffd5b5081356001600160401b03811115611540575f5ffd5b6020830191508360208260051b850101111561122e575f5ffd5b5f5f6020838503121561156b575f5ffd5b82356001600160401b03811115611580575f5ffd5b61148f8582860161151a565b5f5f5f5f5f5f5f5f6080898b0312156115a3575f5ffd5b88356001600160401b038111156115b8575f5ffd5b6115c48b828c0161151a565b90995097505060208901356001600160401b038111156115e2575f5ffd5b6115ee8b828c0161151a565b90975095505060408901356001600160401b0381111561160c575f5ffd5b6116188b828c0161151a565b90955093505060608901356001600160401b03811115611636575f5ffd5b6116428b828c0161151a565b999c989b5096995094979396929594505050565b81835281816020850137505f828201602090810191909152601f909101601f19169091010190565b838152604060208201525f611697604083018486611656565b95945050505050565b808201808211156110b657634e487b7160e01b5f52601160045260245ffd5b60018060a01b0386168152846020820152608060408201525f6116e6608083018587611656565b90508260608301529695505050505050565b602081525f61170b602083018486611656565b949350505050565b8183525f6001600160fb1b0383111561172a575f5ffd5b8260051b80836020870137939093016020019392505050565b602081525f61170b602083018486611713565b608081525f611769608083018a8c611713565b828103602084015261177c81898b611713565b90508281036040840152611791818789611713565b905082810360608401526117a6818587611713565b9b9a5050505050505050505050565b5f602082840312156117c5575f5ffd5b5051919050565b5f82518060208501845e5f92019182525091905056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca264697066735822122068794a549ea260f5e8f76d22f28e70332219328d651409bff011f20bc678ea9864736f6c634300081b0033"
