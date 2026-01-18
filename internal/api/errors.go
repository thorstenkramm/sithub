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

// WriteForbidden writes a JSON:API forbidden error response.
func WriteForbidden(c echo.Context) error {
	errResp := NewError(http.StatusForbidden, "Forbidden", "Access denied", "forbidden")
	c.Response().Header().Set(echo.HeaderContentType, JSONAPIContentType)
	if err := c.JSON(http.StatusForbidden, errResp); err != nil {
		return fmt.Errorf("write forbidden response: %w", err)
	}
	return nil
}
