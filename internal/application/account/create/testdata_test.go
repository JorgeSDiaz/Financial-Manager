package create_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/account/create/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

const (
	fixedID        = "fixed-uuid-0001"
	fixedTimestamp = "2026-02-23T10:00:00Z"
)

// validAccount is the expected account produced by a successful create with the standard valid input.
var validAccount = domainaccount.Account{
	ID:             fixedID,
	Name:           "Efectivo",
	Type:           domainaccount.AccountTypeCash,
	InitialBalance: 1000.0,
	CurrentBalance: 1000.0,
	Currency:       "USD",
	Color:          "#00FF00",
	Icon:           "wallet",
	IsActive:       true,
	CreatedAt:      fixedTime(),
	UpdatedAt:      fixedTime(),
}

// errorAccount is the account passed to the repo when input is minimal (name "X", type cash, no balance).
var errorAccount = domainaccount.Account{
	ID:        fixedID,
	Name:      "X",
	Type:      domainaccount.AccountTypeCash,
	IsActive:  true,
	CreatedAt: fixedTime(),
	UpdatedAt: fixedTime(),
}

// buildMockRepo creates a mocks.Repository pre-configured to accept one Create call
// with the given account and return the given error.
func buildMockRepo(account domainaccount.Account, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("Create", mock.Anything, account).Return(err).Once()
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
