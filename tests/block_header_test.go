package tests

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	sonictool "github.com/Fantom-foundation/go-opera/cmd/sonictool/app"
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

	tempDir := t.TempDir()
	net, err := StartIntegrationTestNet(tempDir)
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

	runTests := func() {
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

	// Run test three times:
	// - once before a node restart
	// - once after a node restart.
	// - once after a node restart with export/import of the genesis file.

	runTests()
	require.NoError(net.Restart())
	runTests()
	restartWithExportImport(t, net, tempDir)
	runTests()

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
		t.Logf("block %d: %s parent time: %v", i, currentTime, previousTime)
		require.Greater(currentTime, previousTime, "time is not monotonically increasing. block %d", i)
	}
}

func testHeaders_MixDigestDiffersForAllBlocks(t *testing.T, headers []*types.Header) {
	require := require.New(t)

	seen := map[common.Hash]struct{}{}

	for i := 1; i < len(headers); i++ {
		_, ok := seen[headers[i].MixDigest]
		require.False(ok, "mix digest is not unique")
		seen[headers[i].MixDigest] = struct{}{}
	}
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

func restartWithExportImport(t *testing.T, net *IntegrationTestNet, tempDir string) {
	require := require.New(t)

	net.Stop()
	t.Log("Network stopped. Exporting genesis file...")

	// export
	os.Args = []string{
		"sonictool",
		"--datadir", tempDir,
		"genesis", "export", "testGenesis.g",
	}
	err := sonictool.Run()
	require.NoError(err, "failed to import genesis file")

	// cleant tempDir
	err = os.RemoveAll(tempDir)
	require.NoError(err, "failed to delete temp directory contents")

	t.Log("Temp directory cleaned. Importing genesis file...")

	// import genesis file
	os.Args = []string{
		"sonictool",
		"--datadir", tempDir,
		"genesis", "--experimental", "testGenesis.g",
	}
	err = sonictool.Run()
	require.NoError(err, "failed to import genesis file")

	t.Log("Genesis file imported. Starting network...")

	// start network again
	err = net.start()
	require.NoError(err)
}
