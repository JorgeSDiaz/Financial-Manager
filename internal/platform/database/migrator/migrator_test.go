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

	"github.com/financial-manager/api/internal/platform/database/migrator"
)

func TestMigrator_Run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		buildDB         func(t *testing.T) *sql.DB
		buildFS         func() fs.FS
		dir             string
		wantErr         bool
		wantErrContains string
		check           func(t *testing.T, db *sql.DB)
	}{
		{
			name:    "applies migration on empty database",
			buildDB: func(t *testing.T) *sql.DB { return openMemDB(t) },
			buildFS: func() fs.FS { return singleMigrationFS() },
			dir:     "migrations",
			check: func(t *testing.T, db *sql.DB) {
				t.Helper()

				var tableName string
				row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='items'")
				require.NoError(t, row.Scan(&tableName))
				assert.Equal(t, "items", tableName)

				var id string
				row = db.QueryRow("SELECT id FROM schema_migrations WHERE id='001_create_items.up.sql'")
				require.NoError(t, row.Scan(&id))
				assert.Equal(t, "001_create_items.up.sql", id)
			},
		},
		{
			name:    "is idempotent — second run does not fail or duplicate",
			buildDB: func(t *testing.T) *sql.DB { return openMemDB(t) },
			buildFS: func() fs.FS { return singleMigrationFS() },
			dir:     "migrations",
			check: func(t *testing.T, db *sql.DB) {
				t.Helper()

				// Run a second time on the same DB — must not error (tested via outer require.NoError)
				// and must not duplicate the migration record.
				m := migrator.New()
				require.NoError(t, m.Run(context.Background(), db, singleMigrationFS(), "migrations"))

				var count int
				row := db.QueryRow("SELECT COUNT(*) FROM schema_migrations")
				require.NoError(t, row.Scan(&count))
				assert.Equal(t, 1, count)
			},
		},
		{
			name:    "applies new migration on second run without re-applying existing",
			buildDB: func(t *testing.T) *sql.DB { return openMemDB(t) },
			buildFS: func() fs.FS { return singleMigrationFS() },
			dir:     "migrations",
			check: func(t *testing.T, db *sql.DB) {
				t.Helper()

				fs2 := fstest.MapFS{
					"migrations/001_create_items.up.sql": {
						Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);"),
					},
					"migrations/002_create_tags.up.sql": {
						Data: []byte("CREATE TABLE tags (id TEXT PRIMARY KEY, label TEXT NOT NULL);"),
					},
				}

				m := migrator.New()
				require.NoError(t, m.Run(context.Background(), db, fs2, "migrations"))

				var count int
				row := db.QueryRow("SELECT COUNT(*) FROM schema_migrations")
				require.NoError(t, row.Scan(&count))
				assert.Equal(t, 2, count)

				var tableName string
				row = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tags'")
				require.NoError(t, row.Scan(&tableName))
				assert.Equal(t, "tags", tableName)
			},
		},
		{
			name:    "applies migrations in ascending filename order",
			buildDB: func(t *testing.T) *sql.DB { return openMemDB(t) },
			buildFS: func() fs.FS {
				return fstest.MapFS{
					"migrations/001_create_items.up.sql": {
						Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);"),
					},
					"migrations/002_seed_items.up.sql": {
						Data: []byte("INSERT INTO items (id, name) VALUES ('1', 'first');"),
					},
				}
			},
			dir: "migrations",
			check: func(t *testing.T, db *sql.DB) {
				t.Helper()

				var count int
				row := db.QueryRow("SELECT COUNT(*) FROM items")
				require.NoError(t, row.Scan(&count))
				assert.Equal(t, 1, count)
			},
		},
		{
			name:    "ignores files that do not end in .up.sql",
			buildDB: func(t *testing.T) *sql.DB { return openMemDB(t) },
			buildFS: func() fs.FS {
				return fstest.MapFS{
					"migrations/001_create_items.up.sql":   {Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);")},
					"migrations/001_create_items.down.sql": {Data: []byte("DROP TABLE items;")},
					"migrations/README.md":                 {Data: []byte("# migrations")},
				}
			},
			dir: "migrations",
			check: func(t *testing.T, db *sql.DB) {
				t.Helper()

				var count int
				row := db.QueryRow("SELECT COUNT(*) FROM schema_migrations")
				require.NoError(t, row.Scan(&count))
				assert.Equal(t, 1, count, "only .up.sql files should be tracked")
			},
		},
		{
			name:    "skips directories inside migrations dir",
			buildDB: func(t *testing.T) *sql.DB { return openMemDB(t) },
			buildFS: func() fs.FS {
				return fstest.MapFS{
					"migrations/subdir":                  {Mode: fs.ModeDir},
					"migrations/001_create_items.up.sql": {Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);")},
				}
			},
			dir: "migrations",
			check: func(t *testing.T, db *sql.DB) {
				t.Helper()

				var count int
				row := db.QueryRow("SELECT COUNT(*) FROM schema_migrations")
				require.NoError(t, row.Scan(&count))
				assert.Equal(t, 1, count, "only files should be applied")
			},
		},
		{
			name: "returns error when db is closed before Run",
			buildDB: func(t *testing.T) *sql.DB {
				db, err := sql.Open("sqlite", ":memory:")
				require.NoError(t, err)
				require.NoError(t, db.Close())
				return db
			},
			buildFS:         func() fs.FS { return singleMigrationFS() },
			dir:             "migrations",
			wantErr:         true,
			wantErrContains: "migrator: create schema_migrations table",
		},
		{
			name:            "returns error when migrations directory does not exist in fs",
			buildDB:         func(t *testing.T) *sql.DB { return openMemDB(t) },
			buildFS:         func() fs.FS { return fstest.MapFS{} },
			dir:             "nonexistent",
			wantErr:         true,
			wantErrContains: "migrator: read directory nonexistent",
		},
		{
			name:    "returns error when reading migration file fails",
			buildDB: func(t *testing.T) *sql.DB { return openMemDB(t) },
			buildFS: func() fs.FS {
				return buildErrorReadFileFS(
					fstest.MapFS{
						"migrations/001_create_items.up.sql": {Data: []byte("CREATE TABLE items (id TEXT PRIMARY KEY, name TEXT NOT NULL);")},
					},
					"migrations/001_create_items.up.sql",
				)
			},
			dir:             "migrations",
			wantErr:         true,
			wantErrContains: "migrator: read file 001_create_items.up.sql",
		},
		{
			name:    "returns error when migration SQL is invalid",
			buildDB: func(t *testing.T) *sql.DB { return openMemDB(t) },
			buildFS: func() fs.FS {
				return fstest.MapFS{
					"migrations/001_bad.up.sql": {Data: []byte("THIS IS NOT VALID SQL !!!")},
				}
			},
			dir:             "migrations",
			wantErr:         true,
			wantErrContains: "migrator: execute 001_bad.up.sql",
		},
		{
			name: "returns error when recording applied migration fails",
			buildDB: func(t *testing.T) *sql.DB {
				db := openMemDB(t)
				setupBrokenSchemaMigrations(t, db)
				return db
			},
			buildFS:         func() fs.FS { return singleMigrationFS() },
			dir:             "migrations",
			wantErr:         true,
			wantErrContains: "migrator: record 001_create_items.up.sql",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := tc.buildDB(t)
			migrationFS := tc.buildFS()

			m := migrator.New()
			err := m.Run(ctx, db, migrationFS, tc.dir)

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrContains)
				return
			}

			require.NoError(t, err)

			if tc.check != nil {
				tc.check(t, db)
			}
		})
	}
}
