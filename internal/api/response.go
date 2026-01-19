// Package api defines JSON:API response helpers.
//
//revive:disable-next-line var-naming
package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// JSONAPIContentType is the JSON:API media type.
const JSONAPIContentType = "application/vnd.api+json"

// Resource represents a JSON:API resource object.
type Resource struct {
	ID         string      `json:"id,omitempty"`
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes,omitempty"`
}

// SingleResponse wraps a single resource.
type SingleResponse struct {
	Data Resource `json:"data"`
}

// CollectionResponse wraps a collection of resources.
type CollectionResponse struct {
	Data []Resource `json:"data"`
}

// MapResources maps items into JSON:API resources.
func MapResources[T any](items []T, build func(T) Resource) []Resource {
	resources := make([]Resource, 0, len(items))
	for _, item := range items {
		resources = append(resources, build(item))
	}
	return resources
}

// ErrorObject represents a JSON:API error.
type ErrorObject struct {
	Status string `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Code   string `json:"code,omitempty"`
}

// ErrorResponse wraps one or more errors.
type ErrorResponse struct {
	Errors []ErrorObject `json:"errors"`
}

// NewError builds a JSON:API error response.
func NewError(status int, title, detail, code string) ErrorResponse {
	return ErrorResponse{
		Errors: []ErrorObject{
			{
				Status: strconv.Itoa(status),
				Title:  title,
				Detail: detail,
				Code:   code,
			},
		},
	}
}

// ParseBookingDate parses a date query parameter, defaulting to today if empty.
func ParseBookingDate(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return time.Now().Format(time.DateOnly), nil
	}
	parsed, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return "", fmt.Errorf("parse booking date: %w", err)
	}
	return parsed.Format(time.DateOnly), nil
}

// RoomRequestParams contains common params extracted from room-related requests.
type RoomRequestParams struct {
	RoomID      string
	BookingDate string
}

// ParseRoomRequest extracts roomID and booking date from a request.
// Returns the params or an error if the date is invalid.
func ParseRoomRequest(roomID, dateParam string) (*RoomRequestParams, error) {
	bookingDate, err := ParseBookingDate(dateParam)
	if err != nil {
		return nil, err
	}
	return &RoomRequestParams{
		RoomID:      roomID,
		BookingDate: bookingDate,
	}, nil
}

// BuildINClause creates SQL placeholders and args for an IN clause.
// Returns the placeholders string (e.g., "?,?,?") and args slice with IDs.
func BuildINClause(ids []string) (placeholders string, args []any) {
	ph := make([]string, len(ids))
	args = make([]any, len(ids))
	for i, id := range ids {
		ph[i] = "?"
		args[i] = id
	}
	return strings.Join(ph, ","), args
}
