package ethapi

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/Fantom-foundation/go-opera/txtrace"
	"github.com/Fantom-foundation/go-opera/utils/signers/gsignercache"
)

// PublicTxTraceAPI provides an API to access transaction tracing
// It offers only methods that operate on public data that is freely available to anyone
type PublicTxTraceAPI struct {
	b Backend
}

// NewPublicTxTraceAPI creates a new transaction trace API
func NewPublicTxTraceAPI(b Backend) *PublicTxTraceAPI {
	return &PublicTxTraceAPI{b}
}

// Transaction - trace_transaction function returns transaction inner traces
func (s *PublicTxTraceAPI) Transaction(ctx context.Context, hash common.Hash) (*[]txtrace.ActionTrace, error) {
	defer func(start time.Time) {
		log.Debug("Executing trace_transaction call finished", "txHash", hash.String(), "runtime", time.Since(start))
	}(time.Now())
	return s.traceTxHash(ctx, hash, nil)
}

// Block - trace_block function returns transaction traces in given block
func (s *PublicTxTraceAPI) Block(ctx context.Context, numberOrHash rpc.BlockNumberOrHash) (*[]txtrace.ActionTrace, error) {

	blockNumber, _ := numberOrHash.Number()

	if uint64(blockNumber.Int64()) > s.b.CurrentBlock().NumberU64() {
		return nil, fmt.Errorf("requested block nr %v > current node block nr %v", blockNumber.Int64(), s.b.CurrentBlock().NumberU64())
	}

	defer func(start time.Time) {
		log.Debug("Executing trace_block call finished", "block", blockNumber.Int64(), "runtime", time.Since(start))
	}(time.Now())

	block, err := s.b.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("cannot get block %v from db got %v", blockNumber.Int64(), err.Error())
	}

	traces, err := s.replayBlock(ctx, block, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot trace block %v got %v", blockNumber.Int64(), err.Error())
	}

	return traces, nil
}

// Get - trace_get function returns transaction traces on specified index position of the traces
// If index is nil, then just root trace is returned
func (s *PublicTxTraceAPI) Get(ctx context.Context, hash common.Hash, traceIndex []hexutil.Uint) (*[]txtrace.ActionTrace, error) {
	defer func(start time.Time) {
		log.Debug("Executing trace_get call finished", "txHash", hash.String(), "index", traceIndex, "runtime", time.Since(start))
	}(time.Now())
	return s.traceTxHash(ctx, hash, &traceIndex)
}

// traceTxHash looks for a block of this transaction hash and trace it
func (s *PublicTxTraceAPI) traceTxHash(ctx context.Context, hash common.Hash, traceIndex *[]hexutil.Uint) (*[]txtrace.ActionTrace, error) {
	_, blockNumber, _, _ := s.b.GetTransaction(ctx, hash)
	blkNr := rpc.BlockNumber(blockNumber)
	block, err := s.b.BlockByNumber(ctx, blkNr)
	if err != nil {
		log.Debug("Cannot get block from db", "blockNr", blkNr, "error", err.Error())
		return nil, err
	}

	return s.replayBlock(ctx, block, &hash, traceIndex)
}

// Replays block and returns traces acording to parameters
//
// txHash
//   - if is nil, all transaction traces in the block are collected
//   - is value, then only trace for that transaction is returned
//
// traceIndex - when specified, then only trace on that index is returned
func (s *PublicTxTraceAPI) replayBlock(ctx context.Context, block *evmcore.EvmBlock, txHash *common.Hash, traceIndex *[]hexutil.Uint) (*[]txtrace.ActionTrace, error) {

	if block == nil {
		return nil, fmt.Errorf("invalid block for tracing")
	}

	if block.NumberU64() == 0 {
		return nil, fmt.Errorf("genesis block is not traceable")
	}

	blockNumber := block.Number.Int64()
	parentBlockNr := rpc.BlockNumber(blockNumber - 1)
	callTrace := txtrace.CallTrace{
		Actions: make([]txtrace.ActionTrace, 0),
	}

	signer := gsignercache.Wrap(types.MakeSigner(s.b.ChainConfig(), block.Number))

	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, rpc.BlockNumberOrHash{BlockNumber: &parentBlockNr})
	if err != nil {
		return nil, fmt.Errorf("cannot get state for block %v, error: %v", block.NumberU64(), err.Error())
	}
	defer state.Release()

	receipts, err := s.b.GetReceiptsByNumber(ctx, rpc.BlockNumber(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("cannot get receipts for block %v, error: %v", block.NumberU64(), err.Error())
	}

	// loop thru all transactions in the block and process them
	for i, tx := range block.Transactions {

		// replay only needed transaction if specified
		if txHash == nil || *txHash == tx.Hash() {

			log.Debug("Replaying transaction", "txHash", tx.Hash().String())

			tx, _, index, err := s.b.GetTransaction(ctx, tx.Hash())
			if err != nil {
				return nil, fmt.Errorf("cannot get tranasction info %s, error %s", tx.Hash().String(), err)
			}

			msg, err := evmcore.TxAsMessage(tx, signer, block.BaseFee)
			if err != nil {
				return nil, fmt.Errorf("cannot get message from transaction %s, error %s", tx.Hash().String(), err)
			}

			txTraces, err := s.traceTx(ctx, s.b, block.Header(), msg, state, block, tx, index, receipts[i].Status, s.b.ChainConfig())
			if err != nil {
				return nil, fmt.Errorf("cannot get transaction trace for transaction %s, error %s", tx.Hash().String(), err)
			} else {
				callTrace.AddTraces(txTraces, traceIndex)
			}

			// already replayed specified transaction so end loop
			if txHash != nil {
				break
			}

		} else if txHash != nil {

			// Replay transaction without tracing to prepare state for next transaction
			log.Debug("Replaying transaction without trace", "txHash", tx.Hash().String())
			msg, err := evmcore.TxAsMessage(tx, signer, block.BaseFee)
			if err != nil {
				return nil, fmt.Errorf("cannot get message from transaction %s, error %s", tx.Hash().String(), err)
			}

			state.Prepare(tx.Hash(), i)
			vmConfig := opera.DefaultVMConfig
			vmConfig.NoBaseFee = true
			vmConfig.Debug = false
			vmConfig.Tracer = nil

			vmenv, _, err := s.b.GetEVM(ctx, msg, state, block.Header(), &vmConfig)
			if err != nil {
				return nil, fmt.Errorf("cannot initialize vm for transaction %s, error: %s", tx.Hash().String(), err.Error())
			}

			res, err := evmcore.ApplyMessage(vmenv, msg, new(evmcore.GasPool).AddGas(msg.Gas()))
			failed := false
			if err != nil {
				failed = true
				log.Error("Cannot replay transaction", "txHash", tx.Hash().String(), "err", err.Error())
			}

			if res != nil && res.Err != nil {
				failed = true
				log.Debug("Error replaying transaction", "txHash", tx.Hash().String(), "err", res.Err.Error())
			}

			state.Finalise(true)

			// Check correct replay status according to receipt data
			if (failed && receipts[i].Status == 1) || (!failed && receipts[i].Status == 0) {
				return nil, fmt.Errorf("invalid transaction replay state at %s", tx.Hash().String())
			}
		}
	}

	// In case of empty result create empty trace for empty block
	if len(callTrace.Actions) == 0 {
		if traceIndex != nil || txHash != nil {
			return nil, nil
		} else {
			return getEmptyBlockTrace(block.Hash, *block.Number), nil
		}
	}

	return &callTrace.Actions, nil
}

// traceTx trace transaction with EVM replay and return processed result
func (s *PublicTxTraceAPI) traceTx(
	ctx context.Context, b Backend, header *evmcore.EvmHeader, msg types.Message,
	state *state.StateDB, block *evmcore.EvmBlock, tx *types.Transaction, index uint64,
	status uint64, chainConfig *params.ChainConfig) (*[]txtrace.ActionTrace, error) {

	// Providing default config with tracer
	cfg := opera.DefaultVMConfig
	cfg.Debug = true
	txTracer := txtrace.NewTraceStructLogger(block, tx, msg, uint(index))
	cfg.Tracer = txTracer
	cfg.NoBaseFee = true

	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var timeout time.Duration = 5 * time.Second
	if s.b.RPCEVMTimeout() > 0 {
		timeout = s.b.RPCEVMTimeout()
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)

	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	vmenv, _, err := b.GetEVM(ctx, msg, state, header, &cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize vm for transaction %s, error: %s", tx.Hash().String(), err.Error())
	}

	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		vmenv.Cancel()
	}()

	// Setup the gas pool and stateDB
	gp := new(evmcore.GasPool).AddGas(msg.Gas())
	state.Prepare(tx.Hash(), int(index))
	result, err := evmcore.ApplyMessage(vmenv, msg, gp)

	traceActions := txTracer.GetResult()
	state.Finalise(true)

	// err is error occured before EVM execution
	if err != nil {
		errTrace := txtrace.GetErrorTraceFromMsg(&msg, block.Hash, *block.Number, tx.Hash(), index, err)
		at := make([]txtrace.ActionTrace, 0)
		at = append(at, *errTrace)
		// check correct replay state
		if status == 1 {
			return nil, fmt.Errorf("invalid transaction replay state at %s", tx.Hash().String())
		}
		return &at, nil
	}
	// If the timer caused an abort, return an appropriate error message
	if vmenv.Cancelled() {
		log.Debug("EVM was canceled due to timeout when replaying transaction ", "txHash", tx.Hash().String())
		return nil, fmt.Errorf("timeout when replaying tx")
	}

	// result.Err is error during EVM execution
	if result != nil && result.Err != nil {
		if len(*traceActions) == 0 {
			log.Error("error in result when replaying transaction:", "txHash", tx.Hash().String(), " err", result.Err.Error())
			errTrace := txtrace.GetErrorTraceFromMsg(&msg, block.Hash, *block.Number, tx.Hash(), index, result.Err)
			at := make([]txtrace.ActionTrace, 0)
			at = append(at, *errTrace)
			return &at, nil
		}
		// check correct replay state
		if status == 1 {
			return nil, fmt.Errorf("invalid transaction replay state at %s", tx.Hash().String())
		}
		return traceActions, nil
	}

	// check correct replay state
	if status == 0 {
		return nil, fmt.Errorf("invalid transaction replay state at %s", tx.Hash().String())
	}
	return traceActions, nil
}

// getEmptyBlockTrace returns trace for empty block
func getEmptyBlockTrace(blockHash common.Hash, blockNumber big.Int) *[]txtrace.ActionTrace {
	emptyTrace := txtrace.CallTrace{
		Actions: make([]txtrace.ActionTrace, 0),
	}
	blockTrace := txtrace.NewActionTrace(blockHash, blockNumber, common.Hash{}, 0, "empty")
	txAction := txtrace.NewAddressAction(common.Address{}, 0, []byte{}, nil, hexutil.Big{}, nil)
	blockTrace.Action = txAction
	blockTrace.Error = "Empty block"
	emptyTrace.AddTrace(blockTrace)
	return &emptyTrace.Actions
}

// FilterArgs represents the arguments for specifiing trace targets
type FilterArgs struct {
	FromAddress *[]common.Address      `json:"fromAddress"`
	ToAddress   *[]common.Address      `json:"toAddress"`
	FromBlock   *rpc.BlockNumberOrHash `json:"fromBlock"`
	ToBlock     *rpc.BlockNumberOrHash `json:"toBlock"`
	After       uint                   `json:"after"`
	Count       uint                   `json:"count"`
}

// Filter is function for trace_filter rpc call
func (s *PublicTxTraceAPI) Filter(ctx context.Context, args FilterArgs) (*[]txtrace.ActionTrace, error) {
	// add log after execution
	defer func(start time.Time) {

		var data []interface{}
		if args.FromBlock != nil {
			data = append(data, "fromBlock", args.FromBlock.BlockNumber.Int64())
		}
		if args.ToBlock != nil {
			data = append(data, "toBlock", args.ToBlock.BlockNumber.Int64())
		}
		if args.FromAddress != nil {
			adresses := make([]string, 0)
			for _, addr := range *args.FromAddress {
				adresses = append(adresses, addr.String())
			}
			data = append(data, "fromAddr", adresses)
		}
		if args.ToAddress != nil {
			adresses := make([]string, 0)
			for _, addr := range *args.ToAddress {
				adresses = append(adresses, addr.String())
			}
			data = append(data, "toAddr", adresses)
		}
		data = append(data, "time", time.Since(start))
		log.Debug("Executing trace_filter call finished", data...)
	}(time.Now())

	// process arguments
	var (
		fromBlock, toBlock rpc.BlockNumber
		mainErr            error
	)
	if args.FromBlock != nil {
		fromBlock = *args.FromBlock.BlockNumber
	}
	if args.ToBlock != nil {
		toBlock = *args.ToBlock.BlockNumber
		if toBlock == rpc.LatestBlockNumber || toBlock == rpc.PendingBlockNumber {
			toBlock = rpc.BlockNumber(s.b.CurrentBlock().NumberU64())
		}
	} else {
		toBlock = rpc.BlockNumber(s.b.CurrentBlock().NumberU64())
	}

	// counter of processed traces
	var traceAdded, traceCount uint
	var fromAddresses, toAddresses map[common.Address]struct{}
	if args.FromAddress != nil {
		fromAddresses = make(map[common.Address]struct{})
		for _, addr := range *args.FromAddress {
			fromAddresses[addr] = struct{}{}
		}
	}
	if args.ToAddress != nil {
		toAddresses = make(map[common.Address]struct{})
		for _, addr := range *args.ToAddress {
			toAddresses[addr] = struct{}{}
		}
	}

	// check for context timeout
	contextDone := false
	go func() {
		<-ctx.Done()
		contextDone = true
	}()

	// struct for collecting result traces
	callTrace := txtrace.CallTrace{
		Actions: make([]txtrace.ActionTrace, 0),
	}

	// count of traces doesn't matter so use parallel workers
	if args.Count == 0 {
		workerCount := runtime.NumCPU() / 2
		blocks := make(chan rpc.BlockNumber, 1000)
		results := make(chan txtrace.ActionTrace, 100000)

		// create workers and their sync group
		var wg sync.WaitGroup
		for w := 0; w < workerCount; w++ {
			wg.Add(1)
			wId := w
			go func() {
				defer wg.Done()
				worker(wId, s, ctx, blocks, results, fromAddresses, toAddresses)
			}()
		}

		// add all blocks in specified range for processing
		for i := fromBlock; i <= toBlock; i++ {
			blocks <- i
		}
		close(blocks)

		var wgResult sync.WaitGroup
		wgResult.Add(1)
		go func() {
			defer wgResult.Done()
			// collect results
			for trace := range results {
				callTrace.AddTrace(&trace)
			}
		}()

		// wait for proccessing all blocks
		wg.Wait()
		close(results)

		wgResult.Wait()
	} else {
	blocks:
		// go thru all blocks in specified range
		for i := fromBlock; i <= toBlock; i++ {
			block, err := s.b.BlockByNumber(ctx, i)
			if err != nil {
				mainErr = err
				break
			}

			// when block has any transaction, then process it
			if block != nil && block.Transactions.Len() > 0 {
				traces, err := s.replayBlock(ctx, block, nil, nil)
				if err != nil {
					mainErr = err
					break
				}

				// loop thru all traces from the block
				// and check
				for _, trace := range *traces {

					if args.Count == 0 || traceAdded < args.Count {
						addTrace := true

						if args.FromAddress != nil || args.ToAddress != nil {
							if args.FromAddress != nil {
								if trace.Action.From == nil {
									addTrace = false
								} else {
									if _, ok := fromAddresses[*trace.Action.From]; !ok {
										addTrace = false
									}
								}
							}
							if args.ToAddress != nil {
								if trace.Action.To == nil {
									addTrace = false
								} else if _, ok := toAddresses[*trace.Action.To]; !ok {
									addTrace = false
								}
							}
						}
						if addTrace {
							if traceCount >= args.After {
								callTrace.AddTrace(&trace)
								traceAdded++
							}
							traceCount++
						}
					} else {
						// already reached desired count of traces in batch
						break blocks
					}
				}
			}
			if contextDone {
				break
			}
		}
	}

	//when timeout occured or another error
	if contextDone || mainErr != nil {
		if mainErr != nil {
			return nil, mainErr
		}
		return nil, fmt.Errorf("timeout when scanning blocks")
	}

	return &callTrace.Actions, nil
}

func worker(id int,
	s *PublicTxTraceAPI,
	ctx context.Context,
	blocks <-chan rpc.BlockNumber,
	results chan<- txtrace.ActionTrace,
	fromAddresses map[common.Address]struct{},
	toAddresses map[common.Address]struct{}) {

	for i := range blocks {
		block, err := s.b.BlockByNumber(ctx, i)
		if err != nil {
			break
		}

		// when block has any transaction, then process it
		if block != nil && block.Transactions.Len() > 0 {
			traces, err := s.replayBlock(ctx, block, nil, nil)
			if err != nil {
				break
			}
			for _, trace := range *traces {
				addTrace := true

				if len(fromAddresses) > 0 {

					if trace.Action.From == nil {
						addTrace = false
					} else {
						if _, ok := fromAddresses[*trace.Action.From]; !ok {
							addTrace = false
						}
					}
				}
				if len(toAddresses) > 0 {
					if trace.Action.To == nil {
						addTrace = false
					} else if _, ok := toAddresses[*trace.Action.To]; !ok {
						addTrace = false
					}
				}
				if addTrace {
					results <- trace
				}
			}
		}
	}
}
