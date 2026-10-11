package users

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/testutil/dbtest"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	return dbtest.New(t)
}

func TestCreateLocalUser(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	rec, err := CreateLocalUser(ctx, db, "alice@example.com", "Alice", hash, false)
	require.NoError(t, err)

	assert.NotEmpty(t, rec.ID)
	assert.Equal(t, "alice@example.com", rec.Email)
	assert.Equal(t, "Alice", rec.DisplayName)
	assert.Equal(t, "internal", rec.UserSource)
	assert.False(t, rec.IsAdmin)
	assert.NotEmpty(t, rec.CreatedAt)
}

func TestCreateLocalUser_DuplicateEmail(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	_, err = CreateLocalUser(ctx, db, "alice@example.com", "Alice", hash, false)
	require.NoError(t, err)

	_, err = CreateLocalUser(ctx, db, "alice@example.com", "Alice2", hash, false)
	assert.ErrorIs(t, err, ErrEmailConflict)
}

func TestCreateLocalUser_Admin(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	rec, err := CreateLocalUser(ctx, db, "admin@example.com", "Admin", hash, true)
	require.NoError(t, err)
	assert.True(t, rec.IsAdmin)
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	created, err := CreateLocalUser(ctx, db, "bob@example.com", "Bob", hash, false)
	require.NoError(t, err)

	found, err := FindByID(ctx, db, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, "Bob", found.DisplayName)
}

func TestFindByID_NotFound(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	_, err := FindByID(ctx, db, "nonexistent")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestFindByEmail(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	_, err = CreateLocalUser(ctx, db, "carol@example.com", "Carol", hash, false)
	require.NoError(t, err)

	found, err := FindByEmail(ctx, db, "carol@example.com")
	require.NoError(t, err)
	assert.Equal(t, "Carol", found.DisplayName)
}

func TestFindByEmail_NotFound(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	_, err := FindByEmail(ctx, db, "nobody@example.com")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestFindDisplayNames_DoesNotInjectViaIDs(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	first, err := CreateLocalUser(ctx, db, "alice@example.com", "Alice", hash, false)
	require.NoError(t, err)
	_, err = CreateLocalUser(ctx, db, "bob@example.com", "Bob", hash, false)
	require.NoError(t, err)

	names, err := FindDisplayNames(ctx, db, []string{first.ID, `x') OR 1=1 --`})
	require.NoError(t, err)
	require.Len(t, names, 1)
	assert.Equal(t, "Alice", names[first.ID])
	assert.NotContains(t, names, `x') OR 1=1 --`)
}

func TestUpsertEntraIDUser_Insert(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	rec, err := UpsertEntraIDUser(ctx, db, "entra-123", "dave@example.com", "Dave", false)
	require.NoError(t, err)
	assert.Equal(t, "dave@example.com", rec.Email)
	assert.Equal(t, "entraid", rec.UserSource)
	assert.Equal(t, "entra-123", rec.EntraID)
}

func TestUpsertEntraIDUser_Update(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	_, err := UpsertEntraIDUser(ctx, db, "entra-456", "eve@example.com", "Eve Original", false)
	require.NoError(t, err)

	updated, err := UpsertEntraIDUser(ctx, db, "entra-456", "eve@example.com", "Eve Updated", true)
	require.NoError(t, err)
	assert.Equal(t, "Eve Updated", updated.DisplayName)
	assert.True(t, updated.IsAdmin)
}

func TestFindByEntraID(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	_, err := UpsertEntraIDUser(ctx, db, "entra-789", "frank@example.com", "Frank", false)
	require.NoError(t, err)

	found, err := FindByEntraID(ctx, db, "entra-789")
	require.NoError(t, err)
	assert.Equal(t, "Frank", found.DisplayName)
}

func TestFindByEntraID_NotFound(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	_, err := FindByEntraID(ctx, db, "nonexistent")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

// insertUserRow inserts a user directly, bypassing bcrypt, so ordering tests stay fast.
func insertUserRow(t *testing.T, db *sql.DB, id, displayName string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO users (id, email, display_name, user_source, created_at, updated_at)
		VALUES (?, ?, ?, 'internal', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`,
		id, id+"@example.com", displayName)
	require.NoError(t, err)
}

func TestListPage(t *testing.T) {
	t.Parallel()

	// Insertion order differs from the expected sort order on purpose.
	seed := []struct{ id, name string }{
		{"u5", "zach"}, {"u2", "Bob"}, {"u4", "bob"}, {"u1", "alice"}, {"u3", "Bob"},
	}
	// Case-insensitive by name, ties broken by id.
	sorted := []string{"u1", "u2", "u3", "u4", "u5"}

	cases := []struct {
		name   string
		limit  int
		offset int
		want   []string
	}{
		{name: "first page", limit: 2, offset: 0, want: sorted[0:2]},
		{name: "middle page with name ties", limit: 2, offset: 2, want: sorted[2:4]},
		{name: "last partial page", limit: 2, offset: 4, want: sorted[4:5]},
		{name: "offset beyond total", limit: 2, offset: 10, want: []string{}},
		{name: "all rows", limit: 100, offset: 0, want: sorted},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := setupTestDB(t)
			for _, u := range seed {
				insertUserRow(t, db, u.id, u.name)
			}

			records, err := ListPage(context.Background(), db, tc.limit, tc.offset)
			require.NoError(t, err)
			require.NotNil(t, records)

			ids := make([]string, 0, len(records))
			for _, rec := range records {
				ids = append(ids, rec.ID)
			}
			assert.Equal(t, tc.want, ids)
		})
	}
}

func TestListPage_Empty(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	records, err := ListPage(context.Background(), db, 20, 0)
	require.NoError(t, err)
	assert.NotNil(t, records)
	assert.Empty(t, records)
}

func TestCount(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	total, err := Count(ctx, db)
	require.NoError(t, err)
	assert.Equal(t, 0, total)

	insertUserRow(t, db, "u1", "Alice")
	insertUserRow(t, db, "u2", "Bob")

	total, err = Count(ctx, db)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	created, err := CreateLocalUser(ctx, db, "grace@example.com", "Grace", hash, false)
	require.NoError(t, err)

	newName := "Grace Updated"
	newAdmin := true
	updated, err := UpdateUser(ctx, db, created.ID, UpdateFields{
		DisplayName: &newName,
		IsAdmin:     &newAdmin,
	})
	require.NoError(t, err)
	assert.Equal(t, "Grace Updated", updated.DisplayName)
	assert.True(t, updated.IsAdmin)
}

func TestUpdateUser_EmailConflict(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	_, err = CreateLocalUser(ctx, db, "existing@example.com", "Existing", hash, false)
	require.NoError(t, err)

	created, err := CreateLocalUser(ctx, db, "other@example.com", "Other", hash, false)
	require.NoError(t, err)

	conflict := "existing@example.com"
	_, err = UpdateUser(ctx, db, created.ID, UpdateFields{Email: &conflict})
	assert.ErrorIs(t, err, ErrEmailConflict)
}

func TestUpdateUser_NotFound(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	newName := "Nobody"
	_, err := UpdateUser(ctx, db, "nonexistent", UpdateFields{DisplayName: &newName})
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestUpdateUser_TreatsFieldValuesAsData(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	target, err := CreateLocalUser(ctx, db, "ivy@example.com", "Ivy", hash, false)
	require.NoError(t, err)
	other, err := CreateLocalUser(ctx, db, "jules@example.com", "Jules", hash, false)
	require.NoError(t, err)

	maliciousEmail := `ivy@example.com', is_admin = 1 WHERE id = '` + other.ID
	updated, err := UpdateUser(ctx, db, target.ID, UpdateFields{Email: &maliciousEmail})
	require.NoError(t, err)

	assert.Equal(t, maliciousEmail, updated.Email)
	assert.False(t, updated.IsAdmin)

	otherFound, err := FindByID(ctx, db, other.ID)
	require.NoError(t, err)
	assert.Equal(t, "jules@example.com", otherFound.Email)
	assert.False(t, otherFound.IsAdmin)
}

func TestUpdateUser_TreatsIDAsData(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	created, err := CreateLocalUser(ctx, db, "karl@example.com", "Karl", hash, false)
	require.NoError(t, err)

	newName := "Injected"
	_, err = UpdateUser(ctx, db, created.ID+`' OR 1=1 --`, UpdateFields{DisplayName: &newName})
	assert.ErrorIs(t, err, ErrUserNotFound)

	found, err := FindByID(ctx, db, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Karl", found.DisplayName)
}

func TestUpdatePasswordHash(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	oldHash, err := HashPassword("OldPassword14!!")
	require.NoError(t, err)

	created, err := CreateLocalUser(ctx, db, "iris@example.com", "Iris", oldHash, false)
	require.NoError(t, err)

	newHash, err := HashPassword("NewPassword14!!")
	require.NoError(t, err)

	err = UpdatePasswordHash(ctx, db, created.ID, newHash)
	require.NoError(t, err)

	found, err := FindByID(ctx, db, created.ID)
	require.NoError(t, err)
	assert.Equal(t, newHash, found.PasswordHash)
}

func TestUpdatePasswordHash_NotFound(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	err := UpdatePasswordHash(ctx, db, "nonexistent", "hash")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestHashAndVerifyPassword(t *testing.T) {
	t.Parallel()

	hash, err := HashPassword("MySecretPassword!")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = VerifyPassword(hash, "MySecretPassword!")
	assert.NoError(t, err)

	err = VerifyPassword(hash, "WrongPassword")
	assert.Error(t, err)
}
