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

func TestCreate(t *testing.T) {
	t.Parallel()
	for _, months := range []int{3, 12, 24} {
		t.Run(time.Month(months).String(), func(t *testing.T) {
			t.Parallel()
			db := dbtest.New(t)
			now := time.Date(2026, 11, 30, 9, 15, 0, 0, time.FixedZone("offset", 3600))
			k, plaintext, err := Create(context.Background(), db, now, CreateInput{
				Name: " Export ", Scope: ScopeReadBookings, ExpiresInMonths: months, CreatedBy: "admin",
			})
			require.NoError(t, err)
			parsed, err := uuid.Parse(plaintext)
			require.NoError(t, err)
			assert.Equal(t, uuid.Version(4), parsed.Version())
			assert.NotEqual(t, k.ID, plaintext)
			assert.Equal(t, "Export", k.Name)
			assert.Equal(t, plaintext[:10], k.KeyPrefix)
			assert.Equal(t, now.UTC().AddDate(0, months, 0).Format(time.RFC3339), k.ExpiresAt)
			if months == 3 {
				assert.Equal(t, "2027-03-02T08:15:00Z", k.ExpiresAt)
			}
			var values [8]string
			require.NoError(t, db.QueryRow(`SELECT * FROM api_keys WHERE id = ?`, k.ID).Scan(
				&values[0], &values[1], &values[2], &values[3], &values[4], &values[5], &values[6], &values[7]))
			assert.Equal(t, HashKey(plaintext), values[3])
			assert.Regexp(t, "^[0-9a-f]{64}$", values[3])
			for _, value := range values {
				assert.NotContains(t, value, plaintext)
			}
		})
	}
}

func TestCreateValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, scope string
		months      int
		want        error
	}{
		{"", ScopeReadBookings, 12, ErrInvalidName},
		{" \t ", ScopeReadBookings, 12, ErrInvalidName},
		{strings.Repeat("a", 101), ScopeReadBookings, 12, ErrInvalidName},
		{strings.Repeat("ü", 101), ScopeReadBookings, 12, ErrInvalidName},
		{"Name", "write", 12, ErrInvalidScope},
		{"Name", ScopeReadBookings, 0, ErrInvalidExpiry},
		{"Name", ScopeReadBookings, 1, ErrInvalidExpiry},
		{"Name", ScopeReadBookings, 6, ErrInvalidExpiry},
		{"Name", ScopeReadBookings, 36, ErrInvalidExpiry},
		{"Name", ScopeReadBookings, -3, ErrInvalidExpiry},
	} {
		t.Run(tc.name+tc.scope+time.Month(tc.months).String(), func(t *testing.T) {
			t.Parallel()
			_, plaintext, err := Create(context.Background(), nil, time.Now(), CreateInput{
				Name: tc.name, Scope: tc.scope, ExpiresInMonths: tc.months,
			})
			assert.ErrorIs(t, err, tc.want)
			assert.Empty(t, plaintext)
		})
	}
}

func TestListAndDelete(t *testing.T) {
	t.Parallel()
	db := dbtest.New(t)
	ctx := context.Background()
	keys, err := List(ctx, db)
	require.NoError(t, err)
	assert.Empty(t, keys)
	now := time.Now()
	for _, name := range []string{"Z", "B", "A"} {
		created := now
		if name == "Z" {
			created = now.Add(-time.Hour)
		}
		_, _, err := Create(ctx, db, created, CreateInput{Name: name, Scope: ScopeReadBookings, ExpiresInMonths: 12})
		require.NoError(t, err)
	}
	keys, err = List(ctx, db)
	require.NoError(t, err)
	require.Len(t, keys, 3)
	assert.Equal(t, []string{"A", "B", "Z"}, []string{keys[0].Name, keys[1].Name, keys[2].Name})
	require.NoError(t, Delete(ctx, db, keys[0].ID))
	assert.ErrorIs(t, Delete(ctx, db, keys[0].ID), ErrNotFound)
	keys, err = List(ctx, db)
	require.NoError(t, err)
	assert.Len(t, keys, 2)
}

func TestHashKey(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", HashKey("abc"))
}
