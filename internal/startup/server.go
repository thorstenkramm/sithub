// Package startup wires the HTTP server and dependencies.
package startup

import (
	"context"
	"database/sql"
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
	"github.com/thorstenkramm/sithub/internal/bookings"
	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/db"
	"github.com/thorstenkramm/sithub/internal/desks"
	"github.com/thorstenkramm/sithub/internal/middleware"
	"github.com/thorstenkramm/sithub/internal/notifications"
	"github.com/thorstenkramm/sithub/internal/rooms"
	"github.com/thorstenkramm/sithub/internal/spaces"
	"github.com/thorstenkramm/sithub/internal/system"
	"github.com/thorstenkramm/sithub/internal/users"
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

	// Load spaces configuration from YAML (single source of truth)
	spacesConfig, err := spaces.Load(cfg.Spaces.ConfigFile)
	if err != nil {
		return fmt.Errorf("load spaces config: %w", err)
	}

	authService, err := auth.NewService(cfg, store)
	if err != nil {
		return fmt.Errorf("init auth service: %w", err)
	}

	notifier := notifications.NewNotifier(cfg.Notifications.WebhookURL)

	e.Use(middleware.LoadUser(authService))
	e.Use(middleware.RedirectForbidden(authService))

	staticDir := "assets/web"
	indexPath := filepath.Join(staticDir, "index.html")

	//nolint:contextcheck // Echo handlers use request context.
	registerRoutes(e, authService, spacesConfig, store, notifier)
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

func registerRoutes(
	e *echo.Echo, authService *auth.Service, spacesConfig *spaces.Config,
	store *sql.DB, notifier notifications.Notifier,
) {
	// Helper to get current config (returns the same config, loaded at startup)
	getConfig := func() *spaces.Config { return spacesConfig }

	// OAuth routes
	e.GET("/oauth/login", auth.LoginHandler(authService))
	e.GET("/oauth/callback", auth.CallbackHandler(authService))

	// Auth routes (no auth middleware required)
	loginLimiter := middleware.NewRateLimiter(60, time.Minute)
	e.POST("/api/v1/auth/login", auth.LocalLoginHandler(authService),
		middleware.RateLimit(loginLimiter))
	e.POST("/api/v1/auth/logout", auth.LogoutHandler())

	// Public
	e.GET("/api/v1/ping", system.Ping)

	// Authenticated routes
	requireAuth := middleware.RequireAuth(authService)
	e.GET("/api/v1/me", auth.MeHandler(), requireAuth)
	e.PATCH("/api/v1/me", auth.UpdateMeHandler(authService), requireAuth)
	e.GET("/api/v1/areas", areas.ListHandlerDynamic(getConfig), requireAuth)
	e.GET("/api/v1/areas/:area_id/rooms", rooms.ListHandlerDynamic(getConfig), requireAuth)
	e.GET("/api/v1/areas/:area_id/presence",
		areas.PresenceHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/rooms/:room_id/desks",
		desks.ListHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/rooms/:room_id/bookings",
		rooms.BookingsHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/bookings", bookings.ListHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/bookings/history",
		bookings.HistoryHandlerDynamic(getConfig, store), requireAuth)
	e.POST("/api/v1/bookings",
		bookings.CreateHandlerDynamic(getConfig, store, notifier), requireAuth)
	e.DELETE("/api/v1/bookings/:id", bookings.DeleteHandler(store, notifier), requireAuth)

	// User management routes
	requireAdmin := middleware.RequireAdmin()
	e.GET("/api/v1/users", users.ListHandler(store), requireAuth, requireAdmin)
	e.GET("/api/v1/users/:id", users.GetHandler(store), requireAuth, requireAdmin)
	e.POST("/api/v1/users", users.CreateHandler(store), requireAuth, requireAdmin)
	e.PATCH("/api/v1/users/:id", users.UpdateHandler(store), requireAuth, requireAdmin)
	e.DELETE("/api/v1/users/:id", users.DeleteHandler(store), requireAuth, requireAdmin)
}

func registerSPAHandlers(e *echo.Echo, staticDir, indexPath string) {
	e.Static("/", staticDir)

	defaultErrorHandler := e.HTTPErrorHandler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if req := c.Request(); req != nil && req.Method == http.MethodGet {
			var httpErr *echo.HTTPError
			if errors.As(err, &httpErr) && httpErr.Code == http.StatusNotFound {
				path := req.URL.Path
				apiPath := strings.HasPrefix(path, "/api/")
				oauthPath := strings.HasPrefix(path, "/oauth/")
				authPath := strings.HasPrefix(path, "/auth/")
				if !apiPath && !oauthPath && !authPath {
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
