package get_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/account/get/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// buildMockRepo creates a mocks.Repository pre-configured to return the given
// account and error for one GetByID call with the specified id.
func buildMockRepo(id string, account domainaccount.Account, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("GetByID", mock.Anything, id).Return(account, err).Once()
	return m
}

// activeAccount is a valid active account used as the canonical fixture in get tests.
var activeAccount = buildActiveAccount("acc-1", "Efectivo")

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
