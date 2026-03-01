package pdfexport_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/pdfexport/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// buildMockRepo creates a mocks.Repository pre-configured with the given data.
func buildMockRepo(
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
	m.On("ListTransactions", mock.Anything, domaintransaction.TransactionTypeIncome, mock.Anything, mock.Anything).Return(incomes, nil).Once()
	m.On("ListTransactions", mock.Anything, domaintransaction.TransactionTypeExpense, mock.Anything, mock.Anything).Return(expenses, nil).Once()

	return m
}
