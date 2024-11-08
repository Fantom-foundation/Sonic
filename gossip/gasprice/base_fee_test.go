package gasprice

import (
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/evmcore"
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