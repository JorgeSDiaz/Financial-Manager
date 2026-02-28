package update_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/category/update/mocks"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

const fixedTimestamp = "2026-02-23T10:00:00Z"

// seeded is the canonical existing category used as the pre-update state in update tests.
var seeded = buildActiveCategory("cat-1", "Old Name")

// buildMockRepoGetByID creates a mocks.Repository pre-configured for one GetByID call.
func buildMockRepoGetByID(id string, category domaincategory.Category, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("GetByID", mock.Anything, id).Return(category, err).Once()
	return m
}

// buildMockRepoFull creates a mocks.Repository pre-configured for one GetByID and one Update call.
func buildMockRepoFull(id string, fetched, updated domaincategory.Category, updateErr error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("GetByID", mock.Anything, id).Return(fetched, nil).Once()
	m.On("Update", mock.Anything, updated).Return(updateErr).Once()
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
