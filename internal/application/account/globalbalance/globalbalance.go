// Package globalbalance implements the get global balance use case.
package globalbalance

import (
	"context"
	"fmt"
)

// UseCase implements the get global balance use case (US-AC-006).
type UseCase struct {
	repo Repository
}

// New creates a new UseCase.
func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute sums the CurrentBalance of all active accounts.
func (uc *UseCase) Execute(ctx context.Context) (float64, error) {
	accounts, err := uc.repo.List(ctx)
	if err != nil {
		return 0, fmt.Errorf("get global balance: %w", err)
	}

	var total float64
	for _, acc := range accounts {
		total += acc.CurrentBalance
	}

	return total, nil
}
