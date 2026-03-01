// Package pdfexport_test contains tests for the PDF export use case.
package pdfexport_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/pdfexport"
	"github.com/financial-manager/api/internal/application/pdfexport/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repo    *mocks.Repository
		input   pdfexport.Input
		wantErr error
	}{
		{
			name: "generates PDF report successfully",
			repo: buildMockRepo(
				[]domainaccount.Account{
					{ID: "acc-1", Name: "Banco", Type: domainaccount.AccountTypeBank, CurrentBalance: 1200.00},
				},
				[]domaincategory.Category{
					{ID: "cat-1", Name: "Alimentación"},
				},
				[]domaintransaction.Transaction{
					buildIncome("tx-1", 1000.00),
				},
				[]domaintransaction.Transaction{
					buildExpense("tx-2", 50.00, "cat-1"),
				},
				nil,
			),
			input: pdfexport.Input{
				Month:         "2026-02",
				IncludeCharts: false,
			},
			wantErr: nil,
		},
		{
			name:    "invalid month format returns error",
			repo:    &mocks.Repository{},
			input:   pdfexport.Input{Month: "invalid"},
			wantErr: fmt.Errorf("invalid month format: %w", errors.New("parsing time")),
		},
		{
			name: "empty month generates report for current month",
			repo: buildMockRepo(
				[]domainaccount.Account{},
				[]domaincategory.Category{},
				[]domaintransaction.Transaction{},
				[]domaintransaction.Transaction{},
				nil,
			),
			input:   pdfexport.Input{Month: "2026-01"},
			wantErr: nil,
		},
		{
			name:    "repository error is propagated",
			repo:    buildMockRepo(nil, nil, nil, nil, errors.New("db error")),
			input:   pdfexport.Input{Month: "2026-02"},
			wantErr: fmt.Errorf("export pdf: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := pdfexport.New(tc.repo)
			pdf, err := uc.Execute(context.Background(), tc.input)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error()[:20])
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pdf)
				assert.Greater(t, len(pdf), 0)
				// PDF files start with %PDF
				assert.Equal(t, "%PDF", string(pdf[:4]))
			}
			tc.repo.AssertExpectations(t)
		})
	}
}

// buildIncome creates an income transaction fixture.
func buildIncome(id string, amount float64) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", "2026-02-15")
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   "acc-1",
		Type:        domaintransaction.TransactionTypeIncome,
		Amount:      amount,
		Description: "Test Income",
		Date:        date,
		IsActive:    true,
	}
}

// buildExpense creates an expense transaction fixture.
func buildExpense(id string, amount float64, categoryID string) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", "2026-02-15")
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   "acc-1",
		CategoryID:  categoryID,
		Type:        domaintransaction.TransactionTypeExpense,
		Amount:      amount,
		Description: "Test Expense",
		Date:        date,
		IsActive:    true,
	}
}
