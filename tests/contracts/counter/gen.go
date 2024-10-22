package counter

//go:generate solc --bin counter.sol --abi counter.sol -o build --overwrite
//go:generate abigen --bin=build/Counter.bin --abi=build/Counter.abi --pkg=counter --out=counter.go
