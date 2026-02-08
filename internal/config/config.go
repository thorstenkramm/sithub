// Package config loads SitHub configuration files.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// ErrMissingEntraIDConfig indicates incomplete Entra ID settings.
var ErrMissingEntraIDConfig = errors.New("missing Entra ID configuration")

// ErrMissingSpacesConfig indicates missing spaces configuration settings.
var ErrMissingSpacesConfig = errors.New("missing spaces configuration")

// Config holds the full application configuration.
type Config struct {
	Main          MainConfig          `mapstructure:"main"`
	Log           LogConfig           `mapstructure:"log"`
	EntraID       EntraIDConfig       `mapstructure:"entraid"`
	Spaces        SpacesConfig        `mapstructure:"spaces"`
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

// SpacesConfig contains space configuration settings.
type SpacesConfig struct {
	ConfigFile string `mapstructure:"config_file"`
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
	v.SetDefault("spaces.config_file", "")
	v.SetDefault("notifications.webhook_url", "")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	for key, value := range overrides {
		v.Set(key, value)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := validateEntraIDConfig(&cfg.EntraID); err != nil {
		return nil, err
	}

	if strings.TrimSpace(cfg.Spaces.ConfigFile) == "" {
		return nil, fmt.Errorf("validate spaces: %w", ErrMissingSpacesConfig)
	}
	if _, err := os.Stat(cfg.Spaces.ConfigFile); err != nil {
		return nil, fmt.Errorf("validate spaces: %w", err)
	}

	return &cfg, nil
}

// EntraIDConfigured returns true if Entra ID OAuth is configured.
func (c *Config) EntraIDConfigured() bool {
	e := c.EntraID
	return e.AuthorizeURL != "" && e.TokenURL != "" && e.RedirectURI != "" &&
		e.ClientID != "" && e.ClientSecret != ""
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
