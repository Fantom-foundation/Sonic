package block_hash

//go:generate solc --bin block_hash.sol --abi block_hash.sol -o build --overwrite
//go:generate abigen --bin=build/BlockHash.bin --abi=build/BlockHash.abi --pkg=block_hash --out=block_hash.go
