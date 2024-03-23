package statedb

import (
	"encoding/json"

	cc "github.com/Fantom-foundation/Carmen/go/common"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/go-opera/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

func CreateCarmenStateDb(carmenStateDb carmen.VmStateDB, state carmen.State) state.State {
	return &carmenStateDB{
		db:    carmenStateDb,
		state: state,
	}
}

type carmenStateDB struct {
	db    carmen.VmStateDB
	state carmen.State // required for Copy method

	// current block - set by BeginBlock
	blockNum uint64

	// current transaction - set by Prepare
	txHash  common.Hash
	txIndex int
}

func (c *carmenStateDB) StartPrefetcher(namespace string) {
	// ignored
}

func (c *carmenStateDB) StopPrefetcher() {
	// ignored
}

func (c *carmenStateDB) Error() error {
	return nil
}

func (c *carmenStateDB) AddLog(log *types.Log) {
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

func (c *carmenStateDB) GetLogs(txHash common.Hash, blockNumber uint64, blockHash common.Hash) []*types.Log {
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

func (c *carmenStateDB) Logs() []*types.Log {
	panic("not supported")
}

func (c *carmenStateDB) AddPreimage(hash common.Hash, preimage []byte) {
	// ignored - preimages of keys hashes are relevant only for geth trie
}

func (c *carmenStateDB) Preimages() map[common.Hash][]byte {
	return nil // preimages of keys hashes are relevant only for geth trie
}

func (c *carmenStateDB) AddRefund(gas uint64) {
	c.db.AddRefund(gas)
}

func (c *carmenStateDB) SubRefund(gas uint64) {
	c.db.SubRefund(gas)
}

func (c *carmenStateDB) Exist(addr common.Address) bool {
	return c.db.Exist(cc.Address(addr))
}

func (c *carmenStateDB) Empty(addr common.Address) bool {
	return c.db.Empty(cc.Address(addr))
}

func (c *carmenStateDB) GetBalance(addr common.Address) *uint256.Int {
	return utils.BigIntToUint256(c.db.GetBalance(cc.Address(addr)))
}

func (c *carmenStateDB) GetNonce(addr common.Address) uint64 {
	return c.db.GetNonce(cc.Address(addr))
}

func (c *carmenStateDB) TxIndex() int {
	return c.txIndex
}

func (c *carmenStateDB) GetCode(addr common.Address) []byte {
	return c.db.GetCode(cc.Address(addr))
}

func (c *carmenStateDB) GetCodeSize(addr common.Address) int {
	return c.db.GetCodeSize(cc.Address(addr))
}

func (c *carmenStateDB) GetCodeHash(addr common.Address) common.Hash {
	return common.Hash(c.db.GetCodeHash(cc.Address(addr)))
}

func (c *carmenStateDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	return common.Hash(c.db.GetState(cc.Address(addr), cc.Key(hash)))
}

func (c *carmenStateDB) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	panic("not implemented")
}

func (c *carmenStateDB) SetTransientState(addr common.Address, key, value common.Hash) {
	panic("not implemented")
}

func (c *carmenStateDB) GetCommittedState(addr common.Address, hash common.Hash) common.Hash {
	return common.Hash(c.db.GetCommittedState(cc.Address(addr), cc.Key(hash)))
}

func (c *carmenStateDB) Database() state.Database {
	panic("not supported")
}

func (c *carmenStateDB) HasSelfDestructed(addr common.Address) bool {
	return c.db.HasSuicided(cc.Address(addr))
}

func (c *carmenStateDB) AddBalance(addr common.Address, amount *uint256.Int) {
	c.db.AddBalance(cc.Address(addr), utils.Uint256ToBigInt(amount))
}

func (c *carmenStateDB) SubBalance(addr common.Address, amount *uint256.Int) {
	c.db.SubBalance(cc.Address(addr), utils.Uint256ToBigInt(amount))
}

func (c *carmenStateDB) SetBalance(addr common.Address, amount *uint256.Int) {
	panic("not supported")
}

func (c *carmenStateDB) SetNonce(addr common.Address, nonce uint64) {
	c.db.SetNonce(cc.Address(addr), nonce)
}

func (c *carmenStateDB) SetCode(addr common.Address, code []byte) {
	c.db.SetCode(cc.Address(addr), code)
}

func (c *carmenStateDB) SetState(addr common.Address, key, value common.Hash) {
	c.db.SetState(cc.Address(addr), cc.Key(key), cc.Value(value))
}

func (c *carmenStateDB) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	panic("not supported")
}

func (c *carmenStateDB) SelfDestruct(addr common.Address) {
	c.db.Suicide(cc.Address(addr))
}

func (c *carmenStateDB) Selfdestruct6780(addr common.Address) {
	panic("not supported")
}

func (c *carmenStateDB) CreateAccount(addr common.Address) {
	c.db.CreateAccount(cc.Address(addr))
}

func (c *carmenStateDB) Copy() *state.StateDB {
	return state.NewStateDB(CreateCarmenStateDb(carmen.CreateNonCommittableStateDBUsing(c.state), c.state))
}

func (c *carmenStateDB) Snapshot() int {
	return c.db.Snapshot()
}

func (c *carmenStateDB) RevertToSnapshot(revid int) {
	c.db.RevertToSnapshot(revid)
}

func (c *carmenStateDB) GetRefund() uint64 {
	return c.db.GetRefund()
}

func (c *carmenStateDB) Finalise(deleteEmptyObjects bool) {
	c.db.EndTransaction()
}

func (c *carmenStateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	c.db.EndTransaction()
	return common.Hash(c.db.GetHash())
}

func (c *carmenStateDB) SetTxContext(txHash common.Hash, txIndex int) {
	c.txHash = txHash
	c.txIndex = txIndex
	c.db.ClearAccessList()
}

func (c *carmenStateDB) BeginBlock(number uint64) error {
	c.blockNum = number
	if db, ok := c.db.(carmen.StateDB); ok {
		db.BeginBlock()
	}
	return nil
}

func (c *carmenStateDB) EndBlock(number uint64) error {
	if db, ok := c.db.(carmen.StateDB); ok {
		db.EndBlock(number)
	}
	return nil
}

func (c *carmenStateDB) Commit(_ uint64, deleteEmptyObjects bool) (common.Hash, error) {
	// get state hash only
	return common.Hash(c.db.GetHash()), nil
}

func (c *carmenStateDB) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	// TODO: update preparation based on revision
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

func (c *carmenStateDB) AddAddressToAccessList(addr common.Address) {
	c.db.AddAddressToAccessList(cc.Address(addr))
}

func (c *carmenStateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	c.db.AddSlotToAccessList(cc.Address(addr), cc.Key(slot))
}

func (c *carmenStateDB) AddressInAccessList(addr common.Address) bool {
	return c.db.IsAddressInAccessList(cc.Address(addr))
}

func (c *carmenStateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressPresent bool, slotPresent bool) {
	return c.db.IsSlotInAccessList(cc.Address(addr), cc.Key(slot))
}

func (c *carmenStateDB) GetStorageRoot(addr common.Address) common.Hash {
	panic("not supported")
}

func (c *carmenStateDB) RawDump(opts *state.DumpConfig) state.Dump {
	panic("not supported")
}

func (c *carmenStateDB) IterativeDump(opts *state.DumpConfig, output *json.Encoder) {
	panic("not supported")
}

func (c *carmenStateDB) Dump(opts *state.DumpConfig) []byte {
	panic("not supported")
}

func (c *carmenStateDB) DumpToCollector(dc state.DumpCollector, conf *state.DumpConfig) (nextKey []byte) {
	panic("not supported")
}

// TODO: expose this as an opera interface
func (c *carmenStateDB) Release() {
	if db, ok := c.db.(carmen.NonCommittableStateDB); ok {
		db.Release()
	}
}
