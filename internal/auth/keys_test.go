package auth

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestLoadOrCreateKeysGeneratesAndPersists(t *testing.T) {
	dir := t.TempDir()

	h1, b1, err := LoadOrCreateKeys(dir)
	require.NoError(t, err)
	require.Len(t, h1, cookieKeyLen)
	require.Len(t, b1, cookieKeyLen)

	// The key file exists and is not world/group readable (0600).
	info, err := os.Stat(filepath.Join(dir, cookieKeyFileName))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	// A second load returns the SAME keys (persistence across "restarts").
	h2, b2, err := LoadOrCreateKeys(dir)
	require.NoError(t, err)
	assert.Equal(t, h1, h2)
	assert.Equal(t, b1, b2)
}

func TestLoadOrCreateKeysMalformedFileFailsLoudly(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, cookieKeyFileName), []byte("garbage"), 0o600))

	_, _, err := LoadOrCreateKeys(dir)
	require.ErrorIs(t, err, ErrInvalidKeyFile)
}

func TestLoadOrCreateKeysWrongKeyLengthFails(t *testing.T) {
	dir := t.TempDir()
	// Two valid base64 tokens but the wrong byte length.
	require.NoError(t, os.WriteFile(filepath.Join(dir, cookieKeyFileName), []byte("YWJj\nZGVm\n"), 0o600))

	_, _, err := LoadOrCreateKeys(dir)
	require.ErrorIs(t, err, ErrInvalidKeyFile)
}

func TestLoadOrCreateKeysEmptyDirIsEphemeral(t *testing.T) {
	h, b, err := LoadOrCreateKeys("")
	require.NoError(t, err)
	require.Len(t, h, cookieKeyLen)
	require.Len(t, b, cookieKeyLen)

	// Nothing is written to the working directory.
	_, statErr := os.Stat(cookieKeyFileName)
	assert.True(t, os.IsNotExist(statErr))
}

func newDataDirConfig(dir string) *config.Config {
	cfg := &config.Config{}
	cfg.Main.DataDir = dir
	return cfg
}

// AC #3: a session cookie encoded before a restart still decodes afterwards.
func TestServiceSessionSurvivesRestart(t *testing.T) {
	dir := t.TempDir()
	cfg := newDataDirConfig(dir)

	svc1, err := NewService(cfg, nil)
	require.NoError(t, err)
	user := &User{ID: "u-1", Name: "Ada Lovelace", Email: "ada@example.com", IsAdmin: true}
	cookie, err := svc1.EncodeUser(user)
	require.NoError(t, err)

	// Simulate a restart: a brand-new Service reading the same data_dir.
	svc2, err := NewService(cfg, nil)
	require.NoError(t, err)
	got, err := svc2.DecodeUser(cookie)
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Email, got.Email)
	assert.Equal(t, user.IsAdmin, got.IsAdmin)
}

// AC #4: rotating (removing) the key file invalidates existing session cookies.
func TestServiceSessionInvalidatedOnKeyRotation(t *testing.T) {
	dir := t.TempDir()
	cfg := newDataDirConfig(dir)

	svc1, err := NewService(cfg, nil)
	require.NoError(t, err)
	cookie, err := svc1.EncodeUser(&User{ID: "u-1"})
	require.NoError(t, err)

	// Rotate the key: remove the file so the next Service generates fresh keys.
	require.NoError(t, os.Remove(filepath.Join(dir, cookieKeyFileName)))

	svc2, err := NewService(cfg, nil)
	require.NoError(t, err)
	_, err = svc2.DecodeUser(cookie)
	require.Error(t, err)
}
