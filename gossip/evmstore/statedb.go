package evmstore

import (
	"fmt"
	"math/big"

	"github.com/Fantom-foundation/Carmen/go/carmen"
	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/inter/state"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
)

// GetLiveStateDb obtains StateDB for block processing - the live writable state
func (s *Store) GetLiveStateDb(stateRoot hash.Hash) (state.StateDB, error) {
	if s.carmenDb == nil {
		return nil, fmt.Errorf("unable to get live StateDb - EvmStore is not open")
	}
	stateDb := CreateCarmenStateDb(s.carmenDb).(*carmenStateDB)
	actualHash := stateDb.GetStateHash()
	if actualHash != common.Hash(stateRoot) {
		stateDb.Release()
		return nil, fmt.Errorf("unable to get Carmen live StateDB - unexpected state root (%x != %x)", actualHash, stateRoot)
	}
	return stateDb, nil
}

// GetReadOnlyHeadStateDb obtains StateDB for block processing without committing
// results to the chain at the end of a block. This is only supposed to be used in tests.
func (s *Store) GetReadOnlyHeadStateDb() (state.StateDB, error) {
	// Flush changes to sync head and archive.
	if err := s.carmenDb.Flush(); err != nil {
		return nil, err
	}
	height, empty, err := s.GetArchiveBlockHeight()
	if err != nil {
		return nil, err
	}
	if empty {
		return nil, fmt.Errorf("unable to obtain read-only state for empty chain")
	}
	return createHistoricStateDb(height, s.carmenDb), nil
}

type TxPoolStateDB interface {
	evmcore.TxPoolStateDB
	GetState(common.Address, common.Hash) common.Hash
}

// GetTxPoolStateDB obtains StateDB for TxPool evaluation - the latest finalized, read-only.
// It is also used in emitter for emitterdriver contract reading at the start of an epoch.
func (s *Store) GetTxPoolStateDB() (TxPoolStateDB, error) {
	// for TxPool and emitter it is ok to provide the newest state (and ignore the expected hash)
	if s.carmenDb == nil {
		return nil, fmt.Errorf("unable to get TxPool StateDb - EvmStore is not open")
	}
	return createTxPoolStateDB(s.carmenDb), nil
}

// GetArchiveBlockHeight provides the last block number available in the archive. Returns 0 if not known.
func (s *Store) GetArchiveBlockHeight() (height uint64, empty bool, err error) {
	if s.carmenDb == nil {
		return 0, true, fmt.Errorf("unable to get archive block height - EvmStore is not open")
	}
	res, err := s.carmenDb.GetArchiveBlockHeight()
	if err != nil {
		return 0, false, err
	}
	if res < 0 {
		return 0, true, nil
	}
	return uint64(res), false, nil
}

// GetRpcStateDb obtains archive StateDB for RPC requests evaluation
func (s *Store) GetRpcStateDb(blockNum *big.Int, stateRoot common.Hash) (state.StateDB, error) {
	// always use archive state (live state may mix data from various block heights)
	if s.carmenDb == nil {
		return nil, fmt.Errorf("unable to get RPC StateDb - EvmStore is not open")
	}
	stateDb := createHistoricStateDb(blockNum.Uint64(), s.carmenDb).(*carmenStateDB)
	actualHash := stateDb.GetStateHash()
	if actualHash != stateRoot {
		stateDb.Release()
		return nil, fmt.Errorf("unable to get Carmen archive StateDB - unexpected state root (%x != %x)", actualHash, stateRoot)
	}
	return stateDb, nil
}

// CheckLiveStateHash returns if the hash of the current live StateDB hash matches (and fullsync is possible)
func (s *Store) CheckLiveStateHash(blockNum idx.Block, root hash.Hash) error {
	if s.carmenDb == nil {
		return fmt.Errorf("unable to get live state - EvmStore is not open")
	}
	var stateHash carmen.Hash
	err := s.carmenDb.QueryHeadState(func(query carmen.QueryContext) {
		stateHash = query.GetStateHash()
	})
	if err != nil {
		return err
	}
	if hash.Hash(stateHash) != root {
		return fmt.Errorf("hash of the EVM state is incorrect: blockNum: %d expected: %x reproducedHash: %x", blockNum, root, stateHash)
	}
	return nil
}

// CheckArchiveStateHash returns if the hash of the given archive StateDB hash matches
func (s *Store) CheckArchiveStateHash(blockNum idx.Block, root hash.Hash) error {
	if s.carmenDb == nil {
		return fmt.Errorf("unable to get live state - EvmStore is not open")
	}
	var stateHash carmen.Hash
	err := s.carmenDb.QueryHistoricState(uint64(blockNum), func(query carmen.QueryContext) {
		stateHash = query.GetStateHash()
	})
	if err != nil {
		return err
	}
	if hash.Hash(stateHash) != root {
		return fmt.Errorf("hash of the archive EVM state is incorrect: blockNum: %d expected: %x actual: %x", blockNum, root, stateHash)
	}
	return nil
}
