// Copyright 2024 Fantom Foundation
// This file is part of Aida Testing Infrastructure for Sonic
//
// Aida is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Aida is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Aida. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/stateless"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie/utils"
	"github.com/holiman/uint256"
	"log"
	"slices"
	"strings"
)

// newShadowProxy creates a StateDB instance bundling two other instances and running each
// operation on both of them, cross checking results. If the results are not equal, an error
// is logged and the result of the primary instance is returned.
func newShadowProxy(prime, shadow stateDb, compareStateHash bool) *shadowStateDb {
	return &shadowStateDb{
		shadowVmStateDb: shadowVmStateDb{
			prime:            prime,
			shadow:           shadow,
			snapshots:        []snapshotPair{},
			err:              nil,
			compareStateHash: compareStateHash,
		},
	}
}

type shadowVmStateDb struct {
	prime            stateDb
	shadow           stateDb
	snapshots        []snapshotPair
	err              error
	compareStateHash bool
}

type shadowStateDb struct {
	shadowVmStateDb
}

type snapshotPair struct {
	prime, shadow int
}

func (s *shadowVmStateDb) CreateAccount(addr common.Address) {
	s.err = errors.Join(s.err, s.run("CreateAccount", func(s vm.StateDB) error {
		s.CreateAccount(addr)
		return nil
	}))
}

func (s *shadowVmStateDb) Exist(addr common.Address) bool {
	return s.getBool("Exist", func(s vm.StateDB) bool { return s.Exist(addr) }, addr)
}

func (s *shadowVmStateDb) Empty(addr common.Address) bool {
	return s.getBool("Empty", func(s vm.StateDB) bool { return s.Empty(addr) }, addr)
}

func (s *shadowVmStateDb) SelfDestruct(addr common.Address) {
	s.err = errors.Join(s.run("SelfDestruct", func(s vm.StateDB) error {
		s.SelfDestruct(addr)
		return nil
	}))
}

func (s *shadowVmStateDb) HasSelfDestructed(addr common.Address) bool {
	return s.getBool("HasSelfDestructed", func(s vm.StateDB) bool { return s.HasSelfDestructed(addr) }, addr)
}

func (s *shadowVmStateDb) GetBalance(addr common.Address) *uint256.Int {
	return s.getUint256("GetBalance", func(s vm.StateDB) *uint256.Int { return s.GetBalance(addr) }, addr)
}

func (s *shadowVmStateDb) AddBalance(addr common.Address, value *uint256.Int, reason tracing.BalanceChangeReason) {
	s.err = errors.Join(s.run("AddBalance", func(s vm.StateDB) error {
		s.AddBalance(addr, value, reason)
		return nil
	}))
}

func (s *shadowVmStateDb) SubBalance(addr common.Address, value *uint256.Int, reason tracing.BalanceChangeReason) {
	s.err = errors.Join(s.run("SubBalance", func(s vm.StateDB) error {
		s.SubBalance(addr, value, reason)
		return nil
	}))
}

func (s *shadowVmStateDb) GetNonce(addr common.Address) uint64 {
	return s.getUint64("GetNonce", func(s vm.StateDB) uint64 { return s.GetNonce(addr) }, addr)
}

func (s *shadowVmStateDb) SetNonce(addr common.Address, value uint64) {
	s.err = errors.Join(s.run("SetNonce", func(s vm.StateDB) error {
		s.SetNonce(addr, value)
		return nil
	}))
}

func (s *shadowVmStateDb) GetCommittedState(addr common.Address, key common.Hash) common.Hash {
	// error here cannot happen
	return s.getHash("GetCommittedState", func(s vm.StateDB) common.Hash { return s.GetCommittedState(addr, key) }, addr, key)
}

func (s *shadowVmStateDb) GetState(addr common.Address, key common.Hash) common.Hash {
	return s.getHash("GetState", func(s vm.StateDB) common.Hash { return s.GetState(addr, key) }, addr, key)
}

func (s *shadowVmStateDb) SetState(addr common.Address, key common.Hash, value common.Hash) {
	s.err = errors.Join(s.err, s.run("SetState", func(s vm.StateDB) error {
		s.SetState(addr, key, value)
		return nil
	}))
}

func (s *shadowVmStateDb) SetTransientState(addr common.Address, key common.Hash, value common.Hash) {
	s.err = errors.Join(s.err, s.run("SetTransientState", func(s vm.StateDB) error {
		s.SetTransientState(addr, key, value)
		return nil
	}))
}

func (s *shadowVmStateDb) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	return s.getHash("GetTransientState", func(s vm.StateDB) common.Hash { return s.GetTransientState(addr, key) }, addr, key)
}

func (s *shadowVmStateDb) GetCode(addr common.Address) []byte {
	return s.getBytes("GetCode", func(s vm.StateDB) []byte { return s.GetCode(addr) }, addr)
}

func (s *shadowVmStateDb) GetCodeSize(addr common.Address) int {
	return s.getInt("GetCodeSize", func(s vm.StateDB) int { return s.GetCodeSize(addr) }, addr)
}

func (s *shadowVmStateDb) GetCodeHash(addr common.Address) common.Hash {
	return s.getHash("GetCodeHash", func(s vm.StateDB) common.Hash { return s.GetCodeHash(addr) }, addr)
}

func (s *shadowVmStateDb) SetCode(addr common.Address, code []byte) {
	s.err = errors.Join(s.run("SetCode", func(s vm.StateDB) error {
		s.SetCode(addr, code)
		return nil
	}))
}

func (s *shadowVmStateDb) Snapshot() int {
	pair := snapshotPair{
		s.prime.Snapshot(),
		s.shadow.Snapshot(),
	}
	s.snapshots = append(s.snapshots, pair)
	return len(s.snapshots) - 1
}

func (s *shadowVmStateDb) RevertToSnapshot(id int) {
	if id < 0 || len(s.snapshots) <= id {
		panic(fmt.Sprintf("invalid snapshot id: %v, max: %v", id, len(s.snapshots)))
	}
	s.prime.RevertToSnapshot(s.snapshots[id].prime)
	s.shadow.RevertToSnapshot(s.snapshots[id].shadow)
}

func (s *shadowVmStateDb) AddRefund(amount uint64) {
	s.err = errors.Join(s.run("AddRefund", func(s vm.StateDB) error {
		s.AddRefund(amount)
		return nil
	}))
	// check that the update value is the same
	s.getUint64("AddRefund", func(s vm.StateDB) uint64 { return s.GetRefund() })
}

func (s *shadowVmStateDb) SubRefund(amount uint64) {
	s.err = errors.Join(s.run("SubRefund", func(s vm.StateDB) error {
		s.SubRefund(amount)
		return nil
	}))
	// check that the update value is the same
	s.getUint64("SubRefund", func(s vm.StateDB) uint64 { return s.GetRefund() })
}

func (s *shadowVmStateDb) GetRefund() uint64 {
	return s.getUint64("GetRefund", func(s vm.StateDB) uint64 { return s.GetRefund() })
}

func (s *shadowVmStateDb) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	s.err = errors.Join(s.run("Prepare", func(s vm.StateDB) error {
		s.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
		return nil
	}))
}

func (s *shadowVmStateDb) AddressInAccessList(addr common.Address) bool {
	return s.getBool("AddressInAccessList", func(s vm.StateDB) bool { return s.AddressInAccessList(addr) }, addr)
}

func (s *shadowVmStateDb) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	return s.getBoolBool("SlotInAccessList", func(s vm.StateDB) (bool, bool) { return s.SlotInAccessList(addr, slot) }, addr, slot)
}

func (s *shadowVmStateDb) AddAddressToAccessList(addr common.Address) {
	s.err = errors.Join(s.run("AddAddressToAccessList", func(s vm.StateDB) error {
		s.AddAddressToAccessList(addr)
		return nil
	}))
}

func (s *shadowVmStateDb) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	s.err = errors.Join(s.run("AddSlotToAccessList", func(s vm.StateDB) error {
		s.AddSlotToAccessList(addr, slot)
		return nil
	}))
}

func (s *shadowVmStateDb) AddLog(log *types.Log) {
	s.err = errors.Join(s.run("AddPreimage", func(s vm.StateDB) error {
		s.AddLog(log)
		return nil
	}))
}

func (s *shadowVmStateDb) GetStorageRoot(addr common.Address) common.Hash {
	return s.getHash("GetStorageRoot", func(s vm.StateDB) common.Hash { return s.GetStorageRoot(addr) }, addr)
}

func (s *shadowVmStateDb) CreateContract(addr common.Address) {
	s.err = errors.Join(s.run("CreateContract", func(s vm.StateDB) error {
		s.CreateContract(addr)
		return nil
	}))
}

func (s *shadowVmStateDb) Selfdestruct6780(addr common.Address) {
	s.err = errors.Join(s.run("Selfdestruct6780", func(s vm.StateDB) error {
		s.Selfdestruct6780(addr)
		return nil
	}))
}

func (s *shadowVmStateDb) PointCache() *utils.PointCache {
	return s.prime.PointCache()
}

func (s *shadowVmStateDb) Witness() *stateless.Witness {
	return s.prime.Witness()
}

func (s *shadowVmStateDb) AddPreimage(hash common.Hash, plain []byte) {
	s.err = errors.Join(s.run("AddPreimage", func(s vm.StateDB) error {
		s.AddPreimage(hash, plain)
		return nil
	}))
}

func (s *shadowVmStateDb) Close() error {
	return errors.Join(s.err, s.prime.Close(), s.shadow.Close())
}

func (s *shadowVmStateDb) Logs() []*types.Log {
	primary := s.prime.Logs()
	shadow := s.shadow.Logs()

	for i, a := range primary {
		b := shadow[i]

		if got, want := a.Address, b.Address; got != want {
			s.logIssue("Logs", got, want)
			s.err = fmt.Errorf("Logs diverged from shadow DB")
			break
		}

		if got, want := a.Topics, b.Topics; !slices.Equal(got, want) {
			s.logIssue("Logs", got, want)
			s.err = fmt.Errorf("Logs diverged from shadow DB")
			break
		}

		if got, want := a.Data, b.Data; !slices.Equal(got, want) {
			s.logIssue("Logs", got, want)
			s.err = fmt.Errorf("Logs diverged from shadow DB")
			break
		}

	}

	return primary
}

func (s *shadowVmStateDb) Commit() common.Hash {
	return s.getStateHash("GetHash", func(s stateDb) common.Hash {
		return s.Commit()
	})
}

func (s *shadowVmStateDb) Reset() {
	s.prime.Reset()
	s.shadow.Reset()
}

func (s *shadowVmStateDb) run(opName string, op func(s vm.StateDB) error) error {
	if err := op(s.prime); err != nil {
		return fmt.Errorf("prime: %w", err)
	}
	if err := op(s.shadow); err != nil {
		return fmt.Errorf("shadow: %w", err)
	}

	return nil
}

func (s *shadowVmStateDb) getBool(opName string, op func(s vm.StateDB) bool, args ...any) bool {
	resP := op(s.prime)
	resS := op(s.shadow)
	if resP != resS {
		s.logIssue(opName, resP, resS, args)
		s.err = fmt.Errorf("%v diverged from shadow DB.", getOpcodeString(opName, args))
	}
	return resP
}

func (s *shadowVmStateDb) getBoolBool(opName string, op func(s vm.StateDB) (bool, bool), args ...any) (bool, bool) {
	resP1, resP2 := op(s.prime)
	resS1, resS2 := op(s.shadow)
	if resP1 != resS1 || resP2 != resS2 {
		s.logIssue(opName, fmt.Sprintf("(%v,%v)", resP1, resP2), fmt.Sprintf("(%v,%v)", resS1, resS2), args)
		s.err = fmt.Errorf("%v diverged from shadow DB.", getOpcodeString(opName, args))
	}
	return resP1, resP2
}

func (s *shadowVmStateDb) getInt(opName string, op func(s vm.StateDB) int, args ...any) int {
	resP := op(s.prime)
	resS := op(s.shadow)
	if resP != resS {
		s.logIssue(opName, resP, resS, args)
		s.err = fmt.Errorf("%v diverged from shadow DB.", getOpcodeString(opName, args))
	}
	return resP
}

func (s *shadowVmStateDb) getUint64(opName string, op func(s vm.StateDB) uint64, args ...any) uint64 {
	resP := op(s.prime)
	resS := op(s.shadow)
	if resP != resS {
		s.logIssue(opName, resP, resS, args)
		s.err = fmt.Errorf("%v diverged from shadow DB.", getOpcodeString(opName, args))
	}
	return resP
}

func (s *shadowVmStateDb) getStateHash(opName string, op func(s stateDb) common.Hash, args ...any) common.Hash {
	resP := op(s.prime)
	resS := op(s.shadow)
	if resP != resS {
		s.logIssue(opName, resP, resS, args)
		s.err = fmt.Errorf("%v diverged from shadow DB.", getOpcodeString(opName, args))
	}
	return resP
}

func (s *shadowVmStateDb) getHash(opName string, op func(s vm.StateDB) common.Hash, args ...any) common.Hash {
	resP := op(s.prime)
	resS := op(s.shadow)
	if resP != resS {
		s.logIssue(opName, resP, resS, args)
		s.err = fmt.Errorf("%v diverged from shadow DB.", getOpcodeString(opName, args))
	}
	return resP
}

func (s *shadowVmStateDb) getUint256(opName string, op func(s vm.StateDB) *uint256.Int, args ...any) *uint256.Int {
	resP := op(s.prime)
	resS := op(s.shadow)
	if resP.Cmp(resS) != 0 {
		s.logIssue(opName, resP, resS, args)
		s.err = fmt.Errorf("%v diverged from shadow DB.", getOpcodeString(opName, args))
	}
	return resP
}

func (s *shadowVmStateDb) getBytes(opName string, op func(s vm.StateDB) []byte, args ...any) []byte {
	resP := op(s.prime)
	resS := op(s.shadow)
	if bytes.Compare(resP, resS) != 0 {
		s.logIssue(opName, resP, resS, args)
		s.err = fmt.Errorf("%v diverged from shadow DB.", getOpcodeString(opName, args))
	}
	return resP
}

func (s *shadowVmStateDb) getError(opName string, op func(s vm.StateDB) error, args ...any) error {
	resP := op(s.prime)
	resS := op(s.shadow)
	if resP != resS {
		s.logIssue(opName, resP, resS, args)
		s.err = fmt.Errorf("%v diverged from shadow DB.", getOpcodeString(opName, args))
	}
	return resP
}

func getOpcodeString(opName string, args ...any) string {
	var opcode strings.Builder
	opcode.WriteString(fmt.Sprintf("%v(", opName))
	for _, arg := range args {
		opcode.WriteString(fmt.Sprintf("%v ", arg))
	}
	opcode.WriteString(")")
	return opcode.String()
}

func (s *shadowVmStateDb) logIssue(opName string, prime, shadow any, args ...any) {
	log.Printf("Diff for %v\n"+
		"\tPrimary: %v \n"+
		"\tShadow: %v", getOpcodeString(opName, args), prime, shadow)

}
