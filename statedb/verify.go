package statedb

import (
	"fmt"
	cc "github.com/Fantom-foundation/Carmen/go/common"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/Carmen/go/state/mpt"
	io2 "github.com/Fantom-foundation/Carmen/go/state/mpt/io"
	"github.com/ethereum/go-ethereum/common"
)

func (m *StateDbManager) IsWorldStateVerifiable() bool {
	return m.doesUseCarmen() && m.parameters.Schema == carmen.StateSchema(5)
}

func (m *StateDbManager) VerifyWorldState(expectedHash common.Hash, observer mpt.VerificationObserver) error {
	if !m.IsWorldStateVerifiable() {
		return fmt.Errorf("unable to verify world state data - Carmen S5 not used")
	}
	// try to obtain information of the contained MPT
	info, err := io2.CheckMptDirectoryAndGetInfo(m.parameters.Directory)
	if err != nil {
		return err
	}
	// get hash of the live state
	liveState, err := carmen.NewState(m.parameters)
	if err != nil {
		return fmt.Errorf("failed to create carmen state; %s", err)
	}
	defer liveState.Close()
	stateHash, err := liveState.GetHash()
	if err != nil {
		return fmt.Errorf("failed to get state hash; %s", err)
	}
	if stateHash != cc.Hash(expectedHash) {
		return fmt.Errorf("validation failed - the live state hash does not match with the last block state root (%x != %x)", stateHash, expectedHash)
	}
	// verify the world state
	return mpt.VerifyFileLiveTrie(m.parameters.Directory, info.Config, observer)
}
