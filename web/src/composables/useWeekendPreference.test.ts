import { useWeekendPreference } from './useWeekendPreference';

const STORAGE_KEY = 'sithub_show_weekends';

describe('useWeekendPreference', () => {
  beforeEach(() => {
    localStorage.clear();
    // Reset the shared ref to default (false)
    const { showWeekends } = useWeekendPreference();
    showWeekends.value = false;
  });

  it('defaults to false when no stored value', () => {
    const { showWeekends } = useWeekendPreference();
    expect(showWeekends.value).toBe(false);
  });

  it('shares the same ref across multiple calls', () => {
    const a = useWeekendPreference();
    const b = useWeekendPreference();
    a.showWeekends.value = true;
    expect(b.showWeekends.value).toBe(true);
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
