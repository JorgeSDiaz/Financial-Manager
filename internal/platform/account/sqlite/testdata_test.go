package sqlite_test

import (
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/require"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// dbCounter provides unique in-memory database names to avoid shared state between parallel tests.
var dbCounter atomic.Int64

// newTestDB creates an isolated in-memory SQLite database with the accounts schema applied.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	name := fmt.Sprintf("test%d", dbCounter.Add(1))
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", name))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		id              TEXT    PRIMARY KEY,
		name            TEXT    NOT NULL,
		type            TEXT    NOT NULL CHECK(type IN ('cash', 'bank', 'credit_card', 'savings')),
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

	return db
}

// buildTestAccount returns a valid Account fixture for use in repository tests.
func buildTestAccount(id, name string) domainaccount.Account {
	now := time.Now().UTC().Truncate(time.Second)
	return domainaccount.Account{
		ID:             id,
		Name:           name,
		Type:           domainaccount.AccountTypeCash,
		InitialBalance: 100.0,
		CurrentBalance: 100.0,
		Currency:       "USD",
		Color:          "#FFFFFF",
		Icon:           "wallet",
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
