package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// TODO: enable test once genesis is fixed
func DisabledTestGenesis_NetworkCanCreateNewBlocksAfterExportImport(t *testing.T) {
	const numBlocks = 3
	require := require.New(t)

	tempDir := t.TempDir()
	net, err := StartIntegrationTestNet(tempDir)
	require.NoError(err)

	// Produce a few blocks on the network.
	for range numBlocks {
		_, err := net.EndowAccount(common.Address{42}, 100)
		require.NoError(err, "failed to endow account")
	}

	// TODO: check for error once genesis is fixed

	// get client
	client, err := net.GetClient()
	require.NoError(err)
	originalHeaders, err := net.GetHeaders()
	require.NoError(err)
	client.Close()

	originalHashes := []common.Hash{}
	for _, header := range originalHeaders {
		originalHashes = append(originalHashes, header.Hash())
	}

	err = net.RestartWithExportImport()
	require.NoError(err)

	// get a fresh client
	newClient, err := net.GetClient()
	require.NoError(err)
	defer newClient.Close()

	// check address 42 has balance
	balance42, err := newClient.BalanceAt(context.Background(), common.Address{42}, nil)
	require.NoError(err)
	require.Equal(0, balance42.Cmp(big.NewInt(100)), "unexpected balance")

	// check headers are consistent with original hashes
	newHeaders, err := net.GetHeaders()
	require.NoError(err)
	require.Equal(len(originalHashes), len(newHeaders), "unexpected number of headers")
	for i := 0; i < len(originalHashes); i++ {
		require.Equal(originalHashes[i], newHeaders[i].Hash(), "unexpected header")
	}

	// Produce a few blocks on the network
	for range numBlocks {
		_, err := net.EndowAccount(common.Address{42}, 100)
		require.NoError(err, "failed to endow account")
	}

	// get headers for all blocks
	allHeaders, err := net.GetHeaders()
	require.NoError(err)

	// check headers from before the export are still reachable
	require.Equal(numBlocks*2+2, len(allHeaders), "unexpected number of headers")

	// TODO: check blocks timestamps
}
