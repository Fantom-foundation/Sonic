package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
)

type StateDB interface {
	vm.StateDB

	Error() error
	GetLogs(hash common.Hash, blockHash common.Hash) []*types.Log
	SetTxContext(thash common.Hash, ti int)
	TxIndex() int
	GetProof(addr common.Address) ([][]byte, error)
	GetStorageProof(a common.Address, key common.Hash) ([][]byte, error)
	StorageTrie(addr common.Address) state.Trie
	SetBalance(addr common.Address, amount *uint256.Int)
	SetCode(addr common.Address, code []byte)
	SetStorage(addr common.Address, storage map[common.Hash]common.Hash)
	Copy() (StateDB, error)
	Finalise() error
	Commit() (common.Hash, error)
	GetStateHash() common.Hash

	BeginBlock(number uint64)
	EndBlock(number uint64)
	Release()
}
