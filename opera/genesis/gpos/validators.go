package gpos

import (
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-opera/inter/validatorpk"
	"github.com/Fantom-foundation/lachesis-base/ltypes"
)

type (
	// Validator is a helper structure to define genesis validators
	Validator struct {
		ID               ltypes.ValidatorID
		Address          common.Address
		PubKey           validatorpk.PubKey
		CreationTime     inter.Timestamp
		CreationEpoch    ltypes.EpochID
		DeactivatedTime  inter.Timestamp
		DeactivatedEpoch ltypes.EpochID
		Status           uint64
	}

	Validators []Validator
)

// Map converts Validators to map
func (gv Validators) Map() map[ltypes.ValidatorID]Validator {
	validators := map[ltypes.ValidatorID]Validator{}
	for _, val := range gv {
		validators[val.ID] = val
	}
	return validators
}
