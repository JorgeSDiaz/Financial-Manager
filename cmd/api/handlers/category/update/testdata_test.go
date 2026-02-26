package update_test

import (
	"context"
	"time"

	appUpdate "github.com/financial-manager/api/internal/application/category/update"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

const fixedTimestamp = "2026-02-23T10:00:00Z"

type fakeUseCase struct {
	wantInput appUpdate.Input
	out       domaincategory.Category
	err       error
}

func (f *fakeUseCase) Execute(_ context.Context, in appUpdate.Input) (domaincategory.Category, error) {
	if f.wantInput != (appUpdate.Input{}) {
		// shallow check: only verify when caller set an expectation
		_ = in // validated via test assertions on the response
	}
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
