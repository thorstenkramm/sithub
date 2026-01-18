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
	adminsGroup string
	usersGroup  string
}

// User represents an authenticated user.
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsAdmin     bool   `json:"is_admin"`
	IsPermitted bool   `json:"is_permitted"`
	AccessToken string `json:"access_token"`
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
		scopes := []string{"openid", "profile", "email", "User.Read"}
		if cfg.EntraID.AdminsGroupID != "" || cfg.EntraID.UsersGroupID != "" {
			scopes = append(scopes, "GroupMember.Read.All")
		}
		oauthConfig = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
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
		adminsGroup: cfg.EntraID.AdminsGroupID,
		usersGroup:  cfg.EntraID.UsersGroupID,
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://graph.microsoft.com/v1.0/me", http.NoBody)
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

	user := &User{
		ID:          graph.ID,
		Name:        graph.DisplayName,
		IsPermitted: s.usersGroup == "",
	}

	if s.adminsGroup != "" || s.usersGroup != "" {
		groupIDs, err := s.fetchGroupIDs(ctx, client)
		if err != nil {
			return user, nil
		}
		if s.usersGroup != "" {
			user.IsPermitted = isGroupMember(groupIDs, s.usersGroup)
		}
		user.IsAdmin = s.isAdminGroupMember(groupIDs)
	}

	return user, nil
}

// RefreshPermissions re-evaluates group membership for the given user.
func (s *Service) RefreshPermissions(ctx context.Context, user *User) error {
	if user == nil || s.usersGroup == "" {
		return nil
	}
	if s.oauthConfig == nil {
		return fmt.Errorf("refresh permissions: missing oauth config")
	}
	if user.AccessToken == "" {
		return fmt.Errorf("refresh permissions: missing access token")
	}

	client := s.oauthConfig.Client(ctx, &oauth2.Token{AccessToken: user.AccessToken})
	groupIDs, err := s.fetchGroupIDs(ctx, client)
	if err != nil {
		return err
	}
	user.IsPermitted = isGroupMember(groupIDs, s.usersGroup)
	user.IsAdmin = s.isAdminGroupMember(groupIDs)
	return nil
}

type graphMemberOfResponse struct {
	Value    []graphGroup `json:"value"`
	NextLink string       `json:"@odata.nextLink"`
}

type graphGroup struct {
	ODataType string `json:"@odata.type"`
	ID        string `json:"id"`
}

const graphMemberOfURL = "https://graph.microsoft.com/v1.0/me/memberOf?$select=id"

func (s *Service) fetchGroupIDs(ctx context.Context, client *http.Client) ([]string, error) {
	var ids []string
	url := graphMemberOfURL

	for url != "" {
		pageIDs, nextLink, err := s.fetchGroupPage(ctx, client, url)
		if err != nil {
			return nil, err
		}
		ids = append(ids, pageIDs...)
		url = nextLink
	}

	return ids, nil
}

func (s *Service) fetchGroupPage(
	ctx context.Context,
	client *http.Client,
	url string,
) (ids []string, nextLink string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, "", fmt.Errorf("build groups request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("fetch groups: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("fetch groups: status %d", resp.StatusCode)
	}

	var body graphMemberOfResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, "", fmt.Errorf("decode groups: %w", err)
	}

	ids = make([]string, 0, len(body.Value))
	for _, group := range body.Value {
		if group.ID == "" {
			continue
		}
		ids = append(ids, group.ID)
	}

	return ids, body.NextLink, nil
}

func (s *Service) isAdminGroupMember(groupIDs []string) bool {
	adminMatch := false
	userMatch := s.usersGroup == ""

	for _, id := range groupIDs {
		if id == s.adminsGroup {
			adminMatch = true
		}
		if s.usersGroup != "" && id == s.usersGroup {
			userMatch = true
		}
	}

	return adminMatch && userMatch
}

func isGroupMember(groupIDs []string, target string) bool {
	for _, id := range groupIDs {
		if id == target {
			return true
		}
	}
	return false
}

// NewState creates a random OAuth state value.
func NewState() (string, error) {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", fmt.Errorf("generate state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
