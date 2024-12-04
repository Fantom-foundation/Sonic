package concurrent

import (
	"sync"

	"github.com/Fantom-foundation/lachesis-base/ltypes"
)

type ValidatorEventsSet struct {
	sync.RWMutex
	Val map[ltypes.ValidatorID]ltypes.EventHash
}

func WrapValidatorEventsSet(v map[ltypes.ValidatorID]ltypes.EventHash) *ValidatorEventsSet {
	return &ValidatorEventsSet{
		RWMutex: sync.RWMutex{},
		Val:     v,
	}
}
