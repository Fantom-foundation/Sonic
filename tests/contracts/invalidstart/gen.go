package invalidstart

//go:generate solc --bin invalidstart.sol --abi invalidstart.sol -o build --overwrite
//go:generate abigen --bin=build/InvalidStart.bin --abi=build/InvalidStart.abi --pkg=invalidstart --out=invalidstart.go
