package summary_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/transaction/summary"
	"github.com/financial-manager/api/internal/application/transaction/summary/mocks"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		incomeRepo  *mocks.Repository
		expenseRepo *mocks.Repository
		input       summary.Input
		wantErr     error
		wantOut     summary.Summary
	}{
		{
			name:        "empty transactions returns zero summary",
			incomeRepo:  buildMockRepoIncome(nil, nil),
			expenseRepo: buildMockRepoExpense(nil, nil),
			input:       summary.Input{},
			wantOut:     summary.Summary{TotalIncome: 0, TotalExpense: 0, Balance: 0},
		},
		{
			name:        "calculates summary correctly",
			incomeRepo:  buildMockRepoIncome([]domaintransaction.Transaction{income100, income500}, nil),
			expenseRepo: buildMockRepoExpense([]domaintransaction.Transaction{expense50, expense200}, nil),
			input:       summary.Input{},
			wantOut:     summary.Summary{TotalIncome: 600.0, TotalExpense: 250.0, Balance: 350.0},
		},
		{
			name:        "repository error for incomes is propagated",
			incomeRepo:  buildMockRepoIncome(nil, errors.New("db error")),
			expenseRepo: &mocks.Repository{},
			input:       summary.Input{},
			wantErr:     fmt.Errorf("get summary: %w", errors.New("db error")),
		},
		{
			name:        "repository error for expenses is propagated",
			incomeRepo:  buildMockRepoIncome(nil, nil),
			expenseRepo: buildMockRepoExpense(nil, errors.New("db error")),
			input:       summary.Input{},
			wantErr:     fmt.Errorf("get summary: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Since we need to call ListByType twice with different types,
			// we merge the mocks
			repo := &mocks.Repository{}
			if tc.incomeRepo != nil {
				repo = tc.incomeRepo
			}
			if tc.expenseRepo != nil {
				// Copy expectations from expenseRepo
				for _, call := range tc.expenseRepo.ExpectedCalls {
					repo.On(call.Method, call.Arguments...).Return(call.ReturnArguments...)
				}
			}

			uc := summary.New(repo)
			out, err := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, out)
			repo.AssertExpectations(t)
		})
	}
}
