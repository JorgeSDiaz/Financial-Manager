package main

import (
	"github.com/financial-manager/api/internal/application/account/create"
	accountdelete "github.com/financial-manager/api/internal/application/account/delete"
	"github.com/financial-manager/api/internal/application/account/get"
	"github.com/financial-manager/api/internal/application/account/globalbalance"
	accountlist "github.com/financial-manager/api/internal/application/account/list"
	"github.com/financial-manager/api/internal/application/account/update"
	categorycreate "github.com/financial-manager/api/internal/application/category/create"
	categorydelete "github.com/financial-manager/api/internal/application/category/delete"
	categorylist "github.com/financial-manager/api/internal/application/category/list"
	categoryupdate "github.com/financial-manager/api/internal/application/category/update"
	"github.com/financial-manager/api/internal/application/health"
	transactiondelete "github.com/financial-manager/api/internal/application/transaction/delete"
	expensecreate "github.com/financial-manager/api/internal/application/transaction/expense/create"
	expenselist "github.com/financial-manager/api/internal/application/transaction/expense/list"
	incomecreate "github.com/financial-manager/api/internal/application/transaction/income/create"
	incomelist "github.com/financial-manager/api/internal/application/transaction/income/list"
	transactionsummary "github.com/financial-manager/api/internal/application/transaction/summary"
	transactionupdate "github.com/financial-manager/api/internal/application/transaction/update"
	"github.com/financial-manager/api/internal/platform/account/sqlite"
	categorysqlite "github.com/financial-manager/api/internal/platform/category/sqlite"
	"github.com/financial-manager/api/internal/platform/clock"
	"github.com/financial-manager/api/internal/platform/database"
	"github.com/financial-manager/api/internal/platform/idgen"
	transactionsqlite "github.com/financial-manager/api/internal/platform/transaction/sqlite"
)

type (
	// healthServices groups all use cases for the health resource.
	healthServices struct {
		Checker *health.CheckUseCase
	}

	// accountServices groups all use cases for the accounts resource.
	accountServices struct {
		Creator       *create.UseCase
		Getter        *get.UseCase
		Lister        *accountlist.UseCase
		Updater       *update.UseCase
		Deleter       *accountdelete.UseCase
		BalanceGetter *globalbalance.UseCase
	}

	// categoryServices groups all use cases for the categories resource.
	categoryServices struct {
		Creator *categorycreate.UseCase
		Lister  *categorylist.UseCase
		Updater *categoryupdate.UseCase
		Deleter *categorydelete.UseCase
	}

	// transactionServices groups all use cases for the transactions resource.
	transactionServices struct {
		IncomeCreator  *incomecreate.UseCase
		IncomeLister   *incomelist.UseCase
		ExpenseCreator *expensecreate.UseCase
		ExpenseLister  *expenselist.UseCase
		Updater        *transactionupdate.UseCase
		Deleter        *transactiondelete.UseCase
		Summary        *transactionsummary.UseCase
	}

	// services holds all use case groups ready to be injected into the HTTP layer.
	services struct {
		Health       healthServices
		Accounts     accountServices
		Categories   categoryServices
		Transactions transactionServices
	}
)

// buildServices wires all use cases with their dependencies.
func buildServices(dbs *database.Databases) *services {
	accountRepo := sqlite.NewAccountRepository(dbs.Accounts)
	categoryRepo := categorysqlite.NewCategoryRepository(dbs.Categories)
	transactionRepo := transactionsqlite.NewTransactionRepository(dbs.Transactions)

	return &services{
		Health: healthServices{
			Checker: health.NewCheckUseCase(),
		},
		Accounts: accountServices{
			Creator:       create.New(accountRepo, idgen.UUIDGenerator{}, clock.WallClock{}),
			Getter:        get.New(accountRepo),
			Lister:        accountlist.New(accountRepo),
			Updater:       update.New(accountRepo, clock.WallClock{}),
			Deleter:       accountdelete.New(accountRepo),
			BalanceGetter: globalbalance.New(accountRepo),
		},
		Categories: categoryServices{
			Creator: categorycreate.New(categoryRepo, idgen.UUIDGenerator{}, clock.WallClock{}),
			Lister:  categorylist.New(categoryRepo),
			Updater: categoryupdate.New(categoryRepo, clock.WallClock{}),
			Deleter: categorydelete.New(categoryRepo),
		},
		Transactions: transactionServices{
			IncomeCreator:  incomecreate.New(transactionRepo, idgen.UUIDGenerator{}, clock.WallClock{}),
			IncomeLister:   incomelist.New(transactionRepo),
			ExpenseCreator: expensecreate.New(transactionRepo, idgen.UUIDGenerator{}, clock.WallClock{}),
			ExpenseLister:  expenselist.New(transactionRepo),
			Updater:        transactionupdate.New(transactionRepo, clock.WallClock{}),
			Deleter:        transactiondelete.New(transactionRepo, clock.WallClock{}),
			Summary:        transactionsummary.New(transactionRepo),
		},
	}
}
