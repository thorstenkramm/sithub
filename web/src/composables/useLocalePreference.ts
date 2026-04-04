import { ref, watch } from 'vue';
import { i18n } from '../plugins/i18n';
import { getSafeLocalStorage } from './storage';

export type SupportedLocale = 'en' | 'de' | 'es' | 'fr' | 'uk';
export type LocalePreference = 'auto' | SupportedLocale;

const STORAGE_KEY = 'sithub_locale';
const SUPPORTED_LOCALES: SupportedLocale[] = ['en', 'de', 'es', 'fr', 'uk'];

function detectBrowserLocale(): SupportedLocale {
  if (typeof navigator === 'undefined') return 'en';
  const languages = navigator.languages ?? [navigator.language];
  for (const lang of languages) {
    const code = (lang.split('-')[0] ?? '').toLowerCase();
    if (SUPPORTED_LOCALES.includes(code as SupportedLocale)) {
      return code as SupportedLocale;
    }
  }
  return 'en';
}

function resolveLocale(pref: LocalePreference): SupportedLocale {
  return pref === 'auto' ? detectBrowserLocale() : pref;
}

function readStoredPreference(storage: Storage | null): LocalePreference {
  const stored = storage?.getItem(STORAGE_KEY);
  if (stored && SUPPORTED_LOCALES.includes(stored as SupportedLocale)) {
    return stored as SupportedLocale;
  }
  if (stored === 'auto') return 'auto';
  return 'auto';
}

function applyLocale(locale: SupportedLocale) {
  (i18n.global.locale as unknown as { value: string }).value = locale;
}

// Module-level shared state so all components see the same preference
const storage = getSafeLocalStorage();
const preference = ref<LocalePreference>(readStoredPreference(storage));

// Apply immediately on module load
applyLocale(resolveLocale(preference.value));

watch(preference, (pref) => {
  if (storage) {
    storage.setItem(STORAGE_KEY, pref);
  }
  applyLocale(resolveLocale(pref));
}, { flush: 'sync' });

/**
 * Composable providing locale preference management.
 * Persists to localStorage and supports auto-detection from browser language.
 */
export function useLocalePreference() {
  function setPreference(pref: LocalePreference) {
    preference.value = pref;
  }

  return {
    preference,
    setPreference
  };
}
