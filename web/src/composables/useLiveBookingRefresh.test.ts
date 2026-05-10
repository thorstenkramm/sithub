import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { defineComponent, h } from 'vue';

import { useLiveBookingRefresh } from './useLiveBookingRefresh';
import { useLiveFeedStore } from '../stores/useLiveFeedStore';
import { useAuthStore } from '../stores/useAuthStore';
import type { LiveBookingEvent, LiveEventHandler } from '../api/liveFeed';

let lastSubscribeHandler: LiveEventHandler | null = null;

function makeWrapper(refresh: () => void, opts: { isRelevant?: (e: LiveBookingEvent) => boolean } = {}) {
  const Test = defineComponent({
    setup() {
      useLiveBookingRefresh({
        refresh,
        isRelevant: opts.isRelevant
      });
      return () => h('div');
    }
  });
  return mount(Test);
}

beforeEach(() => {
  setActivePinia(createPinia());
  vi.useFakeTimers();
  lastSubscribeHandler = null;

  // Stub the live-feed store's subscribe to capture the handler instead of
  // wiring up the real WebSocket composable.
  const liveFeed = useLiveFeedStore();
  liveFeed.subscribe = (handler: LiveEventHandler) => {
    lastSubscribeHandler = handler;
    return () => {
      lastSubscribeHandler = null;
    };
  };
});

afterEach(() => {
  vi.useRealTimers();
});

const otherUserBooking: LiveBookingEvent = {
  type: 'booking.created',
  booking_id: 'b1',
  item_id: 'desk1',
  user_id: 'bob',
  booking_date: '2026-05-11',
  timestamp: '2026-05-10T12:00:00Z'
};

describe('useLiveBookingRefresh', () => {
  it('debounces booking events into one refresh call', () => {
    const refresh = vi.fn();
    makeWrapper(refresh);

    expect(lastSubscribeHandler).toBeTypeOf('function');
    lastSubscribeHandler!(otherUserBooking);
    lastSubscribeHandler!(otherUserBooking);
    lastSubscribeHandler!(otherUserBooking);

    expect(refresh).not.toHaveBeenCalled();
    vi.advanceTimersByTime(300);
    expect(refresh).toHaveBeenCalledTimes(1);
  });

  it('skips self-events', () => {
    const auth = useAuthStore();
    auth.setUser({
      id: 'alice',
      display_name: 'Alice',
      email: 'a@example.com',
      is_admin: false,
      auth_source: 'internal'
    });

    const refresh = vi.fn();
    makeWrapper(refresh);

    lastSubscribeHandler!({ ...otherUserBooking, user_id: 'alice' });
    vi.advanceTimersByTime(300);
    expect(refresh).not.toHaveBeenCalled();
  });

  it('respects isRelevant filter', () => {
    const refresh = vi.fn();
    const isRelevant = vi.fn().mockReturnValue(false);
    makeWrapper(refresh, { isRelevant });

    lastSubscribeHandler!(otherUserBooking);
    vi.advanceTimersByTime(300);
    expect(isRelevant).toHaveBeenCalledWith(otherUserBooking);
    expect(refresh).not.toHaveBeenCalled();
  });

  it('refreshes on reconnected event without checking isRelevant', () => {
    const refresh = vi.fn();
    const isRelevant = vi.fn().mockReturnValue(false);
    makeWrapper(refresh, { isRelevant });

    lastSubscribeHandler!({ type: 'reconnected' });
    vi.advanceTimersByTime(300);
    expect(refresh).toHaveBeenCalledTimes(1);
  });

  it('unsubscribes on unmount', () => {
    const refresh = vi.fn();
    const wrapper = makeWrapper(refresh);
    expect(lastSubscribeHandler).toBeTypeOf('function');
    wrapper.unmount();
    expect(lastSubscribeHandler).toBeNull();
  });
});
