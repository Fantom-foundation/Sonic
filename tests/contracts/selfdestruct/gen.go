package selfdestruct

//go:generate solc --bin selfdestruct.sol --abi selfdestruct.sol -o build --overwrite
//go:generate abigen --bin=build/SelfDestruct.bin --abi=build/SelfDestruct.abi --pkg=selfdestruct --out=selfdestruct.go
