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
	"github.com/financial-manager/api/internal/platform/account/sqlite"
	categorysqlite "github.com/financial-manager/api/internal/platform/category/sqlite"
	"github.com/financial-manager/api/internal/platform/clock"
	"github.com/financial-manager/api/internal/platform/database"
	"github.com/financial-manager/api/internal/platform/idgen"
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

	// services holds all use case groups ready to be injected into the HTTP layer.
	services struct {
		Health     healthServices
		Accounts   accountServices
		Categories categoryServices
	}
)

// buildServices wires all use cases with their dependencies.
func buildServices(dbs *database.Databases) *services {
	accountRepo := sqlite.NewAccountRepository(dbs.Accounts)
	categoryRepo := categorysqlite.NewCategoryRepository(dbs.Categories)

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
	}
}
