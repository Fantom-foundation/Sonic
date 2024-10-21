package tests

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/integration/makefakegenesis"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestIntegrationTest_ContactNet(t *testing.T) {

	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	// connect to blockchain network
	client, err := net.GetClient()
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	block, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Fatalf("Failed to get block number: %v", err)
	}

	t.Logf("Latest block number: %v\n", block)

	id, err := client.NetworkID(context.Background())
	if err != nil {
		t.Errorf("Failed to get network ID: %v", err)
	} else {
		t.Logf("Network ID: %v\n", id)
	}

	balance, err := client.BalanceAt(context.Background(), common.Address{}, nil)
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	} else {
		t.Logf("Balance: %v\n", balance)
	}

	validators := makefakegenesis.GetFakeValidators(1)
	for _, v := range validators {
		t.Logf("Validator: %v\n", v)
		balance, err := client.BalanceAt(context.Background(), v.Address, nil)
		if err != nil {
			t.Fatalf("Failed to get balance: %v", err)
		} else {
			t.Logf("Balance: %v\n", balance)
		}
	}

	chainId := big.NewInt(int64(opera.FakeNetworkID))

	valdiatorKey := evmcore.FakeKey(1)

	validatorAddress := crypto.PubkeyToAddress(valdiatorKey.PublicKey)

	nonce, err := client.NonceAt(context.Background(), validatorAddress, nil)
	if err != nil {
		t.Fatalf("Failed to get nonce: %v", err)
	}
	t.Logf("Nonce: %v\n", nonce)

	price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatalf("Failed to get gas price: %v", err)
	}
	t.Logf("Gas price: %v\n", price)

	/*
		transaction, err := types.SignTx(types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: price,
			Gas:      21000,
			To:       &common.Address{},
			Value:    big.NewInt(1000),
		}), types.NewEIP155Signer(chainId), valdiatorKey)
		if err != nil {
			t.Fatalf("Failed to sign transaction: %v", err)
		}
	*/

	/*
		// Type 1
		transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
			ChainID:  chainId,
			Gas:      21000,
			GasPrice: price,
			To:       &common.Address{},
			Value:    big.NewInt(1000),
			Nonce:    nonce,
		}), types.NewLondonSigner(chainId), valdiatorKey)
		if err != nil {
			t.Fatalf("Failed to sign transaction: %v", err)
		}
	*/

	// Type 2 -- Dynamic Fee Transactions (London)
	transaction, err := types.SignTx(types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainId,
		Gas:       21000,
		GasFeeCap: new(big.Int).Add(price, big.NewInt(1000)),
		GasTipCap: new(big.Int).Add(price, big.NewInt(1000)),
		To:        &common.Address{},
		Value:     big.NewInt(1000),
		Nonce:     nonce,
	}), types.NewLondonSigner(chainId), valdiatorKey)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	/*
		// Type 3 -- Blob Transactions (Cancun)
		transaction, err := types.SignTx(types.NewTx(&types.BlobTx{
			ChainID:   uint256.MustFromBig(chainId),
			Nonce:     nonce,
			To:        common.Address{},
			Value:     uint256.NewInt(1000),
			Gas:       21000,
			GasFeeCap: uint256.MustFromBig(new(big.Int).Add(price, big.NewInt(1000))),
			GasTipCap: uint256.MustFromBig(new(big.Int).Add(price, big.NewInt(1000))),
			//BlobHashes: []common.Hash{{}},
		}), types.NewCancunSigner(chainId), valdiatorKey)
		if err != nil {
			t.Fatalf("Failed to sign transaction: %v", err)
		}
	*/

	err = client.SendTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("Failed to send transaction: %v", err)
	}

	receipt, err := client.TransactionReceipt(context.Background(), transaction.Hash())
	for i := 0; err != nil && i < 10; i++ {
		time.Sleep(1 * time.Second)
		receipt, err = client.TransactionReceipt(context.Background(), transaction.Hash())
	}
	if err != nil {
		t.Fatalf("Failed to get transaction receipt: %v", err)
	}
	t.Logf("Receipt\n\tstatus: %d\n\tblock %d\n\ttransaction %d\n\tgas: %d\n\tgas price: %v\n",
		receipt.Status, receipt.BlockNumber, receipt.TransactionIndex,
		receipt.GasUsed, receipt.EffectiveGasPrice,
	)

	/*
		for i := 0; i < 10; i++ {

			block, err := client.BlockNumber(context.Background())
			if err != nil {
				t.Fatalf("Failed to get block number: %v", err)
			}

			t.Logf("Latest block number: %v\n", block)

			balance, err := client.BalanceAt(context.Background(), common.Address{}, nil)
			if err != nil {
				t.Fatalf("Failed to get balance: %v", err)
			} else {
				t.Logf("Balance: %v\n", balance)
			}

			time.Sleep(1 * time.Second)
		}
	*/

	//t.Fail()
}

/*
// transferValue transfer a financial value from account identified by given privateKey, to given toAddress.
// It returns when the value is already available on the target account.
func transferValue(rpcClient rpc.RpcClient, from *Account, toAddress common.Address, value *big.Int, gasPrice *big.Int) (err error) {
	signedTx, err := createTx(from, toAddress, value, nil, gasPrice, 21000)
	if err != nil {
		return err
	}
	return rpcClient.SendTransaction(context.Background(), signedTx)
}

func createTx(from *Account, toAddress common.Address, value *big.Int, data []byte, gasPrice *big.Int, gasLimit uint64) (*types.Transaction, error) {
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    from.getNextNonce(),
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	})
	return types.SignTx(tx, types.NewEIP155Signer(from.chainID), from.privateKey)
}
*/
