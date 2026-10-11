package auth

import "github.com/labstack/echo/v4"

const apiKeyContextKey = "api_key"

// APIKeyPrincipal identifies a machine credential, not a user. It has no user roles,
// ownership or account ID and must never be placed in the user context slot.
// It is safe for concurrent read-only use after construction.
type APIKeyPrincipal struct {
	ID    string
	Name  string
	Scope string
}

// HasScope matches exactly one supported scope, without inheriting any user roles.
func (p *APIKeyPrincipal) HasScope(scope string) bool {
	return p != nil && scope != "" && p.Scope == scope
}

// SetAPIKeyInContext stores a machine principal separately from the session user.
// Like Echo context operations, it must be called on the request's handler goroutine.
func SetAPIKeyInContext(c echo.Context, principal *APIKeyPrincipal) {
	c.Set(apiKeyContextKey, principal)
}

// GetAPIKeyFromContext returns the machine principal, or nil for session requests.
func GetAPIKeyFromContext(c echo.Context) *APIKeyPrincipal {
	principal, ok := c.Get(apiKeyContextKey).(*APIKeyPrincipal)
	if !ok {
		return nil
	}
	return principal
}
