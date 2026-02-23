// Package migrator applies SQL migration files to a database in a deterministic,
// idempotent manner.
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

// db is the subset of *sql.DB used by Migrator, defined here per
// interface-at-consumer to allow injection of test doubles.
type db interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

const createSchemaMigrationsSQL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    id         TEXT PRIMARY KEY,
    applied_at TEXT NOT NULL
);`

// Migrator applies SQL migration files to a database in a deterministic,
// idempotent manner. It is stateless and safe to reuse across multiple calls.
type Migrator struct{}

// New creates a new Migrator ready to apply migrations.
func New() *Migrator {
	return &Migrator{}
}

// Run applies all *.up.sql files found in dir within the given filesystem to sqlDB.
// It creates the schema_migrations tracking table if necessary. Files are applied
// in ascending filename order and skipped if already recorded as applied.
// Run is idempotent: repeated calls with the same filesystem are safe.
func (m *Migrator) Run(ctx context.Context, sqlDB *sql.DB, migrationFS fs.FS, dir string) error {
	return m.run(ctx, sqlDB, migrationFS, dir)
}

// run is the internal implementation of Run that accepts the db interface,
// enabling injection of test doubles without changing the public API.
func (m *Migrator) run(ctx context.Context, conn db, migrationFS fs.FS, dir string) error {
	if _, err := conn.ExecContext(ctx, createSchemaMigrationsSQL); err != nil {
		return fmt.Errorf("migrator: create schema_migrations table: %w", err)
	}

	entries, err := fs.ReadDir(migrationFS, dir)
	if err != nil {
		return fmt.Errorf("migrator: read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}

		if err := m.applyMigration(ctx, conn, migrationFS, dir, entry.Name()); err != nil {
			return err
		}
	}

	return nil
}

// applyMigration checks if a migration has already been applied and, if not,
// executes it and records it in schema_migrations.
func (m *Migrator) applyMigration(ctx context.Context, conn db, migrationFS fs.FS, dir, filename string) error {
	applied, err := m.isApplied(ctx, conn, filename)
	if err != nil {
		return fmt.Errorf("migrator: check applied %s: %w", filename, err)
	}

	if applied {
		return nil
	}

	content, err := fs.ReadFile(migrationFS, filepath.Join(dir, filename))
	if err != nil {
		return fmt.Errorf("migrator: read file %s: %w", filename, err)
	}

	if _, err := conn.ExecContext(ctx, string(content)); err != nil {
		return fmt.Errorf("migrator: execute %s: %w", filename, err)
	}

	if _, err := conn.ExecContext(ctx,
		"INSERT INTO schema_migrations (id, applied_at) VALUES (?, ?)",
		filename,
		time.Now().UTC().Format(time.RFC3339),
	); err != nil {
		return fmt.Errorf("migrator: record %s: %w", filename, err)
	}

	return nil
}

// isApplied reports whether a migration with the given id has already been applied.
func (m *Migrator) isApplied(ctx context.Context, conn db, id string) (bool, error) {
	var count int
	row := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE id = ?", id)
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}
