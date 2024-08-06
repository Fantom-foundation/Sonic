package vecclock

import (
	"fmt"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io/ioutil"
	"testing"

	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/dag"
	"github.com/Fantom-foundation/lachesis-base/inter/dag/tdag"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
)

var (
	testASCIIScheme = `
a1.0   b1.0   c1.0   d1.0   e1.0
║      ║      ║      ║      ║
║      ╠──────╫───── d2.0   ║
║      ║      ║      ║      ║
║      b2.1 ──╫──────╣      e2.1
║      ║      ║      ║      ║
║      ╠──────╫───── d3.1   ║
a2.1 ──╣      ║      ║      ║
║      ║      ║      ║      ║
║      b3.2 ──╣      ║      ║
║      ║      ║      ║      ║
║      ╠──────╫───── d4.2   ║
║      ║      ║      ║      ║
║      ╠───── c2.2   ║      e3.2
║      ║      ║      ║      ║
`
)

type eventWithCreationTime struct {
	dag.Event
	creationTime inter.Timestamp
}

func (e *eventWithCreationTime) CreationTime() inter.Timestamp {
	return e.creationTime
}

func BenchmarkIndex_Add(b *testing.B) {
	b.StopTimer()
	ordered := make(dag.Events, 0)
	nodes, _, _ := tdag.ASCIIschemeForEach(testASCIIScheme, tdag.ForEachEvent{
		Process: func(e dag.Event, name string) {
			ordered = append(ordered, e)
		},
	})
	validatorsBuilder := pos.NewBuilder()
	for _, peer := range nodes {
		validatorsBuilder.Set(peer, 1)
	}
	validators := validatorsBuilder.Build()
	events := make(map[hash.Event]dag.Event)
	getEvent := func(id hash.Event) dag.Event {
		return events[id]
	}
	for _, e := range ordered {
		events[e.ID()] = e
	}

	vecClock := NewIndex(makeTmpDB, LiteConfig())
	vecClock.Reset(validators, getEvent)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		vecClock.Reset(validators, getEvent)
		b.StartTimer()
		for _, e := range ordered {
			vecClock.Add(&eventWithCreationTime{e, inter.Timestamp(e.Seq())})
			vecClock.Commit()
			i++
			if i >= b.N {
				break
			}
		}
	}
}

func BenchmarkIndex_Add_MemoryDB(b *testing.B) {
	benchmark_Index_Add(b)
}

func benchmark_Index_Add(b *testing.B) {
	b.StopTimer()

	nodes := tdag.GenNodes(70)
	ordered := make(dag.Events, 0)
	tdag.ForEachRandEvent(nodes, 10, 10, nil, tdag.ForEachEvent{
		Process: func(e dag.Event, name string) {
			ordered = append(ordered, e)
		},
	})

	validatorsBuilder := pos.NewBuilder()
	for _, peer := range nodes {
		validatorsBuilder.Set(peer, 1)
	}
	validators := validatorsBuilder.Build()
	events := make(map[hash.Event]dag.Event)
	getEvent := func(id hash.Event) dag.Event {
		return events[id]
	}
	for _, e := range ordered {
		events[e.ID()] = e
	}

	i := 0
	for {
		b.StopTimer()
		vecClock := NewIndex(makeTmpDB, LiteConfig())
		vecClock.Reset(validators, getEvent)
		b.StartTimer()
		for _, e := range ordered {
			vecClock.Add(e)
			vecClock.Commit()
			i++
			if i >= b.N {
				return
			}
		}
	}
}

func tempLevelDB() (kvdb.Store, error) {
	cache16mb := func(string) (int, int) {
		return 16 * opt.MiB, 64
	}
	dir, err := ioutil.TempDir("", "bench")
	if err != nil {
		panic(fmt.Sprintf("can't create temporary directory %s: %v", dir, err))
	}
	disk := leveldb.NewProducer(dir, cache16mb)
	ldb, _ := disk.OpenDB("0")
	return ldb, nil
}
