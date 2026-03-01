// Package sqlite implements the export repository using SQLite.
package sqlite

import (
	"context"
	"database/sql"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// ExportRepository implements the export repository interface using SQLite.
type ExportRepository struct {
	accountsDB     *sql.DB
	categoriesDB   *sql.DB
	transactionsDB *sql.DB
}

// NewExportRepository creates an ExportRepository with the provided databases.
func NewExportRepository(accountsDB, categoriesDB, transactionsDB *sql.DB) *ExportRepository {
	return &ExportRepository{
		accountsDB:     accountsDB,
		categoriesDB:   categoriesDB,
		transactionsDB: transactionsDB,
	}
}

// ListAccounts returns all active accounts.
func (r *ExportRepository) ListAccounts(ctx context.Context) ([]domainaccount.Account, error) {
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

// ListCategories returns all active categories.
func (r *ExportRepository) ListCategories(ctx context.Context) ([]domaincategory.Category, error) {
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

// ListTransactions returns transactions filtered by type and optional date range.
func (r *ExportRepository) ListTransactions(ctx context.Context, tType domaintransaction.TransactionType, startDate, endDate string) ([]domaintransaction.Transaction, error) {
	q := `SELECT id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at
		FROM transactions WHERE is_active = 1`
	args := []interface{}{}

	if tType != "" {
		q += " AND type = ?"
		args = append(args, string(tType))
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
		var tTypeStr string
		var isActive int
		var date, createdAt, updatedAt string
		err := rows.Scan(&t.ID, &t.AccountID, &t.CategoryID, &tTypeStr, &t.Amount, &t.Description, &date, &isActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		t.Type = domaintransaction.TransactionType(tTypeStr)
		t.IsActive = isActive == 1
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
