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

	// create a client
	client, err := net.GetClient()
	require.NoError(err, "failed to get client")

	chainId := getChainId(t, net)

	testCases := []struct {
		name    string
		txMaker func(t *testing.T, account *Account, nonce uint64, price int64) *types.Transaction
	}{
		{name: "LegacyTx",
			txMaker: func(t *testing.T, account *Account, nonce uint64, price int64) *types.Transaction {
				factory := &txFactory{
					senderKey: account.PrivateKey,
					chainId:   chainId,
				}
				return factory.makeLegacyTransactionWithPrice(t, nonce, price)
			},
		}, {name: "AccessListTx",
			txMaker: func(t *testing.T, account *Account, nonce uint64, price int64) *types.Transaction {
				factory := &txFactory{
					senderKey: account.PrivateKey,
					chainId:   chainId,
				}
				return factory.makeAccessListTransactionWithPrice(t, nonce, price)
			},
		}, {name: "DynamicTx",
			txMaker: func(t *testing.T, account *Account, nonce uint64, price int64) *types.Transaction {
				factory := &txFactory{
					senderKey: account.PrivateKey,
					chainId:   chainId,
				}
				return factory.makeDynamicFeeTransactionWithPrice(t, nonce, price)
			},
		}, {name: "BlobTx",
			txMaker: func(t *testing.T, account *Account, nonce uint64, price int64) *types.Transaction {
				factory := &txFactory{
					senderKey: account.PrivateKey,
					chainId:   chainId,
				}
				return factory.makeBlobTransactionWithPrice(t, nonce, price)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			maxFeeCap := getMaxFee(t, client)
			newAccount := NewAccount()
			nonce := getNonce(t, client, newAccount.Address())
			testRejectedTx(t, net, newAccount.Address(), tc.txMaker(t, newAccount, nonce, maxFeeCap))
		})
	}
}

func testRejectedTx(t *testing.T, net *IntegrationTestNet, account common.Address, tx *types.Transaction) {
	require := require.New(t)

	// verify estimated cost
	estimatedCost := tx.Gas()*tx.GasFeeCap().Uint64() + tx.Value().Uint64()
	require.Equal(tx.Cost().Uint64(), estimatedCost, "cost of transaction is not equal to balance")

	// provide just enough balance to NOT cover the cost
	_, err := net.EndowAccount(account, int64(estimatedCost-1))
	require.NoError(err, "failed to endow account")

	// run transaction to be rejected
	receipt, err := net.Run(tx)
	require.ErrorContains(err, "insufficient funds")
	require.Nil(receipt)

	// provide enough balance to cover the cost
	_, err = net.EndowAccount(account, int64(1))
	require.NoError(err, "failed to endow account")

	// run transaction to be successful
	receipt, err = net.Run(tx)
	require.NoError(err)

	// verify receipt
	require.Equal(receipt.Status, uint64(1))
	require.Equal(tx.Gas(), receipt.GasUsed)
	require.GreaterOrEqual(tx.GasFeeCap().Uint64(), receipt.EffectiveGasPrice.Uint64())
	require.GreaterOrEqual(tx.Cost().Uint64(), receipt.EffectiveGasPrice.Uint64()*receipt.GasUsed)
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
	baseFee := block.BaseFee().Int64()
	maxFeeCap = int64(float64(baseFee) * 1.06)
	return
}

func getNonce(t *testing.T, client *ethclient.Client, account common.Address) (nonce uint64) {
	nonce, err := client.NonceAt(context.Background(), account, nil)
	require.NoError(t, err, "failed to get nonce:")
	return
}
