// Package idgen provides UUID-based ID generation for the platform.
package idgen

import "github.com/google/uuid"

// Generator generates unique string identifiers.
type Generator interface {
	NewID() string
}

// UUIDGenerator is the production implementation that generates UUID v4 strings.
type UUIDGenerator struct{}

// NewID returns a new random UUID string.
func (UUIDGenerator) NewID() string {
	return uuid.NewString()
}
