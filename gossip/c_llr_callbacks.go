package gossip

import (
	"math/big"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"

	"github.com/Fantom-foundation/go-opera/gossip/evmstore"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
	"github.com/Fantom-foundation/go-opera/inter/ibr"
	"github.com/Fantom-foundation/go-opera/inter/ier"
	"github.com/Fantom-foundation/go-opera/opera"
)

func indexRawReceipts(s *Store, receiptsForStorage []*types.ReceiptForStorage, txs types.Transactions, blockIdx idx.Block, atropos hash.Event, config *params.ChainConfig, time uint64, baseFee *big.Int, blobGasPrice *big.Int) {
	s.evm.SetRawReceipts(blockIdx, receiptsForStorage)
	receipts, _ := evmstore.UnwrapStorageReceipts(receiptsForStorage, blockIdx, config, common.Hash(atropos), time, baseFee, blobGasPrice, txs)
	for _, r := range receipts {
		s.evm.IndexLogs(r.Logs...)
	}
}

func (s *Store) WriteFullBlockRecord(baseFee *big.Int, blobGasPrice *big.Int, gasLimit uint64, br ibr.LlrIdxFullBlockRecord) {
	for _, tx := range br.Txs {
		s.EvmStore().SetTx(tx.Hash(), tx)
	}

	if len(br.Receipts) != 0 {
		// Note: it's possible for receipts to get indexed twice by BR and block processing
		indexRawReceipts(s, br.Receipts, br.Txs, br.Idx, hash.Event(br.BlockHash), s.GetEvmChainConfig(), uint64(br.Time.Unix()), baseFee, blobGasPrice)
	}
	for i, tx := range br.Txs {
		s.EvmStore().SetTx(tx.Hash(), tx)
		s.EvmStore().SetTxPosition(tx.Hash(), evmstore.TxPosition{
			Block:       br.Idx,
			BlockOffset: uint32(i),
		})
	}

	parentHash := common.Hash{}
	if parent := s.GetBlock(br.Idx - 1); parent != nil {
		parentHash = parent.Hash()
	}

	// TODO: add bloom log and other fields
	builder := inter.NewBlockBuilder().
		SetNumber(uint64(br.Idx)).
		SetTime(br.Time).
		SetParentHash(parentHash).
		SetStateRoot(common.Hash(br.StateRoot)).
		SetGasLimit(gasLimit).
		SetGasUsed(br.GasUsed).
		SetBaseFee(baseFee)

	for i := range br.Txs {
		copy := types.Receipt(*br.Receipts[i])
		builder.AddTransaction(br.Txs[i], &copy)
	}

	block := builder.Build()
	s.SetBlock(br.Idx, block)
	s.SetBlockIndex(block.Hash(), br.Idx)
}

func (s *Store) WriteFullEpochRecord(er ier.LlrIdxFullEpochRecord) {
	s.SetHistoryBlockEpochState(er.Idx, er.BlockState, er.EpochState)
	s.SetEpochBlock(er.BlockState.LastBlock.Idx+1, er.Idx)
}

func (s *Store) WriteUpgradeHeight(bs iblockproc.BlockState, es iblockproc.EpochState, prevEs *iblockproc.EpochState) {
	if prevEs == nil || es.Rules.Upgrades != prevEs.Rules.Upgrades {
		s.AddUpgradeHeight(opera.UpgradeHeight{
			Upgrades: es.Rules.Upgrades,
			Height:   bs.LastBlock.Idx + 1,
			Time:     bs.LastBlock.Time + 1,
		})
	}
}
