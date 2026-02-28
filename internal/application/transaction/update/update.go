// Package update implements the update transaction use case.
package update

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainshared "github.com/financial-manager/api/internal/domain/shared"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (domaintransaction.Transaction, error)
	Update(ctx context.Context, t domaintransaction.Transaction) error
}

type Clock interface {
	Now() time.Time
}

type UseCase struct {
	repo  Repository
	clock Clock
}

func New(repo Repository, clock Clock) *UseCase {
	return &UseCase{repo: repo, clock: clock}
}

type Input struct {
	ID          string  `json:"id"`
	CategoryID  string  `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func (uc *UseCase) Execute(ctx context.Context, in Input) (domaintransaction.Transaction, error) {
	if err := validateInput(in); err != nil {
		return domaintransaction.Transaction{}, err
	}

	tx, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return domaintransaction.Transaction{}, domainshared.ErrNotFound
	}

	if in.Amount > 0 {
		tx.Amount = in.Amount
	}
	if in.Description != "" {
		tx.Description = in.Description
	}
	if in.CategoryID != "" {
		tx.CategoryID = in.CategoryID
	}
	if in.Date != "" {
		date, err := time.Parse("2006-01-02", in.Date)
		if err != nil {
			return domaintransaction.Transaction{}, errors.New("invalid date format, use YYYY-MM-DD")
		}
		tx.Date = date
	}

	tx.UpdatedAt = uc.clock.Now().UTC()

	if err := uc.repo.Update(ctx, tx); err != nil {
		return domaintransaction.Transaction{}, fmt.Errorf("update transaction: %w", err)
	}

	return tx, nil
}

func validateInput(in Input) error {
	if in.ID == "" {
		return errors.New("id is required")
	}
	return nil
}
