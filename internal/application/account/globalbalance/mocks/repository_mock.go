// Package mocks contains testify mock implementations for the globalbalance use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// Repository is a testify mock for the globalbalance.Repository interface.
type Repository struct {
	mock.Mock
}

// List mocks Repository.List.
func (m *Repository) List(ctx context.Context) ([]domainaccount.Account, error) {
	args := m.Called(ctx)
	accounts, _ := args.Get(0).([]domainaccount.Account)
	return accounts, args.Error(1)
}
