// Package sqlite implements the dashboard repository using SQLite.
package sqlite

import (
	"context"
	"database/sql"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// DashboardRepository implements the dashboard repository interface using SQLite.
type DashboardRepository struct {
	accountsDB     *sql.DB
	transactionsDB *sql.DB
	categoriesDB   *sql.DB
}

// NewDashboardRepository creates a DashboardRepository with the provided databases.
func NewDashboardRepository(accountsDB, transactionsDB, categoriesDB *sql.DB) *DashboardRepository {
	return &DashboardRepository{
		accountsDB:     accountsDB,
		transactionsDB: transactionsDB,
		categoriesDB:   categoriesDB,
	}
}

// ListAccounts returns all active accounts.
func (r *DashboardRepository) ListAccounts(ctx context.Context) ([]domainaccount.Account, error) {
	const q = `SELECT id, name, type, initial_balance, current_balance, currency, is_active, created_at, updated_at
		FROM accounts WHERE is_active = 1 ORDER BY name`

	rows, err := r.accountsDB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]domainaccount.Account, 0)
	for rows.Next() {
		var a domainaccount.Account
		var isActive int
		var createdAt, updatedAt string
		err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.InitialBalance, &a.CurrentBalance, &a.Currency, &isActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		a.IsActive = isActive == 1
		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// ListRecentTransactions returns the most recent transactions up to the limit.
func (r *DashboardRepository) ListRecentTransactions(ctx context.Context, limit int) ([]domaintransaction.Transaction, error) {
	const q = `SELECT id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at
		FROM transactions WHERE is_active = 1 ORDER BY date DESC LIMIT ?`

	rows, err := r.transactionsDB.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]domaintransaction.Transaction, 0)
	for rows.Next() {
		var t domaintransaction.Transaction
		var tType string
		var isActive int
		var date, createdAt, updatedAt string
		err := rows.Scan(&t.ID, &t.AccountID, &t.CategoryID, &tType, &t.Amount, &t.Description, &date, &isActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		t.Type = domaintransaction.TransactionType(tType)
		t.IsActive = isActive == 1
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// ListExpenseTransactions returns expense transactions filtered by optional criteria.
func (r *DashboardRepository) ListExpenseTransactions(ctx context.Context, accountID, categoryID, startDate, endDate string) ([]domaintransaction.Transaction, error) {
	return r.listByType(ctx, domaintransaction.TransactionTypeExpense, accountID, categoryID, startDate, endDate)
}

// ListIncomeTransactions returns income transactions filtered by optional criteria.
func (r *DashboardRepository) ListIncomeTransactions(ctx context.Context, accountID, categoryID, startDate, endDate string) ([]domaintransaction.Transaction, error) {
	return r.listByType(ctx, domaintransaction.TransactionTypeIncome, accountID, categoryID, startDate, endDate)
}

// listByType returns transactions filtered by type and optional criteria.
func (r *DashboardRepository) listByType(ctx context.Context, tType domaintransaction.TransactionType, accountID, categoryID, startDate, endDate string) ([]domaintransaction.Transaction, error) {
	q := `SELECT id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at
		FROM transactions WHERE type = ? AND is_active = 1`
	args := []interface{}{string(tType)}

	if accountID != "" {
		q += " AND account_id = ?"
		args = append(args, accountID)
	}
	if categoryID != "" {
		q += " AND category_id = ?"
		args = append(args, categoryID)
	}
	if startDate != "" {
		q += " AND date >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		q += " AND date <= ?"
		args = append(args, endDate)
	}
	q += " ORDER BY date DESC"

	rows, err := r.transactionsDB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]domaintransaction.Transaction, 0)
	for rows.Next() {
		var t domaintransaction.Transaction
		var tType string
		var isActive int
		var date, createdAt, updatedAt string
		err := rows.Scan(&t.ID, &t.AccountID, &t.CategoryID, &tType, &t.Amount, &t.Description, &date, &isActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		t.Type = domaintransaction.TransactionType(tType)
		t.IsActive = isActive == 1
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// ListCategories returns all active categories.
func (r *DashboardRepository) ListCategories(ctx context.Context) ([]domaincategory.Category, error) {
	const q = `SELECT id, name, type, color, icon, is_system, is_active, created_at, updated_at
		FROM categories WHERE is_active = 1 ORDER BY name`

	rows, err := r.categoriesDB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]domaincategory.Category, 0)
	for rows.Next() {
		var c domaincategory.Category
		var isSystem, isActive int
		var createdAt, updatedAt string
		err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.Color, &c.Icon, &isSystem, &isActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		c.IsSystem = isSystem == 1
		c.IsActive = isActive == 1
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
