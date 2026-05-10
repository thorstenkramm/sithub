import { describe, expect, it } from 'vitest';

import { isBookingEvent, liveFeedUrl } from './liveFeed';

describe('liveFeedUrl', () => {
  it('uses ws:// for http origins', () => {
    expect(liveFeedUrl({ protocol: 'http:', host: 'localhost:5173' }))
      .toBe('ws://localhost:5173/api/v1/live');
  });

  it('uses wss:// for https origins', () => {
    expect(liveFeedUrl({ protocol: 'https:', host: 'sithub.example.com' }))
      .toBe('wss://sithub.example.com/api/v1/live');
  });
});

describe('isBookingEvent', () => {
  it('returns true for booking events', () => {
    expect(isBookingEvent({
      type: 'booking.created',
      booking_id: 'b1',
      item_id: 'i1',
      user_id: 'u1',
      booking_date: '2026-05-11',
      timestamp: '2026-05-10T12:00:00Z'
    })).toBe(true);
    expect(isBookingEvent({
      type: 'booking.canceled',
      booking_id: 'b1',
      item_id: 'i1',
      user_id: 'u1',
      booking_date: '2026-05-11',
      timestamp: '2026-05-10T12:00:00Z'
    })).toBe(true);
  });

  it('returns false for the synthetic reconnect event', () => {
    expect(isBookingEvent({ type: 'reconnected' })).toBe(false);
  });
});
