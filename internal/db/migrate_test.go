package db

import (
	"testing"
)

func TestRunMigrations(t *testing.T) {
	dir := t.TempDir()
	store, err := Open(dir)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close db: %v", err)
		}
	}()

	if err := RunMigrations(store); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
}
