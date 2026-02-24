package get_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/account/get"
	"github.com/financial-manager/api/internal/application/account/get/mocks"
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
		wantOut domainaccount.Account
	}{
		{
			name:    "existing ID returns correct account",
			id:      "acc-1",
			repo:    buildMockRepo("acc-1", activeAccount, nil),
			wantOut: activeAccount,
		},
		{
			name:    "nonexistent ID returns ErrNotFound",
			id:      "missing",
			repo:    buildMockRepo("missing", domainaccount.Account{}, domainshared.ErrNotFound),
			wantErr: fmt.Errorf("get account: %w", domainshared.ErrNotFound),
		},
		{
			name:    "repository error is propagated",
			id:      "any",
			repo:    buildMockRepo("any", domainaccount.Account{}, errors.New("db error")),
			wantErr: fmt.Errorf("get account: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := get.New(tc.repo)
			out, err := uc.Execute(context.Background(), tc.id)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, out)
			tc.repo.AssertExpectations(t)
		})
	}
}
