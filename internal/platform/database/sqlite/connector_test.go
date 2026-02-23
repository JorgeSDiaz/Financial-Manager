// Package sqlite_test tests the sqlite package.
package sqlite_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/financial-manager/api/internal/platform/database/sqlite"
)

func TestConnector_Open(t *testing.T) {
	c := sqlite.NewConnector()
	ctx := context.Background()

	t.Run("opens valid path and returns pingable db handle", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.db")

		db, err := c.Open(ctx, path)
		require.NoError(t, err)
		require.NotNil(t, db)
		t.Cleanup(func() { _ = db.Close() })

		assert.NoError(t, db.Ping())
	})

	t.Run("creates nested directories automatically", func(t *testing.T) {
		dir := t.TempDir()
		nestedPath := filepath.Join(dir, "a", "b", "c", "test.db")

		db, err := c.Open(ctx, nestedPath)
		require.NoError(t, err)
		require.NotNil(t, db)
		t.Cleanup(func() { _ = db.Close() })

		_, statErr := os.Stat(filepath.Dir(nestedPath))
		assert.NoError(t, statErr, "nested directories should exist")
	})

	t.Run("expands tilde to home directory", func(t *testing.T) {
		home, err := os.UserHomeDir()
		require.NoError(t, err)

		subDir := filepath.Join(home, ".financial-manager-test-"+t.Name())
		t.Cleanup(func() { _ = os.RemoveAll(subDir) })

		path := filepath.Join("~", ".financial-manager-test-"+t.Name(), "test.db")

		db, err := c.Open(ctx, path)
		require.NoError(t, err)
		require.NotNil(t, db)
		t.Cleanup(func() { _ = db.Close() })

		expanded := filepath.Join(subDir, "test.db")
		_, statErr := os.Stat(expanded)
		assert.NoError(t, statErr, "expanded path should exist")
	})

	t.Run("returns error when parent path is a file not a directory", func(t *testing.T) {
		dir := t.TempDir()
		blockingFile := filepath.Join(dir, "not-a-dir")
		err := createFile(t, blockingFile)
		require.NoError(t, err)

		dbPath := filepath.Join(blockingFile, "test.db")

		db, err := c.Open(ctx, dbPath)
		assert.Nil(t, db)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "sqlite: create directory"),
			"error should contain 'sqlite: create directory', got: %s", err.Error())
	})

	t.Run("returns error when db path is a directory not a file", func(t *testing.T) {
		dir := t.TempDir()
		dbPath := filepath.Join(dir, "test.db")
		err := os.MkdirAll(dbPath, 0o755)
		require.NoError(t, err)

		db, err := c.Open(ctx, dbPath)
		assert.Nil(t, db)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "sqlite: ping"),
			"error should contain 'sqlite: ping', got: %s", err.Error())
	})

	t.Run("returns error when HOME is not set for tilde expansion", func(t *testing.T) {
		origHome := os.Getenv("HOME")
		require.NoError(t, os.Unsetenv("HOME"))
		t.Cleanup(func() {
			if origHome != "" {
				_ = os.Setenv("HOME", origHome)
			}
		})

		db, err := c.Open(ctx, "~/test.db")
		assert.Nil(t, db)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sqlite: expand path ~/test.db")
	})
}
