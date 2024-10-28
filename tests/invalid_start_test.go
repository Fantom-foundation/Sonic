package tests

import (
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/invalidstart"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestInvalidStart_IdentifiesInvalidStartContract(t *testing.T) {
	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	// Deploy the invalid start contract.
	contract, _, err := DeployContract(net, invalidstart.DeployInvalidstart)
	if err != nil {
		t.Fatalf("failed to deploy contract; %v", err)
	}

	// attempt to create a contract with invalid code
	receipt, err := net.Apply(contract.CreateWithInvalidCode)
	if err != nil {
		t.Fatalf("failed to log current base fee; %v", err)
	}
	if receipt.Status != types.ReceiptStatusFailed {
		t.Error("contract deployment succeeded unexpectedly")
	}

	// attempt to create a contract with invalid code
	receipt, err = net.Apply(contract.Create2WithInvalidCode)
	if err != nil {
		t.Fatalf("failed to log current base fee; %v", err)
	}
	if receipt.Status != types.ReceiptStatusFailed {
		t.Error("contract deployment succeeded unexpectedly")
	}

	// create an empty contract
	receipt, err = net.Apply(contract.CreateEmptyContractAndTransferToIt)
	if err != nil {
		t.Fatalf("failed to log current base fee; %v", err)
	}
	if receipt.Status != types.ReceiptStatusFailed {
		t.Error("contract deployment succeeded unexpectedly")
	}
}
