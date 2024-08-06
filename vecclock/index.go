package vecclock

import (
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/vecclock/branchid"
	"github.com/Fantom-foundation/go-opera/vecclock/highestbefore"
	"github.com/Fantom-foundation/go-opera/vecclock/lowestafter"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/dag"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/utils/simplewlru"
	"github.com/ethereum/go-ethereum/log"
)

type Index struct {
	validators    *pos.Validators
	validatorIdxs map[idx.ValidatorID]idx.Validator

	cfg Config

	bi *BranchesInfo

	la       lowestafter.Map
	hb       highestbefore.Map
	branchid branchid.Map

	lastCommit Commit

	getEvent func(hash.Event) dag.Event
	getTmpDB func(name string) kvdb.Store

	cache struct {
		ForklessCause *simplewlru.Cache
	}
}

// NewIndex creates Index instance.
func NewIndex(getTmpDB func(name string) kvdb.Store, cfg Config) *Index {
	vi := &Index{
		getTmpDB: getTmpDB,
		cfg:      cfg,
	}
	vi.cache.ForklessCause, _ = simplewlru.New(uint(vi.cfg.Caches.ForklessCausePairs), vi.cfg.Caches.ForklessCausePairs)
	vi.hb = highestbefore.NewMap(vi.cfg.Caches.HighestBefore, vi.getTmpDB)
	vi.la = lowestafter.NewMap(vi.cfg.Caches.LowestAfter, vi.getTmpDB)
	vi.branchid = branchid.NewMap(vi.cfg.Caches.LowestAfter, vi.getTmpDB)

	return vi
}

// Reset the state
func (vi *Index) Reset(validators *pos.Validators, getEvent func(hash.Event) dag.Event) {
	vi.getEvent = getEvent
	vi.validators = validators
	vi.validatorIdxs = validators.Idxs()
	vi.bi = newInitialBranchesInfo(vi.validators)
	_ = vi.hb.Close()
	vi.hb = highestbefore.NewMap(vi.cfg.Caches.HighestBefore, vi.getTmpDB)
	_ = vi.la.Close()
	vi.la = lowestafter.NewMap(vi.cfg.Caches.LowestAfter, vi.getTmpDB)
	_ = vi.branchid.Close()
	vi.branchid = branchid.NewMap(vi.cfg.Caches.BranchID, vi.getTmpDB)
	vi.cache.ForklessCause.Purge()
}

func (vi *Index) Close() {
	_ = vi.hb.Close()
	_ = vi.la.Close()
	_ = vi.branchid.Close()
	vi.cache.ForklessCause = nil
}

// Add calculates vector clocks for the event
func (vi *Index) Add(e dag.Event) {
	if vi.lastCommit.index != nil {
		log.Crit("Last VecClock operation wasn't commited or reverted")
	}
	c := newCommit(vi)
	getHB := func(id hash.Event) highestbefore.Types {
		v, _ := vi.hb.Get(id)
		return v
	}
	setHB := func(id hash.Event, v highestbefore.Types) {
		c.diffHB[id] = v
	}
	getLA := func(id hash.Event) lowestafter.Types {
		if v := c.diffLA[id]; v != nil {
			return v
		}
		v, _ := vi.la.Get(id)
		return v
	}
	setLA := func(id hash.Event, v lowestafter.Types) {
		c.diffLA[id] = v
	}
	branchID := vi.fillEventVectors(e, getHB, setHB, getLA, setLA)
	c.diffEvent = e.ID()
	c.diffBranchID = branchID
	vi.lastCommit = c
}

func (vi *Index) setForkDetected(before highestbefore.Types, branchID idx.Validator) {
	creatorIdx := vi.bi.BranchIDCreatorIdxs[branchID]
	for _, branchID := range vi.bi.BranchIDByCreators[creatorIdx] {
		for i := len(before); i <= int(branchID); i++ {
			before = append(before, highestbefore.Type{})
		}
		before[branchID] = highestbefore.ForkDetectedSeq
	}
}

func (vi *Index) fillGlobalBranchID(e dag.Event, meIdx idx.Validator) idx.Validator {
	// sanity checks
	if len(vi.bi.BranchIDCreatorIdxs) != len(vi.bi.BranchIDLastSeq) {
		log.Crit("inconsistent BranchIDCreators len (inconsistent DB)")
	}
	if idx.Validator(len(vi.bi.BranchIDCreatorIdxs)) < vi.validators.Len() {
		log.Crit("inconsistent BranchIDCreators len (inconsistent DB)")
	}

	if e.SelfParent() == nil {
		// is it first event indeed?
		if vi.bi.BranchIDLastSeq[meIdx] == 0 {
			// OK, not a new fork
			vi.bi.BranchIDLastSeq[meIdx] = e.Seq()
			return meIdx
		}
	} else {
		selfParentBranchID, _ := vi.branchid.Get(*e.SelfParent())
		// sanity checks
		if len(vi.bi.BranchIDCreatorIdxs) != len(vi.bi.BranchIDLastSeq) {
			log.Crit("inconsistent BranchIDCreators len (inconsistent DB)")
		}

		if vi.bi.BranchIDLastSeq[selfParentBranchID]+1 == e.Seq() {
			vi.bi.BranchIDLastSeq[selfParentBranchID] = e.Seq()
			// OK, not a new fork
			return selfParentBranchID
		}
	}

	// if we're here, then new fork is observed (only globally), create new branchID due to a new fork
	vi.bi.BranchIDLastSeq = append(vi.bi.BranchIDLastSeq, e.Seq())
	vi.bi.BranchIDCreatorIdxs = append(vi.bi.BranchIDCreatorIdxs, meIdx)
	newBranchID := idx.Validator(len(vi.bi.BranchIDLastSeq) - 1)
	vi.bi.BranchIDByCreators[meIdx] = append(vi.bi.BranchIDByCreators[meIdx], newBranchID)
	return newBranchID
}

type CreationTimer interface {
	CreationTime() inter.Timestamp
}

func collectFrom(self highestbefore.Types, other highestbefore.Types, num idx.Validator) {
	if num > idx.Validator(len(other)) {
		num = idx.Validator(len(other))
	}
	for branchID := idx.Validator(0); branchID < num; branchID++ {
		hisSeq := other[branchID]
		if hisSeq.Seq == 0 && !hisSeq.IsForkDetected() {
			// hisSeq doesn't observe anything about this branchID
			continue
		}
		mySeq := self[branchID]

		if mySeq.IsForkDetected() {
			// mySeq observes the maximum already
			continue
		}
		if hisSeq.IsForkDetected() {
			// set fork detected
			self[branchID] = highestbefore.ForkDetectedSeq
		} else {
			if mySeq.Seq == 0 || mySeq.MinSeq > hisSeq.MinSeq {
				// take hisSeq.MinSeq
				mySeq.MinSeq = hisSeq.MinSeq
				self[branchID] = mySeq
			}
			if mySeq.Seq < hisSeq.Seq {
				// take hisSeq.Seq
				mySeq.Seq = hisSeq.Seq
				mySeq.Seq = hisSeq.Seq
				mySeq.Time = hisSeq.Time
				self[branchID] = mySeq
			}
		}
	}
}

// fillEventVectors calculates (and stores) event's vectors, and updates LowestAfter of newly-observed events.
func (vi *Index) fillEventVectors(e dag.Event,
	getHB func(hash.Event) highestbefore.Types, setHB func(hash.Event, highestbefore.Types),
	getLA func(hash.Event) lowestafter.Types, setLA func(hash.Event, lowestafter.Types)) idx.Validator {
	meIdx := vi.validatorIdxs[e.Creator()]

	meBranchID := vi.fillGlobalBranchID(e, meIdx)
	myVecs := allVecs{
		before: make(highestbefore.Types, idx.Validator(len(vi.bi.BranchIDCreatorIdxs))),
		after:  make([]idx.Event, idx.Validator(len(vi.bi.BranchIDCreatorIdxs))),
	}

	// pre-load parents into RAM for quick access
	parentsVecs := make([]highestbefore.Types, len(e.Parents()))
	parentsBranchIDs := make([]idx.Validator, len(e.Parents()))
	for i, p := range e.Parents() {
		parentsBranchIDs[i], _ = vi.branchid.Get(p)
		parentsVecs[i] = getHB(p)
		if parentsVecs[i] == nil {
			log.Crit("processed out of order, parent not found (inconsistent DB)", "parent", p.String())
		}
	}

	// observed by himself
	myVecs.after[meBranchID] = e.Seq()
	myVecs.before[meBranchID].Seq = e.Seq()
	myVecs.before[meBranchID].MinSeq = e.Seq()
	if et, ok := e.(CreationTimer); ok {
		myVecs.before[meBranchID].Time = et.CreationTime()
	}

	for _, pVec := range parentsVecs {
		// calculate HighestBefore, Detect forks for a case when parent observes a fork
		collectFrom(myVecs.before, pVec, idx.Validator(len(vi.bi.BranchIDCreatorIdxs)))
	}
	// Detect forks, which were not observed by parents
	if vi.AtLeastOneFork() {
		for n := idx.Validator(0); n < vi.validators.Len(); n++ {
			if len(vi.bi.BranchIDByCreators[n]) <= 1 {
				continue
			}
			for _, branchID := range vi.bi.BranchIDByCreators[n] {
				if myVecs.before[branchID].IsForkDetected() {
					// if one branch observes a fork, mark all the branches as observing the fork
					vi.setForkDetected(myVecs.before, n)
					break
				}
			}
		}

	nextCreator:
		for n := idx.Validator(0); n < vi.validators.Len(); n++ {
			if myVecs.before[n].IsForkDetected() {
				continue
			}
			for _, branchID1 := range vi.bi.BranchIDByCreators[n] {
				for _, branchID2 := range vi.bi.BranchIDByCreators[n] {
					a := branchID1
					b := branchID2
					if a == b {
						continue
					}

					if myVecs.before[a].IsEmpty() || myVecs.before[b].IsEmpty() {
						continue
					}
					if myVecs.before[a].MinSeq <= myVecs.before[b].Seq && myVecs.before[b].MinSeq <= myVecs.before[a].Seq {
						vi.setForkDetected(myVecs.before, n)
						continue nextCreator
					}
				}
			}
		}
	}

	// graph traversal starting from e, but excluding e
	onWalk := func(walk hash.Event) (godeeper bool) {
		wLowestAfterSeq := getLA(walk)

		if wLowestAfterSeq.Get(meBranchID) != 0 {
			return false
		}

		wLowestAfterSeq.Alloc(meBranchID)
		wLowestAfterSeq[meBranchID] = e.Seq()
		setLA(walk, wLowestAfterSeq)
		return true
	}
	err := vi.DfsSubgraph(e, onWalk)
	if err != nil {
		log.Crit("DFS error", "err", err)
	}

	// store calculated vectors
	setHB(e.ID(), myVecs.before)
	setLA(e.ID(), myVecs.after)

	return meBranchID
}

func gatherFrom(self highestbefore.Types, to idx.Validator, other highestbefore.Types, from []idx.Validator) {
	// read all branches to find highest event
	highestBranch := highestbefore.Type{}
	for _, branchID := range from {
		vother := other.Get(branchID)
		if vother.IsForkDetected() {
			highestBranch = vother
			break
		}
		if vother.Seq > highestBranch.Seq {
			highestBranch.Seq = vother.Seq
			highestBranch.Time = vother.Time
		}
	}
	self.Alloc(to)
	self[to] = highestBranch
}

func (vi *Index) GetMergedHighestBefore(id hash.Event) highestbefore.Types {
	if vi.AtLeastOneFork() {
		scatteredBefore := vi.getHB(id)
		if scatteredBefore == nil {
			return nil
		}

		mergedBefore := make(highestbefore.Types, vi.validators.Len())

		for creatorIdx, branches := range vi.bi.BranchIDByCreators {
			gatherFrom(mergedBefore, idx.Validator(creatorIdx), scatteredBefore, branches)
		}

		return mergedBefore
	}
	v := vi.getHB(id)
	return v
}
