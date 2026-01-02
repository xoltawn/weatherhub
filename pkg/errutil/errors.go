package errutil

import (
	"fmt"
)

// Wrap returns a new error with a formatted message that wraps the original error.
// If the original error (err) is nil, it returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf is like Wrap but supports format specifiers.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	message := fmt.Sprintf(format, args...)

	return fmt.Errorf("%s: %w", message, err)
}
