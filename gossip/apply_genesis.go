package gossip

import (
	"errors"
	"fmt"
	"github.com/Fantom-foundation/go-opera/statedb"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb/batched"
	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
	"github.com/Fantom-foundation/go-opera/inter/ibr"
	"github.com/Fantom-foundation/go-opera/inter/ier"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/utils/dbutil/autocompact"
)

// ApplyGenesis writes initial state.
func (s *Store) ApplyGenesis(g genesis.Genesis) (genesisHash hash.Hash, err error) {
	// use batching wrapper for hot tables
	unwrap := s.WrapTablesAsBatched()
	defer unwrap()

	// write epochs
	var topEr *ier.LlrIdxFullEpochRecord
	g.Epochs.ForEach(func(er ier.LlrIdxFullEpochRecord) bool {
		if er.EpochState.Rules.NetworkID != g.NetworkID || er.EpochState.Rules.Name != g.NetworkName {
			err = errors.New("network ID/name mismatch")
			return false
		}
		if topEr == nil {
			topEr = &er
		}
		s.WriteFullEpochRecord(er)
		return true
	})
	if err != nil {
		return genesisHash, err
	}
	if topEr == nil {
		return genesisHash, errors.New("no ERs in genesis")
	}
	var prevEs *iblockproc.EpochState
	s.ForEachHistoryBlockEpochState(func(bs iblockproc.BlockState, es iblockproc.EpochState) bool {
		s.WriteUpgradeHeight(bs, es, prevEs)
		prevEs = &es
		return true
	})
	s.SetBlockEpochState(topEr.BlockState, topEr.EpochState)
	s.FlushBlockEpochState()

	// write blocks
	var lastBlock idx.Block
	g.Blocks.ForEach(func(br ibr.LlrIdxFullBlockRecord) bool {
		s.WriteFullBlockRecord(br)
		if br.Idx > lastBlock {
			lastBlock = br.Idx
		}
		return true
	})

	// write EVM items
	err = s.evm.ApplyGenesis(g)
	if err != nil {
		return genesisHash, err
	}

	// write EVM items into Carmen
	if statedb.IsExternalStateDbUsed() {
		for blockNum := idx.Block(1); blockNum <= lastBlock; blockNum++ {
			block := s.GetBlock(blockNum)
			if block == nil {
				s.Log.Trace("Skipping missing block from import into StateDB", "index", blockNum)
				continue
			}
			err = statedb.ImportTrieIntoStateDb(s.evm.EvmDb, s.evm.EVMDB(), uint64(blockNum), common.Hash(block.Root))
			if err != nil {
				return genesisHash, fmt.Errorf("genesis import into StateDB failed at block %d; %v", blockNum, err)
			}
			stateHash := statedb.GetLiveStateHash()
			if common.Hash(block.Root) == stateHash {
				s.Log.Info("Imported block into StateDB", "index", blockNum, "root", block.Root)
			} else {
				s.Log.Warn("Imported block into StateDB with not-matching state hash", "index", blockNum, "root", block.Root, "realHash", stateHash)
			}
		}
	}

	// write LLR state
	s.setLlrState(LlrState{
		LowestEpochToDecide: topEr.Idx + 1,
		LowestEpochToFill:   topEr.Idx + 1,
		LowestBlockToDecide: topEr.BlockState.LastBlock.Idx + 1,
		LowestBlockToFill:   topEr.BlockState.LastBlock.Idx + 1,
	})
	s.FlushLlrState()

	s.SetGenesisID(g.GenesisID)
	s.SetGenesisBlockIndex(topEr.BlockState.LastBlock.Idx)

	return genesisHash, err
}

func (s *Store) WrapTablesAsBatched() (unwrap func()) {
	origTables := s.table

	batchedBlocks := batched.Wrap(autocompact.Wrap2M(s.table.Blocks, opt.GiB, 16*opt.GiB, false, "blocks"))
	s.table.Blocks = batchedBlocks

	batchedBlockHashes := batched.Wrap(s.table.BlockHashes)
	s.table.BlockHashes = batchedBlockHashes

	unwrapEVM := s.evm.WrapTablesAsBatched()
	return func() {
		unwrapEVM()
		_ = batchedBlocks.Flush()
		_ = batchedBlockHashes.Flush()
		s.table = origTables
	}
}
