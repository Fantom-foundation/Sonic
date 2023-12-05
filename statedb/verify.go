package statedb

import (
	"fmt"
	cc "github.com/Fantom-foundation/Carmen/go/common"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/Carmen/go/state/mpt"
	"github.com/Fantom-foundation/Carmen/go/state/mpt/io"
	"github.com/Fantom-foundation/go-opera/logger"
	"github.com/ethereum/go-ethereum/common"
	"path/filepath"
)

func (m *StateDbManager) IsWorldStateVerifiable() bool {
	return m.doesUseCarmen() && m.parameters.Schema == carmen.StateSchema(5)
}

func (m *StateDbManager) VerifyWorldState(expectedBlockNum uint64, expectedHash common.Hash) error {
	if m.carmenState != nil {
		return fmt.Errorf("carmen state must be closed for the world state verification")
	}
	if !m.IsWorldStateVerifiable() {
		return fmt.Errorf("unable to verify world state data - Carmen S5 not used")
	}

	observer := verificationObserver{m.logger}

	// check hash of the live state / last state in the archive
	if err := verifyLastState(m.parameters, expectedBlockNum, expectedHash); err != nil {
		return fmt.Errorf("verification of the last block failed: %w", err)
	}
	m.logger.Log.Info("State hash matches the last block state root.")

	// verify the live world state
	info, err := io.CheckMptDirectoryAndGetInfo(m.parameters.Directory)
	if err != nil {
		return fmt.Errorf("failed to check live state dir: %w", err)
	}
	if err := mpt.VerifyFileLiveTrie(m.parameters.Directory, info.Config, observer); err != nil {
		return fmt.Errorf("live state verification failed: %w", err)
	}
	m.logger.Log.Info("Live state verified successfully.")

	// verify the archive
	if m.parameters.Archive != carmen.S5Archive {
		return nil // skip archive checks when S5 archive is not used
	}
	archiveDir := m.parameters.Directory + string(filepath.Separator) + "archive"
	archiveInfo, err := io.CheckMptDirectoryAndGetInfo(archiveDir)
	if err != nil {
		return fmt.Errorf("failed to check archive dir: %w", err)
	}
	if err := mpt.VerifyArchive(archiveDir, archiveInfo.Config, observer); err != nil {
		return fmt.Errorf("archive verification failed: %w", err)
	}
	m.logger.Log.Info("Archive verified successfully.")
	return nil
}

func verifyLastState(params carmen.Parameters, expectedBlockNum uint64, expectedHash common.Hash) error {
	liveState, err := carmen.NewState(params)
	if err != nil {
		return fmt.Errorf("failed to open carmen live state in %s: %w", params.Directory, err)
	}
	defer liveState.Close()
	if err := checkStateHash(liveState, expectedHash); err != nil {
		return fmt.Errorf("live state check failed; %w", err)
	}

	lastArchiveBlock, _, err := liveState.GetArchiveBlockHeight()
	if err != nil {
		return fmt.Errorf("failed to get last archive block height; %w", err)
	}
	if lastArchiveBlock != expectedBlockNum {
		return fmt.Errorf("the last archive block height does not match (%d != %d)", lastArchiveBlock, expectedBlockNum)
	}

	if params.Archive == carmen.NoArchive {
		return nil // skip archive checks when archive is not enabled
	}
	archiveState, err := liveState.GetArchiveState(lastArchiveBlock)
	if err != nil {
		return fmt.Errorf("failed to get carmen archive state; %w", err)
	}
	defer archiveState.Close()
	if err := checkStateHash(archiveState, expectedHash); err != nil {
		return fmt.Errorf("archive state check failed; %w", err)
	}
	return nil
}

func checkStateHash(state carmen.State, expectedHash common.Hash) error {
	stateHash, err := state.GetHash()
	if err != nil {
		return fmt.Errorf("failed to get state hash; %w", err)
	}
	if stateHash != cc.Hash(expectedHash) {
		return fmt.Errorf("state hash does not match (%x != %x)", stateHash, expectedHash)
	}
	return nil
}

type verificationObserver struct {
	logger.Instance
}

func (o verificationObserver) StartVerification() {}

func (o verificationObserver) Progress(msg string) {
	o.Log.Info(msg)
}

func (o verificationObserver) EndVerification(res error) {}
