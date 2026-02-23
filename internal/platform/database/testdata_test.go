// Package database_test contains test helpers for the database package.
package database_test

import (
	"github.com/financial-manager/api/internal/platform/database"
	"github.com/financial-manager/api/internal/platform/database/migrator"
	"github.com/financial-manager/api/internal/platform/database/sqlite"
)

// buildDatabases creates a Databases instance backed by the real sqlite connector
// and migrator, suitable for happy-path integration tests.
func buildDatabases() *database.Databases {
	return database.New(sqlite.NewConnector(), migrator.New())
}
