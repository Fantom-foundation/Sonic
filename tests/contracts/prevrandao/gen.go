package prevrandao

//go:generate solc --bin prevrandao.sol --abi prevrandao.sol -o build --overwrite
//go:generate abigen --bin=build/Prevrandao.bin --abi=build/Prevrandao.abi --pkg=prevrandao --out=prevrandao.go
