package auth

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	// Register JPEG decoder for image.Decode.
	_ "image/jpeg"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

const maxAvatarSize = 512 * 1024 // 512 KB

// ServeAvatarHandler serves user avatar images.
// GET /api/v1/avatars/:user_id
func ServeAvatarHandler(avatarsDir string) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("user_id")
		if userID == "" {
			return api.WriteBadRequest(c, "user_id is required")
		}

		avatarPath := filepath.Join(avatarsDir, userID+".png")

		// Prevent path traversal
		if filepath.Dir(avatarPath) != avatarsDir {
			return api.WriteNotFound(c, "Avatar not found")
		}

		if _, err := os.Stat(avatarPath); os.IsNotExist(err) {
			return api.WriteNotFound(c, "Avatar not found")
		}

		c.Response().Header().Set("Cache-Control", "max-age=300")
		return c.File(avatarPath)
	}
}

// UploadAvatarHandler handles avatar upload for the current user.
// POST /api/v1/me/avatar (multipart form, field "avatar")
func UploadAvatarHandler(avatarsDir string) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		file, err := c.FormFile("avatar")
		if err != nil {
			return api.WriteBadRequest(c, "avatar file is required")
		}

		if file.Size > maxAvatarSize {
			return api.WriteBadRequest(c,
				fmt.Sprintf("avatar must be smaller than %d KB", maxAvatarSize/1024))
		}

		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("open uploaded file: %w", err)
		}
		defer src.Close() //nolint:errcheck // Best-effort cleanup

		// Decode as image (supports PNG and JPEG)
		img, _, err := image.Decode(src)
		if err != nil {
			return api.WriteBadRequest(c, "invalid image format (PNG or JPEG required)")
		}

		// Save as PNG
		avatarPath := filepath.Join(avatarsDir, user.ID+".png")
		// #nosec G304 -- path is constructed from trusted avatarsDir + user ID
		out, err := os.Create(avatarPath)
		if err != nil {
			return fmt.Errorf("create avatar file: %w", err)
		}
		defer out.Close() //nolint:errcheck // Best-effort cleanup

		if err := png.Encode(out, img); err != nil {
			return fmt.Errorf("encode avatar PNG: %w", err)
		}

		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		return c.JSON(http.StatusOK, api.SingleResponse{
			Data: api.Resource{
				Type: "avatars",
				ID:   user.ID,
				Attributes: map[string]string{
					"status": "uploaded",
				},
			},
		})
	}
}

// DeleteAvatarHandler removes the current user's avatar.
// DELETE /api/v1/me/avatar
func DeleteAvatarHandler(avatarsDir string) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		avatarPath := filepath.Join(avatarsDir, user.ID+".png")
		if err := os.Remove(avatarPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove avatar: %w", err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// SyncAvatar downloads the user's profile photo from Microsoft Graph
// and saves it as a PNG. Errors are logged but not propagated — avatar
// sync must never block login.
func SyncAvatar(ctx context.Context, client HTTPClient, userID, avatarsDir string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://graph.microsoft.com/v1.0/me/photo/$value", http.NoBody)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close() //nolint:errcheck // Best-effort cleanup

	if resp.StatusCode != http.StatusOK {
		return
	}

	// Read up to maxAvatarSize + 1 byte to detect oversized photos
	lr := io.LimitReader(resp.Body, maxAvatarSize+1)
	img, _, err := image.Decode(lr)
	if err != nil {
		return
	}

	avatarPath := filepath.Join(avatarsDir, userID+".png")
	// #nosec G304 -- path is constructed from trusted avatarsDir + user ID
	out, err := os.Create(avatarPath)
	if err != nil {
		return
	}
	defer out.Close() //nolint:errcheck // Best-effort cleanup

	if err := png.Encode(out, img); err != nil {
		os.Remove(avatarPath) //nolint:errcheck // Clean up on failure
	}
}

// HTTPClient is the interface needed for avatar sync HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
