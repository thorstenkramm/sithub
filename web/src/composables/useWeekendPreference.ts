import { ref, watch } from 'vue';
import { getSafeLocalStorage } from './storage';

const STORAGE_KEY = 'sithub_show_weekends';

// Module-level shared ref so all callers share the same reactive instance.
const storage = getSafeLocalStorage();
const showWeekends = ref(storage?.getItem(STORAGE_KEY) === 'true');

watch(showWeekends, (val) => {
  if (storage) {
    storage.setItem(STORAGE_KEY, String(val));
  }
});

/**
 * Composable providing weekend preference for booking views.
 * Persists to localStorage, default is false (Mon-Fri only).
 * Uses a shared module-level ref so all consumers stay in sync.
 */
export function useWeekendPreference() {
  return { showWeekends };
}
