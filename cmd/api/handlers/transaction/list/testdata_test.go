package list_test

import (
	"context"
	"time"

	expenselist "github.com/financial-manager/api/internal/application/transaction/expense/list"
	incomelist "github.com/financial-manager/api/internal/application/transaction/income/list"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const fixedTimestamp = "2026-02-28T10:00:00Z"

type fakeIncomeLister struct {
	out []domaintransaction.Transaction
	err error
}

func (f *fakeIncomeLister) Execute(_ context.Context, _ incomelist.Input) ([]domaintransaction.Transaction, error) {
	return f.out, f.err
}

type fakeExpenseLister struct {
	out []domaintransaction.Transaction
	err error
}

func (f *fakeExpenseLister) Execute(_ context.Context, _ expenselist.Input) ([]domaintransaction.Transaction, error) {
	return f.out, f.err
}

func buildDomainTransaction(id, accountID string, txType domaintransaction.TransactionType) domaintransaction.Transaction {
	t, _ := time.Parse("2006-01-02T15:04:05Z", fixedTimestamp)
	date, _ := time.Parse("2006-01-02", "2026-02-28")
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   accountID,
		CategoryID:  "cat-001",
		Type:        txType,
		Amount:      100.0,
		Description: "Test transaction",
		Date:        date,
		IsActive:    true,
		CreatedAt:   t,
		UpdatedAt:   t,
	}
}
