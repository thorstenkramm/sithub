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
	"github.com/thorstenkramm/sithub/internal/notifications"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// CreateRequest represents a booking create JSON:API payload.
type CreateRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			DeskID       string   `json:"desk_id"`
			BookingDate  string   `json:"booking_date"`
			BookingDates []string `json:"booking_dates,omitempty"`
			ForUserID    string   `json:"for_user_id,omitempty"`
			ForUserName  string   `json:"for_user_name,omitempty"`
			IsGuest      bool     `json:"is_guest,omitempty"`
			GuestEmail   string   `json:"guest_email,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// BookingAttributes represents booking resource attributes.
type BookingAttributes struct {
	DeskID           string `json:"desk_id"`
	UserID           string `json:"user_id"`
	BookingDate      string `json:"booking_date"`
	CreatedAt        string `json:"created_at"`
	BookedByUserID   string `json:"booked_by_user_id,omitempty"`
	BookedByUserName string `json:"booked_by_user_name,omitempty"`
	IsGuest          bool   `json:"is_guest,omitempty"`
	GuestEmail       string `json:"guest_email,omitempty"`
}

// MultiDayBookingResult represents the result of a multi-day booking request.
type MultiDayBookingResult struct {
	Created   []api.Resource `json:"created"`
	Conflicts []string       `json:"conflicts,omitempty"`
}

// MyBookingAttributes represents booking resource attributes with location info.
type MyBookingAttributes struct {
	DeskID           string `json:"desk_id"`
	DeskName         string `json:"desk_name"`
	RoomID           string `json:"room_id"`
	RoomName         string `json:"room_name"`
	AreaID           string `json:"area_id"`
	AreaName         string `json:"area_name"`
	BookingDate      string `json:"booking_date"`
	CreatedAt        string `json:"created_at"`
	BookedByUserID   string `json:"booked_by_user_id,omitempty"`
	BookedByUserName string `json:"booked_by_user_name,omitempty"`
	BookedForMe      bool   `json:"booked_for_me,omitempty"`
	IsGuest          bool   `json:"is_guest,omitempty"`
	GuestEmail       string `json:"guest_email,omitempty"`
}

// DeleteHandler returns a handler for canceling a booking.
// Users can cancel their own bookings or bookings made for them;
// The person who booked on behalf can also cancel; admins can cancel any booking.
func DeleteHandler(store *sql.DB, notifier notifications.Notifier) echo.HandlerFunc {
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

		// Check if booking exists
		booking, err := FindBookingByID(ctx, store, bookingID)
		if err != nil {
			return fmt.Errorf("find booking: %w", err)
		}
		if booking == nil {
			return api.WriteNotFound(c, "Booking not found")
		}

		// Check authorization: owner, booker, or admin
		isOwner := booking.UserID == user.ID
		isBooker := booking.BookedByUserID == user.ID
		canCancel := isOwner || isBooker || user.IsAdmin
		if !canCancel {
			return api.WriteNotFound(c, "Booking not found")
		}

		// Delete the booking
		if err := DeleteBooking(ctx, store, bookingID); err != nil {
			return fmt.Errorf("delete booking: %w", err)
		}

		logFields := []any{
			"booking_id", bookingID,
			"canceled_by", user.ID,
			"desk_id", booking.DeskID,
			"booking_date", booking.BookingDate,
		}
		if !isOwner {
			logFields = append(logFields, "booking_owner", booking.UserID)
			if user.IsAdmin && !isBooker {
				logFields = append(logFields, "admin_action", true)
			}
		}
		slog.Info("booking canceled", logFields...)

		// Send notification asynchronously
		notifier.NotifyAsync(&notifications.BookingEvent{
			Event:            notifications.EventBookingCanceled,
			BookingID:        bookingID,
			DeskID:           booking.DeskID,
			UserID:           booking.UserID,
			UserName:         booking.UserName,
			BookingDate:      booking.BookingDate,
			IsGuest:          booking.IsGuest,
			GuestEmail:       booking.GuestEmail,
			CanceledByUserID: user.ID,
			Timestamp:        time.Now().UTC().Format(time.RFC3339),
		})

		return c.NoContent(http.StatusNoContent)
	}
}

// ListHandler returns a handler for listing the current user's future bookings.
// Includes bookings made by the user AND bookings made for the user by others.
func ListHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return ListHandlerDynamic(func() *spaces.Config { return cfg }, store)
}

// ListHandlerDynamic returns a handler for listing the current user's future bookings.
func ListHandlerDynamic(getConfig spaces.ConfigGetter, store *sql.DB) echo.HandlerFunc {
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

		return writeBookingsCollection(c, getConfig(), user.ID, records)
	}
}

// HistoryHandler returns a handler for listing the current user's past bookings.
// Accepts optional query params: from (start date), to (end date) in YYYY-MM-DD format.
func HistoryHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return HistoryHandlerDynamic(func() *spaces.Config { return cfg }, store)
}

// HistoryHandlerDynamic returns a handler for listing the current user's past bookings.
func HistoryHandlerDynamic(getConfig spaces.ConfigGetter, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := auth.GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		// Default: last 30 days
		today := time.Now().UTC()
		defaultFrom := today.AddDate(0, 0, -30).Format(time.DateOnly)
		defaultTo := today.AddDate(0, 0, -1).Format(time.DateOnly) // Yesterday (past only)

		fromDate := c.QueryParam("from")
		toDate := c.QueryParam("to")

		if fromDate == "" {
			fromDate = defaultFrom
		}
		if toDate == "" {
			toDate = defaultTo
		}

		// Validate dates
		if _, err := time.Parse(time.DateOnly, fromDate); err != nil {
			return api.WriteBadRequest(c, "Invalid 'from' date. Use YYYY-MM-DD format.")
		}
		if _, err := time.Parse(time.DateOnly, toDate); err != nil {
			return api.WriteBadRequest(c, "Invalid 'to' date. Use YYYY-MM-DD format.")
		}

		ctx := c.Request().Context()
		records, err := ListUserBookingsRange(ctx, store, user.ID, fromDate, toDate)
		if err != nil {
			return fmt.Errorf("list booking history: %w", err)
		}

		return writeBookingsCollection(c, getConfig(), user.ID, records)
	}
}

func writeBookingsCollection(
	c echo.Context, cfg *spaces.Config, currentUserID string, records []BookingRecord,
) error {
	resources := make([]api.Resource, 0, len(records))
	for i := range records {
		rec := &records[i]
		loc, found := cfg.FindDeskLocation(rec.DeskID)
		if !found {
			slog.Warn("booking references unknown desk", "booking_id", rec.ID, "desk_id", rec.DeskID)
			continue
		}

		attrs := MyBookingAttributes{
			DeskID:      rec.DeskID,
			DeskName:    loc.Desk.Name,
			RoomID:      loc.Room.ID,
			RoomName:    loc.Room.Name,
			AreaID:      loc.Area.ID,
			AreaName:    loc.Area.Name,
			BookingDate: rec.BookingDate,
			CreatedAt:   rec.CreatedAt,
		}

		// Include booked_by info if different from user_id
		if rec.BookedByUserID != "" && rec.BookedByUserID != rec.UserID {
			attrs.BookedByUserID = rec.BookedByUserID
			attrs.BookedByUserName = rec.BookedByUserName
		}
		// Mark if this booking was made for the current user by someone else
		if rec.UserID == currentUserID && rec.BookedByUserID != "" && rec.BookedByUserID != currentUserID {
			attrs.BookedForMe = true
		}
		// Include guest info
		if rec.IsGuest {
			attrs.IsGuest = true
			attrs.GuestEmail = rec.GuestEmail
		}

		resources = append(resources, api.Resource{
			Type:       "bookings",
			ID:         rec.ID,
			Attributes: attrs,
		})
	}

	resp := api.CollectionResponse{Data: resources}
	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	//nolint:wrapcheck // Terminal response
	return c.JSON(http.StatusOK, resp)
}

// CreateHandler returns a handler for creating bookings.
// Supports single-day, multi-day, booking on behalf, and guest bookings.
func CreateHandler(cfg *spaces.Config, store *sql.DB, notifier notifications.Notifier) echo.HandlerFunc {
	return CreateHandlerDynamic(func() *spaces.Config { return cfg }, store, notifier)
}

// CreateHandlerDynamic returns a handler for creating bookings using dynamic config.
func CreateHandlerDynamic(
	getConfig spaces.ConfigGetter, store *sql.DB, notifier notifications.Notifier,
) echo.HandlerFunc {
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
			return handleValidationError(c, err)
		}

		deskID, dates, err := validateRequestFieldsMultiDay(req)
		if err != nil {
			return handleValidationError(c, err)
		}

		cfg := getConfig()
		if _, exists := cfg.FindDesk(deskID); !exists {
			return api.WriteNotFound(c, "Desk not found")
		}

		params, err := resolveBookingParticipants(user, req)
		if err != nil {
			return handleValidationError(c, err)
		}

		// Single day booking
		if len(dates) == 1 {
			return processBooking(
				c, store, notifier, deskID, params.targetUserID, params.targetUserName,
				params.bookedByUserID, params.bookedByUserName, dates[0],
				params.isGuest, params.guestEmail,
			)
		}

		// Multi-day booking
		return processMultiDayBooking(c, store, notifier, deskID, params, dates)
	}
}

// bookingParticipants holds the resolved user info for a booking.
type bookingParticipants struct {
	targetUserID     string
	targetUserName   string
	bookedByUserID   string
	bookedByUserName string
	isGuest          bool
	guestEmail       string
}

// resolveBookingParticipants determines the target user and booker for a booking.
func resolveBookingParticipants(user *auth.User, req *CreateRequest) (*bookingParticipants, error) {
	params := &bookingParticipants{
		targetUserID:     user.ID,
		targetUserName:   user.Name,
		bookedByUserID:   user.ID,
		bookedByUserName: user.Name,
	}

	// Handle guest booking
	if req.Data.Attributes.IsGuest {
		guestName := strings.TrimSpace(req.Data.Attributes.ForUserName)
		guestEmail := strings.TrimSpace(req.Data.Attributes.GuestEmail)
		if guestName == "" {
			return nil, errBadRequest("for_user_name (guest name) is required for guest bookings")
		}
		// Generate a unique guest ID
		params.targetUserID = "guest-" + uuid.New().String()[:8]
		params.targetUserName = guestName
		params.isGuest = true
		params.guestEmail = guestEmail
		return params, nil
	}

	// Handle booking on behalf of another user
	forUserID := strings.TrimSpace(req.Data.Attributes.ForUserID)
	forUserName := strings.TrimSpace(req.Data.Attributes.ForUserName)
	if forUserID != "" {
		if forUserName == "" {
			return nil, errBadRequest("for_user_name is required when for_user_id is provided")
		}
		params.targetUserID = forUserID
		params.targetUserName = forUserName
	}

	return params, nil
}

func handleValidationError(c echo.Context, err error) error {
	var valErr validationError
	if errors.As(err, &valErr) {
		//nolint:wrapcheck // Terminal response, no wrapping needed
		return api.WriteBadRequest(c, valErr.detail)
	}
	return err
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

func validateRequestFieldsMultiDay(req *CreateRequest) (deskID string, dates []string, err error) {
	deskID = strings.TrimSpace(req.Data.Attributes.DeskID)
	if deskID == "" {
		return "", nil, errBadRequest("desk_id is required")
	}

	// Collect dates from either booking_date or booking_dates
	var allDates []string
	if singleDate := strings.TrimSpace(req.Data.Attributes.BookingDate); singleDate != "" {
		allDates = append(allDates, singleDate)
	}
	for _, d := range req.Data.Attributes.BookingDates {
		if trimmed := strings.TrimSpace(d); trimmed != "" {
			allDates = append(allDates, trimmed)
		}
	}

	if len(allDates) == 0 {
		return "", nil, errBadRequest("booking_date or booking_dates is required")
	}

	// Validate and dedupe dates
	today := time.Now().UTC().Truncate(24 * time.Hour)
	seen := make(map[string]struct{})
	var validDates []string

	for _, dateStr := range allDates {
		parsedDate, parseErr := time.Parse(time.DateOnly, dateStr)
		if parseErr != nil {
			return "", nil, errBadRequest("booking_date must be in YYYY-MM-DD format: " + dateStr)
		}
		if parsedDate.Before(today) {
			return "", nil, errBadRequest("booking_date cannot be in the past: " + dateStr)
		}
		if _, exists := seen[dateStr]; !exists {
			seen[dateStr] = struct{}{}
			validDates = append(validDates, dateStr)
		}
	}

	return deskID, validDates, nil
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

func processBooking(
	c echo.Context, store *sql.DB, notifier notifications.Notifier,
	deskID, userID, userName, bookedByUserID, bookedByUserName, bookingDate string,
	isGuest bool, guestEmail string,
) error {
	ctx := c.Request().Context()

	// Skip duplicate check for guests (they have unique IDs)
	if !isGuest {
		existingBookingID, err := FindUserBooking(ctx, store, deskID, userID, bookingDate)
		if err != nil {
			return fmt.Errorf("check existing booking: %w", err)
		}
		if existingBookingID != "" {
			if userID == bookedByUserID {
				//nolint:wrapcheck // Terminal response, no wrapping needed
				return api.WriteConflict(c, "You already have this desk booked for this date")
			}
			//nolint:wrapcheck // Terminal response, no wrapping needed
			return api.WriteConflict(c, "This user already has this desk booked for this date")
		}
	}

	booking, err := CreateBooking(
		ctx, store, deskID, userID, userName,
		bookedByUserID, bookedByUserName, bookingDate,
		isGuest, guestEmail,
	)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			slog.Warn("booking conflict",
				"desk_id", deskID,
				"user_id", userID,
				"booked_by", bookedByUserID,
				"booking_date", bookingDate,
			)
			//nolint:wrapcheck // Terminal response, no wrapping needed
			return api.WriteConflict(c, "Desk is already booked for this date")
		}
		return fmt.Errorf("create booking: %w", err)
	}

	logFields := []any{
		"booking_id", booking.ID,
		"desk_id", deskID,
		"user_id", userID,
		"booking_date", bookingDate,
	}
	if bookedByUserID != userID {
		logFields = append(logFields, "booked_by", bookedByUserID)
	}
	if isGuest {
		logFields = append(logFields, "is_guest", true)
	}
	slog.Info("booking created", logFields...)

	// Send notification asynchronously
	sendBookingCreatedNotification(notifier, booking)

	return writeBookingResponse(c, booking)
}

// processMultiDayBooking creates bookings for multiple dates.
// Returns created bookings and reports conflicts per day.
func processMultiDayBooking(
	c echo.Context, store *sql.DB, notifier notifications.Notifier,
	deskID string, params *bookingParticipants, dates []string,
) error {
	ctx := c.Request().Context()

	created := make([]api.Resource, 0, len(dates))
	var conflicts []string

	for _, bookingDate := range dates {
		// Skip duplicate check for guests (they have unique IDs)
		if !params.isGuest {
			existingBookingID, err := FindUserBooking(
				ctx, store, deskID, params.targetUserID, bookingDate,
			)
			if err != nil {
				return fmt.Errorf("check existing booking: %w", err)
			}
			if existingBookingID != "" {
				conflicts = append(conflicts, bookingDate+": user already has this desk booked")
				continue
			}
		}

		booking, err := CreateBooking(
			ctx, store, deskID, params.targetUserID, params.targetUserName,
			params.bookedByUserID, params.bookedByUserName, bookingDate,
			params.isGuest, params.guestEmail,
		)
		if err != nil {
			if errors.Is(err, ErrConflict) {
				conflicts = append(conflicts, bookingDate+": desk already booked")
				continue
			}
			return fmt.Errorf("create booking: %w", err)
		}

		slog.Info("booking created",
			"booking_id", booking.ID,
			"desk_id", deskID,
			"user_id", params.targetUserID,
			"booking_date", bookingDate,
		)

		// Send notification asynchronously
		sendBookingCreatedNotification(notifier, booking)

		attrs := BookingAttributes{
			DeskID:      booking.DeskID,
			UserID:      booking.UserID,
			BookingDate: booking.BookingDate,
			CreatedAt:   booking.CreatedAt,
		}
		if booking.BookedByUserID != "" && booking.BookedByUserID != booking.UserID {
			attrs.BookedByUserID = booking.BookedByUserID
			attrs.BookedByUserName = booking.BookedByUserName
		}
		if booking.IsGuest {
			attrs.IsGuest = true
			attrs.GuestEmail = booking.GuestEmail
		}

		created = append(created, api.Resource{
			Type:       "bookings",
			ID:         booking.ID,
			Attributes: attrs,
		})
	}

	// Return multi-day response
	result := MultiDayBookingResult{
		Created:   created,
		Conflicts: conflicts,
	}

	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	//nolint:wrapcheck // Terminal response
	return c.JSON(http.StatusCreated, result)
}

func writeBookingResponse(c echo.Context, booking *Booking) error {
	attrs := BookingAttributes{
		DeskID:      booking.DeskID,
		UserID:      booking.UserID,
		BookingDate: booking.BookingDate,
		CreatedAt:   booking.CreatedAt,
	}
	// Include booked_by info if booking was made on behalf
	if booking.BookedByUserID != "" && booking.BookedByUserID != booking.UserID {
		attrs.BookedByUserID = booking.BookedByUserID
		attrs.BookedByUserName = booking.BookedByUserName
	}
	// Include guest info
	if booking.IsGuest {
		attrs.IsGuest = true
		attrs.GuestEmail = booking.GuestEmail
	}

	resp := api.SingleResponse{
		Data: api.Resource{
			Type:       "bookings",
			ID:         booking.ID,
			Attributes: attrs,
		},
	}

	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	//nolint:wrapcheck // Terminal response, no wrapping needed
	return c.JSON(http.StatusCreated, resp)
}

// Booking represents a booking record.
type Booking struct {
	ID               string
	DeskID           string
	UserID           string
	UserName         string
	BookingDate      string
	BookedByUserID   string
	BookedByUserName string
	IsGuest          bool
	GuestEmail       string
	CreatedAt        string
	UpdatedAt        string
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
func CreateBooking(
	ctx context.Context, store *sql.DB,
	deskID, userID, userName, bookedByUserID, bookedByUserName, bookingDate string,
	isGuest bool, guestEmail string,
) (*Booking, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.New().String()

	isGuestInt := 0
	if isGuest {
		isGuestInt = 1
	}

	_, err := store.ExecContext(ctx, `
		INSERT INTO bookings 
		(id, desk_id, user_id, user_name, booked_by_user_id, booked_by_user_name, 
		 booking_date, is_guest, guest_email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, deskID, userID, userName, bookedByUserID, bookedByUserName,
		bookingDate, isGuestInt, guestEmail, now, now,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("insert booking: %w", err)
	}

	return &Booking{
		ID:               id,
		DeskID:           deskID,
		UserID:           userID,
		UserName:         userName,
		BookedByUserID:   bookedByUserID,
		BookedByUserName: bookedByUserName,
		BookingDate:      bookingDate,
		IsGuest:          isGuest,
		GuestEmail:       guestEmail,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// sendBookingCreatedNotification sends an async notification for a created booking.
func sendBookingCreatedNotification(notifier notifications.Notifier, booking *Booking) {
	event := notifications.BookingEvent{
		Event:       notifications.EventBookingCreated,
		BookingID:   booking.ID,
		DeskID:      booking.DeskID,
		UserID:      booking.UserID,
		UserName:    booking.UserName,
		BookingDate: booking.BookingDate,
		IsGuest:     booking.IsGuest,
		GuestEmail:  booking.GuestEmail,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	if booking.BookedByUserID != "" && booking.BookedByUserID != booking.UserID {
		event.BookedByUserID = booking.BookedByUserID
		event.BookedByUserName = booking.BookedByUserName
	}
	notifier.NotifyAsync(&event)
}
