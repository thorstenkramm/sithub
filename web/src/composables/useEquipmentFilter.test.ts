import { describe, it, expect } from 'vitest';
import { parseFilter, matchesFilter, matchesParsedFilter } from './useEquipmentFilter';

describe('parseFilter', () => {
  it('returns empty array for empty string', () => {
    expect(parseFilter('')).toEqual([]);
    expect(parseFilter('   ')).toEqual([]);
  });

  it('parses simple keywords as OR group', () => {
    const groups = parseFilter('webcam monitor');
    expect(groups).toHaveLength(1);
    expect(groups[0].keywords).toEqual(['webcam', 'monitor']);
    expect(groups[0].exact).toEqual([]);
  });

  it('parses AND groups with +', () => {
    const groups = parseFilter('webcam + docking');
    expect(groups).toHaveLength(2);
    expect(groups[0].keywords).toEqual(['webcam']);
    expect(groups[1].keywords).toEqual(['docking']);
  });

  it('parses quoted exact phrases with double quotes', () => {
    const groups = parseFilter('"27 inch display"');
    expect(groups).toHaveLength(1);
    expect(groups[0].exact).toEqual(['27 inch display']);
    expect(groups[0].keywords).toEqual([]);
  });

  it('parses quoted exact phrases with single quotes', () => {
    const groups = parseFilter("'USB-C dock'");
    expect(groups).toHaveLength(1);
    expect(groups[0].exact).toEqual(['USB-C dock']);
  });

  it('parses mixed exact and keyword AND groups', () => {
    const groups = parseFilter('"27 inch display" + webcam');
    expect(groups).toHaveLength(2);
    expect(groups[0].exact).toEqual(['27 inch display']);
    expect(groups[1].keywords).toEqual(['webcam']);
  });

  it('ignores empty AND segments', () => {
    const groups = parseFilter('webcam + + monitor');
    expect(groups).toHaveLength(2);
    expect(groups[0].keywords).toEqual(['webcam']);
    expect(groups[1].keywords).toEqual(['monitor']);
  });
});

describe('matchesFilter', () => {
  const equipment = ['24 inch display', 'webcam', 'USB-C dock', 'standing desk'];

  it('returns true for empty filter', () => {
    expect(matchesFilter(equipment, '')).toBe(true);
    expect(matchesFilter(equipment, '   ')).toBe(true);
  });

  it('matches single keyword (case-insensitive)', () => {
    expect(matchesFilter(equipment, 'webcam')).toBe(true);
    expect(matchesFilter(equipment, 'WEBCAM')).toBe(true);
    expect(matchesFilter(equipment, 'Webcam')).toBe(true);
  });

  it('returns false when keyword not found', () => {
    expect(matchesFilter(equipment, 'projector')).toBe(false);
  });

  it('matches OR keywords (any keyword matches)', () => {
    expect(matchesFilter(equipment, 'projector webcam')).toBe(true);
    expect(matchesFilter(equipment, 'projector whiteboard')).toBe(false);
  });

  it('matches AND groups (all groups must match)', () => {
    expect(matchesFilter(equipment, 'webcam + dock')).toBe(true);
    expect(matchesFilter(equipment, 'webcam + projector')).toBe(false);
  });

  it('matches exact quoted phrases', () => {
    expect(matchesFilter(equipment, '"24 inch display"')).toBe(true);
    expect(matchesFilter(equipment, '"27 inch display"')).toBe(false);
  });

  it('does not treat exact matches as substrings', () => {
    expect(matchesFilter(['27 inch display stand'], '"27 inch display"')).toBe(false);
  });

  it('matches mixed exact and keyword without requiring both', () => {
    expect(matchesFilter(['webcam'], '"27 inch display" webcam')).toBe(true);
  });

  it('matches mixed exact AND keyword', () => {
    expect(matchesFilter(equipment, '"24 inch display" + webcam')).toBe(true);
    expect(matchesFilter(equipment, '"27 inch display" + webcam')).toBe(false);
  });

  it('matches partial keywords within equipment strings', () => {
    expect(matchesFilter(equipment, 'USB')).toBe(true);
    expect(matchesFilter(equipment, 'standing')).toBe(true);
    expect(matchesFilter(equipment, 'inch')).toBe(true);
  });

  it('returns true for empty equipment when filter is empty', () => {
    expect(matchesFilter([], '')).toBe(true);
  });

  it('returns false for empty equipment with active filter', () => {
    expect(matchesFilter([], 'webcam')).toBe(false);
  });
});

describe('matchesParsedFilter', () => {
  it('reuses parsed groups for matching', () => {
    const groups = parseFilter('"27 inch display" + webcam');

    expect(matchesParsedFilter(['27 inch display', 'webcam'], groups)).toBe(true);
    expect(matchesParsedFilter(['27 inch display'], groups)).toBe(false);
  });
});
