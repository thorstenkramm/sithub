package livefeed

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// clientSendBuffer bounds outbound events per client. Beyond this, the
	// hub drops the client to protect itself.
	clientSendBuffer = 16

	// pingInterval determines how often we ping idle clients.
	pingInterval = 30 * time.Second

	// writeWait is the deadline for a single websocket write.
	writeWait = 10 * time.Second

	// pongWait is how long we wait for a pong in response to a ping. After
	// this with no read, the read pump exits and triggers cleanup.
	pongWait = 60 * time.Second
)

// Client wraps a single WebSocket connection. The hub never writes to the
// connection directly; it pushes events onto the send channel which the
// per-client write pump drains. Writes from multiple goroutines are not safe
// in gorilla/websocket, so the write pump is the single writer.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan Event
	userID string

	closed sync.Once
}

// newClient constructs a client bound to a websocket connection.
func newClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan Event, clientSendBuffer),
		userID: userID,
	}
}

// close releases the connection and the send channel exactly once.
func (c *Client) close() {
	c.closed.Do(func() {
		close(c.send)
		if c.conn != nil {
			_ = c.conn.Close() //nolint:errcheck // best-effort cleanup
		}
	})
}

// readPump consumes incoming frames purely to keep the connection responsive
// (pong frames extend the read deadline). The frontend is not expected to
// send anything; any received text/binary frame is read and discarded. When
// reading fails (close, timeout, network error) the client is unregistered.
func (c *Client) readPump() {
	defer c.hub.Unregister(c)

	c.conn.SetReadLimit(512)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait)) //nolint:errcheck // deadline-only
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			return
		}
	}
}

// writePump owns all writes to the websocket connection: incoming events from
// the hub and periodic pings. Exits when the send channel is closed (hub
// shutdown or client drop) or when a write fails.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close() //nolint:errcheck // best-effort cleanup
	}()

	for {
		select {
		case ev, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //nolint:errcheck // deadline-only
			if !ok {
				// Hub closed the channel; signal a clean close to the peer.
				//nolint:errcheck // best-effort close frame
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(ev); err != nil {
				slog.Debug("livefeed write failed", "user_id", c.userID, "err", err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //nolint:errcheck // deadline-only
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
