package list_test

import (
	"context"
	"time"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

const fixedTimestamp = "2026-02-23T10:00:00Z"

type fakeLister struct {
	out []domainaccount.Account
	err error
}

func (f *fakeLister) Execute(_ context.Context) ([]domainaccount.Account, error) {
	return f.out, f.err
}

type fakeBalanceGetter struct {
	out float64
	err error
}

func (f *fakeBalanceGetter) Execute(_ context.Context) (float64, error) {
	return f.out, f.err
}

func buildDomainAccount(id, name string) domainaccount.Account {
	t, _ := time.Parse("2006-01-02T15:04:05Z", fixedTimestamp)
	return domainaccount.Account{
		ID:             id,
		Name:           name,
		Type:           domainaccount.AccountTypeCash,
		InitialBalance: 1000.0,
		CurrentBalance: 1000.0,
		Currency:       "USD",
		Color:          "#00FF00",
		Icon:           "wallet",
		IsActive:       true,
		CreatedAt:      t,
		UpdatedAt:      t,
	}
}

func buildListOutput(accounts ...domainaccount.Account) []domainaccount.Account {
	return accounts
}

func buildBalanceOutput(total float64) float64 {
	return total
}
