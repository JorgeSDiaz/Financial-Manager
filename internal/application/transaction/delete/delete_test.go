package delete_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	transactiondelete "github.com/financial-manager/api/internal/application/transaction/delete"
	"github.com/financial-manager/api/internal/application/transaction/delete/mocks"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
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
			name: "existing transaction is soft-deleted",
			id:   "tx-1",
			repo: buildMockRepoDelete("tx-1", nil),
		},
		{
			name:    "missing ID returns validation error",
			id:      "",
			repo:    &mocks.Repository{},
			wantErr: errors.New("id is required"),
		},
		{
			name:    "nonexistent ID returns ErrNotFound",
			id:      "missing",
			repo:    buildMockRepoGetByID("missing", domaintransaction.Transaction{}, domainshared.ErrNotFound),
			wantErr: domainshared.ErrNotFound,
		},
		{
			name:    "repository error is wrapped and propagated",
			id:      "any",
			repo:    buildMockRepoDelete("any", errors.New("db error")),
			wantErr: fmt.Errorf("delete transaction: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clock := &mocks.Clock{}
			uc := transactiondelete.New(tc.repo, clock)
			err := uc.Execute(context.Background(), tc.id)

			assert.Equal(t, tc.wantErr, err)
			tc.repo.AssertExpectations(t)
		})
	}
}
