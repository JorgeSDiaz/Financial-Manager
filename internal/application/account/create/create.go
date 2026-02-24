// Package create implements the create account use case.
package create

import (
	"context"
	"errors"
	"fmt"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// Input carries the data required to create a new account.
type Input struct {
	Name           string
	Type           string
	InitialBalance float64
	Currency       string
	Color          string
	Icon           string
}

// UseCase implements the create account use case (US-AC-001).
type UseCase struct {
	repo  Repository
	idGen IDGenerator
	clock Clock
}

// New creates a new UseCase.
func New(repo Repository, idGen IDGenerator, clock Clock) *UseCase {
	return &UseCase{repo: repo, idGen: idGen, clock: clock}
}

// Execute validates input, creates a new Account, and persists it.
func (uc *UseCase) Execute(ctx context.Context, in Input) (domainaccount.Account, error) {
	if err := validateInput(in); err != nil {
		return domainaccount.Account{}, err
	}

	now := uc.clock.Now().UTC()
	acc := domainaccount.Account{
		ID:             uc.idGen.NewID(),
		Name:           in.Name,
		Type:           domainaccount.AccountType(in.Type),
		InitialBalance: in.InitialBalance,
		CurrentBalance: in.InitialBalance,
		Currency:       in.Currency,
		Color:          in.Color,
		Icon:           in.Icon,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := uc.repo.Create(ctx, acc); err != nil {
		return domainaccount.Account{}, fmt.Errorf("create account: %w", err)
	}

	return acc, nil
}

var validAccountTypes = map[domainaccount.AccountType]struct{}{
	domainaccount.AccountTypeCash:       {},
	domainaccount.AccountTypeBank:       {},
	domainaccount.AccountTypeCreditCard: {},
	domainaccount.AccountTypeSavings:    {},
}

func validateInput(in Input) error {
	if in.Name == "" {
		return errors.New("account name is required")
	}
	if _, ok := validAccountTypes[domainaccount.AccountType(in.Type)]; !ok {
		return fmt.Errorf("invalid account type %q: must be cash, bank, credit_card, or savings", in.Type)
	}
	if in.InitialBalance < 0 {
		return errors.New("initial balance must be zero or positive")
	}
	return nil
}
