import { beforeEach, describe, expect, it } from 'vitest';
import { useWarningSuppression } from '../useWarningSuppression';

const STORAGE_KEY = 'sithub_warning_suppressed';

describe('useWarningSuppression', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('suppresses a warning after dismissal when the text is unchanged (FR165)', () => {
    const { isWarningSuppressed, suppressWarning } = useWarningSuppression();
    expect(isWarningSuppressed('desk-1', 'Apple only')).toBe(false);
    suppressWarning('desk-1', 'Apple only');
    expect(isWarningSuppressed('desk-1', 'Apple only')).toBe(true);
  });

  it('re-shows the warning when the text changes (FR165)', () => {
    const { isWarningSuppressed, suppressWarning } = useWarningSuppression();
    suppressWarning('desk-1', 'Apple only');
    // Same item, different warning text -> not suppressed, so the dialog reappears.
    expect(isWarningSuppressed('desk-1', 'Apple only, Thunderbolt display')).toBe(false);
  });

  it('keys the stored dismissal on itemId + a hash of the warning text, not itemId alone', () => {
    const { suppressWarning } = useWarningSuppression();
    suppressWarning('desk-1', 'Apple only');
    const stored = JSON.parse(localStorage.getItem(STORAGE_KEY) ?? '[]') as string[];
    expect(stored).toHaveLength(1);
    const key = stored[0]!;
    // Format: `${itemId}::${hash}` — item id present, plus a non-empty hash segment.
    expect(key.startsWith('desk-1::')).toBe(true);
    expect(key.split('::')[1]).toBeTruthy();
    // The key is not the bare item id (legacy per-item scheme).
    expect(key).not.toBe('desk-1');
  });

  it('does not let a legacy item-only key suppress the current warning (AC #3)', () => {
    // Simulate a pre-existing item-only dismissal from the old scheme.
    localStorage.setItem(STORAGE_KEY, JSON.stringify(['desk-1']));
    const { isWarningSuppressed } = useWarningSuppression();
    expect(isWarningSuppressed('desk-1', 'Apple only')).toBe(false);
  });

  it('persists dismissals across composable instances', () => {
    useWarningSuppression().suppressWarning('desk-2', 'No monitor');
    // A fresh instance reads the same localStorage-backed state.
    expect(useWarningSuppression().isWarningSuppressed('desk-2', 'No monitor')).toBe(true);
  });
});
