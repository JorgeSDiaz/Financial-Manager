// Package sqlite implements the CategoryRepository using SQLite.
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

const timeLayout = "2006-01-02T15:04:05Z"

// CategoryRepository implements category repository interfaces using SQLite.
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a CategoryRepository with the provided *sql.DB.
// The caller is responsible for opening and closing the database connection.
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create inserts a new category row.
func (r *CategoryRepository) Create(ctx context.Context, c domaincategory.Category) error {
	const q = `INSERT INTO categories
		(id, name, type, color, icon, is_system, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	isSystem := 0
	if c.IsSystem {
		isSystem = 1
	}
	isActive := 0
	if c.IsActive {
		isActive = 1
	}

	_, err := r.db.ExecContext(ctx, q,
		c.ID, c.Name, string(c.Type),
		c.Color, c.Icon,
		isSystem, isActive,
		c.CreatedAt.UTC().Format(timeLayout),
		c.UpdatedAt.UTC().Format(timeLayout),
	)
	if err != nil {
		return fmt.Errorf("category sqlite: create: %w", err)
	}

	return nil
}

// GetByID retrieves a category by its ID regardless of is_active status.
// Returns domainshared.ErrNotFound if no row exists.
func (r *CategoryRepository) GetByID(ctx context.Context, id string) (domaincategory.Category, error) {
	const q = `SELECT id, name, type, color, icon, is_system, is_active, created_at, updated_at
		FROM categories WHERE id = ?`

	row := r.db.QueryRowContext(ctx, q, id)
	cat, err := scanCategory(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domaincategory.Category{}, domainshared.ErrNotFound
	}
	if err != nil {
		return domaincategory.Category{}, fmt.Errorf("category sqlite: get by id: %w", err)
	}

	return cat, nil
}

// List returns all active categories (is_active = 1), optionally filtered by type.
func (r *CategoryRepository) List(ctx context.Context, categoryType *domaincategory.Type) ([]domaincategory.Category, error) {
	var q string
	var args []interface{}

	if categoryType != nil {
		q = `SELECT id, name, type, color, icon, is_system, is_active, created_at, updated_at
			FROM categories WHERE is_active = 1 AND type = ?`
		args = append(args, string(*categoryType))
	} else {
		q = `SELECT id, name, type, color, icon, is_system, is_active, created_at, updated_at
			FROM categories WHERE is_active = 1`
	}

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("category sqlite: list: %w", err)
	}
	defer rows.Close()

	categories := make([]domaincategory.Category, 0)
	for rows.Next() {
		cat, err := scanCategory(rows)
		if err != nil {
			return nil, fmt.Errorf("category sqlite: list scan: %w", err)
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("category sqlite: list rows: %w", err)
	}

	return categories, nil
}

// Update modifies name, color, icon, and updated_at for an existing category.
func (r *CategoryRepository) Update(ctx context.Context, c domaincategory.Category) error {
	const q = `UPDATE categories SET name = ?, color = ?, icon = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, q,
		c.Name, c.Color, c.Icon,
		c.UpdatedAt.UTC().Format(timeLayout),
		c.ID,
	)
	if err != nil {
		return fmt.Errorf("category sqlite: update: %w", err)
	}

	return nil
}

// Delete soft-deletes a category by setting is_active = 0.
func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	const q = `UPDATE categories SET is_active = 0 WHERE id = ?`

	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("category sqlite: delete: %w", err)
	}

	return nil
}

// HasTransactions always returns false in M3 — transactions table not yet implemented.
// TODO(M4): query the transactions table once it exists.
func (r *CategoryRepository) HasTransactions(_ context.Context, _ string) (bool, error) {
	return false, nil
}

// CountAll returns the total number of categories regardless of is_active status.
func (r *CategoryRepository) CountAll(ctx context.Context) (int, error) {
	const q = `SELECT COUNT(*) FROM categories`

	var count int
	if err := r.db.QueryRowContext(ctx, q).Scan(&count); err != nil {
		return 0, fmt.Errorf("category sqlite: count all: %w", err)
	}

	return count, nil
}

// scanner abstracts *sql.Row and *sql.Rows for the shared scanCategory helper.
type scanner interface {
	Scan(dest ...any) error
}

func scanCategory(s scanner) (domaincategory.Category, error) {
	var (
		c                    domaincategory.Category
		catType              string
		isSystem, isActive   int
		createdAt, updatedAt string
	)

	err := s.Scan(
		&c.ID, &c.Name, &catType,
		&c.Color, &c.Icon,
		&isSystem, &isActive,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return domaincategory.Category{}, err
	}

	c.Type = domaincategory.Type(catType)
	c.IsSystem = isSystem == 1
	c.IsActive = isActive == 1

	c.CreatedAt, err = time.Parse(timeLayout, createdAt)
	if err != nil {
		return domaincategory.Category{}, fmt.Errorf("parse created_at: %w", err)
	}

	c.UpdatedAt, err = time.Parse(timeLayout, updatedAt)
	if err != nil {
		return domaincategory.Category{}, fmt.Errorf("parse updated_at: %w", err)
	}

	return c, nil
}
