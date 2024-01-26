package statedb

import (
	"fmt"
	cc "github.com/Fantom-foundation/Carmen/go/common"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	_ "github.com/Fantom-foundation/Carmen/go/state/gostate"
	"github.com/Fantom-foundation/go-opera/logger"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/pkg/errors"
	"math/big"
	"os"
)

type Config struct {
	Directory string
	EnableArchive bool
	CacheCapacity int64
}

type StateDbManager struct {
	opened bool
	parameters carmen.Parameters
	logger logger.Instance
	carmenState carmen.State
	liveStateDb carmen.StateDB
}

func CreateStateDbManager(cfg Config) *StateDbManager {
	log := logger.New("statedb")
	var archive = carmen.NoArchive
	if cfg.EnableArchive {
		archive = carmen.S5Archive
	}
	log.Log.Info("Carmen configured", "archive", archive, "dir", cfg.Directory)
	return &StateDbManager{
		parameters: carmen.Parameters{
			Directory: cfg.Directory,
			Variant:   "go-file",
			Schema:    carmen.Schema(5),
			Archive:   archive,
			LiveCache: cfg.CacheCapacity / 2,
			ArchiveCache: cfg.CacheCapacity / 2,
		},
		logger: log,
	}
}

func (m *StateDbManager) Open() error {
	if m.opened {
		return m.logAndReturnIntegrationErr("failed to open StateDbManager - already opened")
	}
	m.opened = true
	if m.parameters == (carmen.Parameters{}) {
		return nil // Carmen StateDB not configured to be used
	}

	err := os.MkdirAll(m.parameters.Directory, 0700)
	if err != nil {
		return fmt.Errorf("failed to create carmen dir \"%s\"; %v", m.parameters.Directory, err)
	}

	m.carmenState, err = carmen.NewState(m.parameters)
	if err != nil {
		return fmt.Errorf("failed to create carmen state; %s", err)
	}
	m.liveStateDb = carmen.CreateStateDBUsing(m.carmenState)
	m.logger.Log.Info("Carmen state successfully opened")
	return nil
}

// GetLiveStateDb obtains StateDB for block processing - the live writable state
func (m *StateDbManager) GetLiveStateDb(stateRoot hash.Hash) (state.StateDbInterface, error) {
	if !m.opened {
		return nil, m.logAndReturnIntegrationErr("reading not opened StateDbManager")
	}
	if m.liveStateDb.GetHash() != cc.Hash(stateRoot) {
		return nil, fmt.Errorf("unable to get Carmen live StateDB - unexpected state root (%x != %x)", m.liveStateDb.GetHash(), stateRoot)
	}
	return CreateCarmenStateDb(m.liveStateDb, m.carmenState), nil
}

// GetTxPoolStateDB obtains StateDB for TxPool evaluation - the latest finalized, read-only.
// It is also used in emitter for emitterdriver contract reading at the start of an epoch.
func (m *StateDbManager) GetTxPoolStateDB() (state.StateDbInterface, error) {
	if !m.opened {
		return nil, m.logAndReturnIntegrationErr("reading not opened StateDbManager")
	}
	// for TxPool it is ok to provide a newer state (with a different hash)
	stateDb := carmen.CreateNonCommittableStateDBUsing(m.carmenState)
	return CreateCarmenStateDb(stateDb, m.carmenState), nil
}

// GetArchiveBlockHeight provides the last block number available in the archive. Returns 0 if not known.
func (m *StateDbManager) GetArchiveBlockHeight() (height uint64, empty bool, err error) {
	if !m.opened {
		return 0, true, m.logAndReturnIntegrationErr("reading not opened StateDbManager")
	}
	if m.carmenState != nil {
		return m.liveStateDb.GetArchiveBlockHeight()
	}
	return 0, true, nil
}

// GetRpcStateDb obtains archive StateDB for RPC requests evaluation
func (m *StateDbManager) GetRpcStateDb(blockNum *big.Int, stateRoot common.Hash) (state.StateDbInterface, error) {
	if !m.opened {
		return nil, m.logAndReturnIntegrationErr("reading not opened StateDbManager")
	}
	// always use archive state (live state may mix data from various block heights)
	stateDb, err := m.liveStateDb.GetArchiveStateDB(blockNum.Uint64())
	if err != nil {
		return nil, err
	}
	if stateDb.GetHash() != cc.Hash(stateRoot) {
		return nil, fmt.Errorf("unable to get Carmen archive StateDB - unexpected state root (%x != %x)", stateDb.GetHash(), stateRoot)
	}
	return CreateCarmenStateDb(stateDb, m.carmenState), nil
}

func (m *StateDbManager) Close() error {
	m.opened = false
	if m.carmenState != nil {
		err := m.carmenState.Close()
		if err != nil {
			m.logger.Log.Warn("Failed to close carmen state", "err", err)
			return fmt.Errorf("failed to close carmen state; %s", err)
		}
		m.carmenState = nil
		m.logger.Log.Info("Carmen state successfully closed")
	}
	return nil
}

// logAndReturnIntegrationErr logs an error with its stacktrace, returns the error
func (m *StateDbManager) logAndReturnIntegrationErr(msg string) error {
	err := errors.New(msg) // create the error with a stacktrace
	m.logger.Log.Error(fmt.Sprintf("%+v", err)) // print with the stacktrace
	return err
}
