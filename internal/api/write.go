// Package api defines JSON:API response helpers.
//
//revive:disable-next-line var-naming
package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// WriteCollection writes a JSON:API collection response.
func WriteCollection(c echo.Context, resources []Resource, errLabel string) error {
	resp := CollectionResponse{Data: resources}
	c.Response().Header().Set(echo.HeaderContentType, JSONAPIContentType)
	if err := c.JSON(http.StatusOK, resp); err != nil {
		return fmt.Errorf("%s: %w", errLabel, err)
	}
	return nil
}

// WriteSingle writes a JSON:API single-resource response with the given status.
func WriteSingle(c echo.Context, status int, resource Resource, errLabel string) error {
	resp := SingleResponse{Data: resource}
	c.Response().Header().Set(echo.HeaderContentType, JSONAPIContentType)
	if err := c.JSON(status, resp); err != nil {
		return fmt.Errorf("%s: %w", errLabel, err)
	}
	return nil
}

// WriteInternalError logs the error and writes a JSON:API 500 response.
func WriteInternalError(c echo.Context, label string, origErr error) error {
	_ = origErr // logged by middleware
	errResp := NewError(http.StatusInternalServerError, "Server Error", label, "internal_error")
	c.Response().Header().Set(echo.HeaderContentType, JSONAPIContentType)
	if err := c.JSON(http.StatusInternalServerError, errResp); err != nil {
		return fmt.Errorf("write internal error response: %w", err)
	}
	return nil
}
