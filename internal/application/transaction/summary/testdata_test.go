package summary_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/transaction/summary/mocks"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const fixedDate = "2026-02-28"

// buildMockRepoIncome creates a mocks.Repository pre-configured to return the given
// income transactions for one ListByType call.
func buildMockRepoIncome(txs []domaintransaction.Transaction, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("ListByType", mock.Anything, domaintransaction.TransactionTypeIncome, mock.Anything, "", mock.Anything, mock.Anything).Return(txs, err).Once()
	return m
}

// buildMockRepoExpense creates a mocks.Repository pre-configured to return the given
// expense transactions for one ListByType call.
func buildMockRepoExpense(txs []domaintransaction.Transaction, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("ListByType", mock.Anything, domaintransaction.TransactionTypeExpense, mock.Anything, "", mock.Anything, mock.Anything).Return(txs, err).Once()
	return m
}

// income100 and income500 are income transaction fixtures for summary tests.
var (
	income100 = buildIncome("tx-1", 100.0)
	income500 = buildIncome("tx-2", 500.0)
)

// expense50 and expense200 are expense transaction fixtures for summary tests.
var (
	expense50  = buildExpense("tx-3", 50.0)
	expense200 = buildExpense("tx-4", 200.0)
)

// buildIncome returns a valid income Transaction for use in tests.
func buildIncome(id string, amount float64) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", fixedDate)
	return domaintransaction.Transaction{
		ID:       id,
		Type:     domaintransaction.TransactionTypeIncome,
		Amount:   amount,
		Date:     date,
		IsActive: true,
	}
}

// buildExpense returns a valid expense Transaction for use in tests.
func buildExpense(id string, amount float64) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", fixedDate)
	return domaintransaction.Transaction{
		ID:       id,
		Type:     domaintransaction.TransactionTypeExpense,
		Amount:   amount,
		Date:     date,
		IsActive: true,
	}
}
