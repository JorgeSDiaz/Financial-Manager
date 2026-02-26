package sqlite_test

import (
	"context"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
	categorysqlite "github.com/financial-manager/api/internal/platform/category/sqlite"
)

func TestCategoryRepository_CreateAndGetByID(t *testing.T) {
	t.Parallel()
	repo := categorysqlite.NewCategoryRepository(newTestDB(t))
	ctx := context.Background()

	want := buildTestCategory("cat-1", "Food")
	require.NoError(t, repo.Create(ctx, want))

	got, err := repo.GetByID(ctx, "cat-1")
	require.NoError(t, err)
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Type, got.Type)
	assert.Equal(t, want.Color, got.Color)
	assert.Equal(t, want.Icon, got.Icon)
	assert.Equal(t, want.IsActive, got.IsActive)
}

func TestCategoryRepository_GetByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := categorysqlite.NewCategoryRepository(newTestDB(t))

	_, err := repo.GetByID(context.Background(), "missing")
	require.Error(t, err)
	assert.ErrorIs(t, err, domainshared.ErrNotFound)
}

func TestCategoryRepository_List_OnlyActive(t *testing.T) {
	t.Parallel()
	repo := categorysqlite.NewCategoryRepository(newTestDB(t))
	ctx := context.Background()

	active := buildTestCategory("c1", "Active")
	inactive := buildTestCategory("c2", "Inactive")
	inactive.IsActive = false

	require.NoError(t, repo.Create(ctx, active))
	require.NoError(t, repo.Create(ctx, inactive))

	categories, err := repo.List(ctx, nil)
	require.NoError(t, err)
	require.Len(t, categories, 1)
	assert.Equal(t, "c1", categories[0].ID)
}

func TestCategoryRepository_List_FilterByType(t *testing.T) {
	t.Parallel()
	repo := categorysqlite.NewCategoryRepository(newTestDB(t))
	ctx := context.Background()

	expense := buildTestCategory("c1", "Expense")
	income := buildTestCategory("c2", "Income")
	income.Type = domaincategory.TypeIncome

	require.NoError(t, repo.Create(ctx, expense))
	require.NoError(t, repo.Create(ctx, income))

	expenseType := domaincategory.TypeExpense
	categories, err := repo.List(ctx, &expenseType)
	require.NoError(t, err)
	require.Len(t, categories, 1)
	assert.Equal(t, "c1", categories[0].ID)
	assert.Equal(t, domaincategory.TypeExpense, categories[0].Type)
}

func TestCategoryRepository_Update_ModifiesAllowedFieldsOnly(t *testing.T) {
	t.Parallel()
	repo := categorysqlite.NewCategoryRepository(newTestDB(t))
	ctx := context.Background()

	original := buildTestCategory("cat-1", "Original")
	require.NoError(t, repo.Create(ctx, original))

	updated := original
	updated.Name = "Updated Name"
	updated.Color = "#FF0000"
	updated.Icon = "bank"
	updated.UpdatedAt = time.Now().UTC().Truncate(time.Second).Add(time.Minute)
	// These should NOT change in DB:
	updated.Type = domaincategory.TypeIncome
	updated.IsSystem = true

	require.NoError(t, repo.Update(ctx, updated))

	got, err := repo.GetByID(ctx, "cat-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, "#FF0000", got.Color)
	assert.Equal(t, "bank", got.Icon)
	// Immutable fields unchanged:
	assert.Equal(t, domaincategory.TypeExpense, got.Type)
	assert.Equal(t, false, got.IsSystem)
}

func TestCategoryRepository_Delete_SoftDeleteOnly(t *testing.T) {
	t.Parallel()
	repo := categorysqlite.NewCategoryRepository(newTestDB(t))
	ctx := context.Background()

	cat := buildTestCategory("cat-1", "ToDelete")
	require.NoError(t, repo.Create(ctx, cat))

	require.NoError(t, repo.Delete(ctx, "cat-1"))

	// Should not appear in list
	categories, err := repo.List(ctx, nil)
	require.NoError(t, err)
	assert.Empty(t, categories)

	// But still exists in DB
	got, err := repo.GetByID(ctx, "cat-1")
	require.NoError(t, err)
	assert.Equal(t, "cat-1", got.ID)
	assert.False(t, got.IsActive)
}

func TestCategoryRepository_CountAll(t *testing.T) {
	t.Parallel()
	repo := categorysqlite.NewCategoryRepository(newTestDB(t))
	ctx := context.Background()

	count, err := repo.CountAll(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	require.NoError(t, repo.Create(ctx, buildTestCategory("c1", "One")))
	require.NoError(t, repo.Create(ctx, buildTestCategory("c2", "Two")))

	count, err = repo.CountAll(ctx)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestCategoryRepository_HasTransactions(t *testing.T) {
	t.Parallel()
	repo := categorysqlite.NewCategoryRepository(newTestDB(t))
	ctx := context.Background()

	// In M3, this always returns false
	hasTrans, err := repo.HasTransactions(ctx, "any-id")
	require.NoError(t, err)
	assert.False(t, hasTrans)
}
