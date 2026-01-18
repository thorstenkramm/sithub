// Package config loads SitHub configuration files.
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ErrMissingEntraIDConfig indicates incomplete Entra ID settings.
var ErrMissingEntraIDConfig = errors.New("missing Entra ID configuration")

// Config holds the full application configuration.
type Config struct {
	Main     MainConfig     `mapstructure:"main"`
	Log      LogConfig      `mapstructure:"log"`
	EntraID  EntraIDConfig  `mapstructure:"entraid"`
	TestAuth TestAuthConfig `mapstructure:"test_auth"`
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

// TestAuthConfig configures local test authentication.
type TestAuthConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	UserID   string `mapstructure:"user_id"`
	UserName string `mapstructure:"user_name"`
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
	v.SetDefault("test_auth.enabled", false)
	v.SetDefault("test_auth.user_id", "test-user")
	v.SetDefault("test_auth.user_name", "Test User")

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

	if !cfg.TestAuth.Enabled {
		if cfg.EntraID.AuthorizeURL == "" || cfg.EntraID.TokenURL == "" || cfg.EntraID.RedirectURI == "" ||
			cfg.EntraID.ClientID == "" || cfg.EntraID.ClientSecret == "" {
			return nil, fmt.Errorf("validate entraid: %w", ErrMissingEntraIDConfig)
		}
	}

	return &cfg, nil
}
