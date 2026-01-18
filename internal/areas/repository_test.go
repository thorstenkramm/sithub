package areas

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sithubdb "github.com/thorstenkramm/sithub/internal/db"
)

func TestRepositoryListOrdersAreas(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)

	_, err := db.Exec(
		`INSERT INTO areas (id, name, sort_order, created_at, updated_at)
		VALUES
		  ('b', 'Beta', 2, '2026-01-18T00:00:00Z', '2026-01-18T00:00:00Z'),
		  ('c', 'Gamma', 1, '2026-01-18T00:00:00Z', '2026-01-18T00:00:00Z'),
		  ('a', 'Alpha', 1, '2026-01-18T00:00:00Z', '2026-01-18T00:00:00Z')`,
	)
	require.NoError(t, err)

	repo := NewRepository(db)
	areas, err := repo.List(context.Background())
	require.NoError(t, err)

	require.Len(t, areas, 3)
	assert.Equal(t, "a", areas[0].ID)
	assert.Equal(t, "c", areas[1].ID)
	assert.Equal(t, "b", areas[2].ID)
}

func TestRepositoryListEmpty(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewRepository(db)

	areas, err := repo.List(context.Background())
	require.NoError(t, err)
	assert.Empty(t, areas)
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dir := t.TempDir()
	db, err := sithubdb.Open(dir)
	require.NoError(t, err)

	migrationsPath, err := resolveMigrationsPath()
	require.NoError(t, err)

	require.NoError(t, sithubdb.RunMigrations(db, migrationsPath))

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	return db
}

func resolveMigrationsPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve migrations path")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	return filepath.Join(root, "migrations"), nil
}
