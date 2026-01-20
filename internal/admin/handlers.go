// Package admin provides admin API handlers for managing spaces.
package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// AreaRequest represents an area create/update request.
type AreaRequest struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id,omitempty"`
		Attributes struct {
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
			FloorPlan   string `json:"floor_plan,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// RoomRequest represents a room create/update request.
type RoomRequest struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id,omitempty"`
		Attributes struct {
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
			FloorPlan   string `json:"floor_plan,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// DeskRequest represents a desk create/update request.
type DeskRequest struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id,omitempty"`
		Attributes struct {
			Name      string   `json:"name"`
			Equipment []string `json:"equipment,omitempty"`
			Warning   string   `json:"warning,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// Helper functions to reduce duplication

func checkAdmin(c echo.Context) error {
	if !isAdmin(c) {
		return api.WriteForbidden(c) //nolint:wrapcheck // Response helper
	}
	return nil
}

func requireParam(c echo.Context, name string) (string, error) {
	val := c.Param(name)
	if val == "" {
		return "", api.WriteBadRequest(c, name+" is required") //nolint:wrapcheck // Response helper
	}
	return val, nil
}

func parseJSON[T any](c echo.Context) (*T, error) {
	var req T
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return nil, api.WriteBadRequest(c, "Invalid request body") //nolint:wrapcheck // Response helper
	}
	return &req, nil
}

func requireName(name string, c echo.Context) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", api.WriteBadRequest(c, "name is required") //nolint:wrapcheck // Response helper
	}
	return trimmed, nil
}

func generateID(provided string) string {
	if provided != "" {
		// Sanitize: allow only alphanumeric, hyphens, and underscores
		sanitized := strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '-' || r == '_' {
				return r
			}
			return -1
		}, provided)
		if sanitized != "" {
			return sanitized
		}
	}
	return uuid.New().String()
}

func reloadConfig(ctx context.Context, configHolder *ConfigHolder, store *spaces.Store) error {
	if err := configHolder.Reload(ctx, store); err != nil {
		return fmt.Errorf("reload config: %w", err)
	}
	return nil
}

func handleDeleteError(err error, c echo.Context, entityName string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return api.WriteNotFound(c, entityName+" not found") //nolint:wrapcheck // Response helper
	}
	return fmt.Errorf("delete %s: %w", strings.ToLower(entityName), err)
}

func writeJSONAPIResource(c echo.Context, status int, resource api.Resource) error {
	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	return c.JSON(status, api.SingleResponse{Data: resource}) //nolint:wrapcheck // Echo response
}

func writeJSONAPICollection(c echo.Context, resources []api.Resource) error {
	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	return c.JSON(http.StatusOK, api.CollectionResponse{Data: resources}) //nolint:wrapcheck // Echo response
}

// --- Area Handlers ---

// ListAreasHandler returns all areas (admin view with full details).
func ListAreasHandler(store *spaces.Store) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		areas, err := store.ListAreas(c.Request().Context())
		if err != nil {
			return fmt.Errorf("list areas: %w", err)
		}

		resources := make([]api.Resource, 0, len(areas))
		for i := range areas {
			resources = append(resources, areaToResource(&areas[i]))
		}
		return writeJSONAPICollection(c, resources)
	}
}

func areaToResource(a *spaces.AreaRecord) api.Resource {
	return api.Resource{
		Type: "areas",
		ID:   a.ID,
		Attributes: map[string]interface{}{
			"name":        a.Name,
			"description": a.Description,
			"floor_plan":  a.FloorPlan,
			"created_at":  a.CreatedAt,
			"updated_at":  a.UpdatedAt,
		},
	}
}

// CreateAreaHandler creates a new area.
func CreateAreaHandler(store *spaces.Store, configHolder *ConfigHolder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		req, err := parseJSON[AreaRequest](c)
		if err != nil {
			return err
		}

		if req.Data.Type != "areas" {
			return api.WriteBadRequest(c, "Resource type must be 'areas'")
		}

		name, err := requireName(req.Data.Attributes.Name, c)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		area := &spaces.AreaRecord{
			ID:          generateID(req.Data.ID),
			Name:        name,
			Description: req.Data.Attributes.Description,
			FloorPlan:   req.Data.Attributes.FloorPlan,
		}

		if err := store.CreateArea(ctx, area); err != nil {
			return fmt.Errorf("create area: %w", err)
		}

		if err := reloadConfig(ctx, configHolder, store); err != nil {
			return err
		}

		return writeJSONAPIResource(c, http.StatusCreated, areaToResource(area))
	}
}

// UpdateAreaHandler updates an existing area.
func UpdateAreaHandler(store *spaces.Store, configHolder *ConfigHolder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		areaID, err := requireParam(c, "area_id")
		if err != nil {
			return err
		}

		req, err := parseJSON[AreaRequest](c)
		if err != nil {
			return err
		}

		if req.Data.Type != "" && req.Data.Type != "areas" {
			return api.WriteBadRequest(c, "Resource type must be 'areas'")
		}

		name, err := requireName(req.Data.Attributes.Name, c)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		area := &spaces.AreaRecord{
			ID:          areaID,
			Name:        name,
			Description: req.Data.Attributes.Description,
			FloorPlan:   req.Data.Attributes.FloorPlan,
		}

		if err := store.UpdateArea(ctx, area); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return api.WriteNotFound(c, "Area not found")
			}
			return fmt.Errorf("update area: %w", err)
		}

		if err := reloadConfig(ctx, configHolder, store); err != nil {
			return err
		}

		return writeJSONAPIResource(c, http.StatusOK, areaToResource(area))
	}
}

// DeleteAreaHandler deletes an area.
//
//nolint:dupl // Delete handlers share structure but operate on different entities
func DeleteAreaHandler(store *spaces.Store, configHolder *ConfigHolder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		areaID, err := requireParam(c, "area_id")
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		if err := store.DeleteArea(ctx, areaID); err != nil {
			return handleDeleteError(err, c, "Area")
		}

		if err := reloadConfig(ctx, configHolder, store); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// --- Room Handlers ---

// ListRoomsHandler returns all rooms in an area (admin view).
func ListRoomsHandler(store *spaces.Store) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		areaID, err := requireParam(c, "area_id")
		if err != nil {
			return err
		}

		rooms, err := store.ListRooms(c.Request().Context(), areaID)
		if err != nil {
			return fmt.Errorf("list rooms: %w", err)
		}

		resources := make([]api.Resource, 0, len(rooms))
		for i := range rooms {
			resources = append(resources, roomToResource(&rooms[i]))
		}
		return writeJSONAPICollection(c, resources)
	}
}

func roomToResource(r *spaces.RoomRecord) api.Resource {
	return api.Resource{
		Type: "rooms",
		ID:   r.ID,
		Attributes: map[string]interface{}{
			"name":        r.Name,
			"description": r.Description,
			"floor_plan":  r.FloorPlan,
			"area_id":     r.AreaID,
			"created_at":  r.CreatedAt,
			"updated_at":  r.UpdatedAt,
		},
	}
}

// CreateRoomHandler creates a new room.
//
//nolint:dupl // CRUD handlers share structure but operate on different entities
func CreateRoomHandler(store *spaces.Store, configHolder *ConfigHolder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		areaID, err := requireParam(c, "area_id")
		if err != nil {
			return err
		}

		req, err := parseJSON[RoomRequest](c)
		if err != nil {
			return err
		}

		if req.Data.Type != "rooms" {
			return api.WriteBadRequest(c, "Resource type must be 'rooms'")
		}

		name, err := requireName(req.Data.Attributes.Name, c)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()

		// Check area exists
		area, err := store.GetArea(ctx, areaID)
		if err != nil {
			return fmt.Errorf("get area: %w", err)
		}
		if area == nil {
			return api.WriteNotFound(c, "Area not found")
		}

		room := &spaces.RoomRecord{
			ID:          generateID(req.Data.ID),
			AreaID:      areaID,
			Name:        name,
			Description: req.Data.Attributes.Description,
			FloorPlan:   req.Data.Attributes.FloorPlan,
		}

		if err := store.CreateRoom(ctx, room); err != nil {
			return fmt.Errorf("create room: %w", err)
		}

		if err := reloadConfig(ctx, configHolder, store); err != nil {
			return err
		}

		return writeJSONAPIResource(c, http.StatusCreated, roomToResource(room))
	}
}

// UpdateRoomHandler updates an existing room.
func UpdateRoomHandler(store *spaces.Store, configHolder *ConfigHolder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		roomID, err := requireParam(c, "room_id")
		if err != nil {
			return err
		}

		req, err := parseJSON[RoomRequest](c)
		if err != nil {
			return err
		}

		if req.Data.Type != "" && req.Data.Type != "rooms" {
			return api.WriteBadRequest(c, "Resource type must be 'rooms'")
		}

		name, err := requireName(req.Data.Attributes.Name, c)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()

		// Get existing room to preserve area_id
		existing, err := store.GetRoom(ctx, roomID)
		if err != nil {
			return fmt.Errorf("get room: %w", err)
		}
		if existing == nil {
			return api.WriteNotFound(c, "Room not found")
		}

		room := &spaces.RoomRecord{
			ID:          roomID,
			AreaID:      existing.AreaID,
			Name:        name,
			Description: req.Data.Attributes.Description,
			FloorPlan:   req.Data.Attributes.FloorPlan,
		}

		if err := store.UpdateRoom(ctx, room); err != nil {
			return fmt.Errorf("update room: %w", err)
		}

		if err := reloadConfig(ctx, configHolder, store); err != nil {
			return err
		}

		return writeJSONAPIResource(c, http.StatusOK, roomToResource(room))
	}
}

// DeleteRoomHandler deletes a room.
//
//nolint:dupl // Delete handlers share structure but operate on different entities
func DeleteRoomHandler(store *spaces.Store, configHolder *ConfigHolder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		roomID, err := requireParam(c, "room_id")
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		if err := store.DeleteRoom(ctx, roomID); err != nil {
			return handleDeleteError(err, c, "Room")
		}

		if err := reloadConfig(ctx, configHolder, store); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// --- Desk Handlers ---

// ListDesksHandler returns all desks in a room (admin view).
func ListDesksHandler(store *spaces.Store) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		roomID, err := requireParam(c, "room_id")
		if err != nil {
			return err
		}

		desks, err := store.ListDesks(c.Request().Context(), roomID)
		if err != nil {
			return fmt.Errorf("list desks: %w", err)
		}

		resources := make([]api.Resource, 0, len(desks))
		for i := range desks {
			resources = append(resources, deskToResource(&desks[i]))
		}
		return writeJSONAPICollection(c, resources)
	}
}

func deskToResource(d *spaces.DeskRecord) api.Resource {
	return api.Resource{
		Type: "desks",
		ID:   d.ID,
		Attributes: map[string]interface{}{
			"name":       d.Name,
			"equipment":  d.Equipment,
			"warning":    d.Warning,
			"room_id":    d.RoomID,
			"created_at": d.CreatedAt,
			"updated_at": d.UpdatedAt,
		},
	}
}

// CreateDeskHandler creates a new desk.
//
//nolint:dupl // CRUD handlers share structure but operate on different entities
func CreateDeskHandler(store *spaces.Store, configHolder *ConfigHolder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		roomID, err := requireParam(c, "room_id")
		if err != nil {
			return err
		}

		req, err := parseJSON[DeskRequest](c)
		if err != nil {
			return err
		}

		if req.Data.Type != "desks" {
			return api.WriteBadRequest(c, "Resource type must be 'desks'")
		}

		name, err := requireName(req.Data.Attributes.Name, c)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()

		// Check room exists
		room, err := store.GetRoom(ctx, roomID)
		if err != nil {
			return fmt.Errorf("get room: %w", err)
		}
		if room == nil {
			return api.WriteNotFound(c, "Room not found")
		}

		desk := &spaces.DeskRecord{
			ID:        generateID(req.Data.ID),
			RoomID:    roomID,
			Name:      name,
			Equipment: req.Data.Attributes.Equipment,
			Warning:   req.Data.Attributes.Warning,
		}

		if err := store.CreateDesk(ctx, desk); err != nil {
			return fmt.Errorf("create desk: %w", err)
		}

		if err := reloadConfig(ctx, configHolder, store); err != nil {
			return err
		}

		return writeJSONAPIResource(c, http.StatusCreated, deskToResource(desk))
	}
}

// UpdateDeskHandler updates an existing desk.
func UpdateDeskHandler(store *spaces.Store, configHolder *ConfigHolder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		deskID, err := requireParam(c, "desk_id")
		if err != nil {
			return err
		}

		req, err := parseJSON[DeskRequest](c)
		if err != nil {
			return err
		}

		if req.Data.Type != "" && req.Data.Type != "desks" {
			return api.WriteBadRequest(c, "Resource type must be 'desks'")
		}

		name, err := requireName(req.Data.Attributes.Name, c)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()

		// Get existing desk to preserve room_id
		existing, err := store.GetDesk(ctx, deskID)
		if err != nil {
			return fmt.Errorf("get desk: %w", err)
		}
		if existing == nil {
			return api.WriteNotFound(c, "Desk not found")
		}

		desk := &spaces.DeskRecord{
			ID:        deskID,
			RoomID:    existing.RoomID,
			Name:      name,
			Equipment: req.Data.Attributes.Equipment,
			Warning:   req.Data.Attributes.Warning,
		}

		if err := store.UpdateDesk(ctx, desk); err != nil {
			return fmt.Errorf("update desk: %w", err)
		}

		if err := reloadConfig(ctx, configHolder, store); err != nil {
			return err
		}

		return writeJSONAPIResource(c, http.StatusOK, deskToResource(desk))
	}
}

// DeleteDeskHandler deletes a desk.
//
//nolint:dupl // Delete handlers share structure but operate on different entities
func DeleteDeskHandler(store *spaces.Store, configHolder *ConfigHolder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := checkAdmin(c); err != nil {
			return err
		}

		deskID, err := requireParam(c, "desk_id")
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		if err := store.DeleteDesk(ctx, deskID); err != nil {
			return handleDeleteError(err, c, "Desk")
		}

		if err := reloadConfig(ctx, configHolder, store); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func isAdmin(c echo.Context) bool {
	user := auth.GetUserFromContext(c)
	return user != nil && user.IsAdmin
}
