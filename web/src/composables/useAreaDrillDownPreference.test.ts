import { useAreaDrillDownPreference } from './useAreaDrillDownPreference';

describe('useAreaDrillDownPreference', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('defaults to enabled on compact viewports when no preference stored', () => {
    const { enabled, hasUserChoice, load } = useAreaDrillDownPreference();
    load(false);
    expect(enabled.value).toBe(true);
    expect(hasUserChoice.value).toBe(false);
  });

  it('defaults to disabled on large viewports when no preference stored', () => {
    const { enabled, hasUserChoice, load } = useAreaDrillDownPreference();
    load(true);
    expect(enabled.value).toBe(false);
    expect(hasUserChoice.value).toBe(false);
  });

  it('persists user choice across loads on the same device', () => {
    const first = useAreaDrillDownPreference();
    first.set(true);

    const second = useAreaDrillDownPreference();
    second.load(true);
    expect(second.enabled.value).toBe(true);
    expect(second.hasUserChoice.value).toBe(true);
  });

  it('stored value overrides large-screen default', () => {
    const first = useAreaDrillDownPreference();
    first.set(true);

    const second = useAreaDrillDownPreference();
    second.load(true);
    expect(second.enabled.value).toBe(true);
  });

  it('stored value overrides compact-screen default', () => {
    const first = useAreaDrillDownPreference();
    first.set(false);

    const second = useAreaDrillDownPreference();
    second.load(false);
    expect(second.enabled.value).toBe(false);
  });

  it('set updates the reactive enabled ref', () => {
    const { enabled, set } = useAreaDrillDownPreference();
    set(false);
    expect(enabled.value).toBe(false);
    set(true);
    expect(enabled.value).toBe(true);
  });

  it('ignores corrupted storage values and falls back to viewport default', () => {
    localStorage.setItem('sithub_area_drill_down', 'garbage');
    const { enabled, hasUserChoice, load } = useAreaDrillDownPreference();
    load(true);
    expect(enabled.value).toBe(false);
    expect(hasUserChoice.value).toBe(false);
  });
});
