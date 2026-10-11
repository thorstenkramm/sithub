import { computed, ref, watch, type Ref } from 'vue';

export interface WeekOption {
  label: string;
  value: string;
}

/** Returns Monday of the ISO week containing the given date. */
export function getMondayOfWeek(date: Date): Date {
  const d = new Date(date);
  const day = d.getDay();
  const diff = day === 0 ? -6 : 1 - day;
  d.setDate(d.getDate() + diff);
  d.setHours(0, 0, 0, 0);
  return d;
}

/** Returns the ISO 8601 week number for a date. */
export function getISOWeekNumber(date: Date): number {
  const d = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
  d.setUTCDate(d.getUTCDate() + 4 - (d.getUTCDay() || 7));
  const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
  return Math.ceil(((d.getTime() - yearStart.getTime()) / 86400000 + 1) / 7);
}

/** Returns an ISO week string like "2026-W07". */
export function getISOWeekString(monday: Date): string {
  const weekNum = getISOWeekNumber(monday);
  const thursday = new Date(monday);
  thursday.setDate(monday.getDate() + 3);
  const year = thursday.getFullYear();
  return `${year}-W${String(weekNum).padStart(2, '0')}`;
}

/** Returns the Monday of an ISO week string like "2026-W41", or null if the string is invalid. */
export function getMondayOfISOWeek(isoWeek: string): Date | null {
  const match = isoWeek.match(/^(\d{4})-W(\d{2})$/);
  if (!match) return null;
  const year = parseInt(match[1]!, 10);
  const week = parseInt(match[2]!, 10);
  if (week < 1 || week > 53) return null;
  // Jan 4 is always in ISO week 1
  const jan4Monday = getMondayOfWeek(new Date(year, 0, 4));
  const result = new Date(jan4Monday);
  result.setDate(jan4Monday.getDate() + (week - 1) * 7);
  return getISOWeekString(result) === isoWeek ? result : null;
}

/** Keeps a week within the allowed booking horizon, using its nearest boundary. */
export function clampWeekToOptions(week: string, options: WeekOption[]): string {
  if (options.some(option => option.value === week)) return week;
  const first = options[0]?.value ?? '';
  const last = options[options.length - 1]?.value ?? first;
  return week > last ? last : first;
}

/** Translation function signature accepted by the week helpers (compatible with vue-i18n's `t`). */
export type WeekTranslator = (key: string, named?: Record<string, unknown>) => string;

/** Returns the weekday dates for a given Monday. 5 days (Mon-Fri) or 7 days (Mon-Sun). */
export function getWeekdayDates(monday: Date, includeWeekends = false): string[] {
  const count = includeWeekends ? 7 : 5;
  const dates: string[] = [];
  for (let i = 0; i < count; i++) {
    const d = new Date(monday);
    d.setDate(monday.getDate() + i);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    dates.push(`${year}-${month}-${day}`);
  }
  return dates;
}

const WEEKDAY_KEYS = ['mo', 'tu', 'we', 'th', 'fr', 'sa', 'su'] as const;
const WEEKDAY_SHORT_KEYS = ['moShort', 'tuShort', 'weShort', 'thShort', 'frShort', 'saShort', 'suShort'] as const;

/** Default English labels used when no translation function is provided. */
const WEEKDAY_LABELS_EN = ['MO', 'TU', 'WE', 'TH', 'FR', 'SA', 'SU'];
const WEEKDAY_LABELS_SHORT_EN = ['M', 'T', 'W', 'T', 'F', 'S', 'S'];

/**
 * Returns a localized weekday label for the given day index (0=Monday).
 * When a translation function is provided, uses i18n keys from the `weekdays` section.
 * Falls back to English labels when no translation function is given.
 */
export function getWeekdayLabel(index: number, short = false, t?: (key: string) => string): string {
  if (t) {
    const keys = short ? WEEKDAY_SHORT_KEYS : WEEKDAY_KEYS;
    const key = keys[index];
    return key ? t(`weekdays.${key}`) : '';
  }
  const labels = short ? WEEKDAY_LABELS_SHORT_EN : WEEKDAY_LABELS_EN;
  return labels[index] ?? '';
}

/** Maps a backend weekday abbreviation (MO, TU, ...) to a localized label. */
export function localizeWeekday(backendAbbrev: string, t: (key: string) => string): string {
  const key = backendAbbrev.toLowerCase();
  const validKeys = new Set(WEEKDAY_KEYS as readonly string[]);
  if (validKeys.has(key)) {
    return t(`weekdays.${key}`);
  }
  return backendAbbrev;
}

/** Formats a date as DD.MM. (without year). */
function formatDayMonth(date: Date): string {
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  return `${day}.${month}.`;
}

/** Formats a date as DD.MM.YYYY. */
function formatDayMonthYear(date: Date): string {
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  return `${day}.${month}.${year}`;
}

/**
 * Composable providing week selector state and helpers.
 * Generates the current week plus `maxWeeks` weeks (default 7), defaults to the current week.
 * @param showWeekends - optional reactive ref; when true, selectedWeekDates returns 7 days
 * @param maxWeeks - optional reactive ref with the number of additional bookable weeks
 * @param t - optional translator for the week token ("Week 41" / "KW 41"); English without it
 */
export function useWeekSelector(showWeekends?: Ref<boolean>, maxWeeks?: Ref<number>, t?: WeekTranslator) {
  const weekOptions = computed<WeekOption[]>(() => {
    const options: WeekOption[] = [];
    const now = new Date();
    const monday = getMondayOfWeek(now);

    // Current week + N additional weeks (default 7 additional = 8 total)
    const totalWeeks = (maxWeeks?.value ?? 7) + 1;
    for (let i = 0; i < totalWeeks; i++) {
      const weekMonday = new Date(monday);
      weekMonday.setDate(monday.getDate() + i * 7);
      const weekSunday = new Date(weekMonday);
      weekSunday.setDate(weekMonday.getDate() + 6);
      const isoWeek = getISOWeekString(weekMonday);
      const weekNum = getISOWeekNumber(weekMonday);
      const weekToken = t ? t('weekSelector.weekNumber', { week: weekNum }) : `Week ${weekNum}`;
      const label = `${formatDayMonth(weekMonday)}-${formatDayMonthYear(weekSunday)} - ${weekToken}`;
      options.push({ label, value: isoWeek });
    }
    return options;
  });

  const selectedWeek = ref(weekOptions.value[0]?.value ?? '');

  // Settings load asynchronously and can narrow the provisional booking horizon.
  watch(weekOptions, (options) => {
    selectedWeek.value = clampWeekToOptions(selectedWeek.value, options);
  }, { flush: 'sync' });

  /** Returns the Monday Date object for the currently selected week. */
  const selectedMonday = computed(() => getMondayOfISOWeek(selectedWeek.value) ?? getMondayOfWeek(new Date()));

  /** Returns weekday date strings for the selected week (5 or 7 depending on showWeekends). */
  const selectedWeekDates = computed(() =>
    getWeekdayDates(selectedMonday.value, showWeekends?.value ?? false)
  );

  return {
    weekOptions,
    selectedWeek,
    selectedMonday,
    selectedWeekDates
  };
}
