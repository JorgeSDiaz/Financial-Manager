package list_test

import (
	"context"
	"time"

	appList "github.com/financial-manager/api/internal/application/category/list"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

const fixedTimestamp = "2026-02-23T10:00:00Z"

type fakeUseCase struct {
	out []domaincategory.Category
	err error
}

func (f *fakeUseCase) Execute(_ context.Context, _ appList.Input) ([]domaincategory.Category, error) {
	return f.out, f.err
}

func buildDomainCategories() []domaincategory.Category {
	t, _ := time.Parse("2006-01-02T15:04:05Z", fixedTimestamp)
	return []domaincategory.Category{
		{
			ID:        "cat-1",
			Name:      "Food",
			Type:      domaincategory.TypeExpense,
			Color:     "#FF5733",
			Icon:      "restaurant",
			IsActive:  true,
			CreatedAt: t,
			UpdatedAt: t,
		},
		{
			ID:        "cat-2",
			Name:      "Salary",
			Type:      domaincategory.TypeIncome,
			Color:     "#00FF00",
			Icon:      "money",
			IsActive:  true,
			CreatedAt: t,
			UpdatedAt: t,
		},
	}
}

func buildFailingUseCase(err error) *fakeUseCase {
	return &fakeUseCase{err: err}
}
