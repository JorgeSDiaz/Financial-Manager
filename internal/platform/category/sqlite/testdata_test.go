package sqlite_test

import (
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/require"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// dbCounter provides unique in-memory database names to avoid shared state between parallel tests.
var dbCounter atomic.Int64

// newTestDB creates an isolated in-memory SQLite database with the categories schema applied.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	name := fmt.Sprintf("test%d", dbCounter.Add(1))
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", name))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS categories (
		id          TEXT PRIMARY KEY,
		name        TEXT NOT NULL,
		type        TEXT NOT NULL CHECK(type IN ('expense', 'income')),
		color       TEXT NOT NULL,
		icon        TEXT NOT NULL,
		is_system   INTEGER NOT NULL DEFAULT 0,
		is_active   INTEGER NOT NULL DEFAULT 1,
		created_at  TEXT NOT NULL,
		updated_at  TEXT NOT NULL
	)`)
	require.NoError(t, err)

	return db
}

// buildTestCategory returns a valid Category fixture for use in repository tests.
func buildTestCategory(id, name string) domaincategory.Category {
	now := time.Now().UTC().Truncate(time.Second)
	return domaincategory.Category{
		ID:        id,
		Name:      name,
		Type:      domaincategory.TypeExpense,
		Color:     "#FFFFFF",
		Icon:      "wallet",
		IsSystem:  false,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
