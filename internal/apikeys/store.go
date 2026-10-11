// Package apikeys provides API key storage and admin management handlers.
// Package functions are stateless and safe for concurrent use.
package apikeys

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// ScopeReadBookings is the supported read-only booking scope.
const ScopeReadBookings = "read-bookings"

// PrefixLength is the number of plaintext characters retained for identification.
const PrefixLength = 10

// Validation and lookup errors allow callers to distinguish expected failures.
var (
	ErrNotFound      = errors.New("API key not found")
	ErrInvalidName   = errors.New("invalid name")
	ErrInvalidScope  = errors.New("scope must be read-bookings")
	ErrInvalidExpiry = errors.New("expires_in_months must be one of 3, 12, 24")
)

// Key contains public metadata only. It is safe for concurrent read-only use.
type Key struct {
	ID        string
	Name      string
	KeyPrefix string
	Scope     string
	ExpiresAt string
	CreatedAt string
	CreatedBy string
}

// CreateInput describes a new key. It is safe for concurrent read-only use.
type CreateInput struct {
	Name            string
	Scope           string
	ExpiresInMonths int
	CreatedBy       string
}

// HashKey returns the lowercase SHA-256 hex digest of the exact plaintext key.
func HashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// validate checks the supported scope, expiry and trimmed name before key generation.
func validate(in CreateInput) error {
	if in.Name == "" {
		return fmt.Errorf("name is required: %w", ErrInvalidName)
	}
	if utf8.RuneCountInString(in.Name) > 100 {
		return fmt.Errorf("name must be at most 100 characters: %w", ErrInvalidName)
	}
	if in.Scope != ScopeReadBookings {
		return ErrInvalidScope
	}
	switch in.ExpiresInMonths {
	case 3, 12, 24:
		return nil
	default:
		return ErrInvalidExpiry
	}
}

// Create stores only a hash and prefix, returning the plaintext key once to its caller.
func Create(ctx context.Context, db *sql.DB, now time.Time, in CreateInput) (Key, string, error) {
	in.Name = strings.TrimSpace(in.Name)
	if err := validate(in); err != nil {
		return Key{}, "", err
	}
	plaintext := uuid.NewString()
	k := Key{
		ID: uuid.NewString(), Name: in.Name, KeyPrefix: plaintext[:PrefixLength], Scope: in.Scope,
		ExpiresAt: now.UTC().AddDate(0, in.ExpiresInMonths, 0).Format(time.RFC3339),
		CreatedAt: now.UTC().Format(time.RFC3339), CreatedBy: in.CreatedBy,
	}
	_, err := db.ExecContext(ctx, `INSERT INTO api_keys
		(id, name, key_prefix, key_hash, scope, expires_at, created_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		k.ID, k.Name, k.KeyPrefix, HashKey(plaintext), k.Scope, k.ExpiresAt, k.CreatedAt, k.CreatedBy)
	if err != nil {
		return Key{}, "", fmt.Errorf("create API key: %w", err)
	}
	return k, plaintext, nil
}

// List returns public metadata, newest first, without reading any hashes.
func List(ctx context.Context, db *sql.DB) ([]Key, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, name, key_prefix, scope, expires_at, created_at, created_by
		FROM api_keys ORDER BY created_at DESC, name`)
	if err != nil {
		return nil, fmt.Errorf("list API keys: %w", err)
	}
	defer func() {
		_ = rows.Close() //nolint:errcheck // Read-only rows; iteration errors are checked below
	}()
	keys := make([]Key, 0)
	for rows.Next() {
		var k Key
		if err := rows.Scan(&k.ID, &k.Name, &k.KeyPrefix, &k.Scope,
			&k.ExpiresAt, &k.CreatedAt, &k.CreatedBy); err != nil {
			return nil, fmt.Errorf("scan API key: %w", err)
		}
		keys = append(keys, k)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate API keys: %w", err)
	}
	return keys, nil
}

// Delete permanently revokes a key by its public record ID.
func Delete(ctx context.Context, db *sql.DB, id string) error {
	result, err := db.ExecContext(ctx, `DELETE FROM api_keys WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete API key: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count deleted API keys: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
