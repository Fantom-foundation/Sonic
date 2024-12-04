package block_hash

//go:generate solc --bin block_ltypes.sol --abi block_ltypes.sol -o build --overwrite
//go:generate abigen --bin=build/BlockHash.bin --abi=build/BlockHash.abi --pkg=block_hash --out=block_ltypes.go
