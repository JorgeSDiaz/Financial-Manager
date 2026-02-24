// Package list implements the list accounts use case.
package list

import (
	"context"
	"fmt"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// UseCase implements the list accounts use case (US-AC-003).
type UseCase struct {
	repo Repository
}

// New creates a new UseCase.
func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute returns all active accounts.
func (uc *UseCase) Execute(ctx context.Context) ([]domainaccount.Account, error) {
	accounts, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	return accounts, nil
}
