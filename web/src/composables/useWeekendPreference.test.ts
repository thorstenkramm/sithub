import { useWeekendPreference } from './useWeekendPreference';

const STORAGE_KEY = 'sithub_show_weekends';

describe('useWeekendPreference', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('defaults to false when no stored value', () => {
    const { showWeekends } = useWeekendPreference();
    expect(showWeekends.value).toBe(false);
  });

  it('reads true from localStorage', () => {
    localStorage.setItem(STORAGE_KEY, 'true');
    const { showWeekends } = useWeekendPreference();
    expect(showWeekends.value).toBe(true);
  });

  it('reads false from localStorage', () => {
    localStorage.setItem(STORAGE_KEY, 'false');
    const { showWeekends } = useWeekendPreference();
    expect(showWeekends.value).toBe(false);
  });

  it('persists changes to localStorage', async () => {
    const { showWeekends } = useWeekendPreference();
    showWeekends.value = true;

    // Vue watchers are async - flush
    await new Promise(r => setTimeout(r, 0));
    expect(localStorage.getItem(STORAGE_KEY)).toBe('true');

    showWeekends.value = false;
    await new Promise(r => setTimeout(r, 0));
    expect(localStorage.getItem(STORAGE_KEY)).toBe('false');
  });
});
