package gossip

import (
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
	return &ibr.LlrFullBlockRecord{
		Atropos:  block.Atropos,
		Root:     block.Root,
		Txs:      txs,
		Receipts: receipts,
		Time:     block.Time,
		GasUsed:  block.GasUsed,
	}
}

func (s *Store) GetBlockRecordHash(n idx.Block) *hash.Hash {
	// Get data from LRU cache first.
	if s.cache.BRHashes != nil {
		if c, ok := s.cache.BRHashes.Get(n); ok {
			h := c.(hash.Hash)
			return &h
		}
	}
	br := s.GetFullBlockRecord(n)
	if br == nil {
		return nil
	}
	brHash := br.Hash()
	// Add to LRU cache.
	s.cache.BRHashes.Add(n, brHash, nominalSize)
	return &brHash
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
