// Package list implements the list income transactions use case.
package list

import (
	"context"
	"fmt"

	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Input is the input for the list income use case.
type Input struct {
	AccountID  string `json:"account_id"`
	CategoryID string `json:"category_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

type Repository interface {
	ListByType(ctx context.Context, tType domaintransaction.TransactionType, accountID, categoryID string, startDate, endDate string) ([]domaintransaction.Transaction, error)
}

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, in Input) ([]domaintransaction.Transaction, error) {
	txs, err := uc.repo.ListByType(ctx, domaintransaction.TransactionTypeIncome, in.AccountID, in.CategoryID, in.StartDate, in.EndDate)
	if err != nil {
		return nil, fmt.Errorf("list incomes: %w", err)
	}

	return txs, nil
}
