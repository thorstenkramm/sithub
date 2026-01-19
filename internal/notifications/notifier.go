// Package notifications provides booking notification services.
package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// EventType represents the type of booking event.
type EventType string

const (
	// EventBookingCreated is sent when a booking is created.
	EventBookingCreated EventType = "booking.created"
	// EventBookingCanceled is sent when a booking is canceled.
	EventBookingCanceled EventType = "booking.canceled"
)

// BookingEvent represents a notification payload for booking events.
type BookingEvent struct {
	Event       EventType `json:"event"`
	BookingID   string    `json:"booking_id"`
	DeskID      string    `json:"desk_id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	BookingDate string    `json:"booking_date"`
	IsGuest     bool      `json:"is_guest,omitempty"`
	GuestEmail  string    `json:"guest_email,omitempty"`
	// BookedByUserID is set when booking was made on behalf of someone.
	BookedByUserID   string `json:"booked_by_user_id,omitempty"`
	BookedByUserName string `json:"booked_by_user_name,omitempty"`
	// CanceledByUserID is set when a booking is canceled.
	CanceledByUserID string `json:"canceled_by_user_id,omitempty"`
	// Timestamp is when the event occurred.
	Timestamp string `json:"timestamp"`
}

// Notifier sends booking notifications.
type Notifier interface {
	// NotifyAsync sends a notification asynchronously.
	// It returns immediately and logs any errors.
	NotifyAsync(event *BookingEvent)
}

// NoopNotifier is a no-op notifier that discards all events.
type NoopNotifier struct{}

// NotifyAsync does nothing.
func (n *NoopNotifier) NotifyAsync(_ *BookingEvent) {}

// WebhookNotifier sends notifications via HTTP webhook.
type WebhookNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewWebhookNotifier creates a new webhook notifier.
func NewWebhookNotifier(webhookURL string) *WebhookNotifier {
	return &WebhookNotifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NotifyAsync sends the event to the webhook URL asynchronously.
func (n *WebhookNotifier) NotifyAsync(event *BookingEvent) {
	go func() {
		if err := n.send(event); err != nil {
			slog.Error("failed to send notification",
				"event", event.Event,
				"booking_id", event.BookingID,
				"error", err,
			)
		} else {
			slog.Info("notification sent",
				"event", event.Event,
				"booking_id", event.BookingID,
			)
		}
	}()
}

func (n *WebhookNotifier) send(event *BookingEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	body, err := json.Marshal(*event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SitHub/1.0")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// NewNotifier creates a notifier based on configuration.
// Returns a NoopNotifier if webhookURL is empty.
func NewNotifier(webhookURL string) Notifier {
	if webhookURL == "" {
		slog.Info("notifications disabled (no webhook URL configured)")
		return &NoopNotifier{}
	}
	slog.Info("notifications enabled", "webhook_url", webhookURL)
	return NewWebhookNotifier(webhookURL)
}
