// Package mocks contains testify mock implementations for the create use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// Repository is a testify mock for the create.Repository interface.
type Repository struct {
	mock.Mock
}

// Create mocks Repository.Create.
func (m *Repository) Create(ctx context.Context, account domainaccount.Account) error {
	return m.Called(ctx, account).Error(0)
}
