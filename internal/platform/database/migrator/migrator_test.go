// Package migrator_test tests the migrator package.
package migrator_test

import (
	"context"
	"database/sql"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"

	"github.com/financial-manager/api/internal/platform/database/migrator"
)

// openMemDB opens an in-memory SQLite database for testing.
func openMemDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestMigrator_Run(t *testing.T) {
	ctx := context.Background()

	t.Run("applies migration on empty database", func(t *testing.T) {
		db := openMemDB(t)
		m := migrator.New()

		fsMap := fstest.MapFS{
			"migrations/001_create_items.up.sql": {
				Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);"),
			},
		}

		err := m.Run(ctx, db, fsMap, "migrations")
		require.NoError(t, err)

		var name string
		row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='items'")
		require.NoError(t, row.Scan(&name))
		assert.Equal(t, "items", name)

		var id string
		row = db.QueryRow("SELECT id FROM schema_migrations WHERE id='001_create_items.up.sql'")
		require.NoError(t, row.Scan(&id))
		assert.Equal(t, "001_create_items.up.sql", id)
	})

	t.Run("is idempotent — second run does not fail or duplicate", func(t *testing.T) {
		db := openMemDB(t)
		m := migrator.New()

		fsMap := fstest.MapFS{
			"migrations/001_create_items.up.sql": {
				Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);"),
			},
		}

		require.NoError(t, m.Run(ctx, db, fsMap, "migrations"))
		require.NoError(t, m.Run(ctx, db, fsMap, "migrations"))

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM schema_migrations")
		require.NoError(t, row.Scan(&count))
		assert.Equal(t, 1, count)
	})

	t.Run("applies new migration on second run without re-applying existing", func(t *testing.T) {
		db := openMemDB(t)
		m := migrator.New()

		fs1 := fstest.MapFS{
			"migrations/001_create_items.up.sql": {
				Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);"),
			},
		}
		require.NoError(t, m.Run(ctx, db, fs1, "migrations"))

		fs2 := fstest.MapFS{
			"migrations/001_create_items.up.sql": {
				Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);"),
			},
			"migrations/002_create_tags.up.sql": {
				Data: []byte("CREATE TABLE tags (id TEXT PRIMARY KEY, label TEXT NOT NULL);"),
			},
		}
		require.NoError(t, m.Run(ctx, db, fs2, "migrations"))

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM schema_migrations")
		require.NoError(t, row.Scan(&count))
		assert.Equal(t, 2, count)

		var name string
		row = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tags'")
		require.NoError(t, row.Scan(&name))
		assert.Equal(t, "tags", name)
	})

	t.Run("propagates descriptive error on invalid SQL", func(t *testing.T) {
		db := openMemDB(t)
		m := migrator.New()

		fsMap := fstest.MapFS{
			"migrations/001_bad.up.sql": {
				Data: []byte("THIS IS NOT VALID SQL !!!"),
			},
		}

		err := m.Run(ctx, db, fsMap, "migrations")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "migrator: execute 001_bad.up.sql")
	})

	t.Run("applies migrations in ascending filename order", func(t *testing.T) {
		db := openMemDB(t)
		m := migrator.New()

		fsMap := fstest.MapFS{
			"migrations/001_create_items.up.sql": {
				Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);"),
			},
			"migrations/002_seed_items.up.sql": {
				Data: []byte("INSERT INTO items (id, name) VALUES ('1', 'first');"),
			},
		}

		err := m.Run(ctx, db, fsMap, "migrations")
		require.NoError(t, err)

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM items")
		require.NoError(t, row.Scan(&count))
		assert.Equal(t, 1, count)
	})

	t.Run("ignores files that do not end in .up.sql", func(t *testing.T) {
		db := openMemDB(t)
		m := migrator.New()

		fsMap := fstest.MapFS{
			"migrations/001_create_items.up.sql":   {Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);")},
			"migrations/001_create_items.down.sql": {Data: []byte("DROP TABLE items;")},
			"migrations/README.md":                 {Data: []byte("# migrations")},
		}

		err := m.Run(ctx, db, fsMap, "migrations")
		require.NoError(t, err)

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM schema_migrations")
		require.NoError(t, row.Scan(&count))
		assert.Equal(t, 1, count, "only .up.sql files should be tracked")
	})

	t.Run("returns error when migrations directory does not exist in fs", func(t *testing.T) {
		db := openMemDB(t)
		m := migrator.New()

		err := m.Run(ctx, db, fstest.MapFS{}, "nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "migrator: read directory nonexistent")
	})

	t.Run("skips directories inside migrations dir", func(t *testing.T) {
		db := openMemDB(t)
		m := migrator.New()

		fsMap := fstest.MapFS{
			"migrations/subdir":                  {Mode: fs.ModeDir},
			"migrations/001_create_items.up.sql": {Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);")},
		}

		err := m.Run(ctx, db, fsMap, "migrations")
		require.NoError(t, err)

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM schema_migrations")
		require.NoError(t, row.Scan(&count))
		assert.Equal(t, 1, count, "only files should be applied")
	})

	t.Run("returns error when db is closed before Run", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		require.NoError(t, db.Close())

		m := migrator.New()
		fsMap := fstest.MapFS{
			"migrations/001_create_items.up.sql": {Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);")},
		}

		require.Error(t, m.Run(ctx, db, fsMap, "migrations"))
	})

	t.Run("returns error when reading migration file fails", func(t *testing.T) {
		db := openMemDB(t)
		m := migrator.New()

		errFS := buildErrorReadFileFS(fstest.MapFS{
			"migrations/001_create_items.up.sql": {Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);")},
		}, "migrations/001_create_items.up.sql")

		err := m.Run(ctx, db, errFS, "migrations")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "migrator: read file 001_create_items.up.sql")
	})
}
