import { ref, watch, onScopeDispose } from 'vue';
import { useTheme } from 'vuetify';
import { getSafeLocalStorage } from './storage';

export type ThemePreference = 'auto' | 'light' | 'dark';

const STORAGE_KEY = 'sithub_theme';

function readStoredPreference(storage: Storage | null): ThemePreference {
  const stored = storage?.getItem(STORAGE_KEY);
  if (stored === 'light' || stored === 'dark') return stored;
  return 'auto';
}

/**
 * Composable providing theme preference management.
 * Persists to localStorage and reacts to OS preference changes in auto mode.
 */
export function useThemePreference() {
  const theme = useTheme();
  const storage = getSafeLocalStorage();
  const preference = ref<ThemePreference>(readStoredPreference(storage));

  let mediaQuery: MediaQueryList | null = null;
  let mediaHandler: ((e: MediaQueryListEvent) => void) | null = null;

  const applyThemeName = (name: 'light' | 'dark') => {
    const change = (theme as unknown as { change?: (themeName: string) => void }).change
      ?? (theme.global as unknown as { change?: (themeName: string) => void }).change;

    if (typeof change === 'function') {
      change(name);
      return;
    }

    (theme.global as { name: { value: string } }).name.value = name;
  };

  function applyTheme(pref: ThemePreference) {
    cleanupMediaListener();

    if (pref === 'auto') {
      if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') {
        applyThemeName('light');
        return;
      }
      mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      applyThemeName(mediaQuery.matches ? 'dark' : 'light');
      mediaHandler = (e: MediaQueryListEvent) => {
        applyThemeName(e.matches ? 'dark' : 'light');
      };
      mediaQuery.addEventListener('change', mediaHandler);
    } else {
      applyThemeName(pref);
    }
  }

  function cleanupMediaListener() {
    if (mediaQuery && mediaHandler) {
      mediaQuery.removeEventListener('change', mediaHandler);
      mediaQuery = null;
      mediaHandler = null;
    }
  }

  function setPreference(pref: ThemePreference) {
    preference.value = pref;
  }

  watch(preference, (pref) => {
    if (storage) {
      storage.setItem(STORAGE_KEY, pref);
    }
    applyTheme(pref);
  });

  // Apply immediately on creation
  applyTheme(preference.value);

  onScopeDispose(cleanupMediaListener);

  return {
    preference,
    setPreference
  };
}
