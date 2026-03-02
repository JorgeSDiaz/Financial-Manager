package sqlite_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainshared "github.com/financial-manager/api/internal/domain/shared"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
	transactionsqlite "github.com/financial-manager/api/internal/platform/transaction/sqlite"
)

func TestTransactionRepository_CreateAndGetByID(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	want := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	require.NoError(t, repo.Create(ctx, want))

	got, err := repo.GetByID(ctx, "tx-1")
	require.NoError(t, err)
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.AccountID, got.AccountID)
	assert.Equal(t, want.Type, got.Type)
	assert.InDelta(t, want.Amount, got.Amount, 0.001)
	assert.Equal(t, want.IsActive, got.IsActive)
}

func TestTransactionRepository_GetByID_NotFound(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))

	repo := transactionsqlite.NewTransactionRepository(db)

	_, err := repo.GetByID(context.Background(), "missing")
	require.Error(t, err)
	assert.ErrorIs(t, err, domainshared.ErrNotFound)
}

func TestTransactionRepository_Create_UpdatesAccountBalance(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	// Create income - should increase balance
	income := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	require.NoError(t, repo.Create(ctx, income))

	// Verify account balance increased
	var balance float64
	err := db.QueryRow("SELECT current_balance FROM accounts WHERE id = ?", "acc-001").Scan(&balance)
	require.NoError(t, err)
	assert.InDelta(t, 1100.0, balance, 0.001) // 1000 + 100

	// Create expense - should decrease balance
	expense := buildTestTransaction("tx-2", "acc-001", domaintransaction.TransactionTypeExpense, 50.0)
	require.NoError(t, repo.Create(ctx, expense))

	err = db.QueryRow("SELECT current_balance FROM accounts WHERE id = ?", "acc-001").Scan(&balance)
	require.NoError(t, err)
	assert.InDelta(t, 1050.0, balance, 0.001) // 1100 - 50
}

func TestTransactionRepository_ListByType_IncomeOnly(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	income1 := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	income2 := buildTestTransaction("tx-2", "acc-001", domaintransaction.TransactionTypeIncome, 200.0)
	expense := buildTestTransaction("tx-3", "acc-001", domaintransaction.TransactionTypeExpense, 50.0)

	require.NoError(t, repo.Create(ctx, income1))
	require.NoError(t, repo.Create(ctx, income2))
	require.NoError(t, repo.Create(ctx, expense))

	incomes, err := repo.ListByType(ctx, domaintransaction.TransactionTypeIncome, "", "", "", "")
	require.NoError(t, err)
	require.Len(t, incomes, 2)
	assert.Equal(t, domaintransaction.TransactionTypeIncome, incomes[0].Type)
	assert.Equal(t, domaintransaction.TransactionTypeIncome, incomes[1].Type)
}

func TestTransactionRepository_ListByType_ExpenseOnly(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	income := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	expense1 := buildTestTransaction("tx-2", "acc-001", domaintransaction.TransactionTypeExpense, 50.0)
	expense2 := buildTestTransaction("tx-3", "acc-001", domaintransaction.TransactionTypeExpense, 75.0)

	require.NoError(t, repo.Create(ctx, income))
	require.NoError(t, repo.Create(ctx, expense1))
	require.NoError(t, repo.Create(ctx, expense2))

	expenses, err := repo.ListByType(ctx, domaintransaction.TransactionTypeExpense, "", "", "", "")
	require.NoError(t, err)
	require.Len(t, expenses, 2)
	assert.Equal(t, domaintransaction.TransactionTypeExpense, expenses[0].Type)
	assert.Equal(t, domaintransaction.TransactionTypeExpense, expenses[1].Type)
}

func TestTransactionRepository_ListByType_WithFilters(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestAccount(db, "acc-002"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	// Create transactions for different accounts
	tx1 := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	tx2 := buildTestTransaction("tx-2", "acc-002", domaintransaction.TransactionTypeIncome, 200.0)

	require.NoError(t, repo.Create(ctx, tx1))
	require.NoError(t, repo.Create(ctx, tx2))

	// Filter by account_id
	filtered, err := repo.ListByType(ctx, domaintransaction.TransactionTypeIncome, "acc-001", "", "", "")
	require.NoError(t, err)
	require.Len(t, filtered, 1)
	assert.Equal(t, "acc-001", filtered[0].AccountID)
}

func TestTransactionRepository_ListByType_OnlyActive(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	tx := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	require.NoError(t, repo.Create(ctx, tx))

	// Soft delete the transaction
	require.NoError(t, repo.SoftDelete(ctx, "tx-1"))

	// Should not appear in list
	transactions, err := repo.ListByType(ctx, domaintransaction.TransactionTypeIncome, "", "", "", "")
	require.NoError(t, err)
	assert.Empty(t, transactions)
}

func TestTransactionRepository_Update_ModifiesFields(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	original := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	require.NoError(t, repo.Create(ctx, original))

	updated := original
	updated.Description = "Updated description"
	updated.Amount = 200.0
	updated.UpdatedAt = time.Now().UTC().Truncate(time.Second).Add(time.Minute)

	require.NoError(t, repo.Update(ctx, updated))

	got, err := repo.GetByID(ctx, "tx-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated description", got.Description)
	assert.InDelta(t, 200.0, got.Amount, 0.001)
}

func TestTransactionRepository_SoftDelete_RevertsBalance(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	// Create income - increases balance
	income := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	require.NoError(t, repo.Create(ctx, income))

	// Verify balance increased
	var balance float64
	err := db.QueryRow("SELECT current_balance FROM accounts WHERE id = ?", "acc-001").Scan(&balance)
	require.NoError(t, err)
	assert.InDelta(t, 1100.0, balance, 0.001)

	// Soft delete - should revert balance
	require.NoError(t, repo.SoftDelete(ctx, "tx-1"))

	err = db.QueryRow("SELECT current_balance FROM accounts WHERE id = ?", "acc-001").Scan(&balance)
	require.NoError(t, err)
	assert.InDelta(t, 1000.0, balance, 0.001) // Reverted to original

	// Transaction should be marked as inactive
	_, err = repo.GetByID(ctx, "tx-1")
	assert.ErrorIs(t, err, domainshared.ErrNotFound)
}

func TestTransactionRepository_SoftDelete_ExpenseRevertsBalance(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	// Create expense - decreases balance
	expense := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeExpense, 100.0)
	require.NoError(t, repo.Create(ctx, expense))

	// Verify balance decreased
	var balance float64
	err := db.QueryRow("SELECT current_balance FROM accounts WHERE id = ?", "acc-001").Scan(&balance)
	require.NoError(t, err)
	assert.InDelta(t, 900.0, balance, 0.001)

	// Soft delete - should revert balance (add back)
	require.NoError(t, repo.SoftDelete(ctx, "tx-1"))

	err = db.QueryRow("SELECT current_balance FROM accounts WHERE id = ?", "acc-001").Scan(&balance)
	require.NoError(t, err)
	assert.InDelta(t, 1000.0, balance, 0.001) // Reverted to original
}

func TestTransactionRepository_SoftDelete_NotFound(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))

	repo := transactionsqlite.NewTransactionRepository(db)

	err := repo.SoftDelete(context.Background(), "missing")
	require.Error(t, err)
	assert.ErrorIs(t, err, domainshared.ErrNotFound)
}

func TestTransactionRepository_ListRecent_ReturnsLimitedResults(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	// Create multiple transactions with different dates
	oldTx := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	oldTx.Date = time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour)
	require.NoError(t, repo.Create(ctx, oldTx))

	newTx := buildTestTransaction("tx-2", "acc-001", domaintransaction.TransactionTypeIncome, 200.0)
	require.NoError(t, repo.Create(ctx, newTx))

	// ListRecent should return most recent first
	transactions, err := repo.ListRecent(ctx, 1)
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	assert.Equal(t, "tx-2", transactions[0].ID)
}

func TestTransactionRepository_ListRecent_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)

	transactions, err := repo.ListRecent(context.Background(), 10)
	require.NoError(t, err)
	require.Empty(t, transactions)
}

func TestTransactionRepository_ListRecent_ReturnsOnlyActive(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	tx := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	require.NoError(t, repo.Create(ctx, tx))

	// Soft delete
	require.NoError(t, repo.SoftDelete(ctx, "tx-1"))

	// Should not appear in list
	transactions, err := repo.ListRecent(ctx, 10)
	require.NoError(t, err)
	assert.Empty(t, transactions)
}

func TestTransactionRepository_ListRecent_QueryError(t *testing.T) {
	t.Parallel()
	db, err := sql.Open("sqlite", "file:?mode=invalid")
	require.NoError(t, err)
	defer db.Close()

	repo := transactionsqlite.NewTransactionRepository(db)
	_, err = repo.ListRecent(context.Background(), 10)
	require.Error(t, err)
}

func TestTransactionRepository_ListByType_WithDateRange(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	oldTx := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	oldTx.Date = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	require.NoError(t, repo.Create(ctx, oldTx))

	newTx := buildTestTransaction("tx-2", "acc-001", domaintransaction.TransactionTypeIncome, 200.0)
	newTx.Date = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	require.NoError(t, repo.Create(ctx, newTx))

	// Filter by date range
	transactions, err := repo.ListByType(ctx, domaintransaction.TransactionTypeIncome, "", "", "2026-01-01", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	assert.Equal(t, "tx-2", transactions[0].ID)
}

func TestTransactionRepository_ListByType_WithCategoryFilter(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	require.NoError(t, buildTestAccount(db, "acc-001"))
	require.NoError(t, buildTestCategory(db, "cat-001"))
	require.NoError(t, buildTestCategory(db, "cat-002"))

	repo := transactionsqlite.NewTransactionRepository(db)
	ctx := context.Background()

	tx1 := buildTestTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome, 100.0)
	tx1.CategoryID = "cat-001"
	require.NoError(t, repo.Create(ctx, tx1))

	tx2 := buildTestTransaction("tx-2", "acc-001", domaintransaction.TransactionTypeIncome, 200.0)
	tx2.CategoryID = "cat-002"
	require.NoError(t, repo.Create(ctx, tx2))

	// Filter by category
	transactions, err := repo.ListByType(ctx, domaintransaction.TransactionTypeIncome, "", "cat-001", "", "")
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	assert.Equal(t, "cat-001", transactions[0].CategoryID)
}

func TestTransactionRepository_ListByType_QueryError(t *testing.T) {
	t.Parallel()
	db, err := sql.Open("sqlite", "file:?mode=invalid")
	require.NoError(t, err)
	defer db.Close()

	repo := transactionsqlite.NewTransactionRepository(db)
	_, err = repo.ListByType(context.Background(), domaintransaction.TransactionTypeIncome, "", "", "", "")
	require.Error(t, err)
}
