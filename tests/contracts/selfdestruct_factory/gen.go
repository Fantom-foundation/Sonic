package selfdestruct_factory

//go:generate solc --bin selfdestruct_factory.sol --abi selfdestruct_factory.sol -o build --overwrite --base-path ..
//go:generate abigen --bin=build/SelfDestructFactory.bin --abi=build/SelfDestructFactory.abi --pkg=selfdestruct_factory --out=selfdestruct_factory.go
