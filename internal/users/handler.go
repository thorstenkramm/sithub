package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

const (
	minPasswordLength  = 14
	userSourceInternal = "internal"
)

// UserAttributes represents user resource attributes for JSON:API responses.
type UserAttributes struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	IsAdmin     bool   `json:"is_admin"`
	AuthSource  string `json:"auth_source"`
	Role        string `json:"role"`
	LastLogin   string `json:"last_login"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListHandler returns a handler for listing all users.
func ListHandler(store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		records, err := ListAll(ctx, store)
		if err != nil {
			return fmt.Errorf("list users: %w", err)
		}

		resources := api.MapResources(records, func(rec Record) api.Resource {
			return api.Resource{
				Type:       "users",
				ID:         rec.ID,
				Attributes: recordToAttributes(&rec),
			}
		})

		resp := api.CollectionResponse{Data: resources}
		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		return c.JSON(http.StatusOK, resp)
	}
}

// GetHandler returns a handler for getting a single user by ID.
func GetHandler(store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")
		if userID == "" {
			return api.WriteBadRequest(c, "User ID is required")
		}

		ctx := c.Request().Context()
		rec, err := FindByID(ctx, store, userID)
		if errors.Is(err, ErrUserNotFound) {
			return api.WriteNotFound(c, "User not found")
		}
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}

		resp := api.SingleResponse{
			Data: api.Resource{
				Type:       "users",
				ID:         rec.ID,
				Attributes: recordToAttributes(rec),
			},
		}
		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		return c.JSON(http.StatusOK, resp)
	}
}

type createUserRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Email       string `json:"email"`
			DisplayName string `json:"display_name"`
			Password    string `json:"password"`
			IsAdmin     bool   `json:"is_admin"`
		} `json:"attributes"`
	} `json:"data"`
}

// CreateHandler returns a handler for creating a new local user (admin only).
func CreateHandler(store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req createUserRequest
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return api.WriteBadRequest(c, "Invalid request body")
		}

		email := strings.TrimSpace(req.Data.Attributes.Email)
		displayName := strings.TrimSpace(req.Data.Attributes.DisplayName)
		password := req.Data.Attributes.Password

		if email == "" || displayName == "" || password == "" {
			return api.WriteBadRequest(c, "email, display_name, and password are required")
		}
		if _, err := mail.ParseAddress(email); err != nil {
			return api.WriteBadRequest(c, "Invalid email format")
		}
		if len(password) < minPasswordLength {
			return api.WriteBadRequest(c, fmt.Sprintf("Password must be at least %d characters", minPasswordLength))
		}

		hash, err := HashPassword(password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		ctx := c.Request().Context()
		rec, err := CreateLocalUser(ctx, store, email, displayName, hash, req.Data.Attributes.IsAdmin)
		if errors.Is(err, ErrEmailConflict) {
			return api.WriteConflict(c, "A user with this email already exists")
		}
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		resp := api.SingleResponse{
			Data: api.Resource{
				Type:       "users",
				ID:         rec.ID,
				Attributes: recordToAttributes(rec),
			},
		}
		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		return c.JSON(http.StatusCreated, resp)
	}
}

type updateUserRequest struct {
	Data struct {
		Attributes struct {
			Email       *string `json:"email,omitempty"`
			DisplayName *string `json:"display_name,omitempty"`
			IsAdmin     *bool   `json:"is_admin,omitempty"`
			Password    *string `json:"password,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// UpdateHandler returns a handler for updating a user (admin only).
func UpdateHandler(store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")
		if userID == "" {
			return api.WriteBadRequest(c, "User ID is required")
		}

		var req updateUserRequest
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return api.WriteBadRequest(c, "Invalid request body")
		}

		ctx := c.Request().Context()

		if err := handlePasswordResetIfRequested(ctx, c, store, userID, req.Data.Attributes.Password); err != nil {
			return err
		}

		rec, err := applyFieldUpdates(ctx, c, store, userID, req.Data.Attributes)
		if err != nil {
			return err
		}

		return respondWithUpdatedUser(c, rec)
	}
}

func handlePasswordResetIfRequested(
	ctx context.Context, c echo.Context, store *sql.DB, userID string, password *string,
) error {
	if password == nil {
		return nil
	}

	rec, err := FindByID(ctx, store, userID)
	if errors.Is(err, ErrUserNotFound) {
		return api.WriteNotFound(c, "User not found") //nolint:wrapcheck // Terminal response
	}
	if err != nil {
		return fmt.Errorf("find user for password reset: %w", err)
	}

	if rec.UserSource != userSourceInternal {
		//nolint:wrapcheck // Terminal response
		return api.WriteBadRequest(c, "Password reset is only available for local users")
	}

	if len(*password) < minPasswordLength {
		msg := fmt.Sprintf("Password must be at least %d characters", minPasswordLength)
		return api.WriteBadRequest(c, msg) //nolint:wrapcheck // Terminal response
	}

	hash, err := HashPassword(*password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := UpdatePasswordHash(ctx, store, userID, hash); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}

func applyFieldUpdates(ctx context.Context, c echo.Context, store *sql.DB, userID string, attrs struct {
	Email       *string `json:"email,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	IsAdmin     *bool   `json:"is_admin,omitempty"`
	Password    *string `json:"password,omitempty"`
}) (*Record, error) {
	fields := UpdateFields{
		Email:       trimPtr(attrs.Email),
		DisplayName: trimPtr(attrs.DisplayName),
		IsAdmin:     attrs.IsAdmin,
	}

	rec, err := UpdateUser(ctx, store, userID, fields)
	if errors.Is(err, ErrUserNotFound) {
		return nil, api.WriteNotFound(c, "User not found") //nolint:wrapcheck // Terminal response
	}
	if errors.Is(err, ErrEmailConflict) {
		//nolint:wrapcheck // Terminal response
		return nil, api.WriteConflict(c, "A user with this email already exists")
	}
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return rec, nil
}

func respondWithUpdatedUser(c echo.Context, rec *Record) error {
	resp := api.SingleResponse{
		Data: api.Resource{
			Type:       "users",
			ID:         rec.ID,
			Attributes: recordToAttributes(rec),
		},
	}
	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	return c.JSON(http.StatusOK, resp) //nolint:wrapcheck // Terminal response
}

// DeleteHandler returns a handler for deleting a local user (admin only).
func DeleteHandler(store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")
		if userID == "" {
			return api.WriteBadRequest(c, "User ID is required")
		}

		ctx := c.Request().Context()

		// Only local users can be deleted
		rec, err := FindByID(ctx, store, userID)
		if errors.Is(err, ErrUserNotFound) {
			return api.WriteNotFound(c, "User not found")
		}
		if err != nil {
			return fmt.Errorf("find user for delete: %w", err)
		}

		// Prevent admin from deleting themselves.
		// The user struct has a public ID field; use reflection-free approach.
		if caller := currentUserID(c); caller == userID {
			return api.WriteBadRequest(c, "Cannot delete your own account")
		}

		if rec.UserSource != userSourceInternal {
			return api.WriteBadRequest(c, "Only local users can be deleted (Entra ID users are managed externally)")
		}

		if err := DeleteUser(ctx, store, userID); err != nil {
			return fmt.Errorf("delete user: %w", err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func recordToAttributes(rec *Record) UserAttributes {
	role := "user"
	if rec.IsAdmin {
		role = "admin"
	}
	return UserAttributes{
		Email:       rec.Email,
		DisplayName: rec.DisplayName,
		IsAdmin:     rec.IsAdmin,
		AuthSource:  rec.UserSource,
		Role:        role,
		LastLogin:   rec.LastLogin,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}
}

// currentUserID extracts the authenticated user's ID from the Echo context.
// It accesses the "user" key set by the LoadUser middleware without importing
// the auth package (which would create an import cycle).
func currentUserID(c echo.Context) string {
	type hasID interface{ GetID() string }
	if u, ok := c.Get("user").(hasID); ok {
		return u.GetID()
	}
	return ""
}

func trimPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	return &trimmed
}
