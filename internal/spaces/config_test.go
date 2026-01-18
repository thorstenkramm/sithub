package spaces

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spaces.yaml")
	content := `areas:
  - id: area-1
    name: Office
    rooms:
      - id: room-1
        name: Room 1
        desks:
          - id: desk-1
            name: Desk 1
            equipment:
              - Monitor
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write spaces config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load spaces config: %v", err)
	}
	if len(cfg.Areas) != 1 {
		t.Fatalf("expected 1 area, got %d", len(cfg.Areas))
	}
}

func TestLoadConfigMissingAreaID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spaces.yaml")
	content := `areas:
  - name: Office
    rooms: []
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write spaces config: %v", err)
	}

	if _, err := Load(path); err == nil {
		t.Fatalf("expected validation error")
	}
}
