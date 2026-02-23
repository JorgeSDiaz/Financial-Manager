// Package migrator_test provides test helpers for the migrator package.
package migrator_test

import (
	"errors"
	"io/fs"
	"testing/fstest"
)

// errorReadFileFS wraps a MapFS and returns an error when ReadFile is called
// on a file that matches the given target filename.
type errorReadFileFS struct {
	inner  fstest.MapFS
	target string
}

func (e errorReadFileFS) Open(name string) (fs.File, error) {
	return e.inner.Open(name)
}

func (e errorReadFileFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(e.inner, name)
}

func (e errorReadFileFS) ReadFile(name string) ([]byte, error) {
	if name == e.target {
		return nil, errors.New("simulated read error")
	}

	return fs.ReadFile(e.inner, name)
}

// buildErrorReadFileFS creates an errorReadFileFS that fails on the given target file.
func buildErrorReadFileFS(inner fstest.MapFS, target string) errorReadFileFS {
	return errorReadFileFS{inner: inner, target: target}
}
