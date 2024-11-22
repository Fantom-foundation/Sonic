package tests

import (
	"context"
	"math/big"
	"os"
	"testing"

	sonictool "github.com/Fantom-foundation/go-opera/cmd/sonictool/app"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func writeToStdIn(t *testing.T, input string) {
	t.Helper()
	// Open the file to use as stdin
	file, err := os.Create(t.TempDir() + "/input.txt")
	if err != nil {
		t.Error("Error opening file:", err)
	}
	file.Write([]byte(input))
	t.Cleanup(func() { file.Close() })

	originalStdin := os.Stdin
	t.Cleanup(func() { os.Stdin = originalStdin })
	// Redirect stdin
	file.Seek(0, 0)
	os.Stdin = file
}

func TestGenesisExportImport(t *testing.T) {
	const numBlocks = 5
	require := require.New(t)

	tempDir := t.TempDir()
	net, err := StartIntegrationTestNet(tempDir)
	require.NoError(err)
	// defer net.Stop()

	// Produce a few blocks on the network.
	for range numBlocks {
		_, err := net.EndowAccount(common.Address{42}, 100)
		require.NoError(err, "failed to endow account")
	}

	net.Stop()
	t.Log("Network stopped. Exporting genesis file...")

	// export
	os.Args = []string{
		"sonictool",
		"--datadir", tempDir,
		"genesis", "export", "testGenesis.g",
	}
	err = sonictool.Run()
	require.NoError(err, "failed to import genesis file")

	// delete contents of tempDir and re-import the genesis file
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

	err = net.start()
	require.NoError(err)
	defer net.Stop()

	t.Log("Network started. Checking if blocks are still there...")

	// TODO: check address 32 has balance

	// TODO: check blocks timestamps

	// Produce a few blocks on the network.
	for range numBlocks {
		_, err := net.EndowAccount(common.Address{42}, 100)
		require.NoError(err, "failed to endow account")
	}

	// get client
	client, err := net.GetClient()
	require.NoError(err)
	defer client.Close()

	// get headers for all blocks.
	lastBlock, err := client.BlockByNumber(context.Background(), nil)
	require.NoError(err)
	require.GreaterOrEqual(lastBlock.NumberU64(), uint64(numBlocks))

	headers := []*types.Header{}
	for i := int64(0); i < int64(lastBlock.NumberU64()); i++ {
		header, err := client.HeaderByNumber(context.Background(), big.NewInt(i))
		require.NoError(err)
		headers = append(headers, header)
	}

	require.Equal(numBlocks*2+2, len(headers), "unexpected number of headers")
}
