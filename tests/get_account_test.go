package tests

import (
	"github.com/Fantom-foundation/go-opera/ethapi"
	"github.com/Fantom-foundation/go-opera/tests/contracts/transientstorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/holiman/uint256"
	"testing"
)

func TestGetAccount(t *testing.T) {
	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	// Deploy the transient storage contract
	_, deployReceipt, err := DeployContract(net, transientstorage.DeployTransientstorage)
	if err != nil {
		t.Fatalf("failed to deploy contract; %v", err)
	}

	addr := deployReceipt.ContractAddress

	c, err := net.GetClient()
	if err != nil {
		t.Fatalf("failed to get client; %v", err)
	}

	rpcClient := c.Client()

	var res ethapi.GetAccountResult
	err = rpcClient.Call(&res, "eth_getAccount", addr, rpc.LatestBlockNumber)
	if err != nil {
		t.Fatalf("failed to get account; %v", err)
	}

	if res.CodeHash == (common.Hash{}) {
		t.Error("code hash should not be empty")
	}

	if res.StorageRoot == (common.Hash{}) {
		t.Error("storage root should not be empty")
	}

	if got, want := (*uint256.Int)(res.Balance).Uint64(), uint64(0); got != want {
		t.Errorf("balance not as expected, got: %d want: %d", got, want)
	}

	if res.Nonce < 1 {
		t.Errorf("account nonce is expected to by at least 1")
	}
}
