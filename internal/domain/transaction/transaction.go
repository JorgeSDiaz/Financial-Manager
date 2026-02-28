// Package transaction contains the Transaction entity and its value objects.
package transaction

import "time"

type (
	// TransactionType represents the type of a financial transaction.
	TransactionType string

	// Transaction represents a financial transaction (income or expense).
	Transaction struct {
		ID          string
		AccountID   string
		CategoryID  string
		Type        TransactionType
		Amount      float64
		Description string
		Date        time.Time
		IsActive    bool
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)

const (
	// TransactionTypeIncome represents an income transaction.
	TransactionTypeIncome TransactionType = "income"
	// TransactionTypeExpense represents an expense transaction.
	TransactionTypeExpense TransactionType = "expense"
)
