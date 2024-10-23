package coinbase

//go:generate solc --bin coinbase.sol --abi coinbase.sol -o build --overwrite
//go:generate abigen --bin=build/Coinbase.bin --abi=build/Coinbase.abi --pkg=coinbase --out=coinbase.go
