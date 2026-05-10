package livefeed

import (
	"context"
	"log/slog"
	"sync/atomic"

	"github.com/thorstenkramm/sithub/internal/notifications"
)

// broadcastBuffer is the size of the hub's incoming event queue. A few dozen
// in-flight events is plenty for a single-node deployment; if the hub falls
// behind, NotifyAsync drops with a warning rather than blocking the booking
// handler.
const broadcastBuffer = 64

// Hub fans out booking events to all connected clients.
//
// Hub is safe for concurrent use after Run has been started. NotifyAsync,
// Register, and Unregister can be called from any goroutine.
type Hub struct {
	register   chan *Client
	unregister chan *Client
	broadcast  chan Event
	clients    map[*Client]struct{}
	done       chan struct{}
	running    atomic.Bool
}

// NewHub creates an unstarted hub. Call Run in a goroutine to start it.
func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Event, broadcastBuffer),
		clients:    make(map[*Client]struct{}),
		done:       make(chan struct{}),
	}
}

// Run owns the hub's event loop. It blocks until ctx is canceled, then closes
// every connected client. Returning is the signal that the hub has shut down;
// after that, Register returns false and NotifyAsync drops events silently.
func (h *Hub) Run(ctx context.Context) {
	h.running.Store(true)
	slog.Info("livefeed hub started")
	defer func() {
		for c := range h.clients {
			c.close()
		}
		h.clients = nil
		h.running.Store(false)
		close(h.done)
		slog.Info("livefeed hub stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case c := <-h.register:
			h.clients[c] = struct{}{}
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				c.close()
			}
		case ev := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- ev:
				default:
					// Slow client: drop it rather than blocking the hub.
					slog.Warn("livefeed dropping slow client", "user_id", c.userID)
					delete(h.clients, c)
					c.close()
				}
			}
		}
	}
}

// Register adds a client to the hub. Returns false if the hub has shut down.
func (h *Hub) Register(c *Client) bool {
	select {
	case h.register <- c:
		return true
	case <-h.done:
		return false
	}
}

// Unregister removes a client. Safe to call multiple times. The hub closes the
// client when it processes the request.
func (h *Hub) Unregister(c *Client) {
	select {
	case h.unregister <- c:
	case <-h.done:
	}
}

// NotifyAsync implements notifications.Notifier. It is non-blocking: if the
// hub's broadcast queue is full or the hub has shut down, the event is
// dropped with a warning. Booking handlers must never wait on the hub.
func (h *Hub) NotifyAsync(event *notifications.BookingEvent) {
	if event == nil || !h.running.Load() {
		return
	}
	ev := fromBookingEvent(event)
	select {
	case h.broadcast <- ev:
	default:
		slog.Warn("livefeed broadcast queue full, dropping event",
			"event", ev.Type,
			"booking_id", ev.BookingID,
		)
	}
}
