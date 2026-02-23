// Package migrator_test provides test helpers for the migrator package.
package migrator_test

import (
	"database/sql"
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

// openMemDB opens an in-memory SQLite database for testing and registers
// a cleanup to close it when the test ends.
func openMemDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() { _ = db.Close() })

	return db
}

// singleMigrationFS returns a MapFS with one valid migration file.
func singleMigrationFS() fstest.MapFS {
	return fstest.MapFS{
		"migrations/001_create_items.up.sql": {
			Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);"),
		},
	}
}

// errorReadFileFS wraps a MapFS and returns an error when ReadFile is called
// on a file that matches the given target filename.
type errorReadFileFS struct {
	inner  fstest.MapFS
	target string
}

func (e errorReadFileFS) Open(name string) (fs.File, error) {
	return e.inner.Open(name)
}

func (e errorReadFileFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(e.inner, name)
}

func (e errorReadFileFS) ReadFile(name string) ([]byte, error) {
	if name == e.target {
		return nil, errors.New("simulated read error")
	}

	return fs.ReadFile(e.inner, name)
}

// buildErrorReadFileFS creates an errorReadFileFS that fails on the given target file.
func buildErrorReadFileFS(inner fstest.MapFS, target string) errorReadFileFS {
	return errorReadFileFS{inner: inner, target: target}
}

// setupBrokenSchemaMigrations creates a schema_migrations table with a schema
// that causes INSERT to fail (missing applied_at column), allowing tests to
// reach the record error path in applyMigration.
func setupBrokenSchemaMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`CREATE TABLE schema_migrations (id TEXT PRIMARY KEY)`)
	require.NoError(t, err)
}
