package evmstore

import (
	"fmt"
	"path/filepath"

	"github.com/Fantom-foundation/Carmen/go/carmen"
	"github.com/Fantom-foundation/Carmen/go/database/mpt"
	"github.com/Fantom-foundation/Carmen/go/database/mpt/io"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func (s *Store) VerifyWorldState(expectedBlockNum uint64, expectedHash common.Hash) error {
	if s.carmenDb != nil {
		return fmt.Errorf("carmen state must be closed for the world state verification")
	}

	observer := verificationObserver{s.Log}

	// check hash of the live state / last state in the archive
	if err := verifyLastState(s.cfg.Directory, s.cfg.Archive, expectedBlockNum, expectedHash); err != nil {
		return fmt.Errorf("verification of the last block failed: %w", err)
	}
	s.Log.Info("State hash matches the last block state root.")

	// verify the live world state
	liveDir := filepath.Join(s.cfg.Directory, "live")
	info, err := io.CheckMptDirectoryAndGetInfo(liveDir)
	if err != nil {
		return fmt.Errorf("failed to check live state dir: %w", err)
	}
	if err := mpt.VerifyFileLiveTrie(liveDir, info.Config, observer); err != nil {
		return fmt.Errorf("live state verification failed: %w", err)
	}
	s.Log.Info("Live state verified successfully.")

	// verify the archive
	if !s.cfg.Archive {
		return nil // skip archive checks when S5 archive is not used
	}
	archiveDir := filepath.Join(s.cfg.Directory, "archive")
	archiveInfo, err := io.CheckMptDirectoryAndGetInfo(archiveDir)
	if err != nil {
		return fmt.Errorf("failed to check archive dir: %w", err)
	}
	if err := mpt.VerifyArchive(archiveDir, archiveInfo.Config, observer); err != nil {
		return fmt.Errorf("archive verification failed: %w", err)
	}
	s.Log.Info("Archive verified successfully.")
	return nil
}

func verifyLastState(directory string, archive bool, expectedBlockNum uint64, expectedHash common.Hash) error {
	config := carmen.GetCarmenGoS5WithoutArchiveConfiguration()
	if archive {
		config = carmen.GetCarmenGoS5WithArchiveConfiguration()
	}
	db, err := carmen.OpenDatabase(directory, config, carmen.Properties{})
	if err != nil {
		return fmt.Errorf("failed to open carmen live state in %s: %w", directory, err)
	}
	defer db.Close()

	var stateHash common.Hash
	err = db.QueryHeadState(func(ctxt carmen.QueryContext) {
		stateHash = common.Hash(ctxt.GetStateHash())
	})
	if err != nil {
		return fmt.Errorf("failed to get state hash; %w", err)
	}
	if stateHash != expectedHash {
		return fmt.Errorf("state hash does not match (%x != %x)", stateHash, expectedHash)
	}

	if !archive {
		return nil // skip archive checks when archive is not enabled
	}
	lastArchiveBlock, err := db.GetArchiveBlockHeight()
	if err != nil {
		return fmt.Errorf("failed to get last archive block height; %w", err)
	}
	if uint64(lastArchiveBlock) != expectedBlockNum {
		return fmt.Errorf("the last archive block height does not match (%d != %d)", lastArchiveBlock, expectedBlockNum)
	}

	err = db.QueryHistoricState(expectedBlockNum, func(ctxt carmen.QueryContext) {
		stateHash = common.Hash(ctxt.GetStateHash())
	})
	if err != nil {
		return fmt.Errorf("failed to get state hash; %w", err)
	}
	if stateHash != expectedHash {
		return fmt.Errorf("state hash does not match (%x != %x)", stateHash, expectedHash)
	}

	return nil
}

type verificationObserver struct {
	log.Logger
}

func (o verificationObserver) StartVerification() {}

func (o verificationObserver) Progress(msg string) {
	o.Info(msg)
}

func (o verificationObserver) EndVerification(res error) {}
