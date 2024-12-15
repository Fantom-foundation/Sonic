package caution

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestExecuteAndReportError_DoesNothingWhenInputErrorIsNotNil(t *testing.T) {
	err := fmt.Errorf("someError")
	ExecuteAndReportError(&err, func() error {
		panic("mypanic")
	}, "mymessage")
	if err.Error() != "someError" {
		t.Errorf("error not reported, got %v", err)
	}
}

func TestExecuteAndReportError_ExecutesFunAndReportsError(t *testing.T) {
	var err error
	ExecuteAndReportError(&err, func() error {
		return fmt.Errorf("someError")
	}, "message")
	if errors.Is(err, fmt.Errorf("message: someError")) {
		t.Errorf("unexpected error, got %v", err)
	}
}

func TestExecuteAndReportError_ExecutesFunNothingToReport(t *testing.T) {
	var err error
	ExecuteAndReportError(&err, func() error {
		return nil
	}, "message")
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}
}

func TestCloseAndReportError_(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "file")
	file, err := os.Create(tmpfile)
	if err != nil {
		t.Fatal(err)
	}

	CloseAndReportError(&err, file, "message")
	if _, err := file.Read([]byte{0}); err == nil {
		t.Errorf("file not closed, %v", err)
	}
}
