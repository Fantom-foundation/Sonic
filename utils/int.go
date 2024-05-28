package utils

import (
	"math/big"

	"github.com/holiman/uint256"
)

func BigIntToUint256(value *big.Int) *uint256.Int {
	if value.Sign() < 0 {
		panic("unable to convert negative big.Int to uint256")
	}
	bytes := value.Bytes()
	if len(bytes) > 32 {
		panic("unable to convert big.Int exceeding 32 bytes to uint256")
	}
	return new(uint256.Int).SetBytes(bytes)
}

func Uint256ToBigInt(value *uint256.Int) *big.Int {
	return new(big.Int).SetBytes(value.Bytes())
}
