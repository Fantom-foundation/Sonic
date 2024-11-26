package tests

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestGasPrice_SuggestedGasPricesApproximateActualBaseFees(t *testing.T) {
	require := require.New(t)
	net, client := makeNetAndClient(t)

	fees := []uint64{}
	suggestions := []uint64{}
	ctxt := context.Background()
	for i := 0; i < 10; i++ {
		suggestedPrice, err := client.SuggestGasPrice(ctxt)
		require.NoError(err)

		// new block
		receipt, err := net.EndowAccount(common.Address{42}, 100)
		require.NoError(err)

		lastBlock, err := client.BlockByNumber(ctxt, receipt.BlockNumber)
		require.NoError(err)

		// store suggested and actual prices.
		suggestions = append(suggestions, suggestedPrice.Uint64())
		fees = append(fees, lastBlock.BaseFee().Uint64())
	}

	// Suggestions should over-estimate the actual prices by ~10%
	for i := 1; i < int(len(suggestions)); i++ {
		ratio := float64(suggestions[i]) / float64(fees[i-1])
		require.Less(1.09, ratio, "step %d, suggestion %d, fees %d", i, suggestions[i], fees[i-1])
		require.Less(ratio, 1.11, "step %d, suggestion %d, fees %d", i, suggestions[i], fees[i-1])
	}
}

func TestGasPrice_UnderpricedTransactionsAreRejected(t *testing.T) {
	require := require.New(t)

	net, client := makeNetAndClient(t)
	send := func(tx *types.Transaction) error {
		return client.SendTransaction(context.Background(), tx)
	}

	chainId, err := client.ChainID(context.Background())
	require.NoError(err, "failed to get chain ID::")

	nonce, err := client.NonceAt(context.Background(), net.validator.Address(), nil)
	require.NoError(err, "failed to get nonce:")

	factory := &txFactory{
		senderKey: net.validator.PrivateKey,
		chainId:   chainId,
	}

	lastBlock, err := client.BlockByNumber(context.Background(), nil)
	require.NoError(err)

	// Everything below ~5% above the base fee should be rejected.
	baseFee := int(lastBlock.BaseFee().Uint64())
	for _, extra := range []int{-10, 0, baseFee / 100, 4 * baseFee / 100} {
		feeCap := int64(baseFee + extra)

		err = send(factory.makeLegacyTransactionWithPrice(t, nonce, feeCap))
		require.ErrorContains(err, "transaction underpriced")

		err = send(factory.makeAccessListTransactionWithPrice(t, nonce, feeCap))
		require.ErrorContains(err, "transaction underpriced")

		err = send(factory.makeDynamicFeeTransactionWithPrice(t, nonce, feeCap))
		require.ErrorContains(err, "transaction underpriced")

		err = send(factory.makeBlobTransactionWithPrice(t, nonce, feeCap))
		require.ErrorContains(err, "transaction underpriced")
	}

	// Everything over ~5% above the base fee should be accepted.
	feeCap := int64(baseFee + 7*baseFee/100)
	require.NoError(send(factory.makeLegacyTransactionWithPrice(t, nonce, feeCap)))
	require.NoError(send(factory.makeAccessListTransactionWithPrice(t, nonce+1, feeCap)))
	require.NoError(send(factory.makeDynamicFeeTransactionWithPrice(t, nonce+2, feeCap)))
	require.NoError(send(factory.makeBlobTransactionWithPrice(t, nonce+3, feeCap)))
}

func TestGasPrice_SuggestedGasPriceFeeProvidesPriority(t *testing.T) {
	net, client := makeNetAndClient(t)

	chainId, err := client.ChainID(context.Background())
	require.NoError(t, err, "failed to get chain ID:")

	slowAccount := makeAccountWithBalance(t, net, 1e18)
	priorityAccount := makeAccountWithBalance(t, net, 1e18)

	factories := map[common.Address]*txFactory{
		slowAccount.Address(): {
			senderKey: slowAccount.PrivateKey,
			chainId:   chainId,
		},
		priorityAccount.Address(): {
			senderKey: priorityAccount.PrivateKey,
			chainId:   chainId,
		},
	}

	// This test uses 100 transactions because of at time of writing, the block includes
	// approximately 32 transactions. This test could start failing if the block size changes.
	// To restore expected behavior, adjust the number of transactions so that approximately 3
	// blocks are needed to dispatch the batch of transactions.
	slowTransactions := make([]*types.Transaction, 100)
	nonce, err := client.NonceAt(context.Background(), slowAccount.Address(), nil)
	require.NoError(t, err, "failed to get nonce")
	for i := range 100 {
		tx := factories[slowAccount.Address()].makeDynamicFeeTransactionWithPriceAndTip(
			t,
			nonce+uint64(i),
			105e9,
			int64(i*1000),
		)
		slowTransactions[i] = tx
	}
	for _, tx := range slowTransactions {
		err := client.SendTransaction(context.Background(), tx)
		require.NoError(t, err)
	}

	// Waiting for the first transaction receipt waits for the first block to be completed.
	// Because the amount of transactions is large, there will be transactions still in the pool.
	_, err = net.GetReceipt(slowTransactions[0].Hash())
	require.NoError(t, err)

	tipCap, err := client.SuggestGasTipCap(context.Background())
	require.NoError(t, err, "failed to get tip cap suggestion")
	nonce, err = client.NonceAt(context.Background(), priorityAccount.Address(), nil)
	require.NoError(t, err, "failed to get nonce")
	prioritizedTx := factories[priorityAccount.Address()].makeDynamicFeeTransactionWithPriceAndTip(
		t,
		nonce,
		105e9,
		tipCap.Int64(),
	)
	err = client.SendTransaction(context.Background(), prioritizedTx)
	require.NoError(t, err)

	// get the receipt of the last transaction in transactions (synchronizes the complete batch)
	lastTxReceipt, err := net.GetReceipt(slowTransactions[len(slowTransactions)-1].Hash())
	require.NoError(t, err)

	prioritizedTxReceipt, err :=
		client.TransactionReceipt(context.Background(), prioritizedTx.Hash())
	require.NoError(t, err, "prioritized transaction was not executed before the last slow transaction")

	// The prioritized transaction shall be included in an earlier block as the last slow transaction.
	// Because signing transactions is a heavy operation, the actual block may change depending
	// on the machine load. For this reason, no absolute block number is checked.
	require.Less(t,
		prioritizedTxReceipt.BlockNumber.Uint64(),
		lastTxReceipt.BlockNumber.Uint64(),
		"the prioritized transaction was not executed with priority")
}

func makeNetAndClient(t *testing.T) (*IntegrationTestNet, *ethclient.Client) {
	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() { net.Stop() })

	client, err := net.GetClient()
	require.NoError(t, err)
	t.Cleanup(func() { client.Close() })

	return net, client
}

type txFactory struct {
	senderKey *ecdsa.PrivateKey
	chainId   *big.Int
}

func (f *txFactory) makeLegacyTransactionWithPrice(
	t *testing.T,
	nonce uint64,
	price int64,
) *types.Transaction {
	transaction, err := types.SignTx(types.NewTx(&types.LegacyTx{
		Gas:      21_000,
		GasPrice: big.NewInt(price),
		To:       &common.Address{},
		Nonce:    nonce,
	}), types.NewLondonSigner(f.chainId), f.senderKey)
	require.NoError(t, err, "failed to sign transaction")
	return transaction
}

func (f *txFactory) makeAccessListTransactionWithPrice(
	t *testing.T,
	nonce uint64,
	price int64,
) *types.Transaction {
	transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
		ChainID:  f.chainId,
		Gas:      21_000,
		GasPrice: big.NewInt(price),
		To:       &common.Address{},
		Nonce:    nonce,
	}), types.NewLondonSigner(f.chainId), f.senderKey)

	require.NoError(t, err, "failed to sign transaction:")
	return transaction
}

func (f *txFactory) makeDynamicFeeTransactionWithPrice(
	t *testing.T,
	nonce uint64,
	price int64,
) *types.Transaction {
	transaction, err := types.SignTx(types.NewTx(&types.DynamicFeeTx{
		ChainID:   f.chainId,
		Gas:       21_000,
		GasFeeCap: big.NewInt(price),
		GasTipCap: big.NewInt(0),
		To:        &common.Address{},
		Nonce:     nonce,
	}), types.NewLondonSigner(f.chainId), f.senderKey)
	require.NoError(t, err, "failed to sign transaction:")
	return transaction
}

func (f *txFactory) makeDynamicFeeTransactionWithPriceAndTip(
	t *testing.T,
	nonce uint64,
	price int64,
	tip int64,
) *types.Transaction {
	transaction, err := types.SignTx(types.NewTx(&types.DynamicFeeTx{
		ChainID:   f.chainId,
		Gas:       21_000,
		GasFeeCap: big.NewInt(price),
		GasTipCap: big.NewInt(tip),
		To:        &common.Address{},
		Nonce:     nonce,
	}), types.NewLondonSigner(f.chainId), f.senderKey)
	require.NoError(t, err, "failed to sign transaction:")
	return transaction
}

func (f *txFactory) makeBlobTransactionWithPrice(
	t *testing.T,
	nonce uint64,
	price int64,
) *types.Transaction {
	transaction, err := types.SignTx(types.NewTx(&types.BlobTx{
		ChainID:    uint256.MustFromBig(f.chainId),
		Gas:        21_000,
		GasFeeCap:  uint256.MustFromBig(big.NewInt(price)),
		GasTipCap:  uint256.MustFromBig(big.NewInt(0)),
		Nonce:      nonce,
		BlobFeeCap: uint256.NewInt(3e10), // fee cap for the blob data
		BlobHashes: nil,                  // blob hashes in the transaction
		Sidecar:    nil,                  // sidecar data in the transaction
	}), types.NewCancunSigner(f.chainId), f.senderKey)
	require.NoError(t, err, "failed to sign transaction:")
	return transaction
}
