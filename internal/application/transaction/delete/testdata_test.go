package delete_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/transaction/delete/mocks"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const fixedTimestamp = "2026-02-28T10:00:00Z"
const fixedDate = "2026-02-28"

// buildTransaction returns a valid Transaction for use in tests.
func buildTransaction(id string) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", fixedDate)
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   "acc-001",
		CategoryID:  "cat-001",
		Type:        domaintransaction.TransactionTypeIncome,
		Amount:      100.0,
		Description: "Test transaction",
		Date:        date,
		IsActive:    true,
	}
}

// buildMockRepoGetByID creates a mocks.Repository pre-configured for one GetByID call.
func buildMockRepoGetByID(id string, tx domaintransaction.Transaction, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("GetByID", mock.Anything, id).Return(tx, err).Once()
	return m
}

// buildMockRepoDelete creates a mocks.Repository pre-configured for one GetByID
// (returning transaction) and one SoftDelete call.
func buildMockRepoDelete(id string, deleteErr error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("GetByID", mock.Anything, id).Return(buildTransaction(id), nil).Once()
	m.On("SoftDelete", mock.Anything, id).Return(deleteErr).Once()
	return m
}
