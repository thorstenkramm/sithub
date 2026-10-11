package apikeys

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
)

const resourceType = "api-keys"

func toResource(k Key) api.Resource { //nolint:gocritic // Value needed for MapResources
	return api.Resource{Type: resourceType, ID: k.ID, Attributes: keyAttributes(k)}
}

func keyAttributes(k Key) map[string]interface{} { //nolint:gocritic // Small immutable value
	return map[string]interface{}{
		"name": k.Name, "key_prefix": k.KeyPrefix, "scope": k.Scope,
		"expires_at": k.ExpiresAt, "created_at": k.CreatedAt, "created_by": k.CreatedBy,
	}
}

// ListHandler lists public API key metadata for administrators.
func ListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		keys, err := List(c.Request().Context(), db)
		if err != nil {
			return api.WriteInternalError(c, "list API keys", err)
		}
		return api.WriteCollection(c, api.MapResources(keys, toResource), "write API keys response")
	}
}

type createRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Name            string `json:"name"`
			Scope           string `json:"scope"`
			ExpiresInMonths int    `json:"expires_in_months"`
		} `json:"attributes"`
	} `json:"data"`
}

// CreateHandler generates a key and returns its plaintext only in the create response.
func CreateHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := auth.GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}
		var req createRequest
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return api.WriteBadRequest(c, "Invalid request body")
		}
		if req.Data.Type != resourceType {
			return api.WriteBadRequest(c, "data.type must be api-keys")
		}
		a := req.Data.Attributes
		k, plaintext, err := Create(c.Request().Context(), db, time.Now(), CreateInput{
			Name: a.Name, Scope: a.Scope, ExpiresInMonths: a.ExpiresInMonths, CreatedBy: user.ID,
		})
		if err != nil {
			return writeCreateError(c, err)
		}
		attrs := keyAttributes(k)
		attrs["key"] = plaintext
		resource := api.Resource{Type: resourceType, ID: k.ID, Attributes: attrs}
		c.Response().Header().Set("Cache-Control", "no-store")
		return api.WriteSingle(c, http.StatusCreated, resource, "write created API key response")
	}
}

func writeCreateError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrInvalidName):
		detail := strings.TrimSuffix(err.Error(), ": "+ErrInvalidName.Error())
		return api.WriteBadRequest(c, detail) //nolint:wrapcheck // Terminal response
	case errors.Is(err, ErrInvalidScope), errors.Is(err, ErrInvalidExpiry):
		return api.WriteBadRequest(c, err.Error()) //nolint:wrapcheck // Terminal response
	default:
		return api.WriteInternalError(c, "create API key", err) //nolint:wrapcheck // Terminal response
	}
}

// DeleteHandler permanently revokes an API key.
func DeleteHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := Delete(c.Request().Context(), db, c.Param("id")); err != nil {
			if errors.Is(err, ErrNotFound) {
				return api.WriteNotFound(c, "API key not found")
			}
			return api.WriteInternalError(c, "delete API key", err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
