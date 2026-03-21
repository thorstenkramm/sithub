import { isConfiguredIconName, resolveConfiguredIcon } from './icons';

describe('icons utils', () => {
  it('accepts valid mdi icon names', () => {
    expect(isConfiguredIconName('mdi-office-building')).toBe(true);
    expect(isConfiguredIconName('mdi-ev-station')).toBe(true);
  });

  it('rejects invalid icon names', () => {
    expect(isConfiguredIconName(undefined)).toBe(false);
    expect(isConfiguredIconName(null)).toBe(false);
    expect(isConfiguredIconName('office-building')).toBe(false);
    expect(isConfiguredIconName('mdi Office')).toBe(false);
    expect(isConfiguredIconName('mdi-')).toBe(false);
  });

  it('falls back when the configured icon is invalid', () => {
    expect(resolveConfiguredIcon('bad-icon', '$area')).toBe('$area');
  });

  it('returns an mdi font icon when the configured icon is valid', () => {
    expect(resolveConfiguredIcon('mdi-garage', '$area')).toBe('mdiFont:mdi-garage');
  });
});
