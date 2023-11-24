package statedb

import (
	"bytes"
	"fmt"
	cc "github.com/Fantom-foundation/Carmen/go/common"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	io2 "github.com/Fantom-foundation/Carmen/go/state/mpt/io"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/table"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"io"
	"os"
	"path/filepath"
)

var emptyCode = crypto.Keccak256(nil)

// ImportWorldState imports Fantom World State data from the genesis file into the Carmen state.
// Should be called after ConfigureStateDB, but before InitializeStateDB.
func (m *StateDbManager) ImportWorldState(liveReader io.Reader, archiveReader io.Reader, blockNum uint64, root common.Hash) error {
	if m.parameters.Directory == "" || m.parameters.Schema != carmen.StateSchema(5) {
		return fmt.Errorf("unable to import FWS data - Carmen S5 not used")
	}
	if m.carmenState != nil {
		return fmt.Errorf("carmen state must be closed before the FWS data import")
	}

	if err := os.MkdirAll(m.parameters.Directory, 0700); err != nil {
		return fmt.Errorf("failed to create carmen dir during FWS import; %v", err)
	}
	if err := io2.ImportLiveDb(m.parameters.Directory, liveReader); err != nil {
		return fmt.Errorf("failed to import LiveDB; %v", err)
	}

	if m.parameters.Archive == carmen.S5Archive {
		archiveDir := m.parameters.Directory + string(filepath.Separator) + "archive"
		if err := os.MkdirAll(archiveDir, 0700); err != nil {
			return fmt.Errorf("failed to create carmen archive dir during FWS import; %v", err)
		}
		if err := io2.InitializeArchive(archiveDir, archiveReader, blockNum); err != nil {
			return fmt.Errorf("failed to initialize Archive; %v", err)
		}
	}

	if err := m.Open(); err != nil {
		return fmt.Errorf("failed to open Carmen for checking imported state; %v", err)
	}
	defer m.Close() // Carmen must be closed before calling this function - ensure the same state after it
	return m.checkStateHash(blockNum, root)
}

func (m *StateDbManager) ImportLegacyEvmData(chaindb ethdb.Database, evmDb kvdb.Store, blockNum uint64, root common.Hash) error {
	if m.carmenState != nil {
		return fmt.Errorf("carmen state must be closed before the legacy EVM data import")
	}
	if err := m.Open(); err != nil {
		return fmt.Errorf("failed to open StateDbManager for legacy EVM data import; %v", err)
	}
	defer m.Close()
	m.Log.Info("Importing legacy EVM data into Carmen", "index", blockNum, "root", root)

	var currentBlock uint64 = 1
	var accountsCount, slotsCount uint64 = 0, 0
	bulk := m.liveStateDb.StartBulkLoad(currentBlock)

	restartBulkIfNeeded := func () error {
		if (accountsCount + slotsCount) % 1_000_000 == 0 && currentBlock < blockNum {
			if err := bulk.Close(); err != nil {
				return err
			}
			currentBlock++
			bulk = m.liveStateDb.StartBulkLoad(currentBlock)
		}
		return nil
	}

	triedb := trie.NewDatabase(chaindb)
	t, err := trie.NewSecure(root, triedb)
	if err != nil {
		return fmt.Errorf("failed to open trie; %v", err)
	}
	preimages := table.New(evmDb, []byte("secure-key-"))

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

			bulk.CreateAccount(address)
			bulk.SetNonce(address, acc.Nonce)
			bulk.SetBalance(address, acc.Balance)


			if !bytes.Equal(acc.CodeHash, emptyCode) {
				code := rawdb.ReadCode(chaindb, common.BytesToHash(acc.CodeHash))
				if len(code) == 0 {
					return fmt.Errorf("code is missing for account %v", common.BytesToHash(accIter.LeafKey()))
				}
				bulk.SetCode(address, code)
			}

			if acc.Root != types.EmptyRootHash {
				storageTrie, err := trie.NewSecure(acc.Root, triedb)
				if err != nil {
					return fmt.Errorf("failed to open storage trie; %v", err)
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
					return fmt.Errorf("failed to iterate storage trie; %v", storageIt.Error())
				}
			}

			accountsCount++
			if err := restartBulkIfNeeded(); err != nil {
				return err
			}
		}
	}
	if err := bulk.Close(); err != nil {
		return err
	}
	// add the empty genesis block into archive
	if currentBlock < blockNum {
		bulk = m.liveStateDb.StartBulkLoad(blockNum)
		if err := bulk.Close(); err != nil {
			return err
		}
	}

	return m.checkStateHash(blockNum, root)
}

func (m *StateDbManager) checkStateHash(blockNum uint64, root common.Hash) error {
	if m.compatibleHashes {
		stateHash := m.liveStateDb.GetHash()
		if cc.Hash(root) != stateHash {
			return fmt.Errorf("importing StateDB finished with incorrect state hash: blockNum: %d expected: %x reproducedHash: %x", blockNum, root, stateHash)
		} else {
			m.Log.Info( "Importing EVM state into StateDB finished, stateRoot matches", "index", blockNum, "root", root)
		}
	}
	return nil
}
