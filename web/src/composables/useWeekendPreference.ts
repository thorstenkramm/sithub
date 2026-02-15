import { ref, watch } from 'vue';
import { getSafeLocalStorage } from './storage';

const STORAGE_KEY = 'sithub_show_weekends';

/**
 * Composable providing weekend preference for booking views.
 * Persists to localStorage, default is false (Mon-Fri only).
 */
export function useWeekendPreference() {
  const storage = getSafeLocalStorage();
  const showWeekends = ref(storage?.getItem(STORAGE_KEY) === 'true');

  watch(showWeekends, (val) => {
    if (storage) {
      storage.setItem(STORAGE_KEY, String(val));
    }
  });

  return { showWeekends };
}
