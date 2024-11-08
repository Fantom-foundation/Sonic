package selfdestruct

//go:generate solc --bin selfdestruct.sol --abi selfdestruct.sol -o build --overwrite
//go:generate abigen --bin=build/SelfDestruct.bin --abi=build/SelfDestruct.abi --pkg=selfdestruct --type=SelfDestruct --out=selfdestruct.go
//go:generate abigen --bin=build/SelfDestructFactory.bin --abi=build/SelfDestructFactory.abi --pkg=selfdestruct --type=SelfDestructFactory --out=selfdestruct_factory.go
