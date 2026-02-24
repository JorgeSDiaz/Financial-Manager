package create_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/account/create"
	"github.com/financial-manager/api/internal/application/account/create/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   create.Input
		repo    *mocks.Repository
		idGen   *mocks.IDGenerator
		clock   *mocks.Clock
		wantErr error
		wantOut domainaccount.Account
	}{
		{
			name: "valid input creates account with fixed ID and InitialBalance equals CurrentBalance",
			input: create.Input{
				Name:           "Efectivo",
				Type:           "cash",
				InitialBalance: 1000.0,
				Currency:       "USD",
				Color:          "#00FF00",
				Icon:           "wallet",
			},
			repo:    buildMockRepo(validAccount, nil),
			idGen:   buildMockIDGenerator(),
			clock:   buildMockClock(),
			wantOut: validAccount,
		},
		{
			name:    "empty name returns validation error",
			input:   create.Input{Type: "cash", InitialBalance: 0},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("account name is required"),
		},
		{
			name:    "invalid type returns validation error",
			input:   create.Input{Name: "X", Type: "invalid"},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: fmt.Errorf(`invalid account type "invalid": must be cash, bank, credit_card, or savings`),
		},
		{
			name:    "negative initial balance returns validation error",
			input:   create.Input{Name: "X", Type: "cash", InitialBalance: -1},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("initial balance must be zero or positive"),
		},
		{
			name:    "repository error is wrapped and propagated",
			input:   create.Input{Name: "X", Type: "cash"},
			repo:    buildMockRepo(errorAccount, errors.New("db unavailable")),
			idGen:   buildMockIDGenerator(),
			clock:   buildMockClock(),
			wantErr: fmt.Errorf("create account: %w", errors.New("db unavailable")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := create.New(tc.repo, tc.idGen, tc.clock)
			out, err := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, out)
			tc.repo.AssertExpectations(t)
			tc.idGen.AssertExpectations(t)
			tc.clock.AssertExpectations(t)
		})
	}
}
