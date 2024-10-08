// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/tests"
	"path/filepath"
	"reflect"
	"testing"
)

func initMatcher(st *testMatcher) {
	//// Long tests:
	//st.slow(`^stAttackTest/ContractCreationSpam`)
	//st.slow(`^stBadOpcode/badOpcodes`)
	//st.slow(`^stPreCompiledContracts/modexp`)
	//st.slow(`^stQuadraticComplexityTest/`)
	//st.slow(`^stStaticCall/static_Call50000`)
	//st.slow(`^stStaticCall/static_Return50000`)
	//st.slow(`^stSystemOperationsTest/CallRecursiveBomb`)
	//st.slow(`^stTransactionTest/Opcodes_TransactionInit`)
	//// Very time consuming
	//st.skipLoad(`^stTimeConsuming/`)
	//st.skipLoad(`.*vmPerformance/loop.*`)
	//// Uses 1GB RAM per tested fork
	//st.skipLoad(`^stStaticCall/static_Call1MB`)
	//
	//// Broken tests:
	//// EOF is not part of cancun
	st.skipLoad(`^stEOF/`)
	//
	//// The tests under Pyspecs are the ones that are published as execution-spec tests.
	//// We run these tests separately, no need to _also_ run them as part of the
	//// reference tests.
	//st.skipLoad(`^Pyspecs/`)
}

func TestState(t *testing.T) {
	t.Parallel()

	st := new(testMatcher)
	initMatcher(st)
	for _, dir := range []string{
		filepath.Join(baseDir, "EIPTests", "StateTests"),
		stateTestDir,
		benchmarksDir,
	} {
		st.walk(t, dir, func(t *testing.T, name string, test *tests.StateTest) {
			execStateTest(t, st, test)
		})
	}
}

func execStateTest(t *testing.T, st *testMatcher, test *tests.StateTest) {
	for _, subtest := range test.Subtests() {
		subtest := subtest
		key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)

		t.Run(key, func(t *testing.T) {
			withTrace(t, 0, func(vmconfig vm.Config) error {
				var result error
				test.Run(subtest, vmconfig, false, rawdb.HashScheme, func(err error, state *tests.StateTestState) {
					result = st.checkFailure(t, err)
				})
				return result
			})
		})
	}
}

// Transactions with gasLimit above this value will not get a VM trace on failure.
const traceErrorLimit = 400000

func withTrace(t *testing.T, gasLimit uint64, test func(vm.Config) error) {
	// Use config from command line arguments.
	config := vm.Config{}
	err := test(config)
	if err == nil {
		return
	}

	// Test failed, re-run with tracing enabled.
	t.Error(err)
	if gasLimit > traceErrorLimit {
		t.Log("gas limit too high for EVM trace")
		return
	}
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	config.Tracer = logger.NewJSONLogger(&logger.Config{}, w)
	err2 := test(config)
	if !reflect.DeepEqual(err, err2) {
		t.Errorf("different error for second run: %v", err2)
	}
	w.Flush()
	if buf.Len() == 0 {
		t.Log("no EVM operation logs generated")
	} else {
		t.Log("EVM operation log:\n" + buf.String())
	}
}