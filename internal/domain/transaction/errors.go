// Package transaction contains the Transaction entity and its value objects.
package transaction

import "errors"

var ErrTransactionNotFound = errors.New("transaction not found")
var ErrAccountNotFound = errors.New("account not found")
var ErrCategoryNotFound = errors.New("category not found")
var ErrInvalidAmount = errors.New("amount must be positive")
var ErrInsufficientBalance = errors.New("insufficient balance in account")
