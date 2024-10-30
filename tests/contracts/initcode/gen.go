package initcode

//go:generate solc --bin initcode.sol --abi initcode.sol -o build --overwrite
//go:generate abigen --bin=build/InitCode.bin --abi=build/InitCode.abi --pkg=initcode --out=initcode.go
