// Package category contains the Category entity and its value objects.
package category

import "time"

type (
	// Type represents the classification of a category.
	Type string

	// Category represents a transaction category (income or expense).
	Category struct {
		ID        string
		Name      string
		Type      Type
		Color     string
		Icon      string
		IsSystem  bool
		IsActive  bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

const (
	// TypeExpense represents expense categories.
	TypeExpense Type = "expense"
	// TypeIncome represents income categories.
	TypeIncome Type = "income"
)
