package users

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close() //nolint:errcheck // Cleanup function, error not critical
	})

	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			display_name TEXT NOT NULL,
			password_hash TEXT NOT NULL DEFAULT '',
			user_source TEXT NOT NULL CHECK (user_source IN ('internal', 'entraid')),
			entra_id TEXT NOT NULL DEFAULT '',
			is_admin INTEGER NOT NULL DEFAULT 0,
			last_login TEXT NOT NULL DEFAULT '',
			access_token TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE UNIQUE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_users_entra_id ON users(entra_id);
	`)
	require.NoError(t, err)
	return db
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

func TestListAll(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	_, err = CreateLocalUser(ctx, db, "zach@example.com", "Zach", hash, false)
	require.NoError(t, err)
	_, err = CreateLocalUser(ctx, db, "alice@example.com", "Alice", hash, true)
	require.NoError(t, err)

	all, err := ListAll(ctx, db)
	require.NoError(t, err)
	require.Len(t, all, 2)
	// Sorted by display_name
	assert.Equal(t, "Alice", all[0].DisplayName)
	assert.Equal(t, "Zach", all[1].DisplayName)
}

func TestListAll_Empty(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	all, err := ListAll(ctx, db)
	require.NoError(t, err)
	assert.Empty(t, all)
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

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	hash, err := HashPassword("TestPassword14!")
	require.NoError(t, err)

	created, err := CreateLocalUser(ctx, db, "henry@example.com", "Henry", hash, false)
	require.NoError(t, err)

	err = DeleteUser(ctx, db, created.ID)
	require.NoError(t, err)

	_, err = FindByID(ctx, db, created.ID)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestDeleteUser_NotFound(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	err := DeleteUser(ctx, db, "nonexistent")
	assert.ErrorIs(t, err, ErrUserNotFound)
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
