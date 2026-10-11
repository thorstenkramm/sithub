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
	return writeCollection(c, CollectionResponse{Data: resources}, errLabel)
}

// WritePaginatedCollection writes a JSON:API collection response with pagination meta and links
// derived from the request URL.
func WritePaginatedCollection(c echo.Context, resources []Resource, page Page, total int, errLabel string) error {
	links := BuildPaginationLinks(c.Request().URL, page, total)
	return writeCollection(c, CollectionResponse{
		Data:  resources,
		Meta:  &PaginationMeta{Total: total, Limit: page.Limit, Offset: page.Offset},
		Links: &links,
	}, errLabel)
}

func writeCollection(c echo.Context, resp CollectionResponse, errLabel string) error {
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
