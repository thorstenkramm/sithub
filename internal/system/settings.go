// Package system provides system health and settings endpoints.
package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// SettingsAttributes contains public application settings.
type SettingsAttributes struct {
	WeeksInAdvanced int `json:"weeks_in_advanced"`
}

// SettingsHandler returns a handler that exposes public application settings.
func SettingsHandler(weeksInAdvanced int) echo.HandlerFunc {
	return func(c echo.Context) error {
		return api.WriteSingle(c, http.StatusOK, api.Resource{
			Type: "settings",
			ID:   "settings",
			Attributes: SettingsAttributes{
				WeeksInAdvanced: weeksInAdvanced,
			},
		}, "write settings response")
	}
}
