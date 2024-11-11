package gasprice

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/opera"
)

func TestBaseFee_ExamplePriceAdjustments(t *testing.T) {

	approxExp := func(f, n, d int64) uint64 {
		return uint64(float64(f) * math.Exp(float64(n)/float64(d)))
	}

	tests := map[string]struct {
		parentBaseFee  uint64
		parentGasUsed  uint64
		parentDuration time.Duration
		targetRate     uint64
		wantBaseFee    uint64
	}{
		"base fee remains the same": {
			parentBaseFee:  1e8,
			parentGasUsed:  1e6,
			parentDuration: 1 * time.Second,
			targetRate:     1e6,
			wantBaseFee:    approxExp(1e8, 0, 128), // max change rate per second ~1/128
		},
		"base fee increases": {
			parentBaseFee:  1e8,
			parentGasUsed:  2e6,
			parentDuration: 1 * time.Second,
			targetRate:     1e6,
			wantBaseFee:    approxExp(1e8, 1, 128),
		},
		"base fee decreases": {
			parentBaseFee:  1e8,
			parentGasUsed:  0,
			parentDuration: 1 * time.Second,
			targetRate:     1e6,
			wantBaseFee:    approxExp(1e8, -1, 128),
		},
		"long durations are ignored": {
			parentBaseFee:  123456789,
			parentGasUsed:  0, // < no gas used, should reduce the price
			parentDuration: 61 * time.Second,
			targetRate:     1e6, // < since the duration is too long, the price should not change
			wantBaseFee:    123456789,
		},
		"target rate is zero": {
			parentBaseFee:  123456789,
			parentGasUsed:  0, // < no gas used, should reduce the price
			parentDuration: time.Second,
			targetRate:     0, // < since the target rate is zero, the price should not change
			wantBaseFee:    123456789,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			header := &evmcore.EvmHeader{
				BaseFee:  big.NewInt(int64(test.parentBaseFee)),
				GasUsed:  test.parentGasUsed,
				Duration: inter.Duration(test.parentDuration),
			}

			rules := opera.EconomyRules{
				ShortGasPower: opera.GasPowerRules{
					AllocPerSec: 2 * test.targetRate,
				},
			}

			gotBaseFee := GetBaseFeeForNextBlock(header, rules)
			wantBaseFee := big.NewInt(int64(test.wantBaseFee))
			if gotBaseFee.Cmp(wantBaseFee) != 0 {
				t.Fatalf("base fee is incorrect; got %v, want %v, diff %d", gotBaseFee, wantBaseFee, sub(gotBaseFee, wantBaseFee))
			}

			if header.BaseFee == gotBaseFee {
				t.Fatalf("new base fee is not a copy; got %p, want %p", header.BaseFee, gotBaseFee)
			}
		})
	}
}

func TestBaseFee_PriceCanRecoverFromPriceZero(t *testing.T) {

	target := uint64(1e6)
	header := &evmcore.EvmHeader{
		BaseFee:  big.NewInt(0),
		GasUsed:  target + 1,
		Duration: inter.Duration(1e9), // 1 second
	}

	rules := opera.EconomyRules{
		ShortGasPower: opera.GasPowerRules{
			AllocPerSec: 2 * target,
		},
	}

	newPrice := GetBaseFeeForNextBlock(header, rules)
	if newPrice.Cmp(big.NewInt(1)) < 0 {
		t.Errorf("failed to increase price from zero, new price %v", newPrice)
	}
}

func TestBaseFee_DecayTimeFromInitialToZeroIsApproximately35Minutes(t *testing.T) {
	rules := opera.EconomyRules{
		ShortGasPower: opera.GasPowerRules{
			AllocPerSec: 1e6,
		},
	}

	// This property should be true for any block time.
	blockTimes := []time.Duration{
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
		2 * time.Second,
		5 * time.Second,
	}
	for _, blockTime := range blockTimes {
		t.Run(fmt.Sprintf("blockTime=%s", blockTime.String()), func(t *testing.T) {
			header := &evmcore.EvmHeader{
				BaseFee:  GetInitialBaseFee(),
				GasUsed:  0,
				Duration: inter.Duration(blockTime),
			}
			decayDuration := time.Duration(0)
			for header.BaseFee.Sign() > 0 {
				header.BaseFee = GetBaseFeeForNextBlock(header, rules)
				decayDuration += header.Duration.Duration()
			}

			if decayDuration < 30*time.Minute || decayDuration > 40*time.Minute {
				t.Errorf("time to decay from initial to zero is incorrect; got %v", decayDuration)
			}
		})
	}
}

func TestApproximateExponential_KnownValues(t *testing.T) {
	tests := map[string]struct {
		factor      int64
		numerator   int64
		denominator int64
		want        int64
	}{
		"e^0": {
			factor:      1,
			numerator:   0,
			denominator: 1,
			want:        1,
		},
		"e^1": {
			factor:      1,
			numerator:   1,
			denominator: 1,
			want:        2,
		},
		"e^2": {
			factor:      1,
			numerator:   2,
			denominator: 1,
			want:        6, // < should be 7, but the function just approximates
		},
		"e^(1/2)": {
			factor:      1,
			numerator:   1,
			denominator: 2,
			want:        1,
		},
		"e^-1": {
			factor:      1,
			numerator:   -1,
			denominator: 1,
			want:        0,
		},
		"10*e^2": {
			factor:      10,
			numerator:   2,
			denominator: 1,
			want:        71, // < should be 73, but the function just approximates
		},
		"100*e^(1/2)": {
			factor:      100,
			numerator:   1,
			denominator: 2,
			want:        164,
		},
		"100*e^(-1/2)": {
			factor:      100,
			numerator:   -1,
			denominator: 2,
			want:        60,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotBaseFee := approximateExponential(
				big.NewInt(int64(test.factor)),
				big.NewInt(test.numerator),
				big.NewInt(test.denominator),
			)
			wantBaseFee := big.NewInt(int64(test.want))
			if gotBaseFee.Cmp(wantBaseFee) != 0 {
				t.Fatalf("base fee is incorrect; got %v, want %v", gotBaseFee, wantBaseFee)
			}
		})
	}
}

func TestApproximateExponential_RandomInputs(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	for range 100 {
		factor := int64(r.Int31n(100))

		// Our practical use cases for gas computations are fractions in the range [-1,1]
		denominator := int64(r.Int31n(1e9) + 1)
		numerator := int64((r.Float64()*2 - 1) * float64(denominator))

		want := big.NewInt(int64(float64(factor) * math.Exp(float64(numerator)/float64(denominator))))
		got := approximateExponential(big.NewInt(factor), big.NewInt(numerator), big.NewInt(denominator))

		diff := new(big.Int).Abs(sub(got, want))
		if diff.Cmp(big.NewInt(1)) > 0 {
			t.Errorf(
				"incorrect approximation for f=%d, n=%d, d=%d; got %v, want %v, error %v",
				factor, numerator, denominator, got, want, diff,
			)
		}
	}
}

func BenchmarkBaseFeeComputation(b *testing.B) {
	header := &evmcore.EvmHeader{
		BaseFee:  big.NewInt(1e9),
		GasUsed:  1e6,
		Duration: inter.Duration(1e9),
	}
	rules := opera.EconomyRules{
		ShortGasPower: opera.GasPowerRules{
			AllocPerSec: 1e6,
		},
	}
	for range b.N {
		GetBaseFeeForNextBlock(header, rules)
	}
}
