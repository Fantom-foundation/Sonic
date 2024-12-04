package ibr

import (
	"math/big"

	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/lachesis-base/ltypes"
	"github.com/ethereum/go-ethereum/core/types"
)

type LlrFullBlockRecord struct {
	BlockHash  ltypes.Hash
	ParentHash ltypes.Hash
	StateRoot  ltypes.Hash
	Time       inter.Timestamp
	Duration   uint64
	Difficulty uint64
	GasLimit   uint64
	GasUsed    uint64
	BaseFee    *big.Int
	PrevRandao ltypes.Hash
	Epoch      ltypes.EpochID
	Txs        types.Transactions
	Receipts   []*types.ReceiptForStorage
}

type LlrIdxFullBlockRecord struct {
	LlrFullBlockRecord
	Idx ltypes.BlockID
}

// FullBlockRecordFor returns the full block record used in Genesis processing
// for the given block, list of transactions, and list of transaction receipts.
func FullBlockRecordFor(block *inter.Block, txs types.Transactions,
	rawReceipts []*types.ReceiptForStorage) *LlrFullBlockRecord {
	return &LlrFullBlockRecord{
		BlockHash:  ltypes.Hash(block.Hash()),
		ParentHash: ltypes.Hash(block.ParentHash),
		StateRoot:  ltypes.Hash(block.StateRoot),
		Time:       block.Time,
		Duration:   block.Duration,
		Difficulty: block.Difficulty,
		GasLimit:   block.GasLimit,
		GasUsed:    block.GasUsed,
		BaseFee:    block.BaseFee,
		PrevRandao: ltypes.Hash(block.PrevRandao),
		Epoch:      block.Epoch,
		Txs:        txs,
		Receipts:   rawReceipts,
	}
}
