import { formatDateTime, formatMediumDate, formatShortDate, toLocalISODate } from './date';

describe('formatDateTime', () => {
  const iso = '2026-03-05T14:30:00Z';

  it.each([
    { name: 'empty string', value: '' },
    { name: 'unparseable value', value: 'not-a-date' }
  ])('returns null for $name', ({ value }) => {
    expect(formatDateTime(value, 'en-US')).toBeNull();
  });

  it.each(['en-US', 'de-DE', 'fr-FR'])('formats a timestamp like Intl does for %s', (locale) => {
    const expected = new Intl.DateTimeFormat(locale, { dateStyle: 'medium', timeStyle: 'short' })
      .format(new Date(iso));
    expect(formatDateTime(iso, locale)).toBe(expected);
  });

  it('includes the year of the timestamp', () => {
    expect(formatDateTime(iso, 'en-US')).toContain('2026');
  });

  it('falls back to the runtime locale when none is given', () => {
    expect(formatDateTime(iso)).toBe(
      new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(iso))
    );
  });
});

describe('formatShortDate', () => {
  it('returns an empty string for empty input', () => {
    expect(formatShortDate('')).toBe('');
  });

  it('returns unparseable input verbatim', () => {
    expect(formatShortDate('garbage')).toBe('garbage');
  });
});

describe('toLocalISODate', () => {
  it.each([
    { name: 'a regular date', date: new Date(2026, 9, 8, 14, 30), expected: '2026-10-08' },
    { name: 'a late evening time that is already tomorrow in UTC', date: new Date(2026, 0, 31, 23, 59), expected: '2026-01-31' },
    { name: 'a just-after-midnight time that is still yesterday in UTC', date: new Date(2026, 2, 1, 0, 5), expected: '2026-03-01' },
    { name: 'single digit month and day', date: new Date(2026, 0, 2), expected: '2026-01-02' }
  ])('formats $name in local time', ({ date, expected }) => {
    expect(toLocalISODate(date)).toBe(expected);
  });
});

describe('formatMediumDate', () => {
  it('formats the date with the year in the given locale', () => {
    const expected = new Intl.DateTimeFormat('de-DE', { dateStyle: 'medium' }).format(new Date(2026, 9, 8));
    expect(formatMediumDate('2026-10-08', 'de-DE')).toBe(expected);
    expect(formatMediumDate('2026-10-08', 'de-DE')).toContain('2026');
  });

  it('returns an empty string for empty input and unparseable input verbatim', () => {
    expect(formatMediumDate('')).toBe('');
    expect(formatMediumDate('garbage')).toBe('garbage');
  });
});
