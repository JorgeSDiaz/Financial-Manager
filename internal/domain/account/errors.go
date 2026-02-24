// Package account contains the Account entity and its value objects.
package account

import "errors"

// ErrAccountHasTransactions is returned when attempting to delete an account
// that still has associated transactions.
var ErrAccountHasTransactions = errors.New("account has transactions and cannot be deleted")
