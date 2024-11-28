package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/block_hash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	req "github.com/stretchr/testify/require"
)

func TestBlockHash_CorrectBlockHashesAreAccessibleInContracts(t *testing.T) {
	require := req.New(t)
	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	// Deploy the block hash observer contract.
	_, receipt, err := DeployContract(net, block_hash.DeployBlockHash)
	require.NoError(err, "failed to deploy contract; %v", err)
	contractAddress := receipt.ContractAddress
	contractCreationBlock := receipt.BlockNumber.Uint64()

	runTest := func(t *testing.T) {
		t.Run("visible block hash on head", func(t *testing.T) {
			testVisibleBlockHashOnHead(t, net, contractAddress)
		})
		t.Run("visible block hash in archive", func(t *testing.T) {
			testVisibleBlockHashesInArchive(t, net, contractAddress, contractCreationBlock)
		})
	}

	t.Run("fresh network", runTest)
	net.Restart()
	t.Run("restarted network", runTest)
	net.RestartWithExportImport()
	t.Run("reinitialized network", runTest)
}

func testVisibleBlockHashOnHead(
	t *testing.T,
	net *IntegrationTestNet,
	observerContractAddress common.Address,
) {
	require := req.New(t)
	client, err := net.GetClient()
	require.NoError(err, "failed to get client; %v", err)
	defer client.Close()

	contract, err := block_hash.NewBlockHash(observerContractAddress, client)
	require.NoError(err, "failed to instantiate contract")

	for range 3 {
		receipt, err := net.Apply(contract.Observe)
		require.NoError(err, "failed to observe block hash; %v", err)
		require.Equal(types.ReceiptStatusSuccessful, receipt.Status,
			"failed to observe block hash; %v", err,
		)

		blockNumber := receipt.BlockNumber.Uint64()
		require.Len(receipt.Logs, int(blockNumber+6), "unexpected number of logs")

		for _, log := range receipt.Logs {
			entry, err := contract.ParseSeen(*log)
			require.NoError(err, "failed to parse log; %v", err)
			current := entry.CurrentBlock.Uint64()
			require.Equal(blockNumber, current, "unexpected block number")
			observed := entry.ObservedBlock.Uint64()
			seen := common.Hash(entry.BlockHash)

			want := common.Hash{}
			if observed < current {
				hash, err := client.BlockByNumber(context.Background(), entry.ObservedBlock)
				require.NoError(err, "failed to get block hash; %v", err)
				want = hash.Hash()
			}
			require.Equal(want, seen, "block hash mismatch, current: %d, observed: %d", current, observed)
		}
	}
}

func testVisibleBlockHashesInArchive(
	t *testing.T,
	net *IntegrationTestNet,
	observerContractAddress common.Address,
	observerCreationBlock uint64,
) {
	require := req.New(t)
	client, err := net.GetClient()
	require.NoError(err, "failed to get client; %v", err)
	defer client.Close()

	ctxt := context.Background()

	// Get list of all block hashes.
	numBlocks, err := client.BlockNumber(ctxt)
	require.NoError(err, "failed to get block number; %v", err)

	hashes := []common.Hash{}
	for i := uint64(0); i <= numBlocks; i++ {
		hash, err := client.BlockByNumber(ctxt, big.NewInt(int64(i)))
		require.NoError(err, "failed to get block hash; %v", err)
		hashes = append(hashes, hash.Hash())
	}

	// Check that blocks are reported correctly by archive queries.
	numChecks := 0
	observer, err := block_hash.NewBlockHash(observerContractAddress, client)
	require.NoError(err, "failed to instantiate contract")
	for observationBlock := range numBlocks {
		if observationBlock < observerCreationBlock {
			continue
		}
		for observedBlock := range numBlocks {
			hash, err := observer.GetBlockHash(&bind.CallOpts{
				BlockNumber: big.NewInt(int64(observationBlock)),
			}, big.NewInt(int64(observedBlock)))
			require.NoError(err, "failed to get block hash; %v", err)

			want := common.Hash{}
			if observedBlock < observationBlock {
				want = hashes[observedBlock]
			}
			got := common.Hash(hash)
			require.Equal(want, got, "block hash mismatch, observation: %d, observed: %d", observationBlock, observedBlock)
			numChecks++
		}
	}
	require.Greater(numChecks, 0, "no checks performed")
}
