package integration_tests

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"syscall"
	"time"

	sonicd "github.com/Fantom-foundation/go-opera/cmd/sonicd/app"
	sonictool "github.com/Fantom-foundation/go-opera/cmd/sonictool/app"
	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
)

type IntegrationTestNet struct {
	done <-chan struct{}

	validatorKey *ecdsa.PrivateKey
}

// StartIntegrationTestNet starts a single-node test network for integration tests.
// The node serving the network is started in the same process as the caller. This
// is intended to facilitate debugging of client code in the context of a running
// node.
func StartIntegrationTestNet(directory string) (*IntegrationTestNet, error) {
	done := make(chan struct{})
	go func() {
		defer close(done)

		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()

		// initialize the data directory for the single node on the test network
		// equivalent to running `sonictool --datadir <dataDir> genesis fake 1`
		os.Args = []string{"sonictool", "--datadir", directory, "genesis", "fake", "1"}
		sonictool.Run()

		// start the fakenet sonic node
		// equivalent to running `sonicd ...` but in this local process
		os.Args = []string{
			"sonicd",
			"--datadir", directory,
			"--fakenet", "1/1",
			"--http", "--http.addr", "0.0.0.0", "--http.port", "18545",
			"--http.api", "admin,eth,web3,net,txpool,ftm,trace,debug",
			"--ws", "--ws.addr", "0.0.0.0", "--ws.port", "18546", "--ws.api", "admin,eth,ftm",
			"--pprof", "--pprof.addr", "0.0.0.0",
			"--datadir.minfreedisk", "0",
		}
		sonicd.Run()
	}()

	result := &IntegrationTestNet{
		done:         done,
		validatorKey: evmcore.FakeKey(1),
	}

	// connect to blockchain network
	client, err := result.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}
	defer client.Close()

	const timeout = 30 * time.Second
	start := time.Now()

	// wait for the node to be ready to serve requests
	for time.Since(start) < timeout {
		id, err := client.ChainID(context.Background())
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("Managed to get the chain ID: %d\n", id)
		return result, nil
	}

	return nil, fmt.Errorf("failed to successfully start up a test network within %d", timeout)
}

func (n *IntegrationTestNet) stop() {
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	<-n.done
}

func (n *IntegrationTestNet) run(tx *types.Transaction) (*types.Receipt, error) {
	client, err := n.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	for i := 0; err != nil && i < 10; i++ {
		time.Sleep(1 * time.Second)
		receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}
	return receipt, nil
}

func (n *IntegrationTestNet) endowAccount(
	address common.Address,
	value *uint256.Int,
) error {
	client, err := n.getClient()
	if err != nil {
		return fmt.Errorf("failed to connect to the network: %w", err)
	}
	defer client.Close()

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	validatorAddress := crypto.PubkeyToAddress(n.validatorKey.PublicKey)
	nonce, err := client.NonceAt(context.Background(), validatorAddress, nil)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
		ChainID:  chainId,
		Gas:      21000,
		GasPrice: price,
		To:       &address,
		Value:    value.ToBig(),
		Nonce:    nonce,
	}), types.NewLondonSigner(chainId), n.validatorKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}
	_, err = n.run(transaction)
	return err
}

func (n *IntegrationTestNet) getClient() (*ethclient.Client, error) {
	return ethclient.Dial("http://localhost:18545")
}
