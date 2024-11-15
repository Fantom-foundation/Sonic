package tests

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestGasPrice_GasEvolvesAsExpectedCalculates(t *testing.T) {
	require := require.New(t)

	net, client := makeNetAndClient(t)

	suggestions := []uint64{}
	prices := []uint64{}

	for i := 0; i < 10; i++ {

		suggestedPrice, err := client.SuggestGasPrice(context.Background())
		require.NoError(err)

		// new block
		receipt, err := net.EndowAccount(common.Address{42}, 100)
		require.NoError(err)

		lastBlock, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
		require.NoError(err)

		// store suggested and actuaaal prices.
		suggestions = append(suggestions, suggestedPrice.Uint64())
		prices = append(prices, lastBlock.BaseFee().Uint64())
	}

	for i := range suggestions {
		require.True(withinXPercent(suggestions[i], prices[i], 10), "i", i)
	}
}

func TestGasPrice_UnderpricedTransactionIsRejected(t *testing.T) {
	require := require.New(t)

	net, client := makeNetAndClient(t)

	// new block
	receipt, err := net.EndowAccount(common.Address{42}, 100)
	require.NoError(err)

	lastBlock, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
	require.NoError(err)

	lastBlockBaseFee := lastBlock.BaseFee().Uint64()

	// create and run a legacy transaction with a price lower than the suggested price
	transaction := makeLegacyTransactionWithPrice(t, net, int64(lastBlockBaseFee))
	_, err = net.Run(transaction)
	require.ErrorContains(err, "transaction underpriced")

	// create and run a access list transaction with a price lower than the suggested price
	transaction = makeAccessListTransactionWithPrice(t, net, int64(lastBlockBaseFee))
	_, err = net.Run(transaction)
	require.ErrorContains(err, "transaction underpriced")

	// create and run a legacy transaction with a price lower than the suggested price
	transaction = makeDynamicFeeTransactionWithPrice(t, net, int64(lastBlockBaseFee))
	_, err = net.Run(transaction)
	require.ErrorContains(err, "transaction underpriced")

	// create and run a legacy transaction with a price lower than the suggested price
	transaction = makeBlobTransactionWithPrice(t, net, int64(lastBlockBaseFee))
	_, err = net.Run(transaction)
	require.ErrorContains(err, "transaction underpriced")
}

func makeNetAndClient(t *testing.T) (*IntegrationTestNet, *ethclient.Client) {
	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(t, err)

	client, err := net.GetClient()
	require.NoError(t, err)

	return net, client
}

func withinXPercent(a, b uint64, margin float64) bool {
	// calculate the difference
	diff := uint64(0)
	if a > b {
		diff = a - b
	} else {
		diff = b - a
	}

	percentage := uint64(math.Round(float64(a) * margin / 100))

	// check if the difference is less than 10% of a
	return diff <= percentage
}

func makeLegacyTransactionWithPrice(t *testing.T, net *IntegrationTestNet, price int64) *types.Transaction {
	require := require.New(t)

	chainId, nonce := getChainIDAndNonce(t, net)

	transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
		ChainID:  chainId,
		Gas:      21_000,
		GasPrice: big.NewInt(price),
		To:       &common.Address{},
		Nonce:    nonce,
	}), types.NewLondonSigner(chainId), net.validator.PrivateKey)
	require.NoError(err, "failed to sign transaction:")
	return transaction
}

func makeAccessListTransactionWithPrice(t *testing.T, net *IntegrationTestNet, price int64) *types.Transaction {
	require := require.New(t)

	chainId, nonce := getChainIDAndNonce(t, net)

	transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
		ChainID:  chainId,
		Gas:      21_000,
		GasPrice: big.NewInt(price),
		To:       &common.Address{},
		Nonce:    nonce,
	}), types.NewLondonSigner(chainId), net.validator.PrivateKey)
	require.NoError(err, "failed to sign transaction:")
	return transaction
}

func makeDynamicFeeTransactionWithPrice(t *testing.T, net *IntegrationTestNet, price int64) *types.Transaction {
	require := require.New(t)

	chainId, nonce := getChainIDAndNonce(t, net)

	transaction, err := types.SignTx(types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainId,
		Gas:       21_000,
		GasFeeCap: big.NewInt(price),
		GasTipCap: big.NewInt(price),
		To:        &common.Address{},
		Nonce:     nonce,
	}), types.NewLondonSigner(chainId), net.validator.PrivateKey)
	require.NoError(err, "failed to sign transaction:")
	return transaction
}

func makeBlobTransactionWithPrice(t *testing.T, net *IntegrationTestNet, price int64) *types.Transaction {
	require := require.New(t)

	chainId, nonce := getChainIDAndNonce(t, net)

	transaction, err := types.SignTx(types.NewTx(&types.BlobTx{
		ChainID:    uint256.MustFromBig(chainId),
		Gas:        21_000,
		GasFeeCap:  uint256.MustFromBig(big.NewInt(price)),
		GasTipCap:  uint256.MustFromBig(big.NewInt(price)),
		Nonce:      nonce,
		BlobFeeCap: uint256.NewInt(3e10), // fee cap for the blob data
		BlobHashes: nil,                  // blob hashes in the transaction
		Sidecar:    nil,                  // sidecar data in the transaction
	}), types.NewCancunSigner(chainId), net.validator.PrivateKey)
	require.NoError(err, "failed to sign transaction:")
	return transaction
}

func getChainIDAndNonce(t *testing.T, net *IntegrationTestNet) (*big.Int, uint64) {
	require := require.New(t)
	// these values are needed for the transaction but are irrelevant for the test
	client, err := net.GetClient()
	require.NoError(err, "failed to connect to the network:")
	defer client.Close()

	chainId, err := client.ChainID(context.Background())
	require.NoError(err, "failed to get chain ID::")

	nonce, err := client.NonceAt(context.Background(), net.validator.Address(), nil)
	require.NoError(err, "failed to get nonce:")
	return chainId, nonce
}
