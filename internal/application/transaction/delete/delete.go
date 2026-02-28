// Package delete implements the delete transaction use case (soft delete).
package delete

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
	SoftDelete(ctx context.Context, id string) error
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

func (uc *UseCase) Execute(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return domainshared.ErrNotFound
	}

	if err := uc.repo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}

	return nil
}
