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

func TestFindArea(t *testing.T) {
	cfg := &Config{
		Areas: []Area{
			{ID: "a1", Name: "Main"},
			{ID: "a2", Name: "Annex"},
		},
	}

	area, ok := cfg.FindArea("a2")
	if !ok || area.Name != "Annex" {
		t.Fatalf("expected to find area a2")
	}

	if _, ok := cfg.FindArea("missing"); ok {
		t.Fatalf("expected missing area to be false")
	}
}

func TestFindRoom(t *testing.T) {
	cfg := &Config{
		Areas: []Area{
			{
				ID:   "a1",
				Name: "Main",
				Rooms: []Room{
					{ID: "r1", Name: "Room 1"},
				},
			},
		},
	}

	room, ok := cfg.FindRoom("r1")
	if !ok || room.Name != "Room 1" {
		t.Fatalf("expected to find room r1")
	}

	if _, ok := cfg.FindRoom("missing"); ok {
		t.Fatalf("expected missing room to be false")
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
