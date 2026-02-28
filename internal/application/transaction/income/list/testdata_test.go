package list_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/transaction/income/list/mocks"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const fixedDate = "2026-02-28"

// buildMockRepo creates a mocks.Repository pre-configured to return the given
// transactions and error for one ListByType call.
func buildMockRepo(txs []domaintransaction.Transaction, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("ListByType", mock.Anything, domaintransaction.TransactionTypeIncome, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(txs, err).Once()
	return m
}

// income1 and income2 are canonical income transaction fixtures for list tests.
var (
	income1 = buildIncome("tx-1", "acc-001", "Salary", 1000.0)
	income2 = buildIncome("tx-2", "acc-002", "Freelance", 500.0)
)

// buildIncome returns a valid income Transaction for use in tests.
func buildIncome(id, accountID, description string, amount float64) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", fixedDate)
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   accountID,
		CategoryID:  "cat-001",
		Type:        domaintransaction.TransactionTypeIncome,
		Amount:      amount,
		Description: description,
		Date:        date,
		IsActive:    true,
	}
}
