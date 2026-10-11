package middleware

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/apikeys"
	"github.com/thorstenkramm/sithub/internal/auth"
)

// LoadAPIKey authenticates explicit bearer credentials before cookies on API paths.
// Invalid bearer credentials fail closed: they never fall back to a session.
func LoadAPIKey(store *sql.DB) echo.MiddlewareFunc {
	return loadAPIKey(store, time.Now)
}

func loadAPIKey(store *sql.DB, now func() time.Time) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !strings.HasPrefix(c.Request().URL.Path, "/api/") {
				return next(c)
			}
			present, raw := bearerToken(c.Request().Header.Values("Authorization"))
			if !present {
				return next(c)
			}
			key, err := apikeys.Authenticate(c.Request().Context(), store, raw, now())
			if err != nil {
				return rejectAPIKey(c, raw, err)
			}
			slog.DebugContext(c.Request().Context(), "API key authenticated",
				"api_key_id", key.ID, "api_key_name", key.Name, "path", c.Path())
			auth.SetAPIKeyInContext(c, &auth.APIKeyPrincipal{ID: key.ID, Name: key.Name, Scope: key.Scope})
			return next(c)
		}
	}
}

// bearerToken recognizes the scheme independently from credential syntax so malformed
// Bearer headers (including a missing separator like "Bearer<key>") cannot silently fall
// back to cookies. Multiple credentials are ambiguous and rejected when any is Bearer.
// Other authorization schemes are left to the proxy.
func bearerToken(headers []string) (present bool, token string) {
	for _, header := range headers {
		trimmed := strings.TrimLeft(header, " \t")
		if len(trimmed) < 6 || !strings.EqualFold(trimmed[:6], "Bearer") {
			continue
		}
		if len(headers) != 1 || len(header) < 7 || !strings.EqualFold(header[:7], "Bearer ") {
			return true, ""
		}
		return true, header[7:]
	}
	return false, ""
}

// rejectAPIKey logs only fixed reasons and routing metadata, never credentials or
// database errors (which could contain secret material). Storage failures remain 500.
func rejectAPIKey(c echo.Context, raw string, err error) error {
	if !errors.Is(err, apikeys.ErrInvalidAPIKey) && !errors.Is(err, apikeys.ErrAPIKeyExpired) {
		// Do not log the underlying error: driver errors may include SQL parameters.
		cause := "store_error"
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			cause = "request_canceled"
		}
		slog.ErrorContext(c.Request().Context(), "API key lookup failed",
			"cause", cause, "path", c.Path(), "remote_ip", c.RealIP())
		//nolint:wrapcheck // Terminal response
		return api.WriteInternalError(c, "authenticate API key", errors.New("credential store unavailable"))
	}
	reason, detail := "unknown", "Invalid API key"
	if _, parseErr := uuid.Parse(raw); parseErr != nil || len(raw) != 36 {
		reason = "malformed"
	}
	if errors.Is(err, apikeys.ErrAPIKeyExpired) {
		reason, detail = "expired", "API key has expired"
	}
	// c.Path is the route template, not the attacker-controlled URL/query string.
	slog.WarnContext(c.Request().Context(), "API key rejected",
		"reason", reason, "path", c.Path(), "remote_ip", c.RealIP())
	c.Response().Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
	return api.WriteUnauthorizedDetail(c, detail) //nolint:wrapcheck // Terminal response
}
