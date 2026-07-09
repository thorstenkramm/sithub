import { ref } from 'vue';
import { getSafeLocalStorage } from './storage';
import type { AreaView } from './useAreaViewPreference';

const STORAGE_KEY = 'sithub_my_bookings_view';

/**
 * Persists the selected My Bookings view globally (not per area) in localStorage.
 * Differs from useAreaViewPreference in two ways: it uses a single flat key and
 * defaults to 'table' on desktop when nothing is stored. Mobile always returns 'cards'.
 */
export function useMyBookingsViewPreference() {
  const activeView = ref<AreaView>('table');

  function load(isDesktop: boolean): AreaView {
    if (!isDesktop) {
      activeView.value = 'cards';
      return 'cards';
    }
    const storage = getSafeLocalStorage();
    if (!storage) {
      activeView.value = 'table';
      return 'table';
    }
    try {
      const raw = storage.getItem(STORAGE_KEY);
      if (raw === 'cards') {
        activeView.value = 'cards';
        return 'cards';
      }
    } catch {
      // Corrupted data — fall through to default
    }
    activeView.value = 'table';
    return 'table';
  }

  function save(view: AreaView) {
    activeView.value = view;
    const storage = getSafeLocalStorage();
    if (!storage) return;
    try {
      storage.setItem(STORAGE_KEY, view);
    } catch {
      // Storage full or unavailable
    }
  }

  return { activeView, load, save };
}
