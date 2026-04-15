import { useAreaViewPreference } from './useAreaViewPreference';

describe('useAreaViewPreference', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('defaults to cards when no preference stored', () => {
    const { load } = useAreaViewPreference();
    expect(load('area-1', true)).toBe('cards');
  });

  it('saves and restores table preference for an area', () => {
    const { save, load } = useAreaViewPreference();
    save('area-1', 'table');
    expect(load('area-1', true)).toBe('table');
  });

  it('scopes preference by area ID', () => {
    const { save, load } = useAreaViewPreference();
    save('area-1', 'table');
    expect(load('area-2', true)).toBe('cards');
  });

  it('returns cards on mobile even if table was saved', () => {
    const { save, load } = useAreaViewPreference();
    save('area-1', 'table');
    expect(load('area-1', false)).toBe('cards');
  });

  it('removes preference when set back to cards', () => {
    const { save, load } = useAreaViewPreference();
    save('area-1', 'table');
    save('area-1', 'cards');
    expect(load('area-1', true)).toBe('cards');
  });

  it('updates the reactive activeView ref', () => {
    const { activeView, save, load } = useAreaViewPreference();
    expect(activeView.value).toBe('cards');
    save('area-1', 'table');
    expect(activeView.value).toBe('table');
    load('area-1', true);
    expect(activeView.value).toBe('table');
  });
});
