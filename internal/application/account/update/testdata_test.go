package update_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/account/update/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

const fixedTimestamp = "2026-02-23T10:00:00Z"

// seeded is the canonical existing account used as the pre-update state in update tests.
var seeded = buildActiveAccount("acc-1", "Old Name")

// buildMockRepoGetByID creates a mocks.Repository pre-configured for one GetByID call.
func buildMockRepoGetByID(id string, account domainaccount.Account, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("GetByID", mock.Anything, id).Return(account, err).Once()
	return m
}

// buildMockRepoFull creates a mocks.Repository pre-configured for one GetByID and one Update call.
func buildMockRepoFull(id string, fetched, updated domainaccount.Account, updateErr error) *mocks.Repository {
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

// buildActiveAccount returns a valid active Account for use in tests.
func buildActiveAccount(id, name string) domainaccount.Account {
	return domainaccount.Account{
		ID:             id,
		Name:           name,
		Type:           domainaccount.AccountTypeCash,
		InitialBalance: 500.0,
		CurrentBalance: 500.0,
		Currency:       "USD",
		Color:          "#FFFFFF",
		Icon:           "wallet",
		IsActive:       true,
	}
}
