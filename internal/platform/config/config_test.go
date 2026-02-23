// Package config_test tests the config package.
package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/platform/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name            string
		env             map[string]string
		wantPort        string
		wantEnv         string
		wantDatabaseDir string
	}{
		{
			name:            "returns defaults when no env vars set",
			env:             map[string]string{},
			wantPort:        "8080",
			wantEnv:         "development",
			wantDatabaseDir: "~/FinancialManager/databases/",
		},
		{
			name:            "uses PORT env var when set",
			env:             map[string]string{"PORT": "9090"},
			wantPort:        "9090",
			wantEnv:         "development",
			wantDatabaseDir: "~/FinancialManager/databases/",
		},
		{
			name:            "uses ENV env var when set",
			env:             map[string]string{"ENV": "production"},
			wantPort:        "8080",
			wantEnv:         "production",
			wantDatabaseDir: "~/FinancialManager/databases/",
		},
		{
			name:            "uses both PORT and ENV when set",
			env:             map[string]string{"PORT": "3000", "ENV": "staging"},
			wantPort:        "3000",
			wantEnv:         "staging",
			wantDatabaseDir: "~/FinancialManager/databases/",
		},
		{
			name:            "uses DB_DIR env var when set",
			env:             map[string]string{"DB_DIR": "/custom/db/path/"},
			wantPort:        "8080",
			wantEnv:         "development",
			wantDatabaseDir: "/custom/db/path/",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			cfg := config.Load()

			assert.Equal(t, tc.wantPort, cfg.Port)
			assert.Equal(t, tc.wantEnv, cfg.Env)
			assert.Equal(t, tc.wantDatabaseDir, cfg.DatabaseDir)
		})
	}
}
