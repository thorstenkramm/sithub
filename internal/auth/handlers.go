package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// LoginHandler starts the Entra ID authorization flow.
func LoginHandler(svc *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		state, err := NewState()
		if err != nil {
			return jsonAPIError(c, http.StatusInternalServerError, "Server Error", "Failed to start login", "login_init")
		}

		encoded, err := svc.EncodeState(state)
		if err != nil {
			return jsonAPIError(c, http.StatusInternalServerError, "Server Error", "Failed to store login state", "login_state")
		}

		cookie := newCookie(stateCookieName, encoded, c.Scheme() == "https")
		c.SetCookie(cookie)

		authURL := svc.AuthCodeURL(state)
		if authURL == "" {
			detail := "Entra ID login is not configured"
			return jsonAPIError(c, http.StatusServiceUnavailable, "Login Disabled", detail, "login_disabled")
		}

		if err := c.Redirect(http.StatusFound, authURL); err != nil {
			return fmt.Errorf("redirect to provider: %w", err)
		}
		return nil
	}
}

// CallbackHandler handles the OAuth callback from Entra ID.
func CallbackHandler(svc *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if user := svc.TestUser(); user != nil {
			return setUserCookieAndRedirect(svc, c, user)
		}

		state := c.QueryParam("state")
		code := c.QueryParam("code")
		if state == "" || code == "" {
			return jsonAPIError(c, http.StatusBadRequest, "Invalid Request", "Missing state or code", "invalid_request")
		}

		stored, err := c.Cookie(stateCookieName)
		if err != nil {
			return jsonAPIError(c, http.StatusBadRequest, "Invalid Request", "Missing login state", "missing_state")
		}

		decodedState, err := svc.DecodeState(stored.Value)
		if err != nil || decodedState != state {
			return jsonAPIError(c, http.StatusBadRequest, "Invalid Request", "Invalid login state", "invalid_state")
		}

		token, err := svc.Exchange(c.Request().Context(), code)
		if err != nil {
			return jsonAPIError(c, http.StatusBadRequest, "Login Failed", "Token exchange failed", "token_exchange")
		}

		user, err := svc.FetchUser(c.Request().Context(), token)
		if err != nil {
			return jsonAPIError(c, http.StatusBadRequest, "Login Failed", "User lookup failed", "user_lookup")
		}

		return setUserCookieAndRedirect(svc, c, user)
	}
}

func jsonAPIError(c echo.Context, status int, title, detail, code string) error {
	errResp := api.NewError(status, title, detail, code)
	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	if err := c.JSON(status, errResp); err != nil {
		return fmt.Errorf("write error response: %w", err)
	}
	return nil
}

func newCookie(name, value string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func setUserCookieAndRedirect(svc *Service, c echo.Context, user *User) error {
	encodedUser, err := svc.EncodeUser(*user)
	if err != nil {
		return jsonAPIError(c, http.StatusInternalServerError, "Server Error", "Failed to store user", "user_store")
	}
	userCookie := newCookie(userCookieName, encodedUser, c.Scheme() == "https")
	c.SetCookie(userCookie)
	if err := c.Redirect(http.StatusFound, "/"); err != nil {
		return fmt.Errorf("redirect after login: %w", err)
	}
	return nil
}
