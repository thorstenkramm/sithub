import { ref } from 'vue';
import { getMondayOfWeek, getISOWeekString } from './useWeekSelector';

const WEEK_KEY = 'sithub_selected_week';
const DAY_KEY = 'sithub_selected_day';

function getSessionStorage(): Storage | null {
  try {
    return window.sessionStorage;
  } catch {
    return null;
  }
}

function formatDate(date: Date): string {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, '0');
  const d = String(date.getDate()).padStart(2, '0');
  return `${y}-${m}-${d}`;
}

function currentWeek(): string {
  return getISOWeekString(getMondayOfWeek(new Date()));
}

function today(): string {
  return formatDate(new Date());
}

function isWeekInPast(week: string): boolean {
  const match = week.match(/^(\d{4})-W(\d{2})$/);
  if (!match) return true;
  const year = parseInt(match[1]!, 10);
  const weekNum = parseInt(match[2]!, 10);
  const jan4 = new Date(year, 0, 4);
  const jan4Monday = getMondayOfWeek(jan4);
  const weekMonday = new Date(jan4Monday);
  weekMonday.setDate(jan4Monday.getDate() + (weekNum - 1) * 7);
  const weekSunday = new Date(weekMonday);
  weekSunday.setDate(weekMonday.getDate() + 6);
  weekSunday.setHours(23, 59, 59, 999);
  return weekSunday < new Date();
}

function isDayInPast(day: string): boolean {
  const d = new Date(day + 'T23:59:59');
  return d < new Date();
}

const memorizedWeek = ref<string>(loadWeek());
const memorizedDay = ref<string>(loadDay());

function loadWeek(): string {
  const storage = getSessionStorage();
  if (!storage) return currentWeek();
  const stored = storage.getItem(WEEK_KEY);
  if (!stored || isWeekInPast(stored)) return currentWeek();
  return stored;
}

function loadDay(): string {
  const storage = getSessionStorage();
  if (!storage) return today();
  const stored = storage.getItem(DAY_KEY);
  if (!stored || isDayInPast(stored)) return today();
  return stored;
}

export function useDateState() {
  memorizedWeek.value = loadWeek();
  memorizedDay.value = loadDay();

  function persist() {
    const storage = getSessionStorage();
    if (!storage) return;
    storage.setItem(WEEK_KEY, memorizedWeek.value);
    storage.setItem(DAY_KEY, memorizedDay.value);
  }

  function resetDayToToday() {
    memorizedDay.value = today();
    persist();
  }

  function getWeek(): string {
    if (isWeekInPast(memorizedWeek.value)) {
      memorizedWeek.value = currentWeek();
      persist();
    }
    return memorizedWeek.value;
  }

  function setWeek(week: string) {
    memorizedWeek.value = week;
    persist();
  }

  function getDay(): string {
    if (isDayInPast(memorizedDay.value)) {
      memorizedDay.value = today();
      persist();
    }
    return memorizedDay.value;
  }

  function setDay(day: string) {
    memorizedDay.value = day;
    persist();
  }

  return {
    memorizedWeek,
    memorizedDay,
    getWeek,
    setWeek,
    getDay,
    setDay,
    resetDayToToday
  };
}
