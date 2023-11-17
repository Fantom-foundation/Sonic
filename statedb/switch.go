package statedb

import (
	"fmt"
	cc "github.com/Fantom-foundation/Carmen/go/common"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/Carmen/go/state/mpt"
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
)

type Config struct {
	CarmenParameters carmen.Parameters
}

func (cfg *Config) IsCarmen() bool {
	return cfg.CarmenParameters != (carmen.Parameters{})
}

var carmenState carmen.State
var liveStateDb carmen.StateDB
var compatibleHashes bool
var compatibleArchiveHashes bool

// InitializeStateDB initialize configured StateDB, should be called after ConfigureStateDB.
// Can be called multiple times, but once initialized, ImportWorldState is no longer possible.
func InitializeStateDB(cfg Config) error {
	if cfg.CarmenParameters == (carmen.Parameters{}) {
		return nil // Carmen StateDB not configured
	}
	if liveStateDb != nil {
		return nil // Carmen StateDB already initialized
	}

	err := os.MkdirAll(cfg.CarmenParameters.Directory, 0700)
	if err != nil {
		return fmt.Errorf("failed to create carmen dir; %v", err)
	}

	carmenState, err = carmen.NewState(cfg.CarmenParameters)
	if err != nil {
		return fmt.Errorf("failed to create carmen state; %s", err)
	}
	liveStateDb = carmen.CreateStateDBUsing(carmenState)

	// measure the size of carmen directory
	go metrics.MeasureDbDir("statedb/disksize", cfg.CarmenParameters.Directory)
	return nil
}

// ImportWorldState imports Fantom World State data from the genesis file into the Carmen state.
// Should be called after ConfigureStateDB, but before InitializeStateDB.
func ImportWorldState(liveReader io.Reader, archiveReader io.Reader, blockNum uint64, cfg Config) error {
	if liveStateDb != nil {
		return fmt.Errorf("unable to import FWS data - Carmen State already initialized")
	}
	if cfg.CarmenParameters.Directory == "" || cfg.CarmenParameters.Schema != carmen.StateSchema(5) {
		return fmt.Errorf("unable to import FWS data - Carmen S5 not used")
	}

	if err := os.MkdirAll(cfg.CarmenParameters.Directory, 0700); err != nil {
		return fmt.Errorf("failed to create carmen dir during FWS import; %v", err)
	}
	if err := io2.ImportLiveDb(cfg.CarmenParameters.Directory, liveReader); err != nil {
		return fmt.Errorf("failed to import LiveDB; %v", err)
	}

	if cfg.CarmenParameters.Archive == carmen.S5Archive {
		archiveDir := cfg.CarmenParameters.Directory + string(filepath.Separator) + "archive"
		if err := os.MkdirAll(archiveDir, 0700); err != nil {
			return fmt.Errorf("failed to create carmen archive dir during FWS import; %v", err)
		}
		if err := io2.InitializeArchive(archiveDir, archiveReader, blockNum); err != nil {
			return fmt.Errorf("failed to initialize Archive; %v", err)
		}
	}
	return nil
}

func VerifyWorldState(expectedHash common.Hash, observer mpt.VerificationObserver, cfg Config) error {
	if liveStateDb != nil {
		return fmt.Errorf("unable to verify world state data - Carmen State already initialized")
	}
	if cfg.CarmenParameters.Directory == "" || cfg.CarmenParameters.Schema != carmen.StateSchema(5) {
		return fmt.Errorf("unable to verify world state data - Carmen S5 not used")
	}
	// try to obtain information of the contained MPT
	info, err := io2.CheckMptDirectoryAndGetInfo(cfg.CarmenParameters.Directory)
	if err != nil {
		return err
	}
	// get hash of the live state
	liveState, err := carmen.NewState(cfg.CarmenParameters)
	if err != nil {
		return fmt.Errorf("failed to create carmen state; %s", err)
	}
	defer liveState.Close()
	stateHash, err := liveState.GetHash()
	if err != nil {
		return fmt.Errorf("failed to get state hash; %s", err)
	}
	if stateHash != cc.Hash(expectedHash) {
		return fmt.Errorf("validation failed - the live state hash does not match with the last block state root (%x != %x)", stateHash, expectedHash)
	}
	// verify the world state
	return mpt.VerifyFileLiveTrie(cfg.CarmenParameters.Directory, info.Config, observer)
}

// GetStateDbGeneral is used in evmstore, in situations not covered by following methods - read-only latest state
func GetStateDbGeneral(stateRoot hash.Hash, evmState state.Database, snaps *snapshot.Tree) (*state.StateDB, error) {
	if carmenState != nil {
		stateDb := carmen.CreateNonCommittableStateDBUsing(carmenState)
		if compatibleHashes && stateDb.GetHash() != cc.Hash(stateRoot) {
			return nil, fmt.Errorf("unable to get Carmen live StateDB (general) - unexpected state root (%x != %x)", liveStateDb.GetHash(), stateRoot)
		}
		return state.NewWrapper(CreateCarmenStateDb(stateDb)), nil
	} else {
		return state.NewWithSnapLayers(common.Hash(stateRoot), evmState, snaps, 0)
	}
}

// GetLiveStateDb obtains StateDB for block processing - the live writable state
func GetLiveStateDb(stateRoot hash.Hash, evmState state.Database, snaps *snapshot.Tree) (*state.StateDB, error) {
	if liveStateDb != nil {
		if compatibleHashes && liveStateDb.GetHash() != cc.Hash(stateRoot) {
			return nil, fmt.Errorf("unable to get Carmen live StateDB - unexpected state root (%x != %x)", liveStateDb.GetHash(), stateRoot)
		}
		return state.NewWrapper(CreateCarmenStateDb(liveStateDb)), nil
	} else {
		return state.NewWithSnapLayers(common.Hash(stateRoot), evmState, snaps, 0)
	}
}

// GetTxPoolStateDb obtains StateDB for TxPool evaluation - the latest finalized, read-only
func GetTxPoolStateDb(stateRoot common.Hash, evmState state.Database, snaps *snapshot.Tree) (*state.StateDB, error) {
	if carmenState != nil {
		if compatibleHashes && liveStateDb.GetHash() != cc.Hash(stateRoot) {
			return nil, fmt.Errorf("unable to get Carmen live StateDB (txpool) - unexpected state root (%x != %x)", liveStateDb.GetHash(), stateRoot)
		}
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
		if compatibleArchiveHashes && stateDb.GetHash() != cc.Hash(stateRoot) {
			return nil, fmt.Errorf("unable to get Carmen archive StateDB - unexpected state root (%x != %x)", stateDb.GetHash(), stateRoot)
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
