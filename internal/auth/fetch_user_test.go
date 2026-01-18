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

const (
	graphMeURL          = "https://graph.microsoft.com/v1.0/me"
	graphMeBody         = `{"id":"u1","displayName":"Ada"}`
	graphAdminGroupBody = `{"value":[{"@odata.type":"#microsoft.graph.group","id":"admins"}]}`
)

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

	client := newGraphClient(map[string]string{
		graphMeURL: graphMeBody,
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if user.ID != "u1" || user.Name != "Ada" {
		t.Fatalf("unexpected user: %#v", user)
	}
	if user.IsAdmin {
		t.Fatalf("expected non-admin user, got %#v", user)
	}
}

func TestFetchUserSetsAdminWhenInGroup(t *testing.T) {
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

	svc, err := NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := newGraphClient(map[string]string{
		graphMeURL:       graphMeBody,
		graphMemberOfURL: graphAdminGroupBody,
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if !user.IsAdmin {
		t.Fatalf("expected admin user, got %#v", user)
	}
}

func TestFetchUserRequiresUsersGroupForAdmin(t *testing.T) {
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

	svc, err := NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := newGraphClient(map[string]string{
		graphMeURL:       graphMeBody,
		graphMemberOfURL: graphAdminGroupBody,
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	user, err := svc.FetchUser(ctx, &oauth2.Token{AccessToken: "token"})
	if err != nil {
		t.Fatalf("fetch user: %v", err)
	}
	if user.IsAdmin {
		t.Fatalf("expected non-admin user, got %#v", user)
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

func TestFetchUserIgnoresGroupFetchError(t *testing.T) {
	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL:  "https://example.com/auth",
		TokenURL:      "https://example.com/token",
		RedirectURI:   "https://example.com/callback",
		ClientID:      "client",
		ClientSecret:  "secret",
		AdminsGroupID: "admins",
	}}

	svc, err := NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	client := &http.Client{Transport: roundTripper(func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case graphMeURL:
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
