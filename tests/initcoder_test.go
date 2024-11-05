package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/initcode"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

// this constant comes from  https://eips.ethereum.org/EIPS/eip-3860#parameters
const MAX_INIT_CODE_SIZE uint64 = 49152
const sufficientGas = uint64(100_000)

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
		testForVariant(t, net, contract, receipt, contract.CreateContractWith, gasForCreate, wordCostCreate)
	})

	// -- using CREATE2 instruction
	// create2 has extra costs for hashing, this is explained in
	// https://eips.ethereum.org/EIPS/eip-3860 and reflected in evm.codes calculator.
	const wordCostCreate2 uint64 = wordCostCreate + 6
	var gasForCreate2 uint64 = gasForCreate + 44
	t.Run("create2", func(t *testing.T) {
		testForVariant(t, net, contract, receipt, contract.Create2ContractWith, gasForCreate2, wordCostCreate2)
	})

	t.Run("make a transaction that charges for the init code size", func(t *testing.T) {
		// transactions charge 4 gas for each zero byte in data.
		const zeroByteCost uint64 = 4
		// createa a transaction with 1 word of code.
		receipt1, err := runTransactionWithCodeSizeAndGas(require, net, 1, sufficientGas)
		require.NoError(err)
		require.Equal(types.ReceiptStatusSuccessful, receipt1.Status, "failed on transfer to empty receiver with valid code")

		// createa a transaction with 2 byte of code.
		receipt2, err := runTransactionWithCodeSizeAndGas(require, net, 2, sufficientGas)
		require.NoError(err)
		require.Equal(types.ReceiptStatusSuccessful, receipt2.Status, "failed on transfer to empty receiver with valid code")

		difference := receipt2.GasUsed - receipt1.GasUsed
		require.Equal(difference, zeroByteCost, "gas difference between 1 and 2 words should be 4, instead got %d",
			difference)
	})

	t.Run("make a transaction without enough gas for init code size", func(t *testing.T) {
		// as specified in https://eips.ethereum.org/EIPS/eip-3860#rules,
		// this is similar to transactions considered invalid for not meeting the intrinsic gas cost requirement.
		_, err := runTransactionWithCodeSizeAndGas(require, net, MAX_INIT_CODE_SIZE+1, sufficientGas)
		require.ErrorContains(err, "intrinsic gas too low", "unexpectedly succeeded to create contract with init code length greater than MAX_INITCODE_SIZE")
	})
}

func testForVariant(t *testing.T, net *IntegrationTestNet,
	contract *initcode.Initcode, receipt *types.Receipt, variant variant,
	gasForContract, wordCost uint64) {
	require := require.New(t)

	var gasCostFor2Words uint64 = wordCost * 2
	// we use enough gas for all tests to afford cost of 3 words as well.
	var gasFor3Words = gasForContract + gasCostFor2Words
	t.Run("create transaction with enough gas for init code", func(t *testing.T) {
		// we need to provide same gas to all transactions so that we can compare the cost.
		cost1Word, _ := createAndCheckContractWithCodeLenAndGas(require, net, contract, variant, 30, gasFor3Words)
		cost2Words, _ := createAndCheckContractWithCodeLenAndGas(require, net, contract, variant, 42, gasFor3Words)
		cost3Words, _ := createAndCheckContractWithCodeLenAndGas(require, net, contract, variant, 90, gasFor3Words)
		require.Equal(cost2Words-cost1Word, wordCost,
			"cost difference between 1 and 2 words should be %v, instead got %d", wordCost, cost2Words-cost1Word)
		require.Equal(cost3Words-cost1Word, gasCostFor2Words,
			"cost difference between 1 and 3 words should be %v, instead got %d", gasCostFor2Words, cost3Words-cost1Word)
	})

	t.Run("create transaction without enough gas for init code", func(t *testing.T) {
		receipt = createContractWithCodeLenAndGas(require, net, variant, 30, gasForContract-wordCost)
		require.Equal(types.ReceiptStatusFailed, receipt.Status,
			"unexpectedly succeeded to create contract without enough gas")
	})

	t.Run("create transaction with max init code size", func(t *testing.T) {
		_, _ = createAndCheckContractWithCodeLenAndGas(require, net, contract, variant, MAX_INIT_CODE_SIZE, sufficientGas)
	})

	t.Run("create transaction with MAX_INITCODE_SIZE+1", func(t *testing.T) {
		receipt = createContractWithCodeLenAndGas(require, net, variant, MAX_INIT_CODE_SIZE+1, sufficientGas)
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
	require.NoError(err)
	return receipt
}

type variant func(opts *bind.TransactOpts, codeSize *big.Int) (*types.Transaction, error)

func runTransactionWithCodeSizeAndGas(r *require.Assertions, net *IntegrationTestNet, codeSize, gas uint64) (*types.Receipt, error) {
	// these values are needed for the transaction but are irrelevant for the test
	client, err := net.GetClient()
	r.NoError(err, "failed to connect to the network:")

	// defer client.Close()
	chainId, err := client.ChainID(context.Background())
	r.NoError(err, "failed to get chain ID::")

	nonce, err := client.NonceAt(context.Background(), net.validator.Address(), nil)
	r.NoError(err, "failed to get nonce:")

	price, err := client.SuggestGasPrice(context.Background())
	r.NoError(err, "failed to get gas price:")
	client.Close()
	// ---------

	transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
		ChainID:  chainId,
		Gas:      gas,
		GasPrice: price,
		To:       nil,
		Nonce:    nonce,
		Data:     make([]byte, codeSize),
	}), types.NewLondonSigner(chainId), net.validator.PrivateKey)
	r.NoError(err, "failed to sign transaction:")
	return net.Run(transaction)
}
