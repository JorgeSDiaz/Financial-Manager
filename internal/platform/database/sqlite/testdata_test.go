// Package sqlite_test contains test helpers for the sqlite package.
package sqlite_test

import (
	"os"
	"testing"
)

// createFile creates a regular file at the given path to be used in tests
// that need to simulate a path collision where a file blocks directory creation.
func createFile(t *testing.T, path string) error {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	return f.Close()
}
