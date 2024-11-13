package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestBlockHeader_SatisfyInvariants(t *testing.T) {
	const numBlocks = 5
	require := require.New(t)

	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(err)
	defer net.Stop()

	// Produce a few blocks on the network.
	for range numBlocks {
		_, err := net.EndowAccount(common.Address{42}, 100)
		if err != nil {
			t.Fatalf("failed to endow account; %v", err)
		}
	}

	client, err := net.GetClient()
	require.NoError(err)
	defer client.Close()

	lastBlock, err := client.BlockByNumber(context.Background(), nil)
	require.NoError(err)
	require.GreaterOrEqual(lastBlock.NumberU64(), uint64(numBlocks))

	headers := []*types.Header{}
	for i := int64(0); i < int64(lastBlock.NumberU64()); i++ {
		header, err := client.HeaderByNumber(context.Background(), big.NewInt(i))
		require.NoError(err)
		headers = append(headers, header)
	}

	t.Run("NumberMatchesPosition", func(t *testing.T) {
		testHeaders_NumberMatchesPosition(t, headers)
	})

	t.Run("ParentHashCoversParentContent", func(t *testing.T) {
		testHeaders_ParentHashCoversParentContent(t, headers)
	})

	t.Run("GasUsedIsBelowGasLimit", func(t *testing.T) {
		testHeaders_GasUsedIsBelowGasLimit(t, headers)
	})

	// TODO: Add more tests.
	// - check that the transaction root matches the transactions in the block
	// - check that the receipt root matches the receipts in the block
	// - check that the logs bloom matches the logs in the receipts
	// - coinbase is zero for all blocks
	// - difficulty and nonce is set to 0
	// - time is progressing strictly monotonically and approximately matches the current time
	// - the random mixDigest field is different for each block
}

func testHeaders_NumberMatchesPosition(t *testing.T, headers []*types.Header) {
	require := require.New(t)
	for i, header := range headers {
		require.Equal(header.Number.Uint64(), uint64(i))
	}
}

func testHeaders_ParentHashCoversParentContent(t *testing.T, headers []*types.Header) {
	require := require.New(t)

	// All other blocks have a parent hash that matches the previous block's hash.
	// TODO: fix support for genesis blocks 0 and 1 as well;
	for i := 2; i < len(headers); i++ {
		require.Equal(
			headers[i].ParentHash,
			headers[i-1].Hash(),
		)
	}
}

func testHeaders_GasUsedIsBelowGasLimit(t *testing.T, headers []*types.Header) {
	require := require.New(t)
	for i, header := range headers {
		if i < 2 { // TODO: fix support for genesis blocks 0 and 1 as well;
			continue
		}
		require.LessOrEqual(header.GasUsed, header.GasLimit, "block %d", i)
	}
}
