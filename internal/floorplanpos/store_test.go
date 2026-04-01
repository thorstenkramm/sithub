package floorplanpos

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/db"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dir := t.TempDir()
	store, err := db.Open(dir)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, store.Close()) })
	require.NoError(t, db.RunMigrations(store))
	return store
}

func TestCreateAndFindByFloorPlan(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	ctx := context.Background()

	pos, err := Create(ctx, store, &CreateInput{
		FloorPlan: "office.svg",
		ItemID:    "desk-1",
		Label:     "D1",
		X:         10.5, Y: 20.3, Width: 5.0, Height: 3.0,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, pos.ID)
	assert.Equal(t, "office.svg", pos.FloorPlan)
	assert.Equal(t, "D1", pos.Label)

	positions, err := FindByFloorPlan(ctx, store, "office.svg")
	require.NoError(t, err)
	require.Len(t, positions, 1)
	assert.Equal(t, pos.ID, positions[0].ID)
	assert.InDelta(t, 10.5, positions[0].X, 0.01)
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	ctx := context.Background()

	pos, err := Create(ctx, store, &CreateInput{
		FloorPlan: "office.svg", ItemID: "desk-1",
		X: 10, Y: 20, Width: 5, Height: 3,
	})
	require.NoError(t, err)

	newX := 15.0
	newLabel := "Updated"
	updated, err := Update(ctx, store, pos.ID, UpdateInput{
		X: &newX, Label: &newLabel,
	})
	require.NoError(t, err)
	assert.InDelta(t, 15.0, updated.X, 0.01)
	assert.Equal(t, "Updated", updated.Label)
}

func TestUpdateNotFound(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	ctx := context.Background()

	newX := 15.0
	_, err := Update(ctx, store, "nonexistent", UpdateInput{X: &newX})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	ctx := context.Background()

	pos, err := Create(ctx, store, &CreateInput{
		FloorPlan: "office.svg", ItemID: "desk-1",
		X: 10, Y: 20, Width: 5, Height: 3,
	})
	require.NoError(t, err)

	require.NoError(t, Delete(ctx, store, pos.ID))

	positions, err := FindByFloorPlan(ctx, store, "office.svg")
	require.NoError(t, err)
	assert.Empty(t, positions)
}

func TestDeleteNotFound(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	ctx := context.Background()

	err := Delete(ctx, store, "nonexistent")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestFindByFloorPlanEmpty(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	ctx := context.Background()

	positions, err := FindByFloorPlan(ctx, store, "nothing.svg")
	require.NoError(t, err)
	assert.Empty(t, positions)
}
