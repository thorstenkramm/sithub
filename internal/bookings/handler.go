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
	"github.com/thorstenkramm/sithub/internal/areas"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/notifications"
	"github.com/thorstenkramm/sithub/internal/users"
)

// BookingLimits holds the booking limit configuration.
// It is safe for concurrent use after construction.
type BookingLimits struct {
	WeeksInAdvanced      int
	MaxBookingsPerPerson int
}

// ErrBookingLimitExceeded indicates a booking limit was reached.
var ErrBookingLimitExceeded = errors.New("booking limit exceeded")

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
			Note         string   `json:"note,omitempty"`
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
	ForUserID        string `json:"for_user_id,omitempty"`
	ForUserName      string `json:"for_user_name,omitempty"`
	IsGuest          bool   `json:"is_guest,omitempty"`
	GuestName        string `json:"guest_name,omitempty"`
	GuestEmail       string `json:"guest_email,omitempty"`
	Note             string `json:"note"`
}

// maxNoteLength is the maximum allowed length for a booking note.
const (
	maxNoteLength       = 500
	resourceTypeBooking = "bookings"
)

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
	if req.Data.Type != resourceTypeBooking {
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
			Type:       resourceTypeBooking,
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
func ListHandler(cfg *areas.Config, store *sql.DB) echo.HandlerFunc {
	return ListHandlerDynamic(func() *areas.Config { return cfg }, store)
}

// ListHandlerDynamic returns a handler for listing the current user's future bookings.
func ListHandlerDynamic(getConfig areas.ConfigGetter, store *sql.DB) echo.HandlerFunc {
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
func HistoryHandler(cfg *areas.Config, store *sql.DB) echo.HandlerFunc {
	return HistoryHandlerDynamic(func() *areas.Config { return cfg }, store)
}

// HistoryHandlerDynamic returns a handler for listing the current user's past bookings.
func HistoryHandlerDynamic(getConfig areas.ConfigGetter, store *sql.DB) echo.HandlerFunc {
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
	ctx context.Context, c echo.Context, cfg *areas.Config, store *sql.DB,
	currentUserID string, records []BookingRecord,
) error {
	// Collect unique user IDs for display name lookup. This includes the booker
	// (BookedByUserID) as well as the colleague a booking was made FOR (UserID),
	// so on-behalf-by-me bookings can surface the colleague's name (FR168).
	userIDSet := make(map[string]struct{})
	for i := range records {
		if records[i].BookedByUserID != "" {
			userIDSet[records[i].BookedByUserID] = struct{}{}
		}
		if records[i].UserID != "" && records[i].UserID != currentUserID {
			userIDSet[records[i].UserID] = struct{}{}
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

		resources = append(resources, api.Resource{
			Type:       resourceTypeBooking,
			ID:         rec.ID,
			Attributes: buildMyBookingAttributes(rec, loc, currentUserID, displayNames),
		})
	}

	resp := api.CollectionResponse{Data: resources}
	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	//nolint:wrapcheck // Terminal response
	return c.JSON(http.StatusOK, resp)
}

// buildMyBookingAttributes maps a booking record and its resolved item location into the
// JSON:API attributes for the My Bookings collection. currentUserID identifies the requesting
// user so on-behalf and booked-for-me relationships can be derived; displayNames resolves the
// booker and colleague names.
func buildMyBookingAttributes(
	rec *BookingRecord, loc *areas.ItemLocation,
	currentUserID string, displayNames map[string]string,
) MyBookingAttributes {
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
	// For an on-behalf booking made BY the current user, surface the colleague's
	// id and name (rec.UserID) so the UI can render "On behalf of <name>" (FR168)
	// and detect same area/day conflicts for that colleague (FR178).
	if rec.BookedByUserID == currentUserID && rec.UserID != currentUserID {
		attrs.ForUserID = rec.UserID
		if name, ok := displayNames[rec.UserID]; ok {
			attrs.ForUserName = name
		}
	}
	// Include guest info
	if rec.IsGuest {
		attrs.IsGuest = true
		attrs.GuestName = rec.GuestName
		attrs.GuestEmail = rec.GuestEmail
	}

	return attrs
}

// CreateHandler returns a handler for creating bookings.
// Supports single-day, multi-day, booking on behalf, and guest bookings.
func CreateHandler(cfg *areas.Config, store *sql.DB, notifier notifications.Notifier) echo.HandlerFunc {
	return CreateHandlerDynamic(func() *areas.Config { return cfg }, store, notifier, nil)
}

// CreateHandlerDynamic returns a handler for creating bookings using dynamic config.
func CreateHandlerDynamic(
	getConfig areas.ConfigGetter, store *sql.DB, notifier notifications.Notifier,
	limits *BookingLimits,
) echo.HandlerFunc {
	maxWeeks := 0
	if limits != nil {
		maxWeeks = limits.WeeksInAdvanced
	}

	return func(c echo.Context) error {
		user := auth.GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		req, itemID, dates, err := parseAndValidateBooking(c, maxWeeks)
		if err != nil || c.Response().Committed {
			return err
		}

		note := strings.TrimSpace(req.Data.Attributes.Note)
		if len(note) > maxNoteLength {
			return handleValidationError(c, errBadRequest(
				fmt.Sprintf("Note must be at most %d characters", maxNoteLength)))
		}

		cfg := getConfig()
		loc, exists := cfg.FindItemLocation(itemID)
		if !exists {
			return api.WriteNotFound(c, "Item not found")
		}

		params, err := resolveBookingParticipants(c.Request().Context(), store, user, req)
		if err != nil {
			return handleValidationError(c, err)
		}

		// Reservation access applies to the eventual booking target, not just the acting user.
		if err := handleReservation(c, store, params, loc); err != nil || c.Response().Committed {
			return err
		}

		if err := handleBookingLimits(c, store, params, loc, limits); err != nil || c.Response().Committed {
			return err
		}

		if len(dates) == 1 {
			return processBooking(
				c, store, notifier, itemID, params.targetUserID,
				params.bookedByUserID, dates[0], note,
				params.isGuest, params.guestName, params.guestEmail,
			)
		}

		return processMultiDayBooking(c, store, notifier, itemID, params, dates, note)
	}
}

// handleReservation checks if the user is allowed to book the item based on
// reserved_for configuration. Returns nil when access is granted or no
// reservations are configured.
func handleReservation(
	c echo.Context, store *sql.DB, params *bookingParticipants, loc *areas.ItemLocation,
) error {
	// Skip check if no reserved_for at any level
	if len(loc.Item.ReservedFor) == 0 &&
		len(loc.ItemGroup.ReservedFor) == 0 &&
		len(loc.Area.ReservedFor) == 0 {
		return nil
	}

	targetEmail, err := resolveReservationEmail(c.Request().Context(), store, params)
	if err != nil {
		return fmt.Errorf("lookup user for reservation check: %w", err)
	}
	if targetEmail == "" {
		//nolint:errcheck // Response signals via Committed
		api.WriteForbiddenDetail(c, reservationForbiddenMessage(loc))
		return nil
	}

	if areas.IsReserved(loc, targetEmail) {
		//nolint:errcheck // Response signals via Committed
		api.WriteForbiddenDetail(c, reservationForbiddenMessage(loc))
		return nil
	}
	return nil
}

func resolveReservationEmail(
	ctx context.Context, store *sql.DB, params *bookingParticipants,
) (string, error) {
	if params.isGuest {
		return strings.TrimSpace(params.guestEmail), nil
	}

	rec, err := users.FindByID(ctx, store, params.targetUserID)
	if err != nil {
		return "", fmt.Errorf("find reservation target user: %w", err)
	}
	if rec == nil {
		return "", nil
	}
	return strings.TrimSpace(rec.Email), nil
}

func reservationForbiddenMessage(loc *areas.ItemLocation) string {
	kind, label := "item", loc.Item.Name

	switch {
	case len(loc.Item.ReservedFor) > 0:
		kind, label = "item", firstNonEmpty(loc.Item.Name, loc.Item.ID)
	case len(loc.ItemGroup.ReservedFor) > 0:
		kind, label = "item group", firstNonEmpty(loc.ItemGroup.Name, loc.ItemGroup.ID)
	case len(loc.Area.ReservedFor) > 0:
		kind, label = "area", firstNonEmpty(loc.Area.Name, loc.Area.ID)
	}

	return fmt.Sprintf("This %s %q is reserved. You do not have access.", kind, label)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

// handleBookingLimits enforces booking limits and writes the error response if exceeded.
// Returns nil when limits pass, the user is a guest, or a conflict response was written.
func handleBookingLimits(
	c echo.Context, store *sql.DB, params *bookingParticipants,
	loc *areas.ItemLocation, limits *BookingLimits,
) error {
	if params.isGuest {
		return nil
	}
	err := enforceBookingLimits(
		c.Request().Context(), store, params.targetUserID, loc, limits,
	)
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrBookingLimitExceeded) {
		//nolint:errcheck // Error is intentionally ignored; caller checks response state
		api.WriteConflict(c, err.Error())
		return nil
	}
	return fmt.Errorf("check booking limits: %w", err)
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

// parseAndValidateBooking handles content-type check, JSON parsing, and field validation.
// Returns the parsed request, item ID, and dates. On validation failure, the error
// response is written and c.Response().Committed is true.
func parseAndValidateBooking(
	c echo.Context, maxWeeks int,
) (req *CreateRequest, itemID string, dates []string, err error) {
	if err = validateContentType(c); err != nil {
		if errors.Is(err, errResponseWritten) {
			return nil, "", nil, nil
		}
		return nil, "", nil, err
	}

	req, err = parseCreateRequest(c)
	if err != nil {
		return nil, "", nil, handleValidationError(c, err)
	}

	itemID, dates, err = validateRequestFieldsMultiDay(req, maxWeeks)
	if err != nil {
		return nil, "", nil, handleValidationError(c, err)
	}

	return req, itemID, dates, nil
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
	if req.Data.Type != resourceTypeBooking {
		return nil, errBadRequest("Resource type must be 'bookings'")
	}
	return &req, nil
}

func validateRequestFieldsMultiDay(
	req *CreateRequest, maxWeeks int,
) (itemID string, dates []string, err error) {
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
	// Calculate the booking horizon: current week + maxWeeks additional weeks
	var maxDate time.Time
	if maxWeeks > 0 {
		// End of the Nth week from today's week (Sunday of that week)
		weekday := today.Weekday()
		daysUntilMonday := (8 - int(weekday)) % 7
		nextMonday := today.AddDate(0, 0, daysUntilMonday)
		maxDate = nextMonday.AddDate(0, 0, maxWeeks*7)
	}

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
		if maxWeeks > 0 && !parsedDate.Before(maxDate) {
			return "", nil, errBadRequest(fmt.Sprintf(
				"booking_date is too far in the future (maximum %d weeks in advance): %s",
				maxWeeks, dateStr,
			))
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

// enforceBookingLimits checks all applicable booking limits for a user.
// It checks at item, item group, area, and global scope.
// Returns ErrBookingLimitExceeded (wrapped with a descriptive message) if any limit is reached.
func enforceBookingLimits(
	ctx context.Context, store *sql.DB, userID string,
	loc *areas.ItemLocation, limits *BookingLimits,
) error {
	if limits == nil {
		return nil
	}

	// Check item-level limit (most specific)
	if err := checkBookingLimit(
		ctx, store, userID, loc.Item.MaxBookingsPerPerson,
		[]string{loc.Item.ID},
		fmt.Sprintf("'%s, %s'", loc.ItemGroup.Name, loc.Item.Name),
	); err != nil {
		return err
	}

	// Check item group level limit
	if err := checkBookingLimit(
		ctx, store, userID, loc.ItemGroup.MaxBookingsPerPerson,
		collectItemIDs(loc.ItemGroup.Items),
		fmt.Sprintf("'%s'", loc.ItemGroup.Name),
	); err != nil {
		return err
	}

	// Check area-level limit
	if err := checkBookingLimit(
		ctx, store, userID, loc.Area.MaxBookingsPerPerson,
		collectAreaItemIDs(loc.Area),
		fmt.Sprintf("'%s'", loc.Area.Name),
	); err != nil {
		return err
	}

	// Check global limit
	return checkBookingLimit(
		ctx, store, userID, limits.MaxBookingsPerPerson, nil, "",
	)
}

// checkBookingLimit verifies that a user has not reached the given limit
// for the specified item IDs. A limit of 0 means unlimited (no check).
// When scopeLabel is empty, the error message omits the scope.
func checkBookingLimit(
	ctx context.Context, store *sql.DB, userID string,
	limit int, itemIDs []string, scopeLabel string,
) error {
	if limit <= 0 {
		return nil
	}
	count, err := CountUserFutureBookings(ctx, store, userID, itemIDs)
	if err != nil {
		return err
	}
	if count < limit {
		return nil
	}
	if scopeLabel != "" {
		return fmt.Errorf(
			"%w: you have reached the maximum of %d active bookings for %s",
			ErrBookingLimitExceeded, limit, scopeLabel,
		)
	}
	return fmt.Errorf(
		"%w: you have reached the maximum of %d active bookings",
		ErrBookingLimitExceeded, limit,
	)
}

// collectItemIDs returns all item IDs from an item group's items.
func collectItemIDs(items []areas.Item) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}

// collectAreaItemIDs returns all item IDs across all item groups in an area.
func collectAreaItemIDs(area *areas.Area) []string {
	var ids []string
	for i := range area.ItemGroups {
		ig := &area.ItemGroups[i]
		for _, item := range ig.Items {
			ids = append(ids, item.ID)
		}
	}
	return ids
}

// errResponseWritten indicates the HTTP response was already written.
var errResponseWritten = errors.New("response already written")

func processBooking(
	c echo.Context, store *sql.DB, notifier notifications.Notifier,
	itemID, userID, bookedByUserID, bookingDate, note string,
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
		bookedByUserID, bookingDate, note,
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
	itemID string, params *bookingParticipants, dates []string, note string,
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
			params.bookedByUserID, bookingDate, note,
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
			Type:       resourceTypeBooking,
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
			Type:       resourceTypeBooking,
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
	itemID, userID, bookedByUserID, bookingDate, note string,
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
		 is_guest, guest_name, guest_email, note, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, itemID, userID, bookedByUserID,
		bookingDate, isGuestInt, guestName, guestEmail, note, now, now,
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
		Note:           note,
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
