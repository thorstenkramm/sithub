package spaces

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/db"
)

func setupTestStore(t *testing.T) *sql.DB {
	t.Helper()

	dir := t.TempDir()
	store, err := db.Open(dir)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, store.Close())
	})

	migrationsPath := resolveTestMigrationsPath(t)
	require.NoError(t, db.RunMigrations(store, migrationsPath))

	return store
}

func resolveTestMigrationsPath(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)

	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	return filepath.Join(root, "migrations")
}

// --- Area Tests ---

func TestStore_CreateAndGetArea(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	area := &AreaRecord{
		ID:          "area-1",
		Name:        "Office",
		Description: "Main office",
		FloorPlan:   "floor1.png",
	}

	err := spacesStore.CreateArea(ctx, area)
	require.NoError(t, err)
	assert.NotEmpty(t, area.CreatedAt)
	assert.NotEmpty(t, area.UpdatedAt)

	// Get area
	got, err := spacesStore.GetArea(ctx, "area-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Office", got.Name)
	assert.Equal(t, "Main office", got.Description)
}

func TestStore_GetAreaNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	got, err := spacesStore.GetArea(ctx, "missing")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestStore_ListAreas(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	// Create areas
	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-b", Name: "Beta"}))
	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-a", Name: "Alpha"}))

	areas, err := spacesStore.ListAreas(ctx)
	require.NoError(t, err)
	require.Len(t, areas, 2)
	// Should be ordered by name
	assert.Equal(t, "Alpha", areas[0].Name)
	assert.Equal(t, "Beta", areas[1].Name)
}

func TestStore_UpdateArea(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	// Create area
	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))

	// Update
	err := spacesStore.UpdateArea(ctx, &AreaRecord{
		ID:          "area-1",
		Name:        "Updated Office",
		Description: "New description",
	})
	require.NoError(t, err)

	got, err := spacesStore.GetArea(ctx, "area-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated Office", got.Name)
	assert.Equal(t, "New description", got.Description)
}

func TestStore_UpdateAreaNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	err := spacesStore.UpdateArea(ctx, &AreaRecord{ID: "missing", Name: "Test"})
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestStore_DeleteArea(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))

	err := spacesStore.DeleteArea(ctx, "area-1")
	require.NoError(t, err)

	got, err := spacesStore.GetArea(ctx, "area-1")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestStore_DeleteAreaNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	err := spacesStore.DeleteArea(ctx, "missing")
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

// --- Room Tests ---

func TestStore_CreateAndGetRoom(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	// Create area first
	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))

	room := &RoomRecord{
		ID:          "room-1",
		AreaID:      "area-1",
		Name:        "Room 101",
		Description: "Conference room",
	}

	err := spacesStore.CreateRoom(ctx, room)
	require.NoError(t, err)

	got, err := spacesStore.GetRoom(ctx, "room-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Room 101", got.Name)
	assert.Equal(t, "area-1", got.AreaID)
}

func TestStore_GetRoomNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	got, err := spacesStore.GetRoom(ctx, "missing")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestStore_ListRooms(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))
	require.NoError(t, spacesStore.CreateRoom(ctx, &RoomRecord{ID: "room-b", AreaID: "area-1", Name: "Beta"}))
	require.NoError(t, spacesStore.CreateRoom(ctx, &RoomRecord{ID: "room-a", AreaID: "area-1", Name: "Alpha"}))

	rooms, err := spacesStore.ListRooms(ctx, "area-1")
	require.NoError(t, err)
	require.Len(t, rooms, 2)
	assert.Equal(t, "Alpha", rooms[0].Name)
}

func TestStore_UpdateRoom(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))
	require.NoError(t, spacesStore.CreateRoom(ctx, &RoomRecord{ID: "room-1", AreaID: "area-1", Name: "Room"}))

	err := spacesStore.UpdateRoom(ctx, &RoomRecord{ID: "room-1", Name: "Updated Room"})
	require.NoError(t, err)

	got, err := spacesStore.GetRoom(ctx, "room-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated Room", got.Name)
}

func TestStore_DeleteRoom(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))
	require.NoError(t, spacesStore.CreateRoom(ctx, &RoomRecord{ID: "room-1", AreaID: "area-1", Name: "Room"}))

	err := spacesStore.DeleteRoom(ctx, "room-1")
	require.NoError(t, err)

	got, err := spacesStore.GetRoom(ctx, "room-1")
	require.NoError(t, err)
	assert.Nil(t, got)
}

// --- Desk Tests ---

func TestStore_CreateAndGetDesk(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))
	require.NoError(t, spacesStore.CreateRoom(ctx, &RoomRecord{ID: "room-1", AreaID: "area-1", Name: "Room"}))

	desk := &DeskRecord{
		ID:        "desk-1",
		RoomID:    "room-1",
		Name:      "Corner Desk",
		Equipment: []string{"monitor", "keyboard"},
		Warning:   "Near exit",
	}

	err := spacesStore.CreateDesk(ctx, desk)
	require.NoError(t, err)

	got, err := spacesStore.GetDesk(ctx, "desk-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Corner Desk", got.Name)
	assert.Equal(t, []string{"monitor", "keyboard"}, got.Equipment)
	assert.Equal(t, "Near exit", got.Warning)
}

func TestStore_GetDeskNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	got, err := spacesStore.GetDesk(ctx, "missing")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestStore_ListDesks(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))
	require.NoError(t, spacesStore.CreateRoom(ctx, &RoomRecord{ID: "room-1", AreaID: "area-1", Name: "Room"}))
	require.NoError(t, spacesStore.CreateDesk(ctx, &DeskRecord{ID: "desk-b", RoomID: "room-1", Name: "Beta"}))
	require.NoError(t, spacesStore.CreateDesk(ctx, &DeskRecord{ID: "desk-a", RoomID: "room-1", Name: "Alpha"}))

	desks, err := spacesStore.ListDesks(ctx, "room-1")
	require.NoError(t, err)
	require.Len(t, desks, 2)
	assert.Equal(t, "Alpha", desks[0].Name)
}

func TestStore_UpdateDesk(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))
	require.NoError(t, spacesStore.CreateRoom(ctx, &RoomRecord{ID: "room-1", AreaID: "area-1", Name: "Room"}))
	require.NoError(t, spacesStore.CreateDesk(ctx, &DeskRecord{ID: "desk-1", RoomID: "room-1", Name: "Desk"}))

	err := spacesStore.UpdateDesk(ctx, &DeskRecord{
		ID:        "desk-1",
		Name:      "Updated Desk",
		Equipment: []string{"lamp"},
		Warning:   "Hot spot",
	})
	require.NoError(t, err)

	got, err := spacesStore.GetDesk(ctx, "desk-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated Desk", got.Name)
	assert.Equal(t, []string{"lamp"}, got.Equipment)
	assert.Equal(t, "Hot spot", got.Warning)
}

func TestStore_DeleteDesk(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))
	require.NoError(t, spacesStore.CreateRoom(ctx, &RoomRecord{ID: "room-1", AreaID: "area-1", Name: "Room"}))
	require.NoError(t, spacesStore.CreateDesk(ctx, &DeskRecord{ID: "desk-1", RoomID: "room-1", Name: "Desk"}))

	err := spacesStore.DeleteDesk(ctx, "desk-1")
	require.NoError(t, err)

	got, err := spacesStore.GetDesk(ctx, "desk-1")
	require.NoError(t, err)
	assert.Nil(t, got)
}

// --- Sync and Load Tests ---

func TestStore_SyncFromConfig(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	cfg := &Config{
		Areas: []Area{
			{
				ID:   "area-1",
				Name: "Office",
				Rooms: []Room{
					{
						ID:   "room-1",
						Name: "Room 101",
						Desks: []Desk{
							{ID: "desk-1", Name: "Desk A", Equipment: []string{"monitor"}},
						},
					},
				},
			},
		},
	}

	err := spacesStore.SyncFromConfig(ctx, cfg)
	require.NoError(t, err)

	// Verify data was synced
	area, err := spacesStore.GetArea(ctx, "area-1")
	require.NoError(t, err)
	assert.Equal(t, "Office", area.Name)

	room, err := spacesStore.GetRoom(ctx, "room-1")
	require.NoError(t, err)
	assert.Equal(t, "Room 101", room.Name)

	desk, err := spacesStore.GetDesk(ctx, "desk-1")
	require.NoError(t, err)
	assert.Equal(t, "Desk A", desk.Name)
}

func TestStore_SyncFromConfig_Idempotent(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	cfg := &Config{
		Areas: []Area{
			{ID: "area-1", Name: "Office"},
		},
	}

	// Sync twice
	require.NoError(t, spacesStore.SyncFromConfig(ctx, cfg))
	require.NoError(t, spacesStore.SyncFromConfig(ctx, cfg))

	areas, err := spacesStore.ListAreas(ctx)
	require.NoError(t, err)
	assert.Len(t, areas, 1)
}

func TestStore_LoadConfig(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := NewStore(store)
	ctx := context.Background()

	// Seed data
	require.NoError(t, spacesStore.CreateArea(ctx, &AreaRecord{ID: "area-1", Name: "Office"}))
	require.NoError(t, spacesStore.CreateRoom(ctx, &RoomRecord{ID: "room-1", AreaID: "area-1", Name: "Room 1"}))
	require.NoError(t, spacesStore.CreateDesk(ctx, &DeskRecord{
		ID: "desk-1", RoomID: "room-1", Name: "Desk", Equipment: []string{"monitor"},
	}))

	cfg, err := spacesStore.LoadConfig(ctx)
	require.NoError(t, err)

	require.Len(t, cfg.Areas, 1)
	assert.Equal(t, "area-1", cfg.Areas[0].ID)
	assert.Equal(t, "Office", cfg.Areas[0].Name)

	require.Len(t, cfg.Areas[0].Rooms, 1)
	assert.Equal(t, "room-1", cfg.Areas[0].Rooms[0].ID)

	require.Len(t, cfg.Areas[0].Rooms[0].Desks, 1)
	assert.Equal(t, "desk-1", cfg.Areas[0].Rooms[0].Desks[0].ID)
	assert.Equal(t, []string{"monitor"}, cfg.Areas[0].Rooms[0].Desks[0].Equipment)
}
