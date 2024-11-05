package tests

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/require"
)

func TestWithdrawalsCanBeRLPEncodedAndDecoded(t *testing.T) {
	require := require.New(t)

	// start network.
	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoErrorf(err, "Failed to start the fake network: %v", err)
	defer net.Stop()

	// run endowment to ensure at least one block exists
	receipt, err := net.EndowAccount(common.Address{42}, 1)
	require.NoError(err)
	require.Equal(receipt.Status, types.ReceiptStatusSuccessful, "failed to endow account")

	// get client and block
	client, err := net.GetClient()
	require.NoError(err, "Failed to get the client: %v", err)
	defer client.Close()
	block, err := client.BlockByNumber(context.Background(), big.NewInt(0))
	require.NoError(err, "Failed to get the block: %v", err)

	// check that if no withdrawals are made, then the block withdrawals hash is the empty hash
	require.Empty(block.Withdrawals())
	require.Empty(block.Header().WithdrawalsHash)

	// create a new block with an empty list of withdrawals
	newBody := types.Body{Withdrawals: []*types.Withdrawal{}}
	newBlock := types.NewBlock(block.Header(), &newBody, nil, trie.NewStackTrie(nil))

	// encode block
	buffer := bytes.NewBuffer(make([]byte, 0))
	err = newBlock.EncodeRLP(buffer)
	require.NoError(err, "failed to encode block %v", err)

	// decode block
	stream := rlp.NewStream(buffer, 0)
	err = newBlock.DecodeRLP(stream)
	require.NoError(err, "failed to decode block header; %v", err)

	// make empty string hash
	err = rlp.Encode(buffer, "")
	require.NoError(err)
	newHash := common.Hash(crypto.Keccak256(buffer.Bytes()))
	require.Equal(types.EmptyWithdrawalsHash, newHash)

	// check that the block has an empty list of withdrawals
	require.Equal(types.Withdrawals{}, newBlock.Withdrawals())
	require.Equal(newHash, *newBlock.Header().WithdrawalsHash)
}
