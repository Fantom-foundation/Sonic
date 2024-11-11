package gasprice

import (
	"math/big"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/opera"
)

// GetInitialBaseFee returns the initial base fee to be used in the genesis block.
func GetInitialBaseFee() *big.Int {
	// The initial base fee is set to 1 Gwei. While a value of 0 would also be valid,
	// this value was chosen to have non-zero prices in low-load test networks at least
	// for the first several minutes. In case of no load on the network, the base fee
	// will decrease to 0 within ~35 minutes.
	const kInitialBaseFee = 1e9
	return big.NewInt(kInitialBaseFee)
}

// GetBaseFeeForNextBlock computes the base fee for the next block based on the parent block.
func GetBaseFeeForNextBlock(parent *evmcore.EvmHeader, rules opera.EconomyRules) *big.Int {
	// In general, this function computes the new base fee based on the following formula:
	//
	//     newPrice := oldPrice * e^(((rate-targetRate/targetRate)*duration)/128)
	//
	// where:
	//   - oldPrice is the base fee of the parent block
	//   - rate is the gas rate per second observed in the parent block
	//   - targetRate is the target gas rate per second at which prices are stable
	//   - duration is the time in seconds between the parent and grand-parent blocks
	//
	// All computations are carried out using integers to avoid floating point errors.
	// To that end, terms are re-arranged to fit the following shape:
	//
	//               newPrice := oldPrice * e^(numerator/denominator)
	//
	// where numerator and denominator are integers. The final value is then computed
	// using an approximation of the this function based on a Taylor expansion around 0.

	oldPrice := new(big.Int).Set(parent.BaseFee)

	// If the time gap between the parent and this block is more than
	// 60 seconds, something significantly disturbed the chain and we
	// keep the BaseFee constant.
	duration := parent.Duration
	if duration == 0 || duration > 60*1e9 {
		return oldPrice
	}

	// If the target rate is zero, the new price is not defined.
	targetRate := big.NewInt(int64(rules.ShortGasPower.AllocPerSec / 2))
	if targetRate.Sign() == 0 {
		return oldPrice
	}

	nanosPerSecond := big.NewInt(1e9)
	usedGas := big.NewInt(int64(parent.GasUsed))

	durationInNanos := big.NewInt(int64(duration)) // 63-bit is enough for a duration of 292 years

	numerator := sub(mul(usedGas, nanosPerSecond), mul(targetRate, durationInNanos))
	denominator := mul(big.NewInt(128), mul(targetRate, nanosPerSecond))

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
		res = new(big.Int)
		acc = new(big.Int).Mul(factor, denominator)
	)
	for i := 1; acc.Sign() != 0; i++ {
		res.Add(res, acc)
		acc.Mul(acc, numerator)
		acc.Div(acc, denominator)
		acc.Div(acc, big.NewInt(int64(i)))
	}
	return res.Div(res, denominator)
}

func sub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func mul(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}
