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
  const date = parseLocalISODate(dateStr);
  if (!date.ok) return date.value;
  return new Intl.DateTimeFormat(locale || undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
  }).format(date.value);
}

/**
 * formatDateTime renders an RFC 3339 timestamp (e.g. "2026-03-05T14:30:00Z") as a
 * locale-aware medium date with short time in the viewer's time zone.
 *
 * Returns null for empty or unparseable input so callers can show their own
 * placeholder (e.g. "never").
 */
export function formatDateTime(iso: string, locale?: string): string | null {
  if (!iso) return null;
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return null;
  return new Intl.DateTimeFormat(locale || undefined, {
    dateStyle: 'medium',
    timeStyle: 'short'
  }).format(date);
}

/**
 * toLocalISODate formats a Date as YYYY-MM-DD in the viewer's local time zone.
 * Unlike `toISOString().slice(0, 10)` it never shifts the day around midnight.
 */
export function toLocalISODate(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

/**
 * formatMediumDate renders an ISO date (YYYY-MM-DD) as a locale-aware medium date including the year,
 * e.g. "Oct 8, 2026". It is parsed at local midnight to avoid a UTC off-by-one day shift.
 * Empty input returns "" and an unparseable input is returned verbatim.
 */
export function formatMediumDate(dateStr: string, locale?: string): string {
  const date = parseLocalISODate(dateStr);
  if (!date.ok) return date.value;
  return new Intl.DateTimeFormat(locale || undefined, { dateStyle: 'medium' }).format(date.value);
}

type ParsedLocalDate = { ok: true; value: Date } | { ok: false; value: string };

function parseLocalISODate(dateStr: string): ParsedLocalDate {
  if (!dateStr) return { ok: false, value: '' };
  const date = new Date(`${dateStr}T00:00:00`);
  if (Number.isNaN(date.getTime())) return { ok: false, value: dateStr };
  return { ok: true, value: date };
}
