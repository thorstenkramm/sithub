import { createApp, defineComponent } from 'vue';
import { createVuetify } from 'vuetify';
import { useThemePreference } from './useThemePreference';

const STORAGE_KEY = 'sithub_theme';

function withVuetify<T>(fn: () => T): { result: T; unmount: () => void } {
  const vuetify = createVuetify({
    theme: {
      defaultTheme: 'light',
      themes: {
        light: { dark: false, colors: {} },
        dark: { dark: true, colors: {} }
      }
    }
  });

  let result: T;
  const comp = defineComponent({
    setup() {
      result = fn();
      return () => null;
    }
  });

  const app = createApp(comp);
  app.use(vuetify);
  const el = document.createElement('div');
  app.mount(el);
  return { result: result!, unmount: () => app.unmount() };
}

describe('useThemePreference', () => {
  let originalMatchMedia: typeof window.matchMedia;

  beforeEach(() => {
    localStorage.clear();
    originalMatchMedia = window.matchMedia;
    window.matchMedia = vi.fn().mockReturnValue({
      matches: false,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn()
    });
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
  });

  it.each([
    { stored: undefined, expected: 'auto', label: 'defaults to auto when no stored value' },
    { stored: 'light', expected: 'light', label: 'reads stored light preference' },
    { stored: 'dark', expected: 'dark', label: 'reads stored dark preference' },
    { stored: 'invalid', expected: 'auto', label: 'falls back to auto for invalid stored value' }
  ])('$label', ({ stored, expected }) => {
    if (stored !== undefined) localStorage.setItem(STORAGE_KEY, stored);
    const { result, unmount } = withVuetify(() => useThemePreference().preference.value);
    expect(result).toBe(expected);
    unmount();
  });

  it('persists preference via setPreference', async () => {
    const { unmount } = withVuetify(() => {
      const { setPreference } = useThemePreference();
      setPreference('dark');
    });

    await new Promise(r => setTimeout(r, 0));
    expect(localStorage.getItem(STORAGE_KEY)).toBe('dark');
    unmount();
  });

  it('uses matchMedia for auto mode', () => {
    const { unmount } = withVuetify(() => useThemePreference());
    expect(window.matchMedia).toHaveBeenCalledWith('(prefers-color-scheme: dark)');
    unmount();
  });
});
