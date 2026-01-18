// Package startup wires the HTTP server and dependencies.
package startup

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/areas"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/db"
	"github.com/thorstenkramm/sithub/internal/middleware"
	"github.com/thorstenkramm/sithub/internal/system"
)

// Run starts the HTTP server and blocks until it shuts down.
func Run(ctx context.Context, cfg *config.Config) error {
	e := echo.New()
	e.HideBanner = true

	migrationsPath, err := resolveMigrationsPath()
	if err != nil {
		return fmt.Errorf("resolve migrations path: %w", err)
	}

	store, err := db.Open(cfg.Main.DataDir)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			slog.Error("close database", "err", err)
		}
	}()

	if err := db.RunMigrations(store, migrationsPath); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	authService, err := auth.NewService(cfg)
	if err != nil {
		return fmt.Errorf("init auth service: %w", err)
	}

	e.Use(middleware.LoadUser(authService))
	e.Use(middleware.RedirectForbidden(authService))

	staticDir := "assets/web"
	indexPath := filepath.Join(staticDir, "index.html")

	areasRepo := areas.NewRepository(store)

	//nolint:contextcheck // Echo handlers use request context.
	registerRoutes(e, authService, areasRepo)
	registerSPAHandlers(e, staticDir, indexPath)

	addr := fmt.Sprintf("%s:%d", cfg.Main.Listen, cfg.Main.Port)
	server := &http.Server{
		Addr:              addr,
		Handler:           e,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = err
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

func registerRoutes(e *echo.Echo, authService *auth.Service, areasRepo *areas.Repository) {
	e.GET("/oauth/login", auth.LoginHandler(authService))
	e.GET("/oauth/callback", auth.CallbackHandler(authService))

	e.GET("/api/v1/ping", system.Ping)
	e.GET("/api/v1/me", auth.MeHandler(), middleware.RequireAuth(authService))
	e.GET("/api/v1/areas", areas.ListHandler(areasRepo), middleware.RequireAuth(authService))
}

func registerSPAHandlers(e *echo.Echo, staticDir, indexPath string) {
	e.Static("/", staticDir)

	defaultErrorHandler := e.HTTPErrorHandler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if req := c.Request(); req != nil && req.Method == http.MethodGet {
			var httpErr *echo.HTTPError
			if errors.As(err, &httpErr) && httpErr.Code == http.StatusNotFound {
				path := req.URL.Path
				if !strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/oauth/") {
					if fileErr := c.File(indexPath); fileErr == nil {
						return
					}
				}
			}
		}
		defaultErrorHandler(err, c)
	}
}

func resolveMigrationsPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve migrations path")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	return filepath.Join(root, "migrations"), nil
}
