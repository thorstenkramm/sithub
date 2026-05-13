// Package startup wires the HTTP server and dependencies.
package startup

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/assets"
	"github.com/thorstenkramm/sithub/internal/areas"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/bookings"
	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/db"
	"github.com/thorstenkramm/sithub/internal/floorplanpos"
	"github.com/thorstenkramm/sithub/internal/itemgroups"
	"github.com/thorstenkramm/sithub/internal/items"
	"github.com/thorstenkramm/sithub/internal/livefeed"
	"github.com/thorstenkramm/sithub/internal/middleware"
	"github.com/thorstenkramm/sithub/internal/notifications"
	"github.com/thorstenkramm/sithub/internal/system"
	"github.com/thorstenkramm/sithub/internal/users"
)

// Run starts the HTTP server and blocks until it shuts down.
func Run(ctx context.Context, cfg *config.Config) error {
	e := echo.New()
	e.HideBanner = true

	store, err := db.Open(cfg.Main.DataDir)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			slog.Error("close database", "err", err)
		}
	}()

	if err := db.RunMigrations(store); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	areasConfig, err := loadAndValidateAreas(cfg)
	if err != nil {
		return err
	}

	avatarsDir, err := ensureAvatarsDir(cfg.Main.DataDir)
	if err != nil {
		return err
	}

	authService, err := auth.NewService(cfg, store)
	if err != nil {
		return fmt.Errorf("init auth service: %w", err)
	}

	webhookNotifier := notifications.NewNotifier(cfg.Notifications.WebhookURL)
	hub := livefeed.NewHub()
	go hub.Run(ctx)
	notifier := notifications.MultiNotifier{webhookNotifier, hub}

	e.Use(middleware.LoadUser(authService))
	e.Use(middleware.RedirectForbidden(authService))

	webFS, err := fs.Sub(assets.Web, "web")
	if err != nil {
		return fmt.Errorf("open embedded frontend: %w", err)
	}

	bookingLimits := &bookings.BookingLimits{
		WeeksInAdvanced:      cfg.Bookings.WeeksInAdvanced,
		MaxBookingsPerPerson: cfg.Bookings.MaxBookingsPerPerson,
	}

	//nolint:contextcheck // Echo handlers use request context.
	registerRoutes(e, authService, areasConfig, cfg.Areas.FloorPlansDir, avatarsDir, store,
		notifier, hub, bookingLimits)
	registerSPAHandlers(e, webFS)

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
	e *echo.Echo, authService *auth.Service, areasConfig *areas.Config,
	floorPlansDir, avatarsDir string, store *sql.DB, notifier notifications.Notifier,
	liveHub *livefeed.Hub, bookingLimits *bookings.BookingLimits,
) {
	// Helper to get current config (returns the same config, loaded at startup)
	getConfig := func() *areas.Config { return areasConfig }

	// OAuth routes
	e.GET("/oauth/login", auth.LoginHandler(authService))
	e.GET("/oauth/callback", auth.CallbackHandler(authService, avatarsDir))

	// Auth routes (no auth middleware required)
	loginLimiter := middleware.NewRateLimiter(60, time.Minute)
	e.POST("/api/v1/auth/login", auth.LocalLoginHandler(authService),
		middleware.RateLimit(loginLimiter))
	e.POST("/api/v1/auth/logout", auth.LogoutHandler())
	e.GET("/api/v1/auth/providers", auth.ProvidersHandler(authService))

	// Public
	e.GET("/api/v1/ping", system.Ping)

	// Authenticated routes
	requireAuth := middleware.RequireAuth(authService)
	weeksInAdvanced := 5
	if bookingLimits != nil {
		weeksInAdvanced = bookingLimits.WeeksInAdvanced
	}
	e.GET("/api/v1/settings", system.SettingsHandler(weeksInAdvanced), requireAuth)
	e.GET("/api/v1/me", auth.MeHandler(), requireAuth)
	e.PATCH("/api/v1/me", auth.UpdateMeHandler(authService), requireAuth)
	e.GET("/api/v1/areas", areas.ListHandlerDynamic(getConfig), requireAuth)
	e.GET("/api/v1/areas/:area_id/item-groups",
		itemgroups.ListHandlerDynamic(getConfig), requireAuth)
	e.GET("/api/v1/areas/:area_id/item-groups/availability",
		itemgroups.AvailabilityHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/areas/:area_id/item-groups/matrix",
		itemgroups.MatrixHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/areas/:area_id/presence",
		areas.PresenceHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/item-groups/:item_group_id/items",
		items.ListHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/item-groups/:item_group_id/bookings",
		itemgroups.BookingsHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/bookings", bookings.ListHandlerDynamic(getConfig, store), requireAuth)
	e.GET("/api/v1/bookings/history",
		bookings.HistoryHandlerDynamic(getConfig, store), requireAuth)
	e.POST("/api/v1/bookings",
		bookings.CreateHandlerDynamic(getConfig, store, notifier, bookingLimits), requireAuth)
	e.PATCH("/api/v1/bookings/:id", bookings.PatchHandler(store), requireAuth)
	e.DELETE("/api/v1/bookings/:id", bookings.DeleteHandler(store, notifier), requireAuth)

	// Live feed (WebSocket) for real-time booking updates.
	e.GET("/api/v1/live", livefeed.Handler(liveHub), requireAuth)

	// Floor plan images (authenticated)
	e.GET("/api/v1/floor-plans/:filename",
		areas.FloorPlanHandler(floorPlansDir), requireAuth)

	// Avatar routes (authenticated)
	e.GET("/api/v1/avatars/:user_id",
		auth.ServeAvatarHandler(avatarsDir), requireAuth)
	e.POST("/api/v1/me/avatar",
		auth.UploadAvatarHandler(avatarsDir), requireAuth)
	e.DELETE("/api/v1/me/avatar",
		auth.DeleteAvatarHandler(avatarsDir), requireAuth)

	// Colleagues endpoint (all authenticated users)
	e.GET("/api/v1/colleagues", users.ColleaguesHandler(store), requireAuth)

	// User management routes
	requireAdmin := middleware.RequireAdmin()
	e.GET("/api/v1/users", users.ListHandler(store), requireAuth, requireAdmin)
	e.GET("/api/v1/users/:id", users.GetHandler(store), requireAuth, requireAdmin)
	e.POST("/api/v1/users", users.CreateHandler(store), requireAuth, requireAdmin)
	e.PATCH("/api/v1/users/:id", users.UpdateHandler(store), requireAuth, requireAdmin)
	e.DELETE("/api/v1/users/:id", users.DeleteHandler(store), requireAuth, requireAdmin)

	// Floor plan positions (read: any authenticated user, write: admin only)
	e.GET("/api/v1/floor-plan-positions",
		floorplanpos.ListHandler(store), requireAuth)
	e.POST("/api/v1/floor-plan-positions",
		floorplanpos.CreateHandler(store), requireAuth, requireAdmin)
	e.PUT("/api/v1/floor-plan-positions/:id",
		floorplanpos.UpdateHandler(store), requireAuth, requireAdmin)
	e.DELETE("/api/v1/floor-plan-positions/:id",
		floorplanpos.DeleteHandler(store), requireAuth, requireAdmin)
}

func loadAndValidateAreas(cfg *config.Config) (*areas.Config, error) {
	areasConfig, err := areas.Load(cfg.Areas.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("load areas config: %w", err)
	}
	for _, warning := range areas.FindInvalidConfiguredIcons(areasConfig) {
		slog.Warn(
			"invalid configured icon; frontend will fall back to the default icon",
			"location", warning.Location,
			"icon", warning.Icon,
		)
	}
	if cfg.Areas.FloorPlansDir != "" {
		if err := areas.ValidateFloorPlans(areasConfig, cfg.Areas.FloorPlansDir); err != nil {
			return nil, fmt.Errorf("validate floor plans: %w", err)
		}
	}
	if err := areas.ValidateReservations(areasConfig); err != nil {
		return nil, fmt.Errorf("validate reservations: %w", err)
	}
	return areasConfig, nil
}

func ensureAvatarsDir(dataDir string) (string, error) {
	dir := filepath.Join(dataDir, "avatars")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", fmt.Errorf("create avatars directory: %w", err)
	}
	return dir, nil
}

func registerSPAHandlers(e *echo.Echo, webFS fs.FS) {
	e.StaticFS("/", webFS)

	indexHTML, err := fs.ReadFile(webFS, "index.html")
	if err != nil {
		slog.Warn("embedded frontend missing index.html; SPA fallback disabled")
	}

	defaultErrorHandler := e.HTTPErrorHandler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if req := c.Request(); req != nil && req.Method == http.MethodGet && indexHTML != nil {
			var httpErr *echo.HTTPError
			if errors.As(err, &httpErr) && httpErr.Code == http.StatusNotFound {
				path := req.URL.Path
				apiPath := strings.HasPrefix(path, "/api/")
				oauthPath := strings.HasPrefix(path, "/oauth/")
				authPath := strings.HasPrefix(path, "/auth/")
				if !apiPath && !oauthPath && !authPath {
					c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
					if writeErr := c.HTMLBlob(http.StatusOK, indexHTML); writeErr == nil {
						return
					}
				}
			}
		}
		defaultErrorHandler(err, c)
	}
}
