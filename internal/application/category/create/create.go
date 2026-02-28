// Package create implements the create category use case.
package create

import (
	"context"
	"errors"
	"fmt"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// Input carries the data required to create a new category.
type Input struct {
	Name  string
	Type  string
	Color string
	Icon  string
}

// UseCase implements the create category use case (US-CAT-002, US-CAT-003).
type UseCase struct {
	repo  Repository
	idGen IDGenerator
	clock Clock
}

// New creates a new UseCase.
func New(repo Repository, idGen IDGenerator, clock Clock) *UseCase {
	return &UseCase{repo: repo, idGen: idGen, clock: clock}
}

// Execute validates input, creates a new Category, and persists it.
func (uc *UseCase) Execute(ctx context.Context, in Input) (domaincategory.Category, error) {
	if err := validateInput(in); err != nil {
		return domaincategory.Category{}, err
	}

	now := uc.clock.Now().UTC()
	cat := domaincategory.Category{
		ID:        uc.idGen.NewID(),
		Name:      in.Name,
		Type:      domaincategory.Type(in.Type),
		Color:     in.Color,
		Icon:      in.Icon,
		IsSystem:  false,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.repo.Create(ctx, cat); err != nil {
		return domaincategory.Category{}, fmt.Errorf("create category: %w", err)
	}

	return cat, nil
}

var validCategoryTypes = map[domaincategory.Type]struct{}{
	domaincategory.TypeExpense: {},
	domaincategory.TypeIncome:  {},
}

func validateInput(in Input) error {
	if in.Name == "" {
		return errors.New("category name is required")
	}
	if in.Color == "" {
		return errors.New("category color is required")
	}
	if in.Icon == "" {
		return errors.New("category icon is required")
	}
	if _, ok := validCategoryTypes[domaincategory.Type(in.Type)]; !ok {
		return fmt.Errorf("invalid category type %q: must be expense or income", in.Type)
	}
	return nil
}
