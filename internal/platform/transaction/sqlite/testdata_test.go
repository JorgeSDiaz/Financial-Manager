package sqlite_test

import (
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/require"

	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// dbCounter provides unique in-memory database names to avoid shared state between parallel tests.
var dbCounter atomic.Int64

// newTestDB creates an isolated in-memory SQLite database with the transactions schema applied.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	name := fmt.Sprintf("test%d", dbCounter.Add(1))
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", name))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	// Create accounts table for FK constraint
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		id              TEXT    PRIMARY KEY,
		name            TEXT    NOT NULL,
		type            TEXT    NOT NULL,
		initial_balance REAL    NOT NULL DEFAULT 0,
		current_balance REAL    NOT NULL DEFAULT 0,
		currency        TEXT    NOT NULL DEFAULT 'USD',
		color           TEXT    NOT NULL DEFAULT '',
		icon            TEXT    NOT NULL DEFAULT '',
		is_active       INTEGER NOT NULL DEFAULT 1,
		created_at      TEXT    NOT NULL,
		updated_at      TEXT    NOT NULL
	)`)
	require.NoError(t, err)

	// Create categories table for FK constraint
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS categories (
		id          TEXT PRIMARY KEY,
		name        TEXT NOT NULL,
		type        TEXT NOT NULL,
		color       TEXT NOT NULL DEFAULT '',
		icon        TEXT NOT NULL DEFAULT '',
		is_system   INTEGER NOT NULL DEFAULT 0,
		is_active   INTEGER NOT NULL DEFAULT 1,
		created_at  TEXT NOT NULL,
		updated_at  TEXT NOT NULL
	)`)
	require.NoError(t, err)

	// Create transactions table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS transactions (
		id          TEXT PRIMARY KEY,
		account_id  TEXT NOT NULL,
		category_id TEXT,
		type        TEXT NOT NULL CHECK(type IN ('income', 'expense')),
		amount      REAL NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		date        TEXT NOT NULL,
		is_active   INTEGER NOT NULL DEFAULT 1,
		created_at  TEXT NOT NULL,
		updated_at  TEXT NOT NULL,
		FOREIGN KEY (account_id) REFERENCES accounts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	)`)
	require.NoError(t, err)

	return db
}

// buildTestTransaction returns a valid Transaction fixture for use in repository tests.
func buildTestTransaction(id, accountID string, tType domaintransaction.TransactionType, amount float64) domaintransaction.Transaction {
	now := time.Now().UTC().Truncate(time.Second)
	date := now.Truncate(24 * time.Hour)
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   accountID,
		CategoryID:  "cat-001",
		Type:        tType,
		Amount:      amount,
		Description: "Test transaction",
		Date:        date,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// buildTestAccount creates a test account in the database.
func buildTestAccount(db *sql.DB, id string) error {
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	_, err := db.Exec(`INSERT INTO accounts (id, name, type, initial_balance, current_balance, currency, color, icon, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, "Test Account", "cash", 1000.0, 1000.0, "USD", "#FFFFFF", "wallet", 1, now, now)
	return err
}

// buildTestCategory creates a test category in the database.
func buildTestCategory(db *sql.DB, id string) error {
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	_, err := db.Exec(`INSERT INTO categories (id, name, type, color, icon, is_system, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, "Test Category", "expense", "#FF0000", "tag", 0, 1, now, now)
	return err
}
