package globalbalance_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/account/globalbalance"
	"github.com/financial-manager/api/internal/application/account/globalbalance/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repo    *mocks.Repository
		wantErr error
		wantOut float64
	}{
		{
			name:    "empty repository returns zero balance",
			repo:    buildMockRepo(nil, nil),
			wantOut: 0,
		},
		{
			name:    "sums current balance of accounts returned by repository",
			repo:    buildMockRepo([]domainaccount.Account{cashAccountWith100, bankAccountWith250}, nil),
			wantOut: 350.50,
		},
		{
			name:    "repository error is propagated",
			repo:    buildMockRepo(nil, errors.New("db error")),
			wantErr: fmt.Errorf("get global balance: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := globalbalance.New(tc.repo)
			out, err := uc.Execute(context.Background())

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, out)
			tc.repo.AssertExpectations(t)
		})
	}
}
