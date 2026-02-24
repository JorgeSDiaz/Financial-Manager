package main

import (
	"database/sql"

	"github.com/financial-manager/api/internal/application/account/create"
	accountdelete "github.com/financial-manager/api/internal/application/account/delete"
	"github.com/financial-manager/api/internal/application/account/get"
	"github.com/financial-manager/api/internal/application/account/globalbalance"
	"github.com/financial-manager/api/internal/application/account/list"
	"github.com/financial-manager/api/internal/application/account/update"
	"github.com/financial-manager/api/internal/application/health"
	accountsqlite "github.com/financial-manager/api/internal/platform/account/sqlite"
	"github.com/financial-manager/api/internal/platform/clock"
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
		Lister        *list.UseCase
		Updater       *update.UseCase
		Deleter       *accountdelete.UseCase
		BalanceGetter *globalbalance.UseCase
	}

	// services holds all use case groups ready to be injected into the HTTP layer.
	services struct {
		Health   healthServices
		Accounts accountServices
	}
)

// buildServices wires all use cases with their dependencies.
func buildServices(accountsDB *sql.DB) *services {
	repo := accountsqlite.NewAccountRepository(accountsDB)
	return &services{
		Health: healthServices{
			Checker: health.NewCheckUseCase(),
		},
		Accounts: accountServices{
			Creator:       create.New(repo, idgen.UUIDGenerator{}, clock.WallClock{}),
			Getter:        get.New(repo),
			Lister:        list.New(repo),
			Updater:       update.New(repo, clock.WallClock{}),
			Deleter:       accountdelete.New(repo),
			BalanceGetter: globalbalance.New(repo),
		},
	}
}
