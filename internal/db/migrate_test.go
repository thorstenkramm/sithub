package db

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestRunMigrations(t *testing.T) {
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

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("unable to resolve test path")
	}

	migrationsPath := filepath.Join(filepath.Dir(filename), "..", "..", "migrations")
	if err := RunMigrations(db, migrationsPath); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
}
