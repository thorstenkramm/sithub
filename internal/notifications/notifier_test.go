package notifications

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoopNotifierDoesNothing(t *testing.T) {
	t.Parallel()

	notifier := &NoopNotifier{}
	// Should not panic or do anything
	notifier.NotifyAsync(&BookingEvent{
		Event:     EventBookingCreated,
		BookingID: "test-booking",
	})
}

func TestNewNotifierReturnsNoopWhenURLEmpty(t *testing.T) {
	t.Parallel()

	notifier := NewNotifier("")
	_, ok := notifier.(*NoopNotifier)
	assert.True(t, ok, "should return NoopNotifier when webhook URL is empty")
}

func TestNewNotifierReturnsWebhookNotifierWhenURLProvided(t *testing.T) {
	t.Parallel()

	notifier := NewNotifier("https://example.com/webhook")
	_, ok := notifier.(*WebhookNotifier)
	assert.True(t, ok, "should return WebhookNotifier when webhook URL is provided")
}

func TestWebhookNotifierSendsEvent(t *testing.T) {
	t.Parallel()

	var receivedEvent BookingEvent
	var wg sync.WaitGroup
	wg.Add(1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "SitHub/1.0", r.Header.Get("User-Agent"))

		err := json.NewDecoder(r.Body).Decode(&receivedEvent)
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := NewWebhookNotifier(server.URL)
	event := BookingEvent{
		Event:       EventBookingCreated,
		BookingID:   "booking-123",
		DeskID:      "desk-1",
		UserID:      "user-1",
		UserName:    "Test User",
		BookingDate: "2026-01-20",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	notifier.NotifyAsync(&event)

	// Wait for the async call to complete
	wg.Wait()

	assert.Equal(t, EventBookingCreated, receivedEvent.Event)
	assert.Equal(t, "booking-123", receivedEvent.BookingID)
	assert.Equal(t, "desk-1", receivedEvent.DeskID)
	assert.Equal(t, "user-1", receivedEvent.UserID)
	assert.Equal(t, "Test User", receivedEvent.UserName)
}

func TestWebhookNotifierSendsCancelEvent(t *testing.T) {
	t.Parallel()

	var receivedEvent BookingEvent
	var wg sync.WaitGroup
	wg.Add(1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		require.NoError(t, json.NewDecoder(r.Body).Decode(&receivedEvent))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := NewWebhookNotifier(server.URL)
	event := BookingEvent{
		Event:            EventBookingCanceled,
		BookingID:        "booking-123",
		DeskID:           "desk-1",
		UserID:           "user-1",
		UserName:         "Test User",
		BookingDate:      "2026-01-20",
		CanceledByUserID: "admin-1",
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
	}

	notifier.NotifyAsync(&event)
	wg.Wait()

	assert.Equal(t, EventBookingCanceled, receivedEvent.Event)
	assert.Equal(t, "admin-1", receivedEvent.CanceledByUserID)
}

func TestWebhookNotifierHandlesServerError(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	wg.Add(1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		defer wg.Done()
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	notifier := NewWebhookNotifier(server.URL)

	// Should not panic even when server returns error
	notifier.NotifyAsync(&BookingEvent{
		Event:     EventBookingCreated,
		BookingID: "test-booking",
	})

	wg.Wait()
	// No assertion needed - just verify it doesn't panic
}

func TestWebhookNotifierHandlesInvalidURL(t *testing.T) {
	t.Parallel()

	notifier := NewWebhookNotifier("http://invalid-host-that-does-not-exist:9999")

	// Should not panic even when connection fails
	notifier.NotifyAsync(&BookingEvent{
		Event:     EventBookingCreated,
		BookingID: "test-booking",
	})

	// Give it a moment to attempt the request
	time.Sleep(100 * time.Millisecond)
	// No assertion needed - just verify it doesn't panic
}
