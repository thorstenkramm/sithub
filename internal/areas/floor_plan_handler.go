package areas

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// mimeTypes maps file extensions to MIME types for floor plan images.
var mimeTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".svg":  "image/svg+xml",
}

// FloorPlanHandler serves floor plan images from a configured directory.
// Returns 404 if floor plans are not configured or the file does not exist.
func FloorPlanHandler(floorPlansDir string) echo.HandlerFunc {
	return func(c echo.Context) error {
		if floorPlansDir == "" {
			return api.WriteNotFound(c, "Floor plans not configured")
		}

		filename := c.Param("filename")
		if filename == "" || strings.ContainsAny(filename, "/\\") {
			return api.WriteNotFound(c, "Floor plan not found")
		}

		ext := strings.ToLower(filepath.Ext(filename))
		contentType, ok := mimeTypes[ext]
		if !ok {
			return api.WriteNotFound(c, "Floor plan not found")
		}

		fullPath := filepath.Join(floorPlansDir, filename)

		data, err := os.ReadFile(fullPath) // #nosec G304 -- filename validated above; floorPlansDir from config.
		if err != nil {
			if os.IsNotExist(err) {
				return api.WriteNotFound(c, "Floor plan not found")
			}
			return fmt.Errorf("read floor plan: %w", err)
		}

		return c.Blob(http.StatusOK, contentType, data)
	}
}
