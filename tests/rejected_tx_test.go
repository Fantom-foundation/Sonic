package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestRejectedTx(t *testing.T) {
	require := require.New(t)

	// start network
	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(err)
	defer net.Stop()

	newAccount := NewAccount()

	// create a client
	client, err := net.GetClient()
	require.NoError(err, "failed to get client")

	// make a dynamic tx that cannot be afford
	value := big.NewInt(42)
	gas := int64(21000)
	maxFeeCap := getMaxFee(t, client)

	tx := makeDynamicTxWithValue(t, net, newAccount, value, big.NewInt(maxFeeCap), uint64(gas))
	balance := int64(tx.Gas()*tx.GasFeeCap().Uint64()) + tx.Value().Int64()
	theoreticalCost := value.Int64() + gas*maxFeeCap
	require.Equal(balance, theoreticalCost, "balance and theoretical cost should be equal")

	_, err = net.EndowAccount(newAccount.Address(), balance-1-42)
	require.NoError(err, "failed to endow account")

	err = client.SendTransaction(context.Background(), tx)
	require.ErrorContains(err, "insufficient funds")

	_, err = net.EndowAccount(newAccount.Address(), balance)
	require.NoError(err, "failed to endow account")

	err = client.SendTransaction(context.Background(), tx)
	require.NoError(err)

}

func makeDynamicTxWithValue(t *testing.T, net *IntegrationTestNet, account *Account, value, maxFeeCap *big.Int, gas uint64) *types.Transaction {
	chainId, nonce := getChainAndNonce(t, net, account.Address())

	transaction, err := types.SignTx(types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainId,
		Gas:       gas,
		GasFeeCap: maxFeeCap,
		To:        &common.Address{},
		Nonce:     nonce,
		Value:     value,
	}), types.NewLondonSigner(chainId), account.PrivateKey)
	require.NoError(t, err, "failed to sign transaction:")

	return transaction
}

func getChainAndNonce(t *testing.T, net *IntegrationTestNet, address common.Address) (chainId *big.Int, nonce uint64) {
	// create a client
	client, err := net.GetClient()
	require.NoError(t, err, "failed to get client")

	chainId, err = client.ChainID(context.Background())
	require.NoError(t, err, "failed to get chain ID::")

	nonce, err = client.NonceAt(context.Background(), address, nil)
	require.NoError(t, err, "failed to get nonce:")
	return
}

func getMaxFee(t *testing.T, client *ethclient.Client) (maxFeeCap int64) {
	// get the gas limit and max fee cap from the client
	block, err := client.BlockByNumber(context.Background(), nil)
	require.NoError(t, err, "failed to get block by number")
	baseFee := int64(block.BaseFee().Uint64())
	maxFeeCap = baseFee + int64(float64(baseFee)*0.06)
	return
}
