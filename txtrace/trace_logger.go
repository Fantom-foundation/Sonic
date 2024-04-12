package txtrace

import (
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

const (
	CALL         = "call"
	CREATE       = "create"
	SELFDESTRUCT = "suicide"
)

// TraceStructLogger is a transaction trace data collector
type TraceStructLogger struct {
	from        common.Address
	to          *common.Address
	blockHash   common.Hash
	tx          common.Hash
	txIndex     uint
	blockNumber big.Int
	value       big.Int

	gasLimit     uint64
	gasUsed      uint64
	rootTrace    *CallTrace
	inputData    []byte
	traceAddress []uint32
	output       []byte
}

// CallTrace is struct for holding tracing results
type CallTrace struct {
	Actions []ActionTrace  `json:"result"`
	Stack   []*ActionTrace `json:"-"`
}

// ActionTrace represents single interaction with blockchain
type ActionTrace struct {
	childTraces         []*ActionTrace     `json:"-"`
	Action              *AddressAction     `json:"action"`
	BlockHash           common.Hash        `json:"blockHash"`
	BlockNumber         big.Int            `json:"blockNumber"`
	Result              *TraceActionResult `json:"result,omitempty"`
	Error               string             `json:"error,omitempty"`
	Subtraces           uint64             `json:"subtraces"`
	TraceAddress        []uint32           `json:"traceAddress"`
	TransactionHash     common.Hash        `json:"transactionHash"`
	TransactionPosition uint64             `json:"transactionPosition"`
	TraceType           string             `json:"type"`
}

// AddressAction represents more specific information about
// account interaction
type AddressAction struct {
	CallType      *string         `json:"callType,omitempty"`
	From          *common.Address `json:"from"`
	To            *common.Address `json:"to,omitempty"`
	Value         hexutil.Big     `json:"value"`
	Gas           hexutil.Uint64  `json:"gas"`
	Init          *hexutil.Bytes  `json:"init,omitempty"`
	Input         *hexutil.Bytes  `json:"input,omitempty"`
	Address       *common.Address `json:"address,omitempty"`
	RefundAddress *common.Address `json:"refund_address,omitempty"`
	Balance       *hexutil.Big    `json:"balance,omitempty"`
}

// TraceActionResult holds information related to result of the
// processed transaction
type TraceActionResult struct {
	GasUsed hexutil.Uint64  `json:"gasUsed"`
	Output  *hexutil.Bytes  `json:"output,omitempty"`
	Code    *hexutil.Bytes  `json:"code,omitempty"`
	Address *common.Address `json:"address,omitempty"`
}

// NewTraceStructLogger creates new instance of trace creator
func NewTraceStructLogger(block *evmcore.EvmBlock, tx *types.Transaction, msg types.Message, index uint, gasUsed uint64) *TraceStructLogger {
	traceStructLogger := TraceStructLogger{
		tx:          tx.Hash(),
		from:        msg.From(),
		to:          msg.To(),
		value:       *msg.Value(),
		blockHash:   block.Hash,
		blockNumber: *block.Number,
		txIndex:     index,
		gasLimit:    tx.Gas(),
		gasUsed:     gasUsed,
	}
	return &traceStructLogger
}

// NewActionTrace creates new instance of type ActionTrace
func NewActionTrace(bHash common.Hash, bNumber big.Int, tHash common.Hash, tPos uint64, tType string) *ActionTrace {
	return &ActionTrace{
		BlockHash:           bHash,
		BlockNumber:         bNumber,
		TransactionHash:     tHash,
		TransactionPosition: tPos,
		TraceType:           tType,
		TraceAddress:        make([]uint32, 0),
		Result:              &TraceActionResult{},
	}
}

// NewActionTraceFromTrace creates new instance of type ActionTrace
// based on another trace
func NewActionTraceFromTrace(actionTrace *ActionTrace, tType string, traceAddress []uint32) *ActionTrace {
	trace := NewActionTrace(
		actionTrace.BlockHash,
		actionTrace.BlockNumber,
		actionTrace.TransactionHash,
		actionTrace.TransactionPosition,
		tType)
	trace.TraceAddress = traceAddress
	return trace
}

// NewAddressAction creates specific information about trace addresses
func NewAddressAction(from common.Address, gas uint64, data []byte, to *common.Address, value hexutil.Big, callType *string) *AddressAction {
	action := AddressAction{
		From:     &from,
		To:       to,
		Gas:      hexutil.Uint64(gas),
		Value:    value,
		CallType: callType,
	}
	inputHex := hexutil.Bytes(common.CopyBytes(data))
	if callType == nil {
		action.Init = &inputHex
	} else {
		action.Input = &inputHex
	}
	return &action
}

// CaptureStart implements the tracer interface to initialize the tracing operation.
func (tr *TraceStructLogger) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Tracer CaptureStart failed", r)
		}
	}()

	// Create main trace holder
	txTrace := CallTrace{
		Actions: make([]ActionTrace, 0),
	}

	if value == nil {
		value = big.NewInt(0)
	}

	// Check if To is defined. If not, it is create address call
	callType := CREATE
	var newAddress *common.Address
	if !create {
		callType = CALL
	} else {
		newAddress = &to
	}

	// Store input data
	tr.inputData = common.CopyBytes(input)
	// In new version of go-ethereum setting gas limit is done via callback CaptureTxStart
	if tr.gasLimit == 0 && gas != 0 {
		tr.gasLimit = gas
	}

	// Make transaction trace root object
	blockTrace := NewActionTrace(tr.blockHash, tr.blockNumber, tr.tx, uint64(tr.txIndex), callType)
	var txAction *AddressAction
	if create {
		txAction = NewAddressAction(tr.from, tr.gasLimit, tr.inputData, nil, hexutil.Big(*value), nil)
		if newAddress != nil {
			blockTrace.Result.Address = newAddress
			code := hexutil.Bytes(tr.output)
			blockTrace.Result.Code = &code
		}
	} else {
		txAction = NewAddressAction(tr.from, tr.gasLimit, tr.inputData, tr.to, hexutil.Big(*value), &callType)
		out := hexutil.Bytes(tr.output)
		blockTrace.Result.Output = &out
	}
	blockTrace.Action = txAction

	// Add root object into Tracer
	txTrace.AddTrace(blockTrace)
	tr.rootTrace = &txTrace

	// Init all needed variables
	tr.traceAddress = make([]uint32, 0)
	tr.rootTrace.Stack = append(tr.rootTrace.Stack, &tr.rootTrace.Actions[len(tr.rootTrace.Actions)-1])
}

// CaptureState is not used as transaction tracing doesn't need per instruction resolution
func (tr *TraceStructLogger) CaptureState(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
}

// CaptureEnter implements tracer interface for entering inner contract call
func (tr *TraceStructLogger) CaptureEnter(op vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Tracer CaptureState failed", r)
		}
	}()
	var (
		trace *ActionTrace
	)

	if tr.rootTrace == nil || len(tr.rootTrace.Stack) == 0 {
		log.Debug("There is no root trace or is empty when CaptureEnter", "tx hash", tr.tx.String())
		return
	}

	fromTrace := tr.rootTrace.Stack[len(tr.rootTrace.Stack)-1]

	if value == nil {
		value = big.NewInt(0)
	}

	// Match processed instruction and create trace based on it
	switch op {
	case vm.CREATE, vm.CREATE2:

		trace = NewActionTraceFromTrace(fromTrace, CREATE, tr.traceAddress)
		traceAction := NewAddressAction(from, gas, input, &to, hexutil.Big(*value), nil)
		trace.Action = traceAction

	case vm.CALL, vm.CALLCODE, vm.DELEGATECALL, vm.STATICCALL:

		trace = NewActionTraceFromTrace(fromTrace, CALL, tr.traceAddress)
		callType := strings.ToLower(op.String())
		traceAction := NewAddressAction(from, gas, input, &to, hexutil.Big(*value), &callType)
		trace.Action = traceAction

	case vm.SELFDESTRUCT:

		trace = NewActionTraceFromTrace(fromTrace, SELFDESTRUCT, tr.traceAddress)
		traceAction := NewAddressAction(from, gas, input, nil, hexutil.Big(*value), nil)
		traceAction.Address = &from
		traceAction.RefundAddress = &to
		traceAction.Balance = (*hexutil.Big)(value)
		trace.Action = traceAction
	}
	tr.rootTrace.Stack = append(tr.rootTrace.Stack, trace)
}

// CaptureExit is called when returning from an inner call
func (tr *TraceStructLogger) CaptureExit(output []byte, gasUsed uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Tracer CaptureExit failed", r)
		}
	}()

	if tr.rootTrace == nil {
		log.Debug("There is no root trace when CaptureExit", "tx hash", tr.tx.String())
		return
	}

	size := len(tr.rootTrace.Stack)
	if size <= 1 {
		log.Debug("CaptureExit does not match with number of CaptureEnter", "tx hash", tr.tx.String())
		return
	}

	trace := tr.rootTrace.Stack[size-1]
	tr.rootTrace.Stack = tr.rootTrace.Stack[:size-1]

	parent := tr.rootTrace.Stack[len(tr.rootTrace.Stack)-1]
	parent.childTraces = append(parent.childTraces, trace)

	trace.processOutput(output, err, false)

	result := trace.Result
	if result != nil {
		result.GasUsed = hexutil.Uint64(gasUsed)
	}
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (tr *TraceStructLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Tracer CaptureEnd failed", r)
		}
	}()

	if tr.rootTrace != nil && tr.rootTrace.lastTrace() != nil {

		trace := tr.rootTrace.lastTrace()
		trace.processOutput(output, err, true)
		if trace.Result != nil {
			// set gas used of the root call with the gas from transaction receipt
			// to present all cumulative gas used by this call and its inner calls
			trace.Result.GasUsed = hexutil.Uint64(tr.gasUsed)
		}

		tr.rootTrace.processTraces()
	}
}

// CaptureFault implements the Tracer interface to trace an execution fault
// while running an opcode. Not used for transaction tracing as error is contained
// in CaptureExit or CaptureEnd
func (tr *TraceStructLogger) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

// Handle output data and error
func (trace *ActionTrace) processOutput(output []byte, err error, rootTrace bool) {
	output = common.CopyBytes(output)
	if err == nil {
		switch trace.TraceType {
		case CREATE:
			trace.Result.Code = (*hexutil.Bytes)(&output)
			if !rootTrace {
				trace.Result.Address = trace.Action.To
				trace.Action.To = nil
			}
		case CALL:
			trace.Result.Output = (*hexutil.Bytes)(&output)
		default:
		}
		return
	} else {
		trace.Result = nil
	}

	trace.Error = err.Error()
	if trace.TraceType == CREATE {
		trace.Action.To = nil
	}
	if traceError, ok := traceErrorMapping[err.Error()]; ok {
		trace.Error = traceError
		return
	}

	switch err.(type) {
	case *vm.ErrStackOverflow:
		trace.Error = "Out of stack"
		return
	}

	if !errors.Is(err, vm.ErrExecutionReverted) || len(output) == 0 {
		return
	}
	if len(output) < 4 {
		return
	}
	if unpacked, err := abi.UnpackRevert(output); err == nil {
		trace.Error = unpacked
	}
}

// GetResult returns action traces after recording evm process
func (tr *TraceStructLogger) GetResult() *[]ActionTrace {
	if tr.rootTrace != nil {
		return &tr.rootTrace.Actions
	}
	empty := make([]ActionTrace, 0)
	return &empty
}

// AddTrace Append trace to call trace list
func (callTrace *CallTrace) AddTrace(blockTrace *ActionTrace) {
	if callTrace.Actions == nil {
		callTrace.Actions = make([]ActionTrace, 0)
	}
	callTrace.Actions = append(callTrace.Actions, *blockTrace)
}

// AddTraces Append traces to call trace list
func (callTrace *CallTrace) AddTraces(traces *[]ActionTrace, traceIndex *[]hexutil.Uint) {
	for _, trace := range *traces {
		if traceIndex == nil || isIndexEqual(traceIndex, trace.TraceAddress) {
			callTrace.AddTrace(&trace)
		}
	}
}

// isIndexEqual tells whether index and traceIndex are the same
func isIndexEqual(index *[]hexutil.Uint, traceIndex []uint32) bool {
	if len(*index) != len(traceIndex) {
		return false
	}
	for i, v := range *index {
		if uint32(v) != traceIndex[i] {
			return false
		}
	}
	return true
}

// lastTrace Get last trace in call trace list
func (callTrace *CallTrace) lastTrace() *ActionTrace {
	if len(callTrace.Actions) > 0 {
		return &callTrace.Actions[len(callTrace.Actions)-1]
	}
	return nil
}

// processTraces initiates final information distribution
// accros result traces
func (callTrace *CallTrace) processTraces() {
	trace := &callTrace.Actions[len(callTrace.Actions)-1]
	callTrace.processTrace(trace, []uint32{})
}

// processTrace goes thru all trace results and sets info
func (callTrace *CallTrace) processTrace(trace *ActionTrace, traceAddress []uint32) {
	trace.TraceAddress = traceAddress
	trace.Subtraces = uint64(len(trace.childTraces))
	for i, childTrace := range trace.childTraces {
		childAddress := childTraceAddress(traceAddress, i)
		childTrace.TraceAddress = childAddress
		callTrace.AddTrace(childTrace)
		callTrace.processTrace(callTrace.lastTrace(), childAddress)
	}
}

func childTraceAddress(a []uint32, i int) []uint32 {
	child := make([]uint32, 0, len(a)+1)
	child = append(child, a...)
	child = append(child, uint32(i))
	return child
}

// GetErrorTrace constructs filled error trace
func GetErrorTraceFromMsg(msg *types.Message, blockHash common.Hash, blockNumber big.Int, txHash common.Hash, index uint64, err error) *ActionTrace {
	if msg == nil {
		return createErrorTrace(blockHash, blockNumber, nil, &common.Address{}, txHash, 0, []byte{}, hexutil.Big{}, index, err)
	} else {
		from := msg.From()
		return createErrorTrace(blockHash, blockNumber, &from, msg.To(), txHash, msg.Gas(), msg.Data(), hexutil.Big(*msg.Value()), index, err)
	}
}

// createErrorTrace constructs filled error trace
func createErrorTrace(blockHash common.Hash, blockNumber big.Int,
	from *common.Address, to *common.Address,
	txHash common.Hash, gas uint64, input []byte,
	value hexutil.Big,
	index uint64, err error) *ActionTrace {

	var blockTrace *ActionTrace
	var txAction *AddressAction

	if from == nil {
		from = &common.Address{}
	}

	callType := CALL
	if to != nil {
		blockTrace = NewActionTrace(blockHash, blockNumber, txHash, index, CALL)
		txAction = NewAddressAction(*from, gas, input, to, hexutil.Big{}, &callType)
	} else {
		blockTrace = NewActionTrace(blockHash, blockNumber, txHash, index, CREATE)
		txAction = NewAddressAction(*from, gas, input, nil, hexutil.Big{}, nil)
	}
	blockTrace.Action = txAction
	blockTrace.Result = nil
	if err != nil {
		blockTrace.Error = err.Error()
	} else {
		blockTrace.Error = "Reverted"
	}
	return blockTrace
}

var traceErrorMapping = map[string]string{
	vm.ErrCodeStoreOutOfGas.Error():     "Out of gas",
	vm.ErrOutOfGas.Error():              "Out of gas",
	vm.ErrGasUintOverflow.Error():       "Out of gas",
	vm.ErrMaxCodeSizeExceeded.Error():   "Out of gas",
	vm.ErrInvalidJump.Error():           "Bad jump destination",
	vm.ErrExecutionReverted.Error():     "Reverted",
	vm.ErrReturnDataOutOfBounds.Error(): "Out of bounds",
	"precompiled failed":                "Built-in failed",
	"invalid input length":              "Built-in failed",
}
