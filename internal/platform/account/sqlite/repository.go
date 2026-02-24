// Package sqlite implements the AccountRepository using SQLite.
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

const timeLayout = "2006-01-02T15:04:05Z"

// AccountRepository implements account repository interfaces using SQLite.
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates an AccountRepository with the provided *sql.DB.
// The caller is responsible for opening and closing the database connection.
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create inserts a new account row.
func (r *AccountRepository) Create(ctx context.Context, a domainaccount.Account) error {
	const q = `INSERT INTO accounts
		(id, name, type, initial_balance, current_balance, currency, color, icon, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	active := 0
	if a.IsActive {
		active = 1
	}

	_, err := r.db.ExecContext(ctx, q,
		a.ID, a.Name, string(a.Type),
		a.InitialBalance, a.CurrentBalance,
		a.Currency, a.Color, a.Icon,
		active,
		a.CreatedAt.UTC().Format(timeLayout),
		a.UpdatedAt.UTC().Format(timeLayout),
	)
	if err != nil {
		return fmt.Errorf("account sqlite: create: %w", err)
	}

	return nil
}

// GetByID retrieves an account by its ID regardless of is_active status.
// Returns domainshared.ErrNotFound if no row exists.
func (r *AccountRepository) GetByID(ctx context.Context, id string) (domainaccount.Account, error) {
	const q = `SELECT id, name, type, initial_balance, current_balance, currency, color, icon, is_active, created_at, updated_at
		FROM accounts WHERE id = ?`

	row := r.db.QueryRowContext(ctx, q, id)
	acc, err := scanAccount(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domainaccount.Account{}, domainshared.ErrNotFound
	}
	if err != nil {
		return domainaccount.Account{}, fmt.Errorf("account sqlite: get by id: %w", err)
	}

	return acc, nil
}

// List returns all active accounts (is_active = 1).
func (r *AccountRepository) List(ctx context.Context) ([]domainaccount.Account, error) {
	const q = `SELECT id, name, type, initial_balance, current_balance, currency, color, icon, is_active, created_at, updated_at
		FROM accounts WHERE is_active = 1`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("account sqlite: list: %w", err)
	}
	defer rows.Close()

	accounts := make([]domainaccount.Account, 0)
	for rows.Next() {
		acc, err := scanAccount(rows)
		if err != nil {
			return nil, fmt.Errorf("account sqlite: list scan: %w", err)
		}
		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("account sqlite: list rows: %w", err)
	}

	return accounts, nil
}

// Update modifies name, color, icon, and updated_at for an existing account.
// Type, initial_balance, and current_balance are immutable via this method.
func (r *AccountRepository) Update(ctx context.Context, a domainaccount.Account) error {
	const q = `UPDATE accounts SET name = ?, color = ?, icon = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, q,
		a.Name, a.Color, a.Icon,
		a.UpdatedAt.UTC().Format(timeLayout),
		a.ID,
	)
	if err != nil {
		return fmt.Errorf("account sqlite: update: %w", err)
	}

	return nil
}

// Delete soft-deletes an account by setting is_active = 0.
func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	const q = `UPDATE accounts SET is_active = 0 WHERE id = ?`

	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("account sqlite: delete: %w", err)
	}

	return nil
}

// HasTransactions always returns false in M2 — transactions table not yet implemented.
// TODO(M4): query the transactions table once it exists.
func (r *AccountRepository) HasTransactions(_ context.Context, _ string) (bool, error) {
	return false, nil
}

// scanner abstracts *sql.Row and *sql.Rows for the shared scanAccount helper.
type scanner interface {
	Scan(dest ...any) error
}

func scanAccount(s scanner) (domainaccount.Account, error) {
	var (
		a                    domainaccount.Account
		accType              string
		isActive             int
		createdAt, updatedAt string
	)

	err := s.Scan(
		&a.ID, &a.Name, &accType,
		&a.InitialBalance, &a.CurrentBalance,
		&a.Currency, &a.Color, &a.Icon,
		&isActive, &createdAt, &updatedAt,
	)
	if err != nil {
		return domainaccount.Account{}, err
	}

	a.Type = domainaccount.AccountType(accType)
	a.IsActive = isActive == 1

	a.CreatedAt, err = time.Parse(timeLayout, createdAt)
	if err != nil {
		return domainaccount.Account{}, fmt.Errorf("parse created_at: %w", err)
	}

	a.UpdatedAt, err = time.Parse(timeLayout, updatedAt)
	if err != nil {
		return domainaccount.Account{}, fmt.Errorf("parse updated_at: %w", err)
	}

	return a, nil
}
