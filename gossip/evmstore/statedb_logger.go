package evmstore

import (
	"github.com/Fantom-foundation/go-opera/inter/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"math/big"
)

func WrapStateDbWithLogger(stateDb state.StateDB, logger *tracing.Hooks) state.StateDB {
	return &LoggingStateDB{stateDb, logger}
}

type LoggingStateDB struct {
	state.StateDB
	logger *tracing.Hooks
}

func (l *LoggingStateDB) AddBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) {
	prev := l.StateDB.GetBalance(addr)
	l.StateDB.AddBalance(addr, amount, reason)
	if l.logger.OnBalanceChange != nil && !amount.IsZero() {
		l.logger.OnBalanceChange(addr, prev.ToBig(), l.StateDB.GetBalance(addr).ToBig(), reason)
	}
}

func (l *LoggingStateDB) SubBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) {
	prev := l.StateDB.GetBalance(addr)
	l.StateDB.SubBalance(addr, amount, reason)
	if l.logger.OnBalanceChange != nil && !amount.IsZero() {
		l.logger.OnBalanceChange(addr, prev.ToBig(), l.StateDB.GetBalance(addr).ToBig(), reason)
	}
}

func (l *LoggingStateDB) SetCode(addr common.Address, code []byte) {
	prevCode := l.StateDB.GetCode(addr)
	prevCodeHash := l.StateDB.GetCodeHash(addr)
	l.StateDB.SetCode(addr, code)
	if l.logger.OnCodeChange != nil {
		l.logger.OnCodeChange(addr, prevCodeHash, prevCode, l.StateDB.GetCodeHash(addr), l.StateDB.GetCode(addr))
	}
}

func (l *LoggingStateDB) SetNonce(addr common.Address, nonce uint64) {
	prev := l.StateDB.GetNonce(addr)
	l.StateDB.SetNonce(addr, nonce)
	if l.logger.OnNonceChange != nil {
		l.logger.OnNonceChange(addr, prev, nonce)
	}
}

func (l *LoggingStateDB) AddLog(log *types.Log) {
	l.StateDB.AddLog(log)
	if l.logger.OnLog != nil {
		l.logger.OnLog(log)
	}
}

func (l *LoggingStateDB) Finalise() {
	// TODO OnBalanceChange for deleted accounts
	l.StateDB.Finalise()
}

func (l *LoggingStateDB) SelfDestruct(addr common.Address) {
	prev := l.StateDB.GetBalance(addr)
	l.StateDB.SelfDestruct(addr)
	if l.logger.OnBalanceChange != nil && prev.Sign() > 0 {
		l.logger.OnBalanceChange(addr, prev.ToBig(), new(big.Int), tracing.BalanceDecreaseSelfdestruct)
	}
}
