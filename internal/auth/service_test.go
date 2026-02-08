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

	svc, err := NewService(cfg, nil)
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

	svc, err := NewService(cfg, nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	user := User{ID: "u1", Name: "Test User"}
	encoded, err := svc.EncodeUser(&user)
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

func TestServiceLocalOnlyMode(t *testing.T) {
	cfg := &config.Config{}
	svc, err := NewService(cfg, nil)
	if err != nil {
		t.Fatalf("expected no error for local-only config: %v", err)
	}
	if svc.AuthCodeURL("state") != "" {
		t.Fatal("expected empty auth URL in local-only mode")
	}
}

func TestUserGetID(t *testing.T) {
	user := &User{
		ID:          "user123",
		Name:        "Test User",
		Email:       "test@example.com",
		IsAdmin:     false,
		IsPermitted: true,
		AuthSource:  "internal",
	}

	if user.GetID() != "user123" {
		t.Fatalf("expected ID 'user123', got %q", user.GetID())
	}
}

func TestServiceDecodeStateTampered(t *testing.T) {
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

	_, err = svc.DecodeState("invalid-garbage-value")
	if err == nil {
		t.Fatal("expected error decoding tampered state")
	}
}

func TestServiceDecodeUserTampered(t *testing.T) {
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

	_, err = svc.DecodeUser("invalid-garbage-value")
	if err == nil {
		t.Fatal("expected error decoding tampered user cookie")
	}
}

func TestServiceUserRoundTripPreservesAllFields(t *testing.T) {
	cfg := &config.Config{}
	svc, err := NewService(cfg, nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	user := User{
		ID:          "user-42",
		Name:        "Alice Admin",
		Email:       "alice@example.com",
		IsAdmin:     true,
		IsPermitted: true,
		AuthSource:  "internal",
	}

	encoded, err := svc.EncodeUser(&user)
	if err != nil {
		t.Fatalf("encode user: %v", err)
	}

	decoded, err := svc.DecodeUser(encoded)
	if err != nil {
		t.Fatalf("decode user: %v", err)
	}

	if decoded.Email != user.Email {
		t.Errorf("email mismatch: got %q want %q", decoded.Email, user.Email)
	}
	if decoded.IsAdmin != user.IsAdmin {
		t.Errorf("is_admin mismatch: got %v want %v", decoded.IsAdmin, user.IsAdmin)
	}
	if decoded.AuthSource != user.AuthSource {
		t.Errorf("auth_source mismatch: got %q want %q", decoded.AuthSource, user.AuthSource)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
