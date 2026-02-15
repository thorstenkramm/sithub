import { ref } from 'vue';
import {
  getMondayOfWeek,
  getISOWeekNumber,
  getISOWeekString,
  getWeekdayDates,
  getWeekdayLabel,
  useWeekSelector
} from './useWeekSelector';

describe('getMondayOfWeek', () => {
  it('returns Monday for a Wednesday', () => {
    const wed = new Date(2026, 1, 11); // Wed Feb 11 2026
    const monday = getMondayOfWeek(wed);
    expect(monday.getDay()).toBe(1);
    expect(monday.getDate()).toBe(9);
  });

  it('returns Monday for a Sunday', () => {
    const sun = new Date(2026, 1, 15); // Sun Feb 15 2026
    const monday = getMondayOfWeek(sun);
    expect(monday.getDay()).toBe(1);
    expect(monday.getDate()).toBe(9);
  });

  it('returns same day for a Monday', () => {
    const mon = new Date(2026, 1, 9); // Mon Feb 9 2026
    const monday = getMondayOfWeek(mon);
    expect(monday.getDay()).toBe(1);
    expect(monday.getDate()).toBe(9);
  });
});

describe('getISOWeekNumber', () => {
  it('returns week 1 for Jan 1 2026 (Thursday)', () => {
    expect(getISOWeekNumber(new Date(2026, 0, 1))).toBe(1);
  });

  it('returns week 7 for Feb 9 2026 (Monday)', () => {
    expect(getISOWeekNumber(new Date(2026, 1, 9))).toBe(7);
  });
});

describe('getISOWeekString', () => {
  it('formats correctly', () => {
    const monday = new Date(2026, 1, 9); // Mon Feb 9 2026
    expect(getISOWeekString(monday)).toBe('2026-W07');
  });

  it('pads single-digit weeks', () => {
    const monday = new Date(2026, 0, 5); // Mon Jan 5 2026
    expect(getISOWeekString(monday)).toBe('2026-W02');
  });
});

describe('getWeekdayDates', () => {
  it('returns 5 dates Mon-Fri by default', () => {
    const monday = new Date(2026, 1, 9);
    const dates = getWeekdayDates(monday);
    expect(dates).toHaveLength(5);
    expect(dates[0]).toBe('2026-02-09');
    expect(dates[4]).toBe('2026-02-13');
  });

  it('returns 7 dates Mon-Sun when includeWeekends is true', () => {
    const monday = new Date(2026, 1, 9);
    const dates = getWeekdayDates(monday, true);
    expect(dates).toHaveLength(7);
    expect(dates[0]).toBe('2026-02-09');
    expect(dates[4]).toBe('2026-02-13');
    expect(dates[5]).toBe('2026-02-14'); // Saturday
    expect(dates[6]).toBe('2026-02-15'); // Sunday
  });
});

describe('getWeekdayLabel', () => {
  it('returns full label by default', () => {
    expect(getWeekdayLabel(0)).toBe('MO');
    expect(getWeekdayLabel(4)).toBe('FR');
  });

  it('returns short label when requested', () => {
    expect(getWeekdayLabel(0, true)).toBe('M');
    expect(getWeekdayLabel(4, true)).toBe('F');
  });

  it('returns weekend labels', () => {
    expect(getWeekdayLabel(5)).toBe('SA');
    expect(getWeekdayLabel(6)).toBe('SU');
  });

  it('returns weekend short labels', () => {
    expect(getWeekdayLabel(5, true)).toBe('S');
    expect(getWeekdayLabel(6, true)).toBe('S');
  });

  it('returns empty string for out-of-range', () => {
    expect(getWeekdayLabel(7)).toBe('');
  });
});

describe('useWeekSelector', () => {
  it('generates 8 week options', () => {
    const { weekOptions } = useWeekSelector();
    expect(weekOptions.value).toHaveLength(8);
  });

  it('defaults selectedWeek to first option', () => {
    const { weekOptions, selectedWeek } = useWeekSelector();
    expect(selectedWeek.value).toBe(weekOptions.value[0]!.value);
  });

  it('selectedWeekDates returns 5 dates', () => {
    const { selectedWeekDates } = useWeekSelector();
    expect(selectedWeekDates.value).toHaveLength(5);
  });

  it('week options contain ISO week format', () => {
    const { weekOptions } = useWeekSelector();
    for (const option of weekOptions.value) {
      expect(option.value).toMatch(/^\d{4}-W\d{2}$/);
    }
  });

  it('week option labels contain week number', () => {
    const { weekOptions } = useWeekSelector();
    for (const option of weekOptions.value) {
      expect(option.label).toMatch(/Week \d+/);
    }
  });

  it('selectedWeekDates returns 7 dates when showWeekends is true', () => {
    const showWeekends = ref(true);
    const { selectedWeekDates } = useWeekSelector(showWeekends);
    expect(selectedWeekDates.value).toHaveLength(7);
  });

  it('selectedWeekDates returns 5 dates when showWeekends is false', () => {
    const showWeekends = ref(false);
    const { selectedWeekDates } = useWeekSelector(showWeekends);
    expect(selectedWeekDates.value).toHaveLength(5);
  });
});
