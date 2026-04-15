package auth

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"log/slog"
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
	avatarPath := filepath.Join(avatarsDir, userID+".png")
	logFailure := func(message string, err error, extra ...any) {
		args := []any{"user_id", userID}
		if err != nil {
			args = append(args, "error", err)
		}
		args = append(args, extra...)
		slog.Error(message, args...)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://graph.microsoft.com/v1.0/me/photo/$value", http.NoBody)
	if err != nil {
		logFailure("build avatar sync request", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		logFailure("download avatar", err)
		return
	}
	defer resp.Body.Close() //nolint:errcheck // Best-effort cleanup

	if resp.StatusCode == http.StatusNotFound {
		slog.Info("avatar not found in Microsoft Graph", "user_id", userID)
		if err := os.Remove(avatarPath); err != nil && !os.IsNotExist(err) {
			logFailure("remove stale avatar after not-found response", err, "avatar_path", avatarPath)
			return
		}
		return
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("unexpected avatar sync status", "user_id", userID, "status_code", resp.StatusCode)
		return
	}

	// Read the full body into a buffer (up to maxAvatarSize) to avoid
	// LimitReader truncating mid-decode, which caused "not enough pixel data" errors.
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, maxAvatarSize+1))
	if err != nil {
		logFailure("read avatar body", err)
		return
	}
	contentType := resp.Header.Get("Content-Type")
	if int64(len(bodyBytes)) > maxAvatarSize {
		slog.Warn("avatar exceeds size limit, skipping",
			"user_id", userID, "content_type", contentType, "bytes", len(bodyBytes))
		return
	}

	img, _, err := image.Decode(bytes.NewReader(bodyBytes))
	if err != nil {
		logFailure("decode avatar image", err,
			"content_type", contentType, "bytes", len(bodyBytes))
		return
	}

	// #nosec G304 -- path is constructed from trusted avatarsDir + user ID
	out, err := os.Create(avatarPath)
	if err != nil {
		logFailure("create avatar file", err, "avatar_path", avatarPath)
		return
	}
	defer out.Close() //nolint:errcheck // Best-effort cleanup

	if err := png.Encode(out, img); err != nil {
		logFailure("encode avatar PNG", err, "avatar_path", avatarPath)
		os.Remove(avatarPath) //nolint:errcheck // Clean up on failure
	}
}

// HTTPClient is the interface needed for avatar sync HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
