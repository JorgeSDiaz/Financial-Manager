// Package database manages the SQLite databases and their migrations.
package database

import "embed"

//go:embed migrations
var migrationFiles embed.FS
