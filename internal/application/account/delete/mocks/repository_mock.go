// Package mocks contains testify mock implementations for the delete use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Repository is a testify mock for the delete.Repository interface.
type Repository struct {
	mock.Mock
}

// HasTransactions mocks Repository.HasTransactions.
func (m *Repository) HasTransactions(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// Delete mocks Repository.Delete.
func (m *Repository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
