package create_test

import (
	"context"
	"time"

	appCreate "github.com/financial-manager/api/internal/application/transaction/income/create"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const fixedTimestamp = "2026-02-28T10:00:00Z"

type fakeUseCase struct {
	out domaintransaction.Transaction
	err error
}

func (f *fakeUseCase) Execute(_ context.Context, _ appCreate.Input) (domaintransaction.Transaction, error) {
	return f.out, f.err
}

func buildDomainTransaction(id, accountID string, amount float64) domaintransaction.Transaction {
	t, _ := time.Parse("2006-01-02T15:04:05Z", fixedTimestamp)
	date, _ := time.Parse("2006-01-02", "2026-02-28")
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   accountID,
		CategoryID:  "cat-001",
		Type:        domaintransaction.TransactionTypeIncome,
		Amount:      amount,
		Description: "Test income",
		Date:        date,
		IsActive:    true,
		CreatedAt:   t,
		UpdatedAt:   t,
	}
}
