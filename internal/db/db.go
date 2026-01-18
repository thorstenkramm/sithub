// Package db provides SQLite access and migrations.
package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	// Register sqlite3 driver.
	_ "github.com/mattn/go-sqlite3"
)

// Open opens the SQLite database file in the data directory.
func Open(dataDir string) (*sql.DB, error) {
	path := filepath.Join(dataDir, "sithub.db")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, fmt.Errorf("set wal: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	return db, nil
}
