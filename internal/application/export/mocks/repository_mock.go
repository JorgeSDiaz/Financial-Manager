// Package mocks contains testify mock implementations for the export use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Repository is a testify mock for the export.Repository interface.
type Repository struct {
	mock.Mock
}

// ListAccounts mocks Repository.ListAccounts.
func (m *Repository) ListAccounts(ctx context.Context) ([]domainaccount.Account, error) {
	args := m.Called(ctx)
	accounts, _ := args.Get(0).([]domainaccount.Account)
	return accounts, args.Error(1)
}

// ListCategories mocks Repository.ListCategories.
func (m *Repository) ListCategories(ctx context.Context) ([]domaincategory.Category, error) {
	args := m.Called(ctx)
	categories, _ := args.Get(0).([]domaincategory.Category)
	return categories, args.Error(1)
}

// ListTransactions mocks Repository.ListTransactions.
func (m *Repository) ListTransactions(ctx context.Context, tType domaintransaction.TransactionType, startDate, endDate string) ([]domaintransaction.Transaction, error) {
	args := m.Called(ctx, tType, startDate, endDate)
	transactions, _ := args.Get(0).([]domaintransaction.Transaction)
	return transactions, args.Error(1)
}
