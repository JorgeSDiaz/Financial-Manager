// Package export_test contains tests for the export use cases.
package export_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/export"
	"github.com/financial-manager/api/internal/application/export/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

func TestUseCase_ExportCSV(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repo    *mocks.Repository
		filters export.CSVFilters
		wantErr error
		wantCSV string
	}{
		{
			name: "exports transactions to CSV with headers",
			repo: buildMockRepoForCSV(
				[]domainaccount.Account{
					{ID: "acc-1", Name: "Banco"},
				},
				[]domaincategory.Category{
					{ID: "cat-1", Name: "Alimentación"},
				},
				[]domaintransaction.Transaction{
					buildExpense("tx-1", 50.00, "acc-1", "cat-1", "Groceries"),
				},
				nil,
			),
			filters: export.CSVFilters{},
			wantCSV: "date,type,amount,category,account,description\n2026-02-28,expense,50.00,Alimentación,Banco,Groceries\n",
		},
		{
			name: "exports multiple transactions",
			repo: buildMockRepoForCSV(
				[]domainaccount.Account{
					{ID: "acc-1", Name: "Banco"},
					{ID: "acc-2", Name: "Efectivo"},
				},
				[]domaincategory.Category{
					{ID: "cat-1", Name: "Alimentación"},
					{ID: "cat-2", Name: "Transporte"},
				},
				[]domaintransaction.Transaction{
					buildIncome("tx-1", 1000.00, "acc-1", "Salary"),
					buildExpense("tx-2", 50.00, "acc-1", "cat-1", "Groceries"),
					buildExpense("tx-3", 30.00, "acc-2", "cat-2", "Bus"),
				},
				nil,
			),
			filters: export.CSVFilters{},
			wantCSV: "date,type,amount,category,account,description\n2026-02-28,income,1000.00,Income,Banco,Salary\n2026-02-28,expense,50.00,Alimentación,Banco,Groceries\n2026-02-28,expense,30.00,Transporte,Efectivo,Bus\n",
		},
		{
			name:    "repository error is propagated",
			repo:    buildMockRepoForCSV(nil, nil, nil, errors.New("db error")),
			filters: export.CSVFilters{},
			wantErr: fmt.Errorf("export csv: %w", errors.New("db error")),
		},
		{
			name: "uses unknown for missing account",
			repo: buildMockRepoForCSV(
				[]domainaccount.Account{},
				[]domaincategory.Category{},
				[]domaintransaction.Transaction{
					buildExpense("tx-1", 50.00, "acc-1", "cat-1", "Groceries"),
				},
				nil,
			),
			filters: export.CSVFilters{},
			wantCSV: "date,type,amount,category,account,description\n2026-02-28,expense,50.00,Uncategorized,Unknown,Groceries\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := export.New(tc.repo)
			csv, err := uc.ExportCSV(context.Background(), tc.filters)

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, normalizeCSV(tc.wantCSV), normalizeCSV(csv))
			}
			tc.repo.AssertExpectations(t)
		})
	}
}

func TestUseCase_ExportJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		repo         *mocks.Repository
		wantErr      error
		wantContains []string
	}{
		{
			name: "exports all data to JSON",
			repo: buildMockRepoForJSON(
				[]domainaccount.Account{
					{ID: "acc-1", Name: "Banco"},
				},
				[]domaincategory.Category{
					{ID: "cat-1", Name: "Alimentación"},
				},
				[]domaintransaction.Transaction{
					buildIncome("tx-1", 1000.00, "acc-1", "Salary"),
				},
				[]domaintransaction.Transaction{
					buildExpense("tx-2", 50.00, "acc-1", "cat-1", "Groceries"),
				},
				nil,
			),
			wantContains: []string{
				`"accounts"`,
				`"categories"`,
				`"transactions"`,
				`"Banco"`,
				`"Alimentación"`,
				`"Salary"`,
				`"Groceries"`,
			},
		},
		{
			name:    "repository error is propagated",
			repo:    buildMockRepoForJSON(nil, nil, nil, nil, errors.New("db error")),
			wantErr: fmt.Errorf("export json: %w", errors.New("db error")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := export.New(tc.repo)
			json, err := uc.ExportJSON(context.Background())

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.NoError(t, err)
				jsonStr := string(json)
				for _, want := range tc.wantContains {
					assert.Contains(t, jsonStr, want)
				}
			}
			tc.repo.AssertExpectations(t)
		})
	}
}

// normalizeCSV normalizes CSV string for comparison (handles line endings).
func normalizeCSV(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// buildIncome creates an income transaction fixture.
func buildIncome(id string, amount float64, accountID, description string) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", "2026-02-28")
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   accountID,
		Type:        domaintransaction.TransactionTypeIncome,
		Amount:      amount,
		Description: description,
		Date:        date,
		IsActive:    true,
	}
}

// buildExpense creates an expense transaction fixture.
func buildExpense(id string, amount float64, accountID, categoryID, description string) domaintransaction.Transaction {
	date, _ := time.Parse("2006-01-02", "2026-02-28")
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   accountID,
		CategoryID:  categoryID,
		Type:        domaintransaction.TransactionTypeExpense,
		Amount:      amount,
		Description: description,
		Date:        date,
		IsActive:    true,
	}
}
