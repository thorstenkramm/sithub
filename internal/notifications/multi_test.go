package notifications

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type recordingNotifier struct {
	mu     sync.Mutex
	events []*BookingEvent
}

func (r *recordingNotifier) NotifyAsync(event *BookingEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
}

func TestMultiNotifierFansOut(t *testing.T) {
	t.Parallel()

	a := &recordingNotifier{}
	b := &recordingNotifier{}

	multi := MultiNotifier{a, nil, b}
	event := &BookingEvent{Event: EventBookingCreated, BookingID: "b1"}
	multi.NotifyAsync(event)

	assert.Equal(t, []*BookingEvent{event}, a.events)
	assert.Equal(t, []*BookingEvent{event}, b.events)
}

func TestMultiNotifierEmptyDoesNothing(t *testing.T) {
	t.Parallel()

	var multi MultiNotifier
	// Must not panic on a nil/empty slice.
	multi.NotifyAsync(&BookingEvent{Event: EventBookingCreated})
}
