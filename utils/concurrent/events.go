package concurrent

import (
	"sync"

	"github.com/Fantom-foundation/lachesis-base/ltypes"
)

type EventsSet struct {
	sync.RWMutex
	Val ltypes.EventHashSet
}

func WrapEventsSet(v ltypes.EventHashSet) *EventsSet {
	return &EventsSet{
		RWMutex: sync.RWMutex{},
		Val:     v,
	}
}
