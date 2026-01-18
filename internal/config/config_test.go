package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "sithub.toml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestLoadDefaults(t *testing.T) {
	path := writeConfig(t, `
[entraid]
authorize_url = "https://login"
token_url = "https://token"
redirect_uri = "http://localhost/callback"
client_id = "client"
client_secret = "secret"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Main.Port != 9900 {
		t.Fatalf("expected default port 9900, got %d", cfg.Main.Port)
	}
	if cfg.Main.Listen != "127.0.0.1" {
		t.Fatalf("expected default listen, got %s", cfg.Main.Listen)
	}
	if cfg.TestAuth.Enabled {
		t.Fatalf("expected test auth disabled by default")
	}
	if cfg.TestAuth.UserID != "test-user" {
		t.Fatalf("expected default test auth user id, got %s", cfg.TestAuth.UserID)
	}
	if cfg.TestAuth.UserName != "Test User" {
		t.Fatalf("expected default test auth user name, got %s", cfg.TestAuth.UserName)
	}
}

func TestLoadMissingEntraID(t *testing.T) {
	path := writeConfig(t, `
[entraid]
authorize_url = ""
token_url = ""
redirect_uri = ""
`)

	_, err := Load(path)
	if err == nil {
		t.Fatalf("expected error for missing Entra ID config")
	}
}

func TestLoadMissingEntraIDWithTestAuth(t *testing.T) {
	path := writeConfig(t, `
[test_auth]
enabled = true
`)

	if _, err := Load(path); err != nil {
		t.Fatalf("expected no error with test auth enabled, got %v", err)
	}
}
