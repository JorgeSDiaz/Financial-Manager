package sqlite_test

import (
	"context"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
	accountsqlite "github.com/financial-manager/api/internal/platform/account/sqlite"
)

func TestAccountRepository_CreateAndGetByID(t *testing.T) {
	t.Parallel()
	repo := accountsqlite.NewAccountRepository(newTestDB(t))
	ctx := context.Background()

	want := buildTestAccount("acc-1", "Efectivo")
	require.NoError(t, repo.Create(ctx, want))

	got, err := repo.GetByID(ctx, "acc-1")
	require.NoError(t, err)
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Type, got.Type)
	assert.InDelta(t, want.CurrentBalance, got.CurrentBalance, 0.001)
	assert.Equal(t, want.IsActive, got.IsActive)
}

func TestAccountRepository_GetByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := accountsqlite.NewAccountRepository(newTestDB(t))

	_, err := repo.GetByID(context.Background(), "missing")
	require.Error(t, err)
	assert.ErrorIs(t, err, domainshared.ErrNotFound)
}

func TestAccountRepository_List_OnlyActive(t *testing.T) {
	t.Parallel()
	repo := accountsqlite.NewAccountRepository(newTestDB(t))
	ctx := context.Background()

	active := buildTestAccount("a1", "Active")
	inactive := buildTestAccount("a2", "Inactive")
	inactive.IsActive = false

	require.NoError(t, repo.Create(ctx, active))
	require.NoError(t, repo.Create(ctx, inactive))

	accounts, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	assert.Equal(t, "a1", accounts[0].ID)
}

func TestAccountRepository_Update_ModifiesAllowedFieldsOnly(t *testing.T) {
	t.Parallel()
	repo := accountsqlite.NewAccountRepository(newTestDB(t))
	ctx := context.Background()

	original := buildTestAccount("acc-1", "Original")
	require.NoError(t, repo.Create(ctx, original))

	updated := original
	updated.Name = "Updated Name"
	updated.Color = "#FF0000"
	updated.Icon = "bank"
	updated.UpdatedAt = time.Now().UTC().Truncate(time.Second).Add(time.Minute)
	// These should NOT change in DB:
	updated.Type = domainaccount.AccountTypeBank
	updated.InitialBalance = 9999.0
	updated.CurrentBalance = 9999.0

	require.NoError(t, repo.Update(ctx, updated))

	got, err := repo.GetByID(ctx, "acc-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, "#FF0000", got.Color)
	assert.Equal(t, "bank", got.Icon)
	// Immutable fields unchanged:
	assert.Equal(t, domainaccount.AccountTypeCash, got.Type)
	assert.InDelta(t, 100.0, got.InitialBalance, 0.001)
	assert.InDelta(t, 100.0, got.CurrentBalance, 0.001)
}

func TestAccountRepository_Delete_SoftDeleteOnly(t *testing.T) {
	t.Parallel()
	repo := accountsqlite.NewAccountRepository(newTestDB(t))
	ctx := context.Background()

	acc := buildTestAccount("acc-1", "ToDelete")
	require.NoError(t, repo.Create(ctx, acc))

	require.NoError(t, repo.Delete(ctx, "acc-1"))

	// Row still exists in DB.
	got, err := repo.GetByID(ctx, "acc-1")
	require.NoError(t, err)
	assert.False(t, got.IsActive)

	// Not returned by List.
	accounts, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, accounts)
}

func TestAccountRepository_HasTransactions_ReturnsFalse(t *testing.T) {
	t.Parallel()
	repo := accountsqlite.NewAccountRepository(newTestDB(t))

	has, err := repo.HasTransactions(context.Background(), "any-id")
	require.NoError(t, err)
	assert.False(t, has)
}
