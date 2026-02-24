// Package update implements the update account use case.
package update

import (
	"context"
	"errors"
	"fmt"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// Input carries the data required to update an account.
type Input struct {
	ID    string
	Name  string
	Color string
	Icon  string
}

// UseCase implements the update account use case (US-AC-004).
type UseCase struct {
	repo  Repository
	clock Clock
}

// New creates a new UseCase.
func New(repo Repository, clock Clock) *UseCase {
	return &UseCase{repo: repo, clock: clock}
}

// Execute validates input, fetches the account, applies changes, and persists it.
func (uc *UseCase) Execute(ctx context.Context, in Input) (domainaccount.Account, error) {
	if err := validateInput(in); err != nil {
		return domainaccount.Account{}, err
	}

	acc, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return domainaccount.Account{}, fmt.Errorf("update account: %w", err)
	}

	if in.Name != "" {
		acc.Name = in.Name
	}
	if in.Color != "" {
		acc.Color = in.Color
	}
	if in.Icon != "" {
		acc.Icon = in.Icon
	}
	acc.UpdatedAt = uc.clock.Now().UTC()

	if err := uc.repo.Update(ctx, acc); err != nil {
		return domainaccount.Account{}, fmt.Errorf("update account: %w", err)
	}

	return acc, nil
}

func validateInput(in Input) error {
	if in.ID == "" {
		return errors.New("account ID is required")
	}
	if in.Name == "" {
		return errors.New("account name is required")
	}
	return nil
}
