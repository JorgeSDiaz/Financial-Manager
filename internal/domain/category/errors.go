// Package category contains domain-level errors for the category resource.
package category

import "errors"

var (
	// ErrInvalidType is returned when a category type is not 'income' or 'expense'.
	ErrInvalidType = errors.New("invalid category type: must be 'income' or 'expense'")
	// ErrEmptyName is returned when a category name is empty.
	ErrEmptyName = errors.New("category name cannot be empty")
)
