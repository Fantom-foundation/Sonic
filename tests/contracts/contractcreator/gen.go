package contractcreator

//go:generate solc --bin contractcreator.sol --abi contractcreator.sol -o build --overwrite
//go:generate abigen --bin=build/ContractCreator.bin --abi=build/ContractCreator.abi --pkg=contractcreator --out=contractcreator.go
