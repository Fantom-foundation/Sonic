package tests

import (
	"context"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/counter_event_emitter"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestTransactionOrder(t *testing.T) {
	const (
		numAccounts = uint64(6)
		numPerAcc   = uint64(6)
		numBlocks   = uint64(3)
		numTxs      = numAccounts * numPerAcc
	)
	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(t, err)
	defer net.Stop()

	contract, _, err := DeployContract(net, counter_event_emitter.DeployCounterEventEmitter)
	require.NoError(t, err)

	client, err := net.GetClient()
	require.NoError(t, err)
	defer client.Close()

	accounts := make([]*Account, 0, numAccounts)

	// Only transactions from different accounts can change order.
	for range numAccounts {
		accounts = append(accounts, makeAccountWithMaxBalance(t, net))
	}

	// Repeat the test for X number of blocks
	for range numBlocks {
		blockNrBefore, err := client.BlockNumber(context.Background())
		require.NoError(t, err)

		options := make([]bind.TransactOpts, 0, numTxs)
		// Each account creates M transactions
		for _, acc := range accounts {
			opts, err := net.GetTransactOptions(acc)
			require.NoError(t, err)
			for range numPerAcc {
				options = append(options, *opts)
				opts.Nonce = new(big.Int).SetUint64(opts.Nonce.Uint64() + 1)
			}
		}

		// Pseudo-random shuffle to check that processor correctly orders transactions
		rand.Shuffle(len(options), func(i, j int) {
			options[i], options[j] = options[j], options[i]
		})

		transactions := make(types.Transactions, 0, numTxs)
		// Execute shuffled transactions
		for _, opts := range options {
			tx, err := contract.Increment(&opts)
			require.NoError(t, err)
			transactions = append(transactions, tx)
		}

		// Check that correct number of transactions has been sent
		if got, want := uint64(len(transactions)), numTxs; got != want {
			t.Fatalf("unexpected number of transactions, got: %d, want: %d", got, want)
		}

		// Check that the value in receipt is incremented by one - signals the transactions are ordered
		for _, tx := range transactions {
			receipt, err := net.GetReceipt(tx.Hash()) // first query synchronizes the execution
			require.NoError(t, err)
			count, err := contract.ParseCount(*receipt.Logs[0])
			require.NoError(t, err)
			// Nonce starts at 0 and count starts at 1 per account
			accCount := count.PerAddrCount.Uint64()
			nonce := tx.Nonce() + 1
			if accCount != nonce {
				t.Fatalf("transactions are not ordered, got idx: %d, want idx: %d", accCount, nonce)
			}
		}
		blockNrAfter, err := client.BlockNumber(context.Background())
		require.NoError(t, err)
		// At least one block between iterations must be generated
		// Multiple blocks between iterations can be generated
		if blockNrBefore >= blockNrAfter {
			t.Fatalf("no new block generated between iterations")
		}
	}

	gotCount, err := contract.GetTotalCount(nil)
	require.NoError(t, err)

	if got, want := gotCount.Uint64(), numTxs*numBlocks; got != want {
		t.Errorf("wrong count, got: %d, want: %d", got, want)
	}

	// Check that transactions are ordered correctly in the blockchain and that
	// for each transaction a correct receipt is available.
	globalCounter := uint64(0)
	context := context.Background()
	lastBlock, err := client.BlockNumber(context)
	require.NoError(t, err)
	for i := range lastBlock + 1 {
		block, err := client.BlockByNumber(context, big.NewInt(int64(i)))
		require.NoError(t, err)
		for i, tx := range block.Transactions() {
			receipt, err := client.TransactionReceipt(context, tx.Hash())
			require.NoError(t, err)

			// Check that the receipt matches to the transaction.
			require.Equal(t, receipt.Status, types.ReceiptStatusSuccessful)
			require.Equal(t, receipt.TxHash, tx.Hash())
			require.Equal(t, receipt.BlockHash, block.Hash())
			require.Equal(t, receipt.BlockNumber, block.Number())
			require.Equal(t, receipt.TransactionIndex, uint(i))

			// Check whether the receipt is for a counter transaction.
			if len(receipt.Logs) != 1 {
				continue
			}
			count, err := contract.ParseCount(*receipt.Logs[0])
			if err != nil {
				continue
			}

			// Check that transactions have been processed in order.
			require.Equal(t, count.PerAddrCount.Uint64(), tx.Nonce()+1)
			require.Equal(t, count.TotalCount.Uint64(), globalCounter+1)
			globalCounter++
		}
	}
	require.Equal(t, globalCounter, numTxs*numBlocks)
}

// makeAccountWithMaxBalance creates a new account and endows it with math.MaxInt64 balance.
// Creating the account this way allows to get access to the private key to sign transactions.
func makeAccountWithMaxBalance(t *testing.T, net *IntegrationTestNet) *Account {
	t.Helper()
	account := NewAccount()
	receipt, err := net.EndowAccount(account.Address(), math.MaxInt64)
	require.NoError(t, err)
	require.Equal(t,
		receipt.Status, types.ReceiptStatusSuccessful,
		"endowing account failed")
	return account
}
