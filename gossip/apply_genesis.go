package gossip

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip/gasprice"
	"github.com/Fantom-foundation/go-opera/inter"
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

	s.SetGenesisID(g.GenesisID)
	s.SetGenesisBlockIndex(topEr.BlockState.LastBlock.Idx)

	// write blocks
	rules := s.GetRules()
	gasLimit := rules.Blocks.MaxBlockGas

	s.SetBlock(0, inter.NewBlockBuilder().
		WithNumber(0).
		WithTime(evmcore.FakeGenesisTime-1). // TODO: extend genesis generator to provide time
		WithGasLimit(gasLimit).
		WithStateRoot(common.Hash{}). // TODO: get proper state root from genesis data
		WithBaseFee(big.NewInt(0)).   // TODO: set initial base fee according to the rules
		Build(),
	)

	baseFee := prevEs.Rules.Economy.MinGasPrice
	blobGasPrice := big.NewInt(1) // TODO issue #147
	var lastBlock ibr.LlrIdxFullBlockRecord
	g.Blocks.ForEach(func(br ibr.LlrIdxFullBlockRecord) bool {
		var duration time.Duration
		if rules.Upgrades.Sonic {
			block := s.GetBlock(br.Idx - 1)
			if block == nil {
				block = &inter.Block{
					BaseFee: big.NewInt(0),
				}
			}
			header := &evmcore.EvmHeader{
				GasUsed:  block.GasUsed,
				GasLimit: block.GasLimit,
				BaseFee:  block.BaseFee,
			}
			baseFee = gasprice.GetBaseFeeForNextBlock(header, rules.Economy)
			duration = time.Duration(br.Time-block.Time) * time.Nanosecond
		}
		s.WriteFullBlockRecord(baseFee, blobGasPrice, gasLimit, duration, br)
		if br.Idx > lastBlock.Idx {
			lastBlock = br
		}
		return true
	})

	// write EVM items
	liveReader, err := g.FwsLiveSection.GetReader()
	if err != nil {
		s.Log.Info("Fantom World State Live data not available in the genesis", "err", err)
	}

	if liveReader != nil { // has S5 section - import S5 data
		s.Log.Info("Importing Fantom World State Live data from genesis")
		err = s.evm.ImportLiveWorldState(liveReader)
		if err != nil {
			return fmt.Errorf("failed to import Fantom World State data from genesis; %v", err)
		}

		// import S5 archive
		archiveReader, _ := g.FwsArchiveSection.GetReader()
		if archiveReader != nil { // has archive section
			s.Log.Info("Importing Fantom World State Archive data from genesis")
			err = s.evm.ImportArchiveWorldState(archiveReader)
			if err != nil {
				return fmt.Errorf("failed to import Fantom World State Archive data from genesis; %v", err)
			}
		} else { // no archive section - initialize archive from the live section
			s.Log.Info("No archive in the genesis file - initializing the archive from the live state", "blockNum", lastBlock.Idx)
			liveToArchiveReader, err := g.FwsLiveSection.GetReader() // second reader of the same section for the archive import
			if err != nil {
				return fmt.Errorf("failed to get second FWS section reader; %v", err)
			}
			err = s.evm.InitializeArchiveWorldState(liveToArchiveReader, uint64(lastBlock.Idx))
			if err != nil {
				return fmt.Errorf("failed to import Fantom World State data from genesis; %v", err)
			}
		}
	} else { // no S5 section in the genesis file
		// Import legacy EVM genesis section
		err = s.evm.ImportLegacyEvmData(g.RawEvmItems, uint64(lastBlock.Idx), common.Hash(lastBlock.StateRoot))
		if err != nil {
			return fmt.Errorf("import of legacy genesis data into StateDB failed; %v", err)
		}
	}

	if err := s.evm.Open(); err != nil {
		return fmt.Errorf("unable to open EvmStore to check imported state: %w", err)
	}
	if err := s.evm.CheckLiveStateHash(lastBlock.Idx, lastBlock.StateRoot); err != nil {
		return fmt.Errorf("checking imported live state failed: %w", err)
	} else {
		s.Log.Info("StateDB imported successfully, stateRoot matches", "index", lastBlock.Idx, "root", lastBlock.StateRoot)
	}

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
