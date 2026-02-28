// Package mocks contains testify mock implementations for the list use case interfaces.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// Repository is a testify mock for the list.Repository interface.
type Repository struct {
	mock.Mock
}

// List mocks Repository.List.
func (m *Repository) List(ctx context.Context, categoryType *domaincategory.Type) ([]domaincategory.Category, error) {
	args := m.Called(ctx, categoryType)
	return args.Get(0).([]domaincategory.Category), args.Error(1)
}
