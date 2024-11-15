package counter_event_emitter

//go:generate solc --bin counter_event_emitter.sol --abi counter_event_emitter.sol -o build --overwrite
//go:generate abigen --bin=build/counter_event_emitter.bin --abi=build/counter_event_emitter.abi --pkg=counter_event_emitter --out=counter_event_emitter.go
