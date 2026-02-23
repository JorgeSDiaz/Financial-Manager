// Package health implements the health check use case.
package health

import (
	"context"
	"time"

	domainhealth "github.com/financial-manager/api/internal/domain/health"
)

const appVersion = "1.0.0"

// CheckUseCase implements the health check use case.
type CheckUseCase struct{}

// NewCheckUseCase creates a new CheckUseCase.
func NewCheckUseCase() *CheckUseCase {
	return &CheckUseCase{}
}

// Execute returns the current health status of the application.
func (uc *CheckUseCase) Execute(_ context.Context) (domainhealth.Health, error) {
	return domainhealth.Health{
		Status:    domainhealth.StatusUp,
		Timestamp: time.Now().UTC(),
		Version:   appVersion,
	}, nil
}
