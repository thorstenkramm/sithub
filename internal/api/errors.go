// Package api defines JSON:API response helpers.
//
//revive:disable-next-line var-naming
package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// WriteBadRequest writes a JSON:API bad request error response.
func WriteBadRequest(c echo.Context, detail string) error {
	errResp := NewError(http.StatusBadRequest, "Bad Request", detail, "bad_request")
	c.Response().Header().Set(echo.HeaderContentType, JSONAPIContentType)
	if err := c.JSON(http.StatusBadRequest, errResp); err != nil {
		return fmt.Errorf("write bad request response: %w", err)
	}
	return nil
}

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

// WriteNotFound writes a JSON:API not found error response.
func WriteNotFound(c echo.Context, detail string) error {
	errResp := NewError(http.StatusNotFound, "Not Found", detail, "not_found")
	c.Response().Header().Set(echo.HeaderContentType, JSONAPIContentType)
	if err := c.JSON(http.StatusNotFound, errResp); err != nil {
		return fmt.Errorf("write not found response: %w", err)
	}
	return nil
}

// WriteConflict writes a JSON:API conflict error response.
func WriteConflict(c echo.Context, detail string) error {
	errResp := NewError(http.StatusConflict, "Conflict", detail, "conflict")
	c.Response().Header().Set(echo.HeaderContentType, JSONAPIContentType)
	if err := c.JSON(http.StatusConflict, errResp); err != nil {
		return fmt.Errorf("write conflict response: %w", err)
	}
	return nil
}

// WriteUnsupportedMediaType writes a JSON:API unsupported media type error response.
func WriteUnsupportedMediaType(c echo.Context, detail string) error {
	errResp := NewError(http.StatusUnsupportedMediaType, "Unsupported Media Type", detail, "unsupported_media_type")
	c.Response().Header().Set(echo.HeaderContentType, JSONAPIContentType)
	if err := c.JSON(http.StatusUnsupportedMediaType, errResp); err != nil {
		return fmt.Errorf("write unsupported media type response: %w", err)
	}
	return nil
}
