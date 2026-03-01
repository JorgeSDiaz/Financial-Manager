package export_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/export/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// buildMockRepoForCSV creates a mocks.Repository pre-configured for CSV export tests.
func buildMockRepoForCSV(
	accounts []domainaccount.Account,
	categories []domaincategory.Category,
	transactions []domaintransaction.Transaction,
	err error,
) *mocks.Repository {
	m := &mocks.Repository{}

	if err != nil {
		m.On("ListAccounts", mock.Anything).Return([]domainaccount.Account(nil), err).Once()
		return m
	}

	m.On("ListAccounts", mock.Anything).Return(accounts, nil).Once()
	m.On("ListCategories", mock.Anything).Return(categories, nil).Once()
	m.On("ListTransactions", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(transactions, nil).Once()

	return m
}

// buildMockRepoForJSON creates a mocks.Repository pre-configured for JSON export tests.
func buildMockRepoForJSON(
	accounts []domainaccount.Account,
	categories []domaincategory.Category,
	incomes []domaintransaction.Transaction,
	expenses []domaintransaction.Transaction,
	err error,
) *mocks.Repository {
	m := &mocks.Repository{}

	if err != nil {
		m.On("ListAccounts", mock.Anything).Return([]domainaccount.Account(nil), err).Once()
		return m
	}

	m.On("ListAccounts", mock.Anything).Return(accounts, nil).Once()
	m.On("ListCategories", mock.Anything).Return(categories, nil).Once()
	m.On("ListTransactions", mock.Anything, domaintransaction.TransactionTypeIncome, "", "").Return(incomes, nil).Once()
	m.On("ListTransactions", mock.Anything, domaintransaction.TransactionTypeExpense, "", "").Return(expenses, nil).Once()

	return m
}
