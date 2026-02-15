import { computed, ref, type Ref } from 'vue';

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

const WEEKDAY_LABELS = ['MO', 'TU', 'WE', 'TH', 'FR', 'SA', 'SU'];
const WEEKDAY_LABELS_SHORT = ['M', 'T', 'W', 'T', 'F', 'S', 'S'];

export function getWeekdayLabel(index: number, short = false): string {
  const labels = short ? WEEKDAY_LABELS_SHORT : WEEKDAY_LABELS;
  return labels[index] ?? '';
}

/**
 * Composable providing week selector state and helpers.
 * Generates next 8 weeks, defaults to current week.
 * @param showWeekends - optional reactive ref; when true, selectedWeekDates returns 7 days
 */
export function useWeekSelector(showWeekends?: Ref<boolean>) {
  const dateFormatter = new Intl.DateTimeFormat(undefined, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  });

  const weekOptions = computed<WeekOption[]>(() => {
    const options: WeekOption[] = [];
    const now = new Date();
    const monday = getMondayOfWeek(now);

    for (let i = 0; i < 8; i++) {
      const weekMonday = new Date(monday);
      weekMonday.setDate(monday.getDate() + i * 7);
      const isoWeek = getISOWeekString(weekMonday);
      const weekNum = getISOWeekNumber(weekMonday);
      const dateStr = dateFormatter.format(weekMonday);
      options.push({
        label: `${dateStr} - Week ${weekNum}`,
        value: isoWeek
      });
    }
    return options;
  });

  const selectedWeek = ref(weekOptions.value[0]?.value ?? '');

  /** Returns the Monday Date object for the currently selected week. */
  const selectedMonday = computed(() => {
    const match = selectedWeek.value.match(/^(\d{4})-W(\d{2})$/);
    if (!match) return getMondayOfWeek(new Date());
    const year = parseInt(match[1]!, 10);
    const week = parseInt(match[2]!, 10);
    // Jan 4 is always in ISO week 1
    const jan4 = new Date(year, 0, 4);
    const jan4Monday = getMondayOfWeek(jan4);
    const result = new Date(jan4Monday);
    result.setDate(jan4Monday.getDate() + (week - 1) * 7);
    return result;
  });

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
