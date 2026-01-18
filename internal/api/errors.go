// Package api defines JSON:API response helpers.
//
//revive:disable-next-line var-naming
package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// WriteUnauthorized writes a JSON:API unauthorized error response.
func WriteUnauthorized(c echo.Context) error {
	errResp := NewError(http.StatusUnauthorized, "Unauthorized", "Login required", "auth_required")
	c.Response().Header().Set(echo.HeaderContentType, JSONAPIContentType)
	if err := c.JSON(http.StatusUnauthorized, errResp); err != nil {
		return fmt.Errorf("write unauthorized response: %w", err)
	}
	return nil
}
