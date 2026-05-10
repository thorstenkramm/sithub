import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { useLiveFeed } from './useLiveFeed';
import type { LiveBookingEvent } from '../api/liveFeed';

class StubSocket {
  public static instances: StubSocket[] = [];
  public sentClose = false;
  public onopen: ((ev?: unknown) => void) | null = null;
  public onmessage: ((ev: MessageEvent) => void) | null = null;
  public onerror: ((ev?: unknown) => void) | null = null;
  public onclose: ((ev?: unknown) => void) | null = null;

  constructor(public url: string) {
    StubSocket.instances.push(this);
  }

  open() {
    this.onopen?.();
  }

  recv(data: unknown) {
    this.onmessage?.(new MessageEvent('message', { data: JSON.stringify(data) }));
  }

  recvRaw(data: string) {
    this.onmessage?.(new MessageEvent('message', { data }));
  }

  serverClose() {
    this.onclose?.();
  }

  close() {
    this.sentClose = true;
    this.onclose?.();
  }
}

beforeEach(() => {
  StubSocket.instances = [];
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
});

const event: LiveBookingEvent = {
  type: 'booking.created',
  booking_id: 'b1',
  item_id: 'desk1',
  user_id: 'alice',
  booking_date: '2026-05-11',
  timestamp: '2026-05-10T12:00:00Z'
};

describe('useLiveFeed', () => {
  it('connects, dispatches booking events to subscribers', () => {
    const feed = useLiveFeed({
      socketFactory: (url) => new StubSocket(url) as unknown as WebSocket,
      locationProvider: () => ({ protocol: 'http:', host: 'localhost:5173' })
    });
    const handler = vi.fn();
    feed.onEvent(handler);
    feed.connect();

    const sock = StubSocket.instances[0];
    expect(sock.url).toBe('ws://localhost:5173/api/v1/live');
    sock.open();
    expect(feed.state.value).toBe('open');

    sock.recv(event);
    expect(handler).toHaveBeenCalledWith(event);
  });

  it('ignores invalid payloads', () => {
    const feed = useLiveFeed({
      socketFactory: (url) => new StubSocket(url) as unknown as WebSocket
    });
    const handler = vi.fn();
    feed.onEvent(handler);
    feed.connect();
    const sock = StubSocket.instances[0];
    sock.open();

    sock.recv({ type: 'unknown', booking_id: 'b' });
    sock.recvRaw('not-json');

    expect(handler).not.toHaveBeenCalled();
  });

  it('reconnects with exponential backoff and emits a reconnected event', () => {
    const factory = vi.fn((url: string) => new StubSocket(url) as unknown as WebSocket);
    const feed = useLiveFeed({
      socketFactory: factory,
      backoffBaseMs: 1000,
      maxDelayMs: 8000
    });
    const handler = vi.fn();
    feed.onEvent(handler);
    feed.connect();
    expect(factory).toHaveBeenCalledTimes(1);

    const first = StubSocket.instances[0];
    first.open();
    first.serverClose();
    expect(feed.state.value).toBe('closed');

    // Advance past base delay (with jitter ±20% the upper bound is 1200ms).
    vi.advanceTimersByTime(1500);
    expect(factory).toHaveBeenCalledTimes(2);

    const second = StubSocket.instances[1];
    second.open();

    // Reconnect must dispatch a synthetic reconnected event.
    expect(handler).toHaveBeenCalledWith({ type: 'reconnected' });
  });

  it('does not reconnect after disconnect', () => {
    const factory = vi.fn((url: string) => new StubSocket(url) as unknown as WebSocket);
    const feed = useLiveFeed({ socketFactory: factory });
    feed.connect();
    const sock = StubSocket.instances[0];
    sock.open();

    feed.disconnect();
    expect(sock.sentClose).toBe(true);
    expect(feed.state.value).toBe('idle');

    vi.advanceTimersByTime(60_000);
    expect(factory).toHaveBeenCalledTimes(1);
  });

  it('subscriber unsubscribe stops further dispatch', () => {
    const feed = useLiveFeed({
      socketFactory: (url) => new StubSocket(url) as unknown as WebSocket
    });
    const handler = vi.fn();
    const off = feed.onEvent(handler);
    feed.connect();
    const sock = StubSocket.instances[0];
    sock.open();

    sock.recv(event);
    expect(handler).toHaveBeenCalledTimes(1);

    off();
    sock.recv(event);
    expect(handler).toHaveBeenCalledTimes(1);
  });

  it('isolates a throwing subscriber from others', () => {
    const feed = useLiveFeed({
      socketFactory: (url) => new StubSocket(url) as unknown as WebSocket
    });
    const failing = vi.fn(() => { throw new Error('boom'); });
    const ok = vi.fn();
    feed.onEvent(failing);
    feed.onEvent(ok);
    feed.connect();
    const sock = StubSocket.instances[0];
    sock.open();

    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    sock.recv(event);

    expect(failing).toHaveBeenCalled();
    expect(ok).toHaveBeenCalledWith(event);
    consoleSpy.mockRestore();
  });
});
