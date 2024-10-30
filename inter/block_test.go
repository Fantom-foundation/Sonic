package inter

import (
	"github.com/Fantom-foundation/lachesis-base/hash"
	"math/rand"
	"testing"
)

func TestBlock_ComputePrevRandao_ComputationIsDeterministic(t *testing.T) {
	events := []hash.Event{
		{byte(rand.Int())},
		{byte(rand.Int())},
		{byte(rand.Int())},
		{byte(rand.Int())},
	}
	blk := Block{Events: events}
	randao1 := blk.GetPrevRandao()
	rand.Shuffle(len(blk.Events), func(i, j int) {
		blk.Events[i], blk.Events[j] = blk.Events[j], blk.Events[i]
	})
	randao2 := blk.GetPrevRandao()
	if randao1 != randao2 {
		t.Error("computation is not deterministic")
	}
}
