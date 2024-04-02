package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
)

type StateDB interface {
	vm.StateDB

	Error() error
	GetLogs(hash common.Hash, blockHash common.Hash) []*types.Log
	TxIndex() int
	GetProof(addr common.Address) ([][]byte, error)
	GetStorageProof(a common.Address, key common.Hash) ([][]byte, error)
	StorageTrie(addr common.Address) state.Trie
	SetBalance(addr common.Address, amount *big.Int)
	SetCode(addr common.Address, code []byte)
	SetStorage(addr common.Address, storage map[common.Hash]common.Hash)
	Copy() StateDB
	Finalise()
	Prepare(txHash common.Hash, txIndex int)
	Commit(deleteEmptyObjects bool) (common.Hash, error)

	BeginBlock(number uint64)
	EndBlock(number uint64)
	Release()
}
