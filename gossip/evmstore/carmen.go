package evmstore

import (
	"github.com/Fantom-foundation/Carmen/go/carmen"
	"github.com/Fantom-foundation/go-opera/inter/state"
	"github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// CreateCarmenStateDb creates a state-DB for updating the head state.
func CreateCarmenStateDb(carmenDb carmen.Database) state.StateDB {
	return &carmenStateDB{
		db: carmenDb,
	}
}

// createTxPoolStateDB creates a state-DB for querying the head state.
func createTxPoolStateDB(carmenDb carmen.Database) TxPoolStateDB {
	return &carmenTxPoolStateDB{
		db: carmenDb,
	}
}

// createHistoricStateDb creates a state-DB for querying historic states.
func createHistoricStateDb(block uint64, carmenDb carmen.Database) (state.StateDB, error) {
	ctxt, err := carmenDb.GetHistoricContext(block)
	if err != nil {
		return nil, err
	}
	tx, err := ctxt.BeginTransaction()
	if err != nil {
		return nil, err
	}
	return &carmenStateDB{
		db:           carmenDb,
		blockContext: ctxt,
		tx:           tx,
		blockNum:     block,
	}, nil
}

type blockContext interface {
	BeginTransaction() (carmen.TransactionContext, error)
}

type carmenStateDB struct {
	db           carmen.Database
	blockContext blockContext
	tx           carmen.TransactionContext

	// current block - set by BeginBlock
	blockNum uint64

	// current transaction - set by Prepare
	txHash  common.Hash
	txIndex int

	// last error - set by EndBlock
	err error
}

func (c *carmenStateDB) Error() error {
	return c.err
}

func (c *carmenStateDB) AddLog(log *types.Log) {
	carmenLog := carmen.Log{
		Address: carmen.Address(log.Address),
		Topics:  nil,
		Data:    log.Data,
	}
	for _, topic := range log.Topics {
		carmenLog.Topics = append(carmenLog.Topics, carmen.Hash(topic))
	}
	c.tx.AddLog(&carmenLog)
}

func (c *carmenStateDB) GetLogs(txHash common.Hash, blockHash common.Hash) []*types.Log {
	if txHash != c.txHash {
		panic("obtaining logs of not-current tx not supported")
	}
	carmenLogs := c.tx.GetLogs()
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

func (c *carmenStateDB) AddPreimage(hash common.Hash, preimage []byte) {
	// ignored - preimages of keys hashes are relevant only for geth trie
}

func (c *carmenStateDB) AddRefund(gas uint64) {
	c.tx.AddRefund(gas)
}

func (c *carmenStateDB) SubRefund(gas uint64) {
	c.tx.SubRefund(gas)
}

func (c *carmenStateDB) Exist(addr common.Address) bool {
	return c.tx.Exist(carmen.Address(addr))
}

func (c *carmenStateDB) Empty(addr common.Address) bool {
	return c.tx.Empty(carmen.Address(addr))
}

func (c *carmenStateDB) GetBalance(addr common.Address) *uint256.Int {
	res := c.tx.GetBalance(carmen.Address(addr)).Uint256()
	return &res
}

func (c *carmenStateDB) GetNonce(addr common.Address) uint64 {
	return c.tx.GetNonce(carmen.Address(addr))
}

func (c *carmenStateDB) TxIndex() int {
	return c.txIndex
}

func (c *carmenStateDB) GetCode(addr common.Address) []byte {
	return c.tx.GetCode(carmen.Address(addr))
}

func (c *carmenStateDB) GetCodeSize(addr common.Address) int {
	return c.tx.GetCodeSize(carmen.Address(addr))
}

func (c *carmenStateDB) GetCodeHash(addr common.Address) common.Hash {
	return common.Hash(c.tx.GetCodeHash(carmen.Address(addr)))
}

func (c *carmenStateDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	return common.Hash(c.tx.GetState(carmen.Address(addr), carmen.Key(hash)))
}

func (c *carmenStateDB) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	return common.Hash(c.tx.GetTransientState(carmen.Address(addr), carmen.Key(key)))
}

func (c *carmenStateDB) GetProof(addr common.Address) ([][]byte, error) {
	panic("not supported")
}

func (c *carmenStateDB) GetStorageProof(a common.Address, key common.Hash) ([][]byte, error) {
	panic("not supported")
}

func (c *carmenStateDB) GetStorageRoot(addr common.Address) common.Hash {
	return common.Hash{} // TODO
}

func (c *carmenStateDB) GetCommittedState(addr common.Address, hash common.Hash) common.Hash {
	return common.Hash(c.tx.GetCommittedState(carmen.Address(addr), carmen.Key(hash)))
}

func (c *carmenStateDB) StorageTrie(addr common.Address) ethstate.Trie {
	panic("not supported")
}

func (c *carmenStateDB) HasSelfDestructed(addr common.Address) bool {
	return c.tx.HasSelfDestructed(carmen.Address(addr))
}

func (c *carmenStateDB) AddBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) {
	c.tx.AddBalance(carmen.Address(addr), carmen.NewAmountFromUint256(amount))
}

func (c *carmenStateDB) SubBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) {
	c.tx.SubBalance(carmen.Address(addr), carmen.NewAmountFromUint256(amount))
}

func (c *carmenStateDB) SetBalance(addr common.Address, amount *uint256.Int) {
	// used only in state overrides in RPC
	current := c.tx.GetBalance(carmen.Address(addr)).Uint256()
	if current.Cmp(amount) < 0 {
		c.tx.AddBalance(carmen.Address(addr), carmen.NewAmountFromUint256(current.Sub(amount, &current)))
	} else {
		c.tx.SubBalance(carmen.Address(addr), carmen.NewAmountFromUint256(current.Sub(&current, amount)))
	}
}

func (c *carmenStateDB) SetNonce(addr common.Address, nonce uint64) {
	c.tx.SetNonce(carmen.Address(addr), nonce)
}

func (c *carmenStateDB) SetCode(addr common.Address, code []byte) {
	c.tx.SetCode(carmen.Address(addr), code)
}

func (c *carmenStateDB) SetState(addr common.Address, key, value common.Hash) {
	c.tx.SetState(carmen.Address(addr), carmen.Key(key), carmen.Value(value))
}

func (c *carmenStateDB) SetTransientState(addr common.Address, key, value common.Hash) {
	c.tx.SetTransientState(carmen.Address(addr), carmen.Key(key), carmen.Value(value))
}

func (c *carmenStateDB) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	panic("not supported")
}

func (c *carmenStateDB) SelfDestruct(addr common.Address) {
	c.tx.SelfDestruct(carmen.Address(addr))
}

func (c *carmenStateDB) Selfdestruct6780(addr common.Address) {
	panic("EIP-6780 selfdestruct not implemented")
}

func (c *carmenStateDB) CreateAccount(addr common.Address) {
	c.tx.CreateAccount(carmen.Address(addr))
}

func (c *carmenStateDB) CreateContract(addr common.Address) {
	c.tx.CreateAccount(carmen.Address(addr))
}

func (c *carmenStateDB) ForEachStorage(addr common.Address, cb func(key common.Hash, value common.Hash) bool) error {
	panic("not supported")
}

func (c *carmenStateDB) Copy() (state.StateDB, error) {
	// This is only needed by one API call and in this one
	// no previous state is present in the StateDB. So we can
	// just return a fresh one.
	return createHistoricStateDb(c.blockNum, c.db)
}

func (c *carmenStateDB) Snapshot() int {
	return c.tx.Snapshot()
}

func (c *carmenStateDB) RevertToSnapshot(revid int) {
	c.tx.RevertToSnapshot(revid)
}

func (c *carmenStateDB) GetRefund() uint64 {
	return c.tx.GetRefund()
}

func (c *carmenStateDB) Finalise() error {
	if c.err != nil {
		return c.err
	}
	err := c.tx.Commit()
	if err != nil {
		return err
	}
	tx, err := c.blockContext.BeginTransaction()
	if err != nil {
		return err
	}
	c.tx = tx
	return nil
}

// SetTxContext sets the current transaction hash and index which are
// used when the EVM emits new state logs.
func (c *carmenStateDB) SetTxContext(txHash common.Hash, txIndex int) {
	c.txHash = txHash
	c.txIndex = txIndex
	c.tx.ClearAccessList()
}

func (c *carmenStateDB) BeginBlock(number uint64) {
	if c.err != nil {
		return
	}

	if c.blockContext != nil {
		panic("cannot start block while other block is running")
	}

	c.blockContext, c.err = c.db.BeginBlock(number)
	if c.err != nil {
		return
	}
	c.blockNum = number
	c.tx, c.err = c.blockContext.BeginTransaction()
}

func (c *carmenStateDB) EndBlock(number uint64) {
	if c.err != nil {
		return
	}
	c.err = c.tx.Commit()
	if c.err != nil {
		return
	}
	c.tx = nil
	c.err = c.blockContext.(carmen.HeadBlockContext).Commit()
	if c.err != nil {
		return
	}
	c.blockContext = nil
}

func (c *carmenStateDB) Commit() (common.Hash, error) {
	// get state hash only
	var hash common.Hash
	err := c.db.QueryHeadState(func(ctxt carmen.QueryContext) {
		hash = common.Hash(ctxt.GetStateHash())
	})
	return hash, err
}

func (c *carmenStateDB) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	// TODO: consider rules of Paris and Cancun revisions
	c.tx.ClearAccessList()
	c.tx.AddAddressToAccessList(carmen.Address(sender))
	if dest != nil {
		c.tx.AddAddressToAccessList(carmen.Address(*dest))
	}
	for _, addr := range precompiles {
		c.tx.AddAddressToAccessList(carmen.Address(addr))
	}
	for _, el := range txAccesses {
		c.tx.AddAddressToAccessList(carmen.Address(el.Address))
		for _, key := range el.StorageKeys {
			c.tx.AddSlotToAccessList(carmen.Address(el.Address), carmen.Key(key))
		}
	}
}

func (c *carmenStateDB) AddAddressToAccessList(addr common.Address) {
	c.tx.AddAddressToAccessList(carmen.Address(addr))
}

func (c *carmenStateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	c.tx.AddSlotToAccessList(carmen.Address(addr), carmen.Key(slot))
}

func (c *carmenStateDB) AddressInAccessList(addr common.Address) bool {
	return c.tx.IsAddressInAccessList(carmen.Address(addr))
}

func (c *carmenStateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressPresent bool, slotPresent bool) {
	return c.tx.IsSlotInAccessList(carmen.Address(addr), carmen.Key(slot))
}

func (c *carmenStateDB) GetStateHash() common.Hash {
	var hash common.Hash
	err := c.db.QueryHeadState(func(ctxt carmen.QueryContext) {
		hash = common.Hash(ctxt.GetStateHash())
	})
	if err != nil {
		return common.Hash{} // < no option in state.StateDB interface to signal an error
	}
	return hash
}

func (c *carmenStateDB) Release() {
	if c.err != nil {
		return
	}
	if c.tx != nil {
		if c.err = c.tx.Abort(); c.err != nil {
			return
		}
		c.tx = nil
	}
	if c.blockContext != nil {
		if ctxt, ok := c.blockContext.(carmen.HeadBlockContext); ok {
			if c.err = ctxt.Abort(); c.err != nil {
				return
			}
		}
		if ctxt, ok := c.blockContext.(carmen.HistoricBlockContext); ok {
			if c.err = ctxt.Close(); c.err != nil {
				return
			}
		}
		c.blockContext = nil
	}
	c.db = nil
}

// carmenTxPoolStateDB is a StateDB adapter always targeting the head state
// and only implementing operations required by the transaction pool.
type carmenTxPoolStateDB struct {
	db carmen.Database
}

func (s *carmenTxPoolStateDB) GetNonce(addr common.Address) uint64 {
	res := uint64(0)
	err := s.db.QueryHeadState(func(ctxt carmen.QueryContext) {
		res = ctxt.GetNonce(carmen.Address(addr))
	})
	if err != nil {
		return 0 // TxPoolStateDB ignores errors
	}
	return res
}

func (s *carmenTxPoolStateDB) GetBalance(addr common.Address) *uint256.Int {
	var res uint256.Int
	err := s.db.QueryHeadState(func(ctxt carmen.QueryContext) {
		res = ctxt.GetBalance(carmen.Address(addr)).Uint256()
	})
	if err != nil {
		return nil // TxPoolStateDB ignores errors
	}
	return &res
}

func (s *carmenTxPoolStateDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	var res common.Hash
	err := s.db.QueryHeadState(func(ctxt carmen.QueryContext) {
		res = common.Hash(ctxt.GetState(carmen.Address(addr), carmen.Key(hash)))
	})
	if err != nil {
		return common.Hash{} // TxPoolStateDB ignores errors
	}
	return res
}

func (s *carmenTxPoolStateDB) Release() {
	// ignored
}
