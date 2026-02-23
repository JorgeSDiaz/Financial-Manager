package main

import (
	"context"
	"fmt"

	"github.com/financial-manager/api/internal/application/health"
	domainhealth "github.com/financial-manager/api/internal/domain/health"
	"github.com/financial-manager/api/internal/platform/config"
	"github.com/financial-manager/api/internal/platform/database"
	"github.com/financial-manager/api/internal/platform/database/migrator"
	"github.com/financial-manager/api/internal/platform/database/sqlite"
)

type (
	// healthChecker is the contract for retrieving application health status.
	healthChecker interface {
		Execute(ctx context.Context) (domainhealth.Health, error)
	}

	// dependencies holds all resolved use cases injected into the HTTP layer.
	dependencies struct {
		HealthChecker healthChecker
		Databases     *database.Databases
	}
)

// buildDependencies wires all adapters and use cases together.
// It returns an error if any infrastructure component fails to initialize.
func buildDependencies(cfg *config.Config) (*dependencies, error) {
	dbs := database.New(sqlite.NewConnector(), migrator.New())
	if err := dbs.Open(context.Background(), cfg.DatabaseDir); err != nil {
		return nil, fmt.Errorf("buildDependencies: %w", err)
	}

	return &dependencies{
		HealthChecker: health.NewCheckUseCase(),
		Databases:     dbs,
	}, nil
}
