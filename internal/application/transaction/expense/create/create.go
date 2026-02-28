// Package create implements the create expense transaction use case.
package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

type Repository interface {
	Create(ctx context.Context, t domaintransaction.Transaction) error
}

type IDGenerator interface {
	NewID() string
}

type Clock interface {
	Now() time.Time
}

type UseCase struct {
	repo  Repository
	idGen IDGenerator
	clock Clock
}

func New(repo Repository, idGen IDGenerator, clock Clock) *UseCase {
	return &UseCase{repo: repo, idGen: idGen, clock: clock}
}

type Input struct {
	AccountID   string  `json:"account_id"`
	CategoryID  string  `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func (uc *UseCase) Execute(ctx context.Context, in Input) (domaintransaction.Transaction, error) {
	if err := validateInput(in); err != nil {
		return domaintransaction.Transaction{}, err
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		return domaintransaction.Transaction{}, errors.New("invalid date format, use YYYY-MM-DD")
	}

	now := uc.clock.Now().UTC()
	tx := domaintransaction.Transaction{
		ID:          uc.idGen.NewID(),
		AccountID:   in.AccountID,
		CategoryID:  in.CategoryID,
		Type:        domaintransaction.TransactionTypeExpense,
		Amount:      in.Amount,
		Description: in.Description,
		Date:        date,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.repo.Create(ctx, tx); err != nil {
		return domaintransaction.Transaction{}, fmt.Errorf("create expense: %w", err)
	}

	return tx, nil
}

func validateInput(in Input) error {
	if in.AccountID == "" {
		return errors.New("account_id is required")
	}
	if in.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if in.Date == "" {
		return errors.New("date is required")
	}
	return nil
}
