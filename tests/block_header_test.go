package tests

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip/gasprice"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/Fantom-foundation/go-opera/opera/contracts/driver"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/require"
)

func TestBlockHeader_SatisfiesInvariants(t *testing.T) {
	const numBlocks = 10
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

	originalHeaders, err := net.GetHeaders()
	require.NoError(err)
	originalHashes := []common.Hash{}
	for _, header := range originalHeaders {
		originalHashes = append(originalHashes, header.Hash())
	}

	// Run twice - once before and once after a node restart.
	runTests := func() {
		headers, err := net.GetHeaders()
		require.NoError(err)

		t.Run("CompareHeadersHashes", func(t *testing.T) {
			testHeaders_CompareHeadersHashes(t, originalHashes, headers)
		})

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

		t.Run("TransactionRootMatchesHashOfBlockTxs", func(t *testing.T) {
			testHeaders_TransactionRootMatchesBlockTxsHash(t, headers, client)
		})

		t.Run("ReceiptRootMatchesBlockReceipts", func(t *testing.T) {
			testHeaders_ReceiptRootMatchesBlockReceipts(t, headers, client)
		})

		t.Run("LogsBloomMatchesLogsInReceipts", func(t *testing.T) {
			testHeaders_LogsBloomMatchesLogsInReceipts(t, headers, client)
		})

		t.Run("CoinbaseIsZeroForAllBlocks", func(t *testing.T) {
			testHeaders_CoinbaseIsZeroForAllBlocks(t, headers)
		})

		t.Run("DifficultyIsZeroForAllBlocks", func(t *testing.T) {
			testHeaders_DifficultyIsZeroForAllBlocks(t, headers)
		})

		t.Run("NonceIsZeroForAllBlocks", func(t *testing.T) {
			testHeaders_NonceIsZeroForAllBlocks(t, headers)
		})

		t.Run("TimeProgressesMonotonically", func(t *testing.T) {
			testHeaders_TimeProgressesMonotonically(t, headers)
		})

		t.Run("MixDigestDiffersForAllBlocks", func(t *testing.T) {
			testHeaders_MixDigestDiffersForAllBlocks(t, headers)
		})

		t.Run("LastBlockOfEpochContainsSealingTransaction", func(t *testing.T) {
			testHeaders_LastBlockOfEpochContainsSealingTransaction(t, headers, client)
		})
	}

	runTests()
	require.NoError(net.Restart())
	runTests()

	// TODO: enable when genesis is fixed
	// require.NoError(net.restartWithExportImport())
	// runTests()
}

func testHeaders_CompareHeadersHashes(t *testing.T, hashes []common.Hash, newHeaders []*types.Header) {
	require := require.New(t)

	require.Len(newHeaders, len(hashes), "length mismatch")
	for i, header := range newHeaders {
		require.Equal(hashes[i], header.Hash(), "hash mismatch")
	}
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

	getUnixTime := func(header *types.Header) time.Time {
		t.Helper()
		nanos, _, err := inter.DecodeExtraData(header.Extra)
		require.NoError(err)
		return time.Unix(int64(header.Time), int64(nanos))
	}

	// Check the nano-time and duration encoded in the extra data field.
	for i := 1; i < len(headers); i++ {
		require.Equal(len(headers[i].Extra), 12, "extra data length of block %d", i)
		lastTime := getUnixTime(headers[i-1])
		currentTime := getUnixTime(headers[i])
		wantedDuration := currentTime.Sub(lastTime)
		_, gotDuration, err := inter.DecodeExtraData(headers[i].Extra)
		require.NoError(err, "decoding extra data of block %d", i)
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
		_, duration, err := inter.DecodeExtraData(headers[i-1].Extra)
		require.NoError(err, "decoding extra data of block %d", i-1)
		last := &evmcore.EvmHeader{
			BaseFee:  headers[i-1].BaseFee,
			GasUsed:  headers[i-1].GasUsed,
			Duration: duration,
		}
		require.Equal(
			gasprice.GetBaseFeeForNextBlock(last, rules),
			headers[i].BaseFee,
			"base fee of block %d", i,
		)
	}
}

func testHeaders_TransactionRootMatchesBlockTxsHash(t *testing.T, headers []*types.Header, client *ethclient.Client) {
	require := require.New(t)

	for _, header := range headers {
		block, err := client.BlockByNumber(context.Background(), header.Number)
		require.NoError(err, "failed to get block receipts")

		txsHash := types.DeriveSha(block.Transactions(), trie.NewStackTrie(nil))
		require.Equal(header.TxHash, txsHash, "transaction root hash mismatch")
	}
}

func testHeaders_ReceiptRootMatchesBlockReceipts(t *testing.T, headers []*types.Header, client *ethclient.Client) {
	require := require.New(t)

	for _, header := range headers {
		receipts, err := client.BlockReceipts(context.Background(),
			rpc.BlockNumberOrHashWithHash(header.Hash(), false))
		require.NoError(err, "failed to get block receipts")

		receiptsHash := types.DeriveSha(types.Receipts(receipts), trie.NewStackTrie(nil))
		require.Equal(header.ReceiptHash, receiptsHash, "receipt root hash mismatch")
	}
}

func testHeaders_LogsBloomMatchesLogsInReceipts(t *testing.T, headers []*types.Header, client *ethclient.Client) {
	require := require.New(t)

	for _, header := range headers {
		receipts, err := client.BlockReceipts(context.Background(),
			rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(header.Number.Uint64())))
		require.NoError(err, "failed to get block receipts")

		logsBloom := types.CreateBloom(receipts)
		require.Equal(header.Bloom, logsBloom, "logs bloom mismatch")
	}
}

func testHeaders_CoinbaseIsZeroForAllBlocks(t *testing.T, headers []*types.Header) {
	require := require.New(t)

	for _, header := range headers {
		require.Zero(header.Coinbase, "coinbase is not zero")
	}
}

func testHeaders_DifficultyIsZeroForAllBlocks(t *testing.T, headers []*types.Header) {
	require := require.New(t)

	for _, header := range headers {
		// Cmp returns 0 when the values are equal
		require.Zero(big.NewInt(0).Cmp(header.Difficulty), "difficulty is not zero")
	}
}

func testHeaders_NonceIsZeroForAllBlocks(t *testing.T, headers []*types.Header) {
	require := require.New(t)

	for _, header := range headers {
		require.Zero(header.Nonce.Uint64(), "nonce is not zero")
	}
}

func testHeaders_TimeProgressesMonotonically(t *testing.T, headers []*types.Header) {
	require := require.New(t)

	getTimeFrom := func(header *types.Header) time.Time {
		currentNano, _, err := inter.DecodeExtraData(header.Extra)
		require.NoError(err)
		return time.Unix(int64(header.Time), int64(currentNano))
	}

	for i := 1; i < len(headers); i++ {
		currentTime := getTimeFrom(headers[i])
		previousTime := getTimeFrom(headers[i-1])
		// TODO: enable to print timestamps of blocks
		// t.Logf("block %d: %s parent time: %v", i, currentTime, previousTime)
		require.Greater(currentTime, previousTime, "time is not monotonically increasing. block %d", i)
	}
}

func testHeaders_MixDigestDiffersForAllBlocks(t *testing.T, headers []*types.Header) {
	require := require.New(t)

	seen := map[common.Hash]struct{}{}

	for i := 1; i < len(headers); i++ {
		// We skip empty blocks, since in those cases the MixDigest value is not
		// consumed by any transaction. For those cases, values may be reused.
		// Since the prev-randao value filling this field is computed based on
		// the hash of non-empty lachesis events, the value used for empty blocks
		// is always the same.
		header := headers[i]
		if header.GasUsed == 0 {
			continue
		}
		digest := header.MixDigest
		_, found := seen[digest]
		require.False(found, "mix digest is not unique, block %d, value %x", i, digest)
		seen[digest] = struct{}{}
	}
	require.NotZero(len(seen), "no non-empty blocks in the chain")
}

func testHeaders_LastBlockOfEpochContainsSealingTransaction(t *testing.T, headers []*types.Header, client *ethclient.Client) {
	require := require.New(t)

	maxEpoch := 0
	for i := 0; i < len(headers)-1; i++ {

		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		require.NoError(err, "failed to get block body")

		currentBlockEpoch, err := getEpochOfBlock(client, i)
		require.NoError(err, "failed to get epoch of block %d", i)

		nextBlockEpoch, err := getEpochOfBlock(client, i+1)
		require.NoError(err, "failed to get epoch of block %d", i+1)

		shouldContainSealingTx := currentBlockEpoch != nextBlockEpoch

		containsSealingTx := false
		for _, tx := range block.Transactions() {
			// Any call to the driver is considered a sealing transaction.
			if tx.To() != nil && *tx.To() == driver.ContractAddress {
				containsSealingTx = true
				break
			}
		}

		require.Equal(
			shouldContainSealingTx,
			containsSealingTx,
			"block %d, current epoch %d, next epoch %d",
			i, currentBlockEpoch, nextBlockEpoch,
		)

		if currentBlockEpoch > maxEpoch {
			maxEpoch = currentBlockEpoch
		}
	}
}

func getEpochOfBlock(client *ethclient.Client, blockNumber int) (int, error) {
	var result struct {
		Epoch hexutil.Uint64
	}
	err := client.Client().Call(
		&result,
		"eth_getBlockByNumber",
		fmt.Sprintf("0x%x", blockNumber),
		false,
	)
	if err != nil {
		return 0, err
	}
	return int(result.Epoch), nil
}
