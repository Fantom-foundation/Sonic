package gaspricelimits

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestAddPercentage_AddsRequestedAmountOfPercent(t *testing.T) {
	for v := range 10000 {
		for p := range 20 {
			res := int(addPercentage(big.NewInt(int64(v)), p).Int64())
			expected := v + v*p/100
			fmt.Printf("v: %d, p: %d, res: %d, expected: %d\n", v, p, res, expected)
			if res != expected {
				t.Errorf("Expected %d, got %d", expected, res)
			}
		}
	}
}

func TestAddPercentage_TreatesNilLikeZero(t *testing.T) {
	res := addPercentage(nil, 10)
	if res.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected 0, got %d", res)
	}
}

func TestAddPercentage_CanHandleValuesLargerThanMaxUint64(t *testing.T) {
	value := big.NewInt(math.MaxInt64)
	value = value.Mul(value, value)

	extra := new(big.Int).Div(value, big.NewInt(10))
	want := new(big.Int).Add(value, extra)
	got := addPercentage(value, 10)

	if got.Cmp(want) != 0 {
		t.Errorf("Expected %d, got %d", want, got)
	}
}

func TestGetSuggestedGasPriceForNewTransactions_ReturnsValue10PercentHigherThanBaseFee(t *testing.T) {
	for v := range []int{0, 1, 5, 10, 100, 1000, 10000} {
		baseFee := big.NewInt(int64(v))
		want := addPercentage(baseFee, 10)
		got := GetSuggestedGasPriceForNewTransactions(baseFee)
		if got.Cmp(want) != 0 {
			t.Errorf("Expected %d, got %d", want, got)
		}
	}
}

func TestGetMinimumFeeCapForTransactionPool_ReturnsValue5PercentHigherThanBaseFee(t *testing.T) {
	for v := range []int{0, 1, 5, 10, 100, 1000, 10000} {
		baseFee := big.NewInt(int64(v))
		want := addPercentage(baseFee, 5)
		got := GetMinimumFeeCapForTransactionPool(baseFee)
		if got.Cmp(want) != 0 {
			t.Errorf("Expected %d, got %d", want, got)
		}
	}
}

func TestGetMinimumFeeCapForEventEmitterPool_ReturnsValue2PercentHigherThanBaseFee(t *testing.T) {
	for v := range []int{0, 1, 5, 10, 100, 1000, 10000} {
		baseFee := big.NewInt(int64(v))
		want := addPercentage(baseFee, 5)
		got := GetMinimumFeeCapForEventEmitter(baseFee)
		if got.Cmp(want) != 0 {
			t.Errorf("Expected %d, got %d", want, got)
		}
	}
}
