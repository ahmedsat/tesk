package cmd

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	sqlc "github.com/ahmedsat/tesks/sql"
)

var Migrations embed.FS

var (
	ctx     = context.Background()
	db      *sql.DB
	queries *sqlc.Queries
)

func setupDatabase() (err error) {
	userHome := os.Getenv("HOME")
	if userHome == "" {
		return errors.New("HOME environment variable not set")
	}

	dataDir := filepath.Join(userHome, ".local", "share", "tesks")
	err = os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return
	}

	db, err = sql.Open("sqlite3", filepath.Join(dataDir, "tesks.db"))
	if err != nil {
		return
	}
	queries = sqlc.New(db)

	return
}

func runMigrations() (err error) {

	files, err := Migrations.ReadDir("migrations")
	if err != nil {
		return
	}

	slices.SortFunc(files, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})
	for _, f := range files {
		var b []byte
		b, err = Migrations.ReadFile("migrations/" + f.Name())
		if err != nil {
			return
		}
		_, err = db.ExecContext(ctx, string(b))
		if err != nil {
			return
		}
	}

	return
}

func cleanupTasks() (err error) {
	err = queries.DeleteOldTasks(ctx)
	if err != nil {
		return
	}
	err = queries.DeleteOlderTasks(ctx)
	if err != nil {
		return
	}

	return
}
