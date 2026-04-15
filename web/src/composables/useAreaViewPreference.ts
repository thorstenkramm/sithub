import { ref } from 'vue';
import { getSafeLocalStorage } from './storage';

const STORAGE_KEY = 'sithub_area_view';

export type AreaView = 'cards' | 'table';

/**
 * Persists the selected area view per area ID in localStorage.
 * Returns 'cards' as the default when no preference is stored.
 */
export function useAreaViewPreference() {
  const activeView = ref<AreaView>('cards');

  function load(areaId: string, isDesktop: boolean): AreaView {
    if (!isDesktop) {
      activeView.value = 'cards';
      return 'cards';
    }
    const storage = getSafeLocalStorage();
    if (!storage) {
      activeView.value = 'cards';
      return 'cards';
    }
    try {
      const raw = storage.getItem(STORAGE_KEY);
      if (raw) {
        const prefs = JSON.parse(raw) as Record<string, string>;
        if (prefs[areaId] === 'table') {
          activeView.value = 'table';
          return 'table';
        }
      }
    } catch {
      // Corrupted data — fall through to default
    }
    activeView.value = 'cards';
    return 'cards';
  }

  function save(areaId: string, view: AreaView) {
    activeView.value = view;
    const storage = getSafeLocalStorage();
    if (!storage) return;
    try {
      const raw = storage.getItem(STORAGE_KEY);
      const prefs: Record<string, string> = raw ? JSON.parse(raw) : {};
      if (view === 'cards') {
        delete prefs[areaId];
      } else {
        prefs[areaId] = view;
      }
      storage.setItem(STORAGE_KEY, JSON.stringify(prefs));
    } catch {
      // Storage full or unavailable
    }
  }

  return { activeView, load, save };
}
