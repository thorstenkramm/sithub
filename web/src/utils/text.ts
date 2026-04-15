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
 * Returns a compact display name: first initial + dot + last name part.
 * Example: getShortName("Thorsten Kramm") → "T. Kramm"
 * Example: getShortName("Alexander Seidemann-Klamant") → "A. Seidemann-Klamant"
 * Falls back to the full name if it has only one part.
 */
export function getShortName(name: string | undefined, maxLen = 14): string {
  if (!name || !name.trim()) return '';
  const parts = name.trim().split(/\s+/);
  if (parts.length < 2) return middleTruncate(parts[0]!, maxLen);
  const short = `${parts[0]!.charAt(0).toUpperCase()}. ${parts.slice(1).join(' ')}`;
  return middleTruncate(short, maxLen);
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
