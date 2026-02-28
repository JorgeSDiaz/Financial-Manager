// Package mocks contains testify mock implementations for the create income use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Repository is a testify mock for the create.Repository interface.
type Repository struct {
	mock.Mock
}

// Create mocks Repository.Create.
func (m *Repository) Create(ctx context.Context, t domaintransaction.Transaction) error {
	return m.Called(ctx, t).Error(0)
}
