import { useDateState } from './useDateState';
import { getISOWeekString, getMondayOfWeek } from './useWeekSelector';

function formatDate(date: Date): string {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, '0');
  const d = String(date.getDate()).padStart(2, '0');
  return `${y}-${m}-${d}`;
}

function currentWeek(): string {
  return getISOWeekString(getMondayOfWeek(new Date()));
}

function futureWeek(): string {
  const d = new Date();
  d.setDate(d.getDate() + 14);
  return getISOWeekString(getMondayOfWeek(d));
}

function pastWeek(): string {
  const d = new Date();
  d.setDate(d.getDate() - 14);
  return getISOWeekString(getMondayOfWeek(d));
}

function futureDay(): string {
  const d = new Date();
  d.setDate(d.getDate() + 3);
  return formatDate(d);
}

function pastDay(): string {
  const d = new Date();
  d.setDate(d.getDate() - 3);
  return formatDate(d);
}

describe('useDateState', () => {
  beforeEach(() => {
    sessionStorage.clear();
  });

  it('defaults to current week', () => {
    const { getWeek } = useDateState();
    expect(getWeek()).toBe(currentWeek());
  });

  it('defaults to today', () => {
    const { getDay } = useDateState();
    expect(getDay()).toBe(formatDate(new Date()));
  });

  it('persists week to sessionStorage', () => {
    const { setWeek, getWeek } = useDateState();
    const week = futureWeek();
    setWeek(week);
    expect(getWeek()).toBe(week);
    expect(sessionStorage.getItem('sithub_selected_week')).toBe(week);
  });

  it('persists day to sessionStorage', () => {
    const { setDay, getDay } = useDateState();
    const day = futureDay();
    setDay(day);
    expect(getDay()).toBe(day);
    expect(sessionStorage.getItem('sithub_selected_day')).toBe(day);
  });

  it('resets past week to current week', () => {
    const { setWeek, getWeek } = useDateState();
    setWeek(pastWeek());
    // getWeek detects the past week and resets
    expect(getWeek()).toBe(currentWeek());
  });

  it('resetDayToToday resets the day', () => {
    const { setDay, resetDayToToday, getDay } = useDateState();
    setDay(futureDay());
    resetDayToToday();
    expect(getDay()).toBe(formatDate(new Date()));
  });

  it('resets a past day to today', () => {
    const { setDay, getDay } = useDateState();
    setDay(pastDay());
    expect(getDay()).toBe(formatDate(new Date()));
  });
});
