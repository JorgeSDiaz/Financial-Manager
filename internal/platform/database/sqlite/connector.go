// Package sqlite provides a SQLite database connector using the modernc pure-Go driver.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// Connector opens SQLite database connections.
type Connector struct{}

// NewConnector creates a new Connector.
func NewConnector() *Connector {
	return &Connector{}
}

// Open opens a SQLite database at the given path, creating all necessary parent
// directories. It expands a leading ~ to the user's home directory. The returned
// *sql.DB is already verified with PingContext.
func (c *Connector) Open(ctx context.Context, path string) (*sql.DB, error) {
	expanded, err := expandTilde(path)
	if err != nil {
		return nil, fmt.Errorf("sqlite: expand path %s: %w", path, err)
	}

	dir := filepath.Dir(expanded)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("sqlite: create directory %s: %w", dir, err)
	}

	db, err := sql.Open("sqlite", expanded)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open %s: %w", expanded, err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("sqlite: ping %s: %w", expanded, err)
	}

	return db, nil
}

// expandTilde replaces a leading ~ with the current user's home directory.
func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, path[1:]), nil
}
