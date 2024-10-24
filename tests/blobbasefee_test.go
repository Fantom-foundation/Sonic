package tests

import (
	"context"
	"math"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/blobbasefee"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestBlobBaseFee_CanReadBlobBaseFeeFromHeadAndBlockAndHistory(t *testing.T) {
	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	// Deploy the blob base fee contract.
	contract, _, err := DeployContract(net, blobbasefee.DeployBlobbasefee)
	if err != nil {
		t.Fatalf("failed to deploy contract; %v", err)
	}

	// Collect the current blob base fee from the head state.
	receipt, err := net.Apply(contract.LogCurrentBlobBaseFee)
	if err != nil {
		t.Fatalf("failed to log current blob base fee; %v", err)
	}

	if len(receipt.Logs) != 1 {
		t.Fatalf("unexpected number of logs; expected 1, got %d", len(receipt.Logs))
	}

	entry, err := contract.ParseCurrentBlobBaseFee(*receipt.Logs[0])
	if err != nil {
		t.Fatalf("failed to parse log; %v", err)
	}
	fromLog := entry.Fee.Uint64()

	// Collect the blob base fee from the block header.
	client, err := net.GetClient()
	if err != nil {
		t.Fatalf("failed to get client; %v", err)
	}
	defer client.Close()
	block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		t.Fatalf("failed to get block header; %v", err)
	}
	fromBlock := getBlobBaseFeeFrom(block.Header())

	// Collect the blob base fee from the archive.
	fromArchive, err := contract.GetBlobBaseFee(&bind.CallOpts{BlockNumber: receipt.BlockNumber})
	if err != nil {
		t.Fatalf("failed to get blob base fee from archive; %v", err)
	}

	// we check blob base fee is zero because it is not implemented yet
	if fromLog != 0 {
		t.Fatalf("invalid blob base fee from log; %v", fromLog)
	}

	if fromLog != fromArchive.Uint64() {
		t.Fatalf("blob base fee mismatch; from log %v, from archive %v", fromLog, fromArchive)
	}

	if fromLog != fromBlock {
		t.Fatalf("blob base fee mismatch; from log %v, from block %v", fromLog, fromBlock)
	}
}

// helper functions to calculate blob base fee based on https://eips.ethereum.org/EIPS/eip-4844#gas-accounting
func getBlobBaseFeeFrom(header *types.Header) uint64 {
	blobGasUsed := uint64(0)
	if header.BlobGasUsed != nil {
		blobGasUsed = *header.BlobGasUsed
	}
	excessBlobGas := float64(0)
	if header.ExcessBlobGas != nil {
		excessBlobGas = float64(*header.ExcessBlobGas)
	}
	// source for constants: https://eips.ethereum.org/EIPS/eip-4844#parameters
	const MIN_FEE_PER_BLOB_GAS = float64(1)
	const UPDATE_FRACTION = float64(3338477)
	return blobGasUsed * uint64(MIN_FEE_PER_BLOB_GAS*math.Exp(excessBlobGas/UPDATE_FRACTION))
}
