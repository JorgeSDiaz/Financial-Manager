// Package dashboard implements the dashboard use case.
package dashboard

import (
	"context"
	"fmt"
	"time"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Repository is the port required by the dashboard use case.
type Repository interface {
	ListAccounts(ctx context.Context) ([]domainaccount.Account, error)
	ListRecentTransactions(ctx context.Context, limit int) ([]domaintransaction.Transaction, error)
	ListExpenseTransactions(ctx context.Context, accountID, categoryID, startDate, endDate string) ([]domaintransaction.Transaction, error)
	ListIncomeTransactions(ctx context.Context, accountID, categoryID, startDate, endDate string) ([]domaintransaction.Transaction, error)
	ListCategories(ctx context.Context) ([]domaincategory.Category, error)
}

// UseCase implements the dashboard use case.
type UseCase struct {
	repo Repository
}

// Output represents the dashboard response.
type Output struct {
	GlobalBalance      float64             `json:"global_balance"`
	MonthlySummary     MonthlySummary      `json:"monthly_summary"`
	ExpensesByCategory []ExpenseByCategory `json:"expenses_by_category"`
	RecentTransactions []RecentTransaction `json:"recent_transactions"`
}

// MonthlySummary represents the financial summary for the current month.
type MonthlySummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetBalance   float64 `json:"net_balance"`
}

// ExpenseByCategory represents the expense breakdown by category.
type ExpenseByCategory struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Total        float64 `json:"total"`
	Percentage   float64 `json:"percentage"`
}

// RecentTransaction represents a transaction for the dashboard list.
type RecentTransaction struct {
	ID           string  `json:"id"`
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
	Date         string  `json:"date"`
	Description  string  `json:"description"`
	CategoryName string  `json:"category_name"`
}

// New creates a new Dashboard UseCase.
func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute retrieves the dashboard data.
func (uc *UseCase) Execute(ctx context.Context) (Output, error) {
	// Get global balance
	accounts, err := uc.repo.ListAccounts(ctx)
	if err != nil {
		return Output{}, fmt.Errorf("get dashboard: %w", err)
	}

	var globalBalance float64
	for _, acc := range accounts {
		globalBalance += acc.CurrentBalance
	}

	// Get current month period
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	startDateStr := startOfMonth.Format("2006-01-02")
	endDateStr := endOfMonth.Format("2006-01-02")

	// Get monthly summary
	incomes, err := uc.repo.ListIncomeTransactions(ctx, "", "", startDateStr, endDateStr)
	if err != nil {
		return Output{}, fmt.Errorf("get dashboard: %w", err)
	}

	expenses, err := uc.repo.ListExpenseTransactions(ctx, "", "", startDateStr, endDateStr)
	if err != nil {
		return Output{}, fmt.Errorf("get dashboard: %w", err)
	}

	var totalIncome, totalExpense float64
	for _, tx := range incomes {
		totalIncome += tx.Amount
	}
	for _, tx := range expenses {
		totalExpense += tx.Amount
	}

	// Get expenses by category
	categories, err := uc.repo.ListCategories(ctx)
	if err != nil {
		return Output{}, fmt.Errorf("get dashboard: %w", err)
	}

	categoryMap := make(map[string]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	expenseByCategory := make(map[string]float64)
	for _, tx := range expenses {
		expenseByCategory[tx.CategoryID] += tx.Amount
	}

	var expensesByCategory []ExpenseByCategory
	for catID, total := range expenseByCategory {
		var percentage float64
		if totalExpense > 0 {
			percentage = (total / totalExpense) * 100
		}
		catName := categoryMap[catID]
		if catName == "" {
			catName = "Uncategorized"
		}
		expensesByCategory = append(expensesByCategory, ExpenseByCategory{
			CategoryID:   catID,
			CategoryName: catName,
			Total:        total,
			Percentage:   percentage,
		})
	}

	// Get recent transactions (last 10, mixed income and expense)
	recentTxs, err := uc.repo.ListRecentTransactions(ctx, 10)
	if err != nil {
		return Output{}, fmt.Errorf("get dashboard: %w", err)
	}

	var recentTransactions []RecentTransaction
	for _, tx := range recentTxs {
		catName := categoryMap[tx.CategoryID]
		if catName == "" {
			if tx.Type == domaintransaction.TransactionTypeIncome {
				catName = "Income"
			} else {
				catName = "Uncategorized"
			}
		}
		recentTransactions = append(recentTransactions, RecentTransaction{
			ID:           tx.ID,
			Type:         string(tx.Type),
			Amount:       tx.Amount,
			Date:         tx.Date.Format("2006-01-02"),
			Description:  tx.Description,
			CategoryName: catName,
		})
	}

	return Output{
		GlobalBalance: globalBalance,
		MonthlySummary: MonthlySummary{
			TotalIncome:  totalIncome,
			TotalExpense: totalExpense,
			NetBalance:   totalIncome - totalExpense,
		},
		ExpensesByCategory: expensesByCategory,
		RecentTransactions: recentTransactions,
	}, nil
}
