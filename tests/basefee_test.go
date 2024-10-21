package tests

import (
	"fmt"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/basefee"
)

func TestBaseFee_Install(t *testing.T) {
	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	// Deploy the base fee contract.
	contract, _, err := DeployContract(net, basefee.DeployBasefee)
	if err != nil {
		t.Fatalf("failed to deploy contract; %v", err)
	}

	fmt.Printf("Contract deployed successfully\n")

	fee, err := contract.GetBaseFee(nil)
	if err != nil {
		t.Fatalf("failed to get base fee; %v", err)
	}

	fmt.Printf("Base fee from archive: %v\n", fee)

	// Create an account to send test queries.
	account := NewAccount()
	if err := net.EndowAccount(account.Address(), 1e18); err != nil {
		t.Fatalf("failed to endow account; %v", err)
	}

	txOpts, err := net.GetTransactOptions(account)
	if err != nil {
		t.Fatalf("failed to get transact options; %v", err)
	}
	tx, err := contract.LogCurrentBaseFee(txOpts)
	if err != nil {
		t.Fatalf("failed to get current base fee; %v", err)
	}

	receipt, err := net.GetReceipt(tx.Hash())
	if err != nil {
		t.Fatalf("failed to get receipt; %v", err)
	}

	for i, log := range receipt.Logs {
		entry, err := contract.ParseCurrentFee(*log)
		if err != nil {
			t.Fatalf("failed to parse log; %v", err)
		}
		fmt.Printf("Log %d: %v\n", i, entry.Fee)
	}

	t.Fail()
}
