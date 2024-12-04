package emitter

import (
	"time"

	"github.com/Fantom-foundation/lachesis-base/emitter/ancestor"
	"github.com/Fantom-foundation/lachesis-base/ltypes"
)

// buildSearchStrategies returns a strategy for each parent search
func (em *Emitter) buildSearchStrategies(maxParents ltypes.EventID) []ancestor.SearchStrategy {
	strategies := make([]ancestor.SearchStrategy, 0, maxParents)
	if maxParents == 0 {
		return strategies
	}
	payloadStrategy := em.payloadIndexer.SearchStrategy()
	for ltypes.EventID(len(strategies)) < 1 {
		strategies = append(strategies, payloadStrategy)
	}
	randStrategy := ancestor.NewRandomStrategy(nil)
	for ltypes.EventID(len(strategies)) < maxParents/2 {
		strategies = append(strategies, randStrategy)
	}
	if em.fcIndexer != nil {
		quorumStrategy := em.fcIndexer.SearchStrategy()
		for ltypes.EventID(len(strategies)) < maxParents {
			strategies = append(strategies, quorumStrategy)
		}
	} else if em.quorumIndexer != nil {
		quorumStrategy := em.quorumIndexer.SearchStrategy()
		for ltypes.EventID(len(strategies)) < maxParents {
			strategies = append(strategies, quorumStrategy)
		}
	}
	return strategies
}

// chooseParents selects an "optimal" parents set for the validator
func (em *Emitter) chooseParents(epoch ltypes.EpochID, myValidatorID ltypes.ValidatorID) (*ltypes.EventHash, ltypes.EventHashes, bool) {
	selfParent := em.world.GetLastEvent(epoch, myValidatorID)
	if selfParent == nil {
		return nil, nil, true
	}
	if len(em.world.DagIndex().NoCheaters(selfParent, ltypes.EventHashes{*selfParent})) == 0 {
		em.Periodic.Error(time.Second, "Events emitting isn't allowed due to the doublesign", "validator", myValidatorID)
		return nil, nil, false
	}
	parents := ltypes.EventHashes{*selfParent}
	heads := em.world.GetHeads(epoch) // events with no descendants
	parents = ancestor.ChooseParents(parents, heads, em.buildSearchStrategies(em.maxParents-ltypes.EventID(len(parents))))
	return selfParent, parents, true
}
