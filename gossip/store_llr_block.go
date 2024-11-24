package gossip

import (
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/inter/ibr"
	"github.com/Fantom-foundation/go-opera/inter/ier"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/core/types"
)

func (s *Store) GetFullBlockRecord(n idx.Block) *ibr.LlrFullBlockRecord {
	block := s.GetBlock(n)
	if block == nil {
		return nil
	}
	txs := s.GetBlockTxs(n, block)
	receipts, _ := s.EvmStore().GetRawReceipts(n)
	if receipts == nil {
		receipts = []*types.ReceiptForStorage{}
	}
	return FullBlockRecordFor(block, txs, receipts)
}

func (s *Store) GetFullEpochRecord(epoch idx.Epoch) *ier.LlrFullEpochRecord {
	hbs, hes := s.GetHistoryBlockEpochState(epoch)
	if hbs == nil || hes == nil {
		return nil
	}
	return &ier.LlrFullEpochRecord{
		BlockState: *hbs,
		EpochState: *hes,
	}
}

// FullBlockRecordFor returns the full block record used in Genesis processing
// for the given block, list of transactions, and list of transaction receipts.
func FullBlockRecordFor(block *inter.Block, txs types.Transactions,
	rawReceipts []*types.ReceiptForStorage) *ibr.LlrFullBlockRecord {
	return &ibr.LlrFullBlockRecord{
		BlockHash:  hash.Hash(block.Hash()),
		ParentHash: hash.Hash(block.ParentHash),
		StateRoot:  hash.Hash(block.StateRoot),
		Time:       block.Time,
		Duration:   block.Duration,
		Difficulty: block.Difficulty,
		GasLimit:   block.GasLimit,
		GasUsed:    block.GasUsed,
		BaseFee:    block.BaseFee,
		PrevRandao: hash.Hash(block.PrevRandao),
		Epoch:      block.Epoch,
		Txs:        txs,
		Receipts:   rawReceipts,
	}
}
