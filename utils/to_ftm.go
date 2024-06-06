package utils

import (
	"math/big"

	"github.com/holiman/uint256"
)

// ToFtm number of FTM to Wei
func ToFtm(ftm uint64) *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(ftm), big.NewInt(1e18))
}

// ToFtmU256 number of FTM to Wei using the uint256 type
func ToFtmU256(ftm uint64) *uint256.Int {
	return BigIntToUint256(ToFtm(ftm))
}
