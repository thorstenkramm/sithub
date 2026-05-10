// Package livefeed broadcasts booking events to connected WebSocket clients.
//
// The hub is a single in-process fan-out: every connected client receives every
// event. Per-view filtering happens on the frontend. There is intentionally no
// per-area or per-date subscription state on the server — see Story 31.1.
package livefeed

import "github.com/thorstenkramm/sithub/internal/notifications"

// EventType identifies a live event.
type EventType string

const (
	// EventBookingCreated is broadcast when a booking is created.
	EventBookingCreated EventType = "booking.created"
	// EventBookingCanceled is broadcast when a booking is canceled.
	EventBookingCanceled EventType = "booking.canceled"
)

// Event is the over-the-wire payload sent to live-feed clients.
//
// Personally identifying details (guest name, guest email) are intentionally
// omitted so the broadcast payload is safe to deliver to all connected
// clients. Clients only need enough to refresh their visible slice and to
// recognize whether the change originated from their own session.
type Event struct {
	Type        EventType `json:"type"`
	BookingID   string    `json:"booking_id"`
	ItemID      string    `json:"item_id"`
	UserID      string    `json:"user_id"`
	BookingDate string    `json:"booking_date"`
	Timestamp   string    `json:"timestamp"`
}

// fromBookingEvent maps an internal notification event to the public live
// payload, dropping PII fields (guest_name, guest_email).
func fromBookingEvent(src *notifications.BookingEvent) Event {
	userID := src.UserID
	switch src.Event {
	case notifications.EventBookingCreated:
		if src.BookedByUserID != "" {
			userID = src.BookedByUserID
		}
	case notifications.EventBookingCanceled:
		if src.CanceledByUserID != "" {
			userID = src.CanceledByUserID
		}
	}

	return Event{
		Type:        EventType(src.Event),
		BookingID:   src.BookingID,
		ItemID:      src.ItemID,
		UserID:      userID,
		BookingDate: src.BookingDate,
		Timestamp:   src.Timestamp,
	}
}
