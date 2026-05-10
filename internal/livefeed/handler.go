package livefeed

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
)

// devOrigins are accepted in addition to the request Host so the Vite dev
// server (which serves the frontend on a different port) can connect during
// local development.
var devOrigins = map[string]struct{}{
	"http://localhost:5173":  {},
	"http://127.0.0.1:5173":  {},
	"https://localhost:5173": {},
	"https://127.0.0.1:5173": {},
}

// upgrader is shared because gorilla/websocket.Upgrader is documented as
// safe for concurrent use.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
	Error:           writeUpgradeError,
}

func writeUpgradeError(w http.ResponseWriter, _ *http.Request, _ int, reason error) {
	w.Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(
		api.NewError(http.StatusBadRequest, "Bad Request", "WebSocket upgrade failed", "bad_request"),
	); err != nil {
		slog.Debug("livefeed upgrade error write failed", "err", err)
	}
	if reason != nil {
		slog.Debug("livefeed upgrade rejected", "reason", reason)
	}
}

// checkOrigin restricts WebSocket upgrades to same-origin requests, with an
// allow-list for the Vite dev server in local development.
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// Non-browser clients (curl, wscat, tests) don't send an Origin
		// header. The endpoint is auth-gated, so we accept these.
		return true
	}
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if strings.EqualFold(parsed.Host, r.Host) {
		return true
	}
	if _, ok := devOrigins[strings.ToLower(origin)]; ok {
		return true
	}
	return false
}

// Handler returns an Echo handler that upgrades the request to a WebSocket
// connection and registers it with the hub. The route must be guarded by the
// existing auth middleware so the user is already loaded into the Echo
// context by the time this runs.
func Handler(hub *Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := auth.GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			// Upgrader has already written an HTTP error response.
			slog.Debug("livefeed upgrade failed", "user_id", user.ID, "err", err)
			return nil
		}

		client := newClient(hub, conn, user.ID)
		if !hub.Register(client) {
			_ = conn.Close() //nolint:errcheck // best-effort cleanup
			return nil
		}

		// Read pump runs on this goroutine; write pump on its own.
		go client.writePump()
		client.readPump()
		return nil
	}
}
