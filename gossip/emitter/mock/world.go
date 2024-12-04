// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Fantom-foundation/go-opera/gossip/emitter (interfaces: External,TxPool,TxSigner,Signer)

// Package mock is a generated GoMock package.
package mock

import (
	big "math/big"
	reflect "reflect"

	inter "github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/inter/state"
	validatorpk "github.com/Fantom-foundation/go-opera/inter/validatorpk"
	opera "github.com/Fantom-foundation/go-opera/opera"
	vecmt "github.com/Fantom-foundation/go-opera/vecmt"
	hash "github.com/Fantom-foundation/lachesis-base/hash"
	idx "github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/ltypes"
	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	gomock "github.com/golang/mock/gomock"
)

// MockExternal is a mock of External interface.
type MockExternal struct {
	ctrl     *gomock.Controller
	recorder *MockExternalMockRecorder
}

// MockExternalMockRecorder is the mock recorder for MockExternal.
type MockExternalMockRecorder struct {
	mock *MockExternal
}

// NewMockExternal creates a new mock instance.
func NewMockExternal(ctrl *gomock.Controller) *MockExternal {
	mock := &MockExternal{ctrl: ctrl}
	mock.recorder = &MockExternalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExternal) EXPECT() *MockExternalMockRecorder {
	return m.recorder
}

// Broadcast mocks base method.
func (m *MockExternal) Broadcast(arg0 *inter.EventPayload) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Broadcast", arg0)
}

// Broadcast indicates an expected call of Broadcast.
func (mr *MockExternalMockRecorder) Broadcast(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*MockExternal)(nil).Broadcast), arg0)
}

// Build mocks base method.
func (m *MockExternal) Build(arg0 *inter.MutableEventPayload, arg1 func()) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Build", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Build indicates an expected call of Build.
func (mr *MockExternalMockRecorder) Build(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockExternal)(nil).Build), arg0, arg1)
}

// Check mocks base method.
func (m *MockExternal) Check(arg0 *inter.EventPayload, arg1 inter.Events) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockExternalMockRecorder) Check(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockExternal)(nil).Check), arg0, arg1)
}

// DagIndex mocks base method.
func (m *MockExternal) DagIndex() *vecmt.Index {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DagIndex")
	ret0, _ := ret[0].(*vecmt.Index)
	return ret0
}

// DagIndex indicates an expected call of DagIndex.
func (mr *MockExternalMockRecorder) DagIndex() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DagIndex", reflect.TypeOf((*MockExternal)(nil).DagIndex))
}

// GetEpochValidators mocks base method.
func (m *MockExternal) GetEpochValidators() (*ltypes.Validators, idx.EpochID) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEpochValidators")
	ret0, _ := ret[0].(*ltypes.Validators)
	ret1, _ := ret[1].(idx.EpochID)
	return ret0, ret1
}

// GetEpochValidators indicates an expected call of GetEpochValidators.
func (mr *MockExternalMockRecorder) GetEpochValidators() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEpochValidators", reflect.TypeOf((*MockExternal)(nil).GetEpochValidators))
}

// GetEvent mocks base method.
func (m *MockExternal) GetEvent(arg0 hash.EventHash) *inter.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", arg0)
	ret0, _ := ret[0].(*inter.Event)
	return ret0
}

// GetEvent indicates an expected call of GetEvent.
func (mr *MockExternalMockRecorder) GetEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockExternal)(nil).GetEvent), arg0)
}

// GetEventPayload mocks base method.
func (m *MockExternal) GetEventPayload(arg0 hash.EventHash) *inter.EventPayload {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventPayload", arg0)
	ret0, _ := ret[0].(*inter.EventPayload)
	return ret0
}

// GetEventPayload indicates an expected call of GetEventPayload.
func (mr *MockExternalMockRecorder) GetEventPayload(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventPayload", reflect.TypeOf((*MockExternal)(nil).GetEventPayload), arg0)
}

// GetGenesisTime mocks base method.
func (m *MockExternal) GetGenesisTime() inter.Timestamp {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenesisTime")
	ret0, _ := ret[0].(inter.Timestamp)
	return ret0
}

// GetGenesisTime indicates an expected call of GetGenesisTime.
func (mr *MockExternalMockRecorder) GetGenesisTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenesisTime", reflect.TypeOf((*MockExternal)(nil).GetGenesisTime))
}

// GetHeads mocks base method.
func (m *MockExternal) GetHeads(arg0 idx.EpochID) hash.EventHashes {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeads", arg0)
	ret0, _ := ret[0].(hash.EventHashes)
	return ret0
}

// GetHeads indicates an expected call of GetHeads.
func (mr *MockExternalMockRecorder) GetHeads(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeads", reflect.TypeOf((*MockExternal)(nil).GetHeads), arg0)
}

// GetLastEvent mocks base method.
func (m *MockExternal) GetLastEvent(arg0 idx.EpochID, arg1 idx.ValidatorID) *hash.EventHash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastEvent", arg0, arg1)
	ret0, _ := ret[0].(*hash.EventHash)
	return ret0
}

// GetLastEvent indicates an expected call of GetLastEvent.
func (mr *MockExternalMockRecorder) GetLastEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastEvent", reflect.TypeOf((*MockExternal)(nil).GetLastEvent), arg0, arg1)
}

// GetLatestBlockIndex mocks base method.
func (m *MockExternal) GetLatestBlockIndex() idx.BlockID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestBlockIndex")
	ret0, _ := ret[0].(idx.BlockID)
	return ret0
}

// GetLatestBlockIndex indicates an expected call of GetLatestBlockIndex.
func (mr *MockExternalMockRecorder) GetLatestBlockIndex() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestBlockIndex", reflect.TypeOf((*MockExternal)(nil).GetLatestBlockIndex))
}

// GetRules mocks base method.
func (m *MockExternal) GetRules() opera.Rules {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRules")
	ret0, _ := ret[0].(opera.Rules)
	return ret0
}

// GetRules indicates an expected call of GetRules.
func (mr *MockExternalMockRecorder) GetRules() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRules", reflect.TypeOf((*MockExternal)(nil).GetRules))
}

// IsBusy mocks base method.
func (m *MockExternal) IsBusy() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsBusy")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsBusy indicates an expected call of IsBusy.
func (mr *MockExternalMockRecorder) IsBusy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBusy", reflect.TypeOf((*MockExternal)(nil).IsBusy))
}

// IsSynced mocks base method.
func (m *MockExternal) IsSynced() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSynced")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSynced indicates an expected call of IsSynced.
func (mr *MockExternalMockRecorder) IsSynced() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSynced", reflect.TypeOf((*MockExternal)(nil).IsSynced))
}

// Lock mocks base method.
func (m *MockExternal) Lock() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Lock")
}

// Lock indicates an expected call of Lock.
func (mr *MockExternalMockRecorder) Lock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lock", reflect.TypeOf((*MockExternal)(nil).Lock))
}

// PeersNum mocks base method.
func (m *MockExternal) PeersNum() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeersNum")
	ret0, _ := ret[0].(int)
	return ret0
}

// PeersNum indicates an expected call of PeersNum.
func (mr *MockExternalMockRecorder) PeersNum() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeersNum", reflect.TypeOf((*MockExternal)(nil).PeersNum))
}

// Process mocks base method.
func (m *MockExternal) Process(arg0 *inter.EventPayload) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockExternalMockRecorder) Process(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockExternal)(nil).Process), arg0)
}

// StateDB mocks base method.
func (m *MockExternal) StateDB() state.StateDB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateDB")
	ret0, _ := ret[0].(state.StateDB)
	return ret0
}

// StateDB indicates an expected call of StateDB.
func (mr *MockExternalMockRecorder) StateDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateDB", reflect.TypeOf((*MockExternal)(nil).StateDB))
}

// Unlock mocks base method.
func (m *MockExternal) Unlock() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unlock")
}

// Unlock indicates an expected call of Unlock.
func (mr *MockExternalMockRecorder) Unlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unlock", reflect.TypeOf((*MockExternal)(nil).Unlock))
}

// MockTxPool is a mock of TxPool interface.
type MockTxPool struct {
	ctrl     *gomock.Controller
	recorder *MockTxPoolMockRecorder
}

// MockTxPoolMockRecorder is the mock recorder for MockTxPool.
type MockTxPoolMockRecorder struct {
	mock *MockTxPool
}

// NewMockTxPool creates a new mock instance.
func NewMockTxPool(ctrl *gomock.Controller) *MockTxPool {
	mock := &MockTxPool{ctrl: ctrl}
	mock.recorder = &MockTxPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxPool) EXPECT() *MockTxPoolMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockTxPool) Count() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	return ret0
}

// Count indicates an expected call of Count.
func (mr *MockTxPoolMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockTxPool)(nil).Count))
}

// Has mocks base method.
func (m *MockTxPool) Has(arg0 common.Hash) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Has indicates an expected call of Has.
func (mr *MockTxPoolMockRecorder) Has(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockTxPool)(nil).Has), arg0)
}

// Pending mocks base method.
func (m *MockTxPool) Pending(arg0 bool) (map[common.Address]types.Transactions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pending", arg0)
	ret0, _ := ret[0].(map[common.Address]types.Transactions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Pending indicates an expected call of Pending.
func (mr *MockTxPoolMockRecorder) Pending(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pending", reflect.TypeOf((*MockTxPool)(nil).Pending), arg0)
}

// MockTxSigner is a mock of TxSigner interface.
type MockTxSigner struct {
	ctrl     *gomock.Controller
	recorder *MockTxSignerMockRecorder
}

// MockTxSignerMockRecorder is the mock recorder for MockTxSigner.
type MockTxSignerMockRecorder struct {
	mock *MockTxSigner
}

// NewMockTxSigner creates a new mock instance.
func NewMockTxSigner(ctrl *gomock.Controller) *MockTxSigner {
	mock := &MockTxSigner{ctrl: ctrl}
	mock.recorder = &MockTxSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxSigner) EXPECT() *MockTxSignerMockRecorder {
	return m.recorder
}

// ChainID mocks base method.
func (m *MockTxSigner) ChainID() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChainID")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// ChainID indicates an expected call of ChainID.
func (mr *MockTxSignerMockRecorder) ChainID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChainID", reflect.TypeOf((*MockTxSigner)(nil).ChainID))
}

// Equal mocks base method.
func (m *MockTxSigner) Equal(arg0 types.Signer) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Equal", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Equal indicates an expected call of Equal.
func (mr *MockTxSignerMockRecorder) Equal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Equal", reflect.TypeOf((*MockTxSigner)(nil).Equal), arg0)
}

// Hash mocks base method.
func (m *MockTxSigner) Hash(arg0 *types.Transaction) common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash", arg0)
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// Hash indicates an expected call of Hash.
func (mr *MockTxSignerMockRecorder) Hash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockTxSigner)(nil).Hash), arg0)
}

// Sender mocks base method.
func (m *MockTxSigner) Sender(arg0 *types.Transaction) (common.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sender", arg0)
	ret0, _ := ret[0].(common.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sender indicates an expected call of Sender.
func (mr *MockTxSignerMockRecorder) Sender(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sender", reflect.TypeOf((*MockTxSigner)(nil).Sender), arg0)
}

// SignatureValues mocks base method.
func (m *MockTxSigner) SignatureValues(arg0 *types.Transaction, arg1 []byte) (*big.Int, *big.Int, *big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignatureValues", arg0, arg1)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(*big.Int)
	ret2, _ := ret[2].(*big.Int)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// SignatureValues indicates an expected call of SignatureValues.
func (mr *MockTxSignerMockRecorder) SignatureValues(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignatureValues", reflect.TypeOf((*MockTxSigner)(nil).SignatureValues), arg0, arg1)
}

// MockSigner is a mock of Signer interface.
type MockSigner struct {
	ctrl     *gomock.Controller
	recorder *MockSignerMockRecorder
}

// MockSignerMockRecorder is the mock recorder for MockSigner.
type MockSignerMockRecorder struct {
	mock *MockSigner
}

// NewMockSigner creates a new mock instance.
func NewMockSigner(ctrl *gomock.Controller) *MockSigner {
	mock := &MockSigner{ctrl: ctrl}
	mock.recorder = &MockSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSigner) EXPECT() *MockSignerMockRecorder {
	return m.recorder
}

// Sign mocks base method.
func (m *MockSigner) Sign(arg0 validatorpk.PubKey, arg1 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockSignerMockRecorder) Sign(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockSigner)(nil).Sign), arg0, arg1)
}
