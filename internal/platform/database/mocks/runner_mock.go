// Package mocks contains testify mock implementations for the database package.
package mocks

import (
	"context"
	"database/sql"
	"io/fs"

	"github.com/stretchr/testify/mock"
)

// Runner is a testify mock for the database runner interface.
type Runner struct {
	mock.Mock
}

// Run mocks runner.Run.
func (m *Runner) Run(ctx context.Context, db *sql.DB, migrationFS fs.FS, dir string) error {
	args := m.Called(ctx, db, migrationFS, dir)
	return args.Error(0)
}
