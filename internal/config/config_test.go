package config

import (
	"errors"
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

// writeAreasConfigIn creates an areas YAML file inside a specific directory.
func writeAreasConfigIn(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "areas.yaml")
	if err := os.WriteFile(path, []byte("areas: []\n"), 0o600); err != nil {
		t.Fatalf("write areas config: %v", err)
	}
	return path
}

func TestLoadDefaults(t *testing.T) {
	dataDir := t.TempDir()
	areasPath := writeAreasConfigIn(t, dataDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[entraid]
authorize_url = "https://login"
token_url = "https://token"
redirect_uri = "http://localhost/callback"
client_id = "client"
client_secret = "secret"

[areas]
config_file = "`+areasPath+`"
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
	dataDir := t.TempDir()
	areasPath := writeAreasConfigIn(t, dataDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[entraid]
authorize_url = ""
token_url = ""
redirect_uri = ""

[areas]
config_file = "`+areasPath+`"
`)

	_, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error for empty Entra ID (local-only mode), got %v", err)
	}
}

func TestLoadPartialEntraID(t *testing.T) {
	dataDir := t.TempDir()
	areasPath := writeAreasConfigIn(t, dataDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[entraid]
authorize_url = "https://login"
token_url = ""
redirect_uri = ""

[areas]
config_file = "`+areasPath+`"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatalf("expected error for partial Entra ID config")
	}
}

func TestLoadNoEntraID_LocalOnly(t *testing.T) {
	dataDir := t.TempDir()
	areasPath := writeAreasConfigIn(t, dataDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[areas]
config_file = "`+areasPath+`"
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
	dataDir := t.TempDir()
	areasPath := writeAreasConfigIn(t, dataDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[entraid]
authorize_url = "https://login"
token_url = "https://token"
redirect_uri = "http://localhost/callback"
client_id = "client"
client_secret = "secret"

[areas]
config_file = "`+areasPath+`"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.EntraIDConfigured() {
		t.Fatalf("expected EntraIDConfigured=true")
	}
}

func TestLoadMissingAreasConfig(t *testing.T) {
	path := writeConfig(t, `
[entraid]
authorize_url = "https://login"
token_url = "https://token"
redirect_uri = "http://localhost/callback"
client_id = "client"
client_secret = "secret"
`)

	if _, err := Load(path); err == nil {
		t.Fatalf("expected error for missing areas config")
	}
}

func TestLoadRelativeAreasConfig(t *testing.T) {
	dataDir := t.TempDir()
	writeAreasConfigIn(t, dataDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[areas]
config_file = "areas.yaml"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected relative path to resolve inside data_dir, got %v", err)
	}
	expected := filepath.Join(dataDir, "areas.yaml")
	if cfg.Areas.ConfigFile != expected {
		t.Fatalf("expected resolved path %s, got %s", expected, cfg.Areas.ConfigFile)
	}
}

func TestLoadFloorPlansDirRelative(t *testing.T) {
	dataDir := t.TempDir()
	writeAreasConfigIn(t, dataDir)
	fpDir := filepath.Join(dataDir, "floor_plans")
	if err := os.MkdirAll(fpDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[areas]
config_file = "areas.yaml"
floor_plans = "floor_plans"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Areas.FloorPlansDir != fpDir {
		t.Fatalf("expected %s, got %s", fpDir, cfg.Areas.FloorPlansDir)
	}
}

func TestLoadFloorPlansDirNotExist(t *testing.T) {
	dataDir := t.TempDir()
	writeAreasConfigIn(t, dataDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[areas]
config_file = "areas.yaml"
floor_plans = "nonexistent"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatalf("expected error for missing floor plans dir")
	}
	if !errors.Is(err, ErrFloorPlansDirNotFound) {
		t.Fatalf("expected ErrFloorPlansDirNotFound, got %v", err)
	}
}

func TestLoadFloorPlansDirOutsideDataDir(t *testing.T) {
	dataDir := t.TempDir()
	outsideDir := t.TempDir()
	writeAreasConfigIn(t, dataDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[areas]
config_file = "areas.yaml"
floor_plans = "`+outsideDir+`"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatalf("expected error for floor plans dir outside data_dir")
	}
	if !errors.Is(err, ErrFloorPlansDirOutsideDataDir) {
		t.Fatalf("expected ErrFloorPlansDirOutsideDataDir, got %v", err)
	}
}

func TestLoadFloorPlansDirEmpty(t *testing.T) {
	dataDir := t.TempDir()
	writeAreasConfigIn(t, dataDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[areas]
config_file = "areas.yaml"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error when floor_plans is omitted, got %v", err)
	}
	if cfg.Areas.FloorPlansDir != "" {
		t.Fatalf("expected empty floor_plans, got %s", cfg.Areas.FloorPlansDir)
	}
}

func TestLoadFloorPlansDirLegacyKey(t *testing.T) {
	dataDir := t.TempDir()
	writeAreasConfigIn(t, dataDir)
	fpDir := filepath.Join(dataDir, "floor_plans")
	if err := os.MkdirAll(fpDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[areas]
config_file = "areas.yaml"
floor_plans_dir = "floor_plans"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected legacy key to resolve, got %v", err)
	}
	if cfg.Areas.FloorPlansDir != fpDir {
		t.Fatalf("expected %s, got %s", fpDir, cfg.Areas.FloorPlansDir)
	}
}

func TestLoadAbsolutePathOutsideDataDir(t *testing.T) {
	dataDir := t.TempDir()
	outsideDir := t.TempDir()
	outsidePath := writeAreasConfigIn(t, outsideDir)
	path := writeConfig(t, `
[main]
data_dir = "`+dataDir+`"

[areas]
config_file = "`+outsidePath+`"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatalf("expected error for areas config outside data_dir")
	}
	if !errors.Is(err, ErrAreasConfigOutsideDataDir) {
		t.Fatalf("expected ErrAreasConfigOutsideDataDir, got %v", err)
	}
}
