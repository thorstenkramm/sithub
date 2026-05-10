// Types for the WebSocket live feed at /api/v1/live.
// Mirrors internal/livefeed/event.go on the backend.

export type LiveBookingEventType = 'booking.created' | 'booking.canceled';

export type LiveEventType = LiveBookingEventType | 'reconnected';

export interface LiveBookingEvent {
  type: LiveBookingEventType;
  booking_id: string;
  item_id: string;
  user_id: string;
  booking_date: string;
  timestamp: string;
}

/** Synthetic event fired on a successful reconnect so listeners can refetch. */
export interface LiveReconnectedEvent {
  type: 'reconnected';
}

export type LiveEvent = LiveBookingEvent | LiveReconnectedEvent;

export type LiveEventHandler = (event: LiveEvent) => void;

/** Build the absolute WebSocket URL for the live feed from window.location. */
export function liveFeedUrl(loc: { protocol: string; host: string }): string {
  const scheme = loc.protocol === 'https:' ? 'wss' : 'ws';
  return `${scheme}://${loc.host}/api/v1/live`;
}

/** Type guard for booking events (excludes the synthetic reconnect event). */
export function isBookingEvent(event: LiveEvent): event is LiveBookingEvent {
  return event.type === 'booking.created' || event.type === 'booking.canceled';
}
