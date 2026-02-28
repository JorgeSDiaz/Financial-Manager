package delete_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/category/delete"
	"github.com/financial-manager/api/internal/application/category/delete/mocks"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repo    *mocks.Repository
		id      string
		wantErr error
	}{
		{
			name: "valid delete of custom category succeeds",
			repo: buildMockRepoFull("cat-1", buildActiveCategory("cat-1", "Test"), false, nil),
			id:   "cat-1",
		},
		{
			name:    "missing ID returns validation error",
			repo:    &mocks.Repository{},
			id:      "",
			wantErr: errors.New("category id is required"),
		},
		{
			name:    "nonexistent ID returns wrapped ErrNotFound",
			repo:    buildMockRepoWithGet("missing", domaincategory.Category{}, domainshared.ErrNotFound),
			id:      "missing",
			wantErr: fmt.Errorf("category not found: %w", domainshared.ErrNotFound),
		},
		{
			name: "system category returns error",
			repo: buildMockRepoWithGet("cat-sys", domaincategory.Category{
				ID:       "cat-sys",
				Name:     "System",
				IsSystem: true,
			}, nil),
			id:      "cat-sys",
			wantErr: errors.New("cannot delete system category"),
		},
		{
			name:    "GetByID error is wrapped and propagated",
			repo:    buildMockRepoWithGet("any", domaincategory.Category{}, errors.New("db error")),
			id:      "any",
			wantErr: fmt.Errorf("get category: %w", errors.New("db error")),
		},
		{
			name: "category with transactions returns error",
			repo: func() *mocks.Repository {
				m := &mocks.Repository{}
				m.On("GetByID", mock.Anything, "cat-trans").Return(buildActiveCategory("cat-trans", "Test"), nil).Once()
				m.On("HasTransactions", mock.Anything, "cat-trans").Return(true, nil).Once()
				return m
			}(),
			id:      "cat-trans",
			wantErr: errors.New("cannot delete category with associated transactions"),
		},
		{
			name: "HasTransactions error is wrapped and propagated",
			repo: func() *mocks.Repository {
				m := &mocks.Repository{}
				m.On("GetByID", mock.Anything, "cat-err").Return(buildActiveCategory("cat-err", "Test"), nil).Once()
				m.On("HasTransactions", mock.Anything, "cat-err").Return(false, errors.New("db error")).Once()
				return m
			}(),
			id:      "cat-err",
			wantErr: fmt.Errorf("check transactions: %w", errors.New("db error")),
		},
		{
			name:    "Delete error is wrapped and propagated",
			repo:    buildMockRepoFull("cat-2", buildActiveCategory("cat-2", "Test"), false, errors.New("db error")),
			id:      "cat-2",
			wantErr: fmt.Errorf("delete category: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := delete.New(tc.repo)
			err := uc.Execute(context.Background(), tc.id)

			assert.Equal(t, tc.wantErr, err)
			tc.repo.AssertExpectations(t)
		})
	}
}
