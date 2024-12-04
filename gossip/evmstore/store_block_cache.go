package evmstore

import (
	"github.com/Fantom-foundation/lachesis-base/ltypes"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-opera/evmcore"
)

func (s *Store) GetCachedEvmBlock(n ltypes.BlockID) *evmcore.EvmBlock {
	c, ok := s.cache.EvmBlocks.Get(n)
	if !ok {
		return nil
	}

	return c.(*evmcore.EvmBlock)
}

func (s *Store) SetCachedEvmBlock(n ltypes.BlockID, b *evmcore.EvmBlock) {
	var empty = common.Hash{}
	if b.EvmHeader.TxHash == empty {
		panic("You have to cache only completed blocks (with txs)")
	}
	s.cache.EvmBlocks.Add(n, b, uint(b.EstimateSize()))
}
