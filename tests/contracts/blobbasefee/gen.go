package blobbasefee

//go:generate solc --bin blobbasefee.sol --abi blobbasefee.sol -o build --overwrite
//go:generate abigen --bin=build/BlobBaseFee.bin --abi=build/BlobBaseFee.abi --pkg=blobbasefee --out=blobbasefee.go
