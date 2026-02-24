// Package get implements the get account by ID use case.
package get

import (
	"context"
	"fmt"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// UseCase implements the get account by ID use case (US-AC-002).
type UseCase struct {
	repo Repository
}

// New creates a new UseCase.
func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute retrieves an account by its ID.
func (uc *UseCase) Execute(ctx context.Context, id string) (domainaccount.Account, error) {
	acc, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return domainaccount.Account{}, fmt.Errorf("get account: %w", err)
	}
	return acc, nil
}
