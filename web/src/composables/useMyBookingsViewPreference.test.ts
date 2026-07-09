import { useMyBookingsViewPreference } from './useMyBookingsViewPreference';

describe('useMyBookingsViewPreference', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('defaults to table on desktop when no preference stored', () => {
    const { load } = useMyBookingsViewPreference();
    expect(load(true)).toBe('table');
  });

  it('returns cards on mobile even if table is the desktop default', () => {
    const { load } = useMyBookingsViewPreference();
    expect(load(false)).toBe('cards');
  });

  it('returns cards on mobile even if table was saved', () => {
    const { save, load } = useMyBookingsViewPreference();
    save('table');
    expect(load(false)).toBe('cards');
  });

  it('saves and restores a cards preference on desktop, overriding the default', () => {
    const { save, load } = useMyBookingsViewPreference();
    save('cards');
    expect(load(true)).toBe('cards');
  });

  it('saves and restores a table preference on desktop', () => {
    const { save, load } = useMyBookingsViewPreference();
    save('table');
    expect(load(true)).toBe('table');
  });

  it('updates the reactive activeView ref', () => {
    const { activeView, save, load } = useMyBookingsViewPreference();
    expect(activeView.value).toBe('table');
    save('cards');
    expect(activeView.value).toBe('cards');
    load(true);
    expect(activeView.value).toBe('cards');
  });

  it('falls back to the desktop default when stored data is corrupted', () => {
    localStorage.setItem('sithub_my_bookings_view', '{not-a-valid-view');
    const { load } = useMyBookingsViewPreference();
    expect(load(true)).toBe('table');
  });

  it('falls back to table when localStorage is unavailable (SSR path)', () => {
    const getItem = vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('unavailable');
    });
    const { load } = useMyBookingsViewPreference();
    expect(load(true)).toBe('table');
    getItem.mockRestore();
  });
});
