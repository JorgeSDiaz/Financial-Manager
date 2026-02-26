// Package main is the entry point for the Financial Manager API.
package main

import (
	"log"

	"github.com/financial-manager/api/internal/platform/config"
)

func main() {
	cfg := config.Load()

	dbs, err := openDatabases(cfg)
	if err != nil {
		log.Fatalf("startup: %v", err)
	}
	defer closeDatabases(dbs)

	svc := buildServices(dbs)

	run(cfg, svc)
}
