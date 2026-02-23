// Package main is the entry point for the Financial Manager API.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/financial-manager/api/internal/platform/config"
)

func main() {
	cfg := config.Load()

	deps := buildDependencies()
	router := registerRoutes(deps)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server starting on %s (env=%s)", addr, cfg.Env)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
