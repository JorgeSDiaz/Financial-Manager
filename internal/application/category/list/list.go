// Package list implements the list categories use case.
package list

import (
	"context"
	"fmt"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// Input carries optional filters for listing categories.
type Input struct {
	Type *string // nil = all types, "expense" or "income"
}

// UseCase implements the list categories use case (US-CAT-001, US-CAT-006).
type UseCase struct {
	repo Repository
}

// New creates a new UseCase.
func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute retrieves categories filtered by type (if specified).
func (uc *UseCase) Execute(ctx context.Context, in Input) ([]domaincategory.Category, error) {
	var catType *domaincategory.Type
	if in.Type != nil {
		t := domaincategory.Type(*in.Type)
		if t != domaincategory.TypeExpense && t != domaincategory.TypeIncome {
			return nil, fmt.Errorf("invalid category type %q: must be expense or income", *in.Type)
		}
		catType = &t
	}

	categories, err := uc.repo.List(ctx, catType)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	return categories, nil
}
