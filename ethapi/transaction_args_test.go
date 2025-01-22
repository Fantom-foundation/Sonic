package ethapi

import (
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestGasCap(t *testing.T) {

	var (
		gas10M       = hexutil.Uint64(10_000_000)
		gasMaxUint64 = hexutil.Uint64(math.MaxUint64)
	)

	tests := []struct {
		name         string
		argGas       *hexutil.Uint64
		globalGasCap uint64
		expectedGas  uint64
	}{
		{
			name:         "gas cap 0 and arg gas nil",
			globalGasCap: 0,
			argGas:       nil,
			expectedGas:  math.MaxInt64,
		}, {
			name:         "gas cap 0 and arg gas 10M",
			globalGasCap: 0,
			argGas:       &gas10M,
			expectedGas:  10_000_000,
		}, {
			name:         "gas cap 0 and arg gas maxUint64",
			globalGasCap: 0,
			argGas:       &gasMaxUint64,
			expectedGas:  math.MaxInt64,
		}, {
			name:         "gas cap 50M and arg gas nil",
			globalGasCap: 50_000_000,
			argGas:       nil,
			expectedGas:  50_000_000,
		}, {
			name:         "gas cap 50M and arg gas 10M",
			globalGasCap: 50_000_000,
			argGas:       &gas10M,
			expectedGas:  10_000_000,
		}, {
			name:         "gas cap 50M and arg gas maxUint64",
			globalGasCap: 50_000_000,
			argGas:       &gasMaxUint64,
			expectedGas:  50_000_000,
		}, {
			name:         "gas cap maxUint64 and arg gas 10M",
			globalGasCap: math.MaxUint64,
			argGas:       &gas10M,
			expectedGas:  10_000_000,
		}, {
			name:         "gas cap maxUint64 and arg gas maxUint64",
			globalGasCap: math.MaxUint64,
			argGas:       &gasMaxUint64,
			expectedGas:  math.MaxInt64,
		},
	}

	for _, test := range tests {
		args := TransactionArgs{Gas: test.argGas}

		msg, err := args.ToMessage(test.globalGasCap, nil)

		require.Nil(t, err)
		require.Equal(t, test.expectedGas, msg.GasLimit, test.name)
	}
}
