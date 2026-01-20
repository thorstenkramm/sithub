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

	"github.com/thorstenkramm/sithub/internal/admin"
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

	// Load YAML config and sync to database
	yamlConfig, err := spaces.Load(cfg.Spaces.ConfigFile)
	if err != nil {
		return fmt.Errorf("load spaces config: %w", err)
	}

	spacesStore := spaces.NewStore(store)
	if err := spacesStore.SyncFromConfig(ctx, yamlConfig); err != nil {
		return fmt.Errorf("sync spaces config: %w", err)
	}

	// Load config from database (now the source of truth)
	spacesConfig, err := spacesStore.LoadConfig(ctx)
	if err != nil {
		return fmt.Errorf("load spaces from db: %w", err)
	}

	// Create config holder for live updates
	configHolder := admin.NewConfigHolder(spacesConfig)

	authService, err := auth.NewService(cfg)
	if err != nil {
		return fmt.Errorf("init auth service: %w", err)
	}

	notifier := notifications.NewNotifier(cfg.Notifications.WebhookURL)

	e.Use(middleware.LoadUser(authService))
	e.Use(middleware.RedirectForbidden(authService))

	staticDir := "assets/web"
	indexPath := filepath.Join(staticDir, "index.html")

	//nolint:contextcheck // Echo handlers use request context.
	registerRoutes(e, authService, configHolder, spacesStore, store, notifier)
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
	e *echo.Echo, authService *auth.Service, configHolder *admin.ConfigHolder,
	spacesStore *spaces.Store, store *sql.DB, notifier notifications.Notifier,
) {
	// Helper to get current config
	getConfig := func() *spaces.Config { return configHolder.Get() }

	e.GET("/oauth/login", auth.LoginHandler(authService))
	e.GET("/oauth/callback", auth.CallbackHandler(authService))

	e.GET("/api/v1/ping", system.Ping)
	e.GET("/api/v1/me", auth.MeHandler(), middleware.RequireAuth(authService))
	e.GET("/api/v1/areas", areas.ListHandlerDynamic(getConfig), middleware.RequireAuth(authService))
	e.GET("/api/v1/areas/:area_id/rooms", rooms.ListHandlerDynamic(getConfig), middleware.RequireAuth(authService))
	e.GET("/api/v1/areas/:area_id/presence",
		areas.PresenceHandlerDynamic(getConfig, store), middleware.RequireAuth(authService))
	e.GET("/api/v1/rooms/:room_id/desks",
		desks.ListHandlerDynamic(getConfig, store), middleware.RequireAuth(authService))
	e.GET("/api/v1/rooms/:room_id/bookings",
		rooms.BookingsHandlerDynamic(getConfig, store), middleware.RequireAuth(authService))
	e.GET("/api/v1/bookings", bookings.ListHandlerDynamic(getConfig, store), middleware.RequireAuth(authService))
	e.GET("/api/v1/bookings/history",
		bookings.HistoryHandlerDynamic(getConfig, store), middleware.RequireAuth(authService))
	e.POST("/api/v1/bookings",
		bookings.CreateHandlerDynamic(getConfig, store, notifier), middleware.RequireAuth(authService))
	e.DELETE("/api/v1/bookings/:id", bookings.DeleteHandler(store, notifier), middleware.RequireAuth(authService))

	// Admin routes
	e.GET("/api/v1/admin/areas",
		admin.ListAreasHandler(spacesStore), middleware.RequireAuth(authService))
	e.POST("/api/v1/admin/areas",
		admin.CreateAreaHandler(spacesStore, configHolder), middleware.RequireAuth(authService))
	e.PATCH("/api/v1/admin/areas/:area_id",
		admin.UpdateAreaHandler(spacesStore, configHolder), middleware.RequireAuth(authService))
	e.DELETE("/api/v1/admin/areas/:area_id",
		admin.DeleteAreaHandler(spacesStore, configHolder), middleware.RequireAuth(authService))
	e.GET("/api/v1/admin/areas/:area_id/rooms",
		admin.ListRoomsHandler(spacesStore), middleware.RequireAuth(authService))
	e.POST("/api/v1/admin/areas/:area_id/rooms",
		admin.CreateRoomHandler(spacesStore, configHolder), middleware.RequireAuth(authService))
	e.PATCH("/api/v1/admin/rooms/:room_id",
		admin.UpdateRoomHandler(spacesStore, configHolder), middleware.RequireAuth(authService))
	e.DELETE("/api/v1/admin/rooms/:room_id",
		admin.DeleteRoomHandler(spacesStore, configHolder), middleware.RequireAuth(authService))
	e.GET("/api/v1/admin/rooms/:room_id/desks",
		admin.ListDesksHandler(spacesStore), middleware.RequireAuth(authService))
	e.POST("/api/v1/admin/rooms/:room_id/desks",
		admin.CreateDeskHandler(spacesStore, configHolder), middleware.RequireAuth(authService))
	e.PATCH("/api/v1/admin/desks/:desk_id",
		admin.UpdateDeskHandler(spacesStore, configHolder), middleware.RequireAuth(authService))
	e.DELETE("/api/v1/admin/desks/:desk_id",
		admin.DeleteDeskHandler(spacesStore, configHolder), middleware.RequireAuth(authService))
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
