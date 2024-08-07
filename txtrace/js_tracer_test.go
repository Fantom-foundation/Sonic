package txtrace

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/params"
)

// TestInitJSTracer Test for javascript tracer initialization
func TestInitJSTracer(t *testing.T) {

	tracerTests := []struct {
		name   string
		isOk   bool
		tracer string
	}{
		{"Build in CallTracer init", true, "callTracer"},
		{"Custom correct js code", true, `{data: [], fault: function(log) {}, step: function(log) { if(log.op.toString() == "CALL") this.data.push(log.stack.peek(0));}, result: function() { return this.data; }}`},
		{"Custom js code not compile", false, `{data: [], fault: function(log) {}, step: function(log) { if(log.op.toString() == "CALL") this.data.push(log.stack.peek(0));}, result: function() { return this.data; }`},
		{"Custom js code missing function", false, `{data: [], fault: function(log) {}, step: function(log) { if(log.op.toString() == "CALL") this.data.push(log.stack.peek(0)); }}`},
	}

	for _, tc := range tracerTests {
		t.Run(tc.name, func(t *testing.T) {
			tracer, errTracer := tracers.New(tc.tracer, &tracers.Context{})

			if tc.isOk && (tracer == nil || errTracer != nil) {
				t.Errorf("tracer creation must not fail but failed, error: %v", errTracer)
			} else if !tc.isOk && tracer != nil {
				t.Errorf("tracer creation must fail but did not fail")
			}
			if tracer != nil {
				tracer.Destroy()
			}
		})

	}
}

// TestCallTracerSimpleCall Create build in callTracer and run it on a simple
func TestCallTracerSimpleCall(t *testing.T) {

	callTracer := getJSTracer("callTracer", t)
	defer callTracer.Destroy()

	callTracer.CaptureStart(getEVMEnv(), from, to, false, inputData, 1000, value)
	callTracer.CaptureEnd(outputData, 100, time.Since(time.Now()), nil)

	want := `{
    "type": "CALL",
    "from": "0x0000000000000000000000000000000000000001",
    "to": "0x0000000000000000000000000000000000000002",
    "value": "0x5",
    "gas": "0x3e8",
    "gasUsed": "0x64",
    "input": "0x2f7468610000000000000000000000000000000000000000000000000000000000000008",
    "output": "0x45"
}`
	result, err := callTracer.GetResult()
	if err != nil {
		t.Fatalf("callTracer GetResult must not fail, error: %v", err.Error())
	}
	checkTracerResult(t, result, want)
}

// TestCallTracerComplexCall Create build in callTracer and run it on a complex inner calls
func TestCallTracerComplexCall(t *testing.T) {

	callTracer := getJSTracer("callTracer", t)
	defer callTracer.Destroy()

	executeCallbacks(callTracer)

	want := `{
    "type": "CALL",
    "from": "0x0000000000000000000000000000000000000001",
    "to": "0x0000000000000000000000000000000000000002",
    "value": "0x5",
    "gas": "0x3e8",
    "gasUsed": "0x64",
    "input": "0x2f7468610000000000000000000000000000000000000000000000000000000000000008",
    "output": "0x45",
    "calls": [
        {
            "type": "CREATE2",
            "from": "0x0000000000000000000000000000000000000002",
            "to": "0x0000000000000000000000000000000000000003",
            "value": "0x5",
            "gas": "0x258",
            "gasUsed": "0x190",
            "input": "0x2f7468610000000000000000000000000000000000000000000000000000000000000002",
            "output": "0x78",
            "calls": [
                {
                    "type": "CALL",
                    "from": "0x0000000000000000000000000000000000000002",
                    "to": "0x0000000000000000000000000000000000000003",
                    "value": "0x5",
                    "gas": "0x259",
                    "gasUsed": "0x191",
                    "input": "0x2f7468610000000000000000000000000000000000000000000000000000000000000002",
                    "output": "0x78"
                },
                {
                    "type": "STATICCALL",
                    "from": "0x0000000000000000000000000000000000000002",
                    "to": "0x0000000000000000000000000000000000000003",
                    "value": "0x5",
                    "gas": "0x25a",
                    "gasUsed": "0x192",
                    "input": "0x2f7468610000000000000000000000000000000000000000000000000000000000000002",
                    "output": "0x78"
                },
                {
                    "type": "DELEGATECALL",
                    "from": "0x0000000000000000000000000000000000000002",
                    "to": "0x0000000000000000000000000000000000000003",
                    "value": "0x5",
                    "gas": "0x25b",
                    "gasUsed": "0x193",
                    "input": "0x2f7468610000000000000000000000000000000000000000000000000000000000000002",
                    "output": "0x78"
                },
                {
                    "type": "SELFDESTRUCT",
                    "from": "0x0000000000000000000000000000000000000002",
                    "to": "0x0000000000000000000000000000000000000003",
                    "value": "0x5",
                    "gas": "0x25c",
                    "gasUsed": "0x194",
                    "input": "0x2f7468610000000000000000000000000000000000000000000000000000000000000002",
                    "output": "0x78"
                }
            ]
        },
        {
            "type": "CREATE",
            "from": "0x0000000000000000000000000000000000000002",
            "to": "0x0000000000000000000000000000000000000003",
            "value": "0x5",
            "gas": "0x258",
            "gasUsed": "0x196",
            "input": "0x2f7468610000000000000000000000000000000000000000000000000000000000000002",
            "output": "0x78"
        }
    ]
}`
	result, err := callTracer.GetResult()
	if err != nil {
		t.Fatalf("callTracer GetResult must not fail, error: %v", err.Error())
	}
	checkTracerResult(t, result, want)
}

// TestCustomTracerCodeComplexCall Create a custom js code for tracer and run it on a complex inner calls
func TestCustomTracerCodeComplexCall(t *testing.T) {

	tracer := getJSTracer(`{data: [], fault: function(log) {}, setup: function(config) { this.data.push(config.getType());}, enter: function(callFrame) { this.data.push(callFrame.getType());}, exit: function(frameResult) { this.data.push("exit");}, result: function() { return this.data; }}`, t)
	defer tracer.Destroy()
	executeCallbacks(tracer)

	want := `[
    "CREATE2",
    "CALL",
    "exit",
    "STATICCALL",
    "exit",
    "DELEGATECALL",
    "exit",
    "SELFDESTRUCT",
    "exit",
    "exit",
    "CREATE",
    "exit"
]`
	result, err := tracer.GetResult()
	if err != nil {
		t.Fatalf("tracer GetResult must not fail, error: %v", err.Error())
	}
	checkTracerResult(t, result, want)
}

// executeCallbacks Simulate EVM callbacks to tracer for complex inner call
func executeCallbacks(tracer *tracers.Tracer) {
	tracer.CaptureStart(getEVMEnv(), from, to, false, inputData, 1000, value)

	tracer.CaptureEnter(vm.CREATE2, to, toInner, inputDataInner, 600, value)

	tracer.CaptureEnter(vm.CALL, to, toInner, inputDataInner, 601, value)
	tracer.CaptureExit(outputDataInner, 401, nil)

	tracer.CaptureEnter(vm.STATICCALL, to, toInner, inputDataInner, 602, value)
	tracer.CaptureExit(outputDataInner, 402, nil)

	tracer.CaptureEnter(vm.DELEGATECALL, to, toInner, inputDataInner, 603, value)
	tracer.CaptureExit(outputDataInner, 403, nil)

	tracer.CaptureEnter(vm.SELFDESTRUCT, to, toInner, inputDataInner, 604, value)
	tracer.CaptureExit(outputDataInner, 404, nil)

	tracer.CaptureExit(outputDataInner, 400, nil)

	tracer.CaptureEnter(vm.CREATE, to, toInner, inputDataInner, 600, value)
	tracer.CaptureExit(outputDataInner, 406, nil)

	tracer.CaptureEnd(outputData, 100, time.Since(time.Now()), nil)
}

// getJSTracer Creates new debug tracer according to defined type or code in parameter tracer
func getJSTracer(tracer string, t *testing.T) *tracers.Tracer {
	jsTracer, errTracer := tracers.New(tracer, &tracers.Context{})
	if errTracer != nil {
		t.Fatalf("js tracer creation must not fail but did fail with error: %v", errTracer.Error())
	}
	return jsTracer
}

// getEVMEnv Create fake EVM environment
func getEVMEnv() *vm.EVM {
	return vm.NewEVM(
		vm.BlockContext{BlockNumber: blockNumber},
		vm.TxContext{GasPrice: gasprice},
		&vm.MockStateDB{},
		&params.ChainConfig{},
		vm.Config{})
}

// checkResult Compare result with expected string
func checkTracerResult(t *testing.T, result json.RawMessage, expectedResult string) {
	result, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		t.Errorf("problem with formating result, got error: %v", err)
	}

	if expectedResult != string(result) {
		t.Errorf("expected result is not the same as output got: %v, want: %v", string(result), expectedResult)
	}
}
