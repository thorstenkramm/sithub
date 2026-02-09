package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
	userColumns = `id, email, display_name, password_hash,
		user_source, entra_id, is_admin, last_login,
		created_at, updated_at`
)

// FindByID returns a user by primary key, or ErrUserNotFound.
func FindByID(ctx context.Context, db *sql.DB, id string) (*Record, error) {
	return scanOne(db.QueryRowContext(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = ?`, id))
}

// FindByEmail returns a user by email address, or ErrUserNotFound.
func FindByEmail(ctx context.Context, db *sql.DB, email string) (*Record, error) {
	return scanOne(db.QueryRowContext(ctx,
		`SELECT `+userColumns+` FROM users WHERE email = ?`, email))
}

// FindByEntraID returns a user by their Entra ID object ID, or ErrUserNotFound.
func FindByEntraID(ctx context.Context, db *sql.DB, entraID string) (*Record, error) {
	return scanOne(db.QueryRowContext(ctx,
		`SELECT `+userColumns+` FROM users WHERE entra_id = ?`, entraID))
}

// UpsertEntraIDUser inserts or updates a user from Entra ID login.
// On conflict (same entra_id), display_name and email are updated.
func UpsertEntraIDUser(
	ctx context.Context, db *sql.DB, entraID, email, displayName string, isAdmin bool,
) (*Record, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.New().String()

	isAdminInt := 0
	if isAdmin {
		isAdminInt = 1
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO users (id, email, display_name, password_hash, user_source, entra_id,
			is_admin, last_login, created_at, updated_at)
		VALUES (?, ?, ?, '', 'entraid', ?, ?, ?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET
			display_name = excluded.display_name,
			entra_id = excluded.entra_id,
			is_admin = excluded.is_admin,
			last_login = excluded.last_login,
			updated_at = excluded.updated_at`,
		id, email, displayName, entraID, isAdminInt, now, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert entra user: %w", err)
	}

	return FindByEmail(ctx, db, email)
}

// CreateLocalUser creates a new user with internal (local) authentication.
func CreateLocalUser(
	ctx context.Context, db *sql.DB, email, displayName, passwordHash string, isAdmin bool,
) (*Record, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.New().String()

	isAdminInt := 0
	if isAdmin {
		isAdminInt = 1
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO users (id, email, display_name, password_hash, user_source, entra_id,
			is_admin, last_login, created_at, updated_at)
		VALUES (?, ?, ?, ?, 'internal', '', ?, '', ?, ?)`,
		id, email, displayName, passwordHash, isAdminInt, now, now,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, ErrEmailConflict
		}
		return nil, fmt.Errorf("create local user: %w", err)
	}

	return &Record{
		ID:           id,
		Email:        email,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
		UserSource:   "internal",
		IsAdmin:      isAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// ListAll returns all users ordered by display_name.
func ListAll(ctx context.Context, db *sql.DB) ([]Record, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT `+userColumns+` FROM users ORDER BY display_name`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer func() {
		_ = rows.Close() //nolint:errcheck // Defer close, error not critical
	}()

	var result []Record
	for rows.Next() {
		rec, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}
	return result, nil
}

// UpdateFields contains optional fields to update on a user.
type UpdateFields struct {
	Email       *string
	DisplayName *string
	IsAdmin     *bool
}

// UpdateUser applies partial updates to a user.
func UpdateUser(ctx context.Context, db *sql.DB, id string, fields UpdateFields) (*Record, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	// Build dynamic SET clause
	setClauses := []string{"updated_at = ?"}
	args := []interface{}{now}

	if fields.Email != nil {
		setClauses = append(setClauses, "email = ?")
		args = append(args, *fields.Email)
	}
	if fields.DisplayName != nil {
		setClauses = append(setClauses, "display_name = ?")
		args = append(args, *fields.DisplayName)
	}
	if fields.IsAdmin != nil {
		isAdminInt := 0
		if *fields.IsAdmin {
			isAdminInt = 1
		}
		setClauses = append(setClauses, "is_admin = ?")
		args = append(args, isAdminInt)
	}

	args = append(args, id)

	query := "UPDATE users SET "
	for i, clause := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += clause
	}
	query += " WHERE id = ?"

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, ErrEmailConflict
		}
		return nil, fmt.Errorf("update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update user rows: %w", err)
	}
	if rows == 0 {
		return nil, ErrUserNotFound
	}

	return FindByID(ctx, db, id)
}

// DeleteUser removes a user by ID. Returns ErrUserNotFound if no row was deleted.
func DeleteUser(ctx context.Context, db *sql.DB, id string) error {
	result, err := db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user rows: %w", err)
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

// UpdatePasswordHash sets a new password hash for a user.
func UpdatePasswordHash(ctx context.Context, db *sql.DB, id, hash string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := db.ExecContext(ctx,
		"UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?",
		hash, now, id,
	)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update password rows: %w", err)
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

// UpdateLastLogin sets the last_login timestamp for a user to the current time.
func UpdateLastLogin(ctx context.Context, db *sql.DB, id string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := db.ExecContext(ctx,
		"UPDATE users SET last_login = ? WHERE id = ?",
		now, id,
	)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update last login rows: %w", err)
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

// UpdateAccessToken stores the Entra ID access token for a user.
//
// SECURITY NOTE: The token is stored as plaintext in SQLite. This is acceptable
// for the current threat model (single-server deployment, DB file access implies
// full system compromise). If the deployment model changes (e.g., shared DB host),
// consider encrypting tokens at rest with a server-side key.
func UpdateAccessToken(ctx context.Context, db *sql.DB, id, token string) error {
	result, err := db.ExecContext(ctx,
		"UPDATE users SET access_token = ? WHERE id = ?",
		token, id,
	)
	if err != nil {
		return fmt.Errorf("update access token: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update access token rows: %w", err)
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

// GetAccessToken retrieves the stored Entra ID access token for a user.
func GetAccessToken(ctx context.Context, db *sql.DB, id string) (string, error) {
	var token string
	err := db.QueryRowContext(ctx,
		"SELECT access_token FROM users WHERE id = ?", id,
	).Scan(&token)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrUserNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get access token: %w", err)
	}
	return token, nil
}

// FindDisplayNames returns a map of user IDs to display names.
// Unknown IDs are silently omitted from the result.
func FindDisplayNames(ctx context.Context, db *sql.DB, userIDs []string) (result map[string]string, err error) {
	if len(userIDs) == 0 {
		return map[string]string{}, nil
	}

	// Deduplicate
	seen := make(map[string]struct{}, len(userIDs))
	unique := make([]string, 0, len(userIDs))
	for _, id := range userIDs {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			unique = append(unique, id)
		}
	}

	placeholders := make([]string, len(unique))
	args := make([]interface{}, len(unique))
	for i, id := range unique {
		placeholders[i] = "?"
		args[i] = id
	}

	//nolint:gosec // G201: placeholders are "?" literals, not user input
	query := fmt.Sprintf(
		"SELECT id, display_name FROM users WHERE id IN (%s)",
		strings.Join(placeholders, ","),
	)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query display names: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close display names rows: %w", closeErr)
		}
	}()

	result = make(map[string]string, len(unique))
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("scan display name: %w", err)
		}
		result[id] = name
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate display names: %w", err)
	}

	return result, nil
}

// VerifyPassword compares a bcrypt hash with a plaintext password.
func VerifyPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("verify password: %w", err)
	}
	return nil
}

// HashPassword generates a bcrypt hash from a plaintext password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

func scanOne(row *sql.Row) (*Record, error) {
	var rec Record
	var isAdminInt int
	err := row.Scan(
		&rec.ID, &rec.Email, &rec.DisplayName, &rec.PasswordHash,
		&rec.UserSource, &rec.EntraID, &isAdminInt, &rec.LastLogin,
		&rec.CreatedAt, &rec.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	rec.IsAdmin = isAdminInt == 1
	return &rec, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanRow(row rowScanner) (*Record, error) {
	var rec Record
	var isAdminInt int
	err := row.Scan(
		&rec.ID, &rec.Email, &rec.DisplayName, &rec.PasswordHash,
		&rec.UserSource, &rec.EntraID, &isAdminInt, &rec.LastLogin,
		&rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan user row: %w", err)
	}
	rec.IsAdmin = isAdminInt == 1
	return &rec, nil
}
