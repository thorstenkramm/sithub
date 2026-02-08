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
	spacePath := writeSpaceConfig(t)
	path := writeConfig(t, `
[entraid]
authorize_url = "https://login"
token_url = "https://token"
redirect_uri = "http://localhost/callback"
client_id = "client"
client_secret = "secret"

[spaces]
config_file = "`+spacePath+`"
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
}

func TestLoadMissingEntraID(t *testing.T) {
	spacePath := writeSpaceConfig(t)
	path := writeConfig(t, `
[entraid]
authorize_url = ""
token_url = ""
redirect_uri = ""

[spaces]
config_file = "`+spacePath+`"
`)

	_, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error for empty Entra ID (local-only mode), got %v", err)
	}
}

func TestLoadPartialEntraID(t *testing.T) {
	spacePath := writeSpaceConfig(t)
	path := writeConfig(t, `
[entraid]
authorize_url = "https://login"
token_url = ""
redirect_uri = ""

[spaces]
config_file = "`+spacePath+`"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatalf("expected error for partial Entra ID config")
	}
}

func TestLoadNoEntraID_LocalOnly(t *testing.T) {
	spacePath := writeSpaceConfig(t)
	path := writeConfig(t, `
[spaces]
config_file = "`+spacePath+`"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error for local-only mode, got %v", err)
	}
	if cfg.EntraIDConfigured() {
		t.Fatalf("expected EntraIDConfigured=false for local-only mode")
	}
}

func TestEntraIDConfigured(t *testing.T) {
	spacePath := writeSpaceConfig(t)
	path := writeConfig(t, `
[entraid]
authorize_url = "https://login"
token_url = "https://token"
redirect_uri = "http://localhost/callback"
client_id = "client"
client_secret = "secret"

[spaces]
config_file = "`+spacePath+`"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.EntraIDConfigured() {
		t.Fatalf("expected EntraIDConfigured=true")
	}
}

func TestLoadMissingSpacesConfig(t *testing.T) {
	path := writeConfig(t, `
[entraid]
authorize_url = "https://login"
token_url = "https://token"
redirect_uri = "http://localhost/callback"
client_id = "client"
client_secret = "secret"
`)

	if _, err := Load(path); err == nil {
		t.Fatalf("expected error for missing spaces config")
	}
}

func writeSpaceConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "spaces.yaml")
	if err := os.WriteFile(path, []byte("areas: []\n"), 0o600); err != nil {
		t.Fatalf("write spaces config: %v", err)
	}
	return path
}
