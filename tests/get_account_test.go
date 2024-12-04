package tests

import (
	"github.com/Fantom-foundation/go-opera/ethapi"
	"github.com/Fantom-foundation/go-opera/tests/contracts/transientstorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAccount(t *testing.T) {
	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(t, err, "failed to start the fake network")
	defer net.Stop()

	// Deploy the transient storage contract
	_, deployReceipt, err := DeployContract(net, transientstorage.DeployTransientstorage)
	require.NoError(t, err, "failed to deploy contract")

	addr := deployReceipt.ContractAddress

	c, err := net.GetClient()
	require.NoError(t, err, "failed to get client")
	defer c.Close()

	rpcClient := c.Client()
	defer rpcClient.Close()

	var res ethapi.GetAccountResult
	err = rpcClient.Call(&res, "eth_getAccount", addr, rpc.LatestBlockNumber)
	require.NoError(t, err, "failed to call get account")

	// Extract proof to find actual StorageHash(Root), Nonce, Balance and CodeHash
	var proofRes struct {
		StorageHash common.Hash
		Nonce       hexutil.Uint64
		Balance     *hexutil.U256
		CodeHash    common.Hash
	}
	err = rpcClient.Call(
		&proofRes,
		"eth_getProof",
		addr,
		nil,
		rpc.LatestBlockNumber,
	)
	require.NoError(t, err, "failed call to get proof")

	require.Equal(t, proofRes.CodeHash, res.CodeHash)
	require.Equal(t, proofRes.StorageHash, res.StorageRoot)
	require.Equal(t, proofRes.Balance, res.Balance)
	require.Equal(t, proofRes.Nonce, res.Nonce)
}
