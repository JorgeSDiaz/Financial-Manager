// Package summary implements the get transaction summary use case.
package summary

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

type Summary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

type Input struct {
	AccountID string `json:"account_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (uc *UseCase) Execute(ctx context.Context, in Input) (Summary, error) {
	incomes, err := uc.repo.ListByType(ctx, domaintransaction.TransactionTypeIncome, in.AccountID, "", in.StartDate, in.EndDate)
	if err != nil {
		return Summary{}, fmt.Errorf("get summary: %w", err)
	}

	expenses, err := uc.repo.ListByType(ctx, domaintransaction.TransactionTypeExpense, in.AccountID, "", in.StartDate, in.EndDate)
	if err != nil {
		return Summary{}, fmt.Errorf("get summary: %w", err)
	}

	var totalIncome, totalExpense float64
	for _, tx := range incomes {
		totalIncome += tx.Amount
	}
	for _, tx := range expenses {
		totalExpense += tx.Amount
	}

	return Summary{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
	}, nil
}
