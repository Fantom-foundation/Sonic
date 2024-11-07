package evmmodule

import (
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
	"github.com/Fantom-foundation/go-opera/inter/state"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/Fantom-foundation/go-opera/utils"
)

type EVMModule struct{}

func New() *EVMModule {
	return &EVMModule{}
}

func (p *EVMModule) Start(block iblockproc.BlockCtx, statedb state.StateDB, reader evmcore.DummyChain, onNewLog func(*types.Log), net opera.Rules, evmCfg *params.ChainConfig) blockproc.EVMProcessor {
	var prevBlockHash common.Hash
	var baseFee *big.Int
	if block.Idx <= 1 {  // < the genesis block is block 1
		baseFee = big.NewInt(1e9) // < TODO: make configurable
	} else {
		header := reader.GetHeader(common.Hash{}, uint64(block.Idx-1))
		prevBlockHash = header.Hash

		// TODO: compute base-fee for the next block
		baseFee = new(big.Int).Add(header.BaseFee, big.NewInt(1e8)) // < TODO: implement actual gas-price calculation
	}

	// Start block
	statedb.BeginBlock(uint64(block.Idx))

	return &OperaEVMProcessor{
		block:         block,
		reader:        reader,
		statedb:       statedb,
		onNewLog:      onNewLog,
		net:           net,
		evmCfg:        evmCfg,
		blockIdx:      utils.U64toBig(uint64(block.Idx)),
		prevBlockHash: prevBlockHash,
		gasBaseFee:    baseFee,
	}
}

type OperaEVMProcessor struct {
	block    iblockproc.BlockCtx
	reader   evmcore.DummyChain
	statedb  state.StateDB
	onNewLog func(*types.Log)
	net      opera.Rules
	evmCfg   *params.ChainConfig

	blockIdx      *big.Int
	prevBlockHash common.Hash
	gasBaseFee    *big.Int

	gasUsed uint64

	incomingTxs types.Transactions
	skippedTxs  []uint32
	receipts    types.Receipts
}

func (p *OperaEVMProcessor) evmBlockWith(txs types.Transactions) *evmcore.EvmBlock {
	baseFee := p.net.Economy.MinGasPrice
	if !p.net.Upgrades.London {
		baseFee = nil
	} else if p.net.Upgrades.Sonic {
		baseFee = p.gasBaseFee
	}

	prevRandao := common.Hash{}
	if p.net.Upgrades.Sonic {
		prevRandao.SetBytes([]byte{1}) // TODO provide pseudorandom data?
	}
	h := &evmcore.EvmHeader{
		Number:     p.blockIdx,
		Hash:       common.Hash(p.block.Atropos),
		ParentHash: p.prevBlockHash,
		Root:       common.Hash{},
		Time:       p.block.Time,
		Coinbase:   common.Address{},
		GasLimit:   math.MaxUint64,
		GasUsed:    p.gasUsed,
		BaseFee:    baseFee,
		PrevRandao: prevRandao,
	}

	return evmcore.NewEvmBlock(h, txs)
}

func (p *OperaEVMProcessor) Execute(txs types.Transactions) types.Receipts {
	evmProcessor := evmcore.NewStateProcessor(p.evmCfg, p.reader)
	txsOffset := uint(len(p.incomingTxs))

	// Process txs
	evmBlock := p.evmBlockWith(txs)
	receipts, _, skipped, err := evmProcessor.Process(evmBlock, p.statedb, opera.DefaultVMConfig, &p.gasUsed, func(l *types.Log) {
		// Note: l.Index is properly set before
		l.TxIndex += txsOffset
		p.onNewLog(l)
	})
	if err != nil {
		log.Crit("EVM internal error", "err", err)
	}

	if txsOffset > 0 {
		for i, n := range skipped {
			skipped[i] = n + uint32(txsOffset)
		}
		for _, r := range receipts {
			r.TransactionIndex += txsOffset
		}
	}

	p.incomingTxs = append(p.incomingTxs, txs...)
	p.skippedTxs = append(p.skippedTxs, skipped...)
	p.receipts = append(p.receipts, receipts...)

	return receipts
}

func (p *OperaEVMProcessor) Finalize() (evmBlock *evmcore.EvmBlock, skippedTxs []uint32, receipts types.Receipts) {
	evmBlock = p.evmBlockWith(
		// Filter skipped transactions. Receipts are filtered already
		inter.FilterSkippedTxs(p.incomingTxs, p.skippedTxs),
	)
	skippedTxs = p.skippedTxs
	receipts = p.receipts

	// Commit block
	p.statedb.EndBlock(evmBlock.Number.Uint64())

	// Get state root
	evmBlock.Root = p.statedb.GetStateHash()

	return
}
