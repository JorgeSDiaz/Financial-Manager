// Package delete implements the delete category use case.
package delete

import (
	"context"
	"errors"
	"fmt"

	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

// UseCase implements the delete category use case (US-CAT-005).
type UseCase struct {
	repo Repository
}

// New creates a new UseCase.
func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute deletes a category if it's not a system category and has no transactions.
func (uc *UseCase) Execute(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("category id is required")
	}

	cat, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainshared.ErrNotFound) {
			return fmt.Errorf("category not found: %w", err)
		}
		return fmt.Errorf("get category: %w", err)
	}

	if cat.IsSystem {
		return errors.New("cannot delete system category")
	}

	hasTransactions, err := uc.repo.HasTransactions(ctx, id)
	if err != nil {
		return fmt.Errorf("check transactions: %w", err)
	}
	if hasTransactions {
		return errors.New("cannot delete category with associated transactions")
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}
