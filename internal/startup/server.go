// Package startup wires the HTTP server and dependencies.
package startup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/middleware"
	"github.com/thorstenkramm/sithub/internal/system"
)

// Run starts the HTTP server and blocks until it shuts down.
func Run(ctx context.Context, cfg *config.Config) error {
	e := echo.New()
	e.HideBanner = true

	authService, err := auth.NewService(cfg)
	if err != nil {
		return fmt.Errorf("init auth service: %w", err)
	}

	e.Use(middleware.LoadUser(authService))
	e.Use(middleware.RedirectForbidden(authService))

	staticDir := "assets/web"
	indexPath := filepath.Join(staticDir, "index.html")

	//nolint:contextcheck // Echo handlers use request context.
	registerRoutes(e, authService)
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

func registerRoutes(e *echo.Echo, authService *auth.Service) {
	e.GET("/oauth/login", auth.LoginHandler(authService))
	e.GET("/oauth/callback", auth.CallbackHandler(authService))

	e.GET("/api/v1/ping", system.Ping)
	e.GET("/api/v1/me", auth.MeHandler(), middleware.RequireAuth(authService))
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
