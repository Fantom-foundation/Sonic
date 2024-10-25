package tests

import (
	"context"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/blobbasefee"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
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
		t.Errorf("failed to get blob base fee from archive; %v", err)
	}

	// call the blob base fee rpc method
	fromRpc := new(hexutil.Uint64)
	err = client.Client().Call(&fromRpc, "eth_blobBaseFee")
	if err != nil {
		t.Fatalf("failed to get blob base fee from rpc; %v", err)
	}

	// we check blob base fee is zero because it is not implemented yet. TODO issue #147
	if fromLog != 1 {
		t.Errorf("invalid blob base fee from log; %v", fromLog)
	}

	if fromLog != fromArchive.Uint64() {
		t.Errorf("blob base fee mismatch; from log %v, from archive %v", fromLog, fromArchive)
	}

	if fromLog != fromBlock {
		t.Errorf("blob base fee mismatch; from log %v, from block %v", fromLog, fromBlock)
	}

	if fromLog != uint64(*fromRpc) {
		t.Errorf("blob base fee mismatch; from log %v, from rpc %v", fromLog, fromRpc)
	}
}

// helper functions to calculate blob base fee based on https://eips.ethereum.org/EIPS/eip-4844#gas-accounting
func getBlobBaseFeeFrom(header *types.Header) uint64 {
	excessBlobGas := uint64(0)
	if header.ExcessBlobGas != nil {
		excessBlobGas = uint64(*header.ExcessBlobGas)
	}
	return eip4844.CalcBlobFee(excessBlobGas).Uint64()
}
