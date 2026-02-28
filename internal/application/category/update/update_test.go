package update_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/category/update"
	"github.com/financial-manager/api/internal/application/category/update/mocks"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	updatedAt := fixedTime()

	tests := []struct {
		name    string
		repo    *mocks.Repository
		clock   *mocks.Clock
		input   update.Input
		wantErr error
		wantOut domaincategory.Category
	}{
		{
			name: "valid update returns updated category",
			repo: buildMockRepoFull("cat-1", seeded, domaincategory.Category{
				ID: "cat-1", Name: "New Name",
				Type: domaincategory.TypeExpense, Color: "#FFFFFF", Icon: "wallet",
				IsSystem: false, IsActive: true, UpdatedAt: updatedAt,
			}, nil),
			clock: buildMockClock(),
			input: update.Input{ID: "cat-1", Name: "New Name", Color: "#FFFFFF", Icon: "wallet"},
			wantOut: domaincategory.Category{
				ID: "cat-1", Name: "New Name",
				Type: domaincategory.TypeExpense, Color: "#FFFFFF", Icon: "wallet",
				IsSystem: false, IsActive: true, UpdatedAt: updatedAt,
			},
		},
		{
			name: "valid update with all optional fields returns fully updated category",
			repo: buildMockRepoFull("cat-1", seeded, domaincategory.Category{
				ID: "cat-1", Name: "New Name",
				Type: domaincategory.TypeExpense, Color: "#000000", Icon: "bank",
				IsSystem: false, IsActive: true, UpdatedAt: updatedAt,
			}, nil),
			clock: buildMockClock(),
			input: update.Input{ID: "cat-1", Name: "New Name", Color: "#000000", Icon: "bank"},
			wantOut: domaincategory.Category{
				ID: "cat-1", Name: "New Name",
				Type: domaincategory.TypeExpense, Color: "#000000", Icon: "bank",
				IsSystem: false, IsActive: true, UpdatedAt: updatedAt,
			},
		},
		{
			name:    "missing ID returns validation error",
			repo:    &mocks.Repository{},
			clock:   &mocks.Clock{},
			input:   update.Input{Name: "Name", Color: "red", Icon: "icon"},
			wantErr: errors.New("category id is required"),
		},
		{
			name:    "missing name returns validation error",
			repo:    &mocks.Repository{},
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "cat-1", Color: "red", Icon: "icon"},
			wantErr: errors.New("category name is required"),
		},
		{
			name:    "missing color returns validation error",
			repo:    &mocks.Repository{},
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "cat-1", Name: "Name", Icon: "icon"},
			wantErr: errors.New("category color is required"),
		},
		{
			name:    "missing icon returns validation error",
			repo:    &mocks.Repository{},
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "cat-1", Name: "Name", Color: "red"},
			wantErr: errors.New("category icon is required"),
		},
		{
			name:    "nonexistent ID returns wrapped ErrNotFound",
			repo:    buildMockRepoGetByID("missing", domaincategory.Category{}, domainshared.ErrNotFound),
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "missing", Name: "Name", Color: "red", Icon: "icon"},
			wantErr: fmt.Errorf("category not found: %w", domainshared.ErrNotFound),
		},
		{
			name:    "GetByID error is wrapped and propagated",
			repo:    buildMockRepoGetByID("any", domaincategory.Category{}, errors.New("db error")),
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "any", Name: "Name", Color: "red", Icon: "icon"},
			wantErr: fmt.Errorf("get category: %w", errors.New("db error")),
		},
		{
			name: "system category returns error",
			repo: buildMockRepoGetByID("cat-sys", domaincategory.Category{
				ID:       "cat-sys",
				Name:     "System",
				IsSystem: true,
			}, nil),
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "cat-sys", Name: "Name", Color: "red", Icon: "icon"},
			wantErr: errors.New("cannot update system category"),
		},
		{
			name: "Update error is wrapped and propagated",
			repo: buildMockRepoFull("cat-2", buildActiveCategory("cat-2", "Existing"), domaincategory.Category{
				ID: "cat-2", Name: "New Name",
				Type: domaincategory.TypeExpense, Color: "#FFFFFF", Icon: "wallet",
				IsSystem: false, IsActive: true, UpdatedAt: updatedAt,
			}, errors.New("db write error")),
			clock:   buildMockClock(),
			input:   update.Input{ID: "cat-2", Name: "New Name", Color: "#FFFFFF", Icon: "wallet"},
			wantErr: fmt.Errorf("update category: %w", errors.New("db write error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := update.New(tc.repo, tc.clock)
			out, err := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, out)
			tc.repo.AssertExpectations(t)
			tc.clock.AssertExpectations(t)
		})
	}
}
