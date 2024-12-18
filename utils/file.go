package utils

import (
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/log"
)

func OpenFile(path string, isSyncMode bool) *os.File {
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		log.Crit("Failed to create file dir", "file", path, "err", err)
	}
	sync := 0
	if isSyncMode {
		sync = os.O_SYNC
	}
	fileHandle, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|sync, 0666)
	if err != nil {
		log.Crit("Failed to open file", "file", path, "err", err)
	}
	return fileHandle
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
