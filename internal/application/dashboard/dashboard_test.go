// Package dashboard_test contains tests for the dashboard use case.
package dashboard_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/dashboard"
	"github.com/financial-manager/api/internal/application/dashboard/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repo    *mocks.Repository
		wantErr error
		wantOut dashboard.Output
	}{
		{
			name: "empty data returns dashboard with zeros",
			repo: buildMockRepo(
				nil, // accounts
				nil, // recent transactions
				nil, // category expenses
				nil, // summary transactions
				nil, // categories
				nil, // errors
			),
			wantOut: dashboard.Output{
				GlobalBalance: 0,
				MonthlySummary: dashboard.MonthlySummary{
					TotalIncome:  0,
					TotalExpense: 0,
					NetBalance:   0,
				},
				ExpensesByCategory: []dashboard.ExpenseByCategory{},
				RecentTransactions: []dashboard.RecentTransaction{},
			},
		},
		{
			name: "calculates dashboard with all data",
			repo: buildMockRepo(
				[]domainaccount.Account{account1, account2},                                              // accounts with balances
				[]domaintransaction.Transaction{tx1, tx2, tx3, tx4, tx5, tx6, tx7, tx8, tx9, tx10, tx11}, // recent (only 10 used)
				[]domaintransaction.Transaction{tx3, tx4, tx5, tx6},                                      // expense transactions
				[]domaintransaction.Transaction{tx1, tx2},                                                // summary incomes
				[]domaincategory.Category{category1, category2},                                          // categories
				nil, // errors
			),
			wantOut: dashboard.Output{
				GlobalBalance: 1500.00,
				MonthlySummary: dashboard.MonthlySummary{
					TotalIncome:  600.00,
					TotalExpense: 150.00,
					NetBalance:   450.00,
				},
				ExpensesByCategory: []dashboard.ExpenseByCategory{
					{CategoryID: "cat-1", CategoryName: "Alimentación", Total: 100.00, Percentage: 66.67},
					{CategoryID: "cat-2", CategoryName: "Transporte", Total: 50.00, Percentage: 33.33},
				},
				RecentTransactions: []dashboard.RecentTransaction{
					{ID: "tx-1", Type: "income", Amount: 500.00, Description: "Salary", CategoryName: "Income"},
					{ID: "tx-2", Type: "income", Amount: 100.00, Description: "Bonus", CategoryName: "Income"},
					{ID: "tx-3", Type: "expense", Amount: 50.00, Description: "Groceries", CategoryName: "Alimentación"},
					{ID: "tx-4", Type: "expense", Amount: 50.00, Description: "Food", CategoryName: "Alimentación"},
					{ID: "tx-5", Type: "expense", Amount: 30.00, Description: "Bus", CategoryName: "Transporte"},
					{ID: "tx-6", Type: "expense", Amount: 20.00, Description: "Taxi", CategoryName: "Transporte"},
					{ID: "tx-7", Type: "income", Amount: 200.00, Description: "Freelance", CategoryName: "Income"},
					{ID: "tx-8", Type: "expense", Amount: 40.00, Description: "Dinner", CategoryName: "Alimentación"},
					{ID: "tx-9", Type: "income", Amount: 150.00, Description: "Refund", CategoryName: "Income"},
					{ID: "tx-10", Type: "expense", Amount: 25.00, Description: "Coffee", CategoryName: "Alimentación"},
					{ID: "tx-11", Type: "expense", Amount: 15.00, Description: "Snack", CategoryName: "Alimentación"},
				},
			},
		},
		{
			name: "repository error for accounts is propagated",
			repo: buildMockRepo(
				nil, nil, nil, nil, nil,
				errors.New("db error"),
			),
			wantErr: fmt.Errorf("get dashboard: %w", errors.New("db error")),
		},
		{
			name: "expenses by category calculates percentages correctly",
			repo: buildMockRepo(
				[]domainaccount.Account{},
				[]domaintransaction.Transaction{},
				[]domaintransaction.Transaction{
					{ID: "tx-e1", CategoryID: "cat-1", Type: domaintransaction.TransactionTypeExpense, Amount: 75.00},
					{ID: "tx-e2", CategoryID: "cat-2", Type: domaintransaction.TransactionTypeExpense, Amount: 25.00},
				},
				nil,
				[]domaincategory.Category{category1, category2},
				nil,
			),
			wantOut: dashboard.Output{
				GlobalBalance: 0,
				MonthlySummary: dashboard.MonthlySummary{
					TotalIncome:  0,
					TotalExpense: 100.00,
					NetBalance:   -100.00,
				},
				ExpensesByCategory: []dashboard.ExpenseByCategory{
					{CategoryID: "cat-1", CategoryName: "Alimentación", Total: 75.00, Percentage: 75.0},
					{CategoryID: "cat-2", CategoryName: "Transporte", Total: 25.00, Percentage: 25.0},
				},
				RecentTransactions: []dashboard.RecentTransaction{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := dashboard.New(tc.repo)
			out, err := uc.Execute(context.Background())

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantOut.GlobalBalance, out.GlobalBalance)
				assert.Equal(t, tc.wantOut.MonthlySummary, out.MonthlySummary)
				assert.Equal(t, len(tc.wantOut.ExpensesByCategory), len(out.ExpensesByCategory))
				for i, want := range tc.wantOut.ExpensesByCategory {
					assert.InDelta(t, want.Percentage, out.ExpensesByCategory[i].Percentage, 0.01)
					assert.Equal(t, want.Total, out.ExpensesByCategory[i].Total)
					assert.Equal(t, want.CategoryID, out.ExpensesByCategory[i].CategoryID)
				}
				assert.Equal(t, len(tc.wantOut.RecentTransactions), len(out.RecentTransactions))
				for i, want := range tc.wantOut.RecentTransactions {
					assert.Equal(t, want.ID, out.RecentTransactions[i].ID)
					assert.Equal(t, want.Amount, out.RecentTransactions[i].Amount)
					assert.Equal(t, want.Type, out.RecentTransactions[i].Type)
					assert.Equal(t, want.CategoryName, out.RecentTransactions[i].CategoryName)
				}
			}
			tc.repo.AssertExpectations(t)
		})
	}
}

func TestUseCase_Execute_RecentTransactionsReturns10(t *testing.T) {
	t.Parallel()

	// Create 10 transactions - repo should limit them, use case just passes them through
	transactions := make([]domaintransaction.Transaction, 10)
	now := time.Now()
	for i := 0; i < 10; i++ {
		transactions[i] = domaintransaction.Transaction{
			ID:       fmt.Sprintf("tx-%d", i+1),
			Type:     domaintransaction.TransactionTypeIncome,
			Amount:   float64(i+1) * 10,
			Date:     now,
			IsActive: true,
		}
	}

	repo := &mocks.Repository{}
	repo.On("ListAccounts", mock.Anything).Return([]domainaccount.Account{}, nil).Once()
	repo.On("ListRecentTransactions", mock.Anything, 10).Return(transactions, nil).Once()
	repo.On("ListExpenseTransactions", mock.Anything, "", "", mock.Anything, mock.Anything).Return(nil, nil).Once()
	repo.On("ListIncomeTransactions", mock.Anything, "", "", mock.Anything, mock.Anything).Return(nil, nil).Once()
	repo.On("ListCategories", mock.Anything).Return(nil, nil).Once()

	uc := dashboard.New(repo)
	out, err := uc.Execute(context.Background())

	assert.NoError(t, err)
	assert.Len(t, out.RecentTransactions, 10)
}
