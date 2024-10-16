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
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/Carmen/go/state/gostate"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/tests"
	"path/filepath"
	"reflect"
	"testing"
)

func initMatcher(st *testMatcher) {
	st.skipLoad(`^stEOF/`)
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
				factory := createCarmenFactory(t)
				var result error
				test.RunWith(subtest, vmconfig, factory, func(err error, state *tests.StateTestState) {
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

// createCarmenFactory creates a new factory, that initialises
// carmen implementation of the state database.
func createCarmenFactory(t *testing.T) carmenFactory {
	dir := t.TempDir()
	parameters := carmen.Parameters{
		Variant:   gostate.VariantGoMemory,
		Schema:    carmen.Schema(5),
		Archive:   carmen.NoArchive,
		Directory: dir,
	}

	st, err := carmen.NewState(parameters)
	if err != nil {
		t.Fatalf("cannot create state: %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("cannot close state: %v", err)
		}
	})

	return carmenFactory{st: st}
}
