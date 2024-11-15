package tests

import (
	"context"
	"encoding/binary"
	"math/big"
	"testing"
	"time"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip/gasprice"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestBlockHeader_SatisfiesInvariants(t *testing.T) {
	const numBlocks = 5
	require := require.New(t)

	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(err)
	defer net.Stop()

	// Produce a few blocks on the network.
	for range numBlocks {
		_, err := net.EndowAccount(common.Address{42}, 100)
		require.NoError(err, "failed to endow account")
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

	t.Run("BlockNumberEqualsPositionInChain", func(t *testing.T) {
		testHeaders_BlockNumberEqualsPositionInChain(t, headers)
	})

	t.Run("ParentHashCoversParentContent", func(t *testing.T) {
		testHeaders_ParentHashCoversParentContent(t, headers)
	})

	t.Run("GasUsedIsBelowGasLimit", func(t *testing.T) {
		testHeaders_GasUsedIsBelowGasLimit(t, headers)
	})

	t.Run("EncodesDurationAndNanoTimeInExtraData", func(t *testing.T) {
		testHeaders_EncodesDurationAndNanoTimeInExtraData(t, headers)
	})

	t.Run("BaseFeeEvolutionFollowsPricingRules", func(t *testing.T) {
		testHeaders_BaseFeeEvolutionFollowsPricingRules(t, headers)
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

func testHeaders_BlockNumberEqualsPositionInChain(t *testing.T, headers []*types.Header) {
	require := require.New(t)
	for i, header := range headers {
		require.Equal(header.Number.Uint64(), uint64(i))
	}
}

func testHeaders_ParentHashCoversParentContent(t *testing.T, headers []*types.Header) {
	require := require.New(t)

	// The parent hash of block 0 is expected to be zero.
	require.Equal(
		headers[0].ParentHash, common.Hash{},
		"invalid parent hash for block 0",
	)

	// All other blocks have a parent hash that matches the previous block's hash.
	for i := 1; i < len(headers); i++ {
		require.Equal(
			headers[i].ParentHash,
			headers[i-1].Hash(),
			"invalid hash stored in block %d for block %d", i, i-1,
		)
	}
}

func testHeaders_GasUsedIsBelowGasLimit(t *testing.T, headers []*types.Header) {
	require := require.New(t)
	for i, header := range headers {
		require.LessOrEqual(header.GasUsed, header.GasLimit, "block %d", i)
	}
}

func testHeaders_EncodesDurationAndNanoTimeInExtraData(t *testing.T, headers []*types.Header) {
	require := require.New(t)
	// Check the nano-time and duration encoded in the extra data field.
	for i := 1; i < len(headers); i++ {
		require.Equal(len(headers[i].Extra), 16, "extra data length of block %d", i)
		lastTime := binary.BigEndian.Uint64(headers[i-1].Extra[:8])
		currentTime := binary.BigEndian.Uint64(headers[i].Extra[:8])
		wantedDuration := currentTime - lastTime
		gotDuration := binary.BigEndian.Uint64(headers[i].Extra[8:])
		require.Equal(wantedDuration, gotDuration, "duration of block %d", i)
	}
}

func testHeaders_BaseFeeEvolutionFollowsPricingRules(t *testing.T, headers []*types.Header) {
	require := require.New(t)

	// The genesis block must use the initial base fee.
	rules := opera.FakeEconomyRules()
	require.Equal(
		gasprice.GetInitialBaseFee(rules),
		headers[0].BaseFee,
	)

	// All other blocks compute the base-fee based on the previous block.
	for i := 1; i < len(headers); i++ {
		last := &evmcore.EvmHeader{
			BaseFee:  headers[i-1].BaseFee,
			GasLimit: headers[i-1].GasLimit,
			GasUsed:  headers[i-1].GasUsed,
			Duration: time.Duration(
				binary.BigEndian.Uint64(headers[i-1].Extra[8:]),
			),
		}
		require.Equal(
			gasprice.GetBaseFeeForNextBlock(last, rules),
			headers[i].BaseFee,
			"base fee of block %d", i,
		)
	}
}
