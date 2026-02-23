package main

import (
	"context"

	"github.com/financial-manager/api/internal/application/health"
	domainhealth "github.com/financial-manager/api/internal/domain/health"
)

type (
	// healthChecker is the contract for retrieving application health status.
	healthChecker interface {
		Execute(ctx context.Context) (domainhealth.Health, error)
	}

	// dependencies holds all resolved use cases injected into the HTTP layer.
	dependencies struct {
		HealthChecker healthChecker
	}
)

// buildDependencies wires all adapters and use cases together.
func buildDependencies() *dependencies {
	return &dependencies{
		HealthChecker: health.NewCheckUseCase(),
	}
}
