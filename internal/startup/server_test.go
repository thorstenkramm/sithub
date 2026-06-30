package startup

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/thorstenkramm/sithub/internal/areas"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/db"
	"github.com/thorstenkramm/sithub/internal/livefeed"
	"github.com/thorstenkramm/sithub/internal/middleware"
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

func TestSecurityHeadersPresent(t *testing.T) {
	e := echo.New()
	e.Use(echomw.SecureWithConfig(secureConfig()))
	e.Use(strictTransportSecurity())
	e.Use(contentSecurityPolicy())
	e.GET("/sec-test", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "https://sithub.example.com/sec-test", http.NoBody)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	h := rec.Result().Header
	checks := map[string]string{
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"Strict-Transport-Security": "max-age=31536000",
	}
	for header, want := range checks {
		if got := h.Get(header); got != want {
			t.Fatalf("header %s = %q, want %q", header, got, want)
		}
	}
	csp := h.Get("Content-Security-Policy")
	if csp == "" {
		t.Fatal("expected Content-Security-Policy header to be set")
	}
	if !strings.Contains(csp, "default-src 'self'") {
		t.Fatalf("CSP missing default-src 'self': %q", csp)
	}
	if !strings.Contains(csp, "https://fonts.googleapis.com") {
		t.Fatalf("CSP must allow Google Fonts stylesheet host: %q", csp)
	}
	if !strings.Contains(csp, "connect-src 'self' ws://sithub.example.com wss://sithub.example.com") {
		t.Fatalf("CSP must explicitly allow same-host live-feed WebSockets: %q", csp)
	}
	connectSrc := ""
	for _, directive := range strings.Split(csp, ";") {
		directive = strings.TrimSpace(directive)
		if strings.HasPrefix(directive, "connect-src ") {
			connectSrc = directive
			break
		}
	}
	if connectSrc == "" {
		t.Fatalf("CSP missing connect-src directive: %q", csp)
	}
	for _, source := range strings.Fields(connectSrc)[1:] {
		if source == "ws:" || source == "wss:" {
			t.Fatalf("CSP must not allow arbitrary WebSocket origins: %q", csp)
		}
	}
}

func TestBodyLimitRejectsOversizedBody(t *testing.T) {
	e := echo.New()
	e.Use(echomw.SecureWithConfig(secureConfig()))
	e.Use(strictTransportSecurity())
	e.Use(contentSecurityPolicy())
	e.Use(echomw.BodyLimitWithConfig(bodyLimitConfig()))
	e.POST("/bl-test", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	body := strings.NewReader(strings.Repeat("a", 3<<20)) // 3 MB > 2 MB limit
	req := httptest.NewRequest(http.MethodPost, "/bl-test", body)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413 for oversized body, got %d", rec.Code)
	}
	if got := rec.Result().Header.Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("oversized-response X-Frame-Options = %q, want DENY", got)
	}
	if got := rec.Result().Header.Get("Strict-Transport-Security"); got != "max-age=31536000" {
		t.Fatalf("oversized-response HSTS = %q, want max-age=31536000", got)
	}
}

func TestHTTPServerReadTimeoutClosesSlowRequestBody(t *testing.T) {
	e := echo.New()
	bodyRead := make(chan error, 1)
	e.POST("/slow", func(c echo.Context) error {
		_, err := io.ReadAll(c.Request().Body)
		bodyRead <- err
		if err != nil {
			return fmt.Errorf("read slow request body: %w", err)
		}
		return c.NoContent(http.StatusNoContent)
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := newHTTPServerWithTimeouts(ln.Addr().String(), e, serverTimeouts{
		ReadHeader: time.Second,
		Read:       50 * time.Millisecond,
		Write:      time.Second,
		Idle:       time.Second,
	})
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Serve(ln)
	}()
	t.Cleanup(func() {
		if err := server.Close(); err != nil {
			t.Errorf("close test server: %v", err)
		}
		<-errCh
	})

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close slow-client connection: %v", err)
		}
	}()

	_, err = conn.Write([]byte("POST /slow HTTP/1.1\r\nHost: example.test\r\nContent-Length: 10\r\n\r\n"))
	if err != nil {
		t.Fatalf("write headers: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	if _, err := conn.Write([]byte("0123456789")); err != nil {
		t.Logf("slow-client body write failed after timeout closed the connection: %v", err)
	}

	select {
	case err := <-bodyRead:
		if err == nil {
			t.Fatal("expected slow request body read to fail after ReadTimeout")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for slow request body to be closed")
	}
}

func TestBodyLimitSkipsAvatarUpload(t *testing.T) {
	e := echo.New()
	e.Use(echomw.BodyLimitWithConfig(bodyLimitConfig()))
	authService := newTestAuthService(t)
	e.Use(middleware.LoadUser(authService))
	avatarsDir := t.TempDir()
	registerRoutes(
		e, authService, &areas.Config{},
		t.TempDir(), avatarsDir, nil,
		notifications.NewNotifier(""), livefeed.NewHub(), nil,
	)

	body, contentType := multipartAvatarBody(t, paddedPNG(t, 3<<20))
	req := httptest.NewRequest(http.MethodPost, avatarUploadPath, body)
	req.Header.Set(echo.HeaderContentType, contentType)
	req.AddCookie(testUserCookie(t, authService, &auth.User{
		ID:          "user-1",
		Name:        "Avatar User",
		AuthSource:  "internal",
		IsPermitted: true,
	}))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusRequestEntityTooLarge {
		t.Fatalf("avatar upload must bypass the global body limit, got 413")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected authenticated 3 MB avatar upload to succeed, got %d", rec.Code)
	}
	if _, err := os.Stat(filepath.Join(avatarsDir, "user-1.png")); err != nil {
		t.Fatalf("expected avatar file to be saved: %v", err)
	}
}

func TestAvatarUploadHasRouteSpecificBodyLimit(t *testing.T) {
	e := echo.New()
	e.Use(echomw.BodyLimitWithConfig(bodyLimitConfig()))
	authService := newTestAuthService(t)
	e.Use(middleware.LoadUser(authService))
	registerRoutes(
		e, authService, &areas.Config{},
		t.TempDir(), t.TempDir(), nil,
		notifications.NewNotifier(""), livefeed.NewHub(), nil,
	)

	body, contentType := multipartAvatarBody(t, paddedPNG(t, 5<<20))
	req := httptest.NewRequest(http.MethodPost, avatarUploadPath, body)
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413 for avatar upload above 4 MB route cap, got %d", rec.Code)
	}
}

func TestBodyLimitAllowsNormalBookingRequest(t *testing.T) {
	e := echo.New()
	e.Use(echomw.BodyLimitWithConfig(bodyLimitConfig()))
	authService := newTestAuthService(t)
	e.Use(middleware.LoadUser(authService))
	store := setupStartupTestStore(t)
	registerRoutes(
		e, authService, testAreasConfig(),
		t.TempDir(), t.TempDir(), store,
		notifications.NewNotifier(""), livefeed.NewHub(), nil,
	)

	bookingDate := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	body := `{"data":{"type":"bookings","attributes":{"item_id":"desk-1","booking_date":"` + bookingDate + `"}}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, "application/vnd.api+json")
	req.AddCookie(testUserCookie(t, authService, &auth.User{
		ID:          "user-1",
		Name:        "Booking User",
		AuthSource:  "internal",
		IsPermitted: true,
	}))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected normal booking request to succeed under body limit, got %d: %s", rec.Code, rec.Body.String())
	}
}

func multipartAvatarBody(t *testing.T, avatar []byte) (body *bytes.Buffer, contentType string) {
	t.Helper()
	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("avatar", "avatar.png")
	if err != nil {
		t.Fatalf("create multipart avatar field: %v", err)
	}
	if _, err := part.Write(avatar); err != nil {
		t.Fatalf("write multipart avatar: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	return body, writer.FormDataContentType()
}

func setupStartupTestStore(t *testing.T) *sql.DB {
	t.Helper()
	store, err := db.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open test store: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close test store: %v", err)
		}
	})
	if err := db.RunMigrations(store); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	return store
}

func testAreasConfig() *areas.Config {
	return &areas.Config{Areas: []areas.Area{{
		ID:   "area-1",
		Name: "Area 1",
		ItemGroups: []areas.ItemGroup{{
			ID:   "ig-1",
			Name: "Item Group 1",
			Items: []areas.Item{{
				ID:   "desk-1",
				Name: "Desk 1",
			}},
		}},
	}}}
}

func paddedPNG(t *testing.T, size int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := range 10 {
		for y := range 10 {
			img.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	if buf.Len() < size {
		buf.Write(bytes.Repeat([]byte("x"), size-buf.Len()))
	}
	return buf.Bytes()
}

func writeAreasConfigIn(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "areas.yaml")
	if err := os.WriteFile(path, []byte("areas: []\n"), 0o600); err != nil {
		t.Fatalf("write areas config: %v", err)
	}
	return path
}

func setupTestRouter(t *testing.T) *echo.Echo {
	t.Helper()
	e := echo.New()
	authService := newTestAuthService(t)
	e.Use(middleware.LoadUser(authService))
	registerRoutes(
		e, authService, &areas.Config{},
		t.TempDir(), t.TempDir(), nil,
		notifications.NewNotifier(""), livefeed.NewHub(), nil,
	)
	return e
}

func TestFloorPlanRouteRequiresAuthentication(t *testing.T) {
	e := setupTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/floor-plans/plan.png", http.NoBody)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestFloorPlanPositionsWriteRouteRequiresAuthentication(t *testing.T) {
	e := setupTestRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/floor-plan-positions", http.NoBody)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestFloorPlanPositionsWriteRouteRequiresAdmin(t *testing.T) {
	e := echo.New()
	authService := newTestAuthService(t)
	e.Use(middleware.LoadUser(authService))
	registerRoutes(
		e, authService, &areas.Config{},
		t.TempDir(), t.TempDir(), nil,
		notifications.NewNotifier(""), livefeed.NewHub(), nil,
	)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/floor-plan-positions", http.NoBody)
	req.AddCookie(testUserCookie(t, authService, &auth.User{
		ID:          "user-1",
		Name:        "Regular User",
		AuthSource:  "internal",
		IsPermitted: true,
		IsAdmin:     false,
	}))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
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

func testUserCookie(t *testing.T, svc *auth.Service, user *auth.User) *http.Cookie {
	t.Helper()

	encodedUser, err := svc.EncodeUser(user)
	if err != nil {
		t.Fatalf("encode user: %v", err)
	}

	return &http.Cookie{Name: "sithub_user", Value: encodedUser}
}
