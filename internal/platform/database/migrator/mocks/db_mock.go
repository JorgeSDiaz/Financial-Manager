// Package mocks contains testify mock implementations for the migrator package.
package mocks

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

// DB is a testify mock for the migrator db interface.
type DB struct {
	mock.Mock
}

// ExecContext mocks db.ExecContext.
func (m *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	varArgs := []any{ctx, query}
	varArgs = append(varArgs, args...)
	called := m.Called(varArgs...)

	res, _ := called.Get(0).(sql.Result)

	return res, called.Error(1)
}

// QueryRowContext mocks db.QueryRowContext.
func (m *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	varArgs := []any{ctx, query}
	varArgs = append(varArgs, args...)
	called := m.Called(varArgs...)

	row, _ := called.Get(0).(*sql.Row)

	return row
}
