package main

import (
	"net/http"

	healthhandler "github.com/financial-manager/api/cmd/api/handlers/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// registerRoutes wires HTTP routes to their handlers and returns the router.
func registerRoutes(deps *dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	healthHandler := healthhandler.NewHandler(deps.HealthChecker)

	r.Get("/health", healthHandler.Check)

	return r
}
