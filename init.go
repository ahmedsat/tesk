package main

import (
	"embed"

	"github.com/ahmedsat/tesks/cmd"
	_ "github.com/mattn/go-sqlite3"
)

//go:generate sqlc generate

// Embed migration files
//
//go:embed migrations/*.sql
var migrations embed.FS

func init() {
	cmd.Migrations = migrations
}
