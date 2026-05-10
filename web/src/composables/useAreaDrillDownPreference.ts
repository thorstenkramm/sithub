import { ref } from 'vue';
import { getSafeLocalStorage } from './storage';

const STORAGE_KEY = 'sithub_area_drill_down';

/**
 * Persists the "Area drill-down" toggle in local storage. When no value is stored,
 * the default is on for compact viewports and off for desktop, matching the user's
 * intent that desktop users skip the drill-down step by default.
 */
export function useAreaDrillDownPreference() {
  const enabled = ref(true);
  const hasUserChoice = ref(false);

  function load(isLargeScreen: boolean) {
    const storage = getSafeLocalStorage();
    if (storage) {
      try {
        const raw = storage.getItem(STORAGE_KEY);
        if (raw === 'on' || raw === 'off') {
          enabled.value = raw === 'on';
          hasUserChoice.value = true;
          return;
        }
      } catch {
        // Corrupted data — fall through to default
      }
    }
    enabled.value = !isLargeScreen;
    hasUserChoice.value = false;
  }

  function set(value: boolean) {
    enabled.value = value;
    hasUserChoice.value = true;
    const storage = getSafeLocalStorage();
    if (!storage) return;
    try {
      storage.setItem(STORAGE_KEY, value ? 'on' : 'off');
    } catch {
      // Storage full or unavailable
    }
  }

  return { enabled, hasUserChoice, load, set };
}
