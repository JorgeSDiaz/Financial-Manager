// Package shared contains domain-level errors and types shared across all domain packages.
package shared

import "errors"

// ErrNotFound is returned when a requested resource does not exist.
var ErrNotFound = errors.New("not found")
