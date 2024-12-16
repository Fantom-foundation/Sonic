package encryption

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Fantom-foundation/go-opera/utils/caution"
)

func writeTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", fmt.Errorf("failed to create keystore directory: %w", err)
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := os.CreateTemp(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary key file: %w", err)
	}

	if _, err = f.Write(content); err != nil {
		caution.CloseAndReportError(&err, f, "failed to close key file")
		caution.ExecuteAndReportError(&err, func() error { return os.Remove(f.Name()) },
			"failed to remove temporary key file")
		return "", err
	}

	return f.Name(), f.Close()
}
