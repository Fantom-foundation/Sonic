package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestRejectedTx(t *testing.T) {
	require := require.New(t)

	// start network
	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(err)
	defer net.Stop()

	// create a client
	client, err := net.GetClient()
	require.NoError(err, "failed to get client")

	// make a dynamic tx that cannot be afford
	value := big.NewInt(42)
	gas := int64(21000)

	testCases := map[string]struct {
		nonce   uint64
		txMaker func(t *testing.T, net *IntegrationTestNet, tc testConfig) *types.Transaction
	}{
		"LegacyTx": {
			nonce:   0,
			txMaker: makeLegacyTxWithValue,
		},
		"AccessListTx": {
			nonce:   1,
			txMaker: makeAccessListTxWithValue,
		},
		"DynamicTx": {
			nonce:   2,
			txMaker: makeDynamicTxWithValue,
		},
		"BlobTx": {
			nonce:   3,
			txMaker: makeBlobTxWithValue,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			maxFeeCap := getMaxFee(t, client)
			newAccount := NewAccount()
			testConfig := testConfig{
				account:   newAccount,
				value:     value,
				gas:       uint64(gas),
				maxFeeCap: big.NewInt(maxFeeCap),
				nonce:     tc.nonce,
			}
			testRejectedTx(t, net, testConfig, tc.txMaker(t, net, testConfig))
		})
	}
}

func testRejectedTx(t *testing.T, net *IntegrationTestNet, testConfig testConfig, tx *types.Transaction) {
	require := require.New(t)

	// create a client
	client, err := net.GetClient()
	require.NoError(err, "failed to get client")

	balance := tx.Gas()*tx.GasFeeCap().Uint64() + tx.Value().Uint64()

	_, err = net.EndowAccount(testConfig.account.Address(), int64(balance-1))
	require.NoError(err, "failed to endow account")

	err = client.SendTransaction(context.Background(), tx)
	require.ErrorContains(err, "insufficient funds")

	_, err = net.EndowAccount(testConfig.account.Address(), int64(1))
	require.NoError(err, "failed to endow account")

	err = client.SendTransaction(context.Background(), tx)
	require.NoError(err)
}

func makeLegacyTxWithValue(
	t *testing.T,
	net *IntegrationTestNet,
	tc testConfig,
) *types.Transaction {
	chainId := getChainId(t, net)
	transaction, err := types.SignTx(types.NewTx(&types.LegacyTx{
		Gas:      tc.gas,
		GasPrice: tc.maxFeeCap,
		To:       &common.Address{},
		Nonce:    tc.nonce,
		Value:    tc.value,
	}), types.NewLondonSigner(chainId), tc.account.PrivateKey)
	require.NoError(t, err, "failed to sign transaction:")
	return transaction
}

func makeAccessListTxWithValue(
	t *testing.T,
	net *IntegrationTestNet,
	tc testConfig,
) *types.Transaction {
	chainId := getChainId(t, net)
	transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
		ChainID:  chainId,
		Gas:      tc.gas,
		GasPrice: tc.maxFeeCap,
		To:       &common.Address{},
		Nonce:    tc.nonce,
		Value:    tc.value,
	}), types.NewLondonSigner(chainId), tc.account.PrivateKey)
	require.NoError(t, err, "failed to sign transaction:")
	return transaction
}

func makeDynamicTxWithValue(t *testing.T, net *IntegrationTestNet, tc testConfig) *types.Transaction {
	chainId := getChainId(t, net)
	transaction, err := types.SignTx(types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainId,
		Gas:       tc.gas,
		GasFeeCap: tc.maxFeeCap,
		To:        &common.Address{},
		Nonce:     tc.nonce,
		Value:     tc.value,
	}), types.NewLondonSigner(chainId), tc.account.PrivateKey)
	require.NoError(t, err, "failed to sign transaction:")
	return transaction
}

func makeBlobTxWithValue(t *testing.T, net *IntegrationTestNet, tc testConfig) *types.Transaction {
	chainId := getChainId(t, net)
	transaction, err := types.SignTx(types.NewTx(&types.BlobTx{
		ChainID:    uint256.MustFromBig(chainId),
		Gas:        tc.gas,
		GasFeeCap:  uint256.MustFromBig(tc.maxFeeCap),
		GasTipCap:  uint256.MustFromBig(big.NewInt(0)),
		Nonce:      tc.nonce,
		Value:      uint256.MustFromBig(tc.value),
		BlobFeeCap: uint256.NewInt(3e10), // fee cap for the blob data
		BlobHashes: nil,                  // blob hashes in the transaction
		Sidecar:    nil,                  // sidecar data in the transaction
	}), types.NewCancunSigner(chainId), tc.account.PrivateKey)
	require.NoError(t, err, "failed to sign transaction:")
	return transaction
}

func getChainId(t *testing.T, net *IntegrationTestNet) *big.Int {
	client, err := net.GetClient()
	require.NoError(t, err, "failed to get client")

	chainId, err := client.ChainID(context.Background())
	require.NoError(t, err, "failed to get chain ID::")
	return chainId
}

func getMaxFee(t *testing.T, client *ethclient.Client) (maxFeeCap int64) {
	block, err := client.BlockByNumber(context.Background(), nil)
	require.NoError(t, err, "failed to get block by number")
	baseFee := int64(block.BaseFee().Uint64())
	maxFeeCap = baseFee + int64(float64(baseFee)*0.06)
	return
}

type testConfig struct {
	account   *Account
	value     *big.Int
	gas       uint64
	maxFeeCap *big.Int
	nonce     uint64
}
