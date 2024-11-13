package tests

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

func TestWithdrawalFieldsInBlocks(t *testing.T) {
	requireBase := require.New(t)

	// start network.
	net, err := StartIntegrationTestNet(t.TempDir())
	requireBase.NoErrorf(err, "Failed to start the fake network: ", err)
	defer net.Stop()

	// run endowment to ensure at least one block exists
	receipt, err := net.EndowAccount(common.Address{42}, 1)
	requireBase.NoError(err)
	requireBase.Equal(receipt.Status, types.ReceiptStatusSuccessful, "failed to endow account")

	// get client
	client, err := net.GetClient()
	requireBase.NoError(err, "Failed to get the client: ", err)
	defer client.Close()

	t.Run("verify default values of block's Withdrawals list and hash", func(t *testing.T) {
		require := require.New(t)

		latest, err := client.BlockNumber(context.Background())
		require.NoError(err, "Failed to get the latest block number: ", err)

		// we check from block 1 onward because block 0 does not consider Sonic Upgrade.
		for i := int64(1); i <= int64(latest); i++ {
			block, err := client.BlockByNumber(context.Background(), big.NewInt(i))
			require.NoError(err, "Failed to get the block: ", err)

			// check that the block has an empty list of withdrawals
			require.Empty(block.Withdrawals())
			require.Equal(types.EmptyWithdrawalsHash, *block.Header().WithdrawalsHash, "block %d", i)
		}
	})

	t.Run("blocks are healthy to be RLP encoded and decoded", func(t *testing.T) {
		require := require.New(t)

		// get block
		block, err := client.BlockByNumber(context.Background(), nil)
		requireBase.NoError(err, "Failed to get the block: ", err)

		// encode block
		buffer := bytes.NewBuffer(make([]byte, 0))
		err = block.EncodeRLP(buffer)
		require.NoError(err, "failed to encode block ", err)

		// decode block
		stream := rlp.NewStream(buffer, 0)
		err = block.DecodeRLP(stream)
		require.NoError(err, "failed to decode block header; ", err)

		// check that the block has an empty list of withdrawals
		require.Empty(block.Withdrawals())
		require.Equal(types.EmptyWithdrawalsHash, *block.Header().WithdrawalsHash)
	})
}
