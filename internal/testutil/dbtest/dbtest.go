// Package dbtest provides a migrated SQLite database for tests.
package dbtest

import (
	"database/sql"
	"testing"

	"github.com/thorstenkramm/sithub/internal/db"
)

// New opens a fresh SQLite database in a temporary directory with all migrations applied.
// The database is closed when the test finishes.
func New(t testing.TB) *sql.DB {
	t.Helper()
	store, err := db.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close() //nolint:errcheck // Cleanup function, error not critical
	})
	if err := db.RunMigrations(store); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	return store
}
