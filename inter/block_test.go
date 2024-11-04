package inter

import (
	"bytes"
	"crypto/sha256"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/ethereum/go-ethereum/common"
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

func TestBlock_ComputePrevRandao_ComputationWorks(t *testing.T) {
	tests := []struct {
		name   string
		events hash.Events
	}{
		{
			name:   "empty_events",
			events: []hash.Event{},
		},
		{
			name: "one_events",
			events: []hash.Event{
				{byte(rand.Int())},
			},
		},
		{
			name: "multiple_events",
			events: []hash.Event{
				{byte(rand.Int())},
				{byte(rand.Int())},
				{byte(rand.Int())},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prevRandao := common.Hash{}
			for _, event := range test.events {
				for i := 0; i < 24; i++ {
					prevRandao[i+8] = prevRandao[i+8] ^ event[i+8]
				}
			}
			blk := Block{Events: test.events}
			if want, got := sha256.Sum256(prevRandao.Bytes()), blk.GetPrevRandao(); !bytes.Equal(want[:], got[:]) {
				t.Errorf("unexpected has;, got: %s, want: %s", got, want)
			}
		})
	}
}
