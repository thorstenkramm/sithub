import { ref, computed } from 'vue';
import { getSafeLocalStorage } from './storage';

const STORAGE_KEY = 'sithub_saved_filters';

export function useSavedFilters() {
  const storage = getSafeLocalStorage();

  const savedFilters = ref<string[]>(loadFromStorage());

  function loadFromStorage(): string[] {
    if (!storage) return [];
    try {
      const raw = storage.getItem(STORAGE_KEY);
      if (!raw) return [];
      const parsed = JSON.parse(raw);
      return Array.isArray(parsed) ? parsed.filter((f): f is string => typeof f === 'string') : [];
    } catch {
      return [];
    }
  }

  function persist() {
    if (!storage) return;
    storage.setItem(STORAGE_KEY, JSON.stringify(savedFilters.value));
  }

  function saveFilter(filter: string) {
    const trimmed = filter.trim();
    if (!trimmed) return false;
    if (savedFilters.value.includes(trimmed)) return false;
    savedFilters.value = [...savedFilters.value, trimmed];
    persist();
    return true;
  }

  function deleteFilter(filter: string) {
    savedFilters.value = savedFilters.value.filter(f => f !== filter);
    persist();
  }

  const isSavedFilter = (filter: string) =>
    savedFilters.value.includes(filter.trim());

  const comboboxItems = computed(() => savedFilters.value);

  return {
    savedFilters,
    comboboxItems,
    saveFilter,
    deleteFilter,
    isSavedFilter
  };
}
