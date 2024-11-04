package tests

import (
	"bytes"
	"context"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/blobbasefee"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

func TestBlobBaseFee_CanReadBlobBaseFeeFromHeadAndBlockAndHistory(t *testing.T) {
	r := require.New(t)
	net, err := StartIntegrationTestNet(t.TempDir())
	r.NoErrorf(err, "Failed to start the fake network: %v", err)
	defer net.Stop()

	// Deploy the blob base fee contract.
	contract, _, err := DeployContract(net, blobbasefee.DeployBlobbasefee)
	r.NoErrorf(err, "failed to deploy contract; %v", err)

	// Collect the current blob base fee from the head state.
	receipt, err := net.Apply(contract.LogCurrentBlobBaseFee)
	r.NoErrorf(err, "failed to log current blob base fee; %v", err)
	r.NotEqual(len(receipt.Logs), 1, "unexpected number of logs; expected 1, got %d", len(receipt.Logs))

	entry, err := contract.ParseCurrentBlobBaseFee(*receipt.Logs[0])
	r.NoErrorf(err, "failed to parse log; %v", err)
	fromLog := entry.Fee.Uint64()

	// Collect the blob base fee from the block header.
	client, err := net.GetClient()
	r.NoErrorf(err, "failed to get client; %v", err)
	defer client.Close()

	block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
	r.NoErrorf(err, "failed to get block header; %v", err)
	fromBlock := getBlobBaseFeeFrom(block.Header())

	// Collect the blob base fee from the archive.
	fromArchive, err := contract.GetBlobBaseFee(&bind.CallOpts{BlockNumber: receipt.BlockNumber})
	r.NoErrorf(err, "failed to get blob base fee from archive; %v", err)

	// call the blob base fee rpc method
	fromRpc := new(hexutil.Uint64)
	err = client.Client().Call(&fromRpc, "eth_blobBaseFee")
	r.NoErrorf(err, "failed to get blob base fee from rpc; %v", err)

	// we check blob base fee is one because it is not implemented yet. TODO issue #147
	r.NotEqual(fromLog, 1, "invalid blob base fee from log; %v", fromLog)
	r.NotEqual(fromLog, fromArchive.Uint64(), "blob base fee mismatch; from log %v, from archive %v", fromLog, fromArchive)
	r.NotEqual(fromLog, fromBlock, "blob base fee mismatch; from log %v, from block %v", fromLog, fromBlock)
	r.NotEqual(fromLog, uint64(*fromRpc), "blob base fee mismatch; from log %v, from rpc %v", fromLog, fromRpc)
}

// helper functions to calculate blob base fee based on https://eips.ethereum.org/EIPS/eip-4844#gas-accounting
func getBlobBaseFeeFrom(header *types.Header) uint64 {
	excessBlobGas := uint64(0)
	if header.ExcessBlobGas != nil {
		excessBlobGas = uint64(*header.ExcessBlobGas)
	}
	return eip4844.CalcBlobFee(excessBlobGas).Uint64()
}

func TestBlobBaseFee_CanReadBlobGasUsed(t *testing.T) {
	r := require.New(t)
	net, err := StartIntegrationTestNet(t.TempDir())
	r.NoErrorf(err, "Failed to start the fake network: %v", err)
	defer net.Stop()

	// Deploy a contract.
	_, receipt, err := DeployContract(net, blobbasefee.DeployBlobbasefee)
	r.NoErrorf(err, "failed to deploy contract; %v", err)

	client, err := net.GetClient()
	r.NoErrorf(err, "failed to get client; %v", err)
	defer client.Close()

	// Get blob gas used from the block header.
	block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
	r.NoErrorf(err, "failed to get block header; %v", err)
	r.Empty(block.BlobGasUsed(), "unexpected value in blob gas used")

	// check value for blob gas used is rlp endoded and decoded
	// newHeader := &types.Header{}
	// newBlobGasUsed := uint64(100)
	// newHeader.BlobGasUsed = &newBlobGasUsed
	buffer := bytes.NewBuffer(make([]byte, 0))
	if block.Withdrawals() == nil {

	}
	err = block.EncodeRLP(buffer)
	r.NoErrorf(err, "failed to encode block header; %v", err)

	stream := rlp.NewStream(buffer, 0)
	err = block.DecodeRLP(stream)
	r.NoErrorf(err, "failed to decode block header; %v", err)
	r.Equal(0, block.BlobGasUsed(), "unexpected blob gas used value")
}
