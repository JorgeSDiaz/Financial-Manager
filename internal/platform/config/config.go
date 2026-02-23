// Package config loads application configuration from environment variables.
package config

import (
	"os"
)

// Config holds all application-level configuration values.
type Config struct {
	Port        string
	Env         string
	DatabaseDir string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		DatabaseDir: getEnv("DB_DIR", "~/FinancialManager/databases/"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
