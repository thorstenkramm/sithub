// Package auth provides authentication helpers for SitHub.
package auth

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/oauth2"

	"github.com/thorstenkramm/sithub/internal/config"
)

type roundTripper func(*http.Request) (*http.Response, error)

func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

func TestFetchUserSuccess(t *testing.T) {
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

	client := &http.Client{Transport: roundTripper(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://graph.microsoft.com/v1.0/me" {
			return nil, nil
		}
		body := `{"id":"u1","displayName":"Ada"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	})}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if user.ID != "u1" || user.Name != "Ada" {
		t.Fatalf("unexpected user: %#v", user)
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

	svc, err := NewService(cfg)
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
