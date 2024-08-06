package vecclock

import (
	"github.com/Fantom-foundation/go-opera/vecclock/highestbefore"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/dag"
	"github.com/Fantom-foundation/lachesis-base/inter/dag/tdag"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/Fantom-foundation/lachesis-base/kvdb/memorydb"

	"github.com/Fantom-foundation/go-opera/inter"
)

func testMedianTimeOnIndex(t *testing.T, cfg Config) {
	nodes := tdag.GenNodes(5)
	weights := []pos.Weight{5, 4, 3, 2, 1}
	validators := pos.ArrayToValidators(nodes, weights)

	vi := NewIndex(makeTmpDB, cfg)
	vi.Reset(validators, nil)

	assertar := assert.New(t)
	{ // seq=0
		e := hash.ZeroEvent
		// validator indexes are sorted by weight amount
		before := make(highestbefore.Types, validators.Len())

		before[0].Seq = 0
		before[0].Time = 100

		before[1].Seq = 0
		before[1].Time = 100

		before[2].Seq = 1
		before[2].Time = 10

		before[3].Seq = 1
		before[3].Time = 10

		before[4].Seq = 1
		before[4].Time = 10

		vi.hb.Set(e, before)
		assertar.Equal(inter.Timestamp(1), vi.MedianTime(e, 1))
	}

	{ // fork seen = true
		e := hash.ZeroEvent
		// validator indexes are sorted by weight amount
		before := make(highestbefore.Types, validators.Len())

		before[0] = highestbefore.ForkDetectedSeq
		before[0].Time = 100

		before[1] = highestbefore.ForkDetectedSeq
		before[1].Time = 100

		before[2].Seq = 1
		before[2].Time = 10

		before[3].Seq = 1
		before[3].Time = 10

		before[4].Seq = 1
		before[4].Time = 10

		vi.hb.Set(e, before)
		assertar.Equal(inter.Timestamp(10), vi.MedianTime(e, 1))
	}

	{ // normal
		e := hash.ZeroEvent
		// validator indexes are sorted by weight amount
		before := make(highestbefore.Types, validators.Len())

		before[0].Seq = 1
		before[0].Time = 11

		before[1].Seq = 2
		before[1].Time = 12

		before[2].Seq = 2
		before[2].Time = 13

		before[3].Seq = 3
		before[3].Time = 14

		before[4].Seq = 4
		before[4].Time = 15

		vi.hb.Set(e, before)
		assertar.Equal(inter.Timestamp(12), vi.MedianTime(e, 1))
	}

}

func TestMedianTimeOnIndex(t *testing.T) {
	testMedianTimeOnIndex(t, LiteConfig())
	// test without cache
	testMedianTimeOnIndex(t, Config{})
}

func TestMedianTimeOnDAG(t *testing.T) {
	dagAscii := `
 ║
 nodeA001
 ║
 nodeA012
 ║            ║
 ║            nodeB001
 ║            ║            ║
 ║            ╠═══════════ nodeC001
 ║║           ║            ║            ║
 ║╚══════════─╫─══════════─╫─══════════ nodeD001
║║            ║            ║            ║
╚ nodeA002════╬════════════╬════════════╣
 ║║           ║            ║            ║
 ║╚══════════─╫─══════════─╫─══════════ nodeD002
 ║            ║            ║            ║
 nodeA003════─╫─══════════─╫─═══════════╣
 ║            ║            ║
 ╠════════════nodeB002     ║
 ║            ║            ║
 ╠════════════╫═══════════ nodeC002
`

	weights := []pos.Weight{3, 4, 2, 1}
	genesisTime := inter.Timestamp(1)
	creationTimes := map[string]inter.Timestamp{
		"nodeA001": inter.Timestamp(111),
		"nodeB001": inter.Timestamp(112),
		"nodeC001": inter.Timestamp(13),
		"nodeD001": inter.Timestamp(14),
		"nodeA002": inter.Timestamp(120),
		"nodeD002": inter.Timestamp(20),
		"nodeA012": inter.Timestamp(120),
		"nodeA003": inter.Timestamp(20),
		"nodeB002": inter.Timestamp(20),
		"nodeC002": inter.Timestamp(35),
	}
	medianTimes := map[string]inter.Timestamp{
		"nodeA001": genesisTime,
		"nodeB001": genesisTime,
		"nodeC001": inter.Timestamp(13),
		"nodeD001": genesisTime,
		"nodeA002": inter.Timestamp(112),
		"nodeD002": genesisTime,
		"nodeA012": genesisTime,
		"nodeA003": inter.Timestamp(20),
		"nodeB002": inter.Timestamp(20),
		"nodeC002": inter.Timestamp(35),
	}
	t.Run("testMedianTimeOnDAG", func(t *testing.T) {
		testMedianTime(t, dagAscii, weights, creationTimes, medianTimes, genesisTime, LiteConfig())
		// test without cache
		testMedianTime(t, dagAscii, weights, creationTimes, medianTimes, genesisTime, Config{})
	})
}

func testMedianTime(t *testing.T, dagAscii string, weights []pos.Weight, creationTimes map[string]inter.Timestamp, medianTimes map[string]inter.Timestamp, genesis inter.Timestamp, cfg Config) {
	assertar := assert.New(t)

	var ordered dag.Events
	nodes, _, named := tdag.ASCIIschemeForEach(dagAscii, tdag.ForEachEvent{
		Process: func(e dag.Event, name string) {
			ordered = append(ordered, &eventWithCreationTime{e, creationTimes[name]})
		},
	})

	validators := pos.ArrayToValidators(nodes, weights)

	events := make(map[hash.Event]dag.Event)
	getEvent := func(id hash.Event) dag.Event {
		return events[id]
	}

	vi := NewIndex(makeTmpDB, cfg)
	vi.Reset(validators, getEvent)

	// push
	for _, e := range ordered {
		events[e.ID()] = e
		vi.Add(e)
		vi.Commit()
	}

	// check
	for name, e := range named {
		expected, ok := medianTimes[name]
		if !ok {
			continue
		}
		assertar.Equal(expected, vi.MedianTime(e.ID(), genesis), name)
	}
}

func makeTmpDB(name string) kvdb.Store {
	return memorydb.New()
}
