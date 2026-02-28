// Package update implements the update category use case.
package update

import (
	"context"
	"errors"
	"fmt"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

// Input carries the data required to update a category.
type Input struct {
	ID    string
	Name  string
	Color string
	Icon  string
}

// UseCase implements the update category use case (US-CAT-004).
type UseCase struct {
	repo  Repository
	clock Clock
}

// New creates a new UseCase.
func New(repo Repository, clock Clock) *UseCase {
	return &UseCase{repo: repo, clock: clock}
}

// Execute validates input, updates the Category, and persists it.
func (uc *UseCase) Execute(ctx context.Context, in Input) (domaincategory.Category, error) {
	if err := validateInput(in); err != nil {
		return domaincategory.Category{}, err
	}

	cat, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, domainshared.ErrNotFound) {
			return domaincategory.Category{}, fmt.Errorf("category not found: %w", err)
		}
		return domaincategory.Category{}, fmt.Errorf("get category: %w", err)
	}

	if cat.IsSystem {
		return domaincategory.Category{}, errors.New("cannot update system category")
	}

	cat.Name = in.Name
	cat.Color = in.Color
	cat.Icon = in.Icon
	cat.UpdatedAt = uc.clock.Now().UTC()

	if err := uc.repo.Update(ctx, cat); err != nil {
		return domaincategory.Category{}, fmt.Errorf("update category: %w", err)
	}

	return cat, nil
}

func validateInput(in Input) error {
	if in.ID == "" {
		return errors.New("category id is required")
	}
	if in.Name == "" {
		return errors.New("category name is required")
	}
	if in.Color == "" {
		return errors.New("category color is required")
	}
	if in.Icon == "" {
		return errors.New("category icon is required")
	}
	return nil
}
