package export_test

import (
	"errors"

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

// buildMockRepoForCSVWithType creates a mocks.Repository for CSV export with type filter.
func buildMockRepoForCSVWithType(
	accounts []domainaccount.Account,
	categories []domaincategory.Category,
	transactions []domaintransaction.Transaction,
	err error,
	txType string,
) *mocks.Repository {
	m := &mocks.Repository{}

	m.On("ListAccounts", mock.Anything).Return(accounts, nil).Once()
	m.On("ListCategories", mock.Anything).Return(categories, nil).Once()

	var tType domaintransaction.TransactionType
	if txType == "income" {
		tType = domaintransaction.TransactionTypeIncome
	} else if txType == "expense" {
		tType = domaintransaction.TransactionTypeExpense
	}
	m.On("ListTransactions", mock.Anything, tType, "", "").Return(transactions, err).Once()

	return m
}

// buildMockRepoForCSVWithCategoriesError creates a mock that returns error on ListCategories.
func buildMockRepoForCSVWithCategoriesError() *mocks.Repository {
	m := &mocks.Repository{}
	m.On("ListAccounts", mock.Anything).Return([]domainaccount.Account{{ID: "acc-1", Name: "Banco"}}, nil).Once()
	m.On("ListCategories", mock.Anything).Return([]domaincategory.Category(nil), errors.New("categories error")).Once()
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

// buildMockRepoForJSONWithCategoriesError creates a mock that returns error on ListCategories.
func buildMockRepoForJSONWithCategoriesError() *mocks.Repository {
	m := &mocks.Repository{}
	m.On("ListAccounts", mock.Anything).Return([]domainaccount.Account{{ID: "acc-1", Name: "Banco"}}, nil).Once()
	m.On("ListCategories", mock.Anything).Return([]domaincategory.Category(nil), errors.New("categories error")).Once()
	return m
}

// buildMockRepoForJSONWithExpensesError creates a mock that returns error on ListTransactions for expenses.
func buildMockRepoForJSONWithExpensesError() *mocks.Repository {
	m := &mocks.Repository{}
	m.On("ListAccounts", mock.Anything).Return([]domainaccount.Account{{ID: "acc-1", Name: "Banco"}}, nil).Once()
	m.On("ListCategories", mock.Anything).Return([]domaincategory.Category{{ID: "cat-1", Name: "Food"}}, nil).Once()
	m.On("ListTransactions", mock.Anything, domaintransaction.TransactionTypeIncome, "", "").Return([]domaintransaction.Transaction{}, nil).Once()
	m.On("ListTransactions", mock.Anything, domaintransaction.TransactionTypeExpense, "", "").Return([]domaintransaction.Transaction(nil), errors.New("expenses error")).Once()
	return m
}
