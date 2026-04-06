import { useWarningSuppression } from './useWarningSuppression';

const STORAGE_KEY = 'sithub_warning_suppressed';

describe('useWarningSuppression', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('returns false for unsuppressed item', () => {
    const { isWarningSuppressed } = useWarningSuppression();
    expect(isWarningSuppressed('item-1', 'Some warning')).toBe(false);
  });

  it('returns true after suppressing an item', () => {
    const { isWarningSuppressed, suppressWarning } = useWarningSuppression();
    suppressWarning('item-1', 'Some warning');
    expect(isWarningSuppressed('item-1', 'Some warning')).toBe(true);
  });

  it('does not affect other items', () => {
    const { isWarningSuppressed, suppressWarning } = useWarningSuppression();
    suppressWarning('item-1', 'Some warning');
    expect(isWarningSuppressed('item-2', 'Some warning')).toBe(false);
  });

  it('persists to localStorage as JSON array', () => {
    const { suppressWarning } = useWarningSuppression();
    suppressWarning('item-1', 'Warning A');
    suppressWarning('item-2', 'Warning B');

    const stored = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]');
    expect(stored).toHaveLength(2);
  });

  it('reads from localStorage across composable instances', () => {
    const first = useWarningSuppression();
    first.suppressWarning('item-1', 'Some warning');

    const second = useWarningSuppression();
    expect(second.isWarningSuppressed('item-1', 'Some warning')).toBe(true);
  });

  it('handles corrupted localStorage gracefully', () => {
    localStorage.setItem(STORAGE_KEY, 'not-json');
    const { isWarningSuppressed } = useWarningSuppression();
    expect(isWarningSuppressed('item-1', 'Some warning')).toBe(false);
  });

  it('handles non-array JSON in localStorage gracefully', () => {
    localStorage.setItem(STORAGE_KEY, '{"key": "value"}');
    const { isWarningSuppressed } = useWarningSuppression();
    expect(isWarningSuppressed('item-1', 'Some warning')).toBe(false);
  });

  it('does not duplicate keys on repeated suppression', () => {
    const { suppressWarning } = useWarningSuppression();
    suppressWarning('item-1', 'Same warning');
    suppressWarning('item-1', 'Same warning');

    const stored = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]');
    expect(stored).toHaveLength(1);
  });

  it('resets suppression when warning text changes', () => {
    const { isWarningSuppressed, suppressWarning } = useWarningSuppression();
    suppressWarning('item-1', 'Old warning text');
    expect(isWarningSuppressed('item-1', 'Old warning text')).toBe(true);
    expect(isWarningSuppressed('item-1', 'New warning text')).toBe(false);
  });
});
