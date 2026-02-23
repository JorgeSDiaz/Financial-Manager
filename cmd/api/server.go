package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/financial-manager/api/internal/platform/config"
	"github.com/financial-manager/api/internal/platform/database"
	"github.com/financial-manager/api/internal/platform/database/migrator"
	"github.com/financial-manager/api/internal/platform/database/sqlite"
)

// openDatabases initializes and opens all application databases.
func openDatabases(cfg *config.Config) (*database.Databases, error) {
	dbs := database.New(sqlite.NewConnector(), migrator.New())
	if err := dbs.Open(context.Background(), cfg.DatabaseDir); err != nil {
		return nil, fmt.Errorf("open databases: %w", err)
	}

	return dbs, nil
}

// closeDatabases closes all application databases, logging any error.
func closeDatabases(dbs *database.Databases) {
	if err := dbs.Close(); err != nil {
		log.Printf("shutdown: closing databases: %v", err)
	}
}

// run starts the HTTP server and blocks until a termination signal is received.
func run(cfg *config.Config, svc *services) {
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server starting on %s (env=%s)", addr, cfg.Env)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(addr, registerRoutes(svc)); err != nil {
			log.Fatalf("server stopped: %v", err)
		}
	}()

	<-stop
	log.Println("shutdown: signal received, exiting")
}
