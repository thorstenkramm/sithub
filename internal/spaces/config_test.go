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
    items:
      - id: room-1
        name: Room 1
        items:
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

func TestFindItemGroup(t *testing.T) {
	cfg := &Config{
		Areas: []Area{
			{
				ID:   "a1",
				Name: "Main",
				ItemGroups: []ItemGroup{
					{ID: "r1", Name: "Room 1"},
				},
			},
		},
	}

	room, ok := cfg.FindItemGroup("r1")
	if !ok || room.Name != "Room 1" {
		t.Fatalf("expected to find room r1")
	}

	if _, ok := cfg.FindItemGroup("missing"); ok {
		t.Fatalf("expected missing room to be false")
	}
}

func TestLoadConfigMissingAreaID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spaces.yaml")
	content := `areas:
  - name: Office
    items: []
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write spaces config: %v", err)
	}

	if _, err := Load(path); err == nil {
		t.Fatalf("expected validation error")
	}
}

// Issue 4: Add FindItem unit test
func TestFindItem(t *testing.T) {
	cfg := &Config{
		Areas: []Area{
			{
				ID:   "a1",
				Name: "Main",
				ItemGroups: []ItemGroup{
					{
						ID:   "r1",
						Name: "Room 1",
						Items: []Item{
							{ID: "d1", Name: "Desk 1"},
							{ID: "d2", Name: "Desk 2"},
						},
					},
					{
						ID:   "r2",
						Name: "Room 2",
						Items: []Item{
							{ID: "d3", Name: "Desk 3"},
						},
					},
				},
			},
		},
	}

	// Test finding desk in first room
	desk, ok := cfg.FindItem("d1")
	if !ok || desk.Name != "Desk 1" {
		t.Fatalf("expected to find desk d1")
	}

	// Test finding desk in second room
	desk, ok = cfg.FindItem("d3")
	if !ok || desk.Name != "Desk 3" {
		t.Fatalf("expected to find desk d3")
	}

	// Test missing desk
	if _, ok := cfg.FindItem("missing"); ok {
		t.Fatalf("expected missing desk to be false")
	}
}

func TestFindItemLocation(t *testing.T) {
	cfg := &Config{
		Areas: []Area{
			{
				ID:   "area-1",
				Name: "Main Office",
				ItemGroups: []ItemGroup{
					{
						ID:   "room-1",
						Name: "Room 101",
						Items: []Item{
							{ID: "desk-1", Name: "Corner Desk"},
						},
					},
				},
			},
			{
				ID:   "area-2",
				Name: "Annex",
				ItemGroups: []ItemGroup{
					{
						ID:   "room-2",
						Name: "Room 201",
						Items: []Item{
							{ID: "desk-2", Name: "Window Desk"},
						},
					},
				},
			},
		},
	}

	// Test finding desk with full location
	loc, ok := cfg.FindItemLocation("desk-2")
	if !ok {
		t.Fatalf("expected to find desk-2 location")
	}
	if loc.Area.ID != "area-2" || loc.Area.Name != "Annex" {
		t.Fatalf("expected area-2/Annex, got %s/%s", loc.Area.ID, loc.Area.Name)
	}
	if loc.ItemGroup.ID != "room-2" || loc.ItemGroup.Name != "Room 201" {
		t.Fatalf("expected room-2/Room 201, got %s/%s", loc.ItemGroup.ID, loc.ItemGroup.Name)
	}
	if loc.Item.ID != "desk-2" || loc.Item.Name != "Window Desk" {
		t.Fatalf("expected desk-2/Window Desk, got %s/%s", loc.Item.ID, loc.Item.Name)
	}

	// Test missing desk
	if _, ok := cfg.FindItemLocation("missing"); ok {
		t.Fatalf("expected missing desk location to be false")
	}
}
