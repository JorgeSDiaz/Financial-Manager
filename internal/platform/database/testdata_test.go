// Package database_test contains test helpers for the database package.
package database_test

import (
	"context"
	"database/sql"
)

// failingConnector is a test double that always returns a fixed error from Open.
type failingConnector struct {
	err error
}

// Open always returns the configured error.
func (f *failingConnector) Open(_ context.Context, _ string) (*sql.DB, error) {
	return nil, f.err
}

// buildFailingConnector creates a connector stub that returns err on every Open call.
func buildFailingConnector(err error) *failingConnector {
	return &failingConnector{err: err}
}
