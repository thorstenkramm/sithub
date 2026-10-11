package apikeys

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thorstenkramm/sithub/internal/testutil/dbtest"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		want error
	}{
		{"valid", nil}, {"malformed", ErrInvalidAPIKey}, {"unknown", ErrInvalidAPIKey},
		{"revoked", ErrInvalidAPIKey}, {"expired", ErrAPIKeyExpired}, {"boundary", ErrAPIKeyExpired},
		{"uppercase", ErrInvalidAPIKey}, {"bad timestamp", ErrInvalidAPIKey},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := dbtest.New(t)
			ctx := context.Background()
			now := time.Date(2026, 10, 8, 9, 0, 0, 0, time.UTC)
			k, raw, err := Create(ctx, db, now, CreateInput{
				Name: "Finance", Scope: ScopeReadBookings, ExpiresInMonths: 3,
			})
			require.NoError(t, err)
			switch tc.name {
			case "malformed":
				raw = "not-a-key"
			case "unknown":
				raw = uuid.NewString()
			case "revoked":
				require.NoError(t, Delete(ctx, db, k.ID))
			case "uppercase":
				raw = strings.ToUpper(raw)
			case "expired", "boundary", "bad timestamp":
				expiry := now.Format(time.RFC3339)
				if tc.name == "expired" {
					expiry = now.Add(-time.Second).Format(time.RFC3339)
				}
				if tc.name == "bad timestamp" {
					expiry = "invalid"
				}
				_, err := db.Exec(`UPDATE api_keys SET expires_at = ? WHERE id = ?`, expiry, k.ID)
				require.NoError(t, err)
			}
			principal, err := Authenticate(ctx, db, raw, now)
			if tc.want != nil {
				assert.ErrorIs(t, err, tc.want)
				assert.Nil(t, principal)
			} else {
				require.NoError(t, err)
				assert.Equal(t, k, *principal)
			}
		})
	}
}

func TestAuthenticateRejectsMalformedBeforeDatabase(t *testing.T) {
	t.Parallel()
	for _, key := range []string{"", "bad", strings.Repeat("a", 1000), " " + uuid.NewString()} {
		_, err := Authenticate(context.Background(), nil, key, time.Now())
		assert.ErrorIs(t, err, ErrInvalidAPIKey)
	}
}

func TestAuthenticateDatabaseFailure(t *testing.T) {
	t.Parallel()
	db := dbtest.New(t)
	require.NoError(t, db.Close())
	key, err := Authenticate(context.Background(), db, uuid.NewString(), time.Now())
	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrInvalidAPIKey, "storage errors must not be classified as bad credentials")
	assert.Nil(t, key)
}
