package parentlesscheck

import (
	"github.com/Fantom-foundation/lachesis-base/ltypes"
)

type Checker struct {
	HeavyCheck HeavyCheck
	LightCheck LightCheck
}

type LightCheck func(ltypes.Event) error

type HeavyCheck interface {
	Enqueue(e ltypes.Event, checked func(error)) error
}

// Enqueue tries to fill gaps the fetcher's future import queue.
func (c *Checker) Enqueue(e ltypes.Event, checked func(error)) {
	// Run light checks right away
	err := c.LightCheck(e)
	if err != nil {
		checked(err)
		return
	}

	// Run heavy check in parallel
	_ = c.HeavyCheck.Enqueue(e, checked)
}
