package globalbalance_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/account/globalbalance/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// buildMockRepo creates a mocks.Repository pre-configured to return the given
// accounts and error for one List call.
func buildMockRepo(accounts []domainaccount.Account, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("List", mock.Anything).Return(accounts, err).Once()
	return m
}

// cashAccountWith100 and bankAccountWith250 are fixtures with distinct current balances
// used to verify that global balance sums all accounts correctly.
var (
	cashAccountWith100 = func() domainaccount.Account {
		a := buildActiveAccount("acc-cash", "Efectivo")
		a.CurrentBalance = 100.0
		return a
	}()
	bankAccountWith250 = func() domainaccount.Account {
		a := buildActiveAccount("acc-bank", "Banco Nacional")
		a.CurrentBalance = 250.50
		return a
	}()
)

// buildActiveAccount returns a valid active Account for use in tests.
func buildActiveAccount(id, name string) domainaccount.Account {
	return domainaccount.Account{
		ID:             id,
		Name:           name,
		Type:           domainaccount.AccountTypeCash,
		InitialBalance: 500.0,
		CurrentBalance: 500.0,
		Currency:       "USD",
		IsActive:       true,
	}
}
