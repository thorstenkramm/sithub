package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	cookieKeyFileName = "cookie.key"
	cookieKeyLen      = 32
)

// ErrInvalidKeyFile indicates the persisted cookie-key file is malformed
// (wrong number of keys, not valid base64, or wrong key length).
var ErrInvalidKeyFile = errors.New("invalid cookie key file")

// LoadOrCreateKeys returns the 32-byte hash and block keys used to sign and
// encrypt session cookies.
//
// Security-critical: these keys protect the sithub_user session cookie and the
// sithub_oauth_state cookie. When dataDir is set, they are persisted to
// {dataDir}/cookie.key (mode 0600) so sessions survive server restarts: on the
// first start the file is created with freshly generated random keys and
// reused on every later start.
//
// Key rotation / invalidation: deleting or replacing cookie.key causes a new
// key pair to be generated on the next start, which invalidates every
// outstanding session and OAuth-state cookie (users are transparently
// redirected to login on their next request). A corrupt or truncated file is a
// hard error (ErrInvalidKeyFile) rather than a silent regenerate, so operators
// never unknowingly log out their whole user base.
//
// If dataDir is empty the keys are generated in memory and NOT persisted
// (ephemeral). This only happens when no data_dir is configured; the server
// default is "." (see config), so production always persists.
func LoadOrCreateKeys(dataDir string) (hashKey, blockKey []byte, err error) {
	if dataDir == "" {
		return newRandomKeyPair()
	}

	path := filepath.Join(dataDir, cookieKeyFileName)
	// #nosec G304 -- path is {dataDir}/cookie.key from trusted config, not user input
	data, readErr := os.ReadFile(path)
	switch {
	case readErr == nil:
		return parseKeyFile(data)
	case errors.Is(readErr, os.ErrNotExist):
		return generateAndPersistKeys(dataDir, path)
	default:
		return nil, nil, fmt.Errorf("read cookie key file: %w", readErr)
	}
}

func newRandomKeyPair() (hashKey, blockKey []byte, err error) {
	hashKey = make([]byte, cookieKeyLen)
	blockKey = make([]byte, cookieKeyLen)
	if _, err = io.ReadFull(rand.Reader, hashKey); err != nil {
		return nil, nil, fmt.Errorf("generate hash key: %w", err)
	}
	if _, err = io.ReadFull(rand.Reader, blockKey); err != nil {
		return nil, nil, fmt.Errorf("generate block key: %w", err)
	}
	return hashKey, blockKey, nil
}

// generateAndPersistKeys creates a fresh key pair and writes it to path as two
// base64 lines with 0600 perms, creating dataDir (0750) if needed.
func generateAndPersistKeys(dataDir, path string) (hashKey, blockKey []byte, err error) {
	hashKey, blockKey, err = newRandomKeyPair()
	if err != nil {
		return nil, nil, err
	}
	if err = os.MkdirAll(dataDir, 0o750); err != nil {
		return nil, nil, fmt.Errorf("create data dir: %w", err)
	}
	content := base64.StdEncoding.EncodeToString(hashKey) + "\n" +
		base64.StdEncoding.EncodeToString(blockKey) + "\n"
	if err = os.WriteFile(path, []byte(content), 0o600); err != nil {
		return nil, nil, fmt.Errorf("write cookie key file: %w", err)
	}
	return hashKey, blockKey, nil
}

// parseKeyFile validates that data holds exactly two 32-byte base64-encoded
// keys, returning ErrInvalidKeyFile on any deviation.
func parseKeyFile(data []byte) (hashKey, blockKey []byte, err error) {
	fields := strings.Fields(string(data))
	if len(fields) != 2 {
		return nil, nil, fmt.Errorf("%w: expected 2 keys, got %d", ErrInvalidKeyFile, len(fields))
	}
	if hashKey, err = decodeKey(fields[0]); err != nil {
		return nil, nil, err
	}
	if blockKey, err = decodeKey(fields[1]); err != nil {
		return nil, nil, err
	}
	return hashKey, blockKey, nil
}

func decodeKey(s string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("%w: base64 decode failed", ErrInvalidKeyFile)
	}
	if len(key) != cookieKeyLen {
		return nil, fmt.Errorf("%w: key length %d, want %d", ErrInvalidKeyFile, len(key), cookieKeyLen)
	}
	return key, nil
}
