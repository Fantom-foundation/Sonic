package inter

import (
	"math/big"
	"unsafe"

	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/holiman/uint256"
)

// Block represents the on-disk storage format of a block. It contains all
// fields required to reconstruct the block header, as well as a list of
// hashes of the transactions been executed as part of the represented block.
//
// This struct should be considered immutable. No fields should be modified,
// directly or indirectly. Ideally, all fields should be private, but that
// would invalidate support for RLP encoding as it is used to store instances
// on the disk. However, future updates may make fields inaccessible.
//
// To create a new block, use the BlockBuilder, handling the computation of
// key properties implicitly.
type Block struct {
	// Fields required for the block header.
	Number               uint64
	ParentHash           common.Hash
	StateRoot            common.Hash
	Time                 Timestamp
	Difficulty           uint64
	GasLimit             uint64
	GasUsed              uint64
	BaseFee              uint256.Int
	PrevRandao           common.Hash
	TransactionsHashRoot common.Hash
	ReceiptsHashRoot     common.Hash
	LogBloom             types.Bloom

	// Fields required for linking blocks to contained transactions.
	TransactionHashes []common.Hash

	// Fields required for linking the block internally to a lachesis epoch.
	Epoch idx.Epoch

	// The hash of this block, cached on first access.
	hash common.Hash
}

// Hash computes the hash of this block, committing all its fields.
func (b *Block) Hash() common.Hash {
	if b.hash == (common.Hash{}) {
		b.hash = b.GetEthereumHeader().Hash()
	}
	return b.hash
}

// uncleHash is the hash to be used for the uncle field in Ethereum headers if
// there are no uncles. See https://eips.ethereum.org/EIPS/eip-4844.
var uncleHash = common.BytesToHash(hexutil.MustDecode(
	"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
))

// GetEthereumHeader returns the Ethereum header corresponding to this block.
func (b *Block) GetEthereumHeader() *types.Header {
	return &types.Header{
		ParentHash:  b.ParentHash,
		UncleHash:   uncleHash,
		Coinbase:    common.Address{}, // < in Sonic, the coinbase is always 0
		Root:        b.StateRoot,
		TxHash:      b.TransactionsHashRoot,
		ReceiptHash: b.ReceiptsHashRoot,
		Bloom:       b.LogBloom,
		Difficulty:  big.NewInt(int64(b.Difficulty)),
		Number:      big.NewInt(int64(b.Number)),
		GasLimit:    b.GasLimit,
		GasUsed:     b.GasUsed,
		Time:        uint64(b.Time.Time().Unix()),
		Extra:       nil, // TODO: fill in extra data required for gas computation
		MixDigest:   b.PrevRandao,
		Nonce:       types.BlockNonce{}, // constant 0 in Ethereum
		BaseFee:     b.BaseFee.ToBig(),

		// Sonic does not have a beacon chain and no withdrawals.
		WithdrawalsHash: &types.EmptyWithdrawalsHash,

		// Sonic does not support blobs, so no blob gas is used and there is
		// no excess blob gas.
		BlobGasUsed:   new(uint64), // = 0
		ExcessBlobGas: new(uint64), // = 0
	}
}

func (b *Block) EstimateSize() int {
	return int(unsafe.Sizeof(*b)) +
		len(b.TransactionHashes)*int(unsafe.Sizeof(common.Hash{}))
}

// ----------------------------------------------------------------------------
// BlockBuilder
// ----------------------------------------------------------------------------

type BlockBuilder struct {
	block        Block
	transactions types.Transactions
	receipts     types.Receipts
}

func NewBlockBuilder() *BlockBuilder {
	return &BlockBuilder{}
}

func (b *BlockBuilder) SetNumber(number uint64) *BlockBuilder {
	b.block.Number = number
	return b
}

func (b *BlockBuilder) SetParentHash(hash common.Hash) *BlockBuilder {
	b.block.ParentHash = hash
	return b
}

func (b *BlockBuilder) SetStateRoot(hash common.Hash) *BlockBuilder {
	b.block.StateRoot = hash
	return b
}

func (b *BlockBuilder) AddTransaction(
	transaction *types.Transaction,
	receipt *types.Receipt,
) *BlockBuilder {
	b.transactions = append(b.transactions, transaction)
	b.receipts = append(b.receipts, receipt)
	return b
}

func (b *BlockBuilder) SetTime(time Timestamp) *BlockBuilder {
	b.block.Time = time
	return b
}

func (b *BlockBuilder) SetDifficulty(difficulty uint64) *BlockBuilder {
	b.block.Difficulty = difficulty
	return b
}

func (b *BlockBuilder) SetGasLimit(gasLimit uint64) *BlockBuilder {
	b.block.GasLimit = gasLimit
	return b
}

func (b *BlockBuilder) SetGasUsed(gasUsed uint64) *BlockBuilder {
	b.block.GasUsed = gasUsed
	return b
}

func (b *BlockBuilder) SetBaseFee(baseFee uint256.Int) *BlockBuilder {
	b.block.BaseFee = baseFee
	return b
}

func (b *BlockBuilder) SetPrevRandao(prevRandao common.Hash) *BlockBuilder {
	b.block.PrevRandao = prevRandao
	return b
}

func (b *BlockBuilder) SetEpoch(epoch idx.Epoch) *BlockBuilder {
	b.block.Epoch = epoch
	return b
}

func (b *BlockBuilder) Build() *Block {
	res := new(Block)
	*res = b.block

	res.TransactionsHashRoot = types.DeriveSha(
		b.transactions,
		trie.NewStackTrie(nil),
	)
	res.ReceiptsHashRoot = types.DeriveSha(b.receipts, trie.NewStackTrie(nil))
	res.LogBloom = types.CreateBloom(b.receipts)

	for _, tx := range b.transactions {
		res.TransactionHashes = append(res.TransactionHashes, tx.Hash())
	}

	return res
}
