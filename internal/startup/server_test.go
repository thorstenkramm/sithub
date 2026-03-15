package startup

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/areas"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/notifications"
)

func TestRunShutsDownOnContextCancel(t *testing.T) {
	dataDir := t.TempDir()
	areasFile := writeAreasConfigIn(t, dataDir)

	cfg := &config.Config{
		Main: config.MainConfig{
			Listen:  "127.0.0.1",
			Port:    0,
			DataDir: dataDir,
		},
		EntraID: config.EntraIDConfig{
			AuthorizeURL: "https://login",
			TokenURL:     "https://token",
			RedirectURI:  "http://localhost/callback",
			ClientID:     "client",
			ClientSecret: "secret",
		},
		Areas: config.AreasConfig{
			ConfigFile: areasFile,
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

func writeAreasConfigIn(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "areas.yaml")
	if err := os.WriteFile(path, []byte("areas: []\n"), 0o600); err != nil {
		t.Fatalf("write areas config: %v", err)
	}
	return path
}

func TestFloorPlanRouteRequiresAuthentication(t *testing.T) {
	e := echo.New()
	authService := newTestAuthService(t)

	registerRoutes(
		e,
		authService,
		&areas.Config{},
		t.TempDir(),
		nil,
		notifications.NewNotifier(""),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/floor-plans/plan.png", http.NoBody)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func newTestAuthService(t *testing.T) *auth.Service {
	t.Helper()

	svc, err := auth.NewService(&config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}}, nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return svc
}
