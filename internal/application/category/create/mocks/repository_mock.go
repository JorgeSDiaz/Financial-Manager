package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// Repository is a testify mock for the create.Repository interface.
type Repository struct {
	mock.Mock
}

// Create mocks Repository.Create.
func (m *Repository) Create(ctx context.Context, category domaincategory.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}
