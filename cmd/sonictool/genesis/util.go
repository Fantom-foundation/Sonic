package genesis


import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/fdlimit"
	"os"
	"path"
	"path/filepath"
)

const (
	// DefaultCacheSize is calculated as memory consumption in a worst case scenario with default configuration
	// Average memory consumption might be 3-5 times lower than the maximum
	DefaultCacheSize  = 3600
	ConstantCacheSize = 400
)

// makeDatabaseHandles raises out the number of allowed file handles per process and returns allowance for db.
func makeDatabaseHandles() uint64 {
	limit, err := fdlimit.Maximum()
	if err != nil {
		panic(fmt.Errorf("failed to retrieve file descriptor allowance: %v", err))
	}
	raised, err := fdlimit.Raise(uint64(limit))
	if err != nil {
		panic(fmt.Errorf("failed to raise file descriptor allowance: %v", err))
	}
	return raised / 6 + 1
}

func removeDatabase(dataDir string) error {
	err1 := os.RemoveAll(filepath.Join(dataDir, "chaindata"))
	err2 := os.RemoveAll(filepath.Join(dataDir, "carmen"))
	err3 := os.RemoveAll(filepath.Join(dataDir, "errlock"))
	return errors.Join(err1, err2, err3)
}

func setGenesisProcessing(chaindataDir string) {
	f, _ := os.Create(path.Join(chaindataDir, "unfinished"))
	if f != nil {
		_ = f.Close()
	}
}

func setGenesisComplete(chaindataDir string) {
	_ = os.Remove(path.Join(chaindataDir, "unfinished"))
}
