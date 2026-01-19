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

// Issue 4: Add FindDesk unit test
func TestFindDesk(t *testing.T) {
	cfg := &Config{
		Areas: []Area{
			{
				ID:   "a1",
				Name: "Main",
				Rooms: []Room{
					{
						ID:   "r1",
						Name: "Room 1",
						Desks: []Desk{
							{ID: "d1", Name: "Desk 1"},
							{ID: "d2", Name: "Desk 2"},
						},
					},
					{
						ID:   "r2",
						Name: "Room 2",
						Desks: []Desk{
							{ID: "d3", Name: "Desk 3"},
						},
					},
				},
			},
		},
	}

	// Test finding desk in first room
	desk, ok := cfg.FindDesk("d1")
	if !ok || desk.Name != "Desk 1" {
		t.Fatalf("expected to find desk d1")
	}

	// Test finding desk in second room
	desk, ok = cfg.FindDesk("d3")
	if !ok || desk.Name != "Desk 3" {
		t.Fatalf("expected to find desk d3")
	}

	// Test missing desk
	if _, ok := cfg.FindDesk("missing"); ok {
		t.Fatalf("expected missing desk to be false")
	}
}
