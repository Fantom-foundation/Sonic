package concurrent

import (
	"sync"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
)

type ValidatorEventsSet struct {
	sync.RWMutex
	Val map[idx.ValidatorID]hash.EventHash
}

func WrapValidatorEventsSet(v map[idx.ValidatorID]hash.EventHash) *ValidatorEventsSet {
	return &ValidatorEventsSet{
		RWMutex: sync.RWMutex{},
		Val:     v,
	}
}
