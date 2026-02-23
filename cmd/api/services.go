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

	// services holds all use cases ready to be injected into the HTTP layer.
	services struct {
		HealthChecker healthChecker
	}
)

// buildServices wires all use cases with their dependencies.
func buildServices() *services {
	return &services{
		HealthChecker: health.NewCheckUseCase(),
	}
}
