package makegenesis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/drivermodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/eventmodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/evmmodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/sealmodule"
	"github.com/Fantom-foundation/go-opera/gossip/evmstore"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
	"github.com/Fantom-foundation/go-opera/inter/ibr"
	"github.com/Fantom-foundation/go-opera/inter/ier"
	"github.com/Fantom-foundation/go-opera/inter/state"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"

	mptIo "github.com/Fantom-foundation/Carmen/go/database/mpt/io"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
)

type GenesisBuilder struct {
	tmpStateDB state.StateDB
	carmenDir  string
	carmenStateDb carmen.StateDB

	totalSupply *big.Int

	blocks       []ibr.LlrIdxFullBlockRecord
	epochs       []ier.LlrIdxFullEpochRecord
	currentEpoch ier.LlrIdxFullEpochRecord
}

type BlockProc struct {
	SealerModule     blockproc.SealerModule
	TxListenerModule blockproc.TxListenerModule
	PreTxTransactor  blockproc.TxTransactor
	PostTxTransactor blockproc.TxTransactor
	EventsModule     blockproc.ConfirmedEventsModule
	EVMModule        blockproc.EVM
}

func DefaultBlockProc() BlockProc {
	return BlockProc{
		SealerModule:     sealmodule.New(),
		TxListenerModule: drivermodule.NewDriverTxListenerModule(),
		PreTxTransactor:  drivermodule.NewDriverTxPreTransactor(),
		PostTxTransactor: drivermodule.NewDriverTxTransactor(),
		EventsModule:     eventmodule.New(),
		EVMModule:        evmmodule.New(),
	}
}

func (b *GenesisBuilder) AddBalance(acc common.Address, balance *big.Int) {
	b.tmpStateDB.AddBalance(acc, balance)
	b.totalSupply.Add(b.totalSupply, balance)
}

func (b *GenesisBuilder) SetCode(acc common.Address, code []byte) {
	b.tmpStateDB.SetCode(acc, code)
}

func (b *GenesisBuilder) SetNonce(acc common.Address, nonce uint64) {
	b.tmpStateDB.SetNonce(acc, nonce)
}

func (b *GenesisBuilder) SetStorage(acc common.Address, key, val common.Hash) {
	b.tmpStateDB.SetState(acc, key, val)
}

func (b *GenesisBuilder) AddBlock(br ibr.LlrIdxFullBlockRecord) {
	b.blocks = append(b.blocks, br)
}

func (b *GenesisBuilder) AddEpoch(er ier.LlrIdxFullEpochRecord) {
	b.epochs = append(b.epochs, er)
}

func (b *GenesisBuilder) SetCurrentEpoch(er ier.LlrIdxFullEpochRecord) {
	b.currentEpoch = er
}

func (b *GenesisBuilder) TotalSupply() *big.Int {
	return b.totalSupply
}

func (b *GenesisBuilder) CurrentHash() hash.Hash {
	er := b.epochs[len(b.epochs)-1]
	return er.Hash()
}

func NewGenesisBuilder() *GenesisBuilder {
	carmenDir, err := os.MkdirTemp("", "opera-tmp-genesis")
	if err != nil {
		panic(fmt.Errorf("failed to create temporary dir for GenesisBuilder: %v", err))
	}
	carmenState, err := carmen.NewState(carmen.Parameters{
		Variant:   "go-file",
		Schema:    carmen.Schema(5),
		Directory: carmenDir,
		LiveCache: 1, // use minimum cache (not default)
	})
	if err != nil {
		panic(fmt.Errorf("failed to create carmen state; %s", err))
	}
	carmenStateDb := carmen.CreateStateDBUsing(carmenState)
	tmpStateDB := evmstore.CreateCarmenStateDb(carmenStateDb)
	return &GenesisBuilder{
		tmpStateDB:    tmpStateDB,
		carmenDir:     carmenDir,
		carmenStateDb: carmenStateDb,
		totalSupply:   new(big.Int),
	}
}

type dummyHeaderReturner struct {
}

func (d dummyHeaderReturner) GetHeader(common.Hash, uint64) *evmcore.EvmHeader {
	return &evmcore.EvmHeader{}
}

func (b *GenesisBuilder) ExecuteGenesisTxs(blockProc BlockProc, genesisTxs types.Transactions) error {
	bs, es := b.currentEpoch.BlockState.Copy(), b.currentEpoch.EpochState.Copy()

	blockCtx := iblockproc.BlockCtx{
		Idx:     bs.LastBlock.Idx + 1,
		Time:    bs.LastBlock.Time + 1,
		Atropos: hash.Event{},
	}

	sealer := blockProc.SealerModule.Start(blockCtx, bs, es)
	sealing := true
	txListener := blockProc.TxListenerModule.Start(blockCtx, bs, es, b.tmpStateDB)
	evmProcessor := blockProc.EVMModule.Start(blockCtx, b.tmpStateDB, dummyHeaderReturner{}, func(l *types.Log) {
		txListener.OnNewLog(l)
	}, es.Rules, es.Rules.EvmChainConfig([]opera.UpgradeHeight{
		{
			Upgrades: es.Rules.Upgrades,
			Height:   0,
		},
	}))

	// Execute genesis transactions
	evmProcessor.Execute(genesisTxs)
	bs = txListener.Finalize()

	// Execute pre-internal transactions
	preInternalTxs := blockProc.PreTxTransactor.PopInternalTxs(blockCtx, bs, es, sealing, b.tmpStateDB)
	evmProcessor.Execute(preInternalTxs)
	bs = txListener.Finalize()

	// Seal epoch if requested
	if sealing {
		sealer.Update(bs, es)
		bs, es = sealer.SealEpoch()
		txListener.Update(bs, es)
	}

	// Execute post-internal transactions
	internalTxs := blockProc.PostTxTransactor.PopInternalTxs(blockCtx, bs, es, sealing, b.tmpStateDB)
	evmProcessor.Execute(internalTxs)

	evmBlock, skippedTxs, receipts := evmProcessor.Finalize()
	for _, r := range receipts {
		if r.Status == 0 {
			return errors.New("genesis transaction reverted")
		}
	}
	if len(skippedTxs) != 0 {
		return errors.New("genesis transaction is skipped")
	}
	bs = txListener.Finalize()
	bs.FinalizedStateRoot = hash.Hash(evmBlock.Root)

	bs.LastBlock = blockCtx

	prettyHash := func(root hash.Hash) hash.Event {
		e := inter.MutableEventPayload{}
		// for nice-looking ID
		e.SetEpoch(es.Epoch)
		e.SetLamport(1)
		// actual data hashed
		e.SetExtra(root[:])

		return e.Build().ID()
	}
	receiptsStorage := make([]*types.ReceiptForStorage, len(receipts))
	for i, r := range receipts {
		receiptsStorage[i] = (*types.ReceiptForStorage)(r)
	}
	// add block
	b.blocks = append(b.blocks, ibr.LlrIdxFullBlockRecord{
		LlrFullBlockRecord: ibr.LlrFullBlockRecord{
			Atropos:  prettyHash(bs.FinalizedStateRoot),
			Root:     bs.FinalizedStateRoot,
			Txs:      evmBlock.Transactions,
			Receipts: receiptsStorage,
			Time:     blockCtx.Time,
			GasUsed:  evmBlock.GasUsed,
		},
		Idx: blockCtx.Idx,
	})
	// add epoch
	b.currentEpoch = ier.LlrIdxFullEpochRecord{
		LlrFullEpochRecord: ier.LlrFullEpochRecord{
			BlockState: bs,
			EpochState: es,
		},
		Idx: es.Epoch,
	}
	b.epochs = append(b.epochs, b.currentEpoch)
	return nil
}

type memFile struct {
	*bytes.Buffer
}

func (f *memFile) Close() error {
	*f = memFile{}
	return nil
}

func (b *GenesisBuilder) Build(head genesis.Header) *genesisstore.Store {
	err := b.carmenStateDb.Close()
	if err != nil {
		panic(fmt.Errorf("failed to close genesis carmen state; %s", err))
	}
	return genesisstore.NewStore(func(name string) (io.Reader, error) {
		buf := &memFile{bytes.NewBuffer(nil)}
		if name == genesisstore.BlocksSection(0) {
			for i := len(b.blocks) - 1; i >= 0; i-- {
				_ = rlp.Encode(buf, b.blocks[i])
			}
			return buf, nil
		}
		if name == genesisstore.EpochsSection(0) {
			for i := len(b.epochs) - 1; i >= 0; i-- {
				_ = rlp.Encode(buf, b.epochs[i])
			}
			return buf, nil
		}
		if name == genesisstore.FwsLiveSection(0) {
			err := mptIo.Export(context.Background(), filepath.Join(b.carmenDir, "live"), buf)
			if err != nil {
				return nil, err
			}
		}
		if buf.Len() == 0 {
			return nil, errors.New("not found")
		}
		return buf, nil
	}, head, func() error {
		err := os.RemoveAll(b.carmenDir)
		*b = GenesisBuilder{}
		return err
	})
}
