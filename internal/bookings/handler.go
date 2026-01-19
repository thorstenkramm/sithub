// Package bookings provides booking handlers and store helpers.
package bookings

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mattn/go-sqlite3"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// CreateRequest represents a booking create JSON:API payload.
type CreateRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			DeskID      string `json:"desk_id"`
			BookingDate string `json:"booking_date"`
		} `json:"attributes"`
	} `json:"data"`
}

// BookingAttributes represents booking resource attributes.
type BookingAttributes struct {
	DeskID      string `json:"desk_id"`
	UserID      string `json:"user_id"`
	BookingDate string `json:"booking_date"`
	CreatedAt   string `json:"created_at"`
}

// MyBookingAttributes represents booking resource attributes with location info.
type MyBookingAttributes struct {
	DeskID      string `json:"desk_id"`
	DeskName    string `json:"desk_name"`
	RoomID      string `json:"room_id"`
	RoomName    string `json:"room_name"`
	AreaID      string `json:"area_id"`
	AreaName    string `json:"area_name"`
	BookingDate string `json:"booking_date"`
	CreatedAt   string `json:"created_at"`
}

// DeleteHandler returns a handler for canceling a user's own booking.
func DeleteHandler(store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := auth.GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		bookingID := c.Param("id")
		if bookingID == "" {
			return api.WriteBadRequest(c, "Booking ID is required")
		}

		ctx := c.Request().Context()

		// Check if booking exists and belongs to user
		booking, err := FindBookingByID(ctx, store, bookingID)
		if err != nil {
			return fmt.Errorf("find booking: %w", err)
		}
		if booking == nil {
			return api.WriteNotFound(c, "Booking not found")
		}
		if booking.UserID != user.ID {
			return api.WriteNotFound(c, "Booking not found")
		}

		// Delete the booking
		if err := DeleteBooking(ctx, store, bookingID); err != nil {
			return fmt.Errorf("delete booking: %w", err)
		}

		slog.Info("booking canceled",
			"booking_id", bookingID,
			"user_id", user.ID,
			"desk_id", booking.DeskID,
			"booking_date", booking.BookingDate,
		)

		return c.NoContent(http.StatusNoContent)
	}
}

// ListHandler returns a handler for listing the current user's future bookings.
func ListHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := auth.GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		today := time.Now().UTC().Format(time.DateOnly)
		ctx := c.Request().Context()

		records, err := ListUserBookings(ctx, store, user.ID, today)
		if err != nil {
			return fmt.Errorf("list user bookings: %w", err)
		}

		resources := make([]api.Resource, 0, len(records))
		for _, rec := range records {
			loc, found := cfg.FindDeskLocation(rec.DeskID)
			if !found {
				slog.Warn("booking references unknown desk", "booking_id", rec.ID, "desk_id", rec.DeskID)
				continue
			}

			resources = append(resources, api.Resource{
				Type: "bookings",
				ID:   rec.ID,
				Attributes: MyBookingAttributes{
					DeskID:      rec.DeskID,
					DeskName:    loc.Desk.Name,
					RoomID:      loc.Room.ID,
					RoomName:    loc.Room.Name,
					AreaID:      loc.Area.ID,
					AreaName:    loc.Area.Name,
					BookingDate: rec.BookingDate,
					CreatedAt:   rec.CreatedAt,
				},
			})
		}

		resp := api.CollectionResponse{Data: resources}
		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		return c.JSON(http.StatusOK, resp)
	}
}

// CreateHandler returns a handler for creating a single-day booking.
func CreateHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := auth.GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		if err := validateContentType(c); err != nil {
			if errors.Is(err, errResponseWritten) {
				return nil
			}
			return err
		}

		req, err := parseCreateRequest(c)
		if err != nil {
			var valErr validationError
			if errors.As(err, &valErr) {
				return api.WriteBadRequest(c, valErr.detail)
			}
			return err
		}

		deskID, bookingDate, err := validateRequestFields(req)
		if err != nil {
			var valErr validationError
			if errors.As(err, &valErr) {
				return api.WriteBadRequest(c, valErr.detail)
			}
			return err
		}

		if _, exists := cfg.FindDesk(deskID); !exists {
			return api.WriteNotFound(c, "Desk not found")
		}

		return processBooking(c, store, deskID, user.ID, bookingDate)
	}
}

func validateContentType(c echo.Context) error {
	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if !strings.Contains(contentType, api.JSONAPIContentType) {
		//nolint:errcheck // Error is intentionally ignored; response signals via errResponseWritten
		api.WriteUnsupportedMediaType(c, "Content-Type must be application/vnd.api+json")
		return errResponseWritten
	}
	return nil
}

func parseCreateRequest(c echo.Context) (*CreateRequest, error) {
	var req CreateRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return nil, errBadRequest("Invalid request body")
	}
	if req.Data.Type != "bookings" {
		return nil, errBadRequest("Resource type must be 'bookings'")
	}
	return &req, nil
}

func validateRequestFields(req *CreateRequest) (deskID, bookingDate string, err error) {
	deskID = strings.TrimSpace(req.Data.Attributes.DeskID)
	bookingDate = strings.TrimSpace(req.Data.Attributes.BookingDate)

	if deskID == "" {
		return "", "", errBadRequest("desk_id is required")
	}
	if bookingDate == "" {
		return "", "", errBadRequest("booking_date is required")
	}

	parsedDate, parseErr := time.Parse(time.DateOnly, bookingDate)
	if parseErr != nil {
		return "", "", errBadRequest("booking_date must be in YYYY-MM-DD format")
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	if parsedDate.Before(today) {
		return "", "", errBadRequest("booking_date cannot be in the past")
	}

	return deskID, bookingDate, nil
}

// validationError is a sentinel for validation errors that need WriteBadRequest.
type validationError struct {
	detail string
}

func (e validationError) Error() string { return e.detail }

func errBadRequest(detail string) error {
	return validationError{detail: detail}
}

// errResponseWritten indicates the HTTP response was already written.
var errResponseWritten = errors.New("response already written")

func processBooking(c echo.Context, store *sql.DB, deskID, userID, bookingDate string) error {
	ctx := c.Request().Context()

	existingBookingID, err := FindUserBooking(ctx, store, deskID, userID, bookingDate)
	if err != nil {
		return fmt.Errorf("check existing booking: %w", err)
	}
	if existingBookingID != "" {
		//nolint:wrapcheck // Terminal response, no wrapping needed
		return api.WriteConflict(c, "You already have this desk booked for this date")
	}

	booking, err := CreateBooking(ctx, store, deskID, userID, bookingDate)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			slog.Warn("booking conflict",
				"desk_id", deskID,
				"user_id", userID,
				"booking_date", bookingDate,
			)
			//nolint:wrapcheck // Terminal response, no wrapping needed
			return api.WriteConflict(c, "Desk is already booked for this date")
		}
		return fmt.Errorf("create booking: %w", err)
	}

	slog.Info("booking created",
		"booking_id", booking.ID,
		"desk_id", deskID,
		"user_id", userID,
		"booking_date", bookingDate,
	)

	return writeBookingResponse(c, booking)
}

func writeBookingResponse(c echo.Context, booking *Booking) error {
	resp := api.SingleResponse{
		Data: api.Resource{
			Type: "bookings",
			ID:   booking.ID,
			Attributes: BookingAttributes{
				DeskID:      booking.DeskID,
				UserID:      booking.UserID,
				BookingDate: booking.BookingDate,
				CreatedAt:   booking.CreatedAt,
			},
		},
	}

	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	//nolint:wrapcheck // Terminal response, no wrapping needed
	return c.JSON(http.StatusCreated, resp)
}

// Booking represents a booking record.
type Booking struct {
	ID          string
	DeskID      string
	UserID      string
	BookingDate string
	CreatedAt   string
	UpdatedAt   string
}

// ErrConflict indicates a booking conflict (desk already booked).
var ErrConflict = errors.New("booking conflict")

// FindUserBooking checks if a user already has a booking for a specific desk and date.
// Returns the booking ID if found, empty string otherwise.
func FindUserBooking(ctx context.Context, store *sql.DB, deskID, userID, bookingDate string) (string, error) {
	var bookingID string
	err := store.QueryRowContext(ctx,
		"SELECT id FROM bookings WHERE desk_id = ? AND user_id = ? AND booking_date = ?",
		deskID, userID, bookingDate,
	).Scan(&bookingID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("query user booking: %w", err)
	}
	return bookingID, nil
}

// CreateBooking inserts a new booking record.
func CreateBooking(ctx context.Context, store *sql.DB, deskID, userID, bookingDate string) (*Booking, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.New().String()

	_, err := store.ExecContext(ctx,
		`INSERT INTO bookings (id, desk_id, user_id, booking_date, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		id, deskID, userID, bookingDate, now, now,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("insert booking: %w", err)
	}

	return &Booking{
		ID:          id,
		DeskID:      deskID,
		UserID:      userID,
		BookingDate: bookingDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
