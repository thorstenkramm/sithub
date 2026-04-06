// Package areas provides area configuration, handlers, and domain types.
package areas

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var mdiIconNamePattern = regexp.MustCompile(`^mdi-[a-z0-9-]+$`)

// Config holds the areas configuration.
type Config struct {
	Areas []Area `yaml:"areas"`
}

// ConfigGetter is a function that returns the current areas config.
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

// BaseAttributes returns common attributes for named area resources.
func BaseAttributes(name, description, floorPlan, icon string) map[string]interface{} {
	attrs := map[string]interface{}{
		"name": name,
	}
	if description != "" {
		attrs["description"] = description
	}
	if floorPlan != "" {
		attrs["floor_plan"] = floorPlan
	}
	if icon != "" {
		attrs["icon"] = icon
	}
	return attrs
}

// ItemAttributes returns attributes for item resources.
func ItemAttributes(name string, equipment []string, warning, availability, icon string) map[string]interface{} {
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
	if icon != "" {
		attrs["icon"] = icon
	}
	return attrs
}

// Area describes a bookable area.
type Area struct {
	ID                   string      `yaml:"id"`
	Name                 string      `yaml:"name"`
	Description          string      `yaml:"description,omitempty"`
	FloorPlan            string      `yaml:"floor_plan,omitempty"`
	Icon                 string      `yaml:"icon,omitempty"`
	MaxBookingsPerPerson int         `yaml:"max_bookings_per_person,omitempty"`
	ReservedFor          []string    `yaml:"reserved_for,omitempty"`
	ItemGroups           []ItemGroup `yaml:"items"`
}

// ItemGroup describes a group of bookable items within an area.
type ItemGroup struct {
	ID                   string   `yaml:"id"`
	Name                 string   `yaml:"name"`
	Description          string   `yaml:"description,omitempty"`
	FloorPlan            string   `yaml:"floor_plan,omitempty"`
	Icon                 string   `yaml:"icon,omitempty"`
	MaxBookingsPerPerson int      `yaml:"max_bookings_per_person,omitempty"`
	ReservedFor          []string `yaml:"reserved_for,omitempty"`
	Items                []Item   `yaml:"items"`
}

// Item describes a bookable item within an item group.
type Item struct {
	ID                   string   `yaml:"id"`
	Name                 string   `yaml:"name"`
	Equipment            []string `yaml:"equipment"`
	Warning              string   `yaml:"warning,omitempty"`
	Icon                 string   `yaml:"icon,omitempty"`
	MaxBookingsPerPerson int      `yaml:"max_bookings_per_person,omitempty"`
	ReservedFor          []string `yaml:"reserved_for,omitempty"`
}

// IconWarning describes an invalid configured icon reference.
type IconWarning struct {
	Location string
	Icon     string
}

// Load reads and parses an areas configuration file.
func Load(path string) (*Config, error) {
	// #nosec G304 -- path comes from explicit configuration.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read areas config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse areas config: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("validate areas config: %w", err)
	}

	return &cfg, nil
}

// supportedFloorPlanExts lists allowed floor plan image extensions.
var supportedFloorPlanExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".svg":  true,
}

// ValidateFloorPlans checks that all floor_plan references in the config
// point to existing files with supported formats inside floorPlansDir.
func ValidateFloorPlans(cfg *Config, floorPlansDir string) error {
	for i := range cfg.Areas {
		area := &cfg.Areas[i]
		if area.FloorPlan != "" {
			if err := validateFloorPlanFile(area.FloorPlan, floorPlansDir); err != nil {
				return fmt.Errorf("area %q: %w", area.ID, err)
			}
		}
		for j := range area.ItemGroups {
			ig := &area.ItemGroups[j]
			if ig.FloorPlan != "" {
				if err := validateFloorPlanFile(ig.FloorPlan, floorPlansDir); err != nil {
					return fmt.Errorf("item group %q: %w", ig.ID, err)
				}
			}
		}
	}
	return nil
}

// FindInvalidConfiguredIcons returns non-fatal warnings for icon values that cannot
// be rendered safely by the frontend's configured icon resolver.
func FindInvalidConfiguredIcons(cfg *Config) []IconWarning {
	warnings := make([]IconWarning, 0)

	for i := range cfg.Areas {
		area := &cfg.Areas[i]
		if area.Icon != "" && !isValidConfiguredIcon(area.Icon) {
			warnings = append(warnings, IconWarning{
				Location: fmt.Sprintf("area %q", area.ID),
				Icon:     area.Icon,
			})
		}

		for j := range area.ItemGroups {
			ig := &area.ItemGroups[j]
			if ig.Icon != "" && !isValidConfiguredIcon(ig.Icon) {
				warnings = append(warnings, IconWarning{
					Location: fmt.Sprintf("item group %q", ig.ID),
					Icon:     ig.Icon,
				})
			}

			for _, item := range ig.Items {
				if item.Icon != "" && !isValidConfiguredIcon(item.Icon) {
					warnings = append(warnings, IconWarning{
						Location: fmt.Sprintf("item %q", item.ID),
						Icon:     item.Icon,
					})
				}
			}
		}
	}

	return warnings
}

// validateFloorPlanFile checks that a floor plan filename exists in the
// directory and has a supported extension.
func validateFloorPlanFile(filename, floorPlansDir string) error {
	if strings.ContainsAny(filename, `/\`) || filepath.Base(filename) != filename {
		return fmt.Errorf("floor plan must be a filename only: %q", filename)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if !supportedFloorPlanExts[ext] {
		return fmt.Errorf("unsupported floor plan format %q (allowed: jpg, png, svg)", filename)
	}

	fullPath := filepath.Join(floorPlansDir, filename)

	if _, err := os.Stat(fullPath); err != nil {
		return fmt.Errorf("floor plan not found: %s", fullPath)
	}
	return nil
}

// ErrReservationConflict indicates a hierarchical reservation conflict.
var ErrReservationConflict = errors.New("reservation conflict")

// ValidateReservations checks that child reserved_for lists are subsets of parent lists.
func ValidateReservations(cfg *Config) error {
	for i := range cfg.Areas {
		area := &cfg.Areas[i]
		areaSet := toStringSet(area.ReservedFor)
		for j := range area.ItemGroups {
			ig := &area.ItemGroups[j]
			if err := checkSubset(areaSet, ig.ReservedFor, area.ID, ig.ID, "item group"); err != nil {
				return err
			}
			igSet := toStringSet(ig.ReservedFor)
			if len(ig.ReservedFor) == 0 {
				igSet = areaSet
			}
			for k := range ig.Items {
				item := &ig.Items[k]
				if err := checkSubset(igSet, item.ReservedFor, ig.ID, item.ID, "item"); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func toStringSet(s []string) map[string]struct{} {
	set := make(map[string]struct{}, len(s))
	for _, v := range s {
		set[v] = struct{}{}
	}
	return set
}

// checkSubset verifies that all entries in child are present in parent (if parent is non-empty).
func checkSubset(parentSet map[string]struct{}, child []string, parentID, childID, childType string) error {
	if len(parentSet) == 0 || len(child) == 0 {
		return nil
	}
	for _, email := range child {
		if _, ok := parentSet[email]; !ok {
			return fmt.Errorf(
				"%w: %s %q reserves for %q but parent %q does not include this user",
				ErrReservationConflict, childType, childID, email, parentID,
			)
		}
	}
	return nil
}

// IsReserved checks if the item at the given location is reserved and the user is excluded.
// Returns true if the user cannot book (is excluded from a reserved_for list).
func IsReserved(loc *ItemLocation, userEmail string) bool {
	// Check item level
	if len(loc.Item.ReservedFor) > 0 {
		return !containsString(loc.Item.ReservedFor, userEmail)
	}
	// Check item group level
	if len(loc.ItemGroup.ReservedFor) > 0 {
		return !containsString(loc.ItemGroup.ReservedFor, userEmail)
	}
	// Check area level
	if len(loc.Area.ReservedFor) > 0 {
		return !containsString(loc.Area.ReservedFor, userEmail)
	}
	return false
}

func containsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func validateConfig(cfg *Config) error {
	for i := range cfg.Areas {
		area := &cfg.Areas[i]
		if area.ID == "" || area.Name == "" {
			return fmt.Errorf("area requires id and name")
		}
		for j := range area.ItemGroups {
			ig := &area.ItemGroups[j]
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

func isValidConfiguredIcon(icon string) bool {
	return mdiIconNamePattern.MatchString(strings.TrimSpace(icon))
}
