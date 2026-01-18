package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"

	"github.com/thorstenkramm/sithub/internal/config"
)

// ErrInvalidState indicates an OAuth state mismatch.
var ErrInvalidState = errors.New("invalid oauth state")

const (
	stateCookieName = "sithub_oauth_state"
	userCookieName  = "sithub_user"
)

// Service handles Entra ID authentication and cookie encoding.
type Service struct {
	oauthConfig *oauth2.Config
	cookieCodec *securecookie.SecureCookie
	testAuth    config.TestAuthConfig
}

// User represents an authenticated user.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type graphUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// NewService configures an authentication service from config.
func NewService(cfg *config.Config) (*Service, error) {
	authURL := cfg.EntraID.AuthorizeURL
	tokenURL := cfg.EntraID.TokenURL
	redirectURL := cfg.EntraID.RedirectURI
	clientID := cfg.EntraID.ClientID
	clientSecret := cfg.EntraID.ClientSecret

	missingAuthConfig := authURL == "" || tokenURL == "" || redirectURL == "" || clientID == "" || clientSecret == ""
	if missingAuthConfig && !cfg.TestAuth.Enabled {
		return nil, fmt.Errorf("auth config: %w", config.ErrMissingEntraIDConfig)
	}

	var oauthConfig *oauth2.Config
	if !missingAuthConfig {
		oauthConfig = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email", "User.Read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  authURL,
				TokenURL: tokenURL,
			},
		}
	}

	hashKey := make([]byte, 32)
	blockKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, hashKey); err != nil {
		return nil, fmt.Errorf("generate hash key: %w", err)
	}
	if _, err := io.ReadFull(rand.Reader, blockKey); err != nil {
		return nil, fmt.Errorf("generate block key: %w", err)
	}

	return &Service{
		oauthConfig: oauthConfig,
		cookieCodec: securecookie.New(hashKey, blockKey),
		testAuth:    cfg.TestAuth,
	}, nil
}

// AuthCodeURL returns the authorization URL for the given state.
func (s *Service) AuthCodeURL(state string) string {
	if s.oauthConfig == nil {
		return ""
	}
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange trades the authorization code for an access token.
func (s *Service) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange token: %w", err)
	}
	return token, nil
}

// EncodeState signs and encodes an OAuth state value.
func (s *Service) EncodeState(state string) (string, error) {
	encoded, err := s.cookieCodec.Encode(stateCookieName, state)
	if err != nil {
		return "", fmt.Errorf("encode state: %w", err)
	}
	return encoded, nil
}

// DecodeState decodes a signed OAuth state value.
func (s *Service) DecodeState(value string) (string, error) {
	var state string
	if err := s.cookieCodec.Decode(stateCookieName, value, &state); err != nil {
		return "", fmt.Errorf("decode state: %w", err)
	}
	return state, nil
}

// EncodeUser encodes a user into a signed cookie value.
func (s *Service) EncodeUser(user User) (string, error) {
	encoded, err := s.cookieCodec.Encode(userCookieName, user)
	if err != nil {
		return "", fmt.Errorf("encode user: %w", err)
	}
	return encoded, nil
}

// DecodeUser decodes a user from a signed cookie value.
func (s *Service) DecodeUser(value string) (*User, error) {
	var user User
	if err := s.cookieCodec.Decode(userCookieName, value, &user); err != nil {
		return nil, fmt.Errorf("decode user: %w", err)
	}
	return &user, nil
}

// FetchUser retrieves the current user profile from Microsoft Graph.
func (s *Service) FetchUser(ctx context.Context, token *oauth2.Token) (*User, error) {
	client := s.oauthConfig.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return nil, fmt.Errorf("build user request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch user: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch user: status %d", resp.StatusCode)
	}

	var graph graphUser
	if err := json.NewDecoder(resp.Body).Decode(&graph); err != nil {
		return nil, fmt.Errorf("decode user: %w", err)
	}

	return &User{ID: graph.ID, Name: graph.DisplayName}, nil
}

// NewState creates a random OAuth state value.
func NewState() (string, error) {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", fmt.Errorf("generate state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
