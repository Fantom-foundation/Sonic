package statedb

import (
	"fmt"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	io2 "github.com/Fantom-foundation/Carmen/go/state/mpt/io"
	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher/metrics"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

var carmenParams carmen.Parameters
var carmenState carmen.State
var liveStateDb carmen.StateDB

// ConfigureStateDB sets carmenParams, should be called during config parsing,
// before any other method in this package call.
func ConfigureStateDB(stateImpl string, archiveImpl string, datadir string) error {
	if stateImpl == "" || stateImpl == "geth" {
		if archiveImpl != "" {
			return fmt.Errorf("using geth statedb with Carmen archive is not supported")
		}
		return nil
	}

	var schema carmen.StateSchema
	switch strings.ToLower(stateImpl) {
	case "carmen-s3", "go-file": // "go-file" deprecated, use "carmen-s3"
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
	case "ldb", "":
		archiveType = carmen.LevelDbArchive
	case "s5":
		archiveType = carmen.S5Archive
	default:
		return fmt.Errorf("unsupported archive impl %s", archiveImpl)
	}

	carmenParams = carmen.Parameters{
		Directory: filepath.Join(datadir, "carmen"),
		Variant:   "go-file",
		Schema:    schema,
		Archive:   archiveType,
	}
	return nil
}

// InitializeStateDB initialize configured StateDB, should be called after ConfigureStateDB.
// Can be called multiple times, but once initialized, ImportFws is no longer possible.
func InitializeStateDB() error {
	if (carmenParams == carmen.Parameters{}) {
		return nil // Carmen StateDB not configured
	}
	if liveStateDb != nil {
		return nil // Carmen StateDB already initialized
	}

	err := os.MkdirAll(carmenParams.Directory, 0700)
	if err != nil {
		return fmt.Errorf("failed to create carmen dir; %v", err)
	}

	carmenState, err = carmen.NewState(carmenParams)
	if err != nil {
		return fmt.Errorf("failed to create carmen state; %s", err)
	}
	liveStateDb = carmen.CreateStateDBUsing(carmenState)

	// measure the size of carmen directory
	go metrics.MeasureDbDir("statedb/disksize", carmenParams.Directory)
	return nil
}

// ImportFws imports Fantom World State data from the genesis file into the Carmen state.
// Should be called after ConfigureStateDB, but before InitializeStateDB.
func ImportFws(reader io.Reader) error {
	if liveStateDb != nil {
		return fmt.Errorf("unable to import FWS data - Carmen State already initialized")
	}
	if carmenParams.Directory == "" || carmenParams.Schema != carmen.StateSchema(5) {
		return fmt.Errorf("unable to import FWS data - Carmen S5 not used")
	}
	err := os.MkdirAll(carmenParams.Directory, 0700)
	if err != nil {
		return fmt.Errorf("failed to create carmen dir during FWS import; %v", err)
	}
	return io2.Import(carmenParams.Directory, reader)
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

// GetArchiveBlockHeight provides the last block number available in the archive. Returns 0 if not known.
func GetArchiveBlockHeight() (height uint64, empty bool, err error) {
	if carmenState != nil {
		return liveStateDb.GetArchiveBlockHeight()
	}
	return 0, true, nil
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

// GetGenesisStateDb obtains StateDB for fake genesis generation
// Should be writable, but independent on the production live StateDb
func GetGenesisStateDb(evmState state.Database, snaps *snapshot.Tree) (*state.StateDB, error) {
	return state.NewWithSnapLayers(common.Hash(hash.Zero), evmState, snaps, 0)
}

func ShutdownStateDB() error {
	if carmenState != nil {
		err := carmenState.Close()
		if err != nil {
			return fmt.Errorf("failed to close carmen state; %s", err)
		}
		carmenState = nil
		liveStateDb = nil
	}
	return nil
}
