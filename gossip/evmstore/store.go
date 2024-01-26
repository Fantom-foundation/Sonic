package evmstore

import (
	"github.com/Fantom-foundation/go-opera/statedb"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/nokeyiserr"
	"github.com/Fantom-foundation/lachesis-base/kvdb/table"
	"github.com/Fantom-foundation/lachesis-base/utils/wlru"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/Fantom-foundation/go-opera/logger"
	"github.com/Fantom-foundation/go-opera/topicsdb"
	"github.com/Fantom-foundation/go-opera/utils/adapters/kvdb2ethdb"
	"github.com/Fantom-foundation/go-opera/utils/rlpstore"
)

const nominalSize uint = 1

// Store is a node persistent storage working over physical key-value database.
type Store struct {
	cfg StoreConfig

	mainDB kvdb.Store
	table struct {
		Evm kvdb.Store `table:"M"`
		// API-only tables
		Receipts    kvdb.Store `table:"r"`
		TxPositions kvdb.Store `table:"x"`
		Txs         kvdb.Store `table:"X"`
	}

	EvmDb    ethdb.Database
	EvmState state.Database
	EvmLogs  topicsdb.Index

	backend Backend

	cache struct {
		TxPositions *wlru.Cache `cache:"-"` // store by pointer
		Receipts    *wlru.Cache `cache:"-"` // store by value
		EvmBlocks   *wlru.Cache `cache:"-"` // store by pointer
	}

	rlp rlpstore.Helper

	triegc *prque.Prque // Priority queue mapping block numbers to tries to gc

	logger.Instance
	sdbm *statedb.StateDbManager
}

const (
	TriesInMemory = 16
)

// NewStore creates store over key-value db.
func NewStore(mainDB kvdb.Store, cfg StoreConfig, sdbm *statedb.StateDbManager) *Store {
	s := &Store{
		cfg:      cfg,
		mainDB:   mainDB,
		Instance: logger.New("evm-store"),
		rlp:      rlpstore.Helper{logger.New("rlp")},
		triegc:   prque.New(nil),
		sdbm:     sdbm,
	}

	if cfg.CarmenEvmStore != nil {
		s.backend = carmenBackend{cfg.CarmenEvmStore}
	} else {
		s.backend = legacyBackend{s}
	}

	table.MigrateTables(&s.table, s.mainDB)

	s.initEVMDB()
	if cfg.DisableLogsIndexing {
		s.EvmLogs = topicsdb.NewDummy()
	} else {
		s.EvmLogs = topicsdb.NewWithThreadPool(mainDB)
	}
	s.initCache()

	return s
}

// Close closes underlying database.
func (s *Store) Close() {
	// set all table/cache fields to nil
	table.MigrateTables(&s.table, nil)
	table.MigrateCaches(&s.cache, func() interface{} {
		return nil
	})
	s.EvmLogs.Close()
}

func (s *Store) initCache() {
	s.cache.Receipts = s.makeCache(s.cfg.Cache.ReceiptsSize, s.cfg.Cache.ReceiptsBlocks)
	s.cache.TxPositions = s.makeCache(nominalSize*uint(s.cfg.Cache.TxPositions), s.cfg.Cache.TxPositions)
	s.cache.EvmBlocks = s.makeCache(s.cfg.Cache.EvmBlocksSize, s.cfg.Cache.EvmBlocksNum)
}

func (s *Store) initEVMDB() {
	s.EvmDb = rawdb.NewDatabase(
		kvdb2ethdb.Wrap(
			nokeyiserr.Wrap(
				s.table.Evm)))
	s.EvmState = state.NewDatabaseWithConfig(s.EvmDb, &trie.Config{
		Cache:     s.cfg.Cache.EvmDatabase / opt.MiB,
		Preimages: s.cfg.EnablePreimageRecording,
	})
}

func (s *Store) EVMDB() kvdb.Store {
	return s.table.Evm
}

// Commit changes.
func (s *Store) Commit(block idx.Block, root hash.Hash) error {
	triedb := s.EvmState.TrieDB()
	stateRoot := common.Hash(root)
	err := triedb.Commit(stateRoot, false, nil)
	if err != nil {
		s.Log.Error("Failed to flush trie DB into main DB", "err", err)
	}
	return err
}

// Cap flush matured singleton nodes to disk
func (s *Store) Cap() {
	triedb := s.EvmState.TrieDB()
	var (
		nodes, imgs = triedb.Size()
		limit       = common.StorageSize(s.cfg.Cache.TrieDirtyLimit)
	)
	// If we exceeded our memory allowance, flush matured singleton nodes to disk
	if nodes > limit+ethdb.IdealBatchSize || imgs > 4*1024*1024 {
		triedb.Cap(limit)
	}
}

// StateDB returns state database.
func (s *Store) StateDB(from hash.Hash) (state.StateDbInterface, error) {
	return s.sdbm.GetTxPoolStateDB()
}

// CheckLiveStateDbHash returns if the hash of the current live StateDB hash matches (and fullsync is possible)
func (s *Store) CheckLiveStateDbHash(blockNum idx.Block, from hash.Hash) bool {
	err := s.sdbm.CheckLiveStateHash(uint64(blockNum), common.Hash(from))
	return err == nil
}

// IndexLogs indexes EVM logs
func (s *Store) IndexLogs(recs ...*types.Log) {
	err := s.EvmLogs.Push(recs...)
	if err != nil {
		s.Log.Crit("DB logs index error", "err", err)
	}
}

/*
 * Utils:
 */

func (s *Store) makeCache(weight uint, size int) *wlru.Cache {
	cache, err := wlru.New(weight, size)
	if err != nil {
		s.Log.Crit("Failed to create LRU cache", "err", err)
		return nil
	}
	return cache
}
