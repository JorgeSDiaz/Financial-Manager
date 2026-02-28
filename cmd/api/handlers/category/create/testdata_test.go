package create_test

import (
	"context"
	"time"

	appCreate "github.com/financial-manager/api/internal/application/category/create"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

const fixedTimestamp = "2026-02-23T10:00:00Z"

type fakeUseCase struct {
	out domaincategory.Category
	err error
}

func (f *fakeUseCase) Execute(_ context.Context, _ appCreate.Input) (domaincategory.Category, error) {
	return f.out, f.err
}

func buildDomainCategory(id, name string) domaincategory.Category {
	t, _ := time.Parse("2006-01-02T15:04:05Z", fixedTimestamp)
	return domaincategory.Category{
		ID:        id,
		Name:      name,
		Type:      domaincategory.TypeExpense,
		Color:     "#FF5733",
		Icon:      "restaurant",
		IsSystem:  false,
		IsActive:  true,
		CreatedAt: t,
		UpdatedAt: t,
	}
}
