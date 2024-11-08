package tests

import (
	"bytes"
	"context"
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

	// get client and block
	client, err := net.GetClient()
	requireBase.NoError(err, "Failed to get the client: ", err)
	defer client.Close()
	block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
	requireBase.NoError(err, "Failed to get the block: ", err)

	t.Run("verify default values of block's Withdrawals list and hash", func(t *testing.T) {
		require := require.New(t)
		// check that if no withdrawals are made, then the block withdrawals hash is the empty hash
		require.Equal(types.EmptyWithdrawalsHash, *block.Header().WithdrawalsHash)
		require.Empty(block.Withdrawals())
	})

	t.Run("encode and decode works properly", func(t *testing.T) {
		require := require.New(t)
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
