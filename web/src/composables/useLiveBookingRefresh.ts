import { onBeforeUnmount, onMounted } from 'vue';

import { useLiveFeedStore } from '../stores/useLiveFeedStore';
import { useAuthStore } from '../stores/useAuthStore';
import { isBookingEvent, type LiveBookingEvent, type LiveEvent } from '../api/liveFeed';

export interface LiveBookingRefreshOptions {
  /**
   * Called when an external event suggests the view's data is stale. The
   * implementation should refetch whatever the view renders. Multiple events
   * arriving within `debounceMs` collapse into a single call.
   */
  refresh: () => void | Promise<void>;
  /**
   * Optional filter: return true to consider this event relevant to the
   * current view. Self-events (events triggered by the current user) are
   * always filtered out before this is called.
   *
   * Defaults to "all booking events are relevant".
   */
  isRelevant?: (event: LiveBookingEvent) => boolean;
  /**
   * Coalesce events arriving close together into a single refresh call.
   * Defaults to 250ms.
   */
  debounceMs?: number;
}

/**
 * useLiveBookingRefresh subscribes to the live feed store on mount and calls
 * the given refresh function when a relevant booking event arrives. It
 * unsubscribes on unmount.
 *
 * Self-events (events whose `user_id` matches the authenticated user) are
 * filtered out so the user does not see flicker from their own actions.
 *
 * The synthetic `'reconnected'` event always triggers a refresh: after a
 * dropout, any number of events may have been missed and the safest answer
 * is to refetch.
 */
export function useLiveBookingRefresh(options: LiveBookingRefreshOptions) {
  const liveFeed = useLiveFeedStore();
  const authStore = useAuthStore();

  const debounceMs = options.debounceMs ?? 250;
  const isRelevant = options.isRelevant ?? (() => true);

  let pending: ReturnType<typeof setTimeout> | null = null;
  let unsubscribe: (() => void) | null = null;

  function scheduleRefresh() {
    if (pending !== null) return;
    pending = setTimeout(async () => {
      pending = null;
      try {
        await options.refresh();
      } catch (err) {
         
        console.error('useLiveBookingRefresh refresh threw', err);
      }
    }, debounceMs);
  }

  function handleEvent(event: LiveEvent) {
    if (event.type === 'reconnected') {
      scheduleRefresh();
      return;
    }
    if (!isBookingEvent(event)) return;
    if (event.user_id && event.user_id === authStore.userId) return;
    if (!isRelevant(event)) return;
    scheduleRefresh();
  }

  onMounted(() => {
    unsubscribe = liveFeed.subscribe(handleEvent);
  });

  onBeforeUnmount(() => {
    if (unsubscribe) {
      unsubscribe();
      unsubscribe = null;
    }
    if (pending !== null) {
      clearTimeout(pending);
      pending = null;
    }
  });
}
