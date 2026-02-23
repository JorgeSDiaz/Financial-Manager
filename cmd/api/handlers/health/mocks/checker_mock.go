// Package mocks contains testify mock implementations for the health handler interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domainhealth "github.com/financial-manager/api/internal/domain/health"
)

// Checker is a testify mock for the health handler checker interface.
type Checker struct {
	mock.Mock
}

// Execute mocks the checker.Execute method.
func (m *Checker) Execute(ctx context.Context) (domainhealth.Health, error) {
	args := m.Called(ctx)
	return args.Get(0).(domainhealth.Health), args.Error(1)
}
