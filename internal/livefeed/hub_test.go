package livefeed

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/notifications"
)

const warmupBookingID = "warmup"

// dialClient spins up a tiny upgrade-only test server and dials it. The
// returned ws conn is the client side; the server side is registered with the
// hub.
func dialClient(t *testing.T, hub *Hub) *websocket.Conn {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c := newClient(hub, conn, "user-test")
		if !hub.Register(c) {
			_ = conn.Close() //nolint:errcheck // best-effort cleanup
			return
		}
		go c.writePump()
		c.readPump()
	}))
	t.Cleanup(srv.Close)

	wsURL := strings.Replace(srv.URL, "http://", "ws://", 1) + "/"
	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 2 * time.Second
	conn, resp, err := dialer.Dial(wsURL, nil)
	require.NoError(t, err)
	if resp != nil {
		require.NoError(t, resp.Body.Close())
	}
	t.Cleanup(func() { _ = conn.Close() }) //nolint:errcheck // best-effort cleanup
	return conn
}

// waitForRegistration blocks until the hub delivers events to conn. The hub
// registers clients asynchronously and silently drops events sent before the
// registration lands, so a background ticker keeps sending warmup events while
// the test blocks on a SINGLE read with a generous deadline.
// gorilla/websocket treats a timed-out read as FATAL for the connection, so
// polling with short read deadlines (e.g. ReadJSON inside Eventually) poisons
// the connection on its first miss and can never recover under -race on CI.
func waitForRegistration(t *testing.T, hub *Hub, conn *websocket.Conn) {
	t.Helper()

	stopWarmup := make(chan struct{})
	warmupDone := make(chan struct{})
	go func() {
		defer close(warmupDone)
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopWarmup:
				return
			case <-ticker.C:
				hub.NotifyAsync(&notifications.BookingEvent{
					Event:       notifications.EventBookingCreated,
					BookingID:   warmupBookingID,
					BookingDate: "2026-05-11",
					Timestamp:   "2026-05-10T12:00:00Z",
				})
			}
		}
	}()
	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second)) //nolint:errcheck // deadline-only
	var warm Event
	err := conn.ReadJSON(&warm)
	close(stopWarmup)
	<-warmupDone
	require.NoError(t, err, "websocket client never registered with hub")
}

// readNonWarmup reads events from conn, skipping leftover warmups queued by
// waitForRegistration. Events enqueued after waitForRegistration returns come
// after any leftover warmups, so the loop is guaranteed to terminate.
func readNonWarmup(t *testing.T, conn *websocket.Conn) Event {
	t.Helper()

	var got Event
	for {
		require.NoError(t, conn.ReadJSON(&got))
		if got.BookingID != warmupBookingID {
			return got
		}
	}
}

func TestHubBroadcastReachesClient(t *testing.T) {
	t.Parallel()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	conn := dialClient(t, hub)

	waitForRegistration(t, hub, conn)

	hub.NotifyAsync(&notifications.BookingEvent{
		Event:       notifications.EventBookingCreated,
		BookingID:   "b1",
		ItemID:      "desk1",
		UserID:      "alice",
		BookingDate: "2026-05-11",
		Timestamp:   "2026-05-10T12:00:00Z",
		GuestName:   "should-be-stripped",
		GuestEmail:  "should-be-stripped@example.com",
	})

	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second)) //nolint:errcheck // deadline-only
	got := readNonWarmup(t, conn)
	assert.Equal(t, EventBookingCreated, got.Type)
	assert.Equal(t, "b1", got.BookingID)
	assert.Equal(t, "desk1", got.ItemID)
	assert.Equal(t, "alice", got.UserID)
	assert.Equal(t, "2026-05-11", got.BookingDate)
	assert.Equal(t, "2026-05-10T12:00:00Z", got.Timestamp)
}

func TestHubExitsOnContextCancel(t *testing.T) {
	t.Parallel()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())

	stopped := make(chan struct{})
	go func() {
		hub.Run(ctx)
		close(stopped)
	}()

	cancel()

	select {
	case <-stopped:
	case <-time.After(5 * time.Second):
		t.Fatal("hub did not stop after context cancel")
	}

	// After shutdown, NotifyAsync must not panic and Register must return
	// false.
	hub.NotifyAsync(&notifications.BookingEvent{Event: notifications.EventBookingCreated})
	assert.False(t, hub.Register(&Client{}))
}

func TestHubDropsSlowClient(t *testing.T) {
	t.Parallel()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	// A client whose writePump never runs: its send buffer fills after
	// clientSendBuffer events, then the hub drops it on the next event.
	c := &Client{
		hub:    hub,
		send:   make(chan Event, clientSendBuffer),
		userID: "stuck",
	}
	require.True(t, hub.Register(c))

	for i := 0; i < clientSendBuffer*4; i++ {
		hub.NotifyAsync(&notifications.BookingEvent{
			Event:     notifications.EventBookingCreated,
			BookingID: "b",
		})
	}

	// The hub drops the client by closing its send channel. Wait for that.
	require.Eventually(t, func() bool {
		select {
		case _, ok := <-c.send:
			return !ok
		default:
			return false
		}
	}, 5*time.Second, 10*time.Millisecond, "expected hub to close slow client's send channel")
}

func TestHubConcurrentRegisterUnregister(t *testing.T) {
	t.Parallel()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			c := &Client{hub: hub, send: make(chan Event, clientSendBuffer)}
			if hub.Register(c) {
				// Drain so the hub never has to drop us.
				go func() {
					for range c.send {
					}
				}()
				hub.NotifyAsync(&notifications.BookingEvent{
					Event:     notifications.EventBookingCreated,
					BookingID: "b",
				})
				hub.Unregister(c)
			}
		}()
	}
	wg.Wait()
}

func TestNotifyAsyncDropsWhenQueueFull(t *testing.T) {
	t.Parallel()

	hub := NewHub()
	// Do not start Run: NotifyAsync should still not panic, but events go
	// nowhere. Setting running=true via Run start would let the queue fill;
	// we simulate "queue full" by manually flagging running and filling
	// the broadcast channel.
	hub.running.Store(true)
	for i := 0; i < broadcastBuffer; i++ {
		hub.broadcast <- Event{Type: EventBookingCreated}
	}

	// One more must not block; it is dropped.
	done := make(chan struct{})
	go func() {
		hub.NotifyAsync(&notifications.BookingEvent{Event: notifications.EventBookingCreated})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("NotifyAsync blocked when broadcast queue was full")
	}
}
