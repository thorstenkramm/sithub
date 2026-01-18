// Package spaces loads and stores area configuration.
package spaces

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the space configuration.
type Config struct {
	Areas []Area `yaml:"areas"`
}

// FindArea returns the area matching the provided id.
func (c *Config) FindArea(id string) (*Area, bool) {
	for i := range c.Areas {
		if c.Areas[i].ID == id {
			return &c.Areas[i], true
		}
	}
	return nil, false
}

// BaseAttributes returns common attributes for named space resources.
func BaseAttributes(name, description, floorPlan string) map[string]interface{} {
	attrs := map[string]interface{}{
		"name": name,
	}
	if description != "" {
		attrs["description"] = description
	}
	if floorPlan != "" {
		attrs["floor_plan"] = floorPlan
	}
	return attrs
}

// Area describes a bookable area.
type Area struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	FloorPlan   string `yaml:"floor_plan,omitempty"`
	Rooms       []Room `yaml:"rooms"`
}

// Room describes a room within an area.
type Room struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	FloorPlan   string `yaml:"floor_plan,omitempty"`
	Desks       []Desk `yaml:"desks"`
}

// Desk describes a desk within a room.
type Desk struct {
	ID        string   `yaml:"id"`
	Name      string   `yaml:"name"`
	Equipment []string `yaml:"equipment"`
	Warning   string   `yaml:"warning,omitempty"`
}

// Load reads and parses a space configuration file.
func Load(path string) (*Config, error) {
	// #nosec G304 -- path comes from explicit configuration.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read spaces config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse spaces config: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("validate spaces config: %w", err)
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	for _, area := range cfg.Areas {
		if area.ID == "" || area.Name == "" {
			return fmt.Errorf("area requires id and name")
		}
		for _, room := range area.Rooms {
			if room.ID == "" || room.Name == "" {
				return fmt.Errorf("room requires id and name")
			}
			for _, desk := range room.Desks {
				if desk.ID == "" || desk.Name == "" {
					return fmt.Errorf("desk requires id and name")
				}
			}
		}
	}
	return nil
}
