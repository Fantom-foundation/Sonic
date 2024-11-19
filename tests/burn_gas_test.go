package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/burn_gas"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestLegacyTransaction_BurnGas(t *testing.T) {

	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(t, err)
	defer net.Stop()

	client, err := net.GetClient()
	require.NoError(t, err)
	defer client.Close()

	contract, receipt, err := DeployContract(net, burn_gas.DeployBurnGas)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, types.ReceiptStatusSuccessful)

	// send a simple transaction to query the gas price (no tip)
	tx := makeEip1559Transaction(t, client, 105e9, 0, &net.validator)
	receipt, err = net.Run(tx)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, types.ReceiptStatusSuccessful)
	gasFeeBefore := receipt.EffectiveGasPrice
	t.Log("gasFee before", gasFeeBefore)

	balanceBefore, err := client.BalanceAt(context.Background(), net.validator.Address(), nil)
	require.NoError(t, err)

	receipt, err = net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		opts.GasPrice = nil            // <- no gas price forces eip1559 transaction (to define max tip)
		opts.GasTipCap = big.NewInt(0) // <- no tip should turn effective gas price == basefee
		opts.GasFeeCap = nil           // <- use estimation
		opts.GasLimit = 9_980_000      // <- this seems to be the max allowed gas
		return contract.BurnGas(opts)
	})
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	percentage := float64(100) / float64(gasFeeBefore.Int64())
	percentage *= float64(receipt.EffectiveGasPrice.Int64())
	t.Logf("BurnGas gas used %d gas, effectivePrice %d (%.2f%%)", receipt.GasUsed, receipt.EffectiveGasPrice,
		percentage)

	balanceAfter, err := client.BalanceAt(context.Background(), net.validator.Address(), nil)
	require.NoError(t, err)
	require.Equal(
		t,
		balanceBefore.Uint64()-balanceAfter.Uint64(),
		uint64(receipt.EffectiveGasPrice.Uint64()*receipt.GasUsed),
	)

	// send a simple transaction to query the gas price (no tip)
	tx = makeEip1559Transaction(t, client, 105e9, 0, &net.validator)
	receipt, err = net.Run(tx)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, types.ReceiptStatusSuccessful)
	gasFeeAfter := receipt.EffectiveGasPrice
	percentage = float64(100) / float64(gasFeeBefore.Int64())
	percentage *= float64(gasFeeAfter.Int64())
	t.Logf("gasFee after: %d (%.2f%%)", gasFeeAfter, percentage)
}
