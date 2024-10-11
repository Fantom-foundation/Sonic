// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evmcore

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math"
	"math/big"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/opera"
)

type (
	EvmHeader struct {
		Number     *big.Int
		Hash       common.Hash
		ParentHash common.Hash
		Root       common.Hash
		TxHash     common.Hash
		Time       inter.Timestamp
		Coinbase   common.Address

		GasLimit uint64
		GasUsed  uint64

		BaseFee *big.Int

		Difficulty  *big.Int
		MixDigest   common.Hash
	}

	EvmBlock struct {
		EvmHeader

		Transactions types.Transactions
	}
)

// NewEvmBlock constructor.
func NewEvmBlock(h *EvmHeader, txs types.Transactions) *EvmBlock {
	b := &EvmBlock{
		EvmHeader:    *h,
		Transactions: txs,
	}

	if len(txs) == 0 {
		b.EvmHeader.TxHash = types.EmptyRootHash
	} else {
		b.EvmHeader.TxHash = types.DeriveSha(txs, trie.NewStackTrie(nil))
	}

	return b
}

// ToEvmHeader converts inter.Block to EvmHeader.
func ToEvmHeader(block *inter.Block, index idx.Block, prevHash hash.Event, rules opera.Rules) *EvmHeader {
	baseFee := rules.Economy.MinGasPrice
	if !rules.Upgrades.London {
		baseFee = nil
	}

	difficulty := new(big.Int)
	if !rules.Upgrades.Sonic {
		difficulty.SetUint64(1)
	}

	return &EvmHeader{
		Hash:       common.Hash(block.Atropos),
		ParentHash: common.Hash(prevHash),
		Root:       common.Hash(block.Root),
		Number:     big.NewInt(int64(index)),
		Time:       block.Time,
		GasLimit:   math.MaxUint64,
		GasUsed:    block.GasUsed,
		BaseFee:    baseFee,
		Difficulty: difficulty,
		MixDigest:  common.Hash{}, // TODO provide pseudorandom data?
	}
}

// ConvertFromEthHeader converts ETH-formatted header to Lachesis EVM header
func ConvertFromEthHeader(h *types.Header) *EvmHeader {
	// NOTE: incomplete conversion
	return &EvmHeader{
		Number:     h.Number,
		Coinbase:   h.Coinbase,
		GasLimit:   math.MaxUint64,
		GasUsed:    h.GasUsed,
		Root:       h.Root,
		TxHash:     h.TxHash,
		ParentHash: h.ParentHash,
		Time:       inter.FromUnix(int64(h.Time)),
		Hash:       common.BytesToHash(h.Extra),
		BaseFee:    h.BaseFee,
		Difficulty: h.Difficulty,
		MixDigest:  h.MixDigest,
	}
}

// EthHeader returns header in ETH format
func (h *EvmHeader) EthHeader() *types.Header {
	if h == nil {
		return nil
	}
	// NOTE: incomplete conversion
	ethHeader := &types.Header{
		Number:     h.Number,
		Coinbase:   h.Coinbase,
		GasLimit:   0xffffffffffff, // don't use h.GasLimit (too much bits) here to avoid parsing issues
		GasUsed:    h.GasUsed,
		Root:       h.Root,
		TxHash:     h.TxHash,
		ParentHash: h.ParentHash,
		Time:       uint64(h.Time.Unix()),
		Extra:      h.Hash.Bytes(),
		BaseFee:    h.BaseFee,

		Difficulty: new(big.Int),
		MixDigest:  h.MixDigest,
	}
	// ethHeader.SetExternalHash(h.Hash) < this seems to be an optimization in go-ethereum-substate; skipped for now, needs investigation
	return ethHeader
}

// EvmHeaderJson is simplified version of types.Header, but allowing setting custom hash
type EvmHeaderJson struct {
	ParentHash  common.Hash      `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash      `json:"sha3Uncles"       gencodec:"required"`
	Miner       common.Address   `json:"miner"`
	Root        common.Hash      `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash      `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash      `json:"receiptsRoot"     gencodec:"required"`
	Bloom       types.Bloom      `json:"logsBloom"        gencodec:"required"`
	Difficulty  *hexutil.Big     `json:"difficulty"       gencodec:"required"`
	Number      *hexutil.Big     `json:"number"           gencodec:"required"`
	GasLimit    hexutil.Uint64   `json:"gasLimit"         gencodec:"required"`
	GasUsed     hexutil.Uint64   `json:"gasUsed"          gencodec:"required"`
	Time        hexutil.Uint64   `json:"timestamp"        gencodec:"required"`
	TimeNano    hexutil.Uint64   `json:"timestampNano"`
	Extra       hexutil.Bytes    `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash      `json:"mixHash"`
	Nonce       types.BlockNonce `json:"nonce"`
	BaseFee     *hexutil.Big     `json:"baseFeePerGas"`
	Hash        *common.Hash     `json:"hash"`
	Epoch       hexutil.Uint64   `json:"epoch"`
	TotalDiff   *hexutil.Big     `json:"totalDifficulty"`
}

type EvmBlockJson struct {
	*EvmHeaderJson
	Txs         []interface{}    `json:"transactions"`
	Size        *hexutil.Uint64  `json:"size"` // RLP encoded storage size of the block
	Uncles      []common.Hash    `json:"uncles"`
}

func (h *EvmHeader) ToJson(receipts types.Receipts) *EvmHeaderJson {
	enc := &EvmHeaderJson{
		Number:     (*hexutil.Big)(h.Number),
		Miner:      h.Coinbase,
		GasLimit:   0xffffffffffff, // don't use h.GasLimit (too much bits) here to avoid parsing issues
		GasUsed:    hexutil.Uint64(h.GasUsed),
		Root:       h.Root,
		TxHash:     h.TxHash,
		ParentHash: h.ParentHash,
		UncleHash:  types.EmptyUncleHash,
		Time:       hexutil.Uint64(h.Time.Unix()),
		TimeNano:   hexutil.Uint64(h.Time),
		BaseFee:    (*hexutil.Big)(h.BaseFee),
		Difficulty: (*hexutil.Big)(h.Difficulty),
		MixDigest:  h.MixDigest,
		TotalDiff:  new(hexutil.Big),
		Hash:       &h.Hash,
		Epoch:		hexutil.Uint64(hash.Event(h.Hash).Epoch()),
	}
	if receipts != nil { // if receipts resolution fails, don't set ReceiptsHash at all
		if receipts.Len() != 0 {
			enc.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
			enc.Bloom = types.CreateBloom(receipts)
		} else {
			enc.ReceiptHash = types.EmptyRootHash
		}
	}
	return enc
}

// Header is a copy of EvmBlock.EvmHeader.
func (b *EvmBlock) Header() *EvmHeader {
	if b == nil {
		return nil
	}
	// copy values
	h := b.EvmHeader
	// copy refs
	h.Number = new(big.Int).Set(b.Number)
	if b.BaseFee != nil {
		h.BaseFee = new(big.Int).Set(b.BaseFee)
	}

	return &h
}

func (b *EvmBlock) NumberU64() uint64 {
	return b.Number.Uint64()
}

func (b *EvmBlock) EthBlock() *types.Block {
	if b == nil {
		return nil
	}
	body := types.Body{Transactions: b.Transactions}
	return types.NewBlock(b.EvmHeader.EthHeader(), &body, nil, trie.NewStackTrie(nil))
}

func (b *EvmBlock) EstimateSize() int {
	est := 0
	for _, tx := range b.Transactions {
		est += len(tx.Data())
	}
	return est + b.Transactions.Len()*256
}
