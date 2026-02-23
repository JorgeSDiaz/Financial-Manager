package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
)

// connector opens a single database connection at the given file path.
type connector interface {
	Open(ctx context.Context, path string) (*sql.DB, error)
}

// runner applies SQL migrations to an open database.
type runner interface {
	Run(ctx context.Context, db *sql.DB, migrationFS fs.FS, dir string) error
}

// dbEntry pairs a database filename with its target connection field and migration directory.
type dbEntry struct {
	name   string
	target **sql.DB
	migDir string
}

// Databases holds open connections to each SQLite database used by the application.
type Databases struct {
	connector connector
	runner    runner

	// Categories is the database storing expense and income categories.
	Categories *sql.DB
	// Accounts is the database storing user accounts.
	Accounts *sql.DB
	// Transactions is the database storing financial transactions.
	Transactions *sql.DB
	// Settings is the database storing application settings.
	Settings *sql.DB
}

// New creates a new Databases instance using the provided connector and migration runner.
// Call Open to establish connections and apply migrations.
func New(c connector, r runner) *Databases {
	return &Databases{connector: c, runner: r}
}

// Open opens all application databases in baseDir, applying migrations to each.
// If any database fails to open or migrate, all opened connections are closed
// before returning the error.
func (d *Databases) Open(ctx context.Context, baseDir string) error {
	entries := []dbEntry{
		{"categories.db", &d.Categories, "migrations/categories"},
		{"accounts.db", &d.Accounts, "migrations/accounts"},
		{"transactions.db", &d.Transactions, "migrations/transactions"},
		{"settings.db", &d.Settings, "migrations/settings"},
	}

	for _, entry := range entries {
		db, err := d.connector.Open(ctx, filepath.Join(baseDir, entry.name))
		if err != nil {
			_ = d.Close()
			return fmt.Errorf("database: open %s: %w", entry.name, err)
		}

		*entry.target = db

		if err := d.runner.Run(ctx, db, migrationFiles, entry.migDir); err != nil {
			_ = d.Close()
			return fmt.Errorf("database: migrate %s: %w", entry.name, err)
		}
	}

	return nil
}

// Close closes all database connections, accumulating any errors encountered.
func (d *Databases) Close() error {
	return errors.Join(
		closeDB(d.Categories),
		closeDB(d.Accounts),
		closeDB(d.Transactions),
		closeDB(d.Settings),
	)
}

// closeDB closes a *sql.DB only if it is non-nil.
func closeDB(db *sql.DB) error {
	if db == nil {
		return nil
	}

	return db.Close()
}
