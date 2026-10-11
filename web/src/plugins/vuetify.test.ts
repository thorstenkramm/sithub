const templates = import.meta.glob<string>('../**/*.vue', { query: '?raw', import: 'default', eager: true });

describe('icon usage', () => {
  // The default icon set is mdi-svg: a bare "mdi-*" name is rendered as SVG path data and shows nothing.
  it('should use $aliases instead of mdi-* font class names in templates', () => {
    const offenders = Object.entries(templates).flatMap(([file, source]) =>
      [...source.matchAll(/["'`]mdi-[a-z0-9-]+["'`]/g)].map((match) => `${file}: ${match[0]}`)
    );

    expect(offenders).toEqual([]);
  });
});
