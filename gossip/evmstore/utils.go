package evmstore

import (
	"bytes"
	"errors"
	"io"

	"github.com/Fantom-foundation/go-opera/utils/iodb"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/table"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

var (
	// EmptyCode is the known hash of the empty EVM bytecode.
	EmptyCode = crypto.Keccak256(nil)
)

func (s *Store) CheckEvm(forEachState func(func(root common.Hash) (found bool, err error))) error {
	log.Info("Checking every node hash")
	nodeIt := s.table.Evm.NewIterator(nil, nil)
	defer nodeIt.Release()
	for nodeIt.Next() {
		if len(nodeIt.Key()) != 32 {
			continue
		}
		calcHash := crypto.Keccak256(nodeIt.Value())
		if !bytes.Equal(nodeIt.Key(), calcHash) {
			log.Crit("Malformed node record", "exp", common.Bytes2Hex(calcHash), "got", common.Bytes2Hex(nodeIt.Key()))
		}
	}

	log.Info("Checking every code hash")
	codeIt := table.New(s.table.Evm, []byte("c")).NewIterator(nil, nil)
	defer codeIt.Release()
	for codeIt.Next() {
		if len(codeIt.Key()) != 32 {
			continue
		}
		calcHash := crypto.Keccak256(codeIt.Value())
		if !bytes.Equal(codeIt.Key(), calcHash) {
			log.Crit("Malformed code record", "exp", common.Bytes2Hex(calcHash), "got", common.Bytes2Hex(codeIt.Key()))
		}
	}

	log.Info("Checking every preimage")
	preimageIt := table.New(s.table.Evm, []byte("secure-key-")).NewIterator(nil, nil)
	defer preimageIt.Release()
	for preimageIt.Next() {
		if len(preimageIt.Key()) != 32 {
			continue
		}
		calcHash := crypto.Keccak256(preimageIt.Value())
		if !bytes.Equal(preimageIt.Key(), calcHash) {
			log.Crit("Malformed preimage record", "exp", common.Bytes2Hex(calcHash), "got", common.Bytes2Hex(preimageIt.Key()))
		}
	}

	log.Info("Checking presence of every root")
	forEachState(func(root common.Hash) (found bool, err error) {
		stateTrie, err := s.EvmState.OpenTrie(root)
		return stateTrie != nil && err == nil, err
	})

	return nil
}

func (s *Store) ImportEvm(r io.Reader) error {
	it := iodb.NewIterator(r)
	defer it.Release()
	batch := &restrictedEvmBatch{s.table.Evm.NewBatch()}
	defer batch.Reset()
	for it.Next() {
		err := batch.Put(it.Key(), it.Value())
		if err != nil {
			return err
		}
		if batch.ValueSize() > kvdb.IdealBatchSize {
			err := batch.Write()
			if err != nil {
				return err
			}
			batch.Reset()
		}
	}
	return batch.Write()
}

type restrictedEvmBatch struct {
	kvdb.Batch
}

func IsMptKey(key []byte) bool {
	return len(key) == common.HashLength ||
		(bytes.HasPrefix(key, rawdb.CodePrefix) && len(key) == len(rawdb.CodePrefix)+common.HashLength)
}

func IsPreimageKey(key []byte) bool {
	preimagePrefix := []byte("secure-key-")
	return bytes.HasPrefix(key, preimagePrefix) && len(key) == (len(preimagePrefix)+common.HashLength)
}

func (v *restrictedEvmBatch) Put(key []byte, value []byte) error {
	if !IsMptKey(key) && !IsPreimageKey(key) {
		return errors.New("not expected prefix for EVM history dump")
	}
	return v.Batch.Put(key, value)
}
