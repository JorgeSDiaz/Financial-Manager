// Package sqlite_test tests the sqlite package.
package sqlite_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/financial-manager/api/internal/platform/database/sqlite"
)

func TestConnector_Open(t *testing.T) {
	t.Parallel()

	c := sqlite.NewConnector()
	ctx := context.Background()

	tests := []struct {
		name            string
		buildPath       func(t *testing.T) string
		wantErr         bool
		wantErrContains string
		check           func(t *testing.T, path string)
	}{
		{
			name: "opens valid path and returns pingable db handle",
			buildPath: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "test.db")
			},
		},
		{
			name: "creates nested directories automatically",
			buildPath: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "a", "b", "c", "test.db")
			},
			check: func(t *testing.T, path string) {
				t.Helper()
				_, err := os.Stat(filepath.Dir(path))
				assert.NoError(t, err, "nested directories should exist")
			},
		},
		{
			name: "expands tilde to home directory",
			buildPath: func(t *testing.T) string {
				home, err := os.UserHomeDir()
				require.NoError(t, err)
				subDir := filepath.Join(home, ".financial-manager-test-connector-open")
				t.Cleanup(func() { _ = os.RemoveAll(subDir) })
				return filepath.Join("~", ".financial-manager-test-connector-open", "test.db")
			},
			check: func(t *testing.T, path string) {
				t.Helper()
				home, err := os.UserHomeDir()
				require.NoError(t, err)
				expanded := filepath.Join(home, ".financial-manager-test-connector-open", "test.db")
				_, err = os.Stat(expanded)
				assert.NoError(t, err, "expanded path should exist on disk")
			},
		},
		{
			name: "returns error when parent path is a file not a directory",
			buildPath: func(t *testing.T) string {
				blockingFile := filepath.Join(t.TempDir(), "not-a-dir")
				require.NoError(t, createFile(t, blockingFile))
				return filepath.Join(blockingFile, "test.db")
			},
			wantErr:         true,
			wantErrContains: "sqlite: create directory",
		},
		{
			name: "returns error when db path is a directory not a file",
			buildPath: func(t *testing.T) string {
				dbPath := filepath.Join(t.TempDir(), "test.db")
				require.NoError(t, os.MkdirAll(dbPath, 0o755))
				return dbPath
			},
			wantErr:         true,
			wantErrContains: "sqlite: ping",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path := tc.buildPath(t)
			db, err := c.Open(ctx, path)

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, db)
				assert.Contains(t, err.Error(), tc.wantErrContains)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, db)
			t.Cleanup(func() { _ = db.Close() })
			assert.NoError(t, db.Ping())

			if tc.check != nil {
				tc.check(t, path)
			}
		})
	}
}

// TestConnector_Open_TildeExpansion_NoHome verifies that Open returns a descriptive
// error when HOME is unset and a tilde path is provided. This test runs sequentially
// (no t.Parallel) because it modifies the HOME environment variable.
func TestConnector_Open_TildeExpansion_NoHome(t *testing.T) {
	origHome := os.Getenv("HOME")
	require.NoError(t, os.Unsetenv("HOME"))
	t.Cleanup(func() {
		if origHome != "" {
			_ = os.Setenv("HOME", origHome)
		}
	})

	c := sqlite.NewConnector()
	db, err := c.Open(context.Background(), "~/test.db")

	require.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "sqlite: expand path ~/test.db")
}
