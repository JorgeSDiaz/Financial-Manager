// Package mocks contains testify mock implementations for the update use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// Repository is a testify mock for the update.Repository interface.
type Repository struct {
	mock.Mock
}

// GetByID mocks Repository.GetByID.
func (m *Repository) GetByID(ctx context.Context, id string) (domainaccount.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domainaccount.Account), args.Error(1)
}

// Update mocks Repository.Update.
func (m *Repository) Update(ctx context.Context, account domainaccount.Account) error {
	return m.Called(ctx, account).Error(0)
}
