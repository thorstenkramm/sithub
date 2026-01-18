package startup

import (
	"context"
	"testing"
	"time"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestRunShutsDownOnContextCancel(t *testing.T) {
	cfg := &config.Config{
		Main: config.MainConfig{
			Listen: "127.0.0.1",
			Port:   0,
		},
		EntraID: config.EntraIDConfig{
			AuthorizeURL: "https://login",
			TokenURL:     "https://token",
			RedirectURI:  "http://localhost/callback",
			ClientID:     "client",
			ClientSecret: "secret",
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
