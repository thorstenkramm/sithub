package db

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestOpenSetsPragmas(t *testing.T) {
	dir := t.TempDir()
	db, err := Open(dir)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close db: %v", err)
		}
	}()

	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		t.Fatalf("query journal_mode: %v", err)
	}
	if journalMode != "wal" {
		t.Fatalf("expected wal journal mode, got %s", journalMode)
	}

	var foreignKeys int
	if err := db.QueryRow("PRAGMA foreign_keys;").Scan(&foreignKeys); err != nil {
		t.Fatalf("query foreign_keys: %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("expected foreign_keys=1, got %d", foreignKeys)
	}

	if _, err := db.Exec("CREATE TABLE sample(id INTEGER PRIMARY KEY)"); err != nil {
		t.Fatalf("create table: %v", err)
	}

	expected := filepath.Join(dir, "sithub.db")
	dbPath := dbStatsPath(db)
	expectedPath, err := filepath.EvalSymlinks(expected)
	if err != nil {
		t.Fatalf("eval expected path: %v", err)
	}
	actualPath, err := filepath.EvalSymlinks(dbPath)
	if err != nil {
		t.Fatalf("eval actual path: %v", err)
	}
	if actualPath != expectedPath {
		t.Fatalf("expected db path %s, got %s", expectedPath, actualPath)
	}
}

func dbStatsPath(db *sql.DB) string {
	var filePath string
	_ = db.QueryRow("PRAGMA database_list;").Scan(new(int), new(string), &filePath)
	return filePath
}
