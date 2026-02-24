package delete_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	accountdelete "github.com/financial-manager/api/internal/application/account/delete"
	"github.com/financial-manager/api/internal/application/account/delete/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		repo    *mocks.Repository
		wantErr error
	}{
		{
			name: "existing account is soft-deleted",
			id:   "acc-1",
			repo: buildMockRepoDelete("acc-1", nil),
		},
		{
			name:    "missing ID returns validation error",
			id:      "",
			repo:    &mocks.Repository{},
			wantErr: errors.New("account ID is required"),
		},
		{
			name:    "nonexistent ID returns wrapped ErrNotFound",
			id:      "missing",
			repo:    buildMockRepoDelete("missing", domainshared.ErrNotFound),
			wantErr: fmt.Errorf("delete account: %w", domainshared.ErrNotFound),
		},
		{
			name:    "account with transactions returns ErrAccountHasTransactions",
			id:      "acc-2",
			repo:    buildMockRepoHasTx("acc-2", true, nil),
			wantErr: domainaccount.ErrAccountHasTransactions,
		},
		{
			name:    "repository error is wrapped and propagated",
			id:      "any",
			repo:    buildMockRepoHasTx("any", false, errors.New("db error")),
			wantErr: fmt.Errorf("delete account: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := accountdelete.New(tc.repo)
			err := uc.Execute(context.Background(), tc.id)

			assert.Equal(t, tc.wantErr, err)
			tc.repo.AssertExpectations(t)
		})
	}
}
