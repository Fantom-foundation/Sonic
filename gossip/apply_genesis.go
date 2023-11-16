package gossip

import (
	"errors"
	"fmt"
	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
	"github.com/Fantom-foundation/go-opera/inter/ibr"
	"github.com/Fantom-foundation/go-opera/inter/ier"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/statedb"
	"github.com/Fantom-foundation/go-opera/utils/dbutil/autocompact"
	"github.com/Fantom-foundation/lachesis-base/kvdb/batched"
	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// ApplyGenesis writes initial state.
func (s *Store) ApplyGenesis(g genesis.Genesis) (err error) {
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
		return err
	}
	if topEr == nil {
		return errors.New("no ERs in genesis")
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
	var lastBlock ibr.LlrIdxFullBlockRecord
	g.Blocks.ForEach(func(br ibr.LlrIdxFullBlockRecord) bool {
		s.WriteFullBlockRecord(br)
		if br.Idx > lastBlock.Idx {
			lastBlock = br
		}
		return true
	})

	// write EVM items
	if reader := g.FwsSection.GetReader(); reader != nil {
		s.Log.Info("Importing Fantom World State data from genesis")
		err := statedb.ImportWorldState(reader)
		if err != nil {
			return fmt.Errorf("failed to import Fantom World State data from genesis; %v", err)
		}
	} else { // no S5 section in the genesis file

		// write EVM data from genesis into leveldb
		err = s.evm.ApplyGenesis(g)
		if err != nil {
			return err
		}

		// write EVM items into Carmen
		if statedb.IsExternalStateDbUsed() {
			s.Log.Info("Importing legacy EVM data into Carmen", "index", lastBlock.Idx, "root", lastBlock.Root)
			if err = statedb.InitializeStateDB(); err != nil {
				return fmt.Errorf("failed to initialize StateDB for the genesis import; %v", err)
			}
			err = statedb.ImportTrieIntoExternalStateDb(s.evm.EvmDb, s.evm.EVMDB(), uint64(lastBlock.Idx), common.Hash(lastBlock.Root))
			if err != nil {
				return fmt.Errorf("genesis import into StateDB failed at block %d; %v", lastBlock.Idx, err)
			}
		}
	}

	// check EVM state hash
	if statedb.IsExternalStateDbUsed() {
		if err = statedb.InitializeStateDB(); err != nil { // make sure the StateDB is available
			return fmt.Errorf("failed to initialize StateDB after the genesis import; %v", err)
		}
		stateHash := statedb.GetExternalStateDbHash()
		if common.Hash(lastBlock.Root) == stateHash {
			s.Log.Info("Imported block into StateDB", "index", lastBlock.Idx, "root", lastBlock.Root)
		} else {
			s.Log.Warn("Imported block into StateDB with not-matching state hash", "index", lastBlock.Idx, "expectedHash", lastBlock.Root, "reproducedHash", stateHash)
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

	return nil
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
