// Package config loads SitHub configuration files.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ErrMissingEntraIDConfig indicates incomplete Entra ID settings.
var ErrMissingEntraIDConfig = errors.New("missing Entra ID configuration")

// ErrMissingAreasConfig indicates missing areas configuration settings.
var ErrMissingAreasConfig = errors.New("missing areas configuration")

// ErrAreasConfigOutsideDataDir indicates the areas config path escapes data_dir.
var ErrAreasConfigOutsideDataDir = errors.New("areas config_file must be inside data_dir")

// ErrFloorPlansDirOutsideDataDir indicates the floor plans dir escapes data_dir.
var ErrFloorPlansDirOutsideDataDir = errors.New("areas floor_plans must be inside data_dir")

// ErrFloorPlansDirNotFound indicates the floor plans directory does not exist.
var ErrFloorPlansDirNotFound = errors.New("floor plans directory not found")

// Config holds the full application configuration.
type Config struct {
	Main          MainConfig          `mapstructure:"main"`
	Log           LogConfig           `mapstructure:"log"`
	EntraID       EntraIDConfig       `mapstructure:"entraid"`
	Areas         AreasConfig         `mapstructure:"areas"`
	Notifications NotificationsConfig `mapstructure:"notifications"`
}

// MainConfig contains main server settings.
type MainConfig struct {
	Listen  string `mapstructure:"listen"`
	Port    int    `mapstructure:"port"`
	DataDir string `mapstructure:"data_dir"`
}

// LogConfig contains logging configuration.
type LogConfig struct {
	File   string `mapstructure:"file"`
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// EntraIDConfig contains Entra ID OAuth configuration.
type EntraIDConfig struct {
	AuthorizeURL  string `mapstructure:"authorize_url"`
	TokenURL      string `mapstructure:"token_url"`
	RedirectURI   string `mapstructure:"redirect_uri"`
	ClientID      string `mapstructure:"client_id"`
	ClientSecret  string `mapstructure:"client_secret"`
	UsersGroupID  string `mapstructure:"users_group_id"`
	AdminsGroupID string `mapstructure:"admins_group_id"`
}

// AreasConfig contains areas configuration settings.
type AreasConfig struct {
	ConfigFile    string `mapstructure:"config_file"`
	FloorPlansDir string `mapstructure:"floor_plans"`
}

// NotificationsConfig contains notification settings.
type NotificationsConfig struct {
	WebhookURL string `mapstructure:"webhook_url"`
}

// Load reads configuration from a TOML file and environment variables.
func Load(path string) (*Config, error) {
	return LoadWithOverrides(path, nil)
}

// LoadWithOverrides reads configuration and applies explicit overrides.
func LoadWithOverrides(path string, overrides map[string]interface{}) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("toml")
	v.SetEnvPrefix("SITHUB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("main.listen", "127.0.0.1")
	v.SetDefault("main.port", 9900)
	v.SetDefault("main.data_dir", ".")
	v.SetDefault("log.file", "")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
	v.SetDefault("areas.config_file", "")
	v.SetDefault("areas.floor_plans", "")
	v.SetDefault("areas.floor_plans_dir", "")
	v.SetDefault("notifications.webhook_url", "")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	for key, value := range overrides {
		v.Set(key, value)
	}

	normalizeLegacyFloorPlansConfig(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := validateEntraIDConfig(&cfg.EntraID); err != nil {
		return nil, err
	}

	if err := resolveAreasConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func normalizeLegacyFloorPlansConfig(v *viper.Viper) {
	if strings.TrimSpace(v.GetString("areas.floor_plans")) != "" {
		return
	}

	legacy := strings.TrimSpace(v.GetString("areas.floor_plans_dir"))
	if legacy != "" {
		v.Set("areas.floor_plans", legacy)
	}
}

// EntraIDConfigured returns true if Entra ID OAuth is configured.
func (c *Config) EntraIDConfigured() bool {
	e := c.EntraID
	return e.AuthorizeURL != "" && e.TokenURL != "" && e.RedirectURI != "" &&
		e.ClientID != "" && e.ClientSecret != ""
}

// resolveAreasConfig validates and resolves the areas config file path
// relative to data_dir. Absolute paths outside data_dir are rejected.
func resolveAreasConfig(cfg *Config) error {
	raw := strings.TrimSpace(cfg.Areas.ConfigFile)
	if raw == "" {
		return fmt.Errorf("validate areas: %w", ErrMissingAreasConfig)
	}

	dataDir, err := filepath.Abs(cfg.Main.DataDir)
	if err != nil {
		return fmt.Errorf("resolve data_dir: %w", err)
	}

	var resolved string
	if filepath.IsAbs(raw) {
		resolved = filepath.Clean(raw)
	} else {
		resolved = filepath.Join(dataDir, raw)
	}

	// Ensure the resolved path is inside data_dir.
	if !strings.HasPrefix(resolved, dataDir+string(filepath.Separator)) && resolved != dataDir {
		return fmt.Errorf("validate areas: %w: %s", ErrAreasConfigOutsideDataDir, resolved)
	}

	if _, err := os.Stat(resolved); err != nil {
		return fmt.Errorf("validate areas: %w", err)
	}

	cfg.Areas.ConfigFile = resolved

	if err := resolveFloorPlansDir(cfg, dataDir); err != nil {
		return err
	}

	return nil
}

// resolveFloorPlansDir validates and resolves the floor plans directory
// relative to data_dir. If not set, floor plan features are disabled.
func resolveFloorPlansDir(cfg *Config, dataDir string) error {
	raw := strings.TrimSpace(cfg.Areas.FloorPlansDir)
	if raw == "" {
		return nil
	}

	var resolved string
	if filepath.IsAbs(raw) {
		resolved = filepath.Clean(raw)
	} else {
		resolved = filepath.Join(dataDir, raw)
	}

	if !strings.HasPrefix(resolved, dataDir+string(filepath.Separator)) && resolved != dataDir {
		return fmt.Errorf("validate areas: %w: %s", ErrFloorPlansDirOutsideDataDir, resolved)
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return fmt.Errorf("validate areas: %w: %s", ErrFloorPlansDirNotFound, resolved)
	}
	if !info.IsDir() {
		return fmt.Errorf("validate areas: %w: %s is not a directory", ErrFloorPlansDirNotFound, resolved)
	}

	cfg.Areas.FloorPlansDir = resolved
	return nil
}

// validateEntraIDConfig checks that either all 5 required Entra ID fields
// are set, or none are set (local-only mode).
func validateEntraIDConfig(e *EntraIDConfig) error {
	fields := []string{e.AuthorizeURL, e.TokenURL, e.RedirectURI, e.ClientID, e.ClientSecret}
	setCount := 0
	for _, f := range fields {
		if f != "" {
			setCount++
		}
	}
	if setCount == 0 || setCount == len(fields) {
		return nil
	}
	return fmt.Errorf("validate entraid: %w (all 5 fields required if any is set)", ErrMissingEntraIDConfig)
}
