package concurrent

import (
	"sync"

	"github.com/Fantom-foundation/lachesis-base/hash"
)

type EventsSet struct {
	sync.RWMutex
	Val hash.EventHashSet
}

func WrapEventsSet(v hash.EventHashSet) *EventsSet {
	return &EventsSet{
		RWMutex: sync.RWMutex{},
		Val:     v,
	}
}
