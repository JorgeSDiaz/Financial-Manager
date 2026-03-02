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
	dashboardsqlite "github.com/financial-manager/api/internal/platform/dashboard/sqlite"
)

var dashboardDbCounter atomic.Int64

func newDashboardTestDB(t *testing.T, schema string) *sql.DB {
	t.Helper()
	name := fmt.Sprintf("test%d", dashboardDbCounter.Add(1))
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", name))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec(schema)
	require.NoError(t, err)

	return db
}

func TestDashboardRepository_ListAccounts_ReturnsOnlyActive(t *testing.T) {
	t.Parallel()
	accountsDB := newDashboardTestDB(t, accountsSchema)
	now := time.Now().UTC().Truncate(time.Second)

	_, _ = accountsDB.Exec(`INSERT INTO accounts (id, name, type, initial_balance, current_balance, currency, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"a1", "Active", "cash", 100.0, 100.0, "USD", 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = accountsDB.Exec(`INSERT INTO accounts (id, name, type, initial_balance, current_balance, currency, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"a2", "Inactive", "cash", 100.0, 100.0, "USD", 0, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := dashboardsqlite.NewDashboardRepository(accountsDB, transactionsDBForDashboardTest(t), categoriesDBForDashboardTest(t))
	accounts, err := repo.ListAccounts(context.Background())
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.Equal(t, "a1", accounts[0].ID)
}

func TestDashboardRepository_ListAccounts_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDBForDashboardTest(t), categoriesDBForDashboardTest(t))

	accounts, err := repo.ListAccounts(context.Background())
	require.NoError(t, err)
	require.Empty(t, accounts)
}

func TestDashboardRepository_ListRecentTransactions_ReturnsLimitedResults(t *testing.T) {
	t.Parallel()
	transactionsDB := newDashboardTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "income", 100.0, "Test1", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a1", "c1", "expense", 50.0, "Test2", now.Add(-time.Hour).Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t3", "a1", "c1", "income", 75.0, "Test3", now.Add(-2*time.Hour).Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDB, categoriesDBForDashboardTest(t))
	transactions, err := repo.ListRecentTransactions(context.Background(), 2)
	require.NoError(t, err)
	require.Len(t, transactions, 2)
	require.Equal(t, "t1", transactions[0].ID)
}

func TestDashboardRepository_ListRecentTransactions_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDBForDashboardTest(t), categoriesDBForDashboardTest(t))

	transactions, err := repo.ListRecentTransactions(context.Background(), 10)
	require.NoError(t, err)
	require.Empty(t, transactions)
}

func TestDashboardRepository_ListRecentTransactions_ReturnsOnlyActive(t *testing.T) {
	t.Parallel()
	transactionsDB := newDashboardTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "income", 100.0, "Active", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a1", "c1", "income", 50.0, "Inactive", now.Format(time.RFC3339), 0, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDB, categoriesDBForDashboardTest(t))
	transactions, err := repo.ListRecentTransactions(context.Background(), 10)
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, "t1", transactions[0].ID)
}

func TestDashboardRepository_ListExpenseTransactions_ReturnsFilteredByType(t *testing.T) {
	t.Parallel()
	transactionsDB := newDashboardTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "income", 100.0, "Income", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a1", "c1", "expense", 50.0, "Expense", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDB, categoriesDBForDashboardTest(t))
	transactions, err := repo.ListExpenseTransactions(context.Background(), "", "", "", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, domaintransaction.TransactionTypeExpense, transactions[0].Type)
}

func TestDashboardRepository_ListExpenseTransactions_WithAccountFilter(t *testing.T) {
	t.Parallel()
	transactionsDB := newDashboardTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "expense", 100.0, "Account1", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a2", "c1", "expense", 50.0, "Account2", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDB, categoriesDBForDashboardTest(t))
	transactions, err := repo.ListExpenseTransactions(context.Background(), "a1", "", "", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, "a1", transactions[0].AccountID)
}

func TestDashboardRepository_ListExpenseTransactions_WithDateRange(t *testing.T) {
	t.Parallel()
	transactionsDB := newDashboardTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "expense", 100.0, "Old", "2025-01-01T00:00:00Z", 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a1", "c1", "expense", 50.0, "New", "2026-01-01T00:00:00Z", 1, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDB, categoriesDBForDashboardTest(t))
	transactions, err := repo.ListExpenseTransactions(context.Background(), "", "", "2026-01-01", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, "t2", transactions[0].ID)
}

func TestDashboardRepository_ListExpenseTransactions_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDBForDashboardTest(t), categoriesDBForDashboardTest(t))

	transactions, err := repo.ListExpenseTransactions(context.Background(), "", "", "", "")
	require.NoError(t, err)
	require.Empty(t, transactions)
}

func TestDashboardRepository_ListIncomeTransactions_ReturnsFilteredByType(t *testing.T) {
	t.Parallel()
	transactionsDB := newDashboardTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "income", 100.0, "Income", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a1", "c1", "expense", 50.0, "Expense", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDB, categoriesDBForDashboardTest(t))
	transactions, err := repo.ListIncomeTransactions(context.Background(), "", "", "", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, domaintransaction.TransactionTypeIncome, transactions[0].Type)
}

func TestDashboardRepository_ListIncomeTransactions_WithCategoryFilter(t *testing.T) {
	t.Parallel()
	transactionsDB := newDashboardTestDB(t, transactionsSchema)
	now := time.Now().UTC().Truncate(time.Second)
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t1", "a1", "c1", "income", 100.0, "Cat1", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = transactionsDB.Exec(`INSERT INTO transactions (id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"t2", "a1", "c2", "income", 50.0, "Cat2", now.Format(time.RFC3339), 1, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDB, categoriesDBForDashboardTest(t))
	transactions, err := repo.ListIncomeTransactions(context.Background(), "", "c1", "", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, "c1", transactions[0].CategoryID)
}

func TestDashboardRepository_ListIncomeTransactions_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDBForDashboardTest(t), categoriesDBForDashboardTest(t))

	transactions, err := repo.ListIncomeTransactions(context.Background(), "", "", "", "")
	require.NoError(t, err)
	require.Empty(t, transactions)
}

func TestDashboardRepository_ListCategories_ReturnsOnlyActive(t *testing.T) {
	t.Parallel()
	categoriesDB := newDashboardTestDB(t, categoriesSchema)
	now := time.Now().UTC().Truncate(time.Second)

	_, _ = categoriesDB.Exec(`INSERT INTO categories (id, name, type, color, icon, is_system, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"c1", "Active", "expense", "#fff", "icon", 0, 1, now.Format(time.RFC3339), now.Format(time.RFC3339))
	_, _ = categoriesDB.Exec(`INSERT INTO categories (id, name, type, color, icon, is_system, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"c2", "Inactive", "expense", "#fff", "icon", 0, 0, now.Format(time.RFC3339), now.Format(time.RFC3339))

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDBForDashboardTest(t), categoriesDB)
	categories, err := repo.ListCategories(context.Background())
	require.NoError(t, err)
	require.Len(t, categories, 1)
	require.Equal(t, "c1", categories[0].ID)
}

func TestDashboardRepository_ListCategories_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDBForDashboardTest(t), categoriesDBForDashboardTest(t))

	categories, err := repo.ListCategories(context.Background())
	require.NoError(t, err)
	require.Empty(t, categories)
}

func TestDashboardRepository_ListAccounts_QueryError(t *testing.T) {
	t.Parallel()
	accountsDB, err := sql.Open("sqlite", "file:?mode=invalid")
	require.NoError(t, err)
	defer accountsDB.Close()

	repo := dashboardsqlite.NewDashboardRepository(accountsDB, transactionsDBForDashboardTest(t), categoriesDBForDashboardTest(t))
	_, err = repo.ListAccounts(context.Background())
	require.Error(t, err)
}

func TestDashboardRepository_ListRecentTransactions_QueryError(t *testing.T) {
	t.Parallel()
	transactionsDB, err := sql.Open("sqlite", "file:?mode=invalid")
	require.NoError(t, err)
	defer transactionsDB.Close()

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDB, categoriesDBForDashboardTest(t))
	_, err = repo.ListRecentTransactions(context.Background(), 10)
	require.Error(t, err)
}

func TestDashboardRepository_ListCategories_QueryError(t *testing.T) {
	t.Parallel()
	categoriesDB, err := sql.Open("sqlite", "file:?mode=invalid")
	require.NoError(t, err)
	defer categoriesDB.Close()

	repo := dashboardsqlite.NewDashboardRepository(accountsDBForDashboardTest(t), transactionsDBForDashboardTest(t), categoriesDB)
	_, err = repo.ListCategories(context.Background())
	require.Error(t, err)
}

func accountsDBForDashboardTest(t *testing.T) *sql.DB {
	return newDashboardTestDB(t, accountsSchema)
}

func categoriesDBForDashboardTest(t *testing.T) *sql.DB {
	return newDashboardTestDB(t, categoriesSchema)
}

func transactionsDBForDashboardTest(t *testing.T) *sql.DB {
	return newDashboardTestDB(t, transactionsSchema)
}

func buildTestAccountForDashboard(id, name string) domainaccount.Account {
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

func buildTestCategoryForDashboard(id, name string, tType domaincategory.Type) domaincategory.Category {
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
