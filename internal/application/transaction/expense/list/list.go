// Package list implements the list expense transactions use case.
package list

import (
	"context"
	"fmt"

	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

type Repository interface {
	ListByType(ctx context.Context, tType domaintransaction.TransactionType, accountID, categoryID string, startDate, endDate string) ([]domaintransaction.Transaction, error)
}

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

type Input struct {
	AccountID  string `json:"account_id"`
	CategoryID string `json:"category_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

func (uc *UseCase) Execute(ctx context.Context, in Input) ([]domaintransaction.Transaction, error) {
	txs, err := uc.repo.ListByType(ctx, domaintransaction.TransactionTypeExpense, in.AccountID, in.CategoryID, in.StartDate, in.EndDate)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}

	return txs, nil
}
