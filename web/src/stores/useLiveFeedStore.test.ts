import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ref } from 'vue';
import { createPinia, setActivePinia } from 'pinia';

import { useLiveFeedStore } from './useLiveFeedStore';
import { useLiveFeed } from '../composables/useLiveFeed';

vi.mock('../composables/useLiveFeed', () => ({
  useLiveFeed: vi.fn()
}));

describe('useLiveFeedStore', () => {
  const connect = vi.fn();
  const disconnect = vi.fn();
  const onEvent = vi.fn();
  let feedHandler: ((event: { type: string }) => void) | null = null;

  beforeEach(() => {
    setActivePinia(createPinia());
    vi.mocked(useLiveFeed).mockReset();
    connect.mockReset();
    disconnect.mockReset();
    onEvent.mockReset();
    feedHandler = null;

    vi.mocked(useLiveFeed).mockImplementation(() => {
      onEvent.mockImplementation((handler: (event: { type: string }) => void) => {
        feedHandler = handler;
        return () => {
          if (feedHandler === handler) {
            feedHandler = null;
          }
        };
      });
      return {
        state: ref('idle'),
        connect,
        disconnect,
        onEvent
      };
    });
  });

  it('reuses one feed instance until reset is called', () => {
    const store = useLiveFeedStore();

    store.start();
    store.start();

    expect(useLiveFeed).toHaveBeenCalledTimes(1);
    expect(connect).toHaveBeenCalledTimes(2);

    store.reset();
    store.start();

    expect(disconnect).toHaveBeenCalledTimes(1);
    expect(useLiveFeed).toHaveBeenCalledTimes(2);
  });

  it('fans out events to subscribers and unsubscribes cleanly', () => {
    const store = useLiveFeedStore();
    const first = vi.fn();
    const second = vi.fn();

    store.start();
    const unsubscribeFirst = store.subscribe(first);
    store.subscribe(second);

    expect(feedHandler).toBeTypeOf('function');
    feedHandler!({ type: 'reconnected' });
    expect(first).toHaveBeenCalledWith({ type: 'reconnected' });
    expect(second).toHaveBeenCalledWith({ type: 'reconnected' });

    unsubscribeFirst();
    feedHandler!({ type: 'reconnected' });
    expect(first).toHaveBeenCalledTimes(1);
    expect(second).toHaveBeenCalledTimes(2);
  });

  it('disconnects the active socket when stop is called', () => {
    const store = useLiveFeedStore();

    store.start();
    store.stop();

    expect(disconnect).toHaveBeenCalledTimes(1);
  });
});
