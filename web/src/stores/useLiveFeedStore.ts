import { defineStore } from 'pinia';

import { useLiveFeed, type UseLiveFeedOptions } from '../composables/useLiveFeed';
import type { LiveEvent, LiveEventHandler } from '../api/liveFeed';

/**
 * useLiveFeedStore wraps the singleton WebSocket connection for /api/v1/live
 * and exposes a small pub/sub API. The store is started by App.vue once the
 * user is authenticated and stopped on logout. Components subscribe to events
 * via `subscribe(handler)` and unsubscribe by calling the returned disposer.
 *
 * The underlying composable is created lazily on first start() so unit tests
 * can install a custom socket factory by passing options.
 */
export const useLiveFeedStore = defineStore('liveFeed', () => {
  const handlers = new Set<LiveEventHandler>();
  let feed: ReturnType<typeof useLiveFeed> | null = null;

  function ensureFeed(options?: UseLiveFeedOptions): ReturnType<typeof useLiveFeed> {
    if (feed) return feed;
    feed = useLiveFeed(options);
    feed.onEvent((event: LiveEvent) => {
      for (const handler of handlers) {
        try {
          handler(event);
        } catch (err) {
           
          console.error('liveFeed store subscriber error', err);
        }
      }
    });
    return feed;
  }

  function start(options?: UseLiveFeedOptions) {
    ensureFeed(options).connect();
  }

  function stop() {
    if (feed) {
      feed.disconnect();
    }
  }

  function subscribe(handler: LiveEventHandler): () => void {
    handlers.add(handler);
    return () => {
      handlers.delete(handler);
    };
  }

  function reset() {
    stop();
    handlers.clear();
    feed = null;
  }

  return { start, stop, subscribe, reset };
});
