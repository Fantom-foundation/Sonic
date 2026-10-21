package evmstore

import (
	"math/big"
	"testing"

	"github.com/Fantom-foundation/Carmen/go/common/amount"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/kvdb/memorydb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func openStoreForTotalSupply(t *testing.T) *Store {
	t.Helper()
	cfg := LiteStoreConfig()
	cfg.StateDb.Directory = t.TempDir()
	store := NewStore(memorydb.New(), cfg)
	require.NoError(t, store.Open())
	t.Cleanup(func() { require.NoError(t, store.Close()) })
	return store
}

func TestTotalSupply_UnavailableBeforeFirstUpdate(t *testing.T) {
	store := openStoreForTotalSupply(t)
	_, err := store.GetTotalSupply()
	require.ErrorIs(t, err, ErrTotalSupplyUnavailable)
}

func TestTotalSupply_StoreNotOpen_ReturnsError(t *testing.T) {
	store := NewStore(memorydb.New(), LiteStoreConfig())
	require.Error(t, store.UpdateTotalSupply(1))
	_, err := store.GetTotalSupply()
	require.ErrorIs(t, err, ErrTotalSupplyUnavailable)
}

func TestTotalSupply_SnapshotReflectsBalances(t *testing.T) {
	store := openStoreForTotalSupply(t)

	// the empty state has the zero hash
	stateDb, err := store.GetLiveStateDb(hash.Hash{})
	if err != nil {
		stateDb, err = store.GetLiveStateDb(hash.Hash(store.liveStateDb.GetHash()))
	}
	require.NoError(t, err)

	alice := common.Address{1}
	dead := common.HexToAddress("0x000000000000000000000000000000000000dEaD")
	stateDb.BeginBlock(1)
	stateDb.AddBalance(alice, big.NewInt(1000))
	stateDb.AddBalance(common.Address{}, big.NewInt(30))
	stateDb.AddBalance(dead, big.NewInt(12))
	stateDb.EndBlock(1)
	stateDb.Release()

	require.NoError(t, store.UpdateTotalSupply(1))
	snapshot, err := store.GetTotalSupply()
	require.NoError(t, err)
	require.EqualValues(t, 1, snapshot.Block)
	require.Equal(t, amount.New(1042), snapshot.Total)
	require.Equal(t, amount.New(30), snapshot.ZeroBalance)
	require.Equal(t, amount.New(12), snapshot.DeadBalance)
}
