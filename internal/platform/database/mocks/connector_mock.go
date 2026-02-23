// Package mocks contains testify mock implementations for the database package.
package mocks

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

// Connector is a testify mock for the database connector interface.
type Connector struct {
	mock.Mock
}

// Open mocks connector.Open.
func (m *Connector) Open(ctx context.Context, path string) (*sql.DB, error) {
	args := m.Called(ctx, path)

	db, _ := args.Get(0).(*sql.DB)

	return db, args.Error(1)
}
