package list_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/account/list"
	"github.com/financial-manager/api/internal/application/account/list/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repo    *mocks.Repository
		wantErr error
		wantOut []domainaccount.Account
	}{
		{
			name:    "empty repository returns nil without error",
			repo:    buildMockRepo(nil, nil),
			wantOut: nil,
		},
		{
			name:    "returns accounts from repository",
			repo:    buildMockRepo([]domainaccount.Account{cashAccount, bankAccount}, nil),
			wantOut: []domainaccount.Account{cashAccount, bankAccount},
		},
		{
			name:    "repository error is propagated",
			repo:    buildMockRepo(nil, errors.New("db error")),
			wantErr: fmt.Errorf("list accounts: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := list.New(tc.repo)
			out, err := uc.Execute(context.Background())

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, out)
			tc.repo.AssertExpectations(t)
		})
	}
}
