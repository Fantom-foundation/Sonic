package inter

import (
	"bytes"
	"crypto/sha256"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/ethereum/go-ethereum/common"
	"math/rand"
	"testing"
)

func TestComputePrevRandao_ComputationIsDeterministic(t *testing.T) {
	events := hash.FakeEvents(5)
	randao1 := ComputePrevRandao(events)
	rand.Shuffle(len(events), func(i, j int) {
		events[i], events[j] = events[j], events[i]
	})
	randao2 := ComputePrevRandao(events)
	if randao1 != randao2 {
		t.Error("computation is not deterministic")
	}
}

func TestComputePrevRandao_ComputationWorks(t *testing.T) {
	tests := []struct {
		name   string
		events hash.Events
	}{
		{
			name:   "empty_events",
			events: []hash.Event{},
		},
		{
			name:   "one_events",
			events: hash.FakeEvents(1),
		},
		{
			name:   "multiple_events",
			events: hash.FakeEvents(3),
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
			if want, got := sha256.Sum256(prevRandao.Bytes()), ComputePrevRandao(test.events); !bytes.Equal(want[:], got[:]) {
				t.Errorf("unexpected has;, got: %s, want: %s", got, want)
			}
		})
	}
}
