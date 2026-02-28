// Package sqlite implements the TransactionRepository using SQLite.
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	domainshared "github.com/financial-manager/api/internal/domain/shared"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const timeLayout = "2006-01-02T15:04:05Z"
const dateLayout = "2006-01-02"

// TransactionRepository implements transaction repository interfaces using SQLite.
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a TransactionRepository with the provided *sql.DB.
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create inserts a new transaction row and updates the account balance.
func (r *TransactionRepository) Create(ctx context.Context, t domaintransaction.Transaction) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction sqlite: begin: %w", err)
	}
	defer tx.Rollback()

	const insertQ = `INSERT INTO transactions
		(id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	active := 0
	if t.IsActive {
		active = 1
	}

	_, err = tx.ExecContext(ctx, insertQ,
		t.ID, t.AccountID, t.CategoryID, string(t.Type),
		t.Amount, t.Description,
		t.Date.Format(dateLayout),
		active,
		t.CreatedAt.UTC().Format(timeLayout),
		t.UpdatedAt.UTC().Format(timeLayout),
	)
	if err != nil {
		return fmt.Errorf("transaction sqlite: create: %w", err)
	}

	// Update account balance
	var balanceDelta float64
	if t.Type == domaintransaction.TransactionTypeIncome {
		balanceDelta = t.Amount
	} else {
		balanceDelta = -t.Amount
	}

	const updateBalanceQ = `UPDATE accounts SET current_balance = current_balance + ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	_, err = tx.ExecContext(ctx, updateBalanceQ, balanceDelta, now.Format(timeLayout), t.AccountID)
	if err != nil {
		return fmt.Errorf("transaction sqlite: update account balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction sqlite: commit: %w", err)
	}

	return nil
}

// GetByID retrieves a transaction by its ID.
func (r *TransactionRepository) GetByID(ctx context.Context, id string) (domaintransaction.Transaction, error) {
	const q = `SELECT id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at
		FROM transactions WHERE id = ? AND is_active = 1`

	row := r.db.QueryRowContext(ctx, q, id)
	t, err := scanTransaction(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domaintransaction.Transaction{}, domainshared.ErrNotFound
	}
	if err != nil {
		return domaintransaction.Transaction{}, fmt.Errorf("transaction sqlite: get by id: %w", err)
	}

	return t, nil
}

// Update modifies an existing transaction.
func (r *TransactionRepository) Update(ctx context.Context, t domaintransaction.Transaction) error {
	const q = `UPDATE transactions SET 
		category_id = ?, amount = ?, description = ?, date = ?, updated_at = ? 
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, q,
		t.CategoryID, t.Amount, t.Description,
		t.Date.Format(dateLayout),
		t.UpdatedAt.UTC().Format(timeLayout),
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("transaction sqlite: update: %w", err)
	}

	return nil
}

// SoftDelete marks a transaction as inactive and reverts the account balance.
func (r *TransactionRepository) SoftDelete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction sqlite: begin: %w", err)
	}
	defer tx.Rollback()

	// Get transaction info before deleting
	const getQ = `SELECT account_id, type, amount FROM transactions WHERE id = ? AND is_active = 1`
	var accountID string
	var tType string
	var amount float64
	err = tx.QueryRowContext(ctx, getQ, id).Scan(&accountID, &tType, &amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domainshared.ErrNotFound
		}
		return fmt.Errorf("transaction sqlite: get for delete: %w", err)
	}

	// Soft delete transaction
	const deleteQ = `UPDATE transactions SET is_active = 0, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	_, err = tx.ExecContext(ctx, deleteQ, now.Format(timeLayout), id)
	if err != nil {
		return fmt.Errorf("transaction sqlite: soft delete: %w", err)
	}

	// Revert account balance
	var balanceDelta float64
	if tType == string(domaintransaction.TransactionTypeIncome) {
		balanceDelta = -amount // Revert income by subtracting
	} else {
		balanceDelta = amount // Revert expense by adding
	}

	const updateBalanceQ = `UPDATE accounts SET current_balance = current_balance + ?, updated_at = ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, updateBalanceQ, balanceDelta, now.Format(timeLayout), accountID)
	if err != nil {
		return fmt.Errorf("transaction sqlite: revert account balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction sqlite: commit: %w", err)
	}

	return nil
}

// ListByType returns transactions filtered by type, account, category, and date range.
func (r *TransactionRepository) ListByType(ctx context.Context, tType domaintransaction.TransactionType, accountID, categoryID, startDate, endDate string) ([]domaintransaction.Transaction, error) {
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "type = ?", "is_active = 1")
	args = append(args, string(tType))

	if accountID != "" {
		conditions = append(conditions, "account_id = ?")
		args = append(args, accountID)
	}
	if categoryID != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, categoryID)
	}
	if startDate != "" {
		conditions = append(conditions, "date >= ?")
		args = append(args, startDate)
	}
	if endDate != "" {
		conditions = append(conditions, "date <= ?")
		args = append(args, endDate)
	}

	q := fmt.Sprintf(`SELECT id, account_id, category_id, type, amount, description, date, is_active, created_at, updated_at
		FROM transactions WHERE %s ORDER BY date DESC`, strings.Join(conditions, " AND "))

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("transaction sqlite: list: %w", err)
	}
	defer rows.Close()

	transactions := make([]domaintransaction.Transaction, 0)
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("transaction sqlite: list scan: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("transaction sqlite: list rows: %w", err)
	}

	return transactions, nil
}

// scanner abstracts *sql.Row and *sql.Rows for the shared scanTransaction helper.
type scanner interface {
	Scan(dest ...any) error
}

func scanTransaction(s scanner) (domaintransaction.Transaction, error) {
	var (
		t                          domaintransaction.Transaction
		tType                      string
		isActive                   int
		date, createdAt, updatedAt string
	)

	err := s.Scan(
		&t.ID, &t.AccountID, &t.CategoryID,
		&tType, &t.Amount, &t.Description,
		&date, &isActive, &createdAt, &updatedAt,
	)
	if err != nil {
		return domaintransaction.Transaction{}, err
	}

	t.Type = domaintransaction.TransactionType(tType)
	t.IsActive = isActive == 1

	var errDate, errCreated, errUpdated error
	t.Date, errDate = time.Parse(dateLayout, date)
	t.CreatedAt, errCreated = time.Parse(timeLayout, createdAt)
	t.UpdatedAt, errUpdated = time.Parse(timeLayout, updatedAt)

	if errDate != nil {
		return domaintransaction.Transaction{}, fmt.Errorf("parse date: %w", errDate)
	}
	if errCreated != nil {
		return domaintransaction.Transaction{}, fmt.Errorf("parse created_at: %w", errCreated)
	}
	if errUpdated != nil {
		return domaintransaction.Transaction{}, fmt.Errorf("parse updated_at: %w", errUpdated)
	}

	return t, nil
}
