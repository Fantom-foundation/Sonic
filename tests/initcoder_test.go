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
	_, gasUsed := createAndCheckContractWithCodeLenAndGas(t, net, contract, contract.MeasureAssignGasCost, 0, 0)

	// -- using CREATE instruction
	const wordCostCreate uint64 = 2
	var gasForCreate uint64 = 32000 + gasUsed
	t.Run("create", func(t *testing.T) {
		testForVariant(t, net, contract, contract.CreateContractWith, gasForCreate, wordCostCreate)
	})

	// -- using CREATE2 instruction
	// create2 has extra costs for hashing, this is explained in
	// https://eips.ethereum.org/EIPS/eip-3860 and reflected in evm.codes calculator.
	const wordCostCreate2 uint64 = wordCostCreate + 6
	var gasForCreate2 uint64 = gasForCreate + 44
	t.Run("create2", func(t *testing.T) {
		testForVariant(t, net, contract, contract.Create2ContractWith, gasForCreate2, wordCostCreate2)
	})

	t.Run("create transaction", func(t *testing.T) {
		t.Run("charges depending on the init code size", func(t *testing.T) {
			// transactions charge 4 gas for each zero byte in data.
			const zeroByteCost uint64 = 4
			// create a transaction with 1 word of code.
			receipt1, err := runTransactionWithCodeSizeAndGas(t, net, 1, sufficientGas)
			require.NoError(err)
			require.Equal(types.ReceiptStatusSuccessful, receipt1.Status,
				"failed on transfer to empty receiver with valid code")

			// createa a transaction with 2 byte of code.
			receipt2, err := runTransactionWithCodeSizeAndGas(t, net, 2, sufficientGas)
			require.NoError(err)
			require.Equal(types.ReceiptStatusSuccessful, receipt2.Status,
				"failed on transfer to empty receiver with valid code")

			difference := receipt2.GasUsed - receipt1.GasUsed
			require.Equal(difference, zeroByteCost,
				"gas difference between 1 and 2 words should be 4, instead got", difference)
		})

		t.Run("aborts with init code size larger than MAX_INITCODE_SIZE", func(t *testing.T) {
			// as specified in https://eips.ethereum.org/EIPS/eip-3860#rules,
			// this is similar to transactions considered invalid for not meeting the intrinsic gas cost requirement.
			_, err := runTransactionWithCodeSizeAndGas(t, net, MAX_INIT_CODE_SIZE+1, sufficientGas)
			require.ErrorContains(err,
				"intrinsic gas too low", "unexpectedly succeeded to create contract with init code length greater than MAX_INITCODE_SIZE")
		})
	})
}

func testForVariant(t *testing.T, net *IntegrationTestNet,
	contract *initcode.Initcode, variant variant,
	gasForContract, wordCost uint64) {

	var gasCostFor2Words uint64 = wordCost * 3
	// we use enough gas for all tests to afford cost of 3 words as well.
	var gasFor3Words = gasForContract + gasCostFor2Words
	t.Run("charges depending on the init code size", func(t *testing.T) {
		// we need to provide same gas to all transactions so that we can compare the cost.
		cost1Word, _ := createAndCheckContractWithCodeLenAndGas(t, net, contract, variant, 30, gasFor3Words)
		cost2Words, _ := createAndCheckContractWithCodeLenAndGas(t, net, contract, variant, 42, gasFor3Words)
		cost3Words, _ := createAndCheckContractWithCodeLenAndGas(t, net, contract, variant, 90, gasFor3Words)
		require := require.New(t)
		require.Equal(wordCost, cost2Words-cost1Word, "cost difference between 1 and 2")
		require.Equal(wordCost*2, cost3Words-cost1Word, "cost difference between 1 and 3")
	})

	t.Run("fails without enough gas", func(t *testing.T) {
		// 4 for a zero byte, 1 to make it fail.
		receipt := createContractWithCodeLenAndGas(t, net, variant, 1, gasForContract-4-1)
		require := require.New(t)
		require.Equal(types.ReceiptStatusFailed, receipt.Status,
			"unexpectedly succeeded to create contract without enough gas")
	})

	t.Run("with max init code size", func(t *testing.T) {
		receipt := createContractWithCodeLenAndGas(t, net, variant, MAX_INIT_CODE_SIZE, sufficientGas)
		require := require.New(t)
		require.Equal(types.ReceiptStatusSuccessful, receipt.Status,
			"failed to create contract with code length ", MAX_INIT_CODE_SIZE)
	})

	t.Run("aborts with init code size larger than MAX_INITCODE_SIZE", func(t *testing.T) {
		receipt := createContractWithCodeLenAndGas(t, net, variant, MAX_INIT_CODE_SIZE+1, sufficientGas)
		require := require.New(t)
		require.Equal(types.ReceiptStatusFailed, receipt.Status,
			"unexpectedly succeeded to create contract with init code length greater than MAX_INITCODE_SIZE")
	})
}

func createAndCheckContractWithCodeLenAndGas(t *testing.T, net *IntegrationTestNet,
	contract *initcode.Initcode, variant variant, codeLen, gasLimit uint64) (uint64, uint64) {
	receipt := createContractWithCodeLenAndGas(t, net, variant, codeLen, gasLimit)
	require := require.New(t)
	require.Equal(types.ReceiptStatusSuccessful, receipt.Status,
		"failed to create contract with code length ", codeLen)
	log, err := contract.ParseLogCost(*receipt.Logs[0])
	require.NoError(err)
	return log.Cost.Uint64(), receipt.GasUsed
}

func createContractWithCodeLenAndGas(t *testing.T, net *IntegrationTestNet, variant variant, codeLen, gasLimit uint64) *types.Receipt {
	receipt, err := net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		opts.GasLimit = gasLimit
		return variant(opts, big.NewInt(int64(codeLen)))
	})
	require := require.New(t)
	require.NoError(err)
	return receipt
}

type variant func(opts *bind.TransactOpts, codeSize *big.Int) (*types.Transaction, error)

func runTransactionWithCodeSizeAndGas(t *testing.T, net *IntegrationTestNet, codeSize, gas uint64) (*types.Receipt, error) {
	require := require.New(t)
	// these values are needed for the transaction but are irrelevant for the test
	client, err := net.GetClient()
	require.NoError(err, "failed to connect to the network:")
	defer client.Close()

	chainId, err := client.ChainID(context.Background())
	require.NoError(err, "failed to get chain ID::")

	nonce, err := client.NonceAt(context.Background(), net.validator.Address(), nil)
	require.NoError(err, "failed to get nonce:")

	price, err := client.SuggestGasPrice(context.Background())
	require.NoError(err, "failed to get gas price:")
	// ---------

	transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
		ChainID:  chainId,
		Gas:      gas,
		GasPrice: price,
		To:       nil,
		Nonce:    nonce,
		Data:     make([]byte, codeSize),
	}), types.NewLondonSigner(chainId), net.validator.PrivateKey)
	require.NoError(err, "failed to sign transaction:")
	return net.Run(transaction)
}
