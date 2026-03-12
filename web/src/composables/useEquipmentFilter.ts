/**
 * Equipment filter composable for matching items by equipment keywords.
 *
 * Syntax:
 * - Space-separated words are OR-combined keywords
 * - `+` separates AND groups
 * - Single or double quotes for exact matching
 * - Case-insensitive
 *
 * Example: `"27 inch display" + webcam`
 *   → items must have exact "27 inch display" AND contain "webcam"
 */

export interface AndGroup {
  exact: string[];
  keywords: string[];
}

/**
 * Parse a filter string into AND groups, each containing exact phrases
 * and/or OR-combined keywords.
 */
export function parseFilter(filterText: string): AndGroup[] {
  const trimmed = filterText.trim();
  if (!trimmed) return [];

  return trimmed.split('+').map(segment => {
    const part = segment.trim();
    if (!part) return { exact: [], keywords: [] };

    const exact: string[] = [];
    // Extract quoted phrases (single or double quotes)
    const withoutQuoted = part.replace(/(['"])(.*?)\1/g, (_match, _quote, phrase) => {
      const p = phrase.trim();
      if (p) exact.push(p);
      return '';
    });

    // Remaining unquoted text: split by whitespace for OR keywords
    const keywords = withoutQuoted
      .split(/\s+/)
      .map(k => k.trim())
      .filter(k => k.length > 0);

    return { exact, keywords };
  }).filter(g => g.exact.length > 0 || g.keywords.length > 0);
}

/**
 * Check if an AND group matches an equipment list.
 * Exact phrases and keywords are OR-combined within a group.
 */
function groupMatches(group: AndGroup, equipmentLower: string[]): boolean {
  const exactMatches = group.exact.some(phrase => {
    const phraseLower = phrase.toLowerCase();
    return equipmentLower.some(e => e === phraseLower);
  });

  const keywordMatches = group.keywords.some(kw => {
    const kwLower = kw.toLowerCase();
    return equipmentLower.some(e => e.includes(kwLower));
  });

  if (group.exact.length > 0 && group.keywords.length > 0) {
    return exactMatches || keywordMatches;
  }
  if (group.exact.length > 0) return exactMatches;
  if (group.keywords.length > 0) return keywordMatches;
  return true;
}

export function matchesParsedFilter(equipment: string[], groups: AndGroup[]): boolean {
  if (groups.length === 0) return true;

  const equipmentLower = equipment.map(e => e.toLowerCase().trim());
  return groups.every(group => groupMatches(group, equipmentLower));
}

/**
 * Check if an item's equipment matches the filter string.
 * Returns true if the filter is empty (all items match).
 */
export function matchesFilter(equipment: string[], filterText: string): boolean {
  const groups = parseFilter(filterText);
  return matchesParsedFilter(equipment, groups);
}
