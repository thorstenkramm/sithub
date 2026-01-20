package admin

import (
	"context"
	"fmt"
	"sync"

	"github.com/thorstenkramm/sithub/internal/spaces"
)

// ConfigHolder holds the current spaces configuration and provides thread-safe access.
type ConfigHolder struct {
	mu     sync.RWMutex
	config *spaces.Config
}

// NewConfigHolder creates a new ConfigHolder with an initial config.
func NewConfigHolder(initial *spaces.Config) *ConfigHolder {
	return &ConfigHolder{config: initial}
}

// Get returns the current config.
func (h *ConfigHolder) Get() *spaces.Config {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config
}

// Reload loads the config from the database and updates the holder.
func (h *ConfigHolder) Reload(ctx context.Context, store *spaces.Store) error {
	cfg, err := store.LoadConfig(ctx)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	h.mu.Lock()
	h.config = cfg
	h.mu.Unlock()
	return nil
}
