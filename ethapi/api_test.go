package ethapi

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

// feeHistoryBackend implements only the Backend methods used by FeeHistory.
type feeHistoryBackend struct {
	Backend
}

func (feeHistoryBackend) ResolveRpcBlockNumberOrHash(_ context.Context, b rpc.BlockNumberOrHash) (idx.Block, error) {
	return idx.Block(*b.BlockNumber), nil
}

func (feeHistoryBackend) MinGasPrice() *big.Int { return big.NewInt(7) }

func (feeHistoryBackend) SuggestGasTipCap(_ context.Context, certainty uint64) *big.Int {
	return new(big.Int).SetUint64(certainty)
}

func TestFeeHistory_ReturnsRequestedBlockRange(t *testing.T) {
	tests := map[string]struct {
		blockCount rpc.DecimalOrHex
		lastBlock  rpc.BlockNumber
		wantOldest uint64
		wantLen    int
	}{
		"single block":                       {blockCount: 1, lastBlock: 10, wantOldest: 10, wantLen: 1},
		"range ending at last block":         {blockCount: 5, lastBlock: 10, wantOldest: 6, wantLen: 5},
		"range starting right after genesis": {blockCount: 5, lastBlock: 5, wantOldest: 1, wantLen: 5},
		"range cut at genesis":               {blockCount: 100, lastBlock: 10, wantOldest: 0, wantLen: 11},
		"huge block count clamped to 1024":   {blockCount: math.MaxUint64, lastBlock: 2000, wantOldest: 977, wantLen: 1024},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			api := NewPublicEthereumAPI(feeHistoryBackend{})
			res, err := api.FeeHistory(context.Background(), test.blockCount, test.lastBlock, []float64{50})
			require.NoError(t, err)
			require.Equal(t, test.wantOldest, res.OldestBlock.ToInt().Uint64())
			require.Len(t, res.Reward, test.wantLen)
			require.Len(t, res.GasUsedRatio, test.wantLen)
			// base fee carries one extra entry for the next block
			require.Len(t, res.BaseFee, test.wantLen+1)
		})
	}
}

func TestFeeHistory_RejectsTooManyPercentiles(t *testing.T) {
	api := NewPublicEthereumAPI(feeHistoryBackend{})
	_, err := api.FeeHistory(context.Background(), math.MaxUint64, rpc.LatestBlockNumber, make([]float64, maxRewardPercentiles+1))
	require.ErrorIs(t, err, errInvalidPercentile)
}
