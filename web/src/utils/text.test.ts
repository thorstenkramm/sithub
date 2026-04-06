import { middleTruncate } from './text';

describe('middleTruncate', () => {
  it('returns original text when shorter than maxLen', () => {
    expect(middleTruncate('short', 10)).toBe('short');
  });

  it('returns original text when equal to maxLen', () => {
    expect(middleTruncate('exact', 5)).toBe('exact');
  });

  it('truncates from the middle with ellipsis', () => {
    const result = middleTruncate('Tisch 1, am Gang, rechts', 20);
    expect(result).toHaveLength(20);
    expect(result).toContain('\u2026');
    expect(result.startsWith('Tisch 1, a')).toBe(true);
    expect(result.endsWith('rechts')).toBe(true);
  });

  it('handles very short maxLen', () => {
    const result = middleTruncate('Hello World', 5);
    expect(result).toHaveLength(5);
    expect(result).toContain('\u2026');
  });
});
