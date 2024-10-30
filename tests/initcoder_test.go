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
	r := require.New(t)

	net, err := StartIntegrationTestNet(t.TempDir())
	r.NoError(err)
	defer net.Stop()

	// Deploy the invalid start contract.
	contract, receipt, err := DeployContract(net, initcode.DeployInitcode)
	r.NoError(err)
	r.Equal(types.ReceiptStatusSuccessful, receipt.Status, "failed to deploy contract")

	const gasForOneWord uint64 = 54977

	t.Run("create transaction with enough gas for init code", func(t *testing.T) {
		// we need to provide same gas to all transactions so that we can compare the cost.
		cost1Word, _ := createAndCheckContractWithCodeLenAndGas(r, net, contract, 30, gasForOneWord+4)
		cost2Words, _ := createAndCheckContractWithCodeLenAndGas(r, net, contract, 42, gasForOneWord+4)
		cost3Words, _ := createAndCheckContractWithCodeLenAndGas(r, net, contract, 90, gasForOneWord+4)
		if cost2Words-cost1Word != 2 {
			t.Errorf("cost difference between 1 and 2 words should be 2, instead got %d", cost2Words-cost1Word)
		}
		if cost3Words-cost1Word != 4 {
			t.Errorf("cost difference between 1 and 3 words should be 4, instead got %d", cost3Words-cost1Word)
		}
	})

	t.Run("create transaction without enough gas for init code", func(t *testing.T) {
		receipt = createContractWithCodeLenAndGas(r, net, contract, 42, gasForOneWord)
		r.Equal(types.ReceiptStatusFailed, receipt.Status,
			"unexpectedly succeeded to create contract without enough gas")
	})

	const MAX_INITCODE_SIZE = 49152
	const MAX_INITCODE_COST = (MAX_INITCODE_SIZE / 32) * 2
	const someCost = gasForOneWord + MAX_INITCODE_COST + 9_199
	t.Run("create transaction with max init code size", func(t *testing.T) {
		_, _ = createAndCheckContractWithCodeLenAndGas(r, net, contract, MAX_INITCODE_SIZE, someCost)
	})

	t.Run("create transaction with MAX_INITCODE_SIZE+1", func(t *testing.T) {
		receipt = createContractWithCodeLenAndGas(r, net, contract, MAX_INITCODE_SIZE+1, someCost+2)
		r.Equal(types.ReceiptStatusFailed, receipt.Status,
			"unexpectedly succeeded to create contract with init code length greater than MAX_INITCODE_SIZE")
	})
}

func createAndCheckContractWithCodeLenAndGas(r *require.Assertions, net *IntegrationTestNet, contract *initcode.Initcode, codeLen, gasLimit uint64) (uint64, uint64) {
	receipt := createContractWithCodeLenAndGas(r, net, contract, codeLen, gasLimit)
	r.Equal(types.ReceiptStatusSuccessful, receipt.Status, "failed to create contract with code length %d", codeLen)

	log, err := contract.ParseLogCost(*receipt.Logs[0])
	r.NoError(err)
	return log.Cost.Uint64(), receipt.GasUsed
}

func createContractWithCodeLenAndGas(r *require.Assertions, net *IntegrationTestNet, contract *initcode.Initcode, codeLen, gasLimit uint64) *types.Receipt {
	receipt, err := net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		opts.GasLimit = gasLimit
		return contract.CreateContractWith(opts, big.NewInt(int64(codeLen)))
	})
	fmt.Printf("gas provided: %v, and gas used: %v\n", gasLimit, receipt.GasUsed) // THIS IS A DEBUG LINE
	r.NoError(err)
	return receipt
}
