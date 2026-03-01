// Package pdfexport implements the PDF export use case.
package pdfexport

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Repository is the port required by the PDF export use case.
type Repository interface {
	ListAccounts(ctx context.Context) ([]domainaccount.Account, error)
	ListCategories(ctx context.Context) ([]domaincategory.Category, error)
	ListTransactions(ctx context.Context, tType domaintransaction.TransactionType, startDate, endDate string) ([]domaintransaction.Transaction, error)
}

// UseCase implements the PDF export use case.
type UseCase struct {
	repo Repository
}

// Input represents the PDF export request.
type Input struct {
	Month         string `json:"month"` // Format: "2026-02"
	IncludeCharts bool   `json:"include_charts"`
}

// New creates a new PDF Export UseCase.
func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute generates a PDF report for the specified month.
func (uc *UseCase) Execute(ctx context.Context, in Input) ([]byte, error) {
	// Validate month format
	monthDate, err := time.Parse("2006-01", in.Month)
	if err != nil {
		return nil, fmt.Errorf("invalid month format: %w", err)
	}

	// Calculate month boundaries
	startOfMonth := time.Date(monthDate.Year(), monthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	startDateStr := startOfMonth.Format("2006-01-02")
	endDateStr := endOfMonth.Format("2006-01-02")

	// Get data
	accounts, err := uc.repo.ListAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("export pdf: %w", err)
	}

	categories, err := uc.repo.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("export pdf: %w", err)
	}

	incomes, err := uc.repo.ListTransactions(ctx, domaintransaction.TransactionTypeIncome, startDateStr, endDateStr)
	if err != nil {
		return nil, fmt.Errorf("export pdf: %w", err)
	}

	expenses, err := uc.repo.ListTransactions(ctx, domaintransaction.TransactionTypeExpense, startDateStr, endDateStr)
	if err != nil {
		return nil, fmt.Errorf("export pdf: %w", err)
	}

	// Calculate summary
	var totalIncome, totalExpense float64
	for _, tx := range incomes {
		totalIncome += tx.Amount
	}
	for _, tx := range expenses {
		totalExpense += tx.Amount
	}
	netBalance := totalIncome - totalExpense

	// Calculate expenses by category
	categoryMap := make(map[string]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	expenseByCategory := make(map[string]float64)
	for _, tx := range expenses {
		expenseByCategory[tx.CategoryID] += tx.Amount
	}

	// Create PDF
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("Financial Report - %s", monthDate.Format("January 2006")))
	pdf.Ln(15)

	// Summary Section
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Monthly Summary")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Total Income: $%.2f", totalIncome))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Total Expense: $%.2f", totalExpense))
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Net Balance: $%.2f", netBalance))
	pdf.Ln(15)

	// Expenses by Category Section
	if len(expenseByCategory) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "Expenses by Category")
		pdf.Ln(10)

		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(80, 8, "Category")
		pdf.Cell(50, 8, "Amount")
		pdf.Cell(50, 8, "Percentage")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 11)
		for catID, amount := range expenseByCategory {
			catName := categoryMap[catID]
			if catName == "" {
				catName = "Uncategorized"
			}
			percentage := 0.0
			if totalExpense > 0 {
				percentage = (amount / totalExpense) * 100
			}

			pdf.Cell(80, 8, catName)
			pdf.Cell(50, 8, fmt.Sprintf("$%.2f", amount))
			pdf.Cell(50, 8, fmt.Sprintf("%.1f%%", percentage))
			pdf.Ln(8)
		}
		pdf.Ln(10)
	}

	// Transactions Section
	allTransactions := append(incomes, expenses...)
	if len(allTransactions) > 0 {
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "Transactions")
		pdf.Ln(10)

		// Table header
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(30, 8, "Date")
		pdf.Cell(25, 8, "Type")
		pdf.Cell(35, 8, "Amount")
		pdf.Cell(50, 8, "Category")
		pdf.Cell(50, 8, "Description")
		pdf.Ln(8)

		// Table rows
		pdf.SetFont("Arial", "", 9)
		for _, tx := range allTransactions {
			catName := categoryMap[tx.CategoryID]
			if catName == "" {
				if tx.Type == domaintransaction.TransactionTypeIncome {
					catName = "Income"
				} else {
					catName = "Uncategorized"
				}
			}

			desc := tx.Description
			if len(desc) > 20 {
				desc = desc[:17] + "..."
			}

			pdf.Cell(30, 7, tx.Date.Format("2006-01-02"))
			pdf.Cell(25, 7, string(tx.Type))
			pdf.Cell(35, 7, fmt.Sprintf("$%.2f", tx.Amount))
			pdf.Cell(50, 7, catName)
			pdf.Cell(50, 7, desc)
			pdf.Ln(7)
		}
	}

	// Accounts Section
	if len(accounts) > 0 {
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "Account Balances")
		pdf.Ln(10)

		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(80, 8, "Account")
		pdf.Cell(50, 8, "Type")
		pdf.Cell(50, 8, "Balance")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 11)
		for _, acc := range accounts {
			pdf.Cell(80, 8, acc.Name)
			pdf.Cell(50, 8, string(acc.Type))
			pdf.Cell(50, 8, fmt.Sprintf("$%.2f", acc.CurrentBalance))
			pdf.Ln(8)
		}
	}

	// Generate PDF bytes
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("export pdf: %w", err)
	}

	return buf.Bytes(), nil
}
