package query_account

//go:generate solc --bin query_account.sol --abi query_account.sol -o build --overwrite
//go:generate abigen --bin=build/QueryAccount.bin --abi=build/QueryAccount.abi --pkg=query_account --out=query_account.go
