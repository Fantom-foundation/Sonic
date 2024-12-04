package heavycheck

import (
	"github.com/Fantom-foundation/lachesis-base/ltypes"

	"github.com/Fantom-foundation/go-opera/inter"
)

type EventsOnly struct {
	*Checker
}

func (c *EventsOnly) Enqueue(e ltypes.Event, onValidated func(error)) error {
	return c.Checker.EnqueueEvent(e.(inter.EventPayloadI), onValidated)
}
