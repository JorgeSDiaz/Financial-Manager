package update_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/transaction/update"
	"github.com/financial-manager/api/internal/application/transaction/update/mocks"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	updatedAt := fixedTime()
	newDate, _ := time.Parse("2006-01-02", "2026-03-01")

	tests := []struct {
		name    string
		repo    *mocks.Repository
		clock   *mocks.Clock
		input   update.Input
		wantErr error
		wantOut domaintransaction.Transaction
	}{
		{
			name: "valid update returns updated transaction",
			repo: buildMockRepoFull("tx-1", seeded, domaintransaction.Transaction{
				ID: "tx-1", AccountID: "acc-001", CategoryID: "cat-001",
				Type: domaintransaction.TransactionTypeIncome, Amount: 100.0,
				Description: "New Description", Date: seeded.Date, IsActive: true,
				UpdatedAt: updatedAt,
			}, nil),
			clock: buildMockClock(),
			input: update.Input{ID: "tx-1", Description: "New Description"},
			wantOut: domaintransaction.Transaction{
				ID: "tx-1", AccountID: "acc-001", CategoryID: "cat-001",
				Type: domaintransaction.TransactionTypeIncome, Amount: 100.0,
				Description: "New Description", Date: seeded.Date, IsActive: true,
				UpdatedAt: updatedAt,
			},
		},
		{
			name: "valid update with all fields returns fully updated transaction",
			repo: buildMockRepoFull("tx-1", seeded, domaintransaction.Transaction{
				ID: "tx-1", AccountID: "acc-001", CategoryID: "cat-002",
				Type: domaintransaction.TransactionTypeIncome, Amount: 200.0,
				Description: "Updated Description", Date: newDate, IsActive: true,
				UpdatedAt: updatedAt,
			}, nil),
			clock: buildMockClock(),
			input: update.Input{ID: "tx-1", CategoryID: "cat-002", Amount: 200.0, Description: "Updated Description", Date: "2026-03-01"},
			wantOut: domaintransaction.Transaction{
				ID: "tx-1", AccountID: "acc-001", CategoryID: "cat-002",
				Type: domaintransaction.TransactionTypeIncome, Amount: 200.0,
				Description: "Updated Description", Date: newDate, IsActive: true,
				UpdatedAt: updatedAt,
			},
		},
		{
			name:    "missing ID returns validation error",
			repo:    &mocks.Repository{},
			clock:   &mocks.Clock{},
			input:   update.Input{Description: "Description"},
			wantErr: errors.New("id is required"),
		},
		{
			name:    "nonexistent ID returns ErrNotFound",
			repo:    buildMockRepoGetByID("missing", domaintransaction.Transaction{}, domainshared.ErrNotFound),
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "missing", Description: "Description"},
			wantErr: domainshared.ErrNotFound,
		},
		{
			name:    "GetByID error is wrapped and propagated",
			repo:    buildMockRepoGetByID("any", domaintransaction.Transaction{}, errors.New("db error")),
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "any", Description: "Description"},
			wantErr: domainshared.ErrNotFound,
		},
		{
			name: "Update error is wrapped and propagated",
			repo: buildMockRepoFull("tx-2", buildTransaction("tx-2", "acc-001", "cat-001", "Existing", 100.0), domaintransaction.Transaction{
				ID: "tx-2", AccountID: "acc-001", CategoryID: "cat-001",
				Type: domaintransaction.TransactionTypeIncome, Amount: 100.0,
				Description: "New Description", Date: buildTransaction("tx-2", "acc-001", "cat-001", "Existing", 100.0).Date,
				IsActive: true, UpdatedAt: updatedAt,
			}, errors.New("db write error")),
			clock:   buildMockClock(),
			input:   update.Input{ID: "tx-2", Description: "New Description"},
			wantErr: fmt.Errorf("update transaction: %w", errors.New("db write error")),
		},
		{
			name:    "invalid date format returns validation error",
			repo:    buildMockRepoGetByID("tx-1", seeded, nil),
			clock:   &mocks.Clock{},
			input:   update.Input{ID: "tx-1", Date: "invalid-date"},
			wantErr: errors.New("invalid date format, use YYYY-MM-DD"),
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
