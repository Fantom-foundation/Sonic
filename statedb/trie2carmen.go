package statedb

import (
	"bytes"
	"fmt"
	cc "github.com/Fantom-foundation/Carmen/go/common"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
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
)

var EmptyCode = crypto.Keccak256(nil)

func IsExternalStateDbUsed() bool {
	return carmenParams != carmen.Parameters{}
}

func GetExternalStateDbHash() common.Hash {
	if liveStateDb == nil {
		return common.Hash{}
	}
	return common.Hash(liveStateDb.GetHash())
}

func ImportTrieIntoExternalStateDb(chaindb ethdb.Database, evmDb kvdb.Store, blockNum uint64, root common.Hash) error {
	if liveStateDb == nil {
		return fmt.Errorf("unable to import into Carmen State - not initialized")
	}
	var currentBlock uint64 = 1
	var accountsCount, slotsCount uint64 = 0, 0
	bulk := liveStateDb.StartBulkLoad(currentBlock)

	restartBulkIfNeeded := func () error {
		if (accountsCount + slotsCount) % 1_000_000 == 0 && currentBlock < blockNum {
			if err := bulk.Close(); err != nil {
				return err
			}
			currentBlock++
			bulk = liveStateDb.StartBulkLoad(currentBlock)
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


			if !bytes.Equal(acc.CodeHash, EmptyCode) {
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
	fmt.Printf("Imported %d accounts and %d slots into %d blocks\n", accountsCount, slotsCount, currentBlock)
	return bulk.Close()
}