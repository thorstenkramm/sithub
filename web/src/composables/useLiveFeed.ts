import { ref, type Ref } from 'vue';
import {
  liveFeedUrl,
  type LiveEvent,
  type LiveEventHandler,
  type LiveBookingEvent
} from '../api/liveFeed';

export type LiveFeedState = 'idle' | 'connecting' | 'open' | 'closed';

export interface UseLiveFeedOptions {
  /**
   * Override the WebSocket constructor (used in tests).
   */
  socketFactory?: (url: string) => WebSocket;
  /**
   * Override the location source (used in tests).
   */
  locationProvider?: () => { protocol: string; host: string };
  /**
   * Backoff base in ms. The actual delay is `base * 2^attempt` plus ±20%
   * jitter, capped at `maxDelay`.
   */
  backoffBaseMs?: number;
  /**
   * Maximum reconnect delay in ms.
   */
  maxDelayMs?: number;
}

export interface UseLiveFeedReturn {
  state: Ref<LiveFeedState>;
  connect: () => void;
  disconnect: () => void;
  onEvent: (handler: LiveEventHandler) => () => void;
}

/**
 * useLiveFeed is the low-level WebSocket client for /api/v1/live. It owns one
 * connection at a time, reconnects with exponential backoff + jitter on close,
 * and dispatches `LiveEvent`s to subscribers.
 *
 * On a successful reconnect (after a previous open → close cycle) it dispatches
 * a synthetic `'reconnected'` event so subscribers can refresh the slice they
 * may have missed during the outage.
 */
export function useLiveFeed(options: UseLiveFeedOptions = {}): UseLiveFeedReturn {
  const state = ref<LiveFeedState>('idle');
  const handlers = new Set<LiveEventHandler>();

  const backoffBase = options.backoffBaseMs ?? 1000;
  const maxDelay = options.maxDelayMs ?? 30_000;

  let socket: WebSocket | null = null;
  let reconnectAttempt = 0;
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  let manuallyClosed = false;
  let hasConnectedAtLeastOnce = false;

  function buildUrl(): string {
    const provider = options.locationProvider ?? (() => ({
      protocol: window.location.protocol,
      host: window.location.host
    }));
    return liveFeedUrl(provider());
  }

  function clearReconnectTimer() {
    if (reconnectTimer !== null) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
  }

  function scheduleReconnect() {
    clearReconnectTimer();
    const cap = Math.min(maxDelay, backoffBase * Math.pow(2, reconnectAttempt));
    const jitter = cap * 0.2 * (Math.random() * 2 - 1);
    const delay = Math.max(0, cap + jitter);
    reconnectAttempt += 1;
    reconnectTimer = setTimeout(openSocket, delay);
  }

  function dispatch(event: LiveEvent) {
    for (const handler of handlers) {
      try {
        handler(event);
      } catch (err) {
        // A subscriber throwing must not break the dispatch loop.
         
        console.error('liveFeed subscriber error', err);
      }
    }
  }

  function openSocket() {
    clearReconnectTimer();
    state.value = 'connecting';

    const factory = options.socketFactory ?? ((url: string) => new WebSocket(url));
    let ws: WebSocket;
    try {
      ws = factory(buildUrl());
    } catch {
      state.value = 'closed';
      if (!manuallyClosed) {
        scheduleReconnect();
      }
      return;
    }

    socket = ws;

    ws.onopen = () => {
      reconnectAttempt = 0;
      state.value = 'open';
      if (hasConnectedAtLeastOnce) {
        dispatch({ type: 'reconnected' });
      }
      hasConnectedAtLeastOnce = true;
    };

    ws.onmessage = (msg: MessageEvent) => {
      let parsed: unknown;
      try {
        parsed = JSON.parse(typeof msg.data === 'string' ? msg.data : '');
      } catch {
        return;
      }
      if (!isLiveBookingPayload(parsed)) {
        return;
      }
      dispatch(parsed);
    };

    ws.onerror = () => {
      // The browser fires onerror before onclose on most failures. We rely
      // on onclose to drive the reconnect path.
    };

    ws.onclose = () => {
      socket = null;
      state.value = 'closed';
      if (!manuallyClosed) {
        scheduleReconnect();
      }
    };
  }

  function connect() {
    manuallyClosed = false;
    if (socket && (state.value === 'connecting' || state.value === 'open')) {
      return;
    }
    openSocket();
  }

  function disconnect() {
    manuallyClosed = true;
    clearReconnectTimer();
    reconnectAttempt = 0;
    hasConnectedAtLeastOnce = false;
    if (socket) {
      socket.close();
      socket = null;
    }
    state.value = 'idle';
  }

  function onEvent(handler: LiveEventHandler): () => void {
    handlers.add(handler);
    return () => {
      handlers.delete(handler);
    };
  }

  return { state, connect, disconnect, onEvent };
}

function isLiveBookingPayload(value: unknown): value is LiveBookingEvent {
  if (typeof value !== 'object' || value === null) return false;
  const v = value as Record<string, unknown>;
  if (v.type !== 'booking.created' && v.type !== 'booking.canceled') return false;
  return (
    typeof v.booking_id === 'string'
    && typeof v.item_id === 'string'
    && typeof v.user_id === 'string'
    && typeof v.booking_date === 'string'
    && typeof v.timestamp === 'string'
  );
}
