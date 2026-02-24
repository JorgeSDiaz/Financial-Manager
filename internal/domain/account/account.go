// Package account contains the Account entity and its value objects.
package account

import "time"

type (
	// AccountType represents the category of a financial account.
	AccountType string

	// Account represents a user's financial account.
	Account struct {
		ID             string
		Name           string
		Type           AccountType
		InitialBalance float64
		CurrentBalance float64
		Currency       string
		Color          string
		Icon           string
		IsActive       bool
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
)

const (
	// AccountTypeCash represents physical cash.
	AccountTypeCash AccountType = "cash"
	// AccountTypeBank represents a bank account.
	AccountTypeBank AccountType = "bank"
	// AccountTypeCreditCard represents a credit card account.
	AccountTypeCreditCard AccountType = "credit_card"
	// AccountTypeSavings represents a savings account.
	AccountTypeSavings AccountType = "savings"
)
