// Package main is the entry point for the Financial Manager API.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/financial-manager/api/internal/platform/config"
)

func main() {
	cfg := config.Load()

	deps, err := buildDependencies(cfg)
	if err != nil {
		log.Fatalf("startup: %v", err)
	}
	defer func() {
		if closeErr := deps.Databases.Close(); closeErr != nil {
			log.Printf("shutdown: closing databases: %v", closeErr)
		}
	}()

	router := registerRoutes(deps)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server starting on %s (env=%s)", addr, cfg.Env)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if serveErr := http.ListenAndServe(addr, router); serveErr != nil {
			log.Fatalf("server stopped: %v", serveErr)
		}
	}()

	<-stop
	log.Println("shutdown: signal received, exiting")
}
