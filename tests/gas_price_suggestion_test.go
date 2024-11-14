package tests

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestGasPrice_GasEvolvesAsExpectedCalculates(t *testing.T) {
	require := require.New(t)

	net, err := StartIntegrationTestNet(t.TempDir())
	require.NoError(err)
	defer net.Stop()

	/*
		// Produce a few blocks on the network.
		for range 10 {
			_, err = net.EndowAccount(common.Address{42}, 100)
			require.NoError(err)
		}
	*/

	client, err := net.GetClient()
	require.NoError(err)
	defer client.Close()

	suggestions := []uint64{}
	prices := []uint64{}

	for i := 0; i < 10; i++ {

		suggestedPrice, err := client.SuggestGasPrice(context.Background())
		require.NoError(err)

		// new block
		receipt, err := net.EndowAccount(common.Address{42}, 100)
		require.NoError(err)

		lastBlock, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
		require.NoError(err)

		suggestions = append(suggestions, suggestedPrice.Uint64())
		prices = append(prices, lastBlock.BaseFee().Uint64())

		diff, ok := within10Percent(suggestedPrice, lastBlock.BaseFee())
		t.Logf("i: %v, last block's base fee (%v) ok:%v, suggested price %v, diff: %v", i, lastBlock.BaseFee(), ok, suggestedPrice, diff)
	}

	for i := range suggestions {
		fmt.Printf("%d, %d, %d\n", i, suggestions[i], prices[i])
	}
}

func within10Percent(a, b *big.Int) (*big.Int, bool) {
	// calculate the difference
	diff := new(big.Int).Sub(a, b)
	diff.Abs(diff)
	// calculate 10% of a
	tenPercent := new(big.Int).Mul(a, big.NewInt(10))
	tenPercent.Div(tenPercent, big.NewInt(100))
	// check if the difference is less than 10% of a
	return diff.Div(diff, a), diff.Cmp(tenPercent) < 0
}
