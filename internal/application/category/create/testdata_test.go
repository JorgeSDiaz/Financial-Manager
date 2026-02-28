package create_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/category/create/mocks"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

const (
	fixedID        = "fixed-uuid-cat001"
	fixedTimestamp = "2026-02-23T10:00:00Z"
)

// validCategory is the expected category produced by a successful create with the standard valid input.
var validCategory = domaincategory.Category{
	ID:        fixedID,
	Name:      "Comida",
	Type:      domaincategory.TypeExpense,
	Color:     "#FF5733",
	Icon:      "restaurant",
	IsSystem:  false,
	IsActive:  true,
	CreatedAt: fixedTime(),
	UpdatedAt: fixedTime(),
}

// errorCategory is the category passed to the repo when input is minimal.
var errorCategory = domaincategory.Category{
	ID:        fixedID,
	Name:      "X",
	Type:      domaincategory.TypeExpense,
	Color:     "red",
	Icon:      "icon",
	IsSystem:  false,
	IsActive:  true,
	CreatedAt: fixedTime(),
	UpdatedAt: fixedTime(),
}

// buildMockRepo creates a mocks.Repository pre-configured to accept one Create call
// with the given category and return the given error.
func buildMockRepo(category domaincategory.Category, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("Create", mock.Anything, category).Return(err).Once()
	return m
}

// buildMockIDGenerator creates a mocks.IDGenerator pre-configured to return fixedID once.
func buildMockIDGenerator() *mocks.IDGenerator {
	m := &mocks.IDGenerator{}
	m.On("NewID").Return(fixedID).Once()
	return m
}

// buildMockClock creates a mocks.Clock pre-configured to return fixedTime once.
func buildMockClock() *mocks.Clock {
	m := &mocks.Clock{}
	m.On("Now").Return(fixedTime()).Once()
	return m
}

// fixedTime parses fixedTimestamp and panics on error (test helper).
func fixedTime() time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", fixedTimestamp)
	if err != nil {
		panic(err)
	}
	return t.UTC()
}
