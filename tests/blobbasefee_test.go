package tests

import (
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/blobbasefee"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func TestBlobBaseFee_CanReadBlobBaseFeeFromHeadAndBlockAndHistory(t *testing.T) {
	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	// Deploy the base fee contract.
	contract, _, err := DeployContract(net, blobbasefee.DeployBlobbasefee)
	if err != nil {
		t.Fatalf("failed to deploy contract; %v", err)
	}

	// Collect the current base fee from the head state.
	receipt, err := net.Apply(contract.LogCurrentBlobBaseFee)
	if err != nil {
		t.Fatalf("failed to log current base fee; %v", err)
	}

	if len(receipt.Logs) != 1 {
		t.Fatalf("unexpected number of logs; expected 1, got %d", len(receipt.Logs))
	}

	entry, err := contract.ParseCurrentBlobBaseFee(*receipt.Logs[0])
	if err != nil {
		t.Fatalf("failed to parse log; %v", err)
	}
	fromLog := entry.Fee

	// Collect the base fee from the archive.
	fromArchive, err := contract.GetBlobBaseFee(&bind.CallOpts{BlockNumber: receipt.BlockNumber})
	if err != nil {
		t.Fatalf("failed to get base fee from archive; %v", err)
	}

	// we check blob base fee is zero because it is not implemented yet
	if fromLog.Sign() != 0 {
		t.Fatalf("invalid base fee from log; %v", fromLog)
	}

	if fromLog.Cmp(fromArchive) != 0 {
		t.Fatalf("base fee mismatch; from log %v, from archive %v", fromLog, fromArchive)
	}
}
