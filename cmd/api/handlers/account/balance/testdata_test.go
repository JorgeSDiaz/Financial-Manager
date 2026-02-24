package balance_test

import (
	"context"
	"time"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

const fixedTimestamp = "2026-02-23T10:00:00Z"

type fakeUseCase struct {
	out domainaccount.Account
	err error
}

func (f *fakeUseCase) Execute(_ context.Context, _ string) (domainaccount.Account, error) {
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

func buildOutput(id, name string) domainaccount.Account {
	return buildDomainAccount(id, name)
}
