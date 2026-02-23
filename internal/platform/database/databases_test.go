// Package database_test tests the database package.
package database_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/financial-manager/api/internal/platform/database"
	"github.com/financial-manager/api/internal/platform/database/migrator"
	"github.com/financial-manager/api/internal/platform/database/mocks"
	"github.com/financial-manager/api/internal/platform/database/sqlite"
)

func TestDatabases_Open(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name            string
		buildDBS        func(t *testing.T) *database.Databases
		buildDir        func(t *testing.T) string
		wantErr         bool
		wantErrContains string
		check           func(t *testing.T, dbs *database.Databases)
	}{
		{
			name:     "opens all databases and creates all expected tables",
			buildDBS: func(t *testing.T) *database.Databases { return buildDatabases() },
			buildDir: func(t *testing.T) string { return t.TempDir() },
			check: func(t *testing.T, dbs *database.Databases) {
				t.Helper()

				assertTableExists(t, dbs.Categories, "categories")
				assertTableExists(t, dbs.Accounts, "accounts")
				assertTableExists(t, dbs.Transactions, "transactions")
				assertTableExists(t, dbs.Settings, "settings")
			},
		},
		{
			name:     "seeds exactly 14 categories all with is_system=1",
			buildDBS: func(t *testing.T) *database.Databases { return buildDatabases() },
			buildDir: func(t *testing.T) string { return t.TempDir() },
			check: func(t *testing.T, dbs *database.Databases) {
				t.Helper()

				var total int
				row := dbs.Categories.QueryRow("SELECT COUNT(*) FROM categories")
				require.NoError(t, row.Scan(&total))
				assert.Equal(t, 14, total)

				var systemCount int
				row = dbs.Categories.QueryRow("SELECT COUNT(*) FROM categories WHERE is_system=1")
				require.NoError(t, row.Scan(&systemCount))
				assert.Equal(t, 14, systemCount)
			},
		},
		{
			name:     "second Open is idempotent — tables and seed not duplicated",
			buildDBS: func(t *testing.T) *database.Databases { return buildDatabases() },
			buildDir: func(t *testing.T) string {
				dir := t.TempDir()
				dbs := buildDatabases()
				require.NoError(t, dbs.Open(context.Background(), dir))
				require.NoError(t, dbs.Close())
				return dir
			},
			check: func(t *testing.T, dbs *database.Databases) {
				t.Helper()

				var count int
				row := dbs.Categories.QueryRow("SELECT COUNT(*) FROM categories")
				require.NoError(t, row.Scan(&count))
				assert.Equal(t, 14, count, "categories should not be duplicated on second Open")
			},
		},
		{
			name: "returns error when connector fails to open",
			buildDBS: func(t *testing.T) *database.Databases {
				c := &mocks.Connector{}
				c.On("Open", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("connector: forced open error")).Once()
				return database.New(c, migrator.New())
			},
			buildDir:        func(t *testing.T) string { return t.TempDir() },
			wantErr:         true,
			wantErrContains: "database: open categories.db",
		},
		{
			name: "returns error when migration runner fails",
			buildDBS: func(t *testing.T) *database.Databases {
				r := &mocks.Runner{}
				r.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("string")).
					Return(errors.New("runner: forced error")).Once()
				return database.New(sqlite.NewConnector(), r)
			},
			buildDir:        func(t *testing.T) string { return t.TempDir() },
			wantErr:         true,
			wantErrContains: "database: migrate categories.db",
		},
		{
			name:     "returns error when base dir is a file not a directory",
			buildDBS: func(t *testing.T) *database.Databases { return buildDatabases() },
			buildDir: func(t *testing.T) string {
				blockingFile := filepath.Join(t.TempDir(), "not-a-dir")
				f, err := os.Create(blockingFile)
				require.NoError(t, err)
				require.NoError(t, f.Close())
				return blockingFile
			},
			wantErr:         true,
			wantErrContains: "database: open categories.db",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dir := tc.buildDir(t)
			dbs := tc.buildDBS(t)

			err := dbs.Open(ctx, dir)

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrContains)
				return
			}

			require.NoError(t, err)
			t.Cleanup(func() { _ = dbs.Close() })

			if tc.check != nil {
				tc.check(t, dbs)
			}
		})
	}
}

func TestDatabases_Open_FilesCreated(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	dbs := buildDatabases()
	require.NoError(t, dbs.Open(context.Background(), dir))
	t.Cleanup(func() { _ = dbs.Close() })

	for _, name := range []string{"categories.db", "accounts.db", "transactions.db", "settings.db"} {
		assert.FileExists(t, filepath.Join(dir, name))
	}
}

func TestDatabases_Close(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	dbs := buildDatabases()
	require.NoError(t, dbs.Open(context.Background(), dir))

	assert.NoError(t, dbs.Close())
}

// assertTableExists asserts that a table with the given name exists in the database.
func assertTableExists(t *testing.T, db *sql.DB, tableName string) {
	t.Helper()

	var name string
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", tableName)
	require.NoError(t, row.Scan(&name))
	assert.Equal(t, tableName, name)
}
