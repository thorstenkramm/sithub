package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/users"
)

type localLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// dummyHash is a pre-computed bcrypt hash used to prevent timing-based user enumeration.
// When a login attempt targets a non-existent email, we still run bcrypt.CompareHashAndPassword
// against this dummy so the response time is indistinguishable from a wrong-password attempt.
var dummyHash = mustHash("timing-safe-dummy-value")

func mustHash(password string) string {
	h, err := users.HashPassword(password)
	if err != nil {
		panic("failed to compute dummy bcrypt hash: " + err.Error())
	}
	return h
}

// LocalLoginHandler handles POST /api/v1/auth/login for local email/password login.
func LocalLoginHandler(svc *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req localLoginRequest
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return api.WriteBadRequest(c, "Invalid request body")
		}

		email := strings.TrimSpace(req.Email)
		password := req.Password
		if email == "" || password == "" {
			return api.WriteBadRequest(c, "Email and password are required")
		}
		if _, err := mail.ParseAddress(email); err != nil {
			return api.WriteBadRequest(c, "Invalid email format")
		}

		ctx := c.Request().Context()

		rec, err := users.FindByEmail(ctx, svc.store, email)
		if errors.Is(err, users.ErrUserNotFound) {
			// Run dummy bcrypt to prevent timing-based user enumeration.
			_ = users.VerifyPassword(dummyHash, password) //nolint:errcheck // Intentional dummy
			return jsonAPIError(
				c, http.StatusUnauthorized, "Unauthorized",
				"Invalid email or password", "invalid_credentials",
			)
		}
		if err != nil {
			return fmt.Errorf("find user by email: %w", err)
		}

		if rec.UserSource != "internal" {
			return jsonAPIError(c, http.StatusUnauthorized, "Unauthorized",
				"This account uses Entra ID. Please sign in with Entra ID.", "wrong_auth_source")
		}

		if err := users.VerifyPassword(rec.PasswordHash, password); err != nil {
			return jsonAPIError(
				c, http.StatusUnauthorized, "Unauthorized",
				"Invalid email or password", "invalid_credentials",
			)
		}

		// Record login timestamp (best-effort, don't fail the login)
		_ = users.UpdateLastLogin(ctx, svc.store, rec.ID) //nolint:errcheck // Best-effort

		user := &User{
			ID:          rec.ID,
			Name:        rec.DisplayName,
			Email:       rec.Email,
			IsAdmin:     rec.IsAdmin,
			IsPermitted: true,
			AuthSource:  "internal",
		}

		encodedUser, err := svc.EncodeUser(user)
		if err != nil {
			return jsonAPIError(
				c, http.StatusInternalServerError, "Server Error",
				"Failed to store session", "session_error",
			)
		}

		userCookie := newCookie(userCookieName, encodedUser, c.Scheme() == schemeHTTPS)
		c.SetCookie(userCookie)

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

func userRole(user *User) string {
	if user.IsAdmin {
		return "admin"
	}
	return "user"
}
