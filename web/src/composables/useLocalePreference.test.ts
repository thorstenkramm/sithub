import { useLocalePreference } from './useLocalePreference';
import type { LocalePreference } from './useLocalePreference';
import { i18n } from '../plugins/i18n';

const STORAGE_KEY = 'sithub_locale';

function getI18nLocale(): string {
  return (i18n.global.locale as unknown as { value: string }).value;
}

describe('useLocalePreference', () => {
  beforeEach(() => {
    localStorage.clear();
    // Reset to known state
    const { setPreference } = useLocalePreference();
    setPreference('auto');
  });

  it('defaults to auto when no stored preference', () => {
    const { preference } = useLocalePreference();
    expect(preference.value).toBe('auto');
  });

  it('sets explicit locale preference', () => {
    const { setPreference, preference } = useLocalePreference();
    setPreference('de');
    expect(preference.value).toBe('de');
  });

  it('sets i18n locale when preference changes', () => {
    const { setPreference } = useLocalePreference();
    setPreference('fr');
    expect(getI18nLocale()).toBe('fr');
  });

  it('persists preference to localStorage', () => {
    const { setPreference } = useLocalePreference();
    setPreference('es');
    expect(localStorage.getItem(STORAGE_KEY)).toBe('es');
  });

  it('resolves auto to English when browser language is unsupported', () => {
    // jsdom navigator.language defaults to empty or en-US
    const { setPreference } = useLocalePreference();
    setPreference('de');
    expect(getI18nLocale()).toBe('de');
    setPreference('auto');
    // In jsdom, navigator.languages is typically empty or ['en-US']
    // so auto should resolve to 'en'
    expect(getI18nLocale()).toBe('en');
  });

  it('accepts all supported locale values', () => {
    const { setPreference, preference } = useLocalePreference();
    const locales: LocalePreference[] = ['en', 'de', 'es', 'fr', 'uk', 'auto'];
    for (const locale of locales) {
      setPreference(locale);
      expect(preference.value).toBe(locale);
    }
  });

  it('switches back to explicit locale after auto', () => {
    const { setPreference } = useLocalePreference();
    setPreference('auto');
    setPreference('uk');
    expect(getI18nLocale()).toBe('uk');
    expect(localStorage.getItem(STORAGE_KEY)).toBe('uk');
  });
});
