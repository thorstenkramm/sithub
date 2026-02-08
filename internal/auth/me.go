package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/users"
)

// GetUserFromContext retrieves the authenticated user from the Echo context.
// Returns nil if no user is present or the type assertion fails.
func GetUserFromContext(c echo.Context) *User {
	user, ok := c.Get("user").(*User)
	if !ok {
		return nil
	}
	return user
}

// MeHandler returns the authenticated user profile.
func MeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		resp := api.SingleResponse{
			Data: api.Resource{
				Type: "users",
				ID:   user.ID,
				Attributes: map[string]interface{}{
					"display_name": user.Name,
					"email":        user.Email,
					"is_admin":     user.IsAdmin,
					"auth_source":  user.AuthSource,
					"role":         userRole(user),
				},
			},
		}

		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		return c.JSON(http.StatusOK, resp)
	}
}

type updateMeRequest struct {
	Data struct {
		Attributes struct {
			CurrentPassword string `json:"current_password"`
			NewPassword     string `json:"new_password"`
		} `json:"attributes"`
	} `json:"data"`
}

const (
	minPasswordLength  = 14
	userSourceInternal = "internal"
)

// UpdateMeHandler handles PATCH /api/v1/me for self-service password change.
func UpdateMeHandler(svc *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		if user.AuthSource != userSourceInternal {
			return api.WriteBadRequest(c, "Password change is only available for local accounts")
		}

		currentPassword, newPassword, err := parsePasswordChangeRequest(c)
		if err != nil {
			return err
		}

		if err := validateAndUpdatePassword(c, svc, user, currentPassword, newPassword); err != nil {
			return err
		}

		return respondWithUserProfile(c, user)
	}
}

func parsePasswordChangeRequest(c echo.Context) (currentPassword, newPassword string, err error) {
	var req updateMeRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return "", "", api.WriteBadRequest(c, "Invalid request body") //nolint:wrapcheck // Terminal response
	}

	currentPassword = strings.TrimSpace(req.Data.Attributes.CurrentPassword)
	newPassword = req.Data.Attributes.NewPassword

	if currentPassword == "" || newPassword == "" {
		//nolint:wrapcheck // Terminal response
		return "", "", api.WriteBadRequest(c, "current_password and new_password are required")
	}
	if len(newPassword) < minPasswordLength {
		msg := fmt.Sprintf("Password must be at least %d characters", minPasswordLength)
		return "", "", api.WriteBadRequest(c, msg) //nolint:wrapcheck // Terminal response
	}

	return currentPassword, newPassword, nil
}

func validateAndUpdatePassword(c echo.Context, svc *Service, user *User, currentPassword, newPassword string) error {
	ctx := c.Request().Context()

	rec, err := users.FindByID(ctx, svc.store, user.ID)
	if errors.Is(err, users.ErrUserNotFound) {
		return api.WriteUnauthorized(c) //nolint:wrapcheck // Terminal response
	}
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}

	if err := users.VerifyPassword(rec.PasswordHash, currentPassword); err != nil {
		return jsonAPIError(
			c, http.StatusUnauthorized, "Unauthorized",
			"Current password is incorrect", "invalid_password",
		)
	}

	hash, err := users.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := users.UpdatePasswordHash(ctx, svc.store, user.ID, hash); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}

func respondWithUserProfile(c echo.Context, user *User) error {
	resp := api.SingleResponse{
		Data: api.Resource{
			Type: "users",
			ID:   user.ID,
			Attributes: map[string]interface{}{
				"display_name": user.Name,
				"email":        user.Email,
				"is_admin":     user.IsAdmin,
				"auth_source":  user.AuthSource,
				"role":         userRole(user),
			},
		},
	}

	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	return c.JSON(http.StatusOK, resp) //nolint:wrapcheck // Terminal response
}
