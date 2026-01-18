package auth

import (
	"testing"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestServiceTestUserDisabled(t *testing.T) {
	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}}

	svc, err := NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	if user := svc.TestUser(); user != nil {
		t.Fatalf("expected no test user, got %v", user)
	}
}

func TestServiceTestUserDefaults(t *testing.T) {
	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}, TestAuth: config.TestAuthConfig{
		Enabled: true,
	}}

	svc, err := NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	user := svc.TestUser()
	if user == nil {
		t.Fatal("expected test user")
	}
	if user.ID != "test-user" {
		t.Fatalf("unexpected user id: %s", user.ID)
	}
	if user.Name != "Test User" {
		t.Fatalf("unexpected user name: %s", user.Name)
	}
}

func TestServiceTestUserOverrides(t *testing.T) {
	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}, TestAuth: config.TestAuthConfig{
		Enabled:  true,
		UserID:   "u-123",
		UserName: "Ada Lovelace",
	}}

	svc, err := NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	user := svc.TestUser()
	if user == nil {
		t.Fatal("expected test user")
	}
	if user.ID != "u-123" {
		t.Fatalf("unexpected user id: %s", user.ID)
	}
	if user.Name != "Ada Lovelace" {
		t.Fatalf("unexpected user name: %s", user.Name)
	}
}
