// Package api defines JSON:API response helpers.
//
//revive:disable-next-line var-naming
package api

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	// DefaultPageLimit is the page size used when page[limit] is absent.
	DefaultPageLimit = 20
	// MaxPageLimit is the largest accepted page size; larger requests are capped.
	MaxPageLimit = 100

	pageLimitParam  = "page[limit]"
	pageOffsetParam = "page[offset]"
)

// ErrInvalidPagination reports a malformed page[limit] or page[offset] value.
var ErrInvalidPagination = errors.New("invalid pagination parameter")

// Page is a validated offset-based page request.
type Page struct {
	Limit  int
	Offset int
}

// PaginationMeta is the JSON:API top-level meta for paginated collections.
type PaginationMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// PaginationLinks holds JSON:API pagination links. Prev and Next are null when unavailable.
type PaginationLinks struct {
	Self  string  `json:"self"`
	First string  `json:"first"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
	Last  string  `json:"last"`
}

// ParsePage reads page[limit] and page[offset] from the request query.
// Absent values default to DefaultPageLimit and 0. A limit above MaxPageLimit is capped so clients
// cannot request unbounded result sets. Non-integer values, a limit below 1 and a negative offset
// return ErrInvalidPagination wrapped with the offending parameter name.
func ParsePage(c echo.Context) (Page, error) {
	query := c.QueryParams()
	page := Page{Limit: DefaultPageLimit}

	if values, ok := query[pageLimitParam]; ok {
		limit, err := strconv.Atoi(values[0])
		if err != nil || limit < 1 {
			return Page{}, fmt.Errorf("%s must be a positive integer: %w", pageLimitParam, ErrInvalidPagination)
		}
		page.Limit = min(limit, MaxPageLimit)
	}

	if values, ok := query[pageOffsetParam]; ok {
		offset, err := strconv.Atoi(values[0])
		if err != nil || offset < 0 {
			return Page{}, fmt.Errorf("%s must be a non-negative integer: %w", pageOffsetParam, ErrInvalidPagination)
		}
		page.Offset = offset
	}

	return page, nil
}

// BuildPaginationLinks builds relative JSON:API pagination links from the request URL.
// All non-page query parameters are preserved; page[limit] always carries the effective limit.
func BuildPaginationLinks(u *url.URL, page Page, total int) PaginationLinks {
	link := func(offset int) string {
		query := u.Query()
		query.Set(pageLimitParam, strconv.Itoa(page.Limit))
		query.Set(pageOffsetParam, strconv.Itoa(offset))
		return u.Path + "?" + query.Encode()
	}

	lastOffset := 0
	if total > 0 {
		lastOffset = ((total - 1) / page.Limit) * page.Limit
	}

	links := PaginationLinks{
		Self:  link(page.Offset),
		First: link(0),
		Last:  link(lastOffset),
	}
	if page.Offset > 0 {
		prev := link(max(page.Offset-page.Limit, 0))
		links.Prev = &prev
	}
	// Compare via subtraction so a huge offset cannot overflow.
	if total-page.Offset > page.Limit {
		next := link(page.Offset + page.Limit)
		links.Next = &next
	}
	return links
}
