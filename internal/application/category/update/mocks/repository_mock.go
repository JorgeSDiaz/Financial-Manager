// Package mocks contains testify mock implementations for the update use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// Repository is a testify mock for the update.Repository interface.
type Repository struct {
	mock.Mock
}

// GetByID mocks Repository.GetByID.
func (m *Repository) GetByID(ctx context.Context, id string) (domaincategory.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domaincategory.Category), args.Error(1)
}

// Update mocks Repository.Update.
func (m *Repository) Update(ctx context.Context, category domaincategory.Category) error {
	return m.Called(ctx, category).Error(0)
}
