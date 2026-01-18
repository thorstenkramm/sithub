// Package startup wires the HTTP server and dependencies.
package startup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

	e.GET("/oauth/login", auth.LoginHandler(authService))
	e.GET("/oauth/callback", auth.CallbackHandler(authService)) //nolint:contextcheck // Echo handler uses request context.

	e.GET("/api/v1/ping", system.Ping)
	e.GET("/api/v1/me", auth.MeHandler(), middleware.RequireAuth)

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
		_ = server.Shutdown(shutdownCtx)
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}
