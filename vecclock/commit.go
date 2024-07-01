package vecclock

import (
	"github.com/Fantom-foundation/go-opera/vecclock/highestbefore"
	"github.com/Fantom-foundation/go-opera/vecclock/lowestafter"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
)

type Commit struct {
	index  *Index
	diffLA map[hash.Event]lowestafter.Types
	diffHB map[hash.Event]highestbefore.Types
	//diffHBKeys   []hash.Event
	//diffHBVals   []highestbefore.Types
	biCopy       *BranchesInfo
	diffEvent    hash.Event
	diffBranchID idx.Validator
}

func newCommit(index *Index) Commit {
	return Commit{
		index:  index,
		diffLA: make(map[hash.Event]lowestafter.Types),
		diffHB: make(map[hash.Event]highestbefore.Types),
		//diffHBKeys: make([]hash.Event, 0, 200),
		//diffHBVals: make([]highestbefore.Types, 0, 200),
		biCopy: index.bi.Copy(),
	}
}

func (c *Commit) Commit() {
	if c.index == nil {
		return
	}
	c.index.la.AddMap(c.diffLA)
	c.index.hb.AddMap(c.diffHB)
	//c.index.hb.AddSlice(c.diffHBKeys, c.diffHBVals)
	c.index.branchid.Set(c.diffEvent, c.diffBranchID)
	*c = Commit{}
}

func (c *Commit) Revert() {
	if c.index == nil {
		return
	}
	c.index.bi = c.biCopy
	*c = Commit{}
}

func (vi *Index) Commit() {
	vi.lastCommit.Commit()
}

func (vi *Index) Revert() {
	vi.lastCommit.Revert()
}

func (vi *Index) getHB(id hash.Event) highestbefore.Types {
	if vi.lastCommit.diffHB != nil {
		if v := vi.lastCommit.diffHB[id]; v != nil {
			return v
		}
	}
	v, _ := vi.hb.Get(id)
	return v
}

func (vi *Index) getLA(id hash.Event) lowestafter.Types {
	if vi.lastCommit.diffLA != nil {
		if v := vi.lastCommit.diffLA[id]; v != nil {
			return v
		}
	}
	v, _ := vi.la.Get(id)
	return v
}
