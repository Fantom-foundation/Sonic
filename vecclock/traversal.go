package vecclock

import (
	"errors"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/dag"
)

// DfsSubgraph iterates all the event which are observed by head, and accepted by a filter
// Excluding head
// filter MAY BE called twice for the same event.
func (vi *Index) DfsSubgraph(head dag.Event, walk func(hash.Event) (godeeper bool)) error {
	stack := make(hash.EventsStack, 0, vi.validators.Len()*5)

	// first element
	stack.PushAll(head.Parents())

	for next := stack.Pop(); next != nil; next = stack.Pop() {
		curr := *next

		// filter
		if !walk(curr) {
			continue
		}

		event := vi.getEvent(curr)
		if event == nil {
			return errors.New("event not found " + curr.String())
		}

		// memorize parents
		stack.PushAll(event.Parents())
	}

	return nil
}
