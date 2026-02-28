package list_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/category/list/mocks"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// buildMockRepo creates a mocks.Repository pre-configured for one List call.
func buildMockRepo(categoryType *domaincategory.Type, categories []domaincategory.Category, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("List", mock.Anything, categoryType).Return(categories, err).Once()
	return m
}

// buildCategories returns a slice of test categories.
func buildCategories() []domaincategory.Category {
	return []domaincategory.Category{
		{ID: "cat-1", Name: "Food", Type: domaincategory.TypeExpense, Color: "#FF0000", Icon: "food", IsActive: true},
		{ID: "cat-2", Name: "Salary", Type: domaincategory.TypeIncome, Color: "#00FF00", Icon: "money", IsActive: true},
	}
}
