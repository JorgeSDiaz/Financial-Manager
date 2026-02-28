package create_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/transaction/expense/create"
	"github.com/financial-manager/api/internal/application/transaction/expense/create/mocks"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
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
		wantOut domaintransaction.Transaction
	}{
		{
			name: "valid input creates expense transaction",
			input: create.Input{
				AccountID:   "acc-001",
				CategoryID:  "cat-001",
				Amount:      100.0,
				Description: "Groceries",
				Date:        fixedDate,
			},
			repo:    buildMockRepo(validExpense, nil),
			idGen:   buildMockIDGenerator(),
			clock:   buildMockClock(),
			wantOut: validExpense,
		},
		{
			name:    "empty account_id returns validation error",
			input:   create.Input{Amount: 100.0, Date: fixedDate},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("account_id is required"),
		},
		{
			name:    "zero amount returns validation error",
			input:   create.Input{AccountID: "acc-001", Amount: 0, Date: fixedDate},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("amount must be positive"),
		},
		{
			name:    "negative amount returns validation error",
			input:   create.Input{AccountID: "acc-001", Amount: -100, Date: fixedDate},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("amount must be positive"),
		},
		{
			name:    "empty date returns validation error",
			input:   create.Input{AccountID: "acc-001", Amount: 100.0},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("date is required"),
		},
		{
			name:    "invalid date format returns validation error",
			input:   create.Input{AccountID: "acc-001", Amount: 100.0, Date: "invalid-date"},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("invalid date format, use YYYY-MM-DD"),
		},
		{
			name:    "repository error is wrapped and propagated",
			input:   create.Input{AccountID: "acc-001", Amount: 100.0, Date: fixedDate},
			repo:    buildMockRepo(errorExpense, errors.New("db unavailable")),
			idGen:   buildMockIDGenerator(),
			clock:   buildMockClock(),
			wantErr: fmt.Errorf("create expense: %w", errors.New("db unavailable")),
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
