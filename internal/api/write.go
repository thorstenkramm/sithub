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
