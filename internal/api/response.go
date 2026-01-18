// Package api defines JSON:API response helpers.
package api

import "strconv"

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
