package iblockproc

import (
	"github.com/Fantom-foundation/lachesis-base/ltypes"

	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/opera"
)

type ValidatorEpochStateV0 struct {
	GasRefund      uint64
	PrevEpochEvent ltypes.EventHash
}

type EpochStateV0 struct {
	Epoch          ltypes.EpochID
	EpochStart     inter.Timestamp
	PrevEpochStart inter.Timestamp

	EpochStateRoot ltypes.Hash

	Validators        *ltypes.Validators
	ValidatorStates   []ValidatorEpochStateV0
	ValidatorProfiles ValidatorProfiles

	Rules opera.Rules
}
