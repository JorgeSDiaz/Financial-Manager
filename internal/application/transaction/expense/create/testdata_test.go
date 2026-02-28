package create_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/transaction/expense/create/mocks"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const (
	fixedID        = "fixed-uuid-0001"
	fixedTimestamp = "2026-02-28T10:00:00Z"
	fixedDate      = "2026-02-28"
)

// validExpense is the expected expense transaction produced by a successful create.
var validExpense = domaintransaction.Transaction{
	ID:          fixedID,
	AccountID:   "acc-001",
	CategoryID:  "cat-001",
	Type:        domaintransaction.TransactionTypeExpense,
	Amount:      100.0,
	Description: "Groceries",
	Date:        fixedDateOnly(),
	IsActive:    true,
	CreatedAt:   fixedTime(),
	UpdatedAt:   fixedTime(),
}

// errorExpense is the transaction passed to the repo when input is minimal.
var errorExpense = domaintransaction.Transaction{
	ID:        fixedID,
	AccountID: "acc-001",
	Type:      domaintransaction.TransactionTypeExpense,
	Amount:    100.0,
	Date:      fixedDateOnly(),
	IsActive:  true,
	CreatedAt: fixedTime(),
	UpdatedAt: fixedTime(),
}

// buildMockRepo creates a mocks.Repository pre-configured to accept one Create call.
func buildMockRepo(t domaintransaction.Transaction, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("Create", mock.Anything, t).Return(err).Once()
	return m
}

// buildMockIDGenerator creates a mocks.IDGenerator pre-configured to return fixedID once.
func buildMockIDGenerator() *mocks.IDGenerator {
	m := &mocks.IDGenerator{}
	m.On("NewID").Return(fixedID).Once()
	return m
}

// buildMockClock creates a mocks.Clock pre-configured to return fixedTime once.
func buildMockClock() *mocks.Clock {
	m := &mocks.Clock{}
	m.On("Now").Return(fixedTime()).Once()
	return m
}

// fixedTime parses fixedTimestamp and panics on error (test helper).
func fixedTime() time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", fixedTimestamp)
	if err != nil {
		panic(err)
	}
	return t.UTC()
}

// fixedDateOnly returns the date portion only.
func fixedDateOnly() time.Time {
	t, err := time.Parse("2006-01-02", fixedDate)
	if err != nil {
		panic(err)
	}
	return t
}
