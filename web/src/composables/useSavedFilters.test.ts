import { useSavedFilters } from './useSavedFilters';

describe('useSavedFilters', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('starts with empty filters', () => {
    const { savedFilters } = useSavedFilters();
    expect(savedFilters.value).toEqual([]);
  });

  it('saves a filter to local storage', () => {
    const { saveFilter, savedFilters } = useSavedFilters();
    saveFilter('webcam');
    expect(savedFilters.value).toEqual(['webcam']);
    expect(JSON.parse(localStorage.getItem('sithub_saved_filters')!)).toEqual(['webcam']);
  });

  it('does not save duplicate filters', () => {
    const { saveFilter, savedFilters } = useSavedFilters();
    saveFilter('webcam');
    const result = saveFilter('webcam');
    expect(result).toBe(false);
    expect(savedFilters.value).toEqual(['webcam']);
  });

  it('does not save empty or whitespace-only filters', () => {
    const { saveFilter, savedFilters } = useSavedFilters();
    expect(saveFilter('')).toBe(false);
    expect(saveFilter('   ')).toBe(false);
    expect(savedFilters.value).toEqual([]);
  });

  it('deletes a saved filter', () => {
    const { saveFilter, deleteFilter, savedFilters } = useSavedFilters();
    saveFilter('webcam');
    saveFilter('monitor');
    deleteFilter('webcam');
    expect(savedFilters.value).toEqual(['monitor']);
  });

  it('isSavedFilter checks if filter exists', () => {
    const { saveFilter, isSavedFilter } = useSavedFilters();
    saveFilter('webcam');
    expect(isSavedFilter('webcam')).toBe(true);
    expect(isSavedFilter('monitor')).toBe(false);
  });

  it('comboboxItems returns saved filters', () => {
    const { saveFilter, comboboxItems } = useSavedFilters();
    saveFilter('webcam');
    saveFilter('monitor');
    expect(comboboxItems.value).toEqual(['webcam', 'monitor']);
  });

  it('loads filters from local storage on init', () => {
    localStorage.setItem('sithub_saved_filters', JSON.stringify(['webcam', 'monitor']));
    const { savedFilters } = useSavedFilters();
    expect(savedFilters.value).toEqual(['webcam', 'monitor']);
  });

  it('handles corrupted local storage gracefully', () => {
    localStorage.setItem('sithub_saved_filters', 'not-json');
    const { savedFilters } = useSavedFilters();
    expect(savedFilters.value).toEqual([]);
  });
});
