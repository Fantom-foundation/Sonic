package gaspricelimits

import "math/big"

// GetSuggestedGasPriceForNewTransactions a gas price that should be suggested
// to users for new transactions based on the current base fee. This function
// returns a value that is 10% higher than the base fee to provide a buffer
// for the price to increase before the transaction is included in a block.
func GetSuggestedGasPriceForNewTransactions(baseFee *big.Int) *big.Int {
	return addPercentage(baseFee, 10)
}

// GetMinimumFeeCapForTransactionPool returns the gas price the transaction pool
// should check for when accepting new transactions. This function returns a
// value that is 5% higher than the base fee to provide a buffer for the price
// to increase before the transaction is included in a block.
func GetMinimumFeeCapForTransactionPool(baseFee *big.Int) *big.Int {
	return addPercentage(baseFee, 5)
}

// GetMinimumFeeCapForEventEmitter returns the gas price the event emitter should
// check for when including transactions in a block. This function returns a
// value that is 2% higher than the base fee to provide a buffer for the price
// to increase before the transaction is included in a block.
func GetMinimumFeeCapForEventEmitter(baseFee *big.Int) *big.Int {
	return addPercentage(baseFee, 2)
}

func addPercentage(a *big.Int, percentage int) *big.Int {
	if a == nil {
		return big.NewInt(0)
	}
	res := new(big.Int).Set(a)
	res.Mul(res, big.NewInt(int64(percentage+100)))
	res.Div(res, big.NewInt(100))
	return res
}
