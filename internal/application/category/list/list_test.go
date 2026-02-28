package list_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/category/list"
	"github.com/financial-manager/api/internal/application/category/list/mocks"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	expenseType := "expense"
	incomeType := "income"
	invalidType := "invalid"

	categories := buildCategories()

	tests := []struct {
		name    string
		repo    *mocks.Repository
		input   list.Input
		wantErr error
		wantOut []domaincategory.Category
	}{
		{
			name:    "list all categories returns all active categories",
			repo:    buildMockRepo(nil, categories, nil),
			input:   list.Input{},
			wantOut: categories,
		},
		{
			name: "list expense categories returns only expense categories",
			repo: buildMockRepo(func() *domaincategory.Type {
				t := domaincategory.TypeExpense
				return &t
			}(), []domaincategory.Category{categories[0]}, nil),
			input:   list.Input{Type: &expenseType},
			wantOut: []domaincategory.Category{categories[0]},
		},
		{
			name: "list income categories returns only income categories",
			repo: buildMockRepo(func() *domaincategory.Type {
				t := domaincategory.TypeIncome
				return &t
			}(), []domaincategory.Category{categories[1]}, nil),
			input:   list.Input{Type: &incomeType},
			wantOut: []domaincategory.Category{categories[1]},
		},
		{
			name:    "invalid category type returns error",
			repo:    &mocks.Repository{},
			input:   list.Input{Type: &invalidType},
			wantErr: fmt.Errorf(`invalid category type "invalid": must be expense or income`),
		},
		{
			name:    "repository error is wrapped and propagated",
			repo:    buildMockRepo(nil, nil, errors.New("db error")),
			input:   list.Input{},
			wantErr: fmt.Errorf("list categories: %w", errors.New("db error")),
		},
		{
			name:    "empty list returns empty slice",
			repo:    buildMockRepo(nil, []domaincategory.Category{}, nil),
			input:   list.Input{},
			wantOut: []domaincategory.Category{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := list.New(tc.repo)
			out, err := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, out)
			tc.repo.AssertExpectations(t)
		})
	}
}
