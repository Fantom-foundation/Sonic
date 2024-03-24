package evmstore

import (
	cc "github.com/Fantom-foundation/Carmen/go/common"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/go-opera/inter/state"
	"github.com/Fantom-foundation/go-opera/utils"
	"github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

func CreateCarmenStateDb(carmenStateDb carmen.VmStateDB) state.StateDB {
	return &CarmenStateDB{
		db: carmenStateDb,
	}
}

type CarmenStateDB struct {
	db carmen.VmStateDB

	// current block - set by BeginBlock
	blockNum uint64

	// current transaction - set by Prepare
	txHash  common.Hash
	txIndex int
}

func (c *CarmenStateDB) Error() error {
	return nil
}

func (c *CarmenStateDB) AddLog(log *types.Log) {
	carmenLog := cc.Log{
		Address: cc.Address(log.Address),
		Topics:  nil,
		Data:    log.Data,
	}
	for _, topic := range log.Topics {
		carmenLog.Topics = append(carmenLog.Topics, cc.Hash(topic))
	}
	c.db.AddLog(&carmenLog)
}

func (c *CarmenStateDB) GetLogs(txHash common.Hash, blockHash common.Hash) []*types.Log {
	if txHash != c.txHash {
		panic("obtaining logs of not-current tx not supported")
	}
	carmenLogs := c.db.GetLogs()
	logs := make([]*types.Log, len(carmenLogs))
	for i, clog := range carmenLogs {
		log := &types.Log{
			Address:     common.Address(clog.Address),
			Topics:      nil,
			Data:        clog.Data,
			BlockNumber: c.blockNum,
			TxHash:      c.txHash,
			TxIndex:     uint(c.txIndex),
			BlockHash:   blockHash,
			Index:       clog.Index,
		}
		for _, topic := range clog.Topics {
			log.Topics = append(log.Topics, common.Hash(topic))
		}
		logs[i] = log
	}
	return logs
}

func (c *CarmenStateDB) AddPreimage(hash common.Hash, preimage []byte) {
	// ignored - preimages of keys hashes are relevant only for geth trie
}

func (c *CarmenStateDB) AddRefund(gas uint64) {
	c.db.AddRefund(gas)
}

func (c *CarmenStateDB) SubRefund(gas uint64) {
	c.db.SubRefund(gas)
}

func (c *CarmenStateDB) Exist(addr common.Address) bool {
	return c.db.Exist(cc.Address(addr))
}

func (c *CarmenStateDB) Empty(addr common.Address) bool {
	return c.db.Empty(cc.Address(addr))
}

func (c *CarmenStateDB) GetBalance(addr common.Address) *uint256.Int {
	return utils.BigIntToUint256(c.db.GetBalance(cc.Address(addr)))
}

func (c *CarmenStateDB) GetNonce(addr common.Address) uint64 {
	return c.db.GetNonce(cc.Address(addr))
}

func (c *CarmenStateDB) TxIndex() int {
	return c.txIndex
}

func (c *CarmenStateDB) GetCode(addr common.Address) []byte {
	return c.db.GetCode(cc.Address(addr))
}

func (c *CarmenStateDB) GetCodeSize(addr common.Address) int {
	return c.db.GetCodeSize(cc.Address(addr))
}

func (c *CarmenStateDB) GetCodeHash(addr common.Address) common.Hash {
	return common.Hash(c.db.GetCodeHash(cc.Address(addr)))
}

func (c *CarmenStateDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	return common.Hash(c.db.GetState(cc.Address(addr), cc.Key(hash)))
}

func (c *CarmenStateDB) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	panic("not implemented")
}

func (c *CarmenStateDB) GetProof(addr common.Address) ([][]byte, error) {
	panic("not supported")
}

func (c *CarmenStateDB) GetStorageProof(a common.Address, key common.Hash) ([][]byte, error) {
	panic("not supported")
}

func (c *CarmenStateDB) GetCommittedState(addr common.Address, hash common.Hash) common.Hash {
	return common.Hash(c.db.GetCommittedState(cc.Address(addr), cc.Key(hash)))
}

func (c *CarmenStateDB) StorageTrie(addr common.Address) ethstate.Trie {
	panic("not supported")
}

func (c *CarmenStateDB) HasSelfDestructed(addr common.Address) bool {
	return c.db.HasSuicided(cc.Address(addr))
}

func (c *CarmenStateDB) AddBalance(addr common.Address, amount *uint256.Int) {
	c.db.AddBalance(cc.Address(addr), utils.Uint256ToBigInt(amount))
}

func (c *CarmenStateDB) SubBalance(addr common.Address, amount *uint256.Int) {
	c.db.SubBalance(cc.Address(addr), utils.Uint256ToBigInt(amount))
}

func (c *CarmenStateDB) SetBalance(addr common.Address, amount *uint256.Int) {
	panic("not supported")
}

func (c *CarmenStateDB) SetNonce(addr common.Address, nonce uint64) {
	c.db.SetNonce(cc.Address(addr), nonce)
}

func (c *CarmenStateDB) SetCode(addr common.Address, code []byte) {
	c.db.SetCode(cc.Address(addr), code)
}

func (c *CarmenStateDB) SetState(addr common.Address, key, value common.Hash) {
	c.db.SetState(cc.Address(addr), cc.Key(key), cc.Value(value))
}

func (c *CarmenStateDB) SetTransientState(addr common.Address, key, value common.Hash) {
	panic("not implemented")
}

func (c *CarmenStateDB) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	panic("not supported")
}

func (c *CarmenStateDB) SelfDestruct(addr common.Address) {
	c.db.Suicide(cc.Address(addr))
}

func (c *CarmenStateDB) Selfdestruct6780(addr common.Address) {
	panic("not implemented")
}

func (c *CarmenStateDB) CreateAccount(addr common.Address) {
	c.db.CreateAccount(cc.Address(addr))
}

func (c *CarmenStateDB) ForEachStorage(addr common.Address, cb func(key common.Hash, value common.Hash) bool) error {
	panic("not supported")
}

func (c *CarmenStateDB) Copy() state.StateDB {
	if db, ok := c.db.(carmen.NonCommittableStateDB); ok {
		return CreateCarmenStateDb(db.Copy())
	} else {
		panic("unable to copy committable (live) StateDB")
	}
}

func (c *CarmenStateDB) Snapshot() int {
	return c.db.Snapshot()
}

func (c *CarmenStateDB) RevertToSnapshot(revid int) {
	c.db.RevertToSnapshot(revid)
}

func (c *CarmenStateDB) GetRefund() uint64 {
	return c.db.GetRefund()
}

func (c *CarmenStateDB) Finalise(deleteEmptyObjects bool) {
	c.db.EndTransaction()
}

func (c *CarmenStateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	c.db.EndTransaction()
	return common.Hash(c.db.GetHash())
}

// SetTxContext sets the current transaction hash and index which are
// used when the EVM emits new state logs.
func (c *CarmenStateDB) SetTxContext(txHash common.Hash, txIndex int) {
	c.txHash = txHash
	c.txIndex = txIndex
	c.db.ClearAccessList()
}

func (c *CarmenStateDB) BeginBlock(number uint64) {
	c.blockNum = number
	if db, ok := c.db.(carmen.StateDB); ok {
		db.BeginBlock()
	}
}

func (c *CarmenStateDB) EndBlock(number uint64) {
	if db, ok := c.db.(carmen.StateDB); ok {
		db.EndBlock(number)
	}
}

func (c *CarmenStateDB) Commit(deleteEmptyObjects bool) (common.Hash, error) {
	// get state hash only
	return common.Hash(c.db.GetHash()), nil
}

func (c *CarmenStateDB) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	// TODO: consider rules of Paris and Cancun revisions
	c.db.ClearAccessList()
	c.db.AddAddressToAccessList(cc.Address(sender))
	if dest != nil {
		c.db.AddAddressToAccessList(cc.Address(*dest))
	}
	for _, addr := range precompiles {
		c.db.AddAddressToAccessList(cc.Address(addr))
	}
	for _, el := range txAccesses {
		c.db.AddAddressToAccessList(cc.Address(el.Address))
		for _, key := range el.StorageKeys {
			c.db.AddSlotToAccessList(cc.Address(el.Address), cc.Key(key))
		}
	}
}

func (c *CarmenStateDB) AddAddressToAccessList(addr common.Address) {
	c.db.AddAddressToAccessList(cc.Address(addr))
}

func (c *CarmenStateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	c.db.AddSlotToAccessList(cc.Address(addr), cc.Key(slot))
}

func (c *CarmenStateDB) AddressInAccessList(addr common.Address) bool {
	return c.db.IsAddressInAccessList(cc.Address(addr))
}

func (c *CarmenStateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressPresent bool, slotPresent bool) {
	return c.db.IsSlotInAccessList(cc.Address(addr), cc.Key(slot))
}

func (c *CarmenStateDB) Release() {
	if db, ok := c.db.(carmen.NonCommittableStateDB); ok {
		db.Release()
	}
}
