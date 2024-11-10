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
	return GetInitialBaseFee_Sonic()
	//return GetInitialBaseFee_Eth()
}

func GetBaseFeeForNextBlock(parent *evmcore.EvmHeader) *big.Int {
	//fmt.Printf("GetBaseFee for block after %d @ %d\n", parent.Number, parent.Time)
	return GetBaseFeeForNextBlock_Sonic(parent, big.NewInt(5e5)) // 1M Gas/s - TODO: adjust to real network rate
	//return GetBaseFeeForNextBlock_Eth(parent)
}

// --- Ethereum model ---

func GetInitialBaseFee_Eth() *big.Int {
	return big.NewInt(kInitialBaseFee)
}

func GetBaseFeeForNextBlock_Eth(parent *evmcore.EvmHeader) *big.Int {
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

// --- Sonic model ---

func GetInitialBaseFee_Sonic() *big.Int {
	return big.NewInt(0)
}

func GetBaseFeeForNextBlock_Sonic(parent *evmcore.EvmHeader, targetRate *big.Int) *big.Int {

	// newPrice := oldPrice * e^(((rate-targetRate/targetRate)*duration)/128)

	duration := parent.Duration
	if duration == 0 || duration > 60*1e9 {
		return parent.BaseFee
	}

	nanosPerSecond := big.NewInt(1e9)
	usedGas := big.NewInt(int64(parent.GasUsed))

	durationInNanos := big.NewInt(int64(duration)) // 63-bit is enough for a duration of 292 years

	numerator := sub(mul(usedGas, nanosPerSecond), mul(targetRate, durationInNanos))
	denominator := mul(big.NewInt(128), mul(targetRate, nanosPerSecond))

	oldPrice := parent.BaseFee
	newPrice := approximateExponential(oldPrice, numerator, denominator)

	// If the gas rate is higher than the target, increase the price by at least 1 wei.
	// This is to ensure that the price is always increasing, even if the old price was 0.
	if oldPrice.Cmp(newPrice) == 0 && numerator.Sign() > 0 {
		newPrice.Add(newPrice, big.NewInt(1))
	}

	return newPrice
}

// approximateExponential approximates f * e ** (n/d) using
// Taylor expansion at a=0:
// f * e^(n/d) = f + af/b + a^2f/b^2/2! + a^3f/b^3/3! + ...
func approximateExponential(factor, numerator, denominator *big.Int) *big.Int {
	var (
		output = new(big.Int)
		accum  = new(big.Int).Mul(factor, denominator)
	)
	for i := 1; accum.Sign() != 0; i++ {
		output.Add(output, accum)
		accum.Mul(accum, numerator)
		accum.Div(accum, denominator)
		accum.Div(accum, big.NewInt(int64(i)))
	}
	return output.Div(output, denominator)
}

func sub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func mul(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}
