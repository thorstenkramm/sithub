package auth

import (
	"os"
	"testing"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestServiceStateRoundTrip(t *testing.T) {
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

	state, err := NewState()
	if err != nil {
		t.Fatalf("new state: %v", err)
	}

	encoded, err := svc.EncodeState(state)
	if err != nil {
		t.Fatalf("encode state: %v", err)
	}

	decoded, err := svc.DecodeState(encoded)
	if err != nil {
		t.Fatalf("decode state: %v", err)
	}

	if decoded != state {
		t.Fatalf("state mismatch: got %q want %q", decoded, state)
	}
}

func TestServiceUserRoundTrip(t *testing.T) {
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

	user := User{ID: "u1", Name: "Test User"}
	encoded, err := svc.EncodeUser(user)
	if err != nil {
		t.Fatalf("encode user: %v", err)
	}

	decoded, err := svc.DecodeUser(encoded)
	if err != nil {
		t.Fatalf("decode user: %v", err)
	}

	if decoded.ID != user.ID || decoded.Name != user.Name {
		t.Fatalf("user mismatch: got %+v want %+v", decoded, user)
	}
}

func TestServiceMissingConfig(t *testing.T) {
	cfg := &config.Config{}
	if _, err := NewService(cfg); err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestServiceMissingConfigWithTestAuth(t *testing.T) {
	cfg := &config.Config{TestAuth: config.TestAuthConfig{
		Enabled: true,
	}}
	if _, err := NewService(cfg); err != nil {
		t.Fatalf("expected service when test auth enabled, got %v", err)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
