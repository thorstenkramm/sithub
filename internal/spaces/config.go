// Package spaces loads area configuration from YAML files.
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

// ConfigGetter is a function that returns the current spaces config.
// This allows handlers to use dynamically reloaded configuration.
type ConfigGetter func() *Config

// FindArea returns the area matching the provided id.
func (c *Config) FindArea(id string) (*Area, bool) {
	for i := range c.Areas {
		if c.Areas[i].ID == id {
			return &c.Areas[i], true
		}
	}
	return nil, false
}

// FindItemGroup returns the item group matching the provided id.
func (c *Config) FindItemGroup(id string) (*ItemGroup, bool) {
	for i := range c.Areas {
		for j := range c.Areas[i].ItemGroups {
			if c.Areas[i].ItemGroups[j].ID == id {
				return &c.Areas[i].ItemGroups[j], true
			}
		}
	}
	return nil, false
}

// FindItem returns the item matching the provided id.
func (c *Config) FindItem(id string) (*Item, bool) {
	for i := range c.Areas {
		for j := range c.Areas[i].ItemGroups {
			for k := range c.Areas[i].ItemGroups[j].Items {
				if c.Areas[i].ItemGroups[j].Items[k].ID == id {
					return &c.Areas[i].ItemGroups[j].Items[k], true
				}
			}
		}
	}
	return nil, false
}

// ItemLocation contains an item with its parent item group and area.
type ItemLocation struct {
	Area      *Area
	ItemGroup *ItemGroup
	Item      *Item
}

// FindItemLocation returns the item and its parent item group and area.
func (c *Config) FindItemLocation(itemID string) (*ItemLocation, bool) {
	for i := range c.Areas {
		for j := range c.Areas[i].ItemGroups {
			for k := range c.Areas[i].ItemGroups[j].Items {
				if c.Areas[i].ItemGroups[j].Items[k].ID == itemID {
					return &ItemLocation{
						Area:      &c.Areas[i],
						ItemGroup: &c.Areas[i].ItemGroups[j],
						Item:      &c.Areas[i].ItemGroups[j].Items[k],
					}, true
				}
			}
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

// ItemAttributes returns attributes for item resources.
func ItemAttributes(name string, equipment []string, warning, availability string) map[string]interface{} {
	attrs := map[string]interface{}{
		"name":      name,
		"equipment": equipment,
	}
	if warning != "" {
		attrs["warning"] = warning
	}
	if availability != "" {
		attrs["availability"] = availability
	}
	return attrs
}

// Area describes a bookable area.
type Area struct {
	ID          string      `yaml:"id"`
	Name        string      `yaml:"name"`
	Description string      `yaml:"description,omitempty"`
	FloorPlan   string      `yaml:"floor_plan,omitempty"`
	ItemGroups  []ItemGroup `yaml:"items"`
}

// ItemGroup describes a group of bookable items within an area.
type ItemGroup struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	FloorPlan   string `yaml:"floor_plan,omitempty"`
	Items       []Item `yaml:"items"`
}

// Item describes a bookable item within an item group.
type Item struct {
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
		for _, ig := range area.ItemGroups {
			if ig.ID == "" || ig.Name == "" {
				return fmt.Errorf("item group requires id and name")
			}
			for _, item := range ig.Items {
				if item.ID == "" || item.Name == "" {
					return fmt.Errorf("item requires id and name")
				}
			}
		}
	}
	return nil
}
