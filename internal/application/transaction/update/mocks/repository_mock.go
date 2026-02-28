// Package mocks contains testify mock implementations for the update use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Repository is a testify mock for the update.Repository interface.
type Repository struct {
	mock.Mock
}

// GetByID mocks Repository.GetByID.
func (m *Repository) GetByID(ctx context.Context, id string) (domaintransaction.Transaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domaintransaction.Transaction), args.Error(1)
}

// Update mocks Repository.Update.
func (m *Repository) Update(ctx context.Context, t domaintransaction.Transaction) error {
	return m.Called(ctx, t).Error(0)
}
