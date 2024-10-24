package tests

import (
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
	// original test:
	//if err := bt.CheckFailure(t, test.Run(false, rawdb.HashScheme, true, nil, nil)); err != nil {
	//	t.Errorf("test with config failed: %v", err)
	//}

	sonicTest := BlockTest{*test}
	factory := createCarmenFactory(t)
	if err := bt.CheckFailure(t, sonicTest.Run(factory)); err != nil {
		t.Errorf("test with config failed: %v", err)
	}

}
