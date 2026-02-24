// Package delete implements the delete account use case.
package delete

import (
	"context"
	"errors"
	"fmt"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// UseCase implements the delete account use case (US-AC-005).
type UseCase struct {
	repo Repository
}

// New creates a new UseCase.
func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute soft-deletes the account if it has no associated transactions.
func (uc *UseCase) Execute(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("account ID is required")
	}

	// TODO(M4): HasTransactions always returns false in M2; wire real check in M4.
	hasTx, err := uc.repo.HasTransactions(ctx, id)
	if err != nil {
		return fmt.Errorf("delete account: %w", err)
	}
	if hasTx {
		return domainaccount.ErrAccountHasTransactions
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete account: %w", err)
	}

	return nil
}
