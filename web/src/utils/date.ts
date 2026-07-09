/**
 * formatShortDate renders an ISO date (YYYY-MM-DD) as a locale-aware short
 * date, e.g. "Mon, Apr 8". It is the single formatter used for the area/day
 * swap guard message so the prompt reads identically across ItemsView, the
 * interactive floor plan, and the weekly matrix popover (story 36.9 P5).
 *
 * The input is parsed at local midnight to avoid a UTC off-by-one day shift.
 * Empty input returns "" and an unparseable input is returned verbatim.
 */
export function formatShortDate(dateStr: string, locale?: string): string {
  if (!dateStr) return '';
  const date = new Date(`${dateStr}T00:00:00`);
  if (Number.isNaN(date.getTime())) return dateStr;
  return new Intl.DateTimeFormat(locale || undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
  }).format(date);
}
