// package caution provides utility functions for handling errors and closing
// resources.
// This package should be only be used for sitations where errors are unexpected
// and should be ignored or panic

package caution

import (
	"fmt"
	"io"
)

func ExecuteAndReportError(err *error, f func() error, message string) {
	if *err == nil {
		*err = f()
		if *err != nil {
			*err = fmt.Errorf("%s: %w", message, *err)
		}
	}
}

func CloseAndReportError(err *error, closer io.Closer, message string) {
	ExecuteAndReportError(err, func() error {
		return closer.Close()
	}, message)
}
