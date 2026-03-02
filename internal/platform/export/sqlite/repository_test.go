package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/require"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
	exportsqlite "github.com/financial-manager/api/internal/platform/export/sqlite"
)

var exportDbCounter atomic.Int64

func newExportTestDB(t *testing.T, schema string) *sql.DB {
	t.Helper()
	name := fmt.Sprintf("test%d", exportDbCounter.Add(1))
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", name))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec(schema)
	require.NoError(t, err)

	return db
}

func TestExportRepository_ListAccounts_ReturnsOnlyActive(t *testing.T) {
	t.Parallel()
	accountsDB := newExportTestDB(t, accountsSchema)
	now := time.Now().UTC().Truncate(time.Second)

	_, _ = accountsDB.Exec(`INSERT INTO accounts (id, name, type, initial_balance, current_balance, currency, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"a1", "Active", "cash", 100.0, 100.0, "USD", 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = accountsDB.Exec(`INSERT INTO accounts (id, name, type, initial_balance, current_balance, currency, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"a2", "Inactive", "cash", 100.0, 100.0, "USD", 0, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := exportsqlite.NewExportRepository(accountsDB, categoriesDBForTest(t), transactionsDBForTest(t))
	accounts, err := repo.ListAccounts(context.Background())
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.Equal(t, "a1", accounts[0].ID)
}

func TestExportRepository_ListAccounts_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDBForTest(t), transactionsDBForTest(t))

	accounts, err := repo.ListAccounts(context.Background())
	require.NoError(t, err)
	require.Empty(t, accounts)
}

func TestExportRepository_ListCategories_ReturnsOnlyActive(t *testing.T) {
	t.Parallel()
	categoriesDB := newExportTestDB(t, categoriesSchema)
	now := time.Now().UTC().Truncate(time.Second)

	_, _ = categoriesDB.Exec(`INSERT INTO categories (id, name, type, color, icon, is_system, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"c1", "Active", "expense", "#fff", "icon", 0, 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = categoriesDB.Exec(`INSERT INTO categories (id, name, type, color, icon, is_system, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"c2", "Inactive", "expense", "#fff", "icon", 0, 0, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDB, transactionsDBForTest(t))
	categories, err := repo.ListCategories(context.Background())
	require.NoError(t, err)
	require.Len(t, categories, 1)
	require.Equal(t, "c1", categories[0].ID)
}

func TestExportRepository_ListCategories_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDBForTest(t), transactionsDBForTest(t))

	categories, err := repo.ListCategories(context.Background())
	require.NoError(t, err)
	require.Empty(t, categories)
}

func TestExportRepository_ListTransactions_WithoutFilters(t *testing.T) {
	t.Parallel()
	transactionsDB := newExportTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "income", 100.0, "Test", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDBForTest(t), transactionsDB)
	transactions, err := repo.ListTransactions(context.Background(), "", "", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, "t1", transactions[0].ID)
}

func TestExportRepository_ListTransactions_WithTypeFilter(t *testing.T) {
	t.Parallel()
	transactionsDB := newExportTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "income", 100.0, "Income", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a1", "c1", "expense", 50.0, "Expense", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDBForTest(t), transactionsDB)
	transactions, err := repo.ListTransactions(context.Background(), domaintransaction.TransactionTypeIncome, "", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, domaintransaction.TransactionTypeIncome, transactions[0].Type)
}

func TestExportRepository_ListTransactions_WithDateRange(t *testing.T) {
	t.Parallel()
	transactionsDB := newExportTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "income", 100.0, "Old", "2025-01-01T00:00:00Z", 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a1", "c1", "income", 200.0, "New", "2026-01-01T00:00:00Z", 1, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDBForTest(t), transactionsDB)
	transactions, err := repo.ListTransactions(context.Background(), "", "2026-01-01", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, "t2", transactions[0].ID)
}

func TestExportRepository_ListTransactions_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDBForTest(t), transactionsDBForTest(t))

	transactions, err := repo.ListTransactions(context.Background(), "", "", "")
	require.NoError(t, err)
	require.Empty(t, transactions)
}

func TestExportRepository_ListTransactions_ReturnsOnlyActive(t *testing.T) {
	t.Parallel()
	transactionsDB := newExportTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "income", 100.0, "Active", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a1", "c1", "income", 50.0, "Inactive", now.Format(time.RFC3339), 0, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDBForTest(t), transactionsDB)
	transactions, err := repo.ListTransactions(context.Background(), "", "", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, "t1", transactions[0].ID)
}

func TestExportRepository_ListAccounts_QueryError(t *testing.T) {
	t.Parallel()
	accountsDB, err := sql.Open("sqlite", "file:?mode=invalid")
	require.NoError(t, err)
	defer accountsDB.Close()

	repo := exportsqlite.NewExportRepository(accountsDB, categoriesDBForTest(t), transactionsDBForTest(t))
	_, err = repo.ListAccounts(context.Background())
	require.Error(t, err)
}

func TestExportRepository_ListCategories_QueryError(t *testing.T) {
	t.Parallel()
	categoriesDB, err := sql.Open("sqlite", "file:?mode=invalid")
	require.NoError(t, err)
	defer categoriesDB.Close()

	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDB, transactionsDBForTest(t))
	_, err = repo.ListCategories(context.Background())
	require.Error(t, err)
}

func TestExportRepository_ListTransactions_QueryError(t *testing.T) {
	t.Parallel()
	transactionsDB, err := sql.Open("sqlite", "file:?mode=invalid")
	require.NoError(t, err)
	defer transactionsDB.Close()

	repo := exportsqlite.NewExportRepository(accountsDBForTest(t), categoriesDBForTest(t), transactionsDB)
	_, err = repo.ListTransactions(context.Background(), "", "", "")
	require.Error(t, err)
}

func accountsDBForTest(t *testing.T) *sql.DB {
	return newExportTestDB(t, accountsSchema)
}

func categoriesDBForTest(t *testing.T) *sql.DB {
	return newExportTestDB(t, categoriesSchema)
}

func transactionsDBForTest(t *testing.T) *sql.DB {
	return newExportTestDB(t, transactionsSchema)
}

const accountsSchema = `CREATE TABLE IF NOT EXISTS accounts (
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
)`

const categoriesSchema = `CREATE TABLE IF NOT EXISTS categories (
	id          TEXT PRIMARY KEY,
	name        TEXT NOT NULL,
	type        TEXT NOT NULL CHECK(type IN ('expense', 'income')),
	color       TEXT NOT NULL DEFAULT '',
	icon        TEXT NOT NULL DEFAULT '',
	is_system   INTEGER NOT NULL DEFAULT 0,
	is_active   INTEGER NOT NULL DEFAULT 1,
	created_at  TEXT NOT NULL,
	updated_at  TEXT NOT NULL
)`

const transactionsSchema = `CREATE TABLE IF NOT EXISTS transactions (
	id          TEXT PRIMARY KEY,
	account_id  TEXT NOT NULL,
	category_id TEXT NOT NULL,
	type        TEXT NOT NULL CHECK(type IN ('income', 'expense')),
	amount      REAL NOT NULL,
	description TEXT NOT NULL,
	date        TEXT NOT NULL,
	is_active   INTEGER NOT NULL DEFAULT 1,
	created_at  TEXT NOT NULL,
	updated_at  TEXT NOT NULL
)`

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

func buildTestCategory(id, name string, tType domaincategory.Type) domaincategory.Category {
	now := time.Now().UTC().Truncate(time.Second)
	return domaincategory.Category{
		ID:        id,
		Name:      name,
		Type:      tType,
		Color:     "#FFFFFF",
		Icon:      "icon",
		IsSystem:  false,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
