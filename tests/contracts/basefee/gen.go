package basefee

//go:generate solc --bin basefee.sol --abi basefee.sol -o build --overwrite
//go:generate abigen --bin=build/BaseFee.bin --abi=build/BaseFee.abi --pkg=basefee --out=basefee.go
