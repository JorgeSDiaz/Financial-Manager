package delete_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/category/delete/mocks"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// buildMockRepoWithGet creates a mocks.Repository pre-configured for one GetByID call.
func buildMockRepoWithGet(id string, category domaincategory.Category, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("GetByID", mock.Anything, id).Return(category, err).Once()
	return m
}

// buildMockRepoFull creates a mocks.Repository pre-configured for GetByID, HasTransactions, and Delete.
func buildMockRepoFull(id string, category domaincategory.Category, hasTrans bool, deleteErr error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("GetByID", mock.Anything, id).Return(category, nil).Once()
	m.On("HasTransactions", mock.Anything, id).Return(hasTrans, nil).Once()
	m.On("Delete", mock.Anything, id).Return(deleteErr).Once()
	return m
}

// buildActiveCategory returns a valid active Category for use in tests.
func buildActiveCategory(id, name string) domaincategory.Category {
	return domaincategory.Category{
		ID:       id,
		Name:     name,
		Type:     domaincategory.TypeExpense,
		Color:    "#FFFFFF",
		Icon:     "wallet",
		IsSystem: false,
		IsActive: true,
	}
}
