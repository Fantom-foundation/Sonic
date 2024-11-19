package driver

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GetContractBin is NodeDriver contract genesis implementation bin code
// Has to be compiled with flag bin-runtime
// Built from opera-sfc e86e5db95f98965f4489ad962565a3850126023a, solc 0.8.27, optimize-runs 200, cancun evmVersion
func GetContractBin() []byte {
	return hexutil.MustDecode("0x608060405260043610610131575f3560e01c80638da5cb5b116100a8578063c0c53b8b1161006d578063c0c53b8b1461035d578063d6a0c7af1461037c578063e08d7e661461039b578063e30443bc146103ba578063ebdf104c146103d9578063f2fde38b146103f8575f5ffd5b80638da5cb5b1461027d578063a4066fbe146102c3578063a8ab09ba146102e2578063ad3cb1cc14610301578063b9cc6b1c1461033e575f5ffd5b806339e503ab116100f957806339e503ab146101d25780634f1ef286146101f157806352d1902d14610204578063715018a61461022b57806376fed43a1461023f57806379bead381461025e575f5ffd5b806307690b2a146101355780630aeeca00146101565780631e702f8314610175578063242a6e3f14610194578063267ab446146101b3575b5f5ffd5b348015610140575f5ffd5b5061015461014f366004611129565b610417565b005b348015610161575f5ffd5b5061015461017036600461115a565b6104a8565b348015610180575f5ffd5b5061015461018f366004611171565b610509565b34801561019f575f5ffd5b506101546101ae3660046111d5565b61055f565b3480156101be575f5ffd5b506101546101cd36600461115a565b6105c8565b3480156101dd575f5ffd5b506101546101ec36600461121c565b610622565b6101546101ff366004611260565b6106ba565b34801561020f575f5ffd5b506102186106d9565b6040519081526020015b60405180910390f35b348015610236575f5ffd5b506101546106f4565b34801561024a575f5ffd5b50610154610259366004611321565b610707565b348015610269575f5ffd5b5061015461027836600461137d565b610791565b348015610288575f5ffd5b507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546040516001600160a01b039091168152602001610222565b3480156102ce575f5ffd5b506101546102dd366004611171565b6107f4565b3480156102ed575f5ffd5b506101546102fc36600461121c565b61085c565b34801561030c575f5ffd5b50610331604051806040016040528060058152602001640352e302e360dc1b81525081565b60405161022291906113a5565b348015610349575f5ffd5b506101546103583660046113da565b6108ba565b348015610368575f5ffd5b50610154610377366004611418565b610921565b348015610387575f5ffd5b50610154610396366004611129565b610a66565b3480156103a6575f5ffd5b506101546103b5366004611498565b610aca565b3480156103c5575f5ffd5b506101546103d436600461137d565b610b1a565b3480156103e4575f5ffd5b506101546103f33660046114ca565b610b7d565b348015610403575f5ffd5b50610154610412366004611594565b610c10565b5f546001600160a01b0316331461044157604051630a31c3dd60e41b815260040160405180910390fd5b6001546040516303b4859560e11b81526001600160a01b0384811660048301528381166024830152909116906307690b2a906044015b5f604051808303815f87803b15801561048e575f5ffd5b505af11580156104a0573d5f5f3e3d5ffd5b505050505050565b5f546001600160a01b031633146104d257604051630a31c3dd60e41b815260040160405180910390fd5b6040518181527f0151256d62457b809bbc891b1f81c6dd0b9987552c70ce915b519750cd434dd1906020015b60405180910390a150565b331561052857604051630b9a4d6d60e31b815260040160405180910390fd5b5f54604051631e702f8360e01b815260048101849052602481018390526001600160a01b0390911690631e702f8390604401610477565b5f546001600160a01b0316331461058957604051630a31c3dd60e41b815260040160405180910390fd5b827f0f0ef1ab97439def0a9d2c6d9dc166207f1b13b99e62b442b2993d6153c63a6e83836040516105bb9291906115d5565b60405180910390a2505050565b5f546001600160a01b031633146105f257604051630a31c3dd60e41b815260040160405180910390fd5b6040518181527f2ccdfd47cf0c1f1069d949f1789bb79b2f12821f021634fc835af1de66ea2feb906020016104fe565b5f546001600160a01b0316331461064c57604051630a31c3dd60e41b815260040160405180910390fd5b6001546040516339e503ab60e01b81526001600160a01b0385811660048301526024820185905260448201849052909116906339e503ab906064015b5f604051808303815f87803b15801561069f575f5ffd5b505af11580156106b1573d5f5f3e3d5ffd5b50505050505050565b6106c2610c52565b6106cb82610cf6565b6106d58282610cfe565b5050565b5f6106e2610dbf565b505f5160206116f95f395f51905f5290565b6106fc610e08565b6107055f610e63565b565b331561072657604051630b9a4d6d60e31b815260040160405180910390fd5b5f54604051633b7f6a1d60e11b81526001600160a01b03909116906376fed43a9061075d90889088908890889088906004016115f0565b5f604051808303815f87803b158015610774575f5ffd5b505af1158015610786573d5f5f3e3d5ffd5b505050505050505050565b5f546001600160a01b031633146107bb57604051630a31c3dd60e41b815260040160405180910390fd5b600154604051630f37d5a760e31b81526001600160a01b03848116600483015260248201849052909116906379bead3890604401610477565b5f546001600160a01b0316331461081e57604051630a31c3dd60e41b815260040160405180910390fd5b817fb975807576e3b1461be7de07ebf7d20e4790ed802d7a0c4fdd0a1a13df72a9358260405161085091815260200190565b60405180910390a25050565b331561087b57604051630b9a4d6d60e31b815260040160405180910390fd5b5f5460405163545584dd60e11b81526001600160a01b03858116600483015260248201859052604482018490529091169063a8ab09ba90606401610688565b5f546001600160a01b031633146108e457604051630a31c3dd60e41b815260040160405180910390fd5b7f47d10eed096a44e3d0abc586c7e3a5d6cb5358cc90e7d437cd0627f7e765fb9982826040516109159291906115d5565b60405180910390a15050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff1615906001600160401b03165f811580156109655750825b90505f826001600160401b031660011480156109805750303b155b90508115801561098e575080155b156109ac5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156109d657845460ff60401b1916600160401b1785555b6109df86610ed3565b6109e7610ee4565b5f80546001600160a01b03808b166001600160a01b03199283161790925560018054928a16929091169190911790558315610a5c57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050505050565b5f546001600160a01b03163314610a9057604051630a31c3dd60e41b815260040160405180910390fd5b60015460405163d6a0c7af60e01b81526001600160a01b03848116600483015283811660248301529091169063d6a0c7af90604401610477565b3315610ae957604051630b9a4d6d60e31b815260040160405180910390fd5b5f54604051637046bf3360e11b81526001600160a01b039091169063e08d7e66906104779085908590600401611659565b5f546001600160a01b03163314610b4457604051630a31c3dd60e41b815260040160405180910390fd5b6001546040516338c110ef60e21b81526001600160a01b038481166004830152602482018490529091169063e30443bc90604401610477565b3315610b9c57604051630b9a4d6d60e31b815260040160405180910390fd5b5f54604051633af7c41360e21b81526001600160a01b039091169063ebdf104c90610bd9908b908b908b908b908b908b908b908b9060040161166c565b5f604051808303815f87803b158015610bf0575f5ffd5b505af1158015610c02573d5f5f3e3d5ffd5b505050505050505050505050565b610c18610e08565b6001600160a01b038116610c4657604051631e4fbdf760e01b81525f60048201526024015b60405180910390fd5b610c4f81610e63565b50565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480610cd857507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610ccc5f5160206116f95f395f51905f52546001600160a01b031690565b6001600160a01b031614155b156107055760405163703e46dd60e11b815260040160405180910390fd5b610c4f610e08565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610d58575060408051601f3d908101601f19168201909252610d55918101906116cb565b60015b610d8057604051634c9c8ce360e01b81526001600160a01b0383166004820152602401610c3d565b5f5160206116f95f395f51905f528114610db057604051632a87526960e21b815260048101829052602401610c3d565b610dba8383610eec565b505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107055760405163703e46dd60e11b815260040160405180910390fd5b33610e3a7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146107055760405163118cdaa760e01b8152336004820152602401610c3d565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b610edb610f41565b610c4f81610f8a565b610705610f41565b610ef582610f92565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115610f3957610dba8282610ff5565b6106d5611067565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054600160401b900460ff1661070557604051631afcd79f60e31b815260040160405180910390fd5b610c18610f41565b806001600160a01b03163b5f03610fc757604051634c9c8ce360e01b81526001600160a01b0382166004820152602401610c3d565b5f5160206116f95f395f51905f5280546001600160a01b0319166001600160a01b0392909216919091179055565b60605f5f846001600160a01b03168460405161101191906116e2565b5f60405180830381855af49150503d805f8114611049576040519150601f19603f3d011682016040523d82523d5f602084013e61104e565b606091505b509150915061105e858383611086565b95945050505050565b34156107055760405163b398979f60e01b815260040160405180910390fd5b60608261109b57611096826110e5565b6110de565b81511580156110b257506001600160a01b0384163b155b156110db57604051639996b31560e01b81526001600160a01b0385166004820152602401610c3d565b50805b9392505050565b8051156110f55780518082602001fd5b60405163d6bda27560e01b815260040160405180910390fd5b80356001600160a01b0381168114611124575f5ffd5b919050565b5f5f6040838503121561113a575f5ffd5b6111438361110e565b91506111516020840161110e565b90509250929050565b5f6020828403121561116a575f5ffd5b5035919050565b5f5f60408385031215611182575f5ffd5b50508035926020909101359150565b5f5f83601f8401126111a1575f5ffd5b5081356001600160401b038111156111b7575f5ffd5b6020830191508360208285010111156111ce575f5ffd5b9250929050565b5f5f5f604084860312156111e7575f5ffd5b8335925060208401356001600160401b03811115611203575f5ffd5b61120f86828701611191565b9497909650939450505050565b5f5f5f6060848603121561122e575f5ffd5b6112378461110e565b95602085013595506040909401359392505050565b634e487b7160e01b5f52604160045260245ffd5b5f5f60408385031215611271575f5ffd5b61127a8361110e565b915060208301356001600160401b03811115611294575f5ffd5b8301601f810185136112a4575f5ffd5b80356001600160401b038111156112bd576112bd61124c565b604051601f8201601f19908116603f011681016001600160401b03811182821017156112eb576112eb61124c565b604052818152828201602001871015611302575f5ffd5b816020840160208301375f602083830101528093505050509250929050565b5f5f5f5f5f60808688031215611335575f5ffd5b61133e8661110e565b94506020860135935060408601356001600160401b0381111561135f575f5ffd5b61136b88828901611191565b96999598509660600135949350505050565b5f5f6040838503121561138e575f5ffd5b6113978361110e565b946020939093013593505050565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b5f5f602083850312156113eb575f5ffd5b82356001600160401b03811115611400575f5ffd5b61140c85828601611191565b90969095509350505050565b5f5f5f6060848603121561142a575f5ffd5b6114338461110e565b92506114416020850161110e565b915061144f6040850161110e565b90509250925092565b5f5f83601f840112611468575f5ffd5b5081356001600160401b0381111561147e575f5ffd5b6020830191508360208260051b85010111156111ce575f5ffd5b5f5f602083850312156114a9575f5ffd5b82356001600160401b038111156114be575f5ffd5b61140c85828601611458565b5f5f5f5f5f5f5f5f6080898b0312156114e1575f5ffd5b88356001600160401b038111156114f6575f5ffd5b6115028b828c01611458565b90995097505060208901356001600160401b03811115611520575f5ffd5b61152c8b828c01611458565b90975095505060408901356001600160401b0381111561154a575f5ffd5b6115568b828c01611458565b90955093505060608901356001600160401b03811115611574575f5ffd5b6115808b828c01611458565b999c989b5096995094979396929594505050565b5f602082840312156115a4575f5ffd5b6110de8261110e565b81835281816020850137505f828201602090810191909152601f909101601f19169091010190565b602081525f6115e86020830184866115ad565b949350505050565b60018060a01b0386168152846020820152608060408201525f6116176080830185876115ad565b90508260608301529695505050505050565b8183525f6001600160fb1b03831115611640575f5ffd5b8260051b80836020870137939093016020019392505050565b602081525f6115e8602083018486611629565b608081525f61167f608083018a8c611629565b828103602084015261169281898b611629565b905082810360408401526116a7818789611629565b905082810360608401526116bc818587611629565b9b9a5050505050505050505050565b5f602082840312156116db575f5ffd5b5051919050565b5f82518060208501845e5f92019182525091905056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220b9797b4b462efaed14ce37067eea63553ca12b3b54e3144a585a0f02847285d764736f6c634300081b0033")
}

// ContractAddress is the NodeDriver contract address
var ContractAddress = common.HexToAddress("0xd100a01e00000000000000000000000000000000")
