package emitter

import (
	"testing"

	"github.com/Fantom-foundation/go-opera/gossip/emitter/mock"
	"github.com/Fantom-foundation/go-opera/vecmt"
	"github.com/Fantom-foundation/lachesis-base/emitter/ancestor"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/dag/tdag"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/Fantom-foundation/lachesis-base/kvdb/memorydb"
	"github.com/golang/mock/gomock"
)

func TestChooseParents(t *testing.T) {
	ctrl := gomock.NewController(t)
	external := mock.NewMockExternal(ctrl)
	em := NewEmitter(DefaultConfig(), World{External: external})
	em.maxParents = 3
	em.payloadIndexer = ancestor.NewPayloadIndexer(3)

	epoch := idx.Epoch(1)
	nodes := tdag.GenNodes(2)
	vEmpty, vNonEmpty := nodes[0], nodes[1]
	selfParentHash := sampleProcessedEvent(1, em.payloadIndexer)

	vi := vecmt.NewIndex(nil, vecmt.LiteConfig())
	vi.Reset(pos.ArrayToValidators(nodes, []pos.Weight{1, 1}), memorydb.New(), nil)

	external.EXPECT().GetLastEvent(epoch, vEmpty).
		Return(nil).
		Times(1)
	external.EXPECT().GetLastEvent(epoch, vNonEmpty).
		Return(&selfParentHash).
		Times(1)
	external.EXPECT().GetHeads(epoch).
		Return(hash.Events{sampleProcessedEvent(2, em.payloadIndexer), sampleProcessedEvent(3, em.payloadIndexer)}).
		Times(1)
	external.EXPECT().DagIndex().
		Return(vi).
		Times(1)

	t.Run("emptyChooseParents", func(t *testing.T) {
		selfParent, parents, ok := em.chooseParents(epoch, vEmpty)
		if selfParent != nil {
			t.Error("genesis event must not have self parent")
		}
		if len(parents) > 0 {
			t.Error("genesis event must not have any parents")
		}
		if !ok {
			t.Error("genesis parent assignment must always succeed")
		}
	})

	t.Run("nonEmptyChooseParents", func(t *testing.T) {
		selfParent, parents, ok := em.chooseParents(epoch, vNonEmpty)
		if selfParent == nil {
			t.Error("non-genesis event must have a self parent")
		}
		// strategies sometimes choose the same parent multiple times, test for minimal amount (1 self parent + 1 random/metric)
		if wantMin, got := 2, len(parents); got < wantMin {
			t.Errorf("incorrect number of event parents, expected at least: %d, got: %d", wantMin, got)
		}
		if !ok {
			t.Error("parent assignment must succeed when no cheating is detected")
		}
	})

}

func sampleProcessedEvent(id uint8, payloadIndexer *ancestor.PayloadIndexer) hash.Event {
	event := (&tdag.TestEvent{}).Build([24]byte{id})
	payloadIndexer.ProcessEvent(event, ancestor.Metric(id))
	return event.ID()
}
