package list_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/transaction/income/list"
	"github.com/financial-manager/api/internal/application/transaction/income/list/mocks"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repo    *mocks.Repository
		input   list.Input
		wantErr error
		wantOut []domaintransaction.Transaction
	}{
		{
			name:    "empty repository returns nil without error",
			repo:    buildMockRepo(nil, nil),
			input:   list.Input{},
			wantOut: nil,
		},
		{
			name:    "returns income transactions from repository",
			repo:    buildMockRepo([]domaintransaction.Transaction{income1, income2}, nil),
			input:   list.Input{},
			wantOut: []domaintransaction.Transaction{income1, income2},
		},
		{
			name:    "returns filtered income transactions",
			repo:    buildMockRepo([]domaintransaction.Transaction{income1}, nil),
			input:   list.Input{AccountID: "acc-001"},
			wantOut: []domaintransaction.Transaction{income1},
		},
		{
			name:    "repository error is propagated",
			repo:    buildMockRepo(nil, errors.New("db error")),
			input:   list.Input{},
			wantErr: fmt.Errorf("list incomes: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := list.New(tc.repo)
			out, err := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, out)
			tc.repo.AssertExpectations(t)
		})
	}
}
