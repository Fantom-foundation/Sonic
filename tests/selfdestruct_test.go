package tests

import (
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/query_account"
	"github.com/Fantom-foundation/go-opera/tests/contracts/selfdestruct"
	"github.com/Fantom-foundation/go-opera/tests/contracts/selfdestruct_factory"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestSelfDestruct(t *testing.T) {
	require := require.New(t)

	someBalance := big.NewInt(1234)
	addressCount := 1234 // to avoid beneficiary address collisions

	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(err, "failed to start test network")
	defer net.Stop()

	// Utility contract to query a second account state
	queryAccount, receipt, err := DeployContract(net, query_account.DeployQueryAccount)
	require.NoError(err)
	require.Equal(receipt.Status, types.ReceiptStatusSuccessful, "failed to deploy contract")

	t.Run("constructor", func(t *testing.T) {
		tests := map[string]struct {
			deployTx  deployTxFunction[selfdestruct.Selfdestruct]
			executeTx executeTxFunction[selfdestruct.Selfdestruct]
			effects   map[string]effectFunction
		}{
			// This test checks the testing infrastructure itself
			"sanity check/no selfdestruct": {
				deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.Selfdestruct, error) {
					opts.Value = someBalance
					return selfdestruct.DeploySelfdestruct(opts, backend,
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
				deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.Selfdestruct, error) {
					opts.Value = someBalance
					return selfdestruct.DeploySelfdestruct(opts, backend,
						false, // do not selfdestruct in constructor
						false, // beneficiary is not self
						otherAddress)
				},
				executeTx: func(contract *selfdestruct.Selfdestruct, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
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
				deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.Selfdestruct, error) {
					opts.Value = someBalance
					return selfdestruct.DeploySelfdestruct(opts, backend,
						false, // do not selfdestruct in constructor
						true,  // beneficiary is self
						otherAddress)
				},
				executeTx: func(contract *selfdestruct.Selfdestruct, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
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
				deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.Selfdestruct, error) {
					opts.Value = someBalance
					return selfdestruct.DeploySelfdestruct(opts, backend,
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
				deployTx: func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *selfdestruct.Selfdestruct, error) {
					opts.Value = someBalance
					return selfdestruct.DeploySelfdestruct(opts, backend,
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
					func(to *bind.TransactOpts, cb bind.ContractBackend) (common.Address, *types.Transaction, *selfdestruct.Selfdestruct, error) {
						return test.deployTx(to, cb, otherAddress)
					})
				require.NoError(err)
				require.Equal(deployReceipt.Status, types.ReceiptStatusSuccessful, "failed to deploy contract")
				allLogs := deployReceipt.Logs

				var executionReceipt *types.Receipt
				// Second transaction executes some contract function (if any)
				if test.executeTx != nil {
					executionReceipt, err = net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
						return test.executeTx(contract, opts, otherAddress)
					})
					require.NoError(err)
					require.Equal(executionReceipt.Status, types.ReceiptStatusSuccessful, "failed to execute contract")
					allLogs = append(allLogs, deployReceipt.Logs...)
				}

				effectContext := effectContext{
					qa:                 queryAccount,
					contract:           contract,
					executionReceipt:   executionReceipt,
					allLogs:            allLogs,
					contractAddress:    deployReceipt.ContractAddress,
					beneficiaryAddress: otherAddress,
				}

				// check effects
				for name, effect := range test.effects {
					t.Run(name, func(t *testing.T) {
						effect(require, &effectContext)
					})
				}
			})
		}
	})

	t.Run("nested call", func(t *testing.T) {
		tests := map[string]struct {
			transactions []executeTxFunction[selfdestruct_factory.SelfdestructFactory]
			effects      map[string]effectFunction
		}{
			// This test checks the testing infrastructure itself
			"sanity check/no selfdestruct": {
				transactions: []executeTxFunction[selfdestruct_factory.SelfdestructFactory]{
					func(contract *selfdestruct_factory.SelfdestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
						opts.Value = someBalance
						return contract.Create(opts)
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
				transactions: []executeTxFunction[selfdestruct_factory.SelfdestructFactory]{
					func(contract *selfdestruct_factory.SelfdestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
						opts.Value = someBalance
						return contract.Create(opts)
					},
					func(contract *selfdestruct_factory.SelfdestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
						return contract.Destroy(opts, otherAddress)
					},
				},
				effects: map[string]effectFunction{
					"contract looses balance":   contractBalanceIs(big.NewInt(0)),
					"beneficiary gains balance": beneficiaryBalanceIs(someBalance),
					"storage is not deleted":    contractStorageIs(big.NewInt(123)),
					"code is not deleted":       contractCodeSizeIsNot(big.NewInt(0)),
				},
			},
			"different tx/beneficiary is same account": {
				transactions: []executeTxFunction[selfdestruct_factory.SelfdestructFactory]{
					func(contract *selfdestruct_factory.SelfdestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
						opts.Value = someBalance
						return contract.Create(opts)
					},
					func(contract *selfdestruct_factory.SelfdestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
						return contract.DestroyWithoutBeneficiary(opts)
					},
				},
				effects: map[string]effectFunction{
					"storage is not deleted": contractStorageIs(big.NewInt(123)),
					"code is not deleted":    contractCodeSizeIsNot(big.NewInt(0)),
					"balance is not burned":  contractBalanceIs(someBalance),
				},
			},
			"same tx/beneficiary is other account": {
				transactions: []executeTxFunction[selfdestruct_factory.SelfdestructFactory]{
					func(contract *selfdestruct_factory.SelfdestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
						opts.Value = someBalance
						return contract.CreateAndDestroy(opts, otherAddress)
					},
				},
				effects: map[string]effectFunction{
					"code is deleted":                            contractCodeSizeIs(big.NewInt(0)),
					"contract looses balance":                    contractBalanceIs(big.NewInt(0)),
					"beneficiary gains balance":                  beneficiaryBalanceIs(someBalance),
					"storage is not deleted until the end of tx": nestedContractValueAfterSelfdestructIs(big.NewInt(123)),
				},
			},
			"same tx/beneficiary is same account": {
				transactions: []executeTxFunction[selfdestruct_factory.SelfdestructFactory]{
					func(contract *selfdestruct_factory.SelfdestructFactory, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error) {
						opts.Value = someBalance
						return contract.CreateAndDestroyWithoutBeneficiary(opts)
					},
				},
				effects: map[string]effectFunction{
					"code is deleted":                            contractCodeSizeIs(big.NewInt(0)),
					"contract looses balance":                    contractBalanceIs(big.NewInt(0)),
					"storage is not deleted until the end of tx": nestedContractValueAfterSelfdestructIs(big.NewInt(123)),
				},
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {

				client, err := net.GetClient()
				require.NoError(err)
				defer client.Close()

				otherAddress := common.BigToAddress(big.NewInt(int64(addressCount)))
				addressCount++

				factory, receipt, err := DeployContract(net, selfdestruct_factory.DeploySelfdestructFactory)
				require.NoError(err)
				require.Equal(receipt.Status, types.ReceiptStatusSuccessful, "failed to deploy contract")

				allLogs := []*types.Log{}
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
				require.NotEqual(contractAddress, common.Address{}, "factory failed to construct a contract")
				newContract, err := selfdestruct.NewSelfdestruct(contractAddress, client)
				require.NoError(err, "failed to instantiate contract")

				effectContext := effectContext{
					qa:                 queryAccount,
					cf:                 factory,
					contract:           newContract,
					allLogs:            allLogs,
					contractAddress:    contractAddress,
					beneficiaryAddress: otherAddress,
				}

				// check effects
				for name, effect := range test.effects {
					t.Run(name, func(t *testing.T) {
						effect(require, &effectContext)
					})
				}
			})
		}
	})
}

type effectContext struct {
	qa                 *query_account.QueryAccount
	cf                 *selfdestruct_factory.SelfdestructFactory
	contract           *selfdestruct.Selfdestruct
	executionReceipt   *types.Receipt
	allLogs            []*types.Log
	contractAddress    common.Address
	beneficiaryAddress common.Address
}

type effectFunction func(require *require.Assertions, ctx *effectContext)
type deployTxFunction[T any] func(opts *bind.TransactOpts, backend bind.ContractBackend, otherAddress common.Address) (common.Address, *types.Transaction, *T, error)
type executeTxFunction[T any] func(contract *T, opts *bind.TransactOpts, otherAddress common.Address) (*types.Transaction, error)

func executionHalted() effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		for _, log := range ctx.allLogs {
			_, err := ctx.contract.ParseLogAfterDestruct(*log)
			require.Error(err, "execution should have halted, log after selfdestruct should not be present")
		}
	}
}

func nestedContractValueAfterSelfdestructIs(value *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		for _, log := range ctx.allLogs {
			storage, err := ctx.cf.ParseLogContractStorage(*log)
			if err != nil {
				continue
			}
			require.Conditionf(Equal(storage.Value, value), "storage value differs, got %v: want %v", storage, value)
			return
		}
		require.Fail("no log with storage value found")
	}
}

func contractBalanceIs(expected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		gas, err := ctx.qa.GetBalance(nil, ctx.contractAddress)
		require.NoError(err)
		require.Conditionf(Equal(gas, expected), "balance not expected, got %v: want %v", gas, expected)
	}
}

func contractStorageIs(expected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		storage, err := ctx.contract.SomeData(nil)
		require.NoError(err)
		require.Conditionf(Equal(storage, expected), "storage value differs, got %v: want %v", storage, expected)
	}
}

func beneficiaryBalanceIs(expected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		gas, err := ctx.qa.GetBalance(nil, ctx.beneficiaryAddress)
		require.NoError(err)
		require.Conditionf(Equal(gas, expected), "balance not expected, got %v: want %v", gas, expected)
	}
}

func contractCodeSizeIs(expected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		codeSize, err := ctx.qa.GetCodeSize(nil, ctx.contractAddress)
		require.NoError(err)
		require.Conditionf(Equal(codeSize, expected), "code size not expected, got %v: want %v", codeSize, expected)
	}
}

func contractCodeSizeIsNot(notExpected *big.Int) effectFunction {
	return func(require *require.Assertions, ctx *effectContext) {
		codeSize, err := ctx.qa.GetCodeSize(nil, ctx.contractAddress)
		require.NoError(err)
		require.Conditionf(Not(Equal(codeSize, notExpected)), "code size not expected, got %v: wanted distinct from %v", codeSize, notExpected)
	}
}

type ImplementsCmp[v any] interface {
	Cmp(v) int
}

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
