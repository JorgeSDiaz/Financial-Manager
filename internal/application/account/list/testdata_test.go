package list_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/account/list/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// buildMockRepo creates a mocks.Repository pre-configured to return the given
// accounts and error for one List call.
func buildMockRepo(accounts []domainaccount.Account, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("List", mock.Anything).Return(accounts, err).Once()
	return m
}

// cashAccount and bankAccount are canonical active account fixtures for list tests.
var (
	cashAccount = buildActiveAccount("acc-cash", "Efectivo")
	bankAccount = buildActiveAccount("acc-bank", "Banco Nacional")
)

// buildActiveAccount returns a valid active Account for use in tests.
func buildActiveAccount(id, name string) domainaccount.Account {
	return domainaccount.Account{
		ID:       id,
		Name:     name,
		Type:     domainaccount.AccountTypeCash,
		Currency: "USD",
		IsActive: true,
	}
}
