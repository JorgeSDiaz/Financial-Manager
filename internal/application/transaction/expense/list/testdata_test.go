package list_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/transaction/expense/list/mocks"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const fixedDate = "2026-02-28"

// buildMockRepo creates a mocks.Repository pre-configured to return the given
// transactions and error for one ListByType call.
func buildMockRepo(txs []domaintransaction.Transaction, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("ListByType", mock.Anything, domaintransaction.TransactionTypeExpense, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(txs, err).Once()
	return m
}

// expense1 and expense2 are canonical expense transaction fixtures for list tests.
var (
	expense1 = buildExpense("tx-1", "acc-001", "Groceries", 100.0)
	expense2 = buildExpense("tx-2", "acc-002", "Utilities", 200.0)
)

// buildExpense returns a valid expense Transaction for use in tests.
func buildExpense(id, accountID, description string, amount float64) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", fixedDate)
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   accountID,
		CategoryID:  "cat-001",
		Type:        domaintransaction.TransactionTypeExpense,
		Amount:      amount,
		Description: description,
		Date:        date,
		IsActive:    true,
	}
}
