package update_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/transaction/update/mocks"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const fixedTimestamp = "2026-02-28T10:00:00Z"
const fixedDate = "2026-02-28"

// seeded is the canonical existing transaction used as the pre-update state.
var seeded = buildTransaction("tx-1", "acc-001", "cat-001", "Old Description", 100.0)

// buildTransaction returns a valid Transaction for use in tests.
func buildTransaction(id, accountID, categoryID, description string, amount float64) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", fixedDate)
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   accountID,
		CategoryID:  categoryID,
		Type:        domaintransaction.TransactionTypeIncome,
		Amount:      amount,
		Description: description,
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

// buildMockRepoFull creates a mocks.Repository pre-configured for one GetByID and one Update call.
func buildMockRepoFull(id string, fetched, updated domaintransaction.Transaction, updateErr error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("GetByID", mock.Anything, id).Return(fetched, nil).Once()
	m.On("Update", mock.Anything, updated).Return(updateErr).Once()
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
