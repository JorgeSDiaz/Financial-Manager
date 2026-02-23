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
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"

	"github.com/financial-manager/api/internal/platform/database"
	"github.com/financial-manager/api/internal/platform/database/migrator"
	"github.com/financial-manager/api/internal/platform/database/sqlite"
)

func TestDatabases_Open(t *testing.T) {
	ctx := context.Background()

	t.Run("opens all databases and creates all expected tables", func(t *testing.T) {
		dir := t.TempDir()

		dbs := database.New(sqlite.NewConnector(), migrator.New())
		require.NoError(t, dbs.Open(ctx, dir))
		t.Cleanup(func() { _ = dbs.Close() })

		var name string
		row := dbs.Categories.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='categories'")
		require.NoError(t, row.Scan(&name))
		assert.Equal(t, "categories", name)

		row = dbs.Accounts.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='accounts'")
		require.NoError(t, row.Scan(&name))
		assert.Equal(t, "accounts", name)

		row = dbs.Transactions.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='transactions'")
		require.NoError(t, row.Scan(&name))
		assert.Equal(t, "transactions", name)

		row = dbs.Settings.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='settings'")
		require.NoError(t, row.Scan(&name))
		assert.Equal(t, "settings", name)
	})

	t.Run("seeds exactly 14 categories all with is_system=1", func(t *testing.T) {
		dir := t.TempDir()

		dbs := database.New(sqlite.NewConnector(), migrator.New())
		require.NoError(t, dbs.Open(ctx, dir))
		t.Cleanup(func() { _ = dbs.Close() })

		var total int
		row := dbs.Categories.QueryRow("SELECT COUNT(*) FROM categories")
		require.NoError(t, row.Scan(&total))
		assert.Equal(t, 14, total)

		var systemCount int
		row = dbs.Categories.QueryRow("SELECT COUNT(*) FROM categories WHERE is_system=1")
		require.NoError(t, row.Scan(&systemCount))
		assert.Equal(t, 14, systemCount)
	})

	t.Run("second Open is idempotent — tables and seed not duplicated", func(t *testing.T) {
		dir := t.TempDir()

		dbs1 := database.New(sqlite.NewConnector(), migrator.New())
		require.NoError(t, dbs1.Open(ctx, dir))
		require.NoError(t, dbs1.Close())

		dbs2 := database.New(sqlite.NewConnector(), migrator.New())
		require.NoError(t, dbs2.Open(ctx, dir))
		t.Cleanup(func() { _ = dbs2.Close() })

		var count int
		row := dbs2.Categories.QueryRow("SELECT COUNT(*) FROM categories")
		require.NoError(t, row.Scan(&count))
		assert.Equal(t, 14, count, "categories should not be duplicated on second Open")
	})

	t.Run("creates db files in the specified base directory", func(t *testing.T) {
		dir := t.TempDir()

		dbs := database.New(sqlite.NewConnector(), migrator.New())
		require.NoError(t, dbs.Open(ctx, dir))
		t.Cleanup(func() { _ = dbs.Close() })

		for _, name := range []string{"categories.db", "accounts.db", "transactions.db", "settings.db"} {
			assert.FileExists(t, filepath.Join(dir, name))
		}
	})

	t.Run("returns error when connector fails to open", func(t *testing.T) {
		dbs := database.New(buildFailingConnector(errors.New("connector: forced open error")), migrator.New())

		err := dbs.Open(ctx, t.TempDir())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "database: open categories.db")
	})

	t.Run("returns error and cleans up when migration fails on opened db", func(t *testing.T) {
		dir := t.TempDir()

		// Pre-create categories.db with a broken schema_migrations table so
		// the migration runner fails when checking applied migrations.
		catDB, err := sql.Open("sqlite", filepath.Join(dir, "categories.db"))
		require.NoError(t, err)
		_, err = catDB.Exec(`CREATE TABLE schema_migrations (wrong_col TEXT)`)
		require.NoError(t, err)
		require.NoError(t, catDB.Close())

		dbs := database.New(sqlite.NewConnector(), migrator.New())
		openErr := dbs.Open(ctx, dir)
		require.Error(t, openErr)
		assert.Contains(t, openErr.Error(), "database: migrate categories.db")
	})

	t.Run("returns error when base dir is a file not a directory", func(t *testing.T) {
		tmp := t.TempDir()
		blockingFile := filepath.Join(tmp, "not-a-dir")
		f, err := os.Create(blockingFile)
		require.NoError(t, err)
		require.NoError(t, f.Close())

		dbs := database.New(sqlite.NewConnector(), migrator.New())
		openErr := dbs.Open(ctx, blockingFile)
		require.Error(t, openErr)
		assert.Contains(t, openErr.Error(), "database: open categories.db")
	})
}

func TestDatabases_Close(t *testing.T) {
	t.Run("closes all connections without error", func(t *testing.T) {
		dir := t.TempDir()

		dbs := database.New(sqlite.NewConnector(), migrator.New())
		require.NoError(t, dbs.Open(context.Background(), dir))

		assert.NoError(t, dbs.Close())
	})
}
