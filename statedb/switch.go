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
	"strings"
)

var carmenState carmen.State
var liveStateDb carmen.StateDB

func InitializeStateDB(stateImpl string, archiveImpl string, datadir string) error {
	if stateImpl == "" || stateImpl == "geth" {
		if archiveImpl != "" {
			return fmt.Errorf("using geth statedb with Carmen archive is not supported")
		}
		return nil
	}

	var schema carmen.StateSchema
	switch strings.ToLower(stateImpl) {
	case "carmen-s3":
	case "go-file": // deprecated name, use "carmen-s3"
		schema = carmen.StateSchema(3)
	case "carmen-s5":
		schema = carmen.StateSchema(5)
	default:
		return fmt.Errorf("unsupported statedb impl %s", stateImpl)
	}

	var archiveType carmen.ArchiveType
	switch strings.ToLower(archiveImpl) {
	case "none":
		archiveType = carmen.NoArchive
	case "ldb":
	case "":
		archiveType = carmen.LevelDbArchive
	case "s5":
		archiveType = carmen.S5Archive
	default:
		return fmt.Errorf("unsupported archive impl %s", archiveImpl)
	}

	datadir = filepath.Join(datadir, "carmen")
	err := os.MkdirAll(datadir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create carmen dir; %v", err)
	}

	params := carmen.Parameters{
		Directory: datadir,
		Variant:   "go-file",
		Schema:    schema,
		Archive:   archiveType,
	}
	carmenState, err = carmen.NewState(params)
	if err != nil {
		return fmt.Errorf("failed to create carmen state; %s", err)
	}
	liveStateDb = carmen.CreateStateDBUsing(carmenState)

	// measure the size of carmen directory
	go metrics.MeasureDbDir("statedb/disksize", datadir)
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

// GetLatestRpcBlockNum provides the last block number available in the archive. Returns 0 if not known.
func GetLatestRpcBlockNum() (uint64, error) {
	if carmenState != nil {
		return liveStateDb.GetLastArchiveBlockHeight()
	}
	return 0, nil
}

// GetRpcStateDb obtains archive StateDB for RPC requests evaluation
func GetRpcStateDb(blockNum *big.Int, stateRoot common.Hash, evmState state.Database, snaps *snapshot.Tree) (*state.StateDB, error) {
	if carmenState != nil {
		// always use archive state (live state may mix data from various block heights)
		stateDb, err := liveStateDb.GetArchiveStateDB(blockNum.Uint64())
		if err != nil {
			return nil, err
		}
		return state.NewWrapper(CreateCarmenStateDb(stateDb)), nil
	} else {
		return state.NewWithSnapLayers(stateRoot, evmState, snaps, 0)
	}
}

func ShutdownStateDB() error {
	if carmenState != nil {
		err := carmenState.Close()
		if err != nil {
			return fmt.Errorf("failed to close carmen state; %s", err)
		}
	}
	return nil
}
