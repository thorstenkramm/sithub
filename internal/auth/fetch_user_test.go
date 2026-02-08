// Package auth provides authentication helpers for SitHub.
package auth

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"

	"github.com/thorstenkramm/sithub/internal/config"
)

type roundTripper func(*http.Request) (*http.Response, error)

func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

const (
	graphMeURLWithSelect = "https://graph.microsoft.com/v1.0/me?$select=id,displayName,mail,userPrincipalName"
	graphMeBody          = `{"id":"u1","displayName":"Ada","mail":"ada@example.com",` +
		`"userPrincipalName":"ada@example.com"}`
	graphAdminGroupBody = `{"value":[{"@odata.type":"#microsoft.graph.group","id":"admins"}]}`
	graphUsersGroupBody = `{"value":[{"@odata.type":"#microsoft.graph.group","id":"users"}]}`
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close() //nolint:errcheck // Cleanup function, error not critical
	})

	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			display_name TEXT NOT NULL,
			password_hash TEXT NOT NULL DEFAULT '',
			user_source TEXT NOT NULL CHECK (user_source IN ('internal', 'entraid')),
			entra_id TEXT NOT NULL DEFAULT '',
			is_admin INTEGER NOT NULL DEFAULT 0,
			last_login TEXT NOT NULL DEFAULT '',
			access_token TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE UNIQUE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_users_entra_id ON users(entra_id);
	`)
	if err != nil {
		t.Fatalf("create users table: %v", err)
	}
	return db
}

func TestFetchUserSuccess(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}}

	svc, err := NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := newGraphClient(map[string]string{
		graphMeURLWithSelect: graphMeBody,
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if user.Name != "Ada" {
		t.Fatalf("unexpected user name: %s", user.Name)
	}
	if user.Email != "ada@example.com" {
		t.Fatalf("unexpected email: %s", user.Email)
	}
	if user.AuthSource != "entraid" {
		t.Fatalf("expected entraid auth source, got %s", user.AuthSource)
	}
	if user.IsAdmin {
		t.Fatalf("expected non-admin user, got %#v", user)
	}
	if !user.IsPermitted {
		t.Fatalf("expected permitted user, got %#v", user)
	}
	if user.ID == "" {
		t.Fatalf("expected non-empty DB user ID")
	}
}

func TestFetchUserSetsAdminWhenInGroup(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		EntraID: config.EntraIDConfig{
			AuthorizeURL:  "https://example.com/auth",
			TokenURL:      "https://example.com/token",
			RedirectURI:   "https://example.com/callback",
			ClientID:      "client",
			ClientSecret:  "secret",
			AdminsGroupID: "admins",
		},
	}

	svc, err := NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := newGraphClient(map[string]string{
		graphMeURLWithSelect: graphMeBody,
		graphMemberOfURL:     graphAdminGroupBody,
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if !user.IsAdmin {
		t.Fatalf("expected admin user, got %#v", user)
	}
	if !user.IsPermitted {
		t.Fatalf("expected permitted user, got %#v", user)
	}
}

func TestFetchUserRequiresUsersGroupForAdmin(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		EntraID: config.EntraIDConfig{
			AuthorizeURL:  "https://example.com/auth",
			TokenURL:      "https://example.com/token",
			RedirectURI:   "https://example.com/callback",
			ClientID:      "client",
			ClientSecret:  "secret",
			UsersGroupID:  "users",
			AdminsGroupID: "admins",
		},
	}

	svc, err := NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := newGraphClient(map[string]string{
		graphMeURLWithSelect: graphMeBody,
		graphMemberOfURL:     graphAdminGroupBody,
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if user.IsAdmin {
		t.Fatalf("expected non-admin user, got %#v", user)
	}
	if user.IsPermitted {
		t.Fatalf("expected non-permitted user, got %#v", user)
	}
}

func TestFetchUserPermittedWhenUsersGroupMatches(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		EntraID: config.EntraIDConfig{
			AuthorizeURL: "https://example.com/auth",
			TokenURL:     "https://example.com/token",
			RedirectURI:  "https://example.com/callback",
			ClientID:     "client",
			ClientSecret: "secret",
			UsersGroupID: "users",
		},
	}

	svc, err := NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := newGraphClient(map[string]string{
		graphMeURLWithSelect: graphMeBody,
		graphMemberOfURL:     graphUsersGroupBody,
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if !user.IsPermitted {
		t.Fatalf("expected permitted user, got %#v", user)
	}
}

func TestFetchUserNotPermittedWhenUsersGroupMissing(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		EntraID: config.EntraIDConfig{
			AuthorizeURL: "https://example.com/auth",
			TokenURL:     "https://example.com/token",
			RedirectURI:  "https://example.com/callback",
			ClientID:     "client",
			ClientSecret: "secret",
			UsersGroupID: "users",
		},
	}

	svc, err := NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := newGraphClient(map[string]string{
		graphMeURLWithSelect: graphMeBody,
		graphMemberOfURL:     graphAdminGroupBody,
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if user.IsPermitted {
		t.Fatalf("expected non-permitted user, got %#v", user)
	}
}

func TestFetchUserStatusError(t *testing.T) {
	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}}

	svc, err := NewService(cfg, nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := &http.Client{Transport: roundTripper(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	})}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	if _, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestFetchUserIgnoresGroupFetchError(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL:  "https://example.com/auth",
		TokenURL:      "https://example.com/token",
		RedirectURI:   "https://example.com/callback",
		ClientID:      "client",
		ClientSecret:  "secret",
		AdminsGroupID: "admins",
	}}

	svc, err := NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := &http.Client{Transport: roundTripper(func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case graphMeURLWithSelect:
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(graphMeBody)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		case graphMemberOfURL:
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		}
	})}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if user.IsAdmin {
		t.Fatalf("expected non-admin user, got %#v", user)
	}
	if !user.IsPermitted {
		t.Fatalf("expected permitted user, got %#v", user)
	}
}

func TestFetchGroupIDsHandlesPagination(t *testing.T) {
	firstPage := `{"value":[{"@odata.type":"#microsoft.graph.group","id":"admins"}],` +
		`"@odata.nextLink":"https://graph.microsoft.com/v1.0/me/memberOf?$select=id&$skiptoken=abc"}`
	secondPage := `{"value":[{"@odata.type":"#microsoft.graph.group","id":"users"}]}`
	nextLink := "https://graph.microsoft.com/v1.0/me/memberOf?$select=id&$skiptoken=abc"

	client := newGraphClient(map[string]string{
		graphMemberOfURL: firstPage,
		nextLink:         secondPage,
	})

	svc := &Service{}
	groupIDs, err := svc.fetchGroupIDs(context.Background(), client)
	if err != nil {
		t.Fatalf("fetch groups: %v", err)
	}

	if len(groupIDs) != 2 {
		t.Fatalf("expected 2 groups, got %v", groupIDs)
	}
}

func TestRefreshPermissionsUpdatesFlags(t *testing.T) {
	db := setupTestDB(t)
	seedUserWithToken(t, db, "u1", "ada@example.com", "token")

	cfg := &config.Config{
		EntraID: config.EntraIDConfig{
			AuthorizeURL:  "https://example.com/auth",
			TokenURL:      "https://example.com/token",
			RedirectURI:   "https://example.com/callback",
			ClientID:      "client",
			ClientSecret:  "secret",
			UsersGroupID:  "users",
			AdminsGroupID: "admins",
		},
	}

	svc, err := NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	groupBody := `{"value":[{"@odata.type":"#microsoft.graph.group","id":"admins"},` +
		`{"@odata.type":"#microsoft.graph.group","id":"users"}]}`
	client := newGraphClient(map[string]string{
		graphMemberOfURL: groupBody,
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user := &User{ID: "u1"}

	if err := svc.RefreshPermissions(ctx, user); err != nil {
		t.Fatalf("refresh permissions: %v", err)
	}
	if !user.IsPermitted || !user.IsAdmin {
		t.Fatalf("expected permitted admin, got %#v", user)
	}
}

func TestRefreshPermissionsSkipsLocalUsers(t *testing.T) {
	cfg := &config.Config{
		EntraID: config.EntraIDConfig{
			AuthorizeURL: "https://example.com/auth",
			TokenURL:     "https://example.com/token",
			RedirectURI:  "https://example.com/callback",
			ClientID:     "client",
			ClientSecret: "secret",
			UsersGroupID: "users",
		},
	}

	svc, err := NewService(cfg, nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	user := &User{ID: "u1", AuthSource: "internal", IsPermitted: true}
	if err := svc.RefreshPermissions(context.Background(), user); err != nil {
		t.Fatalf("expected no error for local user: %v", err)
	}
	if !user.IsPermitted {
		t.Fatalf("expected local user to remain permitted")
	}
}

func TestRefreshPermissionsRequiresTokenWhenGroupsConfigured(t *testing.T) {
	db := setupTestDB(t)
	seedUserWithToken(t, db, "u1", "ada@example.com", "")

	cfg := &config.Config{
		EntraID: config.EntraIDConfig{
			AuthorizeURL: "https://example.com/auth",
			TokenURL:     "https://example.com/token",
			RedirectURI:  "https://example.com/callback",
			ClientID:     "client",
			ClientSecret: "secret",
			UsersGroupID: "users",
		},
	}

	svc, err := NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	user := &User{ID: "u1"}
	if err := svc.RefreshPermissions(context.Background(), user); err == nil {
		t.Fatalf("expected error for missing access token")
	}
}

func seedUserWithToken(t *testing.T, db *sql.DB, id, email, token string) {
	t.Helper()
	_, err := db.Exec(`
		INSERT INTO users (id, email, display_name, password_hash,
			user_source, entra_id, is_admin, last_login, access_token,
			created_at, updated_at)
		VALUES (?, ?, 'Test', '', 'entraid', '', 0, '', ?, datetime('now'), datetime('now'))`,
		id, email, token,
	)
	if err != nil {
		t.Fatalf("seed user with token: %v", err)
	}
}

func newGraphClient(responses map[string]string) *http.Client {
	return &http.Client{Transport: roundTripper(func(req *http.Request) (*http.Response, error) {
		body, ok := responses[req.URL.String()]
		if !ok {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	})}
}
