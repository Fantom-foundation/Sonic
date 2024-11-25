package tests

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"syscall"
	"time"

	sonicd "github.com/Fantom-foundation/go-opera/cmd/sonicd/app"
	sonictool "github.com/Fantom-foundation/go-opera/cmd/sonictool/app"
	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// IntegrationTestNet is a in-process test network for integration tests. When
// started, it runs a full Sonic node maintaining a chain within the process
// containing this object. The network can be used to run transactions on and
// to perform queries against.
//
// The main purpose of this network is to facilitate end-to-end debugging of
// client code in the controlled scope of individual unit tests. When running
// tests against an integration test network instance, break-points can be set
// in the client code, thereby facilitating debugging.
//
// A typical use case would look as follows:
//
//	func TestMyClientCode(t *testing.T) {
//	  net, err := StartIntegrationTestNet(t.TempDir())
//	  if err != nil {
//	    t.Fatalf("Failed to start the fake network: %v", err)
//	  }
//	  defer net.Stop()
//	  <run tests against the network>
//	}
//
// Additionally, by providing support for scripting test traffic on a network,
// integration test networks can also be used for automated integration and
// regression tests for client code.
type IntegrationTestNet struct {
	directory string
	done      <-chan struct{}
	validator Account
}

// StartIntegrationTestNet starts a single-node test network for integration tests.
// The node serving the network is started in the same process as the caller. This
// is intended to facilitate debugging of client code in the context of a running
// node.
func StartIntegrationTestNet(directory string) (*IntegrationTestNet, error) {

	// initialize the data directory for the single node on the test network
	// equivalent to running `sonictool --datadir <dataDir> genesis fake 1`
	originalArgs := os.Args
	os.Args = []string{"sonictool", "--datadir", directory, "genesis", "fake", "1"}
	sonictool.Run()
	os.Args = originalArgs

	// start the fakenet sonic node
	result := &IntegrationTestNet{
		directory: directory,
		validator: Account{evmcore.FakeKey(1)},
	}

	if err := result.start(); err != nil {
		return nil, fmt.Errorf("failed to start the test network: %w", err)
	}
	return result, nil
}

func (n *IntegrationTestNet) start() error {
	if n.done != nil {
		return errors.New("network already started")
	}
	done := make(chan struct{})
	go func() {
		defer close(done)

		// start the fakenet sonic node
		// equivalent to running `sonicd ...` but in this local process
		os.Args = []string{
			"sonicd",
			"--datadir", n.directory,
			"--fakenet", "1/1",
			"--http", "--http.addr", "0.0.0.0", "--http.port", "18545",
			"--http.api", "admin,eth,web3,net,txpool,ftm,trace,debug",
			"--ws", "--ws.addr", "0.0.0.0", "--ws.port", "18546", "--ws.api", "admin,eth,ftm",
			"--datadir.minfreedisk", "0",
		}
		sonicd.Run()
	}()

	n.done = done

	// connect to blockchain network
	client, err := n.GetClient()
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}
	defer client.Close()

	const timeout = 300 * time.Second
	start := time.Now()

	// wait for the node to be ready to serve requests
	const maxDelay = 100 * time.Millisecond
	delay := time.Millisecond
	for time.Since(start) < timeout {
		_, err := client.ChainID(context.Background())
		if err != nil {
			time.Sleep(delay)
			delay = 2 * delay
			if delay > maxDelay {
				delay = maxDelay
			}
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to successfully start up a test network within %v", timeout)
}

// Stop shuts the underlying network down.
func (n *IntegrationTestNet) Stop() {
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	<-n.done
	n.done = nil
}

// Stops and restarts the single node on the test network.
func (n *IntegrationTestNet) Restart() error {
	n.Stop()
	return n.start()
}

// EndowAccount sends a requested amount of tokens to the given account. This is
// mainly intended to provide funds to accounts for testing purposes.
func (n *IntegrationTestNet) EndowAccount(
	address common.Address,
	value int64,
) (*types.Receipt, error) {
	client, err := n.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the network: %w", err)
	}
	defer client.Close()

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// The requested funds are moved from the validator account to the target account.
	nonce, err := client.NonceAt(context.Background(), n.validator.Address(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
		ChainID:  chainId,
		Gas:      21000,
		GasPrice: price,
		To:       &address,
		Value:    big.NewInt(value),
		Nonce:    nonce,
	}), types.NewLondonSigner(chainId), n.validator.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	return n.Run(transaction)
}

// Run sends the given transaction to the network and waits for it to be processed.
// The resulting receipt is returned. This function times out after 10 seconds.
func (n *IntegrationTestNet) Run(tx *types.Transaction) (*types.Receipt, error) {
	client, err := n.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}
	defer client.Close()
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}
	return n.GetReceipt(tx.Hash())
}

// GetReceipt waits for the receipt of the given transaction hash to be available.
// The function times out after 10 seconds.
func (n *IntegrationTestNet) GetReceipt(txHash common.Hash) (*types.Receipt, error) {
	client, err := n.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}
	defer client.Close()

	// Wait for the response with some exponential backoff.
	const maxDelay = 100 * time.Millisecond
	now := time.Now()
	delay := time.Millisecond
	for time.Since(now) < 100*time.Second {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if errors.Is(err, ethereum.NotFound) {
			time.Sleep(delay)
			delay = 2 * delay
			if delay > maxDelay {
				delay = maxDelay
			}
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
		}
		return receipt, nil
	}
	return nil, fmt.Errorf("failed to get transaction receipt: timeout")
}

// Apply sends a transaction to the network using the network's validator account
// and waits for the transaction to be processed. The resulting receipt is returned.
func (n *IntegrationTestNet) Apply(
	issue func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Receipt, error) {
	txOpts, err := n.GetTransactOptions(&n.validator)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction options: %w", err)
	}
	transaction, err := issue(txOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return n.GetReceipt(transaction.Hash())
}

// GetTransactOptions provides transaction options to be used to send a transaction
// with the given account. The options include the chain ID, a suggested gas price,
// the next free nonce of the given account, and a hard-coded gas limit of 1e6.
// The main purpose of this function is to provide a convenient way to collect all
// the necessary information required to create a transaction in one place.
func (n *IntegrationTestNet) GetTransactOptions(account *Account) (*bind.TransactOpts, error) {
	client, err := n.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}
	defer client.Close()

	ctxt := context.Background()
	chainId, err := client.ChainID(ctxt)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctxt)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price suggestion: %w", err)
	}

	nonce, err := client.NonceAt(ctxt, account.Address(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	txOpts, err := bind.NewKeyedTransactorWithChainID(account.PrivateKey, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction options: %w", err)
	}
	txOpts.GasPrice = new(big.Int).Mul(gasPrice, big.NewInt(2))
	txOpts.Nonce = big.NewInt(int64(nonce))
	txOpts.GasLimit = 1e6
	return txOpts, nil
}

// GetClient provides raw access to a fresh connection to the network.
// The resulting client must be closed after use.
func (n *IntegrationTestNet) GetClient() (*ethclient.Client, error) {
	return ethclient.Dial("http://localhost:18545")
}

// RestartWithExportImport stops the network, exports the genesis file, cleans the
// temporary directory, imports the genesis file, and starts the network again.
func (n *IntegrationTestNet) RestartWithExportImport() error {
	n.Stop()
	fmt.Println("Network stopped. Exporting genesis file...")

	// save original args
	originalArgs := os.Args

	// export
	os.Args = []string{
		"sonictool",
		"--datadir", n.directory,
		"genesis", "export", "testGenesis.g",
	}
	err := sonictool.Run()
	if err != nil {
		return err
	}

	// cleant tempDir
	err = os.RemoveAll(n.directory)
	if err != nil {
		return err
	}

	fmt.Println("Temp directory cleaned. Importing genesis file...")

	// import genesis file
	os.Args = []string{
		"sonictool",
		"--datadir", n.directory,
		"genesis", "--experimental", "testGenesis.g",
	}
	err = sonictool.Run()
	if err != nil {
		return err
	}

	// restore original args
	os.Args = originalArgs

	fmt.Println("Genesis file imported. Starting network...")

	// start network again
	return n.start()
}

// GetHeaders returns the headers of all blocks on the network from block 0 to the latest block.
func (n *IntegrationTestNet) GetHeaders() ([]*types.Header, error) {
	client, err := n.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}
	defer client.Close()

	lastBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get last block: %w", err)
	}

	headers := []*types.Header{}
	for i := int64(0); i < int64(lastBlock.NumberU64()); i++ {
		header, err := client.HeaderByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			return nil, fmt.Errorf("failed to get header: %w", err)
		}
		headers = append(headers, header)
	}

	return headers, nil
}

// DeployContract is a utility function handling the deployment of a contract on the network.
// The contract is deployed with by the network's validator account. The function returns the
// deployed contract instance and the transaction receipt.
func DeployContract[T any](n *IntegrationTestNet, deploy contractDeployer[T]) (*T, *types.Receipt, error) {
	client, err := n.GetClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}
	defer client.Close()

	transactOptions, err := n.GetTransactOptions(&n.validator)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get transaction options: %w", err)
	}

	_, transaction, contract, err := deploy(transactOptions, client)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy contract: %w", err)
	}

	receipt, err := n.GetReceipt(transaction.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get receipt: %w", err)
	}
	return contract, receipt, nil
}

// contractDeployer is the type of the deployment functions generated by abigen.
type contractDeployer[T any] func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *T, error)
