// Package mocks contains testify mock implementations for the list use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Repository is a testify mock for the list.Repository interface.
type Repository struct {
	mock.Mock
}

// ListByType mocks Repository.ListByType.
func (m *Repository) ListByType(ctx context.Context, tType domaintransaction.TransactionType, accountID, categoryID string, startDate, endDate string) ([]domaintransaction.Transaction, error) {
	args := m.Called(ctx, tType, accountID, categoryID, startDate, endDate)
	return args.Get(0).([]domaintransaction.Transaction), args.Error(1)
}
