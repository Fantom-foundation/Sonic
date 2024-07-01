package vecclock

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/vecclock/highestbefore"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/ethereum/go-ethereum/log"
)

type kv struct {
	a, b hash.Event
}

// ForklessCause calculates "sufficient coherence" between the events.
// The A.HighestBefore array remembers the sequence number of the last
// event by each validator that is an ancestor of A. The array for
// B.LowestAfter remembers the sequence number of the earliest
// event by each validator that is a descendant of B. Compare the two arrays,
// and find how many elements in the A.HighestBefore array are greater
// than or equal to the corresponding element of the B.LowestAfter
// array. If there are more than 2n/3 such matches, then the A and B
// have achieved sufficient coherency.
//
// If B1 and B2 are forks, then they cannot BOTH forkless-cause any specific event A,
// unless more than 1/3W are Byzantine.
// This great property is the reason why this function exists,
// providing the base for the BFT algorithm.
func (vi *Index) ForklessCause(aID, bID hash.Event) bool {
	if res, ok := vi.cache.ForklessCause.Get(kv{aID, bID}); ok {
		return res.(bool)
	}

	res := vi.forklessCause(aID, bID)

	vi.cache.ForklessCause.Add(kv{aID, bID}, res, 1)
	return res
}

func (vi *Index) forklessCause(aID, bID hash.Event) bool {
	branchIDs := vi.bi.BranchIDCreatorIdxs
	// Get events by hash
	a := vi.getHB(aID)
	if a == nil {
		log.Crit("Vector clock error", "err", fmt.Errorf("Event A=%s not found", aID.String()))
		return false
	}

	// check A doesn't observe any forks from B
	if vi.AtLeastOneFork() {
		a.Alloc(idx.Validator(len(branchIDs)) - 1)
		bBranchID, _ := vi.branchid.Get(bID)
		if a[bBranchID].IsForkDetected() { // B is observed as cheater by A
			return false
		}
	}

	// check A observes that {QUORUM} non-cheater-validators observe B
	b := vi.getLA(bID)
	if b == nil {
		log.Crit("Vector clock error", "err", fmt.Errorf("Event B=%s not found", bID.String()))
		return false
	}
	if vi.AtLeastOneFork() {
		b.Alloc(idx.Validator(len(branchIDs)) - 1)
	}

	yes := vi.validators.NewCounter()
	// calculate forkless causing using the indexes
	for branchIDint, creatorIdx := range branchIDs {
		branchID := idx.Validator(branchIDint)

		// bLowestAfter := vi.GetLowestAfterSeq_(bID, branchID)   // lowest event from creator on branchID, which observes B
		bLowestAfter := b[branchID]   // lowest event from creator on branchID, which observes B
		aHighestBefore := a[branchID] // highest event from creator, observed by A

		// if lowest event from branchID which observes B <= highest from branchID observed by A
		// then {highest from branchID observed by A} observes B
		if bLowestAfter <= aHighestBefore.Seq && bLowestAfter != 0 && !aHighestBefore.IsForkDetected() {
			// we may count the same creator multiple times (on different branches)!
			// so not every call increases the counter
			yes.CountByIdx(creatorIdx)
		}
	}
	return yes.HasQuorum()
}

func (vi *Index) ForklessCauseProgress(aID, bID hash.Event, candidateParents, chosenParents hash.Events) (*pos.WeightCounter, []*pos.WeightCounter) {
	// This function is used to determine progress of event bID in forkless causing aID.
	// It may be used to determine progress toward the forkless cause condition for an event not in vi, but whose parents are in vi.
	// To do so, aID should be the self-parent while chosenParents should be the parents of the not-yet-created event.
	// Further, this function can be used to determine the incremental improvement in progress toward satisfying the forkless
	// cause condition beyond the progress of aId and chosenParents, obtained by inclusion of one additional candidate head at a time.
	// This function is useful in parent selection and event creation timing.

	// The first return is ForklessCause(a + chosenParents, b).
	// The second return argument is a slice containing ForklessCause(a + chosenParents + candidateParent, b) with each element in the
	// slice corresponding to each candidate parent in candidateParents.

	// create the counters that measure the forkless cause progress
	branchIDs := vi.bi.BranchIDCreatorIdxs
	candidateParentsFCProgress := make([]*pos.WeightCounter, len(candidateParents))
	for i, _ := range candidateParentsFCProgress {
		candidateParentsFCProgress[i] = vi.validators.NewCounter() // initialise the counter for each candidate parent
	}
	chosenParentsFCProgress := vi.validators.NewCounter() // initialise the counter for chosen parents only

	// Get events by hash
	aHB := vi.getHB(aID)
	if aHB == nil {
		log.Crit("Vector clock error", "err", fmt.Errorf("Event A=%s not found", aID.String()))
		return chosenParentsFCProgress, candidateParentsFCProgress
	}

	candidateParentsHB := make([]highestbefore.Types, len(candidateParents))
	for i, _ := range candidateParents {
		candidateParentsHB[i] = vi.getHB(candidateParents[i])
		if candidateParentsHB[i] == nil {
			log.Crit("Vector clock error", "err", fmt.Errorf("Candidate parent=%s not found", candidateParents[i].String()))
			return chosenParentsFCProgress, candidateParentsFCProgress
		}
	}

	chosenParentsHB := make([]highestbefore.Types, len(chosenParents))
	for i, _ := range chosenParents {
		chosenParentsHB[i] = vi.getHB(chosenParents[i])
		if chosenParentsHB[i] == nil {
			log.Crit("Vector clock error", "err", fmt.Errorf("Chosen parent=%s not found", chosenParents[i].String()))
			return chosenParentsFCProgress, candidateParentsFCProgress
		}
	}

	// check A doesn't observe any forks from B
	if vi.AtLeastOneFork() {
		bBranchID, _ := vi.branchid.Get(bID)
		if aHB.Get(bBranchID).IsForkDetected() { // B is observed as cheater by A
			return chosenParentsFCProgress, candidateParentsFCProgress
		}
	}

	// check chosenParents don't observe any forks from B
	for i := 0; i < len(chosenParentsHB); i++ {
		if vi.AtLeastOneFork() {
			bBranchID, _ := vi.branchid.Get(bID)
			if chosenParentsHB[i].Get(bBranchID).IsForkDetected() { // B is observed as cheater by a chosen parent
				return chosenParentsFCProgress, candidateParentsFCProgress
			}
		}
	}

	// check candidateParents don't observe any forks from B
	for i := 0; i < len(candidateParentsHB); i++ {
		if vi.AtLeastOneFork() {
			bBranchID, _ := vi.branchid.Get(bID)
			if candidateParentsHB[i].Get(bBranchID).IsForkDetected() { // B is observed as cheater by a candidate parent
				return chosenParentsFCProgress, candidateParentsFCProgress
			}
		}
	}

	bLA := vi.getLA(bID)
	if bLA == nil {
		log.Crit("Vector clock error", "err", fmt.Errorf("Event B=%s not found", bID.String()))
		return chosenParentsFCProgress, candidateParentsFCProgress
	}

	// calculate forkless causing using the indexes
	for branchIDint, creatorIdx := range branchIDs {
		branchID := idx.Validator(branchIDint)

		// bLowestAfter := vi.GetLowestAfterSeq_(bID, branchID)   // lowest event from creator on branchID, which observes B
		bLowestAfter := bLA.Get(branchID)  // lowest event from creator on branchID, which observes B
		HighestBefore := aHB.Get(branchID) // highest event from creator, observed by A

		IsForkDetected := HighestBefore.IsForkDetected()

		for i, _ := range chosenParents {
			chosenParentHighestBefore := chosenParentsHB[i].Get(branchID)                  // highest event from creator, observed by a chosen parent
			HighestBefore.Seq = maxEvent(HighestBefore.Seq, chosenParentHighestBefore.Seq) // find HighestBefore as observed by a and all chosen parents
			IsForkDetected = IsForkDetected || chosenParentHighestBefore.IsForkDetected()
		}

		// first do forkless cause for a + chosenParents only
		if bLowestAfter <= HighestBefore.Seq && bLowestAfter != 0 && !IsForkDetected {
			// we may count the same creator multiple times (on different branches)!
			// so not every call increases the counter
			chosenParentsFCProgress.CountByIdx(creatorIdx)
		}
		// now do forkless cause for a + chosenParents + each candidate parent
		for i, _ := range candidateParents {
			candidateParentHighestBefore := candidateParentsHB[i].Get(branchID)
			candidateParentIsForkDetected := IsForkDetected || candidateParentHighestBefore.IsForkDetected()
			candidateParentHighestBefore.Seq = maxEvent(HighestBefore.Seq, candidateParentHighestBefore.Seq)

			if bLowestAfter <= candidateParentHighestBefore.Seq && bLowestAfter != 0 && !candidateParentIsForkDetected {
				// we may count the same creator multiple times (on different branches)!
				// so not every call increases the counter
				candidateParentsFCProgress[i].CountByIdx(creatorIdx)
			}
		}
	}
	// We want FC progress for new candidate events with parents aID + chosenParents + head
	// aID may not contribute to forkless cause without the heads,
	// but may contribute with the heads. HighestBefore and LowestAfter used above do not incorporate
	// these potential new events, so ensure the contribution of aID's creator is checked and made here
	aCreatorID := vi.getEvent(aID).Creator()
	for _, FC := range candidateParentsFCProgress {
		if FC.Sum() > 0 { // if anything in candidate event's subgraph observes bID, then the candidate must too
			FC.Count(aCreatorID)
		}
	}
	return chosenParentsFCProgress, candidateParentsFCProgress
}

func maxEvent(a idx.Event, b idx.Event) idx.Event {
	if a > b {
		return a
	}
	return b
}
