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
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/tests"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var (
	blockTestDir                   = filepath.Join(baseDir, "BlockchainTests")
	executionSpecBlockchainTestDir = filepath.Join(".", "spec-tests", "fixtures", "blockchain_tests")
)

func TestBlockchain(t *testing.T) {
	bt := new(tests.TestMatcher)

	// Very slow test
	bt.SkipLoad(`.*/stTimeConsuming/.*`)
	// test takes a lot for time and goes easily OOM because of sha3 calculation on a huge range,
	// using 4.6 TGas
	bt.SkipLoad(`.*randomStatetest94.json.*`)

	bt.Walk(t, blockTestDir, func(t *testing.T, name string, test *tests.BlockTest) {
		execBlockTest(t, bt, test)
	})
}

// TestExecutionSpecBlocktests runs the test fixtures from execution-spec-tests.
func TestExecutionSpecBlocktests(t *testing.T) {
	if !common.FileExist(executionSpecBlockchainTestDir) {
		t.Skipf("directory %s does not exist", executionSpecBlockchainTestDir)
	}
	bt := new(tests.TestMatcher)

	bt.Walk(t, executionSpecBlockchainTestDir, func(t *testing.T, name string, test *tests.BlockTest) {
		execBlockTest(t, bt, test)
	})
}

func execBlockTest(t *testing.T, bt *tests.TestMatcher, test *tests.BlockTest) {
	if err := bt.CheckFailure(t, test.Run(false, rawdb.HashScheme, true, nil, nil)); err != nil {
		t.Errorf("test with config failed: %v", err)
	}
}
