package update_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/account/update"
	"github.com/financial-manager/api/internal/application/account/update/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
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
		wantOut domainaccount.Account
	}{
		{
			name: "valid update returns updated account",
			repo: buildMockRepoFull("acc-1", seeded, domainaccount.Account{
				ID: "acc-1", Name: "New Name",
				Type: domainaccount.AccountTypeCash, InitialBalance: 500.0, CurrentBalance: 500.0,
				Currency: "USD", Color: "#FFFFFF", Icon: "wallet", IsActive: true, UpdatedAt: updatedAt,
			}, nil),
			clock: buildMockClock(),
			input: update.Input{ID: "acc-1", Name: "New Name"},
			wantOut: domainaccount.Account{
				ID: "acc-1", Name: "New Name",
				Type: domainaccount.AccountTypeCash, InitialBalance: 500.0, CurrentBalance: 500.0,
				Currency: "USD", Color: "#FFFFFF", Icon: "wallet", IsActive: true, UpdatedAt: updatedAt,
			},
		},
		{
			name: "valid update with all optional fields returns fully updated account",
			repo: buildMockRepoFull("acc-1", seeded, domainaccount.Account{
				ID: "acc-1", Name: "New Name",
				Type: domainaccount.AccountTypeCash, InitialBalance: 500.0, CurrentBalance: 500.0,
				Currency: "USD", Color: "#000000", Icon: "bank", IsActive: true, UpdatedAt: updatedAt,
			}, nil),
			clock: buildMockClock(),
			input: update.Input{ID: "acc-1", Name: "New Name", Color: "#000000", Icon: "bank"},
			wantOut: domainaccount.Account{
				ID: "acc-1", Name: "New Name",
				Type: domainaccount.AccountTypeCash, InitialBalance: 500.0, CurrentBalance: 500.0,
				Currency: "USD", Color: "#000000", Icon: "bank", IsActive: true, UpdatedAt: updatedAt,
			},
		},
		{
			name:    "missing ID returns validation error",
			repo:    &mocks.Repository{},
			clock:   &mocks.Clock{},
			input:   update.Input{Name: "Name"},
			wantErr: errors.New("account ID is required"),
		},
		{
			name:    "missing name returns validation error",
			repo:    &mocks.Repository{},
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "acc-1"},
			wantErr: errors.New("account name is required"),
		},
		{
			name:    "nonexistent ID returns wrapped ErrNotFound",
			repo:    buildMockRepoGetByID("missing", domainaccount.Account{}, domainshared.ErrNotFound),
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "missing", Name: "Name"},
			wantErr: fmt.Errorf("update account: %w", domainshared.ErrNotFound),
		},
		{
			name:    "GetByID error is wrapped and propagated",
			repo:    buildMockRepoGetByID("any", domainaccount.Account{}, errors.New("db error")),
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "any", Name: "Name"},
			wantErr: fmt.Errorf("update account: %w", errors.New("db error")),
		},
		{
			name: "Update error is wrapped and propagated",
			repo: buildMockRepoFull("acc-2", buildActiveAccount("acc-2", "Existing"), domainaccount.Account{
				ID: "acc-2", Name: "New Name",
				Type: domainaccount.AccountTypeCash, InitialBalance: 500.0, CurrentBalance: 500.0,
				Currency: "USD", Color: "#FFFFFF", Icon: "wallet", IsActive: true, UpdatedAt: updatedAt,
			}, errors.New("db write error")),
			clock:   buildMockClock(),
			input:   update.Input{ID: "acc-2", Name: "New Name"},
			wantErr: fmt.Errorf("update account: %w", errors.New("db write error")),
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
