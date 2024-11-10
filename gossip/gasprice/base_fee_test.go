package gasprice

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/inter"
)

func TestBaseFee_PriceAdjustments(t *testing.T) {

	tests := map[string]struct {
		parentBaseFee  uint64
		parentGasUsed  uint64
		parentGasLimit uint64
		wantBaseFee    uint64
	}{
		"base fee remains the same": {
			parentBaseFee:  1e9,
			parentGasUsed:  1e6,
			parentGasLimit: 2e6,
			wantBaseFee:    1e9,
		},
		"base fee increases": {
			parentBaseFee:  1e9,
			parentGasUsed:  2e6,
			parentGasLimit: 2e6,
			wantBaseFee:    1e9 + 1e9/8, // +12.5%
		},
		"base fee decreases": {
			parentBaseFee:  1e9,
			parentGasUsed:  0,
			parentGasLimit: 2e6,
			wantBaseFee:    1e9 - 1e9/8, // -12.5%
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			header := &evmcore.EvmHeader{
				BaseFee:  big.NewInt(int64(test.parentBaseFee)),
				GasUsed:  test.parentGasUsed,
				GasLimit: test.parentGasLimit,
			}

			gotBaseFee := GetBaseFeeForNextBlock(header)
			wantBaseFee := big.NewInt(int64(test.wantBaseFee))
			if gotBaseFee.Cmp(wantBaseFee) != 0 {
				t.Fatalf("base fee is incorrect; got %v, want %v", gotBaseFee, wantBaseFee)
			}
		})
	}

	// Test the base fee price adjustments.
	// The base fee is adjusted based on the gas used in the previous block.
	// If the gas used is equal to the target gas, the base fee remains the same.
	// If the gas used is greater than the target gas, the base fee increases.
	// If the gas used is less than the target gas, the base fee decreases.
	// The base fee is adjusted by a fraction of the parent base fee.
	// The adjustment is capped at 1/8 of the parent base fee.
	// The base fee is adjusted by at least 1 wei.
	// The base fee is capped at 2^64 - 1 wei.
	// The base fee is initialized to 1e9 wei.
	// The target gas is the parent gas limit divided by 2.
	// The gas used is the parent gas used.
	// The base fee is calculated using the formula:
	// newBaseFee = parentBaseFee + parentBaseFee * adjustment
	// where adjustment = gasUsedDelta / gasTarget / kBaseFeeMaxChangeDenominator
	// gasUsedDelta = gasUsed - gasTarget
	// gasTarget = parentGasLimit / kElasticMultiplier
	// The base fee is calculated using big.Int to avoid overflow.
	// The base fee is returned as a big.Int.
	// The base fee is returned as a pointer to a big.Int.
	// The base fee is returned as a new big.Int
	// The base fee is returned as a copy of the parent base fee.
	// The base fee is returned as a new big.Int with the parent base fee.

}

// TODO:
//  - make sure gas price can grow again if it is zero

func TestBaseFee_SonicPriceAdjustments(t *testing.T) {

	approxExp := func(f, n, d int64) uint64 {
		return uint64(float64(f) * math.Exp(float64(n)/float64(d)))
	}

	tests := map[string]struct {
		parentBaseFee  uint64
		parentGasUsed  uint64
		parentDuration uint64
		targetRate     uint64
		wantBaseFee    uint64
	}{
		"base fee remains the same": {
			parentBaseFee:  1e8,
			parentGasUsed:  1e6,
			parentDuration: 1e9, // 1 second
			targetRate:     1e6,
			wantBaseFee:    approxExp(1e8, 0, 128), // max change rate per second ~1/128
		},
		"base fee increases": {
			parentBaseFee:  1e8,
			parentGasUsed:  2e6,
			parentDuration: 1e9, // 1 second
			targetRate:     1e6,
			wantBaseFee:    approxExp(1e8, 1, 128),
		},
		"base fee decreases": {
			parentBaseFee:  1e8,
			parentGasUsed:  0,
			parentDuration: 1e9, // 1 second
			targetRate:     1e6,
			wantBaseFee:    approxExp(1e8, -1, 128),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			header := &evmcore.EvmHeader{
				BaseFee:  big.NewInt(int64(test.parentBaseFee)),
				GasUsed:  test.parentGasUsed,
				Duration: inter.Duration(test.parentDuration),
			}

			gotBaseFee := GetBaseFeeForNextBlock_Sonic(header, big.NewInt(int64(test.targetRate)))
			wantBaseFee := big.NewInt(int64(test.wantBaseFee))
			if gotBaseFee.Cmp(wantBaseFee) != 0 {
				t.Fatalf("base fee is incorrect; got %v, want %v, diff %d", gotBaseFee, wantBaseFee, sub(gotBaseFee, wantBaseFee))
			}
		})
	}
}

func TestBaseFee_SonicPriceCanRecoverFromPriceZero(t *testing.T) {

	target := uint64(1e6)
	header := &evmcore.EvmHeader{
		BaseFee:  big.NewInt(0),
		GasUsed:  target + 1,
		Duration: inter.Duration(1e9), // 1 second
	}

	newPrice := GetBaseFeeForNextBlock_Sonic(header, big.NewInt(int64(target)))
	if newPrice.Cmp(big.NewInt(1)) < 0 {
		t.Errorf("failed to increase price from zero, new price %v", newPrice)
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
		GasLimit: 2e6,
		Duration: inter.Duration(1e9),
	}
	for range b.N {
		GetBaseFeeForNextBlock(header)
	}
}
