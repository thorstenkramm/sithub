package startup

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestRunShutsDownOnContextCancel(t *testing.T) {
	cfg := &config.Config{
		Main: config.MainConfig{
			Listen:  "127.0.0.1",
			Port:    0,
			DataDir: t.TempDir(),
		},
		EntraID: config.EntraIDConfig{
			AuthorizeURL: "https://login",
			TokenURL:     "https://token",
			RedirectURI:  "http://localhost/callback",
			ClientID:     "client",
			ClientSecret: "secret",
		},
		Spaces: config.SpacesConfig{
			ConfigFile: writeSpacesConfig(t),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)

	go func() {
		errCh <- Run(ctx, cfg)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for server shutdown")
	}
}

func writeSpacesConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "spaces.yaml")
	if err := os.WriteFile(path, []byte("areas: []\n"), 0o600); err != nil {
		t.Fatalf("write spaces config: %v", err)
	}
	return path
}
