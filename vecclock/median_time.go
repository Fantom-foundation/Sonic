package vecclock

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"sort"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"

	"github.com/Fantom-foundation/go-opera/inter"
)

// medianTimeIndex is a handy index for the MedianTime() func
type medianTimeIndex struct {
	weight       pos.Weight
	creationTime inter.Timestamp
}

// MedianTime calculates weighted median of claimed time within highest observed events.
func (vi *Index) MedianTime(id hash.Event, defaultTime inter.Timestamp) inter.Timestamp {
	// Get event by hash
	before := vi.GetMergedHighestBefore(id)
	if before == nil {
		log.Crit("MedianTime error", "err", fmt.Errorf("event=%s not found", id.String()))
	}

	honestTotalWeight := pos.Weight(0) // isn't equal to validators.TotalWeight(), because doesn't count cheaters
	highests := make([]medianTimeIndex, 0, len(vi.validatorIdxs))
	// convert []HighestBefore -> []medianTimeIndex
	for creatorIdxI := range vi.validators.IDs() {
		creatorIdx := idx.Validator(creatorIdxI)
		highest := medianTimeIndex{}
		highest.weight = vi.validators.GetWeightByIdx(creatorIdx)
		beforeE := before[creatorIdx]
		highest.creationTime = beforeE.Time

		// edge cases
		if beforeE.IsForkDetected() {
			// cheaters don't influence medianTime
			highest.weight = 0
		} else if beforeE.Seq == 0 {
			// if no event was observed from this node, then use defaultTime
			highest.creationTime = defaultTime
		}

		highests = append(highests, highest)
		honestTotalWeight += highest.weight
	}
	// it's technically possible honestTotalWeight == 0 (all validators are cheaters)

	// sort by claimed time (partial order is enough here, because we need only creationTime)
	sort.Slice(highests, func(i, j int) bool {
		a, b := highests[i], highests[j]
		return a.creationTime < b.creationTime
	})

	// Calculate weighted median
	halfWeight := honestTotalWeight / 2
	var currWeight pos.Weight
	var median inter.Timestamp
	for _, highest := range highests {
		currWeight += highest.weight
		if currWeight >= halfWeight {
			median = highest.creationTime
			break
		}
	}

	// sanity check
	if currWeight < halfWeight || currWeight > honestTotalWeight {
		log.Crit("MedianTime sanity check", "err", fmt.Errorf("median wasn't calculated correctly, median=%d, currWeight=%d, totalWeight=%d, len(highests)=%d, id=%s",
			median,
			currWeight,
			honestTotalWeight,
			len(highests),
			id.String(),
		))
	}

	return median
}
