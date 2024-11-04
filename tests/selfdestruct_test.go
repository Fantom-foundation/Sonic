package tests

import (
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/query_account"
	"github.com/Fantom-foundation/go-opera/tests/contracts/selfdestruct"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

var (
	// addressCount is used to generate unique addresses, avoiding clashes of beneficiary addresses
	addressCount = 1234
)

func TestSelfDestruct(t *testing.T) {
	require := require.New(t)

	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(err, "failed to start test network")
	defer net.Stop()

	// Utility contract to query a second account state
	queryAccount, receipt, err := DeployContract(net, query_account.DeployQueryAccount)
	require.NoError(err)
	require.Equal(receipt.Status, types.ReceiptStatusSuccessful, "failed to deploy contract")

	t.Run("constructor", func(t *testing.T) {
		testSelfDestruct_Constructor(t, net, queryAccount)
	})

	t.Run("nested call", func(t *testing.T) {
		testSelfDestruct_NestedCall(t, big.NewInt(1234), net, queryAccount)
	})
}

func testSelfDestruct_Constructor(t *testing.T, net *IntegrationTestNet, queryAccount *query_account.QueryAccount) {
	require := require.New(t)
	someBalance := big.NewInt(1234)

	tests := map[string]struct {
		deployTx  deployTxFunction[selfdestruct.SelfDestruct]
		executeTx executeTxFunction[selfdestruct.SelfDestruct]
		effects   map[string]effectFunction
	}{
		// This test checks the testing infrastructure itself
		"sanity check/no selfdestruct": {
			deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.SelfDestruct, error) {
				opts.Value = someBalance
				return selfdestruct.DeploySelfDestruct(opts, backend,
					false, // do not selfdestruct in constructor
					false, // beneficiary is not self
					otherAddress)
			},
			effects: map[string]effectFunction{
				"contract keeps balance":       contractBalanceIs(someBalance),
				"beneficiary gains no balance": beneficiaryBalanceIs(big.NewInt(0)),
				"storage is not deleted":       contractStorageIs(big.NewInt(123)),
				"code is not deleted":          contractCodeSizeIsNot(big.NewInt(0)),
			},
		},
		"different tx/beneficiary is other account": {
			deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.SelfDestruct, error) {
				opts.Value = someBalance
				return selfdestruct.DeploySelfDestruct(opts, backend,
					false, // do not selfdestruct in constructor
					false, // beneficiary is not self
					otherAddress)
			},
			executeTx: func(contract *selfdestruct.SelfDestruct, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
				return contract.DestroyContract(opts,
					false, // beneficiary is not self
					otherAddress)
			},
			effects: map[string]effectFunction{
				"the current execution frame halts": executionHalted(),
				"contract looses balance":           contractBalanceIs(big.NewInt(0)),
				"beneficiary gains balance":         beneficiaryBalanceIs(someBalance),
				"storage is not deleted":            contractStorageIs(big.NewInt(123)),
				"code is not deleted":               contractCodeSizeIsNot(big.NewInt(0)),
			},
		},
		"different tx/beneficiary is same account": {
			deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.SelfDestruct, error) {
				opts.Value = someBalance
				return selfdestruct.DeploySelfDestruct(opts, backend,
					false, // do not selfdestruct in constructor
					true,  // beneficiary is self
					otherAddress)
			},
			executeTx: func(contract *selfdestruct.SelfDestruct, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
				return contract.DestroyContract(opts, true, otherAddress)
			},
			effects: map[string]effectFunction{
				"the current execution frame halts": executionHalted(),
				"storage is not deleted":            contractStorageIs(big.NewInt(123)),
				"code is not deleted":               contractCodeSizeIsNot(big.NewInt(0)),
				"balance is not burned":             contractBalanceIs(someBalance),
			},
		},
		"same tx/beneficiary is other account": {
			deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.SelfDestruct, error) {
				opts.Value = someBalance
				return selfdestruct.DeploySelfDestruct(opts, backend,
					true,  // selfdestruct in constructor
					false, // beneficiary is not self
					otherAddress)
			},
			effects: map[string]effectFunction{
				"the current execution frame halts": executionHalted(),
				"code is deleted":                   contractCodeSizeIs(big.NewInt(0)),
				"contract looses balance":           contractBalanceIs(big.NewInt(0)),
				"beneficiary gains balance":         beneficiaryBalanceIs(someBalance),
			},
		},
		"same tx/beneficiary is same account": {
			deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.SelfDestruct, error) {
				opts.Value = someBalance
				return selfdestruct.DeploySelfDestruct(opts, backend,
					true, // selfdestruct in constructor
					true, // beneficiary is  self
					otherAddress)
			},
			effects: map[string]effectFunction{
				"the current execution frame halts": executionHalted(),
				"code is deleted":                   contractCodeSizeIs(big.NewInt(0)),
				"contract looses balance":           contractBalanceIs(big.NewInt(0)),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			otherAddress := common.BigToAddress(big.NewInt(int64(addressCount)))
			addressCount++

			// First transaction deploys contract
			contract, deployReceipt, err := DeployContract(net,
				func(to *bind.TransactOpts, cb bind.ContractBackend) (common.Address, *types.Transaction, *selfdestruct.SelfDestruct, error) {
					return test.deployTx(to, cb, otherAddress)
				})
			require.NoError(err)
			require.Equal(deployReceipt.Status,
				types.ReceiptStatusSuccessful,
				"failed to deploy contract",
			)
			allLogs := deployReceipt.Logs

			var executionReceipt *types.Receipt
			// Second transaction executes some contract function (if any)
			if test.executeTx != nil {
				executionReceipt, err = net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
					return test.executeTx(contract, opts, otherAddress)
				})
				require.NoError(err)
				require.Equal(executionReceipt.Status,
					types.ReceiptStatusSuccessful,
					"failed to execute contract",
				)
				allLogs = append(allLogs, deployReceipt.Logs...)
			}

			// check effects
			effectContext := effectContext{
				queryAccount:       queryAccount,
				contract:           contract,
				executionReceipt:   executionReceipt,
				allLogs:            allLogs,
				contractAddress:    deployReceipt.ContractAddress,
				beneficiaryAddress: otherAddress,
			}
			for name, effect := range test.effects {
				t.Run(name, func(t *testing.T) {
					effect(require, &effectContext)
				})
			}
		})
	}
}

func testSelfDestruct_NestedCall(t *testing.T, someBalance *big.Int, net *IntegrationTestNet, queryAccount *query_account.QueryAccount) {
	require := require.New(t)

	tests := map[string]struct {
		transactions []executeTxFunction[selfdestruct.SelfDestructFactory]
		effects      map[string]effectFunction
	}{
		// This test checks the testing infrastructure itself
		"sanity check/no selfdestruct": {
			transactions: []executeTxFunction[selfdestruct.SelfDestructFactory]{
				func(factory *selfdestruct.SelfDestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
					opts.Value = someBalance
					return factory.Create(opts)
				},
			},
			effects: map[string]effectFunction{
				"contract has storage":         contractStorageIs(big.NewInt(123)),
				"contract has code":            contractCodeSizeIsNot(big.NewInt(0)),
				"contract keeps balance":       contractBalanceIs(someBalance),
				"beneficiary gains no balance": beneficiaryBalanceIs(big.NewInt(0)),
			},
		},
		"different tx/beneficiary is other account": {
			transactions: []executeTxFunction[selfdestruct.SelfDestructFactory]{
				func(factory *selfdestruct.SelfDestructFactory, opts *bind.TransactOpts, _ common.Address) (*types.Transaction, error) {
					opts.Value = someBalance
					return factory.Create(opts)
				},
				func(factory *selfdestruct.SelfDestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
					return factory.Destroy(opts, otherAddress)
				},
			},
			effects: map[string]effectFunction{
				"storage is not deleted":       contractStorageIs(big.NewInt(123)),
				"code is not deleted":          contractCodeSizeIsNot(big.NewInt(0)),
				"contract looses balance":      contractBalanceIs(big.NewInt(0)),
				"beneficiary gains balance":    beneficiaryBalanceIs(someBalance),
				"nested execution frame halts": executionHalted(),
			},
		},
		"different tx/beneficiary is same account": {
			transactions: []executeTxFunction[selfdestruct.SelfDestructFactory]{
				func(factory *selfdestruct.SelfDestructFactory, opts *bind.TransactOpts, _ common.Address) (*types.Transaction, error) {
					opts.Value = someBalance
					return factory.Create(opts)
				},
				func(factory *selfdestruct.SelfDestructFactory, opts *bind.TransactOpts, _ common.Address) (*types.Transaction, error) {
					return factory.DestroyWithoutBeneficiary(opts)
				},
			},
			effects: map[string]effectFunction{
				"storage is not deleted":       contractStorageIs(big.NewInt(123)),
				"code is not deleted":          contractCodeSizeIsNot(big.NewInt(0)),
				"balance is not burned":        contractBalanceIs(someBalance),
				"nested execution frame halts": executionHalted(),
			},
		},
		"same tx/beneficiary is other account": {
			transactions: []executeTxFunction[selfdestruct.SelfDestructFactory]{
				func(factory *selfdestruct.SelfDestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
					opts.Value = someBalance
					return factory.CreateAndDestroy(opts, otherAddress)
				},
			},
			effects: map[string]effectFunction{
				"code is deleted":                contractCodeSizeIs(big.NewInt(0)),
				"contract looses balance":        contractBalanceIs(big.NewInt(0)),
				"beneficiary gains balance":      beneficiaryBalanceIs(someBalance),
				"storage exists until end of tx": nestedContractValueAfterSelfDestructIs(big.NewInt(123)),
				"nested execution frame halts":   executionHalted(),
			},
		},
		"same tx/beneficiary is same account": {
			transactions: []executeTxFunction[selfdestruct.SelfDestructFactory]{
				func(factory *selfdestruct.SelfDestructFactory, opts *bind.TransactOpts, _ common.Address) (*types.Transaction, error) {
					opts.Value = someBalance
					return factory.CreateAndDestroyWithoutBeneficiary(opts)
				},
			},
			effects: map[string]effectFunction{
				"code is deleted":                contractCodeSizeIs(big.NewInt(0)),
				"contract looses balance":        contractBalanceIs(big.NewInt(0)),
				"storage exists until end of tx": nestedContractValueAfterSelfDestructIs(big.NewInt(123)),
				"nested execution frame halts":   executionHalted(),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			// generate a new beneficiary address for each test
			otherAddress := common.BigToAddress(big.NewInt(int64(addressCount)))
			addressCount++

			// deploy factory contract
			factory, receipt, err := DeployContract(net, selfdestruct.DeploySelfDestructFactory)
			require.NoError(err)
			require.Equal(receipt.Status,
				types.ReceiptStatusSuccessful,
				"failed to deploy contract",
			)

			allLogs := []*types.Log{}

			// execute all described transactions
			for _, tx := range test.transactions {
				receipt, err := net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
					return tx(factory, opts, otherAddress)
				})
				require.NoError(err)
				allLogs = append(allLogs, receipt.Logs...)
			}
			var contractAddress common.Address
			for _, log := range allLogs {
				parse, err := factory.ParseLogDeployed(*log)
				if err == nil {
					contractAddress = parse.Addr
					break
				}
			}
			require.NotEqual(contractAddress,
				common.Address{},
				"factory failed to construct a contract",
			)

			// get contract bindings (to enable contract interaction)
			client, err := net.GetClient()
			require.NoError(err)
			defer client.Close()
			newContract, err := selfdestruct.NewSelfDestruct(contractAddress, client)
			require.NoError(err, "failed to instantiate contract")

			// check effects
			effectContext := effectContext{
				queryAccount:       queryAccount,
				factory:            factory,
				contract:           newContract,
				allLogs:            allLogs,
				contractAddress:    contractAddress,
				beneficiaryAddress: otherAddress,
			}
			for name, effect := range test.effects {
				t.Run(name, func(t *testing.T) {
					effect(require, &effectContext)
				})
			}
		})
	}
}

// effectContext stores all the information needed to check the effects of a test
// Each test will construct a contract that may selfdestruct.
// Each test will deploy a new contract and generate a new beneficiary account to
// avoid clashes between tests. This struct will be filled with the results of
// the each test setup.
type effectContext struct {
	queryAccount       *query_account.QueryAccount       //< utility contract to query account state
	contract           *selfdestruct.SelfDestruct        //< contract to test (may have selfdestructed)
	factory            *selfdestruct.SelfDestructFactory //< factory contract to deploy new contracts
	executionReceipt   *types.Receipt                    //< receipt of the execution transaction for constructor tests
	allLogs            []*types.Log                      //< all logs generated by all transactions
	contractAddress    common.Address                    //< address of the contract that may selfdestruct
	beneficiaryAddress common.Address                    //< address of the beneficiary account
}

type effectFunction func(require *require.Assertions, ctx *effectContext)
type deployTxFunction[T any] func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *T, error)
type executeTxFunction[T any] func(contract *T, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error)

// executionHalted checks that the execution stopped after selfdestruct
// This is done by looking for logs
func executionHalted() effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		for _, log := range ctx.allLogs {
			_, err := ctx.contract.ParseLogAfterDestruct(*log)
			require.Error(
				err,
				"execution should have halted, log after selfdestruct should not be present",
			)
		}
	}
}

// nestedContractValueAfterSelfDestructIs reads the internal storage of the contract
// and compares its content to the expected value
// use this matcher when the contract is being destroyed by a nested call, but
// the internal value of the contract storage will be emitted before the
// transaction is  completed.
func nestedContractValueAfterSelfDestructIs(value *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		for _, log := range ctx.allLogs {
			storage, err := ctx.factory.ParseLogContractStorage(*log)
			if err != nil {
				continue
			}
			require.Conditionf(
				Equal(storage.Value, value),
				"storage value differs, got %v: want %v",
				storage,
				value,
			)
			return
		}
		require.Fail("no log with storage value found")
	}
}

// contractBalanceIs reads the contract balance and compare it to the expected value
func contractBalanceIs(expected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		gas, err := ctx.queryAccount.GetBalance(nil, ctx.contractAddress)
		require.NoError(err)
		require.Conditionf(
			Equal(gas, expected),
			"balance not expected, got %v: want %v",
			gas,
			expected,
		)
	}
}

// contractStorageIs reads the contract storage and compare it to the expected value
func contractStorageIs(expected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		storage, err := ctx.contract.SomeData(nil)
		require.NoError(err)
		require.Conditionf(
			Equal(storage, expected),
			"storage value differs, got %v: want %v",
			storage,
			expected,
		)
	}
}

func beneficiaryBalanceIs(expected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		gas, err := ctx.queryAccount.GetBalance(nil, ctx.beneficiaryAddress)
		require.NoError(err)
		require.Conditionf(
			Equal(gas, expected),
			"balance not expected, got %v: want %v",
			gas,
			expected,
		)
	}
}

func contractCodeSizeIs(expected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		codeSize, err := ctx.queryAccount.GetCodeSize(nil, ctx.contractAddress)
		require.NoError(err)
		require.Conditionf(
			Equal(codeSize, expected),
			"code size not expected, got %v: want %v",
			codeSize,
			expected,
		)
	}
}

func contractCodeSizeIsNot(notExpected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		codeSize, err := ctx.queryAccount.GetCodeSize(nil, ctx.contractAddress)
		require.NoError(err)
		require.Conditionf(
			Not(Equal(codeSize, notExpected)),
			"code size not expected, got %v: wanted distinct from %v",
			codeSize,
			notExpected,
		)
	}
}

// ImplementsCmp is an interface for types that can be compared to facilitate testing
type ImplementsCmp[v any] interface {
	Cmp(v) int
}

// Equal returns a function that compares two values of types that implement ImplementsCmp
// this is used to compare big.Int values using require.Condition
func Equal[T ImplementsCmp[T]](a, b T) func() bool {
	return func() bool {
		return a.Cmp(b) == 0
	}
}

func Not(f func() bool) func() bool {
	return func() bool {
		return !f()
	}
}
