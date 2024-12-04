package emitter

import (
	"time"

	"github.com/Fantom-foundation/lachesis-base/emitter/ancestor"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
)

// buildSearchStrategies returns a strategy for each parent search
func (em *Emitter) buildSearchStrategies(maxParents idx.EventID) []ancestor.SearchStrategy {
	strategies := make([]ancestor.SearchStrategy, 0, maxParents)
	if maxParents == 0 {
		return strategies
	}
	payloadStrategy := em.payloadIndexer.SearchStrategy()
	for idx.EventID(len(strategies)) < 1 {
		strategies = append(strategies, payloadStrategy)
	}
	randStrategy := ancestor.NewRandomStrategy(nil)
	for idx.EventID(len(strategies)) < maxParents/2 {
		strategies = append(strategies, randStrategy)
	}
	if em.fcIndexer != nil {
		quorumStrategy := em.fcIndexer.SearchStrategy()
		for idx.EventID(len(strategies)) < maxParents {
			strategies = append(strategies, quorumStrategy)
		}
	} else if em.quorumIndexer != nil {
		quorumStrategy := em.quorumIndexer.SearchStrategy()
		for idx.EventID(len(strategies)) < maxParents {
			strategies = append(strategies, quorumStrategy)
		}
	}
	return strategies
}

// chooseParents selects an "optimal" parents set for the validator
func (em *Emitter) chooseParents(epoch idx.EpochID, myValidatorID idx.ValidatorID) (*hash.EventHash, hash.EventHashes, bool) {
	selfParent := em.world.GetLastEvent(epoch, myValidatorID)
	if selfParent == nil {
		return nil, nil, true
	}
	if len(em.world.DagIndex().NoCheaters(selfParent, hash.EventHashes{*selfParent})) == 0 {
		em.Periodic.Error(time.Second, "Events emitting isn't allowed due to the doublesign", "validator", myValidatorID)
		return nil, nil, false
	}
	parents := hash.EventHashes{*selfParent}
	heads := em.world.GetHeads(epoch) // events with no descendants
	parents = ancestor.ChooseParents(parents, heads, em.buildSearchStrategies(em.maxParents-idx.EventID(len(parents))))
	return selfParent, parents, true
}
