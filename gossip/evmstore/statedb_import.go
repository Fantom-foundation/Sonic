package evmstore

import (
	"bytes"
	"context"
	"fmt"
	cc "github.com/Fantom-foundation/Carmen/go/common"
	"github.com/Fantom-foundation/Carmen/go/common/amount"
	io2 "github.com/Fantom-foundation/Carmen/go/database/mpt/io"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/utils/adapters/kvdb2ethdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/nokeyiserr"
	"github.com/Fantom-foundation/lachesis-base/kvdb/pebble"
	"github.com/Fantom-foundation/lachesis-base/kvdb/table"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"io"
	"os"
	"path/filepath"
)

var emptyCodeHash = crypto.Keccak256(nil)

// ImportLiveWorldState imports Fantom World State data from the live state genesis section.
// Must be called before the first Open call.
func (s *Store) ImportLiveWorldState(liveReader io.Reader) error {
	liveDir := filepath.Join(s.parameters.Directory, "live")
	if err := os.MkdirAll(liveDir, 0700); err != nil {
		return fmt.Errorf("failed to create carmen dir during FWS import; %v", err)
	}
	if err := io2.ImportLiveDb(liveDir, liveReader); err != nil {
		return fmt.Errorf("failed to import LiveDB; %v", err)
	}
	return nil
}

// ImportArchiveWorldState imports Fantom World State data from the archive state genesis section.
// Must be called before the first Open call.
func (s *Store) ImportArchiveWorldState(archiveReader io.Reader) error {
	if s.parameters.Archive == carmen.NoArchive {
		return nil // skip if the archive is disabled
	}
	if s.parameters.Archive == carmen.S5Archive {
		archiveDir := filepath.Join(s.parameters.Directory, "archive")
		if err := os.MkdirAll(archiveDir, 0700); err != nil {
			return fmt.Errorf("failed to create carmen archive dir during FWS import; %v", err)
		}
		if err := io2.ImportArchive(archiveDir, archiveReader); err != nil {
			return fmt.Errorf("failed to initialize Archive; %v", err)
		}
		return nil
	}
	return fmt.Errorf("archive is used, but cannot be initialized from FWS live genesis section")
}

// InitializeArchiveWorldState imports Fantom World State data from the live state genesis section.
// Must be called before the first Open call.
func (s *Store) InitializeArchiveWorldState(liveReader io.Reader, blockNum uint64) error {
	if s.parameters.Archive == carmen.NoArchive {
		return nil // skip if the archive is disabled
	}
	if s.parameters.Archive == carmen.S5Archive {
		archiveDir := filepath.Join(s.parameters.Directory, "archive")
		if err := os.MkdirAll(archiveDir, 0700); err != nil {
			return fmt.Errorf("failed to create carmen archive dir during FWS import; %v", err)
		}
		if err := io2.InitializeArchive(archiveDir, liveReader, blockNum); err != nil {
			return fmt.Errorf("failed to initialize Archive; %v", err)
		}
		return nil
	}
	return fmt.Errorf("archive is used, but cannot be initialized from FWS live genesis section")
}

// ExportLiveWorldState exports Fantom World State data for the live state genesis section.
// The Store must be closed during the call.
func (s *Store) ExportLiveWorldState(ctx context.Context, out io.Writer) error {
	liveDir := filepath.Join(s.parameters.Directory, "live")
	if err := io2.Export(ctx, io2.NewLog(), liveDir, out); err != nil {
		return fmt.Errorf("failed to export Live StateDB; %v", err)
	}
	return nil
}

// ExportArchiveWorldState exports Fantom World State data for the archive state genesis section.
// The Store must be closed during the call.
func (s *Store) ExportArchiveWorldState(ctx context.Context, out io.Writer) error {
	archiveDir := filepath.Join(s.parameters.Directory, "archive")
	if err := io2.ExportArchive(ctx, io2.NewLog(), archiveDir, out); err != nil {
		return fmt.Errorf("failed to export Archive StateDB; %v", err)
	}
	return nil
}

func (s *Store) ImportLegacyEvmData(evmItems genesis.EvmItems, blockNum uint64, root common.Hash) error {
	if err := s.Open(); err != nil {
		return fmt.Errorf("failed to open EvmStore for legacy EVM data import; %v", err)
	}
	defer s.Close()

	carmenDir, err := os.MkdirTemp(s.parameters.Directory, "opera-tmp-import-legacy-genesis")
	if err != nil {
		panic(fmt.Errorf("failed to create temporary dir for legacy EVM data import: %v", err))
	}
	defer os.RemoveAll(carmenDir)

	s.Log.Info("Unpacking legacy EVM data into a temporary directory", "dir", carmenDir)
	db, err := pebble.New(carmenDir, 1024, 100, nil, nil)
	if err != nil {
		panic(fmt.Errorf("failed to open temporary database for legacy EVM data import: %v", err))
	}
	evmItems.ForEach(func(key, value []byte) bool {
		err := db.Put(key, value)
		if err != nil {
			return false
		}
		return true
	})

	s.Log.Info("Importing legacy EVM data into Carmen", "index", blockNum, "root", root)

	var currentBlock uint64 = 1
	var accountsCount, slotsCount uint64 = 0, 0
	bulk := s.liveStateDb.StartBulkLoad(currentBlock)

	restartBulkIfNeeded := func() error {
		if (accountsCount+slotsCount)%1_000_000 == 0 && currentBlock < blockNum {
			if err := bulk.Close(); err != nil {
				return err
			}
			currentBlock++
			bulk = s.liveStateDb.StartBulkLoad(currentBlock)
		}
		return nil
	}

	chaindb := rawdb.NewDatabase(kvdb2ethdb.Wrap(nokeyiserr.Wrap(db)))
	triedb := trie.NewDatabase(chaindb)
	t, err := trie.NewSecure(root, triedb)
	if err != nil {
		return fmt.Errorf("failed to open trie; %v", err)
	}
	preimages := table.New(db, []byte("secure-key-"))

	accIter := t.NodeIterator(nil)
	for accIter.Next(true) {
		if accIter.Leaf() {

			addressBytes, err := preimages.Get(accIter.LeafKey())
			if err != nil || addressBytes == nil {
				return fmt.Errorf("missing preimage for account address hash %v; %v", accIter.LeafKey(), err)
			}
			address := cc.Address(common.BytesToAddress(addressBytes))

			var acc state.Account
			if err := rlp.DecodeBytes(accIter.LeafBlob(), &acc); err != nil {
				return fmt.Errorf("invalid account encountered during traversal; %v", err)
			}

			balance, err := amount.NewFromBigInt(acc.Balance)
			if err != nil {
				return fmt.Errorf("failed to convert balance; %w", err)
			}

			bulk.CreateAccount(address)
			bulk.SetNonce(address, acc.Nonce)
			bulk.SetBalance(address, balance)

			if !bytes.Equal(acc.CodeHash, emptyCodeHash) {
				code := rawdb.ReadCode(chaindb, common.BytesToHash(acc.CodeHash))
				if len(code) == 0 {
					return fmt.Errorf("missing code for account %v", address)
				}
				bulk.SetCode(address, code)
			}

			if acc.Root != types.EmptyRootHash {
				storageTrie, err := trie.NewSecure(acc.Root, triedb)
				if err != nil {
					return fmt.Errorf("failed to open storage trie for account %v; %v", address, err)
				}
				storageIt := storageTrie.NodeIterator(nil)
				for storageIt.Next(true) {
					if storageIt.Leaf() {
						keyBytes, err := preimages.Get(storageIt.LeafKey())
						if err != nil || keyBytes == nil {
							return fmt.Errorf("missing preimage for storage key hash %v; %v", storageIt.LeafKey(), err)
						}
						key := cc.Key(common.BytesToHash(keyBytes))

						_, valueBytes, _, err := rlp.Split(storageIt.LeafBlob())
						if err != nil {
							return fmt.Errorf("failed to decode storage; %v", err)
						}
						value := cc.Value(common.BytesToHash(valueBytes))

						bulk.SetState(address, key, value)
						slotsCount++
						if err := restartBulkIfNeeded(); err != nil {
							return err
						}
					}
				}
				if storageIt.Error() != nil {
					return fmt.Errorf("failed to iterate storage trie of account %v; %v", address, storageIt.Error())
				}
			}

			accountsCount++
			if err := restartBulkIfNeeded(); err != nil {
				return err
			}
		}
	}
	if accIter.Error() != nil {
		return fmt.Errorf("failed to iterate accounts trie; %v", accIter.Error())
	}

	if err := bulk.Close(); err != nil {
		return err
	}
	// add the empty genesis block into archive
	if currentBlock < blockNum {
		bulk = s.liveStateDb.StartBulkLoad(blockNum)
		if err := bulk.Close(); err != nil {
			return err
		}
	}
	return nil
}
