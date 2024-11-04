package accessCost

//go:generate solc --bin access_cost.sol --abi access_cost.sol -o build --overwrite
//go:generate abigen --bin=build/AccessCost.bin --abi=build/AccessCost.abi --pkg=accessCost --out=access_cost.go
