/**
 * Truncates text from the middle, preserving the start and end portions.
 * Example: middleTruncate("Tisch 1, am Gang, rechts", 20) → "Tisch 1, a…rechts"
 */
export function middleTruncate(text: string, maxLen: number): string {
  if (text.length <= maxLen) return text;
  const ellipsis = '\u2026';
  const available = maxLen - 1; // 1 char for ellipsis
  const front = Math.ceil(available / 2);
  const back = Math.floor(available / 2);
  return text.slice(0, front) + ellipsis + text.slice(text.length - back);
}

/**
 * Derives initials from a display name by taking the first letter of each
 * space-separated part and uppercasing them.
 * Example: getInitials("Alexander Seidemann-Klamant") → "AS"
 */
export function getInitials(name: string | undefined): string {
  if (!name || !name.trim()) return '';
  return name
    .trim()
    .split(/\s+/)
    .map((part) => part.charAt(0).toUpperCase())
    .join('');
}
