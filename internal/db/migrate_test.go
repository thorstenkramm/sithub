package db

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	store, err := Open(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = store.Close() //nolint:errcheck // Cleanup function, error not critical
	})
	return store
}

func TestRunMigrations(t *testing.T) {
	t.Parallel()
	store := openTestDB(t)

	require.NoError(t, RunMigrations(store))
	// A second run is a no-op.
	require.NoError(t, RunMigrations(store))
}

func TestMigration000006APIKeys(t *testing.T) {
	t.Parallel()
	store := openTestDB(t)
	require.NoError(t, RunMigrations(store))
	const insert = `INSERT INTO api_keys
		(id, name, key_prefix, key_hash, scope, expires_at, created_at)
		VALUES (?, 'Export', '12345678-9', 'hash', 'read-bookings', '2027-10-08T00:00:00Z', '2026-10-08T00:00:00Z')`
	_, err := store.Exec(insert, "key1")
	require.NoError(t, err)
	var name, creator string
	require.NoError(t, store.QueryRow(`SELECT name, created_by FROM api_keys WHERE id = 'key1'`).Scan(&name, &creator))
	assert.Equal(t, "Export", name)
	assert.Empty(t, creator)
	_, err = store.Exec(insert, "key2")
	require.ErrorContains(t, err, "UNIQUE")
	m, err := newMigrator(store)
	require.NoError(t, err)
	require.NoError(t, m.Migrate(5))
	_, err = store.Exec(`SELECT id FROM api_keys`)
	require.Error(t, err)
}

func TestMigration000004KeepsExistingUsersActive(t *testing.T) {
	t.Parallel()
	store := openTestDB(t)

	m, err := newMigrator(store)
	require.NoError(t, err)
	require.NoError(t, m.Migrate(3))

	_, err = store.Exec(`INSERT INTO users (id, email, display_name, user_source, created_at, updated_at)
		VALUES ('u1', 'u1@example.com', 'User One', 'internal', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`)
	require.NoError(t, err)

	require.NoError(t, RunMigrations(store))

	var deleted int
	require.NoError(t, store.QueryRow(`SELECT deleted FROM users WHERE id = 'u1'`).Scan(&deleted))
	assert.Equal(t, 0, deleted)

	var count int
	require.NoError(t, store.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count))
	assert.Equal(t, 1, count)
}

func TestMigration000004Down(t *testing.T) {
	t.Parallel()
	store := openTestDB(t)
	require.NoError(t, RunMigrations(store))

	m, err := newMigrator(store)
	require.NoError(t, err)
	require.NoError(t, m.Migrate(3))

	_, err = store.Exec(`SELECT deleted FROM users`)
	require.Error(t, err)
}

func TestMigration000005NoUserGainsDataExporterRole(t *testing.T) {
	t.Parallel()
	store := openTestDB(t)

	m, err := newMigrator(store)
	require.NoError(t, err)
	require.NoError(t, m.Migrate(4))

	_, err = store.Exec(`INSERT INTO users (id, email, display_name, user_source, is_admin, created_at, updated_at)
		VALUES ('a1', 'a1@example.com', 'Admin', 'internal', 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`)
	require.NoError(t, err)

	require.NoError(t, RunMigrations(store))

	var flag int
	require.NoError(t, store.QueryRow(`SELECT is_data_exporter FROM users WHERE id = 'a1'`).Scan(&flag))
	assert.Equal(t, 0, flag)

	require.NoError(t, m.Migrate(4))
	_, err = store.Exec(`SELECT is_data_exporter FROM users`)
	require.Error(t, err)
}
