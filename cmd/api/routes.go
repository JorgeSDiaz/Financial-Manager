package main

import (
	"net/http"

	accountbalance "github.com/financial-manager/api/cmd/api/handlers/account/balance"
	accountcreate "github.com/financial-manager/api/cmd/api/handlers/account/create"
	accountdelete "github.com/financial-manager/api/cmd/api/handlers/account/delete"
	accountget "github.com/financial-manager/api/cmd/api/handlers/account/get"
	accountlist "github.com/financial-manager/api/cmd/api/handlers/account/list"
	accountupdate "github.com/financial-manager/api/cmd/api/handlers/account/update"
	categorycreate "github.com/financial-manager/api/cmd/api/handlers/category/create"
	categorydelete "github.com/financial-manager/api/cmd/api/handlers/category/delete"
	categorylist "github.com/financial-manager/api/cmd/api/handlers/category/list"
	categoryupdate "github.com/financial-manager/api/cmd/api/handlers/category/update"
	healthhandler "github.com/financial-manager/api/cmd/api/handlers/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// registerRoutes builds the router by composing middlewares and all route groups.
func registerRoutes(svc *services) http.Handler {
	r := chi.NewRouter()
	registerMiddlewares(r)
	registerHealthRoutes(r, svc)
	registerAccountRoutes(r, svc)
	registerCategoryRoutes(r, svc)
	return r
}

// registerMiddlewares attaches global middlewares to the router.
func registerMiddlewares(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
}

// registerHealthRoutes mounts the health check endpoint.
func registerHealthRoutes(r *chi.Mux, svc *services) {
	h := healthhandler.NewHandler(svc.Health.Checker)
	r.Get("/health", h.Check)
}

// registerAccountRoutes mounts the /api/v1/accounts route group.
func registerAccountRoutes(r *chi.Mux, svc *services) {
	createHandler := accountcreate.New(svc.Accounts.Creator)
	listHandler := accountlist.New(svc.Accounts.Lister, svc.Accounts.BalanceGetter)
	getHandler := accountget.New(svc.Accounts.Getter)
	updateHandler := accountupdate.New(svc.Accounts.Updater)
	deleteHandler := accountdelete.New(svc.Accounts.Deleter)
	balanceHandler := accountbalance.New(svc.Accounts.Getter)

	r.Route("/api/v1/accounts", func(r chi.Router) {
		r.Post("/", createHandler.Handle)
		r.Get("/", listHandler.Handle)
		r.Get("/{id}", getHandler.Handle)
		r.Put("/{id}", updateHandler.Handle)
		r.Delete("/{id}", deleteHandler.Handle)
		r.Get("/{id}/balance", balanceHandler.Handle)
	})
}

// registerCategoryRoutes mounts the /api/v1/categories route group.
func registerCategoryRoutes(r *chi.Mux, svc *services) {
	createHandler := categorycreate.New(svc.Categories.Creator)
	listHandler := categorylist.New(svc.Categories.Lister)
	updateHandler := categoryupdate.New(svc.Categories.Updater)
	deleteHandler := categorydelete.New(svc.Categories.Deleter)

	r.Route("/api/v1/categories", func(r chi.Router) {
		r.Post("/", createHandler.Handle)
		r.Get("/", listHandler.Handle)
		r.Put("/{id}", updateHandler.Handle)
		r.Delete("/{id}", deleteHandler.Handle)
	})
}
