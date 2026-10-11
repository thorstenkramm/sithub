package apikeys

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ErrInvalidAPIKey means the credential is malformed, unknown or revoked.
var ErrInvalidAPIKey = errors.New("invalid API key")

// ErrAPIKeyExpired means a matching credential has reached its expiry instant.
var ErrAPIKeyExpired = errors.New("API key has expired")

// Authenticate is security-critical: validate before accessing storage, hash the exact
// input without normalization, look up by its unique hash and confirm in constant time.
// Only unexpired keys authenticate. Returned metadata never contains the hash or secret.
func Authenticate(ctx context.Context, store *sql.DB, rawKey string, now time.Time) (*Key, error) {
	if len(rawKey) != 36 {
		return nil, ErrInvalidAPIKey
	}
	if _, err := uuid.Parse(rawKey); err != nil {
		return nil, ErrInvalidAPIKey
	}
	computed := HashKey(rawKey)
	var k Key
	var stored string
	err := store.QueryRowContext(ctx, `SELECT id, name, key_prefix, scope, expires_at, created_at, created_by, key_hash
		FROM api_keys WHERE key_hash = ?`, computed).Scan(
		&k.ID, &k.Name, &k.KeyPrefix, &k.Scope, &k.ExpiresAt, &k.CreatedAt, &k.CreatedBy, &stored)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidAPIKey
	}
	if err != nil {
		return nil, fmt.Errorf("authenticate API key: %w", err)
	}
	if subtle.ConstantTimeCompare([]byte(stored), []byte(computed)) != 1 {
		return nil, ErrInvalidAPIKey
	}
	expiresAt, err := time.Parse(time.RFC3339, k.ExpiresAt)
	if err != nil {
		return nil, ErrInvalidAPIKey
	}
	if !now.UTC().Before(expiresAt) {
		return nil, ErrAPIKeyExpired
	}
	return &k, nil
}
