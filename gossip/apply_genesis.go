package gossip

import (
	"errors"
	"fmt"
	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
	"github.com/Fantom-foundation/go-opera/inter/ibr"
	"github.com/Fantom-foundation/go-opera/inter/ier"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
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
	if !s.StateDbManager.IsAlreadyImported() {
		liveReader, err := g.FwsSection.GetReader()
		if err != nil {
			s.Log.Info("Fantom World State data not available in the genesis", "err", err)
		}

		if liveReader != nil { // has S5 section - import S5 data
			s.Log.Info("Importing Fantom World State data from genesis", "index", lastBlock.Idx)
			archiveReader, err := g.FwsSection.GetReader() // second reader of the same section for the archive import
			if err != nil {
				return fmt.Errorf("failed to get second FWS section reader; %v", err)
			}
			err = s.StateDbManager.ImportWorldState(liveReader, archiveReader, uint64(lastBlock.Idx))
			if err != nil {
				return fmt.Errorf("failed to import Fantom World State data from genesis; %v", err)
			}
		} else { // no S5 section in the genesis file
			// Import legacy EVM genesis section
			err = s.StateDbManager.ImportLegacyEvmData(g.RawEvmItems, uint64(lastBlock.Idx), common.Hash(lastBlock.Root))
			if err != nil {
				return fmt.Errorf("import of legacy genesis data into StateDB failed; %v", err)
			}
		}
	} else {
		s.Log.Info("EVM data import skipped - data already present")
	}

	if err := s.StateDbManager.CheckLiveStateHash(uint64(lastBlock.Idx), common.Hash(lastBlock.Root)); err != nil {
		return err
	} else {
		s.Log.Info("StateDB imported successfully, stateRoot matches", "index", lastBlock.Idx, "root", lastBlock.Root)
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
