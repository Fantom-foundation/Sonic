package burn_gas

//go:generate solc --bin burn_gas.sol --abi burn_gas.sol -o build --overwrite
//go:generate abigen --bin=build/BurnGas.bin --abi=build/BurnGas.abi --pkg=burn_gas --out=burn_gas.go
