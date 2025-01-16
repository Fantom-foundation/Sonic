package errlock

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Fantom-foundation/go-opera/utils/caution"
)

// Check if errlock is written
func Check() error {
	locked, reason, eLockPath, err := read(datadir)
	if err != nil {
		return fmt.Errorf("Node isn't allowed to start due to an error reading"+
			" the lock file %s.\n Please fix the issue. Error message:\n%w",
			eLockPath, err)
	}

	if locked {
		return fmt.Errorf("Node isn't allowed to start due to a previous error."+
			" Please fix the issue and then delete file \"%s\". Error message:\n%s",
			eLockPath, reason)
	}
	return nil
}

var (
	datadir string
)

// SetDefaultDatadir for errlock files
func SetDefaultDatadir(dir string) {
	datadir = dir
}

// Permanent error
func Permanent(err error) {
	eLockPath, _ := write(datadir, err.Error())
	panic(fmt.Errorf("Node is permanently stopping due to an issue. Please fix"+
		" the issue and then delete file \"%s\". Error message:\n%s",
		eLockPath, err.Error()))
}

func readAll(reader io.Reader, max int) ([]byte, error) {
	buf := make([]byte, max)
	consumed := 0
	for {
		n, err := reader.Read(buf[consumed:])
		consumed += n
		if consumed == len(buf) || err == io.EOF {
			return buf[:consumed], nil
		}
		if err != nil {
			return nil, err
		}
	}
}

// read errlock file
func read(dir string) (locked bool, reason string, eLockPath string, err error) {
	eLockPath = path.Join(dir, "errlock")

	data, err := os.Open(eLockPath)
	if err != nil {
		// if file doesn't exist, directory is not locked and it is not an error
		return false, "", eLockPath, nil
	}
	defer caution.CloseAndReportError(&err, data, "Failed to close errlock file")

	// read no more than N bytes
	maxFileLen := 5000
	eLockBytes, err := readAll(data, maxFileLen)
	if err != nil {
		return true, "", eLockPath, fmt.Errorf("failed to read lock file %v: %w", eLockPath, err)
	}
	return true, string(eLockBytes), eLockPath, nil
}

// write errlock file
func write(dir string, eLockStr string) (string, error) {
	eLockPath := path.Join(dir, "errlock")

	return eLockPath, os.WriteFile(eLockPath, []byte(eLockStr), 0666) // assume no custom encoding needed
}
