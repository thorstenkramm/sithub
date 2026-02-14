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
	"github.com/thorstenkramm/sithub/internal/users"
)

// CreateRequest represents a booking create JSON:API payload.
type CreateRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			ItemID       string   `json:"item_id"`
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
	ItemID         string `json:"item_id"`
	UserID         string `json:"user_id"`
	BookingDate    string `json:"booking_date"`
	CreatedAt      string `json:"created_at"`
	BookedByUserID string `json:"booked_by_user_id,omitempty"`
	IsGuest        bool   `json:"is_guest,omitempty"`
	GuestEmail     string `json:"guest_email,omitempty"`
	Note           string `json:"note"`
}

// MultiDayBookingResult represents the result of a multi-day booking request.
type MultiDayBookingResult struct {
	Created   []api.Resource `json:"created"`
	Conflicts []string       `json:"conflicts,omitempty"`
}

// MyBookingAttributes represents booking resource attributes with location info.
type MyBookingAttributes struct {
	ItemID           string `json:"item_id"`
	ItemName         string `json:"item_name"`
	ItemGroupID      string `json:"item_group_id"`
	ItemGroupName    string `json:"item_group_name"`
	AreaID           string `json:"area_id"`
	AreaName         string `json:"area_name"`
	BookingDate      string `json:"booking_date"`
	CreatedAt        string `json:"created_at"`
	BookedByUserID   string `json:"booked_by_user_id,omitempty"`
	BookedByUserName string `json:"booked_by_user_name,omitempty"`
	BookedForMe      bool   `json:"booked_for_me,omitempty"`
	IsGuest          bool   `json:"is_guest,omitempty"`
	GuestEmail       string `json:"guest_email,omitempty"`
	Note             string `json:"note"`
}

// maxNoteLength is the maximum allowed length for a booking note.
const maxNoteLength = 500

// PatchRequest represents a booking update JSON:API payload.
type PatchRequest struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			Note *string `json:"note"`
		} `json:"attributes"`
	} `json:"data"`
}

// PatchHandler returns a handler for updating a booking's note.
// Authorization: booking owner, the person who booked, or admin.
func PatchHandler(store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := auth.GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		bookingID := c.Param("id")
		if bookingID == "" {
			return api.WriteBadRequest(c, "Booking ID is required")
		}

		note, err := parsePatchNote(c, bookingID)
		if err != nil {
			if errors.Is(err, errResponseWritten) {
				return nil
			}
			return err
		}

		ctx := c.Request().Context()
		booking, err := findAuthorizedBooking(ctx, store, bookingID, user)
		if errors.Is(err, ErrBookingNotFound) {
			return api.WriteNotFound(c, "Booking not found")
		}
		if err != nil {
			return err
		}

		if err := UpdateNote(ctx, store, bookingID, note); err != nil {
			return fmt.Errorf("update note: %w", err)
		}

		slog.Info("booking note updated",
			"booking_id", bookingID,
			"updated_by", user.ID,
		)

		return writePatchResponse(c, booking, note)
	}
}

// parsePatchNote validates content type and parses the note from the PATCH request body.
// Returns errResponseWritten if an error response was already sent.
func parsePatchNote(c echo.Context, bookingID string) (string, error) {
	if err := validateContentType(c); err != nil {
		return "", err
	}

	var req PatchRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		//nolint:errcheck // Error ignored; response already written
		api.WriteBadRequest(c, "Invalid request body")
		return "", errResponseWritten
	}
	if req.Data.Type != "bookings" {
		//nolint:errcheck // Error ignored; response already written
		api.WriteBadRequest(c, "Resource type must be 'bookings'")
		return "", errResponseWritten
	}
	if req.Data.ID == "" {
		//nolint:errcheck // Error ignored; response already written
		api.WriteBadRequest(c, "Resource ID is required")
		return "", errResponseWritten
	}
	if req.Data.ID != bookingID {
		//nolint:errcheck // Error ignored; response already written
		api.WriteBadRequest(c, "Resource ID must match booking ID")
		return "", errResponseWritten
	}
	if req.Data.Attributes.Note == nil {
		//nolint:errcheck // Error ignored; response already written
		api.WriteBadRequest(c, "Note is required")
		return "", errResponseWritten
	}

	note := strings.TrimSpace(*req.Data.Attributes.Note)
	if len(note) > maxNoteLength {
		//nolint:errcheck // Error ignored; response already written
		api.WriteBadRequest(c, fmt.Sprintf(
			"Note must be at most %d characters", maxNoteLength))
		return "", errResponseWritten
	}

	return note, nil
}

// findAuthorizedBooking retrieves a booking and checks that the user is authorized.
// Returns the booking record or writes an error response and returns a terminal error.
func findAuthorizedBooking(
	ctx context.Context, store *sql.DB, bookingID string, user *auth.User,
) (*BookingRecord, error) {
	booking, err := FindBookingByID(ctx, store, bookingID)
	if err != nil {
		return nil, fmt.Errorf("find booking: %w", err)
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}

	isOwner := booking.UserID == user.ID
	isBooker := booking.BookedByUserID == user.ID
	if !isOwner && !isBooker && !user.IsAdmin {
		return nil, ErrBookingNotFound
	}

	return booking, nil
}

// ErrBookingNotFound is a sentinel error for booking not found responses.
var ErrBookingNotFound = errors.New("booking not found")

func writePatchResponse(c echo.Context, booking *BookingRecord, note string) error {
	attrs := BookingAttributes{
		ItemID:      booking.ItemID,
		UserID:      booking.UserID,
		BookingDate: booking.BookingDate,
		CreatedAt:   booking.CreatedAt,
		Note:        note,
	}
	if booking.BookedByUserID != "" && booking.BookedByUserID != booking.UserID {
		attrs.BookedByUserID = booking.BookedByUserID
	}
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
	//nolint:wrapcheck // Terminal response
	return c.JSON(http.StatusOK, resp)
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
			"item_id", booking.ItemID,
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
			ItemID:           booking.ItemID,
			UserID:           booking.UserID,
			BookingDate:      booking.BookingDate,
			IsGuest:          booking.IsGuest,
			GuestName:        booking.GuestName,
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

		return writeBookingsCollection(ctx, c, getConfig(), store, user.ID, records)
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

		return writeBookingsCollection(ctx, c, getConfig(), store, user.ID, records)
	}
}

func writeBookingsCollection(
	ctx context.Context, c echo.Context, cfg *spaces.Config, store *sql.DB,
	currentUserID string, records []BookingRecord,
) error {
	// Collect unique user IDs for display name lookup
	userIDSet := make(map[string]struct{})
	for i := range records {
		if records[i].BookedByUserID != "" {
			userIDSet[records[i].BookedByUserID] = struct{}{}
		}
	}
	userIDs := make([]string, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	displayNames, err := users.FindDisplayNames(ctx, store, userIDs)
	if err != nil {
		slog.Warn("failed to look up display names", "error", err)
		displayNames = map[string]string{}
	}

	resources := make([]api.Resource, 0, len(records))
	for i := range records {
		rec := &records[i]
		loc, found := cfg.FindItemLocation(rec.ItemID)
		if !found {
			slog.Warn("booking references unknown item", "booking_id", rec.ID, "item_id", rec.ItemID)
			continue
		}

		attrs := MyBookingAttributes{
			ItemID:        rec.ItemID,
			ItemName:      loc.Item.Name,
			ItemGroupID:   loc.ItemGroup.ID,
			ItemGroupName: loc.ItemGroup.Name,
			AreaID:        loc.Area.ID,
			AreaName:      loc.Area.Name,
			BookingDate:   rec.BookingDate,
			CreatedAt:     rec.CreatedAt,
			Note:          rec.Note,
		}

		// Include booked_by info if different from user_id
		if rec.BookedByUserID != "" && rec.BookedByUserID != rec.UserID {
			attrs.BookedByUserID = rec.BookedByUserID
			if name, ok := displayNames[rec.BookedByUserID]; ok {
				attrs.BookedByUserName = name
			}
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

		itemID, dates, err := validateRequestFieldsMultiDay(req)
		if err != nil {
			return handleValidationError(c, err)
		}

		cfg := getConfig()
		if _, exists := cfg.FindItem(itemID); !exists {
			return api.WriteNotFound(c, "Item not found")
		}

		params, err := resolveBookingParticipants(c.Request().Context(), store, user, req)
		if err != nil {
			return handleValidationError(c, err)
		}

		// Single day booking
		if len(dates) == 1 {
			return processBooking(
				c, store, notifier, itemID, params.targetUserID,
				params.bookedByUserID, dates[0],
				params.isGuest, params.guestName, params.guestEmail,
			)
		}

		// Multi-day booking
		return processMultiDayBooking(c, store, notifier, itemID, params, dates)
	}
}

// bookingParticipants holds the resolved user info for a booking.
type bookingParticipants struct {
	targetUserID   string
	bookedByUserID string
	isGuest        bool
	guestName      string
	guestEmail     string
}

// resolveBookingParticipants determines the target user and booker for a booking.
// When for_user_id is provided, it validates that the target user exists in the database.
func resolveBookingParticipants(
	ctx context.Context, store *sql.DB, user *auth.User, req *CreateRequest,
) (*bookingParticipants, error) {
	params := &bookingParticipants{
		targetUserID:   user.ID,
		bookedByUserID: user.ID,
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
		params.isGuest = true
		params.guestName = guestName
		params.guestEmail = guestEmail
		return params, nil
	}

	// Handle booking on behalf of another user
	forUserID := strings.TrimSpace(req.Data.Attributes.ForUserID)
	if forUserID != "" {
		if _, err := users.FindByID(ctx, store, forUserID); err != nil {
			return nil, errBadRequest("for_user_id: user not found")
		}
		params.targetUserID = forUserID
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

func validateRequestFieldsMultiDay(req *CreateRequest) (itemID string, dates []string, err error) {
	itemID = strings.TrimSpace(req.Data.Attributes.ItemID)
	if itemID == "" {
		return "", nil, errBadRequest("item_id is required")
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

	return itemID, validDates, nil
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
	itemID, userID, bookedByUserID, bookingDate string,
	isGuest bool, guestName, guestEmail string,
) error {
	ctx := c.Request().Context()

	// Skip duplicate check for guests (they have unique IDs)
	if !isGuest {
		existingBookingID, err := FindUserBooking(ctx, store, itemID, userID, bookingDate)
		if err != nil {
			return fmt.Errorf("check existing booking: %w", err)
		}
		if existingBookingID != "" {
			if userID == bookedByUserID {
				//nolint:wrapcheck // Terminal response, no wrapping needed
				return api.WriteConflict(c, "You already have this item booked for this date")
			}
			//nolint:wrapcheck // Terminal response, no wrapping needed
			return api.WriteConflict(c, "This user already has this item booked for this date")
		}
	}

	booking, err := CreateBooking(
		ctx, store, itemID, userID,
		bookedByUserID, bookingDate,
		isGuest, guestName, guestEmail,
	)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			slog.Warn("booking conflict",
				"item_id", itemID,
				"user_id", userID,
				"booked_by", bookedByUserID,
				"booking_date", bookingDate,
			)
			//nolint:wrapcheck // Terminal response, no wrapping needed
			return api.WriteConflict(c, "Item is already booked for this date")
		}
		return fmt.Errorf("create booking: %w", err)
	}

	logFields := []any{
		"booking_id", booking.ID,
		"item_id", itemID,
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
	itemID string, params *bookingParticipants, dates []string,
) error {
	ctx := c.Request().Context()

	created := make([]api.Resource, 0, len(dates))
	var conflicts []string

	for _, bookingDate := range dates {
		// Skip duplicate check for guests (they have unique IDs)
		if !params.isGuest {
			existingBookingID, err := FindUserBooking(
				ctx, store, itemID, params.targetUserID, bookingDate,
			)
			if err != nil {
				return fmt.Errorf("check existing booking: %w", err)
			}
			if existingBookingID != "" {
				conflicts = append(conflicts, bookingDate+": user already has this item booked")
				continue
			}
		}

		booking, err := CreateBooking(
			ctx, store, itemID, params.targetUserID,
			params.bookedByUserID, bookingDate,
			params.isGuest, params.guestName, params.guestEmail,
		)
		if err != nil {
			if errors.Is(err, ErrConflict) {
				conflicts = append(conflicts, bookingDate+": item already booked")
				continue
			}
			return fmt.Errorf("create booking: %w", err)
		}

		slog.Info("booking created",
			"booking_id", booking.ID,
			"item_id", itemID,
			"user_id", params.targetUserID,
			"booking_date", bookingDate,
		)

		// Send notification asynchronously
		sendBookingCreatedNotification(notifier, booking)

		attrs := BookingAttributes{
			ItemID:      booking.ItemID,
			UserID:      booking.UserID,
			BookingDate: booking.BookingDate,
			CreatedAt:   booking.CreatedAt,
			Note:        booking.Note,
		}
		if booking.BookedByUserID != "" && booking.BookedByUserID != booking.UserID {
			attrs.BookedByUserID = booking.BookedByUserID
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
		ItemID:      booking.ItemID,
		UserID:      booking.UserID,
		BookingDate: booking.BookingDate,
		CreatedAt:   booking.CreatedAt,
		Note:        booking.Note,
	}
	// Include booked_by info if booking was made on behalf
	if booking.BookedByUserID != "" && booking.BookedByUserID != booking.UserID {
		attrs.BookedByUserID = booking.BookedByUserID
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
	ID             string
	ItemID         string
	UserID         string
	BookingDate    string
	BookedByUserID string
	IsGuest        bool
	GuestName      string
	GuestEmail     string
	Note           string
	CreatedAt      string
	UpdatedAt      string
}

// ErrConflict indicates a booking conflict (item already booked).
var ErrConflict = errors.New("booking conflict")

// FindUserBooking checks if a user already has a booking for a specific item and date.
// Returns the booking ID if found, empty string otherwise.
func FindUserBooking(ctx context.Context, store *sql.DB, itemID, userID, bookingDate string) (string, error) {
	var bookingID string
	err := store.QueryRowContext(ctx,
		"SELECT id FROM bookings WHERE item_id = ? AND user_id = ? AND booking_date = ?",
		itemID, userID, bookingDate,
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
	itemID, userID, bookedByUserID, bookingDate string,
	isGuest bool, guestName, guestEmail string,
) (*Booking, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.New().String()

	isGuestInt := 0
	if isGuest {
		isGuestInt = 1
	}

	_, err := store.ExecContext(ctx, `
		INSERT INTO bookings
		(id, item_id, user_id, booked_by_user_id, booking_date,
		 is_guest, guest_name, guest_email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, itemID, userID, bookedByUserID,
		bookingDate, isGuestInt, guestName, guestEmail, now, now,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("insert booking: %w", err)
	}

	return &Booking{
		ID:             id,
		ItemID:         itemID,
		UserID:         userID,
		BookedByUserID: bookedByUserID,
		BookingDate:    bookingDate,
		IsGuest:        isGuest,
		GuestName:      guestName,
		GuestEmail:     guestEmail,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// sendBookingCreatedNotification sends an async notification for a created booking.
func sendBookingCreatedNotification(notifier notifications.Notifier, booking *Booking) {
	event := notifications.BookingEvent{
		Event:       notifications.EventBookingCreated,
		BookingID:   booking.ID,
		ItemID:      booking.ItemID,
		UserID:      booking.UserID,
		BookingDate: booking.BookingDate,
		IsGuest:     booking.IsGuest,
		GuestEmail:  booking.GuestEmail,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	if booking.GuestName != "" {
		event.GuestName = booking.GuestName
	}
	if booking.BookedByUserID != "" && booking.BookedByUserID != booking.UserID {
		event.BookedByUserID = booking.BookedByUserID
	}
	notifier.NotifyAsync(&event)
}
