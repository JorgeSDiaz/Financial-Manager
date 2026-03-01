// Package export implements the export use cases.
package export

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Repository is the port required by the export use cases.
type Repository interface {
	ListAccounts(ctx context.Context) ([]domainaccount.Account, error)
	ListCategories(ctx context.Context) ([]domaincategory.Category, error)
	ListTransactions(ctx context.Context, tType domaintransaction.TransactionType, startDate, endDate string) ([]domaintransaction.Transaction, error)
}

// UseCase implements the export use cases.
type UseCase struct {
	repo Repository
}

// CSVFilters represents the filters for CSV export.
type CSVFilters struct {
	StartDate string
	EndDate   string
	Type      string // "income", "expense", or empty for all
}

// CSVRow represents a row in the CSV export.
type CSVRow struct {
	Date        string  `json:"date"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Account     string  `json:"account"`
	Description string  `json:"description"`
}

// BackupData represents the full backup data.
type BackupData struct {
	Accounts     []domainaccount.Account         `json:"accounts"`
	Categories   []domaincategory.Category       `json:"categories"`
	Transactions []domaintransaction.Transaction `json:"transactions"`
}

// New creates a new Export UseCase.
func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// ExportCSV exports transactions to CSV format.
func (uc *UseCase) ExportCSV(ctx context.Context, filters CSVFilters) (string, error) {
	// Get accounts and categories first for name resolution
	accounts, err := uc.repo.ListAccounts(ctx)
	if err != nil {
		return "", fmt.Errorf("export csv: %w", err)
	}

	categories, err := uc.repo.ListCategories(ctx)
	if err != nil {
		return "", fmt.Errorf("export csv: %w", err)
	}

	var tType domaintransaction.TransactionType
	switch filters.Type {
	case "income":
		tType = domaintransaction.TransactionTypeIncome
	case "expense":
		tType = domaintransaction.TransactionTypeExpense
	}

	transactions, err := uc.repo.ListTransactions(ctx, tType, filters.StartDate, filters.EndDate)
	if err != nil {
		return "", fmt.Errorf("export csv: %w", err)
	}

	accountMap := make(map[string]string)
	for _, acc := range accounts {
		accountMap[acc.ID] = acc.Name
	}

	categoryMap := make(map[string]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	// Build CSV
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// Write header
	if err := writer.Write([]string{"date", "type", "amount", "category", "account", "description"}); err != nil {
		return "", fmt.Errorf("export csv: %w", err)
	}

	// Write rows
	for _, tx := range transactions {
		accountName := accountMap[tx.AccountID]
		if accountName == "" {
			accountName = "Unknown"
		}
		categoryName := categoryMap[tx.CategoryID]
		if categoryName == "" {
			if tx.Type == domaintransaction.TransactionTypeIncome {
				categoryName = "Income"
			} else {
				categoryName = "Uncategorized"
			}
		}

		row := []string{
			tx.Date.Format("2006-01-02"),
			string(tx.Type),
			fmt.Sprintf("%.2f", tx.Amount),
			categoryName,
			accountName,
			tx.Description,
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("export csv: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("export csv: %w", err)
	}

	return sb.String(), nil
}

// ExportJSON exports all data to JSON format.
func (uc *UseCase) ExportJSON(ctx context.Context) ([]byte, error) {
	accounts, err := uc.repo.ListAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("export json: %w", err)
	}

	categories, err := uc.repo.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("export json: %w", err)
	}

	// Get all transactions (both income and expense)
	incomes, err := uc.repo.ListTransactions(ctx, domaintransaction.TransactionTypeIncome, "", "")
	if err != nil {
		return nil, fmt.Errorf("export json: %w", err)
	}

	expenses, err := uc.repo.ListTransactions(ctx, domaintransaction.TransactionTypeExpense, "", "")
	if err != nil {
		return nil, fmt.Errorf("export json: %w", err)
	}

	transactions := append(incomes, expenses...)

	data := BackupData{
		Accounts:     accounts,
		Categories:   categories,
		Transactions: transactions,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("export json: %w", err)
	}

	return jsonData, nil
}
