// Package mocks contains testify mock implementations for the delete use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Repository is a testify mock for the delete.Repository interface.
type Repository struct {
	mock.Mock
}

// GetByID mocks Repository.GetByID.
func (m *Repository) GetByID(ctx context.Context, id string) (domaintransaction.Transaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domaintransaction.Transaction), args.Error(1)
}

// SoftDelete mocks Repository.SoftDelete.
func (m *Repository) SoftDelete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
