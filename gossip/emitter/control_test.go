package emitter

import (
	"testing"
	"time"

	"github.com/Fantom-foundation/go-opera/opera"
)

func TestGetEmitterIntervalLimit_IsOffWhenIntervalIsZero(t *testing.T) {
	rules := opera.EmitterRules{
		Interval:       0,
		StallThreshold: 200,
	}
	ms := time.Microsecond
	for _, delay := range []time.Duration{0, 100 * ms, 199 * ms, 200 * ms, 201 * ms} {
		_, enabled := getEmitterIntervalLimit(rules, delay)
		if enabled {
			t.Fatal("should be disabled")
		}
	}
}

func TestGetEmitterIntervalLimit_SwitchesToStallIfDelayed(t *testing.T) {
	ms := time.Millisecond
	regular := 100 * ms
	stallThreshold := 200 * ms
	stalled := 300 * ms

	rules := opera.EmitterRules{
		Interval:        uint64(regular.Milliseconds()),
		StallThreshold:  uint64(stallThreshold.Milliseconds()),
		StalledInterval: uint64(stalled.Milliseconds()),
	}

	for _, delay := range []time.Duration{0, 100 * ms, 199 * ms, 200 * ms, 201 * ms, 60 * time.Minute} {
		got, enabled := getEmitterIntervalLimit(rules, delay)
		if !enabled {
			t.Fatalf("should be enabled for delay %v", delay)
		}
		want := regular
		if delay > stallThreshold {
			want = stalled
		}
		if want != got {
			t.Fatalf("for delay %v, want %v, got %v", delay, want, got)
		}
	}
}
