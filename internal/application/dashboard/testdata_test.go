package dashboard_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/dashboard/mocks"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

// Get current month dates for tests
func getCurrentMonthDates() (string, string) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	return startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02")
}

var currentDate, _ = time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

// buildMockRepo creates a mocks.Repository pre-configured with the given data.
func buildMockRepo(
	accounts []domainaccount.Account,
	recentTxs []domaintransaction.Transaction,
	expenseTxs []domaintransaction.Transaction,
	summaryIncomes []domaintransaction.Transaction,
	categories []domaincategory.Category,
	err error,
) *mocks.Repository {
	m := &mocks.Repository{}

	if err != nil {
		m.On("ListAccounts", mock.Anything).Return(nil, err).Once()
		return m
	}

	m.On("ListAccounts", mock.Anything).Return(accounts, nil).Once()
	m.On("ListRecentTransactions", mock.Anything, 10).Return(recentTxs, nil).Once()
	m.On("ListExpenseTransactions", mock.Anything, "", "", mock.Anything, mock.Anything).Return(expenseTxs, nil).Once()
	m.On("ListIncomeTransactions", mock.Anything, "", "", mock.Anything, mock.Anything).Return(summaryIncomes, nil).Once()
	m.On("ListCategories", mock.Anything).Return(categories, nil).Once()

	return m
}

// Account fixtures
var (
	account1 = func() domainaccount.Account {
		a := domainaccount.Account{
			ID:             "acc-1",
			Name:           "Banco",
			Type:           domainaccount.AccountTypeBank,
			InitialBalance: 1000.00,
			CurrentBalance: 1200.00,
			Currency:       "USD",
			IsActive:       true,
		}
		return a
	}()

	account2 = func() domainaccount.Account {
		a := domainaccount.Account{
			ID:             "acc-2",
			Name:           "Efectivo",
			Type:           domainaccount.AccountTypeCash,
			InitialBalance: 200.00,
			CurrentBalance: 300.00,
			Currency:       "USD",
			IsActive:       true,
		}
		return a
	}()
)

// Category fixtures
var (
	category1 = domaincategory.Category{
		ID:       "cat-1",
		Name:     "Alimentación",
		Type:     domaincategory.TypeExpense,
		Color:    "#FF0000",
		Icon:     "food",
		IsSystem: false,
		IsActive: true,
	}

	category2 = domaincategory.Category{
		ID:       "cat-2",
		Name:     "Transporte",
		Type:     domaincategory.TypeExpense,
		Color:    "#00FF00",
		Icon:     "transport",
		IsSystem: false,
		IsActive: true,
	}
)

// Transaction fixtures - all using current month dates
var (
	today = currentDate
	tx1   = buildTransaction("tx-1", domaintransaction.TransactionTypeIncome, 500.00, "Salary", today)
	tx2   = buildTransaction("tx-2", domaintransaction.TransactionTypeIncome, 100.00, "Bonus", today.AddDate(0, 0, -1))
	tx3   = buildTransaction("tx-3", domaintransaction.TransactionTypeExpense, 50.00, "Groceries", today.AddDate(0, 0, -2))
	tx4   = buildTransaction("tx-4", domaintransaction.TransactionTypeExpense, 50.00, "Food", today.AddDate(0, 0, -3))
	tx5   = buildTransactionWithCategory("tx-5", domaintransaction.TransactionTypeExpense, 30.00, "Bus", today.AddDate(0, 0, -4), "cat-2")
	tx6   = buildTransactionWithCategory("tx-6", domaintransaction.TransactionTypeExpense, 20.00, "Taxi", today.AddDate(0, 0, -5), "cat-2")
	tx7   = buildTransaction("tx-7", domaintransaction.TransactionTypeIncome, 200.00, "Freelance", today.AddDate(0, 0, -6))
	tx8   = buildTransaction("tx-8", domaintransaction.TransactionTypeExpense, 40.00, "Dinner", today.AddDate(0, 0, -7))
	tx9   = buildTransaction("tx-9", domaintransaction.TransactionTypeIncome, 150.00, "Refund", today.AddDate(0, 0, -8))
	tx10  = buildTransaction("tx-10", domaintransaction.TransactionTypeExpense, 25.00, "Coffee", today.AddDate(0, 0, -9))
	tx11  = buildTransaction("tx-11", domaintransaction.TransactionTypeExpense, 15.00, "Snack", today.AddDate(0, 0, -10))
)

// buildTransaction creates a transaction fixture.
func buildTransactionWithCategory(id string, tType domaintransaction.TransactionType, amount float64, desc string, date time.Time, categoryID string) domaintransaction.Transaction {
	return domaintransaction.Transaction{
		ID:          id,
		AccountID:   "acc-1",
		CategoryID:  categoryID,
		Type:        tType,
		Amount:      amount,
		Description: desc,
		Date:        date,
		IsActive:    true,
	}
}

// buildTransaction creates a transaction fixture with default category.
func buildTransaction(id string, tType domaintransaction.TransactionType, amount float64, desc string, date time.Time) domaintransaction.Transaction {
	categoryID := ""
	if tType == domaintransaction.TransactionTypeExpense {
		categoryID = "cat-1"
	}
	return buildTransactionWithCategory(id, tType, amount, desc, date, categoryID)
}
