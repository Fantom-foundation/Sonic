package ethapi

import (
	"context"
	"errors"
	"testing"

	"github.com/Fantom-foundation/Carmen/go/common/amount"
	"github.com/Fantom-foundation/go-opera/gossip/evmstore"
	"github.com/stretchr/testify/require"
)

type fakeTotalSupplyProvider struct {
	snapshot *evmstore.TotalSupplySnapshot
	err      error
}

func (f fakeTotalSupplyProvider) GetTotalSupply() (*evmstore.TotalSupplySnapshot, error) {
	return f.snapshot, f.err
}

func TestTotalSupply_ReturnsValuesOfTheSameBlock(t *testing.T) {
	api := &PublicSonicAPI{fakeTotalSupplyProvider{snapshot: &evmstore.TotalSupplySnapshot{
		Block:       42,
		Total:       amount.New(1000),
		ZeroBalance: amount.New(30),
		DeadBalance: amount.New(12),
	}}}

	res, err := api.TotalSupply(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 42, res.BlockNumber)
	require.EqualValues(t, 1000, res.TotalSupply.ToInt().Uint64())
	require.EqualValues(t, 958, res.TotalSupplyWithoutBurned.ToInt().Uint64())
	require.EqualValues(t, 30, res.BurnedZero.ToInt().Uint64())
	require.EqualValues(t, 12, res.BurnedDead.ToInt().Uint64())
}

func TestTotalSupply_BackendError_IsForwarded(t *testing.T) {
	issue := errors.New("no snapshot")
	api := &PublicSonicAPI{fakeTotalSupplyProvider{err: issue}}

	_, err := api.TotalSupply(context.Background())
	require.ErrorIs(t, err, issue)
}

func TestTotalSupply_BurnedExceedsTotal_ReturnsError(t *testing.T) {
	api := &PublicSonicAPI{fakeTotalSupplyProvider{snapshot: &evmstore.TotalSupplySnapshot{
		Total:       amount.New(10),
		ZeroBalance: amount.New(8),
		DeadBalance: amount.New(8),
	}}}

	_, err := api.TotalSupply(context.Background())
	require.ErrorContains(t, err, "exceeds the total supply")
}

func TestTotalSupply_BurnedOverflow_ReturnsError(t *testing.T) {
	api := &PublicSonicAPI{fakeTotalSupplyProvider{snapshot: &evmstore.TotalSupplySnapshot{
		Total:       amount.Max(),
		ZeroBalance: amount.Max(),
		DeadBalance: amount.New(1),
	}}}

	_, err := api.TotalSupply(context.Background())
	require.ErrorContains(t, err, "overflows")
}
