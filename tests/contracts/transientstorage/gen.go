package transientstorage

//go:generate solc --bin transientstorage.sol --abi transientstorage.sol -o build --overwrite
//go:generate abigen --bin=build/TransientStorage.bin --abi=build/TransientStorage.abi --pkg=transientstorage --out=transientstorage.go
