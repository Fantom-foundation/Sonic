package evmstore

import (
	"github.com/Fantom-foundation/Carmen/go/evmstore"
	"github.com/Fantom-foundation/lachesis-base/utils/cachescale"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type (
	// StoreCacheConfig is a config for the db.
	StoreCacheConfig struct {
		// Cache size for Receipts (size in bytes).
		ReceiptsSize uint
		// Cache size for Receipts (number of blocks).
		ReceiptsBlocks int
		// Cache size for TxPositions.
		TxPositions int
		// Cache size for EVM database.
		EvmDatabase int
		// Cache size for EVM snapshot.
		EvmSnap int
		// Cache size for EvmBlock (number of blocks).
		EvmBlocksNum int
		// Cache size for EvmBlock (size in bytes).
		EvmBlocksSize uint
		// Memory limit (MB) at which to start flushing dirty trie nodes to disk
		TrieDirtyLimit uint
	}
	// StoreConfig is a config for store db.
	StoreConfig struct {
		Cache StoreCacheConfig
		// Enables tracking of SHA3 preimages in the VM
		EnablePreimageRecording bool
		// Disables EVM logs indexing
		DisableLogsIndexing bool
		// Disables storing of txs positions
		DisableTxHashesIndexing bool
		// Switch to using Carmen EvmStore instead of leveldb.
		CarmenEvmStore evmstore.EvmStore `toml:"-"`
	}
)

// DefaultStoreConfig for product.
func DefaultStoreConfig(scale cachescale.Func) StoreConfig {
	return StoreConfig{
		Cache: StoreCacheConfig{
			ReceiptsSize:      scale.U(4 * opt.MiB),
			ReceiptsBlocks:    scale.I(4000),
			TxPositions:       scale.I(20000),
			EvmDatabase:       scale.I(32 * opt.MiB),
			EvmSnap:           scale.I(32 * opt.MiB),
			EvmBlocksNum:      scale.I(5000),
			EvmBlocksSize:     scale.U(6 * opt.MiB),
			TrieDirtyLimit:    scale.U(256 * opt.MiB),
		},
		EnablePreimageRecording: true,
	}
}

// LiteStoreConfig is for tests or inmemory.
func LiteStoreConfig() StoreConfig {
	return DefaultStoreConfig(cachescale.Ratio{Base: 10, Target: 1})
}
