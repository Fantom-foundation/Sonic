package tests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/initcode"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestInitCodeSizeLimitAndMetered(t *testing.T) {
	require := require.New(t)

	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(err)
	defer net.Stop()

	// Deploy the invalid start contract.
	contract, receipt, err := DeployContract(net, initcode.DeployInitcode)
	require.NoError(err)
	require.Equal(types.ReceiptStatusSuccessful, receipt.Status, "failed to deploy contract")

	// deploy measureGasAndAssign to get cost of deploying a contract without create.
	assignCost, gasUsed := createAndCheckContractWithCodeLenAndGas(require, net, contract, contract.MeasureGasAndAssign, 0, 0_000)

	// -- using CREATE instruction
	const wordCostCreate uint64 = 2
	var gasForCreate uint64 = 32000 + gasUsed + assignCost
	t.Run("create", func(t *testing.T) {
		testForVariant(t, require, net, contract, receipt, contract.CreateContractWith, gasForCreate, wordCostCreate)
	})

	// -- using CREATE2 instruction
	// create2 has extra costs for hashing.
	var gasForCreate2 uint64 = gasForCreate + 44
	const wordCostCreate2 uint64 = wordCostCreate + 6
	t.Run("create2", func(t *testing.T) {
		testForVariant(t, require, net, contract, receipt, contract.Create2ContractWith, gasForCreate2, wordCostCreate2)
	})

}

func testForVariant(t *testing.T, require *require.Assertions, net *IntegrationTestNet,
	contract *initcode.Initcode, receipt *types.Receipt, variant variant,
	gasForContract, wordCost uint64) {

	var gasDifferenceFor2Words uint64 = wordCost * 2
	// we use enough gas for all tests to afford cost of 3 words as well.
	var gasForAllTests = gasForContract + gasDifferenceFor2Words
	t.Run("create transaction with enough gas for init code", func(t *testing.T) {
		// we need to provide same gas to all transactions so that we can compare the cost.
		cost1Word, _ := createAndCheckContractWithCodeLenAndGas(require, net, contract, variant, 30, gasForAllTests)
		cost2Words, _ := createAndCheckContractWithCodeLenAndGas(require, net, contract, variant, 42, gasForAllTests)
		cost3Words, _ := createAndCheckContractWithCodeLenAndGas(require, net, contract, variant, 90, gasForAllTests)
		require.Equal(cost2Words-cost1Word, wordCost,
			"cost difference between 1 and 2 words should be %v, instead got %d", wordCost, cost2Words-cost1Word)
		require.Equal(cost3Words-cost1Word, gasDifferenceFor2Words,
			"cost difference between 1 and 3 words should be %v, instead got %d", gasDifferenceFor2Words, cost3Words-cost1Word)
	})

	t.Run("create transaction without enough gas for init code", func(t *testing.T) {
		receipt = createContractWithCodeLenAndGas(require, net, variant, 30, gasForContract-wordCost)
		require.Equal(types.ReceiptStatusFailed, receipt.Status,
			"unexpectedly succeeded to create contract without enough gas")
	})

	// these two constants come from  https://eips.ethereum.org/EIPS/eip-3860#parameters
	const MAX_INITCODE_SIZE uint64 = 49152
	var MAX_INITCODE_COST uint64 = (MAX_INITCODE_SIZE / 32) * wordCost

	var someCost = gasForContract + MAX_INITCODE_COST + 9_199
	t.Run("create transaction with max init code size", func(t *testing.T) {
		_, _ = createAndCheckContractWithCodeLenAndGas(require, net, contract, variant, MAX_INITCODE_SIZE, someCost)
	})

	t.Run("create transaction with MAX_INITCODE_SIZE+1", func(t *testing.T) {
		receipt = createContractWithCodeLenAndGas(require, net, variant, MAX_INITCODE_SIZE+1, someCost+2)
		require.Equal(types.ReceiptStatusFailed, receipt.Status,
			"unexpectedly succeeded to create contract with init code length greater than MAX_INITCODE_SIZE")
	})
}

func createAndCheckContractWithCodeLenAndGas(require *require.Assertions, net *IntegrationTestNet, contract *initcode.Initcode, variant variant, codeLen, gasLimit uint64) (uint64, uint64) {
	receipt := createContractWithCodeLenAndGas(require, net, variant, codeLen, gasLimit)
	require.Equal(types.ReceiptStatusSuccessful, receipt.Status, "failed to create contract with code length %d", codeLen)

	log, err := contract.ParseLogCost(*receipt.Logs[0])
	require.NoError(err)
	return log.Cost.Uint64(), receipt.GasUsed
}

func createContractWithCodeLenAndGas(require *require.Assertions, net *IntegrationTestNet, variant variant, codeLen, gasLimit uint64) *types.Receipt {
	receipt, err := net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		opts.GasLimit = gasLimit
		return variant(opts, big.NewInt(int64(codeLen))) //
	})
	fmt.Printf("gas provided: %v, and gas used: %v\n", gasLimit, receipt.GasUsed) // THIS IS A DEBUG LINE

	require.NoError(err)
	return receipt
}

type variant func(opts *bind.TransactOpts, codeSize *big.Int) (*types.Transaction, error)
