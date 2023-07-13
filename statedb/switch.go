package statedb

import (
	"fmt"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher/metrics"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"math/big"
	"os"
	"path/filepath"
)

var carmenState carmen.State
var liveStateDb carmen.StateDB

func InitializeStateDB(impl string, datadir string) error {
	if impl == "" || impl == "geth" {
		return nil // no initialization needed
	}
	datadir = filepath.Join(datadir, "carmen")
	go metrics.MeasureDbDir("statedb/disksize", datadir)

	if impl != "go-file" {
		return fmt.Errorf("statedb impl %s not supported", impl)
	}
	err := os.MkdirAll(datadir, 0700)
	if err != nil {
		panic(fmt.Errorf("failed to create carmen dir"))
	}
	params := carmen.Parameters{
		Schema:    carmen.StateSchema(3),
		Directory: datadir,
		Archive:   carmen.LevelDbArchive,
	}
	carmenState, err = carmen.NewGoCachedFileState(params)
	if err != nil {
		panic(fmt.Errorf("failed to create carmen state; %s", err))
	}
	liveStateDb = carmen.CreateStateDBUsing(carmenState)
	return nil
}

// GetStateDbGeneral is used in evmstore, in situations not covered by following methods - read-only latest state
func GetStateDbGeneral(stateRoot hash.Hash, evmState state.Database, snaps *snapshot.Tree) (*state.StateDB, error) {
	if carmenState != nil {
		stateDb := carmen.CreateNonCommittableStateDBUsing(carmenState)
		return state.NewWrapper(CreateCarmenStateDb(stateDb)), nil
	} else {
		return state.NewWithSnapLayers(common.Hash(stateRoot), evmState, snaps, 0)
	}
}

// GetLiveStateDb obtains StateDB for block processing - the live writable state
func GetLiveStateDb(stateRoot hash.Hash, evmState state.Database, snaps *snapshot.Tree) (*state.StateDB, error) {
	if liveStateDb != nil {
		return state.NewWrapper(CreateCarmenStateDb(liveStateDb)), nil
	} else {
		return state.NewWithSnapLayers(common.Hash(stateRoot), evmState, snaps, 0)
	}
}

// GetTxPoolStateDb obtains StateDB for TxPool evaluation - the latest finalized, read-only
func GetTxPoolStateDb(stateRoot common.Hash, evmState state.Database, snaps *snapshot.Tree) (*state.StateDB, error) {
	if carmenState != nil {
		stateDb := carmen.CreateNonCommittableStateDBUsing(carmenState)
		return state.NewWrapper(CreateCarmenStateDb(stateDb)), nil
	} else {
		return state.NewWithSnapLayers(stateRoot, evmState, snaps, 0)
	}
}

// GetRpcStateDb obtains archive StateDB for RPC requests evaluation
func GetRpcStateDb(useLatest bool, blockNum *big.Int, stateRoot common.Hash, evmState state.Database, snaps *snapshot.Tree) (*state.StateDB, error) {
	if carmenState != nil {
		if useLatest { // use read-only live state for latest or pending
			stateDb := carmen.CreateNonCommittableStateDBUsing(carmenState)
			return state.NewWrapper(CreateCarmenStateDb(stateDb)), nil
		} else { // use archive state
			stateDb, err := liveStateDb.GetArchiveStateDB(blockNum.Uint64())
			if err != nil {
				return nil, err
			}
			return state.NewWrapper(CreateCarmenStateDb(stateDb)), nil
		}
	} else {
		return state.NewWithSnapLayers(stateRoot, evmState, snaps, 0)
	}
}
