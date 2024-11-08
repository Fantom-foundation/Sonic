package gasprice

import (
	"math/big"

	"github.com/Fantom-foundation/go-opera/evmcore"
)

const (
	kElasticMultiplier           = 2
	kInitialBaseFee              = 1e9
	kBaseFeeMaxChangeDenominator = 8 // < TODO: adjust to lower block time
)

func GetInitialBaseFee() *big.Int {
	return big.NewInt(kInitialBaseFee)
}

func GetBaseFeeForNextBlock(parent *evmcore.EvmHeader) *big.Int {
	// This function implements the base-fee calculation algorithm of EIP-1559.

	newBaseFee := new(big.Int).Set(parent.BaseFee)

	gasTarget := parent.GasLimit / kElasticMultiplier
	gasUsed := parent.GasUsed

	if gasTarget == 0 {
		return newBaseFee
	}

	if gasUsed == gasTarget {
		// we don't need to adjust the base fee
	} else if gasUsed > gasTarget {
		gasUsedDelta := gasUsed - gasTarget

		baseFeeDelta := new(big.Int).Mul(parent.BaseFee, big.NewInt(int64(gasUsedDelta)))
		baseFeeDelta = new(big.Int).Div(baseFeeDelta, big.NewInt(int64(gasTarget)))
		baseFeeDelta = new(big.Int).Div(baseFeeDelta, big.NewInt(kBaseFeeMaxChangeDenominator))

		// The increase must be at least 1 wei.
		if baseFeeDelta.Sign() == 0 {
			baseFeeDelta.SetUint64(1)
		}
		newBaseFee = new(big.Int).Add(parent.BaseFee, baseFeeDelta)
	} else {
		gasUsedDelta := gasTarget - gasUsed
		baseFeeDelta := new(big.Int).Mul(parent.BaseFee, big.NewInt(int64(gasUsedDelta)))
		baseFeeDelta = new(big.Int).Div(baseFeeDelta, big.NewInt(int64(gasTarget)))
		baseFeeDelta = new(big.Int).Div(baseFeeDelta, big.NewInt(kBaseFeeMaxChangeDenominator))
		newBaseFee = new(big.Int).Sub(parent.BaseFee, baseFeeDelta)
	}

	return newBaseFee
}
